package apollo

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/stack-labs/stack/logger"
	"github.com/stack-labs/stack/pkg/config/source"
	apo "github.com/stack-labs/stack/plugin/config/source/apollo/agollo"
	apoC "github.com/stack-labs/stack/plugin/config/source/apollo/agollo/env/config"
)

type apolloSource struct {
	appID      string
	client     *apo.Client
	namespaces []string
	opts       source.Options
}

var (
	DefaultAppID          = "stack"
	DefaultAddr           = "http://127.0.0.1:8080"
	DefaultCluster        = "dev"
	DefaultIsBackupConfig = false
	DefaultNamespaces     = "application"
	DefaultSecret         = ""
)

func (a *apolloSource) Read() (set *source.ChangeSet, err error) {
	return read(a.namespaces, a.client)
}

func (a *apolloSource) Watch() (source.Watcher, error) {
	return newWatcher(a.namespaces, a.client)
}

// Write is unsupported
func (a *apolloSource) Write(cs *source.ChangeSet) error {
	return nil
}

func (a *apolloSource) String() string {
	return "apollo"
}

func read(ns []string, client *apo.Client) (set *source.ChangeSet, err error) {
	s := map[string]interface{}{}
	set = &source.ChangeSet{}
	for _, namespace := range ns {
		cache := client.GetConfigCache(namespace)
		cache.Range(func(key, value interface{}) bool {
			setValue(s, fmt.Sprintf("%v", value), strings.Split(fmt.Sprintf("%v", key), ".")...)
			return true
		})
	}

	set.Data, _ = json.Marshal(s)
	set.Checksum = set.Sum()
	set.Format = "json"
	set.Source = "file"

	if len(s) == 0 {
		err = fmt.Errorf("apollo data is nill, check the apollo error logs")
		log.Warn(err)
	}

	return
}

func setValue(input map[string]interface{}, v interface{}, keys ...string) {
	if len(keys) == 1 {
		input[keys[0]] = v
		return
	} else {
		var tmpMap map[string]interface{}
		if input[keys[0]] != nil {
			tmpMap = input[keys[0]].(map[string]interface{})
		} else {
			tmpMap = make(map[string]interface{})
		}

		input[keys[0]] = tmpMap
		setValue(tmpMap, v, keys[1:]...)
	}
}

func NewSource(opts ...source.Option) source.Source {
	var options source.Options
	for _, o := range opts {
		o(&options)
	}

	appID := "stack"
	addr := DefaultAddr
	cluster := DefaultCluster
	namespaces := DefaultNamespaces
	secret := DefaultSecret

	if options.Context != nil {
		appIDTemp, ok := options.Context.Value(appIDKey{}).(string)
		if !ok {
			log.Errorf("appId is necessary")
		} else {
			appID = appIDTemp
		}
		clusterTemp, ok := options.Context.Value(clusterKey{}).(string)
		if ok {
			cluster = clusterTemp
		} else if len(os.Getenv("APOLLO_CLUSTER")) > 0 {
			cluster = os.Getenv("APOLLO_CLUSTER")
		}

		addrTemp, ok := options.Context.Value(addrKey{}).(string)
		if ok {
			addr = addrTemp
		} else if len(os.Getenv("APOLLO_ADDRESS")) > 0 {
			addr = os.Getenv("APOLLO_ADDRESS")
		}

		namespaceTemp, ok := options.Context.Value(namespacesKey{}).(string)
		if ok {
			namespaces = namespaceTemp
		}

		secretTemp, ok := options.Context.Value(secretKey{}).(string)
		if ok {
			secret = secretTemp
		} else if len(os.Getenv("APOLLO_SECRET_KEY")) > 0 {
			secret = os.Getenv("APOLLO_SECRET_KEY")
		}
	}

	c := &apoC.AppConfig{
		AppID:          appID,
		Cluster:        cluster,
		IP:             addr,
		NamespaceName:  namespaces,
		IsBackupConfig: false,
		Secret:         secret,
	}

	client, err := apo.StartWithConfig(func() (*apoC.AppConfig, error) {
		return c, nil
	})

	if err != nil {
		log.Errorf("apollo client init error: %s", err)
	}

	return &apolloSource{
		appID:      appID,
		client:     client,
		namespaces: strings.Split(namespaces, ","),
		opts:       options,
	}
}
