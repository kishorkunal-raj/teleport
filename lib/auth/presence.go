/*
Copyright 2021 Gravitational, Inc.

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

import (
	"context"

	"github.com/gravitational/teleport/api/types"
)

// Presence records and reports the presence of all components
// of the cluster - Nodes, Proxies and SSH nodes
type Presence interface {
	// Semaphores is responsible for semaphore handling
	types.Semaphores

	// GetLocalClusterName upserts local domain
	GetLocalClusterName() (string, error)

	// GetNodes returns a list of registered servers. Schema validation can be
	// skipped to improve performance.
	GetNodes(namespace string, opts ...MarshalOption) ([]types.Server, error)

	// DeleteAllNodes deletes all nodes in a namespace.
	DeleteAllNodes(namespace string) error

	// DeleteNode deletes node in a namespace
	DeleteNode(namespace, name string) error

	// UpsertNode registers node presence, permanently if TTL is 0 or for the
	// specified duration with second resolution if it's >= 1 second.
	UpsertNode(server types.Server) (*types.KeepAlive, error)

	// UpsertNodes bulk inserts nodes.
	UpsertNodes(namespace string, servers []types.Server) error

	// DELETE IN: 5.1.0
	//
	// This logic has been moved to KeepAliveServer.
	//
	// KeepAliveNode updates node TTL in the storage
	KeepAliveNode(ctx context.Context, h types.KeepAlive) error

	// GetAuthServers returns a list of registered servers
	GetAuthServers() ([]types.Server, error)

	// UpsertAuthServer registers auth server presence, permanently if ttl is 0 or
	// for the specified duration with second resolution if it's >= 1 second
	UpsertAuthServer(server types.Server) error

	// UpsertProxy registers proxy server presence, permanently if ttl is 0 or
	// for the specified duration with second resolution if it's >= 1 second
	UpsertProxy(server types.Server) error

	// ProxyGetter gets a list of proxies
	types.ProxyGetter

	// DeleteProxy deletes proxy by name
	DeleteProxy(name string) error

	// DeleteAllProxies deletes all proxies
	DeleteAllProxies() error

	// UpsertReverseTunnel upserts reverse tunnel entry temporarily or permanently
	UpsertReverseTunnel(tunnel types.ReverseTunnel) error

	// GetReverseTunnels returns a list of registered servers
	GetReverseTunnels(opts ...MarshalOption) ([]types.ReverseTunnel, error)

	// DeleteReverseTunnel deletes reverse tunnel by it's domain name
	DeleteReverseTunnel(domainName string) error

	// GetNamespaces returns a list of namespaces
	GetNamespaces() ([]types.Namespace, error)

	// GetNamespace returns namespace by name
	GetNamespace(name string) (*types.Namespace, error)

	// UpsertNamespace upserts namespace
	UpsertNamespace(types.Namespace) error

	// DeleteNamespace deletes namespace by name
	DeleteNamespace(name string) error

	// UpsertTrustedCluster creates or updates a TrustedCluster in the backend.
	UpsertTrustedCluster(ctx context.Context, tc types.TrustedCluster) (types.TrustedCluster, error)

	// GetTrustedCluster returns a single TrustedCluster by name.
	GetTrustedCluster(string) (types.TrustedCluster, error)

	// GetTrustedClusters returns all TrustedClusters in the backend.
	GetTrustedClusters() ([]types.TrustedCluster, error)

	// DeleteTrustedCluster removes a TrustedCluster from the backend by name.
	DeleteTrustedCluster(ctx context.Context, name string) error

	// UpsertTunnelConnection upserts tunnel connection
	UpsertTunnelConnection(types.TunnelConnection) error

	// GetTunnelConnections returns tunnel connections for a given cluster
	GetTunnelConnections(clusterName string, opts ...MarshalOption) ([]types.TunnelConnection, error)

	// GetAllTunnelConnections returns all tunnel connections
	GetAllTunnelConnections(opts ...MarshalOption) ([]types.TunnelConnection, error)

	// DeleteTunnelConnection deletes tunnel connection by name
	DeleteTunnelConnection(clusterName string, connName string) error

	// DeleteTunnelConnections deletes all tunnel connections for cluster
	DeleteTunnelConnections(clusterName string) error

	// DeleteAllTunnelConnections deletes all tunnel connections for cluster
	DeleteAllTunnelConnections() error

	// CreateRemoteCluster creates a remote cluster
	CreateRemoteCluster(types.RemoteCluster) error

	// UpdateRemoteCluster updates a remote cluster
	UpdateRemoteCluster(ctx context.Context, rc types.RemoteCluster) error

	// GetRemoteClusters returns a list of remote clusters
	GetRemoteClusters(opts ...MarshalOption) ([]types.RemoteCluster, error)

	// GetRemoteCluster returns a remote cluster by name
	GetRemoteCluster(clusterName string) (types.RemoteCluster, error)

	// DeleteRemoteCluster deletes remote cluster by name
	DeleteRemoteCluster(clusterName string) error

	// DeleteAllRemoteClusters deletes all remote clusters
	DeleteAllRemoteClusters() error

	// UpsertKubeService registers kubernetes service presence.
	UpsertKubeService(context.Context, types.Server) error

	// GetAppServers gets all application servers.
	GetAppServers(context.Context, string, ...MarshalOption) ([]types.Server, error)

	// UpsertAppServer adds an application server.
	UpsertAppServer(context.Context, types.Server) (*types.KeepAlive, error)

	// DeleteAppServer removes an application server.
	DeleteAppServer(context.Context, string, string) error

	// DeleteAllAppServers removes all application servers.
	DeleteAllAppServers(context.Context, string) error

	// GetDatabaseServers returns all registered database proxy servers.
	GetDatabaseServers(context.Context, string, ...MarshalOption) ([]types.DatabaseServer, error)
	// UpsertDatabaseServer creates or updates a new database proxy server.
	UpsertDatabaseServer(context.Context, types.DatabaseServer) (*types.KeepAlive, error)
	// DeleteDatabaseServer removes the specified database proxy server.
	DeleteDatabaseServer(context.Context, string, string, string) error
	// DeleteAllDatabaseServers removes all database proxy servers.
	DeleteAllDatabaseServers(context.Context, string) error

	// KeepAliveServer updates TTL of the server resource in the backend.
	KeepAliveServer(ctx context.Context, h types.KeepAlive) error

	// GetKubeServices returns a list of registered kubernetes services.
	GetKubeServices(context.Context) ([]types.Server, error)

	// DeleteKubeService deletes a named kubernetes service.
	DeleteKubeService(ctx context.Context, name string) error

	// DeleteAllKubeServices deletes all registered kubernetes services.
	DeleteAllKubeServices(context.Context) error
}
