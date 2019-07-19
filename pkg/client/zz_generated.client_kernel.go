/*
	Note: This file is autogenerated! Do not edit it manually!
	Edit client_kernel_template.go instead, and run
	hack/generate-client.sh afterwards.
*/

package client

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/filterer"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// KernelClient is an interface for accessing Kernel-specific API objects
type KernelClient interface {
	// New returns a new Kernel
	New() *api.Kernel
	// Get returns the Kernel matching given UID from the storage
	Get(meta.UID) (*api.Kernel, error)
	// Set saves the given Kernel into persistent storage
	Set(*api.Kernel) error
	// Find returns the Kernel matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	Find(filter filterer.BaseFilter) (*api.Kernel, error)
	// FindAll returns multiple Kernels matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	FindAll(filter filterer.BaseFilter) ([]*api.Kernel, error)
	// Delete deletes the Kernel with the given UID from the storage
	Delete(uid meta.UID) error
	// List returns a list of all Kernels available
	List() ([]*api.Kernel, error)
}

// Kernels returns the KernelClient for the IgniteInternalClient instance
func (c *IgniteInternalClient) Kernels() KernelClient {
	if c.kernelClient == nil {
		c.kernelClient = newKernelClient(c.storage, c.gv)
	}

	return c.kernelClient
}

// kernelClient is a struct implementing the KernelClient interface
// It uses a shared storage instance passed from the Client together with its own Filterer
type kernelClient struct {
	storage  storage.Storage
	filterer *filterer.Filterer
	gvk      schema.GroupVersionKind
}

// newKernelClient builds the kernelClient struct using the storage implementation and a new Filterer
func newKernelClient(s storage.Storage, gv schema.GroupVersion) KernelClient {
	return &kernelClient{
		storage:  s,
		filterer: filterer.NewFilterer(s),
		gvk:      gv.WithKind(api.KindKernel.Title()),
	}
}

// New returns a new Object of its kind
func (c *kernelClient) New() *api.Kernel {
	log.Tracef("Client.New; GVK: %v", c.gvk)
	obj, err := c.storage.New(c.gvk)
	if err != nil {
		panic(fmt.Sprintf("Client.New must not return an error: %v", err))
	}
	return obj.(*api.Kernel)
}

// Find returns a single Kernel based on the given Filter
func (c *kernelClient) Find(filter filterer.BaseFilter) (*api.Kernel, error) {
	log.Tracef("Client.Find; GVK: %v", c.gvk)
	object, err := c.filterer.Find(c.gvk, filter)
	if err != nil {
		return nil, err
	}

	return object.(*api.Kernel), nil
}

// FindAll returns multiple Kernels based on the given Filter
func (c *kernelClient) FindAll(filter filterer.BaseFilter) ([]*api.Kernel, error) {
	log.Tracef("Client.FindAll; GVK: %v", c.gvk)
	matches, err := c.filterer.FindAll(c.gvk, filter)
	if err != nil {
		return nil, err
	}

	results := make([]*api.Kernel, 0, len(matches))
	for _, item := range matches {
		results = append(results, item.(*api.Kernel))
	}

	return results, nil
}

// Get returns the Kernel matching given UID from the storage
func (c *kernelClient) Get(uid meta.UID) (*api.Kernel, error) {
	log.Tracef("Client.Get; UID: %q, GVK: %v", uid, c.gvk)
	object, err := c.storage.Get(c.gvk, uid)
	if err != nil {
		return nil, err
	}

	return object.(*api.Kernel), nil
}

// Set saves the given Kernel into the persistent storage
func (c *kernelClient) Set(kernel *api.Kernel) error {
	log.Tracef("Client.Set; UID: %q, GVK: %v", kernel.GetUID(), c.gvk)
	return c.storage.Set(c.gvk, kernel)
}

// Delete deletes the Kernel from the storage
func (c *kernelClient) Delete(uid meta.UID) error {
	log.Tracef("Client.Delete; UID: %q, GVK: %v", uid, c.gvk)
	return c.storage.Delete(c.gvk, uid)
}

// List returns a list of all Kernels available
func (c *kernelClient) List() ([]*api.Kernel, error) {
	log.Tracef("Client.List; GVK: %v", c.gvk)
	list, err := c.storage.List(c.gvk)
	if err != nil {
		return nil, err
	}

	results := make([]*api.Kernel, 0, len(list))
	for _, item := range list {
		results = append(results, item.(*api.Kernel))
	}

	return results, nil
}
