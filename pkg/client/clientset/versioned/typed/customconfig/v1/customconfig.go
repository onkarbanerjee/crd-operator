/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/onkarbanerjee/crd-custom-config/pkg/apis/customconfig/v1"
	scheme "github.com/onkarbanerjee/crd-custom-config/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CustomConfigsGetter has a method to return a CustomConfigInterface.
// A group's client should implement this interface.
type CustomConfigsGetter interface {
	CustomConfigs(namespace string) CustomConfigInterface
}

// CustomConfigInterface has methods to work with CustomConfig resources.
type CustomConfigInterface interface {
	Create(ctx context.Context, customConfig *v1.CustomConfig, opts metav1.CreateOptions) (*v1.CustomConfig, error)
	Update(ctx context.Context, customConfig *v1.CustomConfig, opts metav1.UpdateOptions) (*v1.CustomConfig, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.CustomConfig, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.CustomConfigList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.CustomConfig, err error)
	CustomConfigExpansion
}

// customConfigs implements CustomConfigInterface
type customConfigs struct {
	client rest.Interface
	ns     string
}

// newCustomConfigs returns a CustomConfigs
func newCustomConfigs(c *MtcilV1Client, namespace string) *customConfigs {
	return &customConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the customConfig, and returns the corresponding customConfig object, and an error if there is any.
func (c *customConfigs) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.CustomConfig, err error) {
	result = &v1.CustomConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("customconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CustomConfigs that match those selectors.
func (c *customConfigs) List(ctx context.Context, opts metav1.ListOptions) (result *v1.CustomConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.CustomConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("customconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested customConfigs.
func (c *customConfigs) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("customconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a customConfig and creates it.  Returns the server's representation of the customConfig, and an error, if there is any.
func (c *customConfigs) Create(ctx context.Context, customConfig *v1.CustomConfig, opts metav1.CreateOptions) (result *v1.CustomConfig, err error) {
	result = &v1.CustomConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("customconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(customConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a customConfig and updates it. Returns the server's representation of the customConfig, and an error, if there is any.
func (c *customConfigs) Update(ctx context.Context, customConfig *v1.CustomConfig, opts metav1.UpdateOptions) (result *v1.CustomConfig, err error) {
	result = &v1.CustomConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("customconfigs").
		Name(customConfig.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(customConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the customConfig and deletes it. Returns an error if one occurs.
func (c *customConfigs) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("customconfigs").
		Name(name).
		Body(&opts).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *customConfigs) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("customconfigs").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do().
		Error()
}

// Patch applies the patch and returns the patched customConfig.
func (c *customConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.CustomConfig, err error) {
	result = &v1.CustomConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("customconfigs").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do().
		Into(result)
	return
}
