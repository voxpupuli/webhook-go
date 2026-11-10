package parsers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/voxpupuli/webhook-go/config"
	customerrors "github.com/voxpupuli/webhook-go/lib"
)

// verificationTestSecrets holds one secret per provider, keyed by the
// matching `secrets` configuration key.
var verificationTestSecrets = map[string]string{
	"github":           "github-test-secret",
	"gitlab":           "gitlab-test-secret",
	"bitbucket":        "bitbucket-test-secret",
	"bitbucket_server": "bitbucket-server-test-secret",
	"gitea":            "gitea-test-secret",
	"azure_devops":     "azure-devops-test-secret",
}

// verificationConfigKeys keeps a stable order when writing test configs.
var verificationConfigKeys = []string{
	"github",
	"gitlab",
	"bitbucket",
	"bitbucket_server",
	"gitea",
	"azure_devops",
}

// verificationProvider describes how one provider authenticates its
// webhooks, so the same verification matrix can be exercised for each.
type verificationProvider struct {
	// configKey is the key of this provider in the `secrets` block
	configKey string
	// payloadFile is the webhook payload fixture used for all requests
	payloadFile string
	// headers identify the provider, as ParseHeaders expects them
	headers []Header
	// authHeaders returns the authentication headers the real provider
	// would send for the given request body
	authHeaders func(body []byte) []Header
	// invalidAuthHeaders returns authentication headers that fail
	// verification for the given request body
	invalidAuthHeaders func(body []byte) []Header
	// expected is the Data a successful parse produces
	expected Data
	// tokenAuth marks providers whose authentication does not cover the
	// payload (GitLab's token, Azure DevOps' basic auth), for which the
	// tampered-payload case does not apply
	tokenAuth bool
}

