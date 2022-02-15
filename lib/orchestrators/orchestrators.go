package orchestrators

import (
	"fmt"

	"github.com/voxpupuli/webhook-go/config"
)

// Deploy reads in a command string and then passes the command
// to the appropriate orchestration tool based on the application
// settings.
//
// A different type is associated with each orchestration tool and
// that type is initialized into a variable based on what is set in
// the configuration file.
//
// Deploy returns an interface of whatever custom result type is returned
// from an orchestration tool, as well as an error
func Deploy(cmd string) (interface{}, error) {
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
		res, err := boltRunner.boltCommand(20000, cmd)
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
	return nil, nil
}
