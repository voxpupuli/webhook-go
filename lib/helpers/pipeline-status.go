package helpers

import (
	"errors"

	"github.com/voxpupuli/webhook-go/lib/parsers"
)

func (h *Helper) CheckPipelineStatus(data parsers.Data, deployOnErr bool) error {
	if !data.Completed {
		return errors.New("received webhook but the job is not complete. Ignoring")
	} else if !data.Succeed && !deployOnErr {
		return errors.New("received webhook but the job failed. Ignoring")
	}
	return nil
}
