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
	"encoding/json"
	"fmt"

	machinev1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/driver"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/codes"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/status"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	apiv1alpha1 "github.com/gardener/machine-controller-manager-provider-local/pkg/api/v1alpha1"
	"github.com/gardener/machine-controller-manager-provider-local/pkg/api/validation"
)

func (d *localDriver) CreateMachine(ctx context.Context, req *driver.CreateMachineRequest) (*driver.CreateMachineResponse, error) {
	if req.MachineClass.Provider != apiv1alpha1.Provider {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("requested for Provider '%s', we only support '%s'", req.MachineClass.Provider, apiv1alpha1.Provider))
	}

	klog.V(3).Infof("Machine creation request has been received for %q", req.Machine.Name)
	defer klog.V(3).Infof("Machine creation request has been processed for %q", req.Machine.Name)

	providerSpec, err := validateProviderSpecAndSecret(req.MachineClass, req.Secret)
	if err != nil {
		return nil, err
	}

	userDataSecret := userDataSecretForMachine(req.Machine)
	userDataSecret.Data = map[string][]byte{"userdata": req.Secret.Data["userData"]}

	if err := controllerutil.SetControllerReference(req.Machine, userDataSecret, d.client.Scheme()); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("could not set userData secret ownership: %s", err.Error()))
	}

	if err := d.client.Patch(ctx, userDataSecret, client.Apply, fieldOwner, client.ForceOwnership); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error applying user data secret: %s", err.Error()))
	}

	if _, err := d.applyService(ctx, req); err != nil {
		return nil, err
	}

	pod, err := d.applyPod(ctx, req, providerSpec, userDataSecret)
	if err != nil {
		return nil, err
	}

	return &driver.CreateMachineResponse{
		ProviderID: pod.Name,
		NodeName:   pod.Name,
	}, nil
}

func (d *localDriver) applyService(ctx context.Context, req *driver.CreateMachineRequest) (*corev1.Service, error) {
	svc := service(req.Machine)
	svc.Spec.Type = corev1.ServiceTypeClusterIP
	svc.Spec.ClusterIP = corev1.ClusterIPNone
	svc.Spec.Ports = []corev1.ServicePort{{
		Port:       10250,
		Protocol:   corev1.ProtocolTCP,
		TargetPort: intstr.FromInt(10250),
	}}
	svc.Spec.Selector = map[string]string{
		labelKeyProvider: apiv1alpha1.Provider,
		labelKeyApp:      labelValueMachine,
	}

	if err := d.client.Patch(ctx, svc, client.Apply, fieldOwner, client.ForceOwnership); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error applying service: %s", err.Error()))
	}

	return svc, nil
}

func (d *localDriver) applyPod(
	ctx context.Context,
	req *driver.CreateMachineRequest,
	providerSpec *apiv1alpha1.ProviderSpec,
	userDataSecret *corev1.Secret,
) (
	*corev1.Pod,
	error,
) {
	pod := podForMachine(req.Machine)
	pod.Labels = map[string]string{
		labelKeyProvider:                   apiv1alpha1.Provider,
		labelKeyApp:                        labelValueMachine,
		"networking.gardener.cloud/to-dns": "allowed",
		"networking.gardener.cloud/to-private-networks":                 "allowed",
		"networking.gardener.cloud/to-public-networks":                  "allowed",
		"networking.gardener.cloud/to-shoot-networks":                   "allowed",
		"networking.gardener.cloud/to-runtime-apiserver":                "allowed", // needed for ManagedSeeds such that gardenlets deployed to these Machines can talk to the seed's kube-apiserver (which is the same like the garden cluster kube-apiserver)
		"networking.resources.gardener.cloud/to-kube-apiserver-tcp-443": "allowed",
	}
	pod.Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:            "node",
				Image:           providerSpec.Image,
				ImagePullPolicy: corev1.PullIfNotPresent,
				SecurityContext: &corev1.SecurityContext{
					Privileged: pointer.Bool(true),
				},
				Env: []corev1.EnvVar{{
					Name: "NODE_NAME",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				}},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "userdata",
						MountPath: "/etc/machine",
					},
					{
						Name:      "containerd",
						MountPath: "/var/lib/containerd",
					},
					{
						Name:      "modules",
						MountPath: "/lib/modules",
						ReadOnly:  true,
					},
				},
				ReadinessProbe: &corev1.Probe{
					ProbeHandler: corev1.ProbeHandler{
						Exec: &corev1.ExecAction{
							Command: []string{"sh", "-c", "/usr/bin/kubectl --kubeconfig /var/lib/kubelet/kubeconfig-real get no $NODE_NAME"},
						},
					},
				},
				Ports: []corev1.ContainerPort{{
					ContainerPort: 30123,
					Name:          "vpn-shoot",
					Protocol:      corev1.ProtocolTCP,
				}},
			},
		},
		Volumes: []corev1.Volume{
			{
				Name: "userdata",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName:  userDataSecret.Name,
						DefaultMode: pointer.Int32(0777),
					},
				},
			},
			{
				Name: "containerd",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
			{
				Name: "modules",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/lib/modules",
					},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(req.Machine, pod, d.client.Scheme()); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("could not set pod ownership: %s", err.Error()))
	}

	if err := d.client.Patch(ctx, pod, client.Apply, fieldOwner, client.ForceOwnership); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error applying pod: %s", err.Error()))
	}

	return pod, nil
}

func validateProviderSpecAndSecret(machineClass *machinev1alpha1.MachineClass, secret *corev1.Secret) (*apiv1alpha1.ProviderSpec, error) {
	if machineClass == nil {
		return nil, status.Error(codes.Internal, "MachineClass ProviderSpec is nil")
	}

	var providerSpec *apiv1alpha1.ProviderSpec
	if err := json.Unmarshal(machineClass.ProviderSpec.Raw, &providerSpec); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	validationErr := validation.ValidateProviderSpec(providerSpec, secret, field.NewPath("providerSpec"))
	if validationErr.ToAggregate() != nil && len(validationErr.ToAggregate().Errors()) > 0 {
		err := fmt.Errorf("error while validating ProviderSpec: %v", validationErr.ToAggregate().Error())
		klog.V(2).Infof("Validation of AWSMachineClass failed %s", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return providerSpec, nil
}
