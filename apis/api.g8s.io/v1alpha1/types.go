/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Name string

// PasswordSpec defines the desired state of Password
type PasswordSpec struct {
	Name   `json:"id,omitempty"`
	Length uint8 `json:"length,omitempty"`
}

// PasswordStatus defines the observed state of Password
type PasswordStatus struct {
	Ready bool `json:"ready"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Password is the Schema for the passwords API
type Password struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PasswordSpec   `json:"spec,omitempty"`
	Status PasswordStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PasswordList contains a list of Password
type PasswordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Password `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Password{}, &PasswordList{})
}