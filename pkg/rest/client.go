/*
Copyright 2020 The go-harbor Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
*/

package rest

import (
	flowcontrol2 "github.com/hujianxiong/go-harbor/pkg/rest/util/flowcontrol"
	"net/http"
	"net/url"
	"strings"
)

// Interface captures the set of operations for generically interacting with Kubernetes REST apis.
type Interface interface {
	Verb(verb string) *Request
	Post() *Request
	Put() *Request
	List() *Request
	Get() *Request
	Delete() *Request
}

// RESTClient imposes common Kubernetes API conventions on a set of resource paths.
// The baseURL is expected to point to an HTTP or HTTPS path that is the parent
// of one or more resources.  The server should return a decodable API resource
// object, or an api.Status object which contains information about the reason for
// any failure.
//
// Most consumers should use client.New() to get a Kubernetes API client.
type RESTClient struct {
	// base is the root URL for all invocations of the client
	base *url.URL
	// versionedAPIPath is a path segment connecting the base URL to the resource root
	versionedAPIPath string

	// contentConfig is the information used to communicate with the server.
	contentConfig ContentConfig

	Throttle flowcontrol2.RateLimiter
	headers  map[string]string
	// Set specific behavior of the client.  If not set http.DefaultClient will be used.
	Client *http.Client
}

func (c *RESTClient) List() *Request {
	return c.Verb("GET")
}

func (c *RESTClient) Post() *Request {
	return c.Verb("POST")
}

func (c *RESTClient) Put() *Request {
	return c.Verb("PUT")
}

// NewRESTClient creates a new RESTClient. This client performs generic REST functions
// such as Get, Put, Post, and Delete on specified paths.  Codec controls encoding and
// decoding of responses from the server.
func NewRESTClient(baseURL *url.URL, versionedAPIPath string, config ContentConfig, headers map[string]string, maxQPS float32, maxBurst int, rateLimiter flowcontrol2.RateLimiter, client *http.Client) (*RESTClient, error) {
	base := *baseURL
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}
	base.RawQuery = ""
	base.Fragment = ""

	/*	if config.GroupVersion == nil {
		config.GroupVersion = &schema.GroupVersion{}
	}*/
	if len(config.ContentType) == 0 {
		config.ContentType = "application/json"
	}
	/*	serializers, err := createSerializers(config)
		if err != nil {
			return nil, err
		}*/

	var throttle flowcontrol2.RateLimiter
	if maxQPS > 0 && rateLimiter == nil {
		throttle = flowcontrol2.NewTokenBucketRateLimiter(maxQPS, maxBurst)
	} else if rateLimiter != nil {
		throttle = rateLimiter
	}
	return &RESTClient{
		base:             &base,
		versionedAPIPath: versionedAPIPath,
		contentConfig:    config,
		Throttle:         throttle,
		headers:          headers,
		Client:           client,
	}, nil
}

// Get begins a GET request. Short for c.Verb("GET").
func (c *RESTClient) Get() *Request {
	return c.Verb("GET")
}

// Delete begins a DELETE request. Short for c.Verb("DELETE").
func (c *RESTClient) Delete() *Request {
	return c.Verb("DELETE")
}

// Verb begins a request with a verb (GET, POST, PUT, DELETE).
//
// Example usage of RESTClient's request building interface:
// c, err := NewRESTClient(...)
// if err != nil { ... }
// resp, err := c.Verb("GET").
//  Path("pods").
//  SelectorParam("labels", "area=staging").
//  Timeout(10*time.Second).
//  Do()
// if err != nil { ... }
// list, ok := resp.(*api.PodList)
//
func (c *RESTClient) Verb(verb string) *Request {
	if c.Client == nil {
		return NewRequest(nil, verb, c.base, c.headers, c.versionedAPIPath, c.contentConfig, c.Throttle, 0)
	}
	return NewRequest(c.Client, verb, c.base, c.headers, c.versionedAPIPath, c.contentConfig, c.Throttle, c.Client.Timeout)
}
