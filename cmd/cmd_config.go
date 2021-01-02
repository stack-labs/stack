package cmd

import (
	"fmt"
	"strings"
	"time"

	au "github.com/stack-labs/stack-rpc/auth"
	br "github.com/stack-labs/stack-rpc/broker"
	cl "github.com/stack-labs/stack-rpc/client"
	cfg "github.com/stack-labs/stack-rpc/config"
	lg "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/plugin"
	reg "github.com/stack-labs/stack-rpc/registry"
	ser "github.com/stack-labs/stack-rpc/server"
	tra "github.com/stack-labs/stack-rpc/transport"
	"github.com/stack-labs/stack-rpc/util/log"
	wb "github.com/stack-labs/stack-rpc/web"
)

type stack struct {
	Includes  string    `json:"includes" sc:"includes"`
	Config    config    `json:"config" sc:"config"`
	Registry  registry  `json:"registry" sc:"registry"`
	Broker    broker    `json:"broker" sc:"broker"`
	Client    client    `json:"client" sc:"client"`
	Profile   string    `json:"profile" sc:"profile"`
	Runtime   string    `json:"runtime" sc:"runtime"`
	Server    server    `json:"server" sc:"server"`
	Selector  selector  `json:"selector" sc:"selector"`
	Transport transport `json:"transport" sc:"transport"`
	Logger    logger    `json:"logger" sc:"logger"`
	Auth      auth      `json:"auth" sc:"auth"`
	Web       web       `json:"web" sc:"web"`
}

type config struct {
	HierarchyMerge bool `json:"hierarchyMerge" sc:"hierarchy-merge"`
	Storage        bool `json:"storage" sc:"storage"`
}

func (c *config) Options() []cfg.Option {
	var cfgOptions []cfg.Option

	cfgOptions = append(cfgOptions, cfg.HierarchyMerge(c.HierarchyMerge))
	cfgOptions = append(cfgOptions, cfg.Storage(c.Storage))

	return cfgOptions
}

type broker struct {
	Address string `json:"address" sc:"address"`
	Name    string `json:"name" sc:"name"`
}

func (b *broker) Options() []br.Option {
	var brOptions []br.Option

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

type client struct {
	Protocol string        `json:"protocol" sc:"protocol"`
	Pool     pool          `json:"pool" sc:"pool"`
	Request  clientRequest `json:"request" sc:"request"`
}

func (c *client) Options() []cl.Option {
	var cliOpts []cl.Option

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

type registry struct {
	Address string `json:"address" sc:"address"`
	Name    string `json:"name" sc:"name"`
}

func (r *registry) Options() []reg.Option {
	var regOptions []reg.Option

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

type server struct {
	Address   string         `json:"address" sc:"address"`
	Advertise string         `json:"advertise" sc:"advertise"`
	ID        string         `json:"id" sc:"id"`
	Metadata  metadata       `json:"metadata" sc:"metadata"`
	Name      string         `json:"name" sc:"name"`
	Protocol  string         `json:"protocol" sc:"protocol"`
	Version   string         `json:"version" sc:"version"`
	Registry  serverRegistry `json:"registry" sc:"registry"`
}

type serverRegistry struct {
	TTL      int `json:"ttl" sc:"ttl"`
	Interval int `json:"interval" sc:"interval"`
}

func (s *server) Options() []ser.Option {
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

type selector struct {
	Name string `json:"name" sc:"name"`
}

type transport struct {
	Name    string `json:"name" sc:"name"`
	Address string `json:"address" sc:"address"`
}

func (t *transport) Options() []tra.Option {
	var traOptions []tra.Option
	if len(t.Address) > 0 {
		traOptions = append(traOptions, tra.Addrs(strings.Split(t.Address, ",")...))
	}

	if plugin.TransportPlugins[t.Name] != nil {
		traOptions = append(traOptions, plugin.TransportPlugins[t.Name].Options()...)
	}

	return traOptions
}

type logger struct {
	Name  string `json:"name" sc:"name"`
	Level string `json:"level" sc:"level"`
	// todo support map settings
	// Fields          map[string]string `json:"fields" sc:"fields"`
	CallerSkipCount int            `json:"caller-skip-count" sc:"caller-skip-count"`
	Persistence     logPersistence `json:"persistence" sc:"persistence"`
}

type logPersistence struct {
	Enable    bool   `sc:"enable"`
	Dir       string `sc:"dir"`
	BackupDir string `sc:"back-dir"`
	// log file max size in megabytes
	MaxFileSize int `sc:"max-file-size"`
	// backup dir max size in megabytes
	MaxBackupSize int `sc:"max-backup-size"`
	// backup files keep max days
	MaxBackupKeepDays int `sc:"max-backup-keep-days"`
	// default pattern is ${serviceName}_${level}.log
	// todo available patterns map
	FileNamePattern string `sc:"file-name-pattern"`
	// default pattern is ${serviceName}_${level}_${yyyyMMdd_HH}_${idx}.zip
	// todo available patterns map
	BackupFileNamePattern string `sc:"backup-file-name-pattern"`
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

func (l *logger) Options() []lg.Option {
	var logOptions []lg.Option

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

type auth struct {
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

func (a *auth) Options() []au.Option {
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

type web struct {
	Id        string            `json:"id" sc:"id"`
	Name      string            `json:"name" sc:"name"`
	Metadata  map[string]string `json:"metadata" sc:"metadata"`
	Address   string            `json:"address" sc:"address"`
	Version   string            `json:"version" sc:"version"`
	Advertise string            `json:"advertise" sc:"advertise"`
	Registry  serverRegistry    `json:"registry" sc:"registry"`
	Secure    bool              `json:"secure" sc:"secure"`
}

func (w *web) Options() []wb.Option {
	var opts []wb.Option

	if len(w.Id) != 0 {
		opts = append(opts, wb.Version(w.Id))
	}

	if len(w.Name) != 0 {
		opts = append(opts, wb.Name(w.Name))
	}

	if len(w.Version) != 0 {
		opts = append(opts, wb.Version(w.Version))
	}

	if len(w.Metadata) != 0 {
		// todo merge?
		opts = append(opts, wb.Metadata(w.Metadata))
	}

	if ttl := time.Duration(w.Registry.TTL); ttl >= 0 {
		opts = append(opts, wb.RegisterTTL(ttl*time.Second))
	}

	if val := time.Duration(w.Registry.Interval); val > 0 {
		opts = append(opts, wb.RegisterInterval(val*time.Second))
	}

	opts = append(opts, wb.Secure(w.Secure))

	return opts
}

type StackConfig struct {
	Stack stack `json:"stack" sc:"stack"`
}
