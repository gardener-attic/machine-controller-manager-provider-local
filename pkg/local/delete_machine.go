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
	"time"

	apiv1alpha1 "github.com/gardener/machine-controller-manager-provider-local/pkg/api/v1alpha1"

	"github.com/gardener/machine-controller-manager/pkg/util/provider/driver"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/codes"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/status"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (d *localDriver) DeleteMachine(ctx context.Context, req *driver.DeleteMachineRequest) (*driver.DeleteMachineResponse, error) {
	if req.MachineClass.Provider != apiv1alpha1.Provider {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("requested for Provider '%s', we only support '%s'", req.MachineClass.Provider, apiv1alpha1.Provider))
	}

	klog.V(3).Infof("Machine deletion request has been received for %q", req.Machine.Name)
	defer klog.V(3).Infof("Machine deletion request has been processed for %q", req.Machine.Name)

	userDataSecret := userDataSecretForMachine(req.Machine)
	if err := d.client.Delete(ctx, userDataSecret); client.IgnoreNotFound(err) != nil {
		// Unknown leads to short retry in machine controller
		return nil, status.Error(codes.Unknown, fmt.Sprintf("error deleting user data secret: %s", err.Error()))
	}

	pod := podForMachine(req.Machine)
	if err := d.client.Delete(ctx, pod); err != nil {
		if !apierrors.IsNotFound(err) {
			// Unknown leads to short retry in machine controller
			return nil, status.Error(codes.Unknown, fmt.Sprintf("error deleting pod: %s", err.Error()))
		}
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Actively wait until pod is deleted since the extension contract in machine-controller-manager expects drivers to
	// do so. If we would not wait until the pod is gone it might happen that the kubelet could re-register the Node
	// object even after it was already deleted by machine-controller-manager.
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	if err := wait.PollUntilWithContext(timeoutCtx, 5*time.Second, func(ctx context.Context) (bool, error) {
		if err := d.client.Get(ctx, client.ObjectKeyFromObject(pod), pod); err != nil {
			if apierrors.IsNotFound(err) {
				return true, nil
			}
			// Unknown leads to short retry in machine controller
			return false, status.Error(codes.Unknown, err.Error())
		}
		return false, nil
	}); err != nil {
		// will be retried with short retry by machine controller
		return nil, status.Error(codes.DeadlineExceeded, err.Error())
	}

	return &driver.DeleteMachineResponse{}, nil
}