// Test_ParseData_Verification exercises, for every provider, that:
//   - unsigned requests are still accepted when no secret is configured,
//     including when secrets are configured for other providers only
//     (backward compatibility)
//   - valid signatures or tokens are accepted once a secret is configured
//   - invalid signatures, tampered payloads and missing authentication
//     headers are rejected with the auth error once a secret is configured
func Test_ParseData_Verification(t *testing.T) {
	// Restore a configuration without verification secrets when done, so
	// test ordering can never affect other tests in this package.
	defer loadSecretsConfig(t, nil)

	t.Run("GitHub", func(t *testing.T) {
		githubHeaders := []Header{{Name: "X-Github-Event", Value: "push"}}
		expected := Data{
			Branch:     "simple-tag",
			Deleted:    true,
			ModuleName: "Hello-World",
			RepoName:   "Codertocat/Hello-World",
			RepoUser:   "Codertocat",
			Completed:  true,
			Succeed:    true,
			Provider:   "github",
		}

		testProviderVerification(t, verificationProvider{
			configKey:   "github",
			payloadFile: "./json/github/push.json",
			headers:     githubHeaders,
			authHeaders: func(body []byte) []Header {
				return []Header{{
					Name:  "X-Hub-Signature-256",
					Value: "sha256=" + hmacHexDigest(sha256.New, verificationTestSecrets["github"], body),
				}}
			},
			invalidAuthHeaders: func(body []byte) []Header {
				return []Header{{
					Name:  "X-Hub-Signature-256",
					Value: "sha256=" + hmacHexDigest(sha256.New, "not-the-configured-secret", body),
				}}
			},
			expected: expected,
		})

		t.Run("parsed with a legacy sha1 signature", func(t *testing.T) {
			loadSecretsConfig(t, map[string]string{"github": verificationTestSecrets["github"]})

			payload := readPayload(t, "./json/github/push.json")
			headers := appendHeaders(githubHeaders, []Header{{
				Name:  "X-Hub-Signature",
				Value: "sha1=" + hmacHexDigest(sha1.New, verificationTestSecrets["github"], payload),
			}})

			d := Data{}
			c := getGinContextFromBody(payload, headers)

			require.NoError(t, d.ParseData(c))
			require.Equal(t, expected, d)
		})
	})

	t.Run("GitLab", func(t *testing.T) {
		testProviderVerification(t, verificationProvider{
			configKey:   "gitlab",
			payloadFile: "./json/gitlab/push.json",
			headers: []Header{
				{Name: "X-Gitlab-Event", Value: "Push Hook"},
			},
			authHeaders: func(_ []byte) []Header {
				return []Header{{Name: "X-Gitlab-Token", Value: verificationTestSecrets["gitlab"]}}
			},
			invalidAuthHeaders: func(_ []byte) []Header {
				return []Header{{Name: "X-Gitlab-Token", Value: "not-the-configured-secret"}}
			},
			expected: Data{
				Branch:     "master",
				Deleted:    false,
				ModuleName: "Diaspora",
				RepoName:   "mike/diaspora",
				RepoUser:   "Mike",
				Completed:  true,
				Succeed:    true,
				Provider:   "gitlab",
			},
			tokenAuth: true,
		})
	})

	t.Run("Bitbucket Cloud", func(t *testing.T) {
		testProviderVerification(t, verificationProvider{
			configKey:   "bitbucket",
			payloadFile: "./json/bitbucket-cloud.json",
			headers: []Header{
				{Name: "X-Event-Key", Value: "repo:push"},
				{Name: "X-Hook-UUID", Value: "aba83f33-f838-4727-aac7-0fc45fac66f7"},
			},
			authHeaders: func(body []byte) []Header {
				return []Header{{
					Name:  "X-Hub-Signature",
					Value: "sha256=" + hmacHexDigest(sha256.New, verificationTestSecrets["bitbucket"], body),
				}}
			},
			invalidAuthHeaders: func(body []byte) []Header {
				return []Header{{
					Name:  "X-Hub-Signature",
					Value: "sha256=" + hmacHexDigest(sha256.New, "not-the-configured-secret", body),
				}}
			},
			expected: Data{
				Branch:     "master",
				Deleted:    false,
				ModuleName: "hello_app",
				RepoName:   "dhollinger/hello_app",
				RepoUser:   "dhollinger",
				Completed:  true,
				Succeed:    true,
				Provider:   "bitbucket-cloud",
			},
		})
	})

	t.Run("Bitbucket Server", func(t *testing.T) {
		bitbucketServerHeaders := []Header{
			{Name: "X-Event-Key", Value: "repo:refs_changed"},
			{Name: "X-Request-Id", Value: "abcde12345"},
		}
		expected := Data{
			Branch:     "master",
			Deleted:    false,
			ModuleName: "repository",
			RepoName:   "project/repository",
			RepoUser:   "project",
			Completed:  true,
			Succeed:    true,
			Provider:   "bitbucket-server",
		}

		testProviderVerification(t, verificationProvider{
			configKey:   "bitbucket_server",
			payloadFile: "./json/bitbucket-server.json",
			headers:     bitbucketServerHeaders,
			authHeaders: func(body []byte) []Header {
				return []Header{{
					Name:  "X-Hub-Signature",
					Value: "sha256=" + hmacHexDigest(sha256.New, verificationTestSecrets["bitbucket_server"], body),
				}}
			},
			invalidAuthHeaders: func(body []byte) []Header {
				return []Header{{
					Name:  "X-Hub-Signature",
					Value: "sha256=" + hmacHexDigest(sha256.New, "not-the-configured-secret", body),
				}}
			},
			expected: expected,
		})

		t.Run("rejected with a malformed signature", func(t *testing.T) {
			loadSecretsConfig(t, map[string]string{"bitbucket_server": verificationTestSecrets["bitbucket_server"]})

			payload := readPayload(t, "./json/bitbucket-server.json")
			headers := appendHeaders(bitbucketServerHeaders, []Header{
				{Name: "X-Hub-Signature", Value: "abc"},
			})

			d := Data{}
			c := getGinContextFromBody(payload, headers)

			require.ErrorIs(t, d.ParseData(c), customerrors.ErrWebhookVerificationFailed)
		})
	})

	t.Run("Gitea", func(t *testing.T) {
		testProviderVerification(t, verificationProvider{
			configKey:   "gitea",
			payloadFile: "./json/gitea/push.json",
			headers: []Header{
				{Name: "X-Gitea-Event", Value: "push"},
			},
			authHeaders: func(body []byte) []Header {
				return []Header{{
					Name:  "X-Gitea-Signature",
					Value: hmacHexDigest(sha256.New, verificationTestSecrets["gitea"], body),
				}}
			},
			invalidAuthHeaders: func(body []byte) []Header {
				return []Header{{
					Name:  "X-Gitea-Signature",
					Value: hmacHexDigest(sha256.New, "not-the-configured-secret", body),
				}}
			},
			expected: Data{
				Branch:     "simple-tag",
				Deleted:    false,
				ModuleName: "Hello-World",
				RepoName:   "Codertocat/Hello-World",
				RepoUser:   "Codertocat",
				Completed:  true,
				Succeed:    true,
				Provider:   "gitea",
			},
		})
	})

	t.Run("Azure DevOps", func(t *testing.T) {
		azureHeaders := []Header{{Name: "X-Azure-DevOps", Value: "git.push"}}
		expected := Data{
			Branch:     "master",
			Deleted:    false,
			ModuleName: "Fabrikam-Fiber-Git",
			RepoName:   "Fabrikam-Fiber-Git",
			RepoUser:   "278d5cd2-584d-4b63-824a-2ba458937249",
			Completed:  true,
			Succeed:    true,
			Provider:   "azuredevops",
		}

		basicAuth := func(scheme string, password string) []Header {
			credentials := base64.StdEncoding.EncodeToString([]byte("arbitrary-username:" + password))
			return []Header{{Name: "Authorization", Value: scheme + " " + credentials}}
		}

		testProviderVerification(t, verificationProvider{
			configKey:   "azure_devops",
			payloadFile: "./json/azure_devops.json",
			headers:     azureHeaders,
			authHeaders: func(_ []byte) []Header {
				return basicAuth("Basic", verificationTestSecrets["azure_devops"])
			},
			invalidAuthHeaders: func(_ []byte) []Header {
				return basicAuth("Basic", "not-the-configured-secret")
			},
			expected:  expected,
			tokenAuth: true,
		})

		t.Run("parsed with a case-insensitive basic scheme", func(t *testing.T) {
			loadSecretsConfig(t, map[string]string{"azure_devops": verificationTestSecrets["azure_devops"]})

			payload := readPayload(t, "./json/azure_devops.json")

			for _, scheme := range []string{"basic", "BaSiC"} {
				d := Data{}
				c := getGinContextFromBody(payload, appendHeaders(
					azureHeaders,
					basicAuth(scheme, verificationTestSecrets["azure_devops"]),
				))

				require.NoError(t, d.ParseData(c), "authentication scheme %q", scheme)
				require.Equal(t, expected, d, "authentication scheme %q", scheme)
			}
		})

		t.Run("parsed with a password containing a colon", func(t *testing.T) {
			loadSecretsConfig(t, map[string]string{"azure_devops": "azure-devops:test-secret"})

			payload := readPayload(t, "./json/azure_devops.json")

			d := Data{}
			c := getGinContextFromBody(payload, appendHeaders(
				azureHeaders,
				basicAuth("Basic", "azure-devops:test-secret"),
			))

			require.NoError(t, d.ParseData(c))
			require.Equal(t, expected, d)
		})
	})
}

