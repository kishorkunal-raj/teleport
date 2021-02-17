/*
Copyright 2017-2019 Gravitational, Inc.

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

package auth

import "github.com/gravitational/teleport/api/types"

// ClusterConfiguration stores the cluster configuration in the backend. All
// the resources modified by this interface can only have a single instance
// in the backend.
type ClusterConfiguration interface {
	// SetClusterName gets services.ClusterName from the backend.
	GetClusterName(opts ...MarshalOption) (types.ClusterName, error)
	// SetClusterName sets services.ClusterName on the backend.
	SetClusterName(types.ClusterName) error

	// GetStaticTokens gets services.StaticTokens from the backend.
	GetStaticTokens() (types.StaticTokens, error)
	// SetStaticTokens sets services.StaticTokens on the backend.
	SetStaticTokens(types.StaticTokens) error
	// DeleteStaticTokens deletes static tokens resource
	DeleteStaticTokens() error

	// GetAuthPreference gets services.AuthPreference from the backend.
	GetAuthPreference() (types.AuthPreference, error)
	// SetAuthPreference sets services.AuthPreference from the backend.
	SetAuthPreference(types.AuthPreference) error

	// GetClusterConfig gets services.ClusterConfig from the backend.
	GetClusterConfig(opts ...MarshalOption) (types.ClusterConfig, error)
	// SetClusterConfig sets services.ClusterConfig on the backend.
	SetClusterConfig(types.ClusterConfig) error
}

// ServerClusterConfiguration manages cluster configuration on auth server
type ServerClusterConfiguration interface {
	ClusterConfiguration

	// UpsertClusterName upserts cluster name
	UpsertClusterName(types.ClusterName) error

	// DeleteClusterConfig deletes cluster config resource
	DeleteClusterConfig() error

	// DeleteClusterName deletes cluster name resource
	DeleteClusterName() error
}
