package orchestrators

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Bolt struct {
	Transport    *string
	Targets      []string
	Concurrency  *int64
	RunAs        *string
	SudoPassword *string
	User         *string
	Password     *string
	HostKeyCheck *bool
}

func (b *Bolt) boltCommand(timeout time.Duration, command string) (*Result, error) {
	cmd := []string{"bolt", "command", "run", "--targets"}

	var targets string
	for i := range b.Targets {
		targets = targets + b.Targets[i] + ","
	}
	targets = strings.TrimSuffix(targets, ",")
	cmd = append(cmd, targets)

	if b.User != nil {
		userArgs := []string{"-u", *b.User}
		cmd = append(cmd, userArgs...)
	}

	if b.Password != nil {
		passArgs := []string{"--password", *b.Password}
		cmd = append(cmd, passArgs...)
	}

	if b.Transport != nil {
		transArgs := []string{"--transport", *b.Transport}
		cmd = append(cmd, transArgs...)
	}

	if b.Concurrency != nil {
		concurrency := []string{"--concurrency", strconv.FormatInt(*b.Concurrency, 10)}
		cmd = append(cmd, concurrency...)
	}

	if b.RunAs != nil {
		runAs := []string{"--run-as", *b.RunAs}
		cmd = append(cmd, runAs...)
	}

	if b.SudoPassword != nil {
		sudoPass := []string{"--sudo-password", *b.SudoPassword}
		cmd = append(cmd, sudoPass...)
	}

	if *b.HostKeyCheck {
		cmd = append(cmd, "--no-host-key-check")
	}

	cmd = append(cmd, "--format", "json", "--connect-timeout", "120", command)

	logrus.Infof("%v", cmd)

	out, err := runCommand(strings.Join(cmd, " "), timeout)
	if err != nil {
		logrus.Errorln(err)
		logrus.Errorln(string(out))
		return nil, fmt.Errorf("Bolt: \"%s\": %s: %s", strings.Join(cmd, " "), string(out), err)
	}

	result := new(Result)
	if err = json.Unmarshal(out, result); err != nil {
		return nil, err
	}

	return result, nil
}

func runCommand(command string, timeout time.Duration) ([]byte, error) {
	var args []string

	if runtime.GOOS == "windows" {
		args = []string{"cmd", "/C"}
	} else {
		args = []string{"/bin/sh", "-c"}
	}
	args = append(args, command)

	// ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// defer cancel()

	cmd := exec.Command(args[0], args[1:]...)
	return cmd.CombinedOutput()
}
