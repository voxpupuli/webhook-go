package helpers

import (
	"testing"

	"github.com/voxpupuli/webhook-go/config"
	"gotest.tools/assert"
)

func Test_GetR10kConfig(t *testing.T) {
	h := Helper{}
	config.Init("./yaml")

	conf := h.GetR10kConfig()
	assert.Equal(t, ConfigFile, conf)
}
