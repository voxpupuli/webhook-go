package parsers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v39/github"
	"github.com/voxpupuli/webhook-go/config"
	customerrors "github.com/voxpupuli/webhook-go/lib"
)

// providerSecret returns the webhook verification secret configured for the
// given provider (as identified by ParseHeaders). An empty secret disables
// verification for that provider, preserving the default behavior of
// accepting unsigned webhooks.
func providerSecret(provider string) string {
	conf := config.GetConfig()

	switch provider {
	case "github":
		return conf.Secrets.Github
	case "gitlab":
		return conf.Secrets.Gitlab
	case "bitbucket-cloud":
		return conf.Secrets.Bitbucket
	case "bitbucket-server":
		return conf.Secrets.BitbucketServer
	case "gitea":
		return conf.Secrets.Gitea
	case "azuredevops":
		return conf.Secrets.AzureDevOps
	default:
		return ""
	}
}

// verifyGithubSignature validates the HMAC signature GitHub delivers in the
// X-Hub-Signature-256 header, falling back to the legacy X-Hub-Signature
// (SHA-1) header when the newer one is absent. The comparison is done in
// constant time by the go-github helper.
func verifyGithubSignature(header http.Header, payload []byte, secret string) error {
	signature := header.Get(github.SHA256SignatureHeader)
	if signature == "" {
		signature = header.Get(github.SHA1SignatureHeader)
	}

	if err := github.ValidateSignature(signature, payload, []byte(secret)); err != nil {
		return customerrors.ErrWebhookVerificationFailed
	}

	return nil
}

// verifyGitlabToken compares the X-Gitlab-Token header with the configured
// secret in constant time.
func verifyGitlabToken(header http.Header, secret string) error {
	if subtle.ConstantTimeCompare([]byte(header.Get("X-Gitlab-Token")), []byte(secret)) != 1 {
		return customerrors.ErrWebhookVerificationFailed
	}

	return nil
}

// verifyBitbucketSignature validates the X-Hub-Signature header Bitbucket
// Cloud delivers ("sha256=<hex digest>") against the HMAC-SHA256 of the raw
// request body, comparing in constant time.
func verifyBitbucketSignature(header http.Header, payload []byte, secret string) error {
	const sigPrefix = "sha256="

	signature := header.Get("X-Hub-Signature")
	if !strings.HasPrefix(signature, sigPrefix) {
		return customerrors.ErrWebhookVerificationFailed
	}

	given, err := hex.DecodeString(strings.TrimPrefix(signature, sigPrefix))
	if err != nil {
		return customerrors.ErrWebhookVerificationFailed
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)

	if !hmac.Equal(given, mac.Sum(nil)) {
		return customerrors.ErrWebhookVerificationFailed
	}

	return nil
}

// verifyGiteaSignature validates the X-Gitea-Signature header, the hex
// encoded HMAC-SHA256 of the raw request body, comparing in constant time.
func verifyGiteaSignature(header http.Header, payload []byte, secret string) error {
	given, err := hex.DecodeString(header.Get("X-Gitea-Signature"))
	if err != nil {
		return customerrors.ErrWebhookVerificationFailed
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)

	if !hmac.Equal(given, mac.Sum(nil)) {
		return customerrors.ErrWebhookVerificationFailed
	}

	return nil
}

// verifyAzureDevopsAuth validates the HTTP Basic authentication Azure DevOps
// uses for webhooks: the password component of the Authorization header
// carries the shared secret, while the username is arbitrary and ignored.
// The scheme is matched case-insensitively, as RFC 7235 requires, and the
// credentials are split on the first colon, so a password may itself
// contain colons.
func verifyAzureDevopsAuth(r *http.Request, secret string) error {
	_, password, ok := r.BasicAuth()
	if !ok {
		return customerrors.ErrWebhookVerificationFailed
	}

	if subtle.ConstantTimeCompare([]byte(password), []byte(secret)) != 1 {
		return customerrors.ErrWebhookVerificationFailed
	}

	return nil
}

// readBodyAndVerify reads the full request body and, when a verification
// secret is configured for the provider, verifies the request with the given
// verification function. The raw payload is returned for further parsing
// and the request body is restored from it, so the webhook parsing libraries
// that read the body themselves can still consume it.
func readBodyAndVerify(c *gin.Context, provider string, verify func(header http.Header, payload []byte, secret string) error) ([]byte, error) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	defer c.Request.Body.Close()

	if secret := providerSecret(provider); secret != "" {
		if err := verify(c.Request.Header, payload, secret); err != nil {
			return nil, err
		}
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(payload))

	return payload, nil
}
