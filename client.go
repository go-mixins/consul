// Package consul provides simple Consul client API wrapper that handles
// registering and deregistering of service
package consul

import (
	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

// Client provides simplified Consul-aware application and access to key-value
// store. The structure members specify configuration parameters and must be
// configured externally, e.g. using envconfig.
type Client struct {
	// Address is the address of the Consul server
	ConsulAddress string `default:"127.0.0.1:8500"`
	// Datacenter to use. If not provided, the default agent datacenter is used.
	ConsulDatacenter string

	// Service data to register. Must be supplied externally before the first
	// call of Init()
	ID      string
	Name    string
	Tags    []string
	Port    int
	Address string

	*consul.Client
	connected bool
}

// Init must be called to register service in Consul and to use KV store
func (client *Client) Init() (err error) {
	config := consul.DefaultConfig()
	config.Address = client.ConsulAddress
	config.Datacenter = client.ConsulDatacenter

	if client.Client, err = consul.NewClient(config); err != nil {
		return errors.Wrap(err, "connecting to consul")
	}
	if err = client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      client.ID,
		Name:    client.Name,
		Tags:    client.Tags,
		Port:    client.Port,
		Address: client.Address,
	}); err != nil {
		return errors.Wrap(err, "registering service")
	}
	client.connected = true
	return
}

// Close must be called to deregister the service and close connection to
// Consul
func (client *Client) Close() (err error) {
	if !client.connected {
		return
	}
	id := client.ID
	if client.ID == "" {
		id = client.Name
	}
	if err = client.Agent().ServiceDeregister(id); err != nil {
		return errors.Wrap(err, "deregistering service")
	}
	return errors.Wrap(client.Close(), "closing client")
}
