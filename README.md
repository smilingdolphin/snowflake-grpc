# snowflake-grpc
A simple snowflake generator in grpc server.

# snowflake grpc service
#### grpc client example
```go
package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/fpay/snowflake-go/pb"
)

const (
	Address = "127.0.0.1:11070" // config.sample.yaml server port
)

func main() {
    conn, err := grpc.Dial(Address, grpc.WithInsecure())
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    c := pb.NewSnowflakeServiceClient(conn)
    r, err := c.Generate(context.Background(), &pb.Request{})
    if err != nil {
        panic(err)
    }
    fmt.Println(r.GetUniqid())
}
```
