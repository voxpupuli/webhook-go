package orchestrators

import (
	"fmt"

	"github.com/voxpupuli/webhook-go/lib/orchestrators/bolt"
)

type Orchestrator struct {
	Name     string
	Nodes    []string
	Protocol *string
	User     string
}

var (
	connInfo map[string]string
	admin    bool
)

func (o *Orchestrator) Deploy(cmd string) error {

	switch o.Name {
	case "bolt":
		connInfo["type"] = *o.Protocol
		connInfo["user"] = o.User
		if o.User == "root" {
			admin = true
		}
		for _, node := range o.Nodes {
			connInfo["host"] = node
			_, err := bolt.Command(connInfo, 120, admin, cmd)
			if err != nil {
				return err
			}
		}
	case "choria":
	default:
		return fmt.Errorf("Orchestration tool `%s` is not supported at this time", o.Name)
	}
	return nil
}
