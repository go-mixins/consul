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
	Address string `default:"localhost:8500"`

	// Scheme is the URI scheme for the Consul server
	Scheme string `default:"http"`

	// Datacenter to use. If not provided, the default agent datacenter is used.
	Datacenter string `default:""`

	// Service data to register. Must be supplied externally before the first
	// call of Init()
	Service consul.AgentServiceRegistration `ignored:"true"`

	*consul.Client
}

// Init must be called to register service in Consul and to use KV store
func (client *Client) Init() (err error) {
	config := consul.DefaultConfig()
	config.Address = client.Address
	config.Scheme = client.Scheme
	config.Datacenter = client.Datacenter
	if client.Client, err = consul.NewClient(config); err != nil {
		return errors.Wrap(err, "connecting to consul")
	}
	if err = client.Agent().ServiceRegister(&client.Service); err != nil {
		return errors.Wrap(err, "registering service")
	}
	return
}

// Close must be called to deregister register service and close connection to
// Consul
func (client *Client) Close() (err error) {
	if err = client.Agent().ServiceDeregister(client.Service.ID); err != nil {
		return errors.Wrap(err, "deregistering service")
	}
	return errors.Wrap(client.Close(), "closing client")
}
