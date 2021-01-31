package main

import (
	"fmt"
	"time"

	"github.com/stack-labs/stack-rpc/plugin/config/source/apollo/agollo"
	"github.com/stack-labs/stack-rpc/plugin/config/source/apollo/agollo/env/config"
)

func main() {
	c := &config.AppConfig{
		AppID:          "test-demo",
		Cluster:        "DEV",
		IP:             "http://127.0.0.1:8080",
		NamespaceName:  "application,test-registry",
		IsBackupConfig: true,
		Secret:         "",
	}

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})

	if err != nil {
		fmt.Println("err:", err)
		panic(err)
	}

	checkKey(c.NamespaceName, client)

	client, err = agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})

	if err != nil {
		fmt.Println("err:", err)
		panic(err)
	}

	checkKey(c.NamespaceName, client)

	time.Sleep(5 * time.Second)
}

func checkKey(namespace string, client *agollo.Client) {
	cache := client.GetConfigCache(namespace)
	count := 0
	cache.Range(func(key, value interface{}) bool {
		fmt.Println("key : ", key, ", value :", value)
		count++
		return true
	})
	if count < 1 {
		panic("config key can not be null")
	}
}
