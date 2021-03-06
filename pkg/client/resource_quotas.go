/*
Copyright 2014 Google Inc. All rights reserved.

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

package client

import (
	"errors"
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

// ResourceQuotasNamespacer has methods to work with ResourceQuota resources in a namespace
type ResourceQuotasNamespacer interface {
	ResourceQuotas(namespace string) ResourceQuotaInterface
}

// ResourceQuotaInterface has methods to work with ResourceQuota resources.
type ResourceQuotaInterface interface {
	List(selector labels.Selector) (*api.ResourceQuotaList, error)
	Get(name string) (*api.ResourceQuota, error)
	Delete(name string) error
	Create(resourceQuota *api.ResourceQuota) (*api.ResourceQuota, error)
	Update(resourceQuota *api.ResourceQuota) (*api.ResourceQuota, error)
	Watch(label, field labels.Selector, resourceVersion string) (watch.Interface, error)
}

// resourceQuotas implements ResourceQuotasNamespacer interface
type resourceQuotas struct {
	r  *Client
	ns string
}

// newResourceQuotas returns a resourceQuotas
func newResourceQuotas(c *Client, namespace string) *resourceQuotas {
	return &resourceQuotas{
		r:  c,
		ns: namespace,
	}
}

// List takes a selector, and returns the list of resourceQuotas that match that selector.
func (c *resourceQuotas) List(selector labels.Selector) (result *api.ResourceQuotaList, err error) {
	result = &api.ResourceQuotaList{}
	err = c.r.Get().Namespace(c.ns).Resource("resourceQuotas").SelectorParam("labels", selector).Do().Into(result)
	return
}

// Get takes the name of the resourceQuota, and returns the corresponding ResourceQuota object, and an error if it occurs
func (c *resourceQuotas) Get(name string) (result *api.ResourceQuota, err error) {
	if len(name) == 0 {
		return nil, errors.New("name is required parameter to Get")
	}

	result = &api.ResourceQuota{}
	err = c.r.Get().Namespace(c.ns).Resource("resourceQuotas").Name(name).Do().Into(result)
	return
}

// Delete takes the name of the resourceQuota, and returns an error if one occurs
func (c *resourceQuotas) Delete(name string) error {
	return c.r.Delete().Namespace(c.ns).Resource("resourceQuotas").Name(name).Do().Error()
}

// Create takes the representation of a resourceQuota.  Returns the server's representation of the resourceQuota, and an error, if it occurs.
func (c *resourceQuotas) Create(resourceQuota *api.ResourceQuota) (result *api.ResourceQuota, err error) {
	result = &api.ResourceQuota{}
	err = c.r.Post().Namespace(c.ns).Resource("resourceQuotas").Body(resourceQuota).Do().Into(result)
	return
}

// Update takes the representation of a resourceQuota to update.  Returns the server's representation of the resourceQuota, and an error, if it occurs.
func (c *resourceQuotas) Update(resourceQuota *api.ResourceQuota) (result *api.ResourceQuota, err error) {
	result = &api.ResourceQuota{}
	if len(resourceQuota.ResourceVersion) == 0 {
		err = fmt.Errorf("invalid update object, missing resource version: %v", resourceQuota)
		return
	}
	err = c.r.Put().Namespace(c.ns).Resource("resourceQuotas").Name(resourceQuota.Name).Body(resourceQuota).Do().Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested resource
func (c *resourceQuotas) Watch(label, field labels.Selector, resourceVersion string) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("resourceQuotas").
		Param("resourceVersion", resourceVersion).
		SelectorParam("labels", label).
		SelectorParam("fields", field).
		Watch()
}
