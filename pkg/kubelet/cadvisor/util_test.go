// +build cgo,linux

/*
Copyright 2017 The Kubernetes Authors.

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
	"fmt"
	info "github.com/google/cadvisor/info/v1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	v1helper "k8s.io/kubernetes/pkg/api/v1/helper"
	"k8s.io/kubernetes/pkg/features"
	"testing"
)

func TestCapacityFromMachineInfo(t *testing.T) {
	machineInfo := &info.MachineInfo{
		NumCores:       2,
		MemoryCapacity: 2048,
		HugePages: []info.HugePagesInfo{
			{
				PageSize: 5,
				NumPages: 5,
			},
		},
	}

	// enable the features.HugePages
	utilfeature.DefaultFeatureGate.Set(fmt.Sprintf("%s=true", features.HugePages))

	resourceList := CapacityFromMachineInfo(machineInfo)

	// assert the cpu and memory
	assert.Equal(t, int64(2000), resourceList.Cpu().MilliValue())
	assert.Equal(t, int64(2048), resourceList.Memory().Value())

	// assert the hugepage
	hugePageKey := int64(5 * 1024)
	value, found := resourceList[v1helper.HugePageResourceName(*resource.NewQuantity(hugePageKey, resource.BinarySI))]
	assert.Equal(t, true, found)
	assert.Equal(t, int64(hugePageKey*5), value.Value())
}
