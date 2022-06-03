package fakes

import "sync"

type EnvironmentSetup struct {
	LinkCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			LayerPath  string
			WorkingDir string
		}
		Returns struct {
			Error error
		}
		Stub func(string, string) error
	}
	ResetLayerCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			LayerPath string
		}
		Returns struct {
			Error error
		}
		Stub func(string) error
	}
	ResetLocalCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			WorkingDir string
		}
		Returns struct {
			Error error
		}
		Stub func(string) error
	}
}

func (f *EnvironmentSetup) Link(param1 string, param2 string) error {
	f.LinkCall.mutex.Lock()
	defer f.LinkCall.mutex.Unlock()
	f.LinkCall.CallCount++
	f.LinkCall.Receives.LayerPath = param1
	f.LinkCall.Receives.WorkingDir = param2
	if f.LinkCall.Stub != nil {
		return f.LinkCall.Stub(param1, param2)
	}
	return f.LinkCall.Returns.Error
}
func (f *EnvironmentSetup) ResetLayer(param1 string) error {
	f.ResetLayerCall.mutex.Lock()
	defer f.ResetLayerCall.mutex.Unlock()
	f.ResetLayerCall.CallCount++
	f.ResetLayerCall.Receives.LayerPath = param1
	if f.ResetLayerCall.Stub != nil {
		return f.ResetLayerCall.Stub(param1)
	}
	return f.ResetLayerCall.Returns.Error
}
func (f *EnvironmentSetup) ResetLocal(param1 string) error {
	f.ResetLocalCall.mutex.Lock()
	defer f.ResetLocalCall.mutex.Unlock()
	f.ResetLocalCall.CallCount++
	f.ResetLocalCall.Receives.WorkingDir = param1
	if f.ResetLocalCall.Stub != nil {
		return f.ResetLocalCall.Stub(param1)
	}
	return f.ResetLocalCall.Returns.Error
}
