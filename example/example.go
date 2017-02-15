package main

import (
    "context"
    "flag"
    "fmt"
    "log"

    "github.com/etix/xdgraph"

    "github.com/dgraph-io/dgraph/client"
    "github.com/dgraph-io/dgraph/query/graph"
    "google.golang.org/grpc"
)

var dgraph = flag.String("d", "127.0.0.1:8080", "Dgraph server address")

func main() {
    flag.Parse()

    conn, err := grpc.Dial(*dgraph, grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    c := graph.NewDgraphClient(conn)

    req := client.Req{}

    req.SetQuery(`
        mutation {
            set {
                <bob> <name> "Bob" .
                <bob> <age> "27"^^<xs:int> .
                <bob> <sign> "gemini" .
                <bob> <address> "somewhere" .

                <clara> <name> "Clara" .
                <clara> <age> "32"^^<xs:int> .
                <clara> <sign> "taurus" .

                <clara> <follows> <bob> .
                <bob> <follow> <clara> .
            }
        }

        query {
            me(id: clara) {
                _uid_
                name
                age
                sign
                address
                follows {
                    _uid_
                    name
                    age
                    sign
                    address
                }
            }
        }
    `)

    resp, err := c.Run(context.Background(), req.Request())
    if err != nil {
        log.Fatal(err)
    }

    xd := xdgraph.ReadResponse(resp)

    fmt.Printf("UID: %d\n", xd.Attribute("me").Uid())
    fmt.Printf("Name: %s\n", xd.Attribute("me").Property("name").ToString())
    fmt.Printf("Sign (using First()): %s\n", xd.First().Property("sign").ToString())

    //fmt.Printf("Age: %d\n", xd.Attribute("me").Property("age").ToInt()) // See https://github.com/dgraph-io/dgraph/issues/594

    fmt.Printf("Clara follows: %s (%d)\n",
        xd.Attribute("me").Attribute("follows").Property("name").ToString(),
        xd.Attribute("me").Attribute("follows").Uid())

    fmt.Printf("Is Clara's address set? %t\n", !xd.Attribute("me").Property("address").IsNil())
    fmt.Printf("Is Bob's address set? %t\n", !xd.Attribute("me").Attribute("follows").Property("address").IsNil())

    fmt.Printf("Raw: %s\n", xd.Attribute("me"))

    fmt.Printf("Json (full graph):\n%s\n", xd.Json())
    fmt.Printf("Json (sub graph):\n%s\n", xd.Attribute("me").Attribute("follows").Json())
}
