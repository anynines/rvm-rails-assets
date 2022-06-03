package fakes

import (
	"sync"

	"github.com/paketo-buildpacks/packit/v2"
)

type EnvironmentConfiguration struct {
	ConfigureCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			LaunchEnv packit.Environment
		}
		Returns struct {
			Error error
		}
		Stub func(packit.Environment) error
	}
}

func (f *EnvironmentConfiguration) Configure(param1 packit.Environment) error {
	f.ConfigureCall.Lock()
	defer f.ConfigureCall.Unlock()
	f.ConfigureCall.CallCount++
	f.ConfigureCall.Receives.LaunchEnv = param1
	if f.ConfigureCall.Stub != nil {
		return f.ConfigureCall.Stub(param1)
	}
	return f.ConfigureCall.Returns.Error
}
