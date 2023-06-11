package main

import (
	"economic_srvs/user_srv/handler"
	"economic_srvs/user_srv/proto"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

/*
*
gRPC Service Layer.
*/
func main() {
	// appoint ip and port by flag
	IP := flag.String("ip", "0.0.0.0", "ip address")
	Port := flag.Int("port", 8088, "ip port")
	flag.Parse()
	fmt.Println("ip:", *IP)
	fmt.Println("port:", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	// launch service
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("Oh No! listen error: " + err.Error())
	}
	err = server.Serve(listener)
	if err != nil {
		panic("Oh No! grpc error: " + err.Error())
	}
}
