/*
 * Copyright 2017 Ludovic Fauvet <etix@l0cal.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package xdgraph provides a simple helper for manipulating Dgraph
// gRPC responses.
package xdgraph

import (
	"encoding/json"
	"time"

	"github.com/dgraph-io/dgraph/protos"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkb"
)

// ReadResponse is the entry point of this package
// It takes in parameter the response from a gRPC call:
//  resp, _ := c.Run(...)
//  xd := xdgraph.ReadResponse(resp)
//  [...]
func ReadResponse(resp *protos.Response) *Response {
	return &Response{
		nodes: []*protos.Node{resp.GetN()[0]},
	}
}

// Response is a struct that carries the current graph.Node.
type Response struct {
	nodes []*protos.Node
}

// First can be used to access the first attribute without explicitely
// giving its name.
func (r Response) First() Response {
	if len(r.nodes) == 0 {
		return Response{}
	}
	if len(r.nodes[0].GetChildren()) == 0 {
		return Response{}
	}
	return Response{nodes: []*protos.Node{r.nodes[0].GetChildren()[0]}}
}

// Attribute moves to the given attribute name. It must be a children of
// the current attribute. This can be asserted using the IsNil() function.
func (r Response) Attribute(name string) Response {
	if len(r.nodes) == 0 {
		return Response{}
	}
	var nodes []*protos.Node
	for _, n := range r.nodes {
		for _, c := range n.GetChildren() {
			if c.GetAttribute() == name {
				nodes = append(nodes, c)
			}
		}
	}

	return Response{nodes: nodes}
}

// Property returns the given property by name.
func (r Response) Property(name string) Property {
	if len(r.nodes) == 0 {
		return Property{}
	}
	for _, p := range r.nodes[0].GetProperties() {
		if p.GetProp() == name {
			return Property{value: p.GetValue()}
		}
	}
	return Property{}
}

// Properties returns a slice of properties by name.
func (r Response) Properties(name string) []Property {
	var properties []Property
	for _, n := range r.nodes {
		for _, p := range n.GetProperties() {
			if p.GetProp() == name {
				properties = append(properties, Property{value: p.GetValue()})
			}
		}
	}
	return properties
}

// Each will run the provided function for each elements contained in the
// response. Example:
//  xd.First().Attribute("follows").Each(func(r xdgraph.Response) {
//      fmt.Println(r.Property("name").ToString())
//  })
func (r Response) Each(fn func(Response)) {
	for _, n := range r.nodes {
		fn(Response{nodes: []*protos.Node{n}})
	}
}

// String returns the attribute content in RAW format.
func (r Response) String() string {
	return r.Json()
}

// Json returns the attribute content in JSON format.
func (r Response) Json() string {
	if len(r.nodes) == 0 {
		return ""
	}
	j, _ := json.MarshalIndent(r.nodes, "", "    ")
	return string(j)
}

// IsNil returns true if the attribute is not available in the response.
func (r Response) IsNil() bool {
	if len(r.nodes) == 0 {
		return true
	}
	return false
}

// Property is a struct that carries the current graph.Value.
type Property struct {
	value *protos.Value
}

// String returns the RAW value
func (p Property) String() string {
	return p.value.String()
}

// ToString returns the property as a string
func (p Property) ToString() string {
	str := p.value.GetStrVal()
	if len(str) > 0 {
		return str
	}
	return p.value.GetDefaultVal()
}

// ToBytes returns the property as []byte
func (p Property) ToBytes() []byte {
	return p.value.GetBytesVal()
}

// ToInt returns the property as an int64
func (p Property) ToInt() int64 {
	return p.value.GetIntVal()
}

// ToBool returns the property as a bool
func (p Property) ToBool() bool {
	return p.value.GetBoolVal()
}

// ToFloat returns the property as a float64
func (p Property) ToFloat() float64 {
	return p.value.GetDoubleVal()
}

// ToGeo returns the property as a geom.T
func (p Property) ToGeo() geom.T {
	t, _ := wkb.Unmarshal(p.value.GetGeoVal())
	return t
}

// ToDate returns the property as a time.Time
func (p Property) ToDate() time.Time {
	var t time.Time
	t.UnmarshalBinary(p.value.GetDateVal())
	return t
}

// ToDateTime returns the property as a time.Time
func (p Property) ToDateTime() time.Time {
	var t time.Time
	t.UnmarshalBinary(p.value.GetDatetimeVal())
	return t
}

// ToPassword returns the property as a string
func (p Property) ToPassword() string {
	return p.value.GetPasswordVal()
}

// ToUid returns the property as a uint64
func (p Property) ToUid() uint64 {
	return p.value.GetUidVal()
}

// IsNil returns true if the property is not available in the response.
func (p Property) IsNil() bool {
	return p.value == nil
}
