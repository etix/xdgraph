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

    "github.com/dgraph-io/dgraph/query/graph"
    "github.com/twpayne/go-geom"
    "github.com/twpayne/go-geom/encoding/wkb"
)

// ReadResponse is the entry point of this package
// It takes in parameter the response from a gRPC call:
//  resp, _ := c.Run(...)
//  xd := xdgraph.ReadResponse(resp)
//  [...]
func ReadResponse(resp *graph.Response) *Response {
    return &Response{
        node: resp.GetN()[0],
    }
}

// Response is a struct that carries the current graph.Node.
type Response struct {
    node *graph.Node
}

// First can be used to access the first attribute without explicitely
// giving its name.
func (r Response) First() Response {
    if len(r.node.GetChildren()) == 0 {
        return Response{}
    }
    return Response{node: r.node.GetChildren()[0]}
}

// Attribute moves to the given attribute name. It must be a children of
// the current attribute. This can be asserted using the IsNil() function.
func (r Response) Attribute(name string) Response {
    for _, c := range r.node.GetChildren() {
        if c.GetAttribute() == name {
            return Response{node: c}
            break
        }
    }
    return Response{}
}

// Property returns the given property by name.
func (r Response) Property(name string) Property {
    for _, p := range r.node.GetProperties() {
        if p.GetProp() == name {
            return Property{value: p.GetValue()}
        }
    }
    return Property{}
}

// Uid returns the UID of the current attribute if contained in the response.
func (r Response) Uid() uint64 {
    return r.node.GetUid()
}

// Xid returns the XID of the current attribute if contained in the response.
func (r Response) Xid() string {
    // BUG(r): GetXid() doesn't seem to be supported by the upstream
    return r.node.GetXid()
}

// String returns the attribute content in RAW format.
func (r Response) String() string {
    return r.Json()
}

// Json returns the attribute content in JSON format.
func (r Response) Json() string {
    j, _ := json.MarshalIndent(r.node, "", "    ")
    return string(j)
}

// IsNil returns true if the attribute is not available in the response.
func (r Response) IsNil() bool {
    return r.node == nil
}

// Property is a struct that carries the current graph.Value.
type Property struct {
    value *graph.Value
}

// String returns the RAW value
func (p Property) String() string {
    return p.value.String()
}

// ToString returns the property as a string
func (p Property) ToString() string {
    return p.value.GetStrVal()
}

// ToBytes returns the property as []byte
func (p Property) ToBytes() []byte {
    return p.value.GetBytesVal()
}

// ToInt returns the property as an int32
func (p Property) ToInt() int32 {
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

// IsNil returns true if the property is not available in the response.
func (p Property) IsNil() bool {
    return p.value == nil
}
