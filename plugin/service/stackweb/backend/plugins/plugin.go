package plugins

import (
	"net/http"
)

// Plugin is the biggest unit of web component.
// a component is a webapp which will be registered on root path by Name function
type Plugin interface {
	// Name of module
	Name() string

	// Path returns the root path of this module
	Path() string

	// Init initializes the module
	Init(...Option) error

	// Handlers returns http handler of this module
	Handlers() map[string]*Handler
}

type Handler struct {
	Name   string
	Func   func(w http.ResponseWriter, r *http.Request)
	Method []string
	Hld    http.Handler
}

func (h Handler) IsFunc() bool {
	return h.Func != nil
}

// Rsp is the struct of http api response
type Rsp struct {
	Code    uint        `json:"code,omitempty"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

var (
	modules []Plugin
)

func Register(m Plugin) {
	modules = append(modules, m)
}

// Plugins returns all of the registered modules
func Plugins() []Plugin {
	return modules
}
