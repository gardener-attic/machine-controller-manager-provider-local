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

	machinev1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/driver"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	fieldOwner        = client.FieldOwner("machine-controller-manager-provider-local")
	labelKeyApp       = "app"
	labelKeyProvider  = "machine-provider"
	labelValueMachine = "machine"
)

// NewDriver returns an empty AWSDriver object
func NewDriver(client client.Client) driver.Driver {
	return &localDriver{client}
}

type localDriver struct {
	client client.Client
}

// GenerateMachineClassForMigration is not implemented.
func (d *localDriver) GenerateMachineClassForMigration(_ context.Context, _ *driver.GenerateMachineClassForMigrationRequest) (*driver.GenerateMachineClassForMigrationResponse, error) {
	return &driver.GenerateMachineClassForMigrationResponse{}, nil
}

func podForMachine(machine *machinev1alpha1.Machine) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName(machine.Name),
			Namespace: machine.Namespace,
		},
	}
}

func podName(machineName string) string {
	return "machine-" + machineName
}

func userDataSecretName(machineName string) string {
	return podName(machineName) + "-userdata"
}
