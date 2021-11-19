// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package local

import (
	"context"
	"fmt"

	apiv1alpha1 "github.com/gardener/machine-controller-manager-provider-local/pkg/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gardener/machine-controller-manager/pkg/util/provider/driver"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/codes"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/status"
	"k8s.io/klog"
)

func (d *localDriver) ListMachines(ctx context.Context, req *driver.ListMachinesRequest) (*driver.ListMachinesResponse, error) {
	if req.MachineClass.Provider != apiv1alpha1.Provider {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("requested for Provider '%s', we only support '%s'", req.MachineClass.Provider, apiv1alpha1.Provider))
	}

	klog.V(3).Infof("Machine list request has been received for %q", req.MachineClass.Name)
	defer klog.V(3).Infof("Machine list request has been processed for %q", req.MachineClass.Name)

	podList := &corev1.PodList{}
	if err := d.client.List(ctx, podList, client.InNamespace(req.MachineClass.Namespace), client.MatchingLabels{labelKeyProvider: apiv1alpha1.Provider}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	machineList := make(map[string]string, len(podList.Items))
	for _, pod := range podList.Items {
		machineList[pod.Name] = pod.Name
	}

	return &driver.ListMachinesResponse{MachineList: machineList}, nil
}
