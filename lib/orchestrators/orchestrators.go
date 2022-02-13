package orchestrators

import (
	"fmt"

	"github.com/voxpupuli/webhook-go/config"
)

type Result struct {
	Items []struct {
		Node   string            `json:"node"`
		Status string            `json:"status"`
		Result map[string]string `json:"result"`
	} `json:"items"`
	NodeCount   int `json:"node_count"`
	ElapsedTime int `json:"elapsed_time"`
}

func Deploy(cmd string) (*Result, error) {
	orch := config.GetConfig().Orchestration

	switch *orch.Type {
	case "bolt":
		boltRunner := Bolt{
			Transport:    orch.Bolt.Transport,
			Targets:      orch.Bolt.Targets,
			RunAs:        orch.Bolt.RunAs,
			SudoPassword: orch.Bolt.SudoPassword,
			User:         orch.User,
			Password:     orch.Password,
			HostKeyCheck: &orch.Bolt.HostKeyCheck,
			Concurrency:  orch.Bolt.Concurrency,
		}
		res, err := boltRunner.boltCommand(1000, cmd)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "choria":
	case "mcollective":
	case "":
		return nil, fmt.Errorf("orchestration tool must be specified, but was nil")
	default:
		return nil, fmt.Errorf("orchestration tool `%s` is not supported at this time", *orch.Type)
	}
	return &Result{}, nil
}
