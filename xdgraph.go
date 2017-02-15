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

package xdgraph

import (
    "encoding/json"
    "time"

    "github.com/dgraph-io/dgraph/query/graph"
    proto "github.com/golang/protobuf/proto"
    "github.com/twpayne/go-geom"
    "github.com/twpayne/go-geom/encoding/wkb"
)

func ReadResponse(resp *graph.Response) *response {
    return &response{
        node: resp.GetN()[0],
    }
}

type response struct {
    node *graph.Node
}

func (r response) First() response {
    if len(r.node.GetChildren()) == 0 {
        return response{}
    }
    return response{node: r.node.GetChildren()[0]}
}

func (r response) Attribute(name string) response {
    for _, c := range r.node.GetChildren() {
        if c.GetAttribute() == name {
            return response{node: c}
            break
        }
    }
    return response{}
}

func (r response) Property(name string) property {
    for _, p := range r.node.GetProperties() {
        if p.GetProp() == name {
            return property{value: p.GetValue()}
        }
    }
    return property{}
}

func (r response) Uid() uint64 {
    return r.node.GetUid()
}

func (r response) Xid() string {
    return r.node.GetXid()
}

func (r response) String() string {
    return proto.CompactTextString(r.node)
}

func (r response) Json() string {
    j, _ := json.MarshalIndent(r.node, "", "    ")
    return string(j)
}

func (r response) IsNil() bool {
    return r.node == nil
}

type property struct {
    value *graph.Value
}

func (v property) String() string {
    return v.value.String()
}

func (v property) ToString() string {
    return v.value.GetStrVal()
}

func (v property) ToBytes() []byte {
    return v.value.GetBytesVal()
}

func (v property) ToInt() int32 {
    return v.value.GetIntVal()
}

func (v property) ToBool() bool {
    return v.value.GetBoolVal()
}

func (v property) ToFloat() float64 {
    return v.value.GetDoubleVal()
}

func (v property) ToGeo() geom.T {
    t, _ := wkb.Unmarshal(v.value.GetGeoVal())
    return t
}

func (v property) ToDate() time.Time {
    var t time.Time
    t.UnmarshalBinary(v.value.GetDateVal())
    return t
}

func (v property) ToDateTime() time.Time {
    var t time.Time
    t.UnmarshalBinary(v.value.GetDatetimeVal())
    return t
}

func (v property) ToPassword() string {
    return v.value.GetPasswordVal()
}

func (v property) IsNil() bool {
    return v.value == nil
}
