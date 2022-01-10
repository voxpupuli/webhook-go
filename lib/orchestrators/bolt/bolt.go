package bolt

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
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

func runCommand(command string, timeout time.Duration) ([]byte, error) {
	var args []string

	if runtime.GOOS == "windows" {
		args = []string{"cmd", "/C"}
	} else {
		args = []string{"/bin/sh", "-c"}
	}
	args = append(args, command)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	return cmd.Output()
}

func Command(connInfo map[string]string, timeout time.Duration, sudo bool, command string) (*Result, error) {
	cmdargs := []string{
		"bolt", "command", "run", "--nodes", connInfo["type"] + "://" + connInfo["host"], "-u", connInfo["user"],
	}

	if connInfo["type"] == "winrm" {
		cmdargs = append(cmdargs, "-p", "\""+connInfo["password"]+"\"", "--no-ssl")
	} else {
		if sudo {
			cmdargs = append(cmdargs, "--run-as", "root")
		}

		cmdargs = append(cmdargs, "--no-host-key-check")
	}

	cmdargs = append(cmdargs, "--format", "json", "--connect-timeout", "120", command)

	out, err := runCommand(strings.Join(cmdargs, " "), timeout)
	if err != nil {
		return nil, fmt.Errorf("Bolt: \"%s\": %s: %s", strings.Join(cmdargs, " "), out, err)
	}

	result := new(Result)
	if err = json.Unmarshal(out, result); err != nil {
		return nil, err
	}

	return result, nil
}
