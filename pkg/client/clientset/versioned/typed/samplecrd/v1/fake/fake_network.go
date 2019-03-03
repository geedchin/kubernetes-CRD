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

package fake

import (
	samplecrdv1 "crd-demo1/pkg/apis/samplecrd/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeNetworks implements NetworkInterface
type FakeNetworks struct {
	Fake *FakeSamplecrdV1
	ns   string
}

var networksResource = schema.GroupVersionResource{Group: "samplecrd.k8s.io", Version: "v1", Resource: "networks"}

var networksKind = schema.GroupVersionKind{Group: "samplecrd.k8s.io", Version: "v1", Kind: "Network"}

// Get takes name of the network, and returns the corresponding network object, and an error if there is any.
func (c *FakeNetworks) Get(name string, options v1.GetOptions) (result *samplecrdv1.Network, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(networksResource, c.ns, name), &samplecrdv1.Network{})

	if obj == nil {
		return nil, err
	}
	return obj.(*samplecrdv1.Network), err
}

// List takes label and field selectors, and returns the list of Networks that match those selectors.
func (c *FakeNetworks) List(opts v1.ListOptions) (result *samplecrdv1.NetworkList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(networksResource, networksKind, c.ns, opts), &samplecrdv1.NetworkList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &samplecrdv1.NetworkList{ListMeta: obj.(*samplecrdv1.NetworkList).ListMeta}
	for _, item := range obj.(*samplecrdv1.NetworkList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested networks.
func (c *FakeNetworks) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(networksResource, c.ns, opts))

}

// Create takes the representation of a network and creates it.  Returns the server's representation of the network, and an error, if there is any.
func (c *FakeNetworks) Create(network *samplecrdv1.Network) (result *samplecrdv1.Network, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(networksResource, c.ns, network), &samplecrdv1.Network{})

	if obj == nil {
		return nil, err
	}
	return obj.(*samplecrdv1.Network), err
}

// Update takes the representation of a network and updates it. Returns the server's representation of the network, and an error, if there is any.
func (c *FakeNetworks) Update(network *samplecrdv1.Network) (result *samplecrdv1.Network, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(networksResource, c.ns, network), &samplecrdv1.Network{})

	if obj == nil {
		return nil, err
	}
	return obj.(*samplecrdv1.Network), err
}

// Delete takes name of the network and deletes it. Returns an error if one occurs.
func (c *FakeNetworks) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(networksResource, c.ns, name), &samplecrdv1.Network{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNetworks) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(networksResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &samplecrdv1.NetworkList{})
	return err
}

// Patch applies the patch and returns the patched network.
func (c *FakeNetworks) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *samplecrdv1.Network, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(networksResource, c.ns, name, pt, data, subresources...), &samplecrdv1.Network{})

	if obj == nil {
		return nil, err
	}
	return obj.(*samplecrdv1.Network), err
}
