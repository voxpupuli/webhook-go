package helpers

import "errors"

func GetPipelineStatus(status bool) error {
	if !status {
		return errors.New("Not allowed to deploy on a failed pipeline")
	}
	return nil
}
