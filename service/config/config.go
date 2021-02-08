package config

import (
	"fmt"
	"strings"
	"time"

	au "github.com/stack-labs/stack/auth"
	br "github.com/stack-labs/stack/broker"
	cl "github.com/stack-labs/stack/client"
	sel "github.com/stack-labs/stack/client/selector"
	cfg "github.com/stack-labs/stack/config"
	lg "github.com/stack-labs/stack/logger"
	"github.com/stack-labs/stack/plugin"
	reg "github.com/stack-labs/stack/registry"
	ser "github.com/stack-labs/stack/server"
	ss "github.com/stack-labs/stack/service"
	sw "github.com/stack-labs/stack/service/web"
	tra "github.com/stack-labs/stack/transport"
	"github.com/stack-labs/stack/util/log"
)

var (
	stackStdConfigFile = "stack.yml"
	stackConfig        = StackConfig{}
)

func init() {
	cfg.RegisterOptions(&stackConfig)
}

type Config struct {
	HierarchyMerge bool `json:"hierarchyMerge" sc:"hierarchy-merge"`
	Storage        bool `json:"storage" sc:"storage"`
}

func (c *Config) Options() []cfg.Option {
	var cfgOptions []cfg.Option

	cfgOptions = append(cfgOptions, cfg.HierarchyMerge(c.HierarchyMerge))
	cfgOptions = append(cfgOptions, cfg.Storage(c.Storage))

	return cfgOptions
}

type Broker struct {
	Address string `json:"address" sc:"address"`
	Name    string `json:"name" sc:"name"`
}

func (b *Broker) Options() []br.Option {
	var brOptions []br.Option

	if len(b.Name) > 0 {
		brOptions = append(brOptions, br.Name(b.Name))
	}

	if len(b.Address) > 0 {
		brOptions = append(brOptions, br.Addrs(strings.Split(b.Address, ",")...))
	}

	// todo adapt options by name

	return brOptions
}

type pool struct {
	Size int `json:"size" sc:"size"`
	TTL  int `json:"ttl" sc:"ttl"`
}

type clientRequest struct {
	Retries int    `json:"retries" sc:"retries"`
	Timeout string `json:"timeout" sc:"timeout"`
}

type Client struct {
	Name     string        `json:"name" sc:"name"`
	Protocol string        `json:"protocol" sc:"protocol"`
	Pool     pool          `json:"pool" sc:"pool"`
	Request  clientRequest `json:"request" sc:"request"`
}

func (c *Client) Options() []cl.Option {
	var cliOpts []cl.Option

	if len(c.Name) > 0 {
		cliOpts = append(cliOpts, cl.Name(c.Name))
	}

	if len(c.Protocol) > 0 {
		cliOpts = append(cliOpts, cl.Protocol(c.Protocol))
	}

	requestRetries := c.Request.Retries
	if requestRetries >= 0 {
		cliOpts = append(cliOpts, cl.Retries(requestRetries))
	}

	if len(c.Request.Timeout) > 0 {
		d, err := time.ParseDuration(c.Request.Timeout)
		if err != nil {
			err = fmt.Errorf("failed to parse client_request_timeout: %v. it shoud be with unit suffix such as 1s, 2m", c.Request.Timeout)
			log.Warn(err)
		} else {
			cliOpts = append(cliOpts, cl.RequestTimeout(d))
		}
	}

	if c.Pool.Size > 0 {
		cliOpts = append(cliOpts, cl.PoolSize(c.Pool.Size))
	}

	if poolTTL := time.Duration(c.Pool.TTL); poolTTL > 0 {
		cliOpts = append(cliOpts, cl.PoolTTL(poolTTL*time.Second))
	}

	return cliOpts
}

type Registry struct {
	Address string `json:"address" sc:"address"`
	Name    string `json:"name" sc:"name"`
}

func (r *Registry) Options() []reg.Option {
	var regOptions []reg.Option

	if len(r.Name) > 0 {
		regOptions = append(regOptions, reg.Name(r.Name))
	}

	if len(r.Address) > 0 {
		regOptions = append(regOptions, reg.Addrs(strings.Split(r.Address, ",")...))
	}

	if plugin.RegistryPlugins[r.Name] != nil {
		regOptions = append(regOptions, plugin.RegistryPlugins[r.Name].Options()...)
	}

	return regOptions
}

type metadata []string

func (m metadata) Value(k string) string {
	for _, s := range m {
		kv := strings.Split(s, "=")
		if len(kv) == 2 && kv[0] == k {
			return kv[1]
		}
	}

	return ""
}

