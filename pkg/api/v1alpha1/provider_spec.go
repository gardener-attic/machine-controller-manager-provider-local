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
