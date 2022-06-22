package railsassets

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate faux --interface Executable --output fakes/executable.go

// Executable defines the interface for executing a program as a child process.
type Executable interface {
	Execute(pexec.Execution) error
}

// PrecompileProcess performs the "rake assets:precompile" build process.
type PrecompileProcess struct {
	executable Executable
	logger     scribe.Emitter
}

// NewPrecompileProcess initializes an instance of PrecompileProcess.
func NewPrecompileProcess(executable Executable, logger scribe.Emitter) PrecompileProcess {
	return PrecompileProcess{
		executable: executable,
		logger:     logger,
	}
}

// Execute runs "bundle exec rake assets:precompile assets:clean" as a child
// process. If the process fails, the error message will include the entire
// output of the child process.
func (p PrecompileProcess) Execute(workingDir string) error {
	buffer := bytes.NewBuffer(nil)

	env := []string{}
	for _, e := range []string{
		"DB_ADAPTER",
		"SECRET_KEY_BASE",
		"RAILS_ENV",
	} {
		if val, ok := os.LookupEnv(e); ok {
			env = append(env, e+"="+val)
		}
	}

	args := []string{
		"--login",
		"-c",
		strings.Join(
			[]string{
				"source",
				filepath.Join(os.ExpandEnv("$rvm_path"), "profile.d", "rvm"),
				"&&",
				strings.Join(env, " "),
				"bundle",
				"exec",
				"rake",
				"assets:precompile",
				"assets:clean",
			},
			" ",
		),
	}

	p.logger.Subprocess("Running 'bash %s'", strings.Join(args, " "))
	err := p.executable.Execute(pexec.Execution{
		Args:   args,
		Env:    append(os.Environ(), env...),
		Stdout: p.logger.ActionWriter,
		Stderr: p.logger.ActionWriter,
	})

	if err != nil {
		return fmt.Errorf("failed to execute bundle exec output:\n%s\nerror: %s", buffer.String(), err)
	}

	return nil
}