type Server struct {
	Address     string         `json:"address" sc:"address"`
	Advertise   string         `json:"advertise" sc:"advertise"`
	ID          string         `json:"id" sc:"id"`
	Metadata    metadata       `json:"metadata" sc:"metadata"`
	Name        string         `json:"name" sc:"name"`
	Protocol    string         `json:"protocol" sc:"protocol"`
	Version     string         `json:"version" sc:"version"`
	Registry    serverRegistry `json:"Registry" sc:"Registry"`
	EnableDebug bool           `json:"enableDebug" sc:"enable-debug"`
}

type serverRegistry struct {
	TTL      int `json:"ttl" sc:"ttl"`
	Interval int `json:"interval" sc:"interval"`
}

func (s *Server) Options() []ser.Option {
	var serverOpts []ser.Option

	// Parse the server options
	metadata := make(map[string]string)
	for _, d := range s.Metadata {
		var key, val string
		parts := strings.Split(d, "=")
		key = parts[0]
		if len(parts) > 1 {
			val = strings.Join(parts[1:], "=")
		}
		metadata[key] = val
	}

	if len(metadata) > 0 {
		serverOpts = append(serverOpts, ser.Metadata(metadata))
	}

	if len(s.Protocol) > 0 {
		serverOpts = append(serverOpts, ser.Protocol(s.Protocol))
	}

	if len(s.Name) > 0 {
		serverOpts = append(serverOpts, ser.Name(s.Name))
	}

	if len(s.Version) > 0 {
		serverOpts = append(serverOpts, ser.Version(s.Version))
	}

	if len(s.ID) > 0 {
		serverOpts = append(serverOpts, ser.Id(s.ID))
	}

	if len(s.Address) > 0 {
		serverOpts = append(serverOpts, ser.Address(s.Address))
	}

	if len(s.Advertise) > 0 {
		serverOpts = append(serverOpts, ser.Advertise(s.Advertise))
	}

	if ttl := time.Duration(s.Registry.TTL); ttl >= 0 {
		serverOpts = append(serverOpts, ser.RegisterTTL(ttl*time.Second))
	}

	if val := time.Duration(s.Registry.Interval); val > 0 {
		serverOpts = append(serverOpts, ser.RegisterInterval(val*time.Second))
	}

	return serverOpts
}

type Selector struct {
	Name string `json:"name" sc:"name"`
}

func (s *Selector) Options() []sel.Option {
	var selOptions []sel.Option

	if len(s.Name) > 0 {
		selOptions = append(selOptions, sel.Name(s.Name))
	}

	if plugin.TransportPlugins[s.Name] != nil {
		selOptions = append(selOptions, plugin.SelectorPlugins[s.Name].Options()...)
	}

	return selOptions
}

type Transport struct {
	Name    string `json:"name" sc:"name"`
	Address string `json:"address" sc:"address"`
}

func (t *Transport) Options() []tra.Option {
	var traOptions []tra.Option

	if len(t.Name) > 0 {
		traOptions = append(traOptions, tra.Name(t.Name))
	}

	if len(t.Address) > 0 {
		traOptions = append(traOptions, tra.Addrs(strings.Split(t.Address, ",")...))
	}

	if plugin.TransportPlugins[t.Name] != nil {
		traOptions = append(traOptions, plugin.TransportPlugins[t.Name].Options()...)
	}

	return traOptions
}

type Logger struct {
	Name  string `json:"name" sc:"name"`
	Level string `json:"level" sc:"level"`
	// todo support map settings
	// Fields          map[string]string `json:"fields" sc:"fields"`
	CallerSkipCount int            `json:"caller-skip-count" sc:"caller-skip-count"`
	Persistence     logPersistence `json:"persistence" sc:"persistence"`
}

type logPersistence struct {
	Enable    bool   `json:"enable" sc:"enable"`
	Dir       string `json:"dir" sc:"dir"`
	BackupDir string `json:"backupDir" sc:"back-dir"`
	// log file max size in megabytes
	MaxFileSize int `json:"maxFileSize" sc:"max-file-size"`
	// backup dir max size in megabytes
	MaxBackupSize int `json:"maxBackupSize" sc:"max-backup-size"`
	// backup files keep max days
	MaxBackupKeepDays int `json:"maxBackupKeepDays" sc:"max-backup-keep-days"`
	// default pattern is ${serviceName}_${level}.log
	// todo available patterns map
	FileNamePattern string `json:"fileNamePattern" sc:"file-name-pattern"`
	// default pattern is ${serviceName}_${level}_${yyyyMMdd_HH}_${idx}.zip
	// todo available patterns map
	BackupFileNamePattern string `json:"backupFileNamePattern" sc:"backup-file-name-pattern"`
}

func (l *logPersistence) Options() *lg.PersistenceOptions {
	o := &lg.PersistenceOptions{
		Enable:                l.Enable,
		Dir:                   l.Dir,
		BackupDir:             l.BackupDir,
		MaxFileSize:           l.MaxFileSize,
		MaxBackupSize:         l.MaxBackupSize,
		MaxBackupKeepDays:     l.MaxBackupKeepDays,
		FileNamePattern:       l.FileNamePattern,
		BackupFileNamePattern: l.BackupFileNamePattern,
	}

	return o
}