// testProviderVerification runs the common verification matrix for a single
// provider, reloading the configuration for every case.
func testProviderVerification(t *testing.T, p verificationProvider) {
	t.Helper()

	payload := readPayload(t, p.payloadFile)
	// The same payload with one extra byte: a signature valid for the
	// original body no longer matches the tampered one.
	tampered := append(bytes.Clone(payload), ' ')

	// A configuration with secrets for every provider except this one.
	othersSecrets := make(map[string]string, len(verificationTestSecrets)-1)
	for key, secret := range verificationTestSecrets {
		if key != p.configKey {
			othersSecrets[key] = secret
		}
	}

	ownSecrets := map[string]string{p.configKey: verificationTestSecrets[p.configKey]}

	t.Run("parsed unsigned when no secret is configured", func(t *testing.T) {
		loadSecretsConfig(t, nil)

		d := Data{}
		c := getGinContextFromBody(payload, p.headers)

		require.NoError(t, d.ParseData(c))
		require.Equal(t, p.expected, d)
	})

	t.Run("parsed unsigned when only other providers have secrets", func(t *testing.T) {
		loadSecretsConfig(t, othersSecrets)

		d := Data{}
		c := getGinContextFromBody(payload, p.headers)

		require.NoError(t, d.ParseData(c))
		require.Equal(t, p.expected, d)
	})

	t.Run("parsed with a valid signature or token", func(t *testing.T) {
		loadSecretsConfig(t, ownSecrets)

		d := Data{}
		c := getGinContextFromBody(payload, appendHeaders(p.headers, p.authHeaders(payload)))

		require.NoError(t, d.ParseData(c))
		require.Equal(t, p.expected, d)
	})

	t.Run("rejected with an invalid signature or token", func(t *testing.T) {
		loadSecretsConfig(t, ownSecrets)

		d := Data{}
		c := getGinContextFromBody(payload, appendHeaders(p.headers, p.invalidAuthHeaders(payload)))

		require.ErrorIs(t, d.ParseData(c), customerrors.ErrWebhookVerificationFailed)
	})

	if !p.tokenAuth {
		t.Run("rejected with a tampered payload", func(t *testing.T) {
			loadSecretsConfig(t, ownSecrets)

			d := Data{}
			// The signature is valid for the original payload, but the
			// request carries the tampered one.
			c := getGinContextFromBody(tampered, appendHeaders(p.headers, p.authHeaders(payload)))

			require.ErrorIs(t, d.ParseData(c), customerrors.ErrWebhookVerificationFailed)
		})
	}

	t.Run("rejected without authentication headers", func(t *testing.T) {
		loadSecretsConfig(t, ownSecrets)

		d := Data{}
		c := getGinContextFromBody(payload, p.headers)

		require.ErrorIs(t, d.ParseData(c), customerrors.ErrWebhookVerificationFailed)
	})
}

