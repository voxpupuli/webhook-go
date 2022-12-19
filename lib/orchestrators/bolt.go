package orchestrators

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// The Bolt contains the command arguments for running a Bolt Command
type Bolt struct {
	Transport    *string
	Targets      []string
	Concurrency  *int64
	HostKeyCheck *bool
}

// A BoltResult returns the JSON result from a Bolt Command Run as a struct
type BoltResult struct {
	Items []struct {
		Target  string `json:"target"`
		Action  string `json:"action"`
		Command string `json:"object"`
		Status  string `json:"status"`
		Value   struct {
			ExitCode     int64  `json:"exit_code"`
			MergedOutput string `json:"merged_output"`
			Stderr       string `json:"stderr"`
			Stdout       string `json:"stdout"`
		} `json:"value"`
	} `json:"items"`
	TargetCount int `json:"target_count"`
	ElapsedTime int `json:"elapsed_time"`
}

// boltCommand assembles the `bolt command run` command and arguments then passes it
// to the runCommand function to execute the command through os/exec as an exec.Command.
// Then it returns the pointer to a BoltResult and an error
func (b *Bolt) boltCommand(timeout time.Duration, command string) (*BoltResult, error) {
	// Create the baseline bolt command for the `command` sub-command
	cmd := []string{"bolt", "command", "run", "--targets"}

	// Convert the targets from a slice of strings to a string with
	// targets delimited by a comma
	var targets string
	for i := range b.Targets {
		targets = targets + b.Targets[i] + ","
	}
	targets = strings.TrimSuffix(targets, ",")
	cmd = append(cmd, targets)

	// If the Bolt Transport is set, then add the bolt transport option
	// to the bolt command
	if b.Transport != nil {
		transArgs := []string{"--transport", *b.Transport}
		cmd = append(cmd, transArgs...)
	}

	// If Concurrency is set, then add the bolt concurrency option to
	// the bolt command
	if b.Concurrency != nil {
		concurrency := []string{"--concurrency", strconv.FormatInt(*b.Concurrency, 10)}
		cmd = append(cmd, concurrency...)
	}

	// If the Bolt HostKeyCheck is set to false, then disable the host key check
	if *b.HostKeyCheck == false {
		cmd = append(cmd, "--no-host-key-check")
	}

	// Format the bolt output as JSON and add a connection timeout
	cmd = append(cmd, "--format", "json", "--connect-timeout", "120", command)

	// Send the command to runCommand and store the returned output and error
	//
	// If the runCommand function fails, then return an error without a result
	out, err := runCommand(strings.Join(cmd, " "), timeout)
	if err != nil {
		return nil, fmt.Errorf("Bolt: \"%s\": %s: %s", strings.Join(cmd, " "), string(out), err)
	}

	// Parse the output into a BoltResult
	result := new(BoltResult)
	if err = json.Unmarshal(out, result); err != nil {
		return nil, err
	}

	return result, nil
}

// runCommand executes any bolt sub-commands passed to it.
func runCommand(command string, timeout time.Duration) ([]byte, error) {
	var args []string

	// Currently the Windows option is not used
	if runtime.GOOS == "windows" {
		args = []string{"cmd", "/C"}
	} else {
		args = []string{"/bin/sh", "-c"}
	}
	args = append(args, command)

	// TODO: Configure context to work correctly
	// ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// defer cancel()

	// Create a new exec.Command from the passed in command and return the CombinedOutput
	cmd := exec.Command(args[0], args[1:]...)
	return cmd.CombinedOutput()
}
