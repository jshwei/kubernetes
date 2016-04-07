// +build windows

/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cadvisor

import (
	"github.com/google/cadvisor/events"
	cadvisorapi "github.com/google/cadvisor/info/v1"
	cadvisorapiv2 "github.com/google/cadvisor/info/v2"
)

// Fake cAdvisor implementation.
type FakeWin struct {
}

var _ Interface = new(FakeWin)

func New(port uint) (Interface, error) {
	return &FakeWin {}, nil
}

func (c *FakeWin) Start() error {
	return nil
}

func (c *FakeWin) ContainerInfo(name string, req *cadvisorapi.ContainerInfoRequest) (*cadvisorapi.ContainerInfo, error) {
	return new(cadvisorapi.ContainerInfo), nil
}

func (c *FakeWin) SubcontainerInfo(name string, req *cadvisorapi.ContainerInfoRequest) (map[string]*cadvisorapi.ContainerInfo, error) {
	return map[string]*cadvisorapi.ContainerInfo{}, nil
}

func (c *FakeWin) DockerContainer(name string, req *cadvisorapi.ContainerInfoRequest) (cadvisorapi.ContainerInfo, error) {
	return cadvisorapi.ContainerInfo{}, nil
}

func (c *FakeWin) MachineInfo() (*cadvisorapi.MachineInfo, error) {
	return new(cadvisorapi.MachineInfo), nil
}

func (c *FakeWin) VersionInfo() (*cadvisorapi.VersionInfo, error) {
	return new(cadvisorapi.VersionInfo), nil
}

func (c *FakeWin) DockerImagesFsInfo() (cadvisorapiv2.FsInfo, error) {
	return cadvisorapiv2.FsInfo{}, nil
}

func (c *FakeWin) RootFsInfo() (cadvisorapiv2.FsInfo, error) {
	return cadvisorapiv2.FsInfo{}, nil
}

func (c *FakeWin) WatchEvents(request *events.Request) (*events.EventChannel, error) {
	return new(events.EventChannel), nil
}
