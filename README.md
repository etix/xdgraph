# xdgraph: a simple Go helper for manipulating Dgraph gRPC responses

### Installation
`go get github.com/etix/xdgraph`

### Usage
```go
package main

import (
    "github.com/etix/xdgraph"
    "..."
)

int main() {
    // Set up the gRPC connection
    conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Initialize the Dgraph gRPC object
    c := graph.NewDgraphClient(conn)
    req := client.Req{}
    req.SetQuery(`<your query>`)

    // Execute the query and fetch the response
    resp, err := c.Run(context.Background(), req.Request())
    if err != nil {
        log.Fatal(err)
    }

    // Let's use xdgraph to manipulate the results!
    xd := xdgraph.ReadResponse(resp)

    // Get the UID of the "me" attribute
    fmt.Println(xd.Attribute("me").Uid())

    // Get the "name" property of the "me" attribute
    fmt.Println(xd.Attribute("me").Property("name").ToString())

    // Skipping the first attribute using First()
    fmt.Println(xd.First().Property("sign").ToString())

    // Traverse multiple attributes
    fmt.Println(xd.Attribute("me").Attribute("follows").Property("name").ToString())

    // Use different types
    fmt.Println(xd.Attribute("me").Attribute("follows").Property("birthdate").ToDate())

    // Check if a property is set (true/false)
    fmt.Println(xd.Attribute("me").Property("address").IsNil())

    // Output the full graph in JSON
    fmt.Println(xd.Json())

    // Output a sub-graph in JSON
    fmt.Println(xd.Attribute("me").Attribute("follows").Json())
}
```
See the example folder for a complete example.

### Supported property types
[See here](https://github.com/dgraph-io/dgraph/blob/master/query/graph/graphresponse.pb.go) for the list of supported types upstream.

| Type        |                            |
| ----------- | -------------------------- |
| string      | .ToString()                |
| []byte      | .ToBytes()                 |
| int32       | .ToInt()                   |
| bool        | .ToBool()                  |
| float64     | .ToFloat()                 |
| geom        | .ToGeo()                   |
| datetime    | .ToDate() or .ToDateTime() |

### Note
Since the [Dgraph](https://github.com/dgraph-io/dgraph/) project is quite young, the APIs can change at any time.
