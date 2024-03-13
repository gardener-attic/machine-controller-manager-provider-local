// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

const (
	// V1Alpha1 is the API version.
	V1Alpha1 = "mcm.gardener.cloud/v1alpha1"
	// Provider is a constant for the provider name.
	Provider = "local"
)

// ProviderSpec is the spec to be used while parsing the calls.
type ProviderSpec struct {
	// APIVersion determines the API version for the provider APIs.
	APIVersion string `json:"apiVersion,omitempty"`
	// Image is the container image to use for the node.
	Image string `json:"image,omitempty"`
}
