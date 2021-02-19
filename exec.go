package railsassets

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit"
)

// RunBashCmd executes a command in an interactive BASH shell
func RunBashCmd(command string, context packit.BuildContext) error {
	logger := NewLogEmitter(os.Stdout)

	cmd := exec.Command("bash")
	cmd.Dir = context.WorkingDir

	cmd.Args = append(
		cmd.Args,
		"--login",
		"-c",
		strings.Join(
			[]string{
				"source",
				filepath.Join(os.ExpandEnv("$rvm_path"), "profile.d", "rvm"),
				"&&",
				command,
			},
			" ",
		),
	)

	cmd.Env = os.Environ()

	logger.Process("Executing: %s", strings.Join(cmd.Args, " "))

	stdoutPipe, _ := cmd.StdoutPipe()
	var stderrBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(&stderrBuf)

	if err := cmd.Start(); err != nil {
		logger.Process("Failed to start command: %s", cmd.String())
		logger.Break()
		return err
	}

	stdoutReader := bufio.NewReader(stdoutPipe)
	stdoutLine, err := stdoutReader.ReadString('\n')
	for err == nil {
		logger.Subprocess(stdoutLine)
		stdoutLine, err = stdoutReader.ReadString('\n')
	}
	err = cmd.Wait()

	if err != nil {
		logger.Process("Command failed: %s", cmd.String())
		logger.Process("Error status code: %s", err.Error())
		if len(stderrBuf.String()) > 0 {
			logger.Process("Command output on stderr:")
			logger.Subprocess(stderrBuf.String())
		}
		return err
	}

	logger.Break()

	return nil
}
