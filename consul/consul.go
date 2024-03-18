package consul

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

func GetRegister(ip, name string, port int) error {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      uuid.NewString(),
		Name:    name,
		Tags:    nil,
		Port:    port,
		Address: ip,
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           fmt.Sprintf("%v:%v", ip, port),
			DeregisterCriticalServiceAfter: "10s",
		},
	})
	if err != nil {
		return err
	}
	return nil
}
