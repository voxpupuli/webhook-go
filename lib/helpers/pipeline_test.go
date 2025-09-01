package helpers

import (
	"testing"

	"gotest.tools/assert"
)

func TestGetPipelineStatus(t *testing.T) {
	success := true
	failure := false

	noErr := GetPipelineStatus(success)
	assert.NilError(t, noErr)

	err := GetPipelineStatus(failure)
	assert.Error(t, err, "Not allowed to deploy on a failed pipeline")
}
