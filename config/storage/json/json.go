package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/stack-labs/stack-rpc/logger"
)

type jsonStorage struct {
	dir string
}

func (j *jsonStorage) Write(namespace string, config []byte) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("json storage write config to file err: %s", err)
		}
	}()

	if config == nil {
		err = fmt.Errorf("json storage write nil config")
		return
	}

	file, err := os.Create(j.fileName(namespace))
	if err != nil && err != os.ErrExist {
		err = fmt.Errorf("json storage write nil config")
		return
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(config)

	return
}

func (j *jsonStorage) Load(namespace string) (config []byte, err error) {
	defer func() {
		if err != nil {
			log.Errorf("json storage load config from file err: %s", err)
		}
	}()

	fs, err := ioutil.ReadFile(j.fileName(namespace))
	if err != nil {
		err = fmt.Errorf("read file error: %s", err)
		return
	}

	err = json.Unmarshal(fs, &config)
	if err != nil {
		err = fmt.Errorf("unmarshal file error: %s", err)
		return
	}

	return config, nil
}

func (j *jsonStorage) fileName(namespace string) string {
	return fmt.Sprintf("%s_%s.json", j.dir, namespace)
}