func (l *Logger) Options() []lg.Option {
	var logOptions []lg.Option

	if len(l.Name) > 0 {
		logOptions = append(logOptions, lg.Name(l.Name))
	}

	if len(l.Level) > 0 {
		level, err := lg.GetLevel(l.Level)
		if err != nil {
			err = fmt.Errorf("ilegal logger level error: %s", err)
			log.Warn(err)
		} else {
			logOptions = append(logOptions, lg.WithLevel(level))
		}
	}

	if l.Persistence.Enable {
		logOptions = append(logOptions, lg.Persistence(l.Persistence.Options()))
	}

	if plugin.LoggerPlugins[l.Name] != nil {
		logOptions = append(logOptions, plugin.LoggerPlugins[l.Name].Options()...)
	} else if len(l.Name) > 0 {
		log.Warnf("seems you declared a logger name:[%s] which stack can't find out.", l.Name)
	}

	return logOptions
}

type Auth struct {
	Name            string          `json:"name" sc:"name"`
	Enable          bool            `json:"enable" sc:"enable"`
	Namespace       string          `json:"namespace" sc:"namespace"`
	AuthCredentials authCredentials `json:"authCredentials" sc:"authCredentials"`
	PublicKey       string          `json:"publicKey" sc:"public-key"`
	PrivateKey      string          `json:"privateKey" sc:"private-key"`
}

type authCredentials struct {
	ID     string `json:"id" sc:"id"`
	Secret string `json:"secret" sc:"secret"`
}

func (a *Auth) Options() []au.Option {
	var opts []au.Option

	opts = append(opts, au.Enable(a.Enable))
	opts = append(opts, au.Namespace(a.Namespace))

	if len(a.AuthCredentials.ID) > 0 {
		opts = append(opts, au.Credentials(a.AuthCredentials.ID, a.AuthCredentials.Secret))
	}

	opts = append(opts, au.PublicKey(a.PublicKey))
	opts = append(opts, au.PrivateKey(a.PrivateKey))

	if plugin.LoggerPlugins[a.Name] != nil {
		opts = append(opts, plugin.AuthPlugins[a.Name].Options()...)
	} else if len(a.Name) > 0 {
		log.Warnf("seems you declared an auth name:[%s] which stack can't find out.", a.Name)
	}

	return opts
}

type Web struct {
	Enable   bool   `json:"enable" sc:"enable"`
	RootPath string `json:"rootPath" sc:"root-path"`
	Static   struct {
		Route string `json:"route" sc:"route"`
		Dir   string `json:"dir" sc:"dir"`
	} `json:"static" sc:"static"`
}

type serviceOpts []ss.Option

func (s serviceOpts) opts() ss.Options {
	opts := ss.Options{}
	for _, o := range s {
		o(&opts)
	}

	return opts
}

type Service struct {
	ID      string `json:"id" sc:"id"`
	Name    string `json:"name" sc:"name"`
	Address string `json:"address" sc:"address"`
	RPC     string `json:"rpc" sc:"rpc"`
	Web     Web    `json:"web" sc:"web"`
}

func (s *Service) Options() serviceOpts {
	var opts serviceOpts

	if len(s.ID) > 0 {
		opts = append(opts, ss.Id(s.ID))
	}

	if len(s.Name) > 0 {
		opts = append(opts, ss.Name(s.Name))
	}

	opts = append(opts, sw.Enable(s.Web.Enable))

	if len(s.Web.RootPath) > 0 {
		opts = append(opts, sw.RootPath(s.Web.RootPath))
	}

	if len(s.Web.Static.Dir) > 0 {
		opts = append(opts, sw.StaticDir(s.Web.Static.Route, s.Web.Static.Dir))
	}

	return opts
}

type StackConfig struct {
	Stack struct {
		Includes  string    `json:"includes" sc:"includes"`
		Config    Config    `json:"config" sc:"config"`
		Registry  Registry  `json:"registry" sc:"registry"`
		Broker    Broker    `json:"broker" sc:"broker"`
		Client    Client    `json:"client" sc:"client"`
		Profile   string    `json:"profile" sc:"profile"`
		Runtime   string    `json:"runtime" sc:"runtime"`
		Server    Server    `json:"server" sc:"server"`
		Selector  Selector  `json:"selector" sc:"selector"`
		Transport Transport `json:"transport" sc:"transport"`
		Logger    Logger    `json:"logger" sc:"logger"`
		Auth      Auth      `json:"auth" sc:"auth"`
		Service   Service   `json:"service" sc:"service"`
	} `json:"stack" sc:"stack"`
}
