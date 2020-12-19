package config

import (
	"context"
	"github.com/stack-labs/stack-rpc/util/log"
	"reflect"
	"sync"
	"time"
)

func injectAutowired(ctx context.Context) {
	refresh := func() {
		var wg sync.WaitGroup
		for s, value := range optionsPool {
			wg.Add(1)

			go func(key string, val reflect.Value) {
				defer wg.Done()
				log.Debugf("setting values for %s", key)

				bindAutowiredValue(val)
			}(s, value)
		}
		wg.Wait()
	}

	// refresh for the first time
	refresh()
	for {
		select {
		// todo configurable, maybe
		case <-time.After(3 * time.Second):
			refresh()
		case data := <-ctx.Done():
			log.Infof("config autowired stop because of %v", data)
		}
	}
}

func bindAutowiredValue(val reflect.Value, path ...string) {
	log.Debugf("setting values for %s", val.Kind().String())
	configV := _sugar.Get(path...)
	v := reflect.Indirect(val)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(configV.Int(0)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(uint64(configV.Int(0)))
	case reflect.String:
		v.SetString(configV.String(""))
	case reflect.Bool:
		v.SetBool(configV.Bool(false))
	case reflect.Struct:
		// Iterate over the struct fields
		fields := v.Type()
		for i := 0; i < fields.NumField(); i++ {
			tag := fields.Field(i).Tag.Get(DefaultOptionsTagName)
			if tag == "" || tag == "-" {
				continue
			}
			value := v.Field(i)
			newPath := append(path, tag)
			bindAutowiredValue(value, newPath...)
		}
	default:
		log.Warnf("unsupported type: %s of %s", v.Kind().String(), v.String())
	}
}
