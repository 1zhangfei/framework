package grpc

import (
	"encoding/json"
	"fmt"
	"github.com/1zhangfei/framework/config"
	"github.com/1zhangfei/framework/consul"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"math/rand"
	"net"
	"time"
)

type Cnf struct {
	App struct {
		Ip   string
		Port int
	} `json:"rpc"`
	Name   string `json:"tokenName"`
	Consul string `json:"consul"`
}

func getGrpcConfig(address string) (*Cnf, error) {
	err := config.ViperInit(address)
	if err != nil {
		return nil, err
	}

	dataId := viper.GetString("Grpc.DataId")
	group := viper.GetString("Grpc.Group")

	if err != nil {
		return nil, err
	}
	var conf Cnf
	getConfig, err := config.GetConfig(dataId, group)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(getConfig), &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func Service(address string, Register func(s *grpc.Server)) error {
	rand.Seed(time.Now().UnixNano())
	conf, err2 := getGrpcConfig(address)
	if err2 != nil {
		return err2
	}
	port := rand.Intn(9000) + 1000
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 服务注册
	s := grpc.NewServer()

	err = consul.GetRegister(conf.App.Ip, conf.Name, port)

	reflection.Register(s)
	// 健康检测
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}
