/*
Copyright 2023 James Riley O'Donnell.

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
	"context"

	v1alpha1 "github.com/the-gizmo-dojo/g8s/pkg/apis/api.g8s.io/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePasswords implements PasswordInterface
type FakePasswords struct {
	Fake *FakeApiV1alpha1
	ns   string
}

var passwordsResource = v1alpha1.SchemeGroupVersion.WithResource("passwords")

var passwordsKind = v1alpha1.SchemeGroupVersion.WithKind("Password")

// Get takes name of the password, and returns the corresponding password object, and an error if there is any.
func (c *FakePasswords) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Password, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(passwordsResource, c.ns, name), &v1alpha1.Password{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Password), err
}

// List takes label and field selectors, and returns the list of Passwords that match those selectors.
func (c *FakePasswords) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.PasswordList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(passwordsResource, passwordsKind, c.ns, opts), &v1alpha1.PasswordList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.PasswordList{ListMeta: obj.(*v1alpha1.PasswordList).ListMeta}
	for _, item := range obj.(*v1alpha1.PasswordList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested passwords.
func (c *FakePasswords) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(passwordsResource, c.ns, opts))

}

// Create takes the representation of a password and creates it.  Returns the server's representation of the password, and an error, if there is any.
func (c *FakePasswords) Create(ctx context.Context, password *v1alpha1.Password, opts v1.CreateOptions) (result *v1alpha1.Password, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(passwordsResource, c.ns, password), &v1alpha1.Password{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Password), err
}

// Update takes the representation of a password and updates it. Returns the server's representation of the password, and an error, if there is any.
func (c *FakePasswords) Update(ctx context.Context, password *v1alpha1.Password, opts v1.UpdateOptions) (result *v1alpha1.Password, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(passwordsResource, c.ns, password), &v1alpha1.Password{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Password), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePasswords) UpdateStatus(ctx context.Context, password *v1alpha1.Password, opts v1.UpdateOptions) (*v1alpha1.Password, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(passwordsResource, "status", c.ns, password), &v1alpha1.Password{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Password), err
}

// Delete takes name of the password and deletes it. Returns an error if one occurs.
func (c *FakePasswords) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(passwordsResource, c.ns, name, opts), &v1alpha1.Password{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePasswords) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(passwordsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.PasswordList{})
	return err
}

// Patch applies the patch and returns the patched password.
func (c *FakePasswords) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Password, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(passwordsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Password{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Password), err
}
