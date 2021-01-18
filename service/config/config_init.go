package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"

	cfg "github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
	cliSource "github.com/stack-labs/stack-rpc/pkg/config/source/cli"
	"github.com/stack-labs/stack-rpc/pkg/config/source/file"
	"github.com/stack-labs/stack-rpc/service"
	uf "github.com/stack-labs/stack-rpc/util/file"
	"github.com/stack-labs/stack-rpc/util/log"
)

func LoadConfig(sOpts *service.Options) (err error) {
	// set the config file path
	if filePath := sOpts.Cmd.App().Context().String("config"); len(filePath) > 0 {
		sOpts.Conf = filePath
	}

	// need to init the special config if specified
	if len(sOpts.Conf) == 0 {
		wkDir, errN := os.Getwd()
		if errN != nil {
			err = fmt.Errorf("stack can't access working wkDir: %s", errN)
			return
		}

		sOpts.Conf = fmt.Sprintf("%s%s%s", wkDir, string(os.PathSeparator), stackStdConfigFile)
	}

	var appendSource []source.Source
	var cfgOption []cfg.Option
	if len(sOpts.Conf) > 0 {
		// check file exists
		exists, err := uf.Exists(sOpts.Conf)
		if err != nil {
			log.Error(fmt.Errorf("config file is not existed %s", err))
		}

		if exists {
			// todo support more types
			val := struct {
				Stack struct {
					Includes string `yaml:"includes"`
					Config   Config `yaml:"config"`
				} `yaml:"stack"`
			}{}
			stdFileSource := file.NewSource(file.WithPath(sOpts.Conf))
			appendSource = append(appendSource, stdFileSource)

			set, errN := stdFileSource.Read()
			if errN != nil {
				err = fmt.Errorf("stack read the stack.yml err: %s", errN)
				return err
			}

			errN = yaml.Unmarshal(set.Data, &val)
			if errN != nil {
				err = fmt.Errorf("unmarshal stack.yml err: %s", errN)
				return err
			}

			if len(val.Stack.Includes) > 0 {
				filePath := sOpts.Conf[:strings.LastIndex(sOpts.Conf, string(os.PathSeparator))+1]
				for _, f := range strings.Split(val.Stack.Includes, ",") {
					log.Infof("load extra config file: %s%s", filePath, f)
					f = strings.TrimSpace(f)
					extraFile := fmt.Sprintf("%s%s", filePath, f)
					extraExists, err := uf.Exists(extraFile)
					if err != nil {
						log.Error(fmt.Errorf("config file is not existed %s", err))
						continue
					} else if !extraExists {
						log.Error(fmt.Errorf("config file [%s] is not existed", extraFile))
						continue
					}

					extraFileSource := file.NewSource(file.WithPath(extraFile))
					appendSource = append(appendSource, extraFileSource)
				}
			}

			// config option
			cfgOption = append(cfgOption, cfg.Storage(val.Stack.Config.Storage), cfg.HierarchyMerge(val.Stack.Config.HierarchyMerge))
		}
	}

	// the last two must be env & stackCmd line
	appendSource = append(appendSource, cliSource.NewSource(sOpts.Cmd.App(), cliSource.Context(sOpts.Cmd.App().Context())))
	cfgOption = append(cfgOption, cfg.Source(appendSource...))
	err = sOpts.Config.Init(cfgOption...)
	if err != nil {
		err = fmt.Errorf("init config err: %s", err)
		return
	}

	return
}

func SetOptions(sOpts *service.Options) (err error) {
	conf := stackConfig.Stack

	sOpts.ServerOptions = append(sOpts.ServerOptions, conf.Server.Options()...)
	sOpts.ClientOptions = append(sOpts.ClientOptions, conf.Client.Options()...)
	sOpts.ConfigOptions = append(sOpts.ConfigOptions, conf.Config.Options()...)
	sOpts.TransportOptions = append(sOpts.TransportOptions, conf.Transport.Options()...)
	sOpts.SelectorOptions = append(sOpts.SelectorOptions, conf.Selector.Options()...)
	sOpts.RegistryOptions = append(sOpts.RegistryOptions, conf.Registry.Options()...)
	sOpts.BrokerOptions = append(sOpts.BrokerOptions, conf.Broker.Options()...)
	sOpts.LoggerOptions = append(sOpts.LoggerOptions, conf.Logger.Options()...)
	sOpts.AuthOptions = append(sOpts.AuthOptions, conf.Auth.Options()...)

	return
}
