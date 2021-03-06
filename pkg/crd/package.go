/*
Copyright 2016 The Fission Authors.

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

package crd

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	fv1 "github.com/fission/fission/pkg/apis/fission.io/v1"
)

type (
	PackageInterface interface {
		Create(*fv1.Package) (*fv1.Package, error)
		Get(name string) (*fv1.Package, error)
		Update(*fv1.Package) (*fv1.Package, error)
		Delete(name string, options *metav1.DeleteOptions) error
		List(opts metav1.ListOptions) (*fv1.PackageList, error)
		Watch(opts metav1.ListOptions) (watch.Interface, error)
	}

	packageClient struct {
		client    *rest.RESTClient
		namespace string
	}
)

func MakePackageInterface(crdClient *rest.RESTClient, namespace string) PackageInterface {
	return &packageClient{
		client:    crdClient,
		namespace: namespace,
	}
}

func (c *packageClient) Create(f *fv1.Package) (*fv1.Package, error) {
	var result fv1.Package
	err := c.client.Post().
		Resource("packages").
		Namespace(c.namespace).
		Body(f).
		Do().Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *packageClient) Get(name string) (*fv1.Package, error) {
	var result fv1.Package
	err := c.client.Get().
		Resource("packages").
		Namespace(c.namespace).
		Name(name).
		Do().Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *packageClient) Update(f *fv1.Package) (*fv1.Package, error) {
	var result fv1.Package
	err := c.client.Put().
		Resource("packages").
		Namespace(c.namespace).
		Name(f.Metadata.Name).
		Body(f).
		Do().Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *packageClient) Delete(name string, opts *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.namespace).
		Resource("packages").
		Name(name).
		Body(opts).
		Do().
		Error()
}

func (c *packageClient) List(opts metav1.ListOptions) (*fv1.PackageList, error) {
	var result fv1.PackageList
	err := c.client.Get().
		Namespace(c.namespace).
		Resource("packages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *packageClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Get().
		Prefix("watch").
		Namespace(c.namespace).
		Resource("packages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
