package helpers

import (
	"testing"

	"github.com/voxpupuli/webhook-go/config"
	"gotest.tools/assert"
)

func Test_GetG10kDeployEnvironmentCommand(t *testing.T) {
	mCfg := "./yaml/webhook-g10k.yaml"
	config.Init(&mCfg)

	h := Helper{}

	t.Run("branch without commit", func(t *testing.T) {
		cmd := h.GetG10kDeployEnvironmentCommand("production", "production", "")
		assert.DeepEqual(t, []string{
			"g10k",
			"-config=/etc/puppetlabs/r10k/r10k.yaml",
			"-branch=production",
			"-verbose",
		}, cmd)
	})

	t.Run("branch with commit", func(t *testing.T) {
		cmd := h.GetG10kDeployEnvironmentCommand("production", "production", "abc123")
		assert.DeepEqual(t, []string{
			"g10k",
			"-config=/etc/puppetlabs/r10k/r10k.yaml",
			"-branch=production",
			"-commit=abc123",
			"-verbose",
		}, cmd)
	})

	t.Run("prefixed environment with commit", func(t *testing.T) {
		cmd := h.GetG10kDeployEnvironmentCommand("main_production", "production", "def456")
		assert.DeepEqual(t, []string{
			"g10k",
			"-config=/etc/puppetlabs/r10k/r10k.yaml",
			"-environment=main_production",
			"-commit=def456",
			"-verbose",
		}, cmd)
	})
}
