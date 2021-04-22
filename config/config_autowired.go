package config

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/stack-labs/stack/util/log"
)

var (
	timeKind   = "time.Time"
	nullString = "null"
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

	go func() {
		for {
			select {
			// todo configurable, maybe
			case <-time.After(3 * time.Second):
				refresh()
			case data := <-ctx.Done():
				log.Infof("config autowired action stopped because of %v", data)
			}
		}
	}()
}

func bindAutowiredValue(obj reflect.Value, path ...string) {
	value := _sugar.Get(path...)
	v := reflect.Indirect(obj)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := int64(value.Int(0))
		if v.OverflowInt(n) {
			log.Errorf("bindAutowiredValue can't assign value due to %s-overflow", v.Kind())
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n := uint64(value.Int(0))
		if v.OverflowInt(int64(n)) {
			log.Errorf("bindAutowiredValue can't assign value due to %s-overflow", v.Kind())
		}
		v.SetUint(n)
	case reflect.String:
		valTmp := value.String("")
		if len(valTmp) == 0 {
			// maybe is other type
			valTmp = string(value.Bytes())
			if valTmp == nullString {
				valTmp = ""
			}
		}
		v.SetString(valTmp)
	case reflect.Bool:
		v.SetBool(value.Bool(false))
		// supports string only now
	case reflect.Slice, reflect.Array:
		values := value.StringSlice([]string{})
		v.Set(reflect.MakeSlice(reflect.SliceOf(v.Type().Elem()), len(values), len(values)))
		for idx, val := range values {
			nvalue := reflect.Indirect(reflect.New(v.Type().Elem()))
			nvalue.SetString(val)
			v.Index(idx).Set(nvalue)
		}
	case reflect.Struct:
		// Iterate over the struct fields
		fields := v.Type()

		// if it's time.Time type
		if fields.String() == timeKind {
			// supports RFC3339 only
			valTmp := value.String("")
			if len(valTmp) > 0 {
				t, err := time.Parse(time.RFC3339, valTmp)
				if err != nil {
					log.Warnf("%s is not standard RFC3339 format", valTmp)
					break
				}
				v.Set(reflect.ValueOf(t))
			}

			break
		}

		for i := 0; i < fields.NumField(); i++ {
			tag := fields.Field(i).Tag.Get(DefaultOptionsTagName)
			if tag == "" || tag == "-" {
				continue
			}

			nextValue := v.Field(i)
			newPath := append(path, tag)
			bindAutowiredValue(nextValue, newPath...)
		}
	default:
		log.Warnf("unsupported type: %s of %s", v.Kind().String(), v.String())
	}
}
