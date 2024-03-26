package redis

import (
	"context"
	"encoding/json"
	"github.com/1zhangfei/framework/config"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"time"
)

func WhitRedislInit(address string, hand func(cli *redis.Client) (string, error)) (string, error) {
	err := config.ViperInit(address)
	if err != nil {
		return "", err
	}
	data := viper.GetString("Database.DataId")
	Group := viper.GetString("Database.Group")
	var App struct {
		R struct {
			Host string
			Port string
		} `json:"Redis"`
	}

	getConfig, err := config.GetConfig(data, Group)
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal([]byte(getConfig), &App); err != nil {
		return "", err
	}

	Rdb := redis.NewClient(&redis.Options{
		Addr: App.R.Host + ":" + App.R.Port,
		DB:   0,
	})
	res, err := hand(Rdb)

	if err != nil {
		return "", err
	}

	return res, nil

}

func GetKey(key, address string) (string, error) {
	init, err := WhitRedislInit(address, func(cli *redis.Client) (string, error) {
		val, err := cli.Get(context.Background(), key).Result()
		if err != nil {
			return "", err
		}
		return val, nil
	})
	if err != nil {
		return "", err
	}
	return init, err

}

func SetByVal(Address, Key string, val interface{}, duration time.Duration) error {
	_, err := WhitRedislInit(Address, func(cli *redis.Client) (string, error) {
		err := cli.Set(context.Background(), Key, val, duration).Err()
		return "", err

	})
	if err != nil {
		return err
	}
	return nil
}

func Exists(Address, key string) (bool, error) {
	res, err := WhitRedislInit(Address, func(cli *redis.Client) (string, error) {
		result, err := cli.Exists(context.Background(), key).Result()
		if err != nil {
			return "", err
		}
		if result != 0 {
			return "1", nil
		}
		return "", nil
	})
	if err != nil {
		return false, err
	}
	if res == "1" {
		return true, nil
	}

	return false, nil

}

func Lock(address, key string, val interface{}, duration time.Duration, isReadOnly bool) (bool, error) {
	res := false
	_, err := WhitRedislInit(address, func(cli *redis.Client) (string, error) {
		if !isReadOnly {
			for {
				result, err := cli.SetNX(context.Background(), key, val, duration).Result()
				if err != nil {
					return "", err
				}
				res = true
				if result {
					return "", nil
				}
			}

		}
		re, err := cli.SetNX(context.Background(), key, val, duration).Result()
		res = re
		return "", err
	})
	return res, err
}

func UnLock(address, key string) error {
	_, err := WhitRedislInit(address, func(cli *redis.Client) (string, error) {
		err := cli.Del(context.Background(), key).Err()
		return "", err
	})
	if err != nil {
		return err
	}
	return nil
}
