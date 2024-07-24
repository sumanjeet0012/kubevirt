/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2024 Red Hat, Inc.
 *
 */

package multus

import (
	"encoding/json"
	"fmt"

	networkv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"

	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/kubevirt/pkg/network/vmispec"
)

type NetworkAnnotationPool struct {
	pool []networkv1.NetworkSelectionElement
}

func (nap *NetworkAnnotationPool) Add(multusNetworkAnnotation networkv1.NetworkSelectionElement) {
	nap.pool = append(nap.pool, multusNetworkAnnotation)
}

func (nap *NetworkAnnotationPool) IsEmpty() bool {
	return len(nap.pool) == 0
}

func (nap *NetworkAnnotationPool) ToString() (string, error) {
	multusNetworksAnnotation, err := json.Marshal(nap.pool)
	if err != nil {
		return "", fmt.Errorf("failed to create JSON list from multus interface pool %v", nap.pool)
	}
	return string(multusNetworksAnnotation), nil
}

func NewAnnotationData(namespace string, interfaces []v1.Interface, network v1.Network, podInterfaceName string) networkv1.NetworkSelectionElement {
	multusIface := vmispec.LookupInterfaceByName(interfaces, network.Name)
	namespace, networkName := GetNamespaceAndNetworkName(namespace, network.Multus.NetworkName)
	var multusIfaceMac string
	if multusIface != nil {
		multusIfaceMac = multusIface.MacAddress
	}
	return networkv1.NetworkSelectionElement{
		InterfaceRequest: podInterfaceName,
		MacRequest:       multusIfaceMac,
		Namespace:        namespace,
		Name:             networkName,
	}
}
