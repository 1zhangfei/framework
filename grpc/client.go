package grpc

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"google.golang.org/grpc"
)

func Client(address string) (*grpc.ClientConn, error) {
	conf, err2 := getGrpcConfig(address)
	if err2 != nil {
		return nil, err2
	}

	return grpc.Dial(
		fmt.Sprintf("consul://%v:%v/%v?wait=14s", conf.App.Ip, conf.Consul, conf.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

}
