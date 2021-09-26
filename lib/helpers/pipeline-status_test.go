package helpers

import (
	"errors"
	"testing"

	"github.com/voxpupuli/webhook-go/lib/parsers"
	"gotest.tools/assert"
)

func Test_CheckPipelineStatus(t *testing.T) {
	h := Helper{}
	running := parsers.Data{
		Completed: false,
		Succeed:   false,
	}

	failed := parsers.Data{
		Completed: true,
		Succeed:   false,
	}

	notCompletedErr := errors.New("received webhook but the job is not complete. Ignoring").Error()
	notSucceededErr := errors.New("received webhook but the job failed. Ignoring").Error()

	err := h.CheckPipelineStatus(running, false)
	assert.Equal(t, notCompletedErr, err.Error())

	err = h.CheckPipelineStatus(failed, false)
	assert.Equal(t, notSucceededErr, err.Error())

	err = h.CheckPipelineStatus(failed, true)
	assert.NilError(t, err)
}
