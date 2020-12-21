package cmd

import "strings"

type stack struct {
	Registry  registry  `json:"registry" sc:"registry"`
	Broker    broker    `json:"broker" sc:"broker"`
	Client    client    `json:"client" sc:"client"`
	Profile   string    `json:"profile" sc:"profile"`
	Runtime   string    `json:"runtime" sc:"runtime"`
	Server    server    `json:"server" sc:"server"`
	Selector  selector  `json:"selector" sc:"selector"`
	Transport transport `json:"transport" sc:"transport"`
	Logger    logger    `json:"logger" sc:"logger"`
}

type broker struct {
	Address string `json:"address" sc:"address"`
	Name    string `json:"name" sc:"name"`
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

type registry struct {
	Address  string `json:"address" sc:"address"`
	Interval int    `json:"interval" sc:"interval"`
	Name     string `json:"name" sc:"name"`
	TTL      int    `json:"ttl" sc:"ttl"`
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
	Address   string   `json:"address" sc:"address"`
	Advertise string   `json:"advertise" sc:"advertise"`
	ID        string   `json:"id" sc:"id"`
	Metadata  metadata `json:"metadata" sc:"metadata"`
	Name      string   `json:"name" sc:"name"`
	Protocol  string   `json:"protocol" sc:"protocol"`
	Version   string   `json:"version" sc:"version"`
}

type selector struct {
	Name string `json:"name" sc:"name"`
}

type transport struct {
	Name    string `json:"name" sc:"name"`
	Address string `json:"address" sc:"address"`
}

type logger struct {
	Name  string `json:"name" sc:"name"`
	Level string `json:"level" sc:"level"`
}

type Value struct {
	Stack stack `json:"stack" sc:"stack"`
}
