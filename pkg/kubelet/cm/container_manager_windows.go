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

package cm

import (
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
	"k8s.io/kubernetes/pkg/kubelet/cadvisor"
	"k8s.io/kubernetes/pkg/util/mount"
)

const (
	// The percent of the machine memory capacity. The value is used to calculate
	// docker memory resource container's hardlimit to workaround docker memory
	// leakage issue. Please see kubernetes/issues/9881 for more detail.
	DockerMemoryLimitThresholdPercent = 70
	// The minimum memory limit allocated to docker container: 150Mi
	MinDockerMemoryLimit = 150 * 1024 * 1024
)

type containerManagerImpl struct {
	cadvisorInterface cadvisor.Interface
	mountUtil         mount.Interface
	NodeConfig
}

var _ ContainerManager = &containerManagerImpl{}

// TODO(vmarmol): Add limits to the system containers.
// Takes the absolute name of the specified containers.
// Empty container name disables use of the specified container.
func NewContainerManager(mountUtil mount.Interface, cadvisorInterface cadvisor.Interface) (ContainerManager, error) {
	return &containerManagerImpl{
		cadvisorInterface: cadvisorInterface,
		mountUtil:         mountUtil,
		NodeConfig:        NodeConfig{},
	}, nil
}

// TODO: plumb this up as a flag to Kubelet in a future PR
type KernelTunableBehavior string

const (
	KernelTunableWarn   KernelTunableBehavior = "warn"
	KernelTunableError  KernelTunableBehavior = "error"
	KernelTunableModify KernelTunableBehavior = "modify"
)

func (cm *containerManagerImpl) setupNode() error {
	if cm.DockerDaemonContainerName != "" {
		info, err := cm.cadvisorInterface.MachineInfo()
		var capacity = api.ResourceList{}
		if err != nil {
		} else {
			capacity = cadvisor.CapacityFromMachineInfo(info)
		}
		memoryLimit := (int64(capacity.Memory().Value() * DockerMemoryLimitThresholdPercent / 100))
		if memoryLimit < MinDockerMemoryLimit {
			glog.Warningf("Memory limit %d for container %s is too small, reset it to %d", memoryLimit, cm.DockerDaemonContainerName, MinDockerMemoryLimit)
			memoryLimit = MinDockerMemoryLimit
		}

		glog.V(2).Infof("Configure resource-only container %s with memory limit: %d", cm.DockerDaemonContainerName, memoryLimit)
	}

	return nil
}

func (cm *containerManagerImpl) Start(nodeConfig NodeConfig) error {
	cm.NodeConfig = nodeConfig

	// Setup the node
	if err := cm.setupNode(); err != nil {
		return err
	}

	return nil
}

func (cm *containerManagerImpl) SystemContainersLimit() api.ResourceList {
	cpuLimit := int64(0)

	return api.ResourceList{
		api.ResourceCPU: *resource.NewMilliQuantity(
			cpuLimit,
			resource.DecimalSI),
	}
}
