// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

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

func service(machine *machinev1alpha1.Machine) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "machines",
			Namespace: machine.Namespace,
		},
	}
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

func userDataSecretForMachine(machine *machinev1alpha1.Machine) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName(machine.Name) + "-userdata",
			Namespace: machine.Namespace,
		},
	}
}

func podName(machineName string) string {
	return "machine-" + machineName
}