// loadSecretsConfig writes the given `secrets` block to a temporary config
// file and loads it, so the parsers read it through the normal config path.
func loadSecretsConfig(t *testing.T, secrets map[string]string) {
	t.Helper()

	var configYAML strings.Builder
	configYAML.WriteString("secrets:\n")
	for _, key := range verificationConfigKeys {
		configYAML.WriteString("  " + key + ": \"" + secrets[key] + "\"\n")
	}

	path := filepath.Join(t.TempDir(), "webhook.yml")
	if err := os.WriteFile(path, []byte(configYAML.String()), 0o600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	config.Init(&path)
}

// hmacHexDigest computes the hex encoded HMAC of body with secret, the way
// the providers sign their webhook payloads.
func hmacHexDigest(hashFunc func() hash.Hash, secret string, body []byte) string {
	mac := hmac.New(hashFunc, []byte(secret))
	_, _ = mac.Write(body)

	return hex.EncodeToString(mac.Sum(nil))
}

// appendHeaders returns a new header list with extra appended to base.
func appendHeaders(base []Header, extra []Header) []Header {
	headers := append([]Header{}, base...)

	return append(headers, extra...)
}

// readPayload reads a webhook payload fixture from disk.
func readPayload(t *testing.T, filename string) []byte {
	t.Helper()

	payload, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read payload fixture %s: %v", filename, err)
	}

	return payload
}

// getGinContextFromBody builds a gin.Context carrying a POST request with
// the given body and headers, like getGinContext does for a payload file.
func getGinContextFromBody(body []byte, headers []Header) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
		Body:   io.NopCloser(bytes.NewReader(body)),
		Method: "POST",
	}

	for _, header := range headers {
		req.Header.Add(header.Name, header.Value)
	}

	c.Request = req

	return c
}
