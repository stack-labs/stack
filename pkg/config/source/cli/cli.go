package cli

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
)

type cliSource struct {
	opts source.Options
	ctx  *cli.Context
}

func (c *cliSource) Read() (*source.ChangeSet, error) {
	changes := make(map[string]interface{})

	//for _, name := range c.ctx.GlobalFlagNames() {
	//	tmp := toEntry(name, c.ctx.GlobalGeneric(name))
	//	_ = mergo.Map(&changes, tmp) // TODO need to sort error handling
	//}
	//
	//for _, name := range c.ctx.FlagNames() {
	//	tmp := toEntry(name, c.ctx.Generic(name))
	//	_ = mergo.Map(&changes, tmp) // TODO need to sort error handling
	//}

	for _, name := range c.ctx.GlobalFlagNames() {
		if c.ctx.GlobalIsSet(name) {
			n := c.ctx.FlagAlias(name)
			v := c.ctx.GlobalGeneric(name)
			c.setValue(changes, v, strings.Split(n, "_")...)
		}
	}

	for _, name := range c.ctx.FlagNames() {
		if c.ctx.IsSet(name) {
			n := c.ctx.FlagAlias(name)
			v := c.ctx.Generic(name)
			c.setValue(changes, v, strings.Split(n, "_")...)
		}
	}

	b, err := c.opts.Encoder.Encode(changes)
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Format:    c.opts.Encoder.String(),
		Data:      b,
		Timestamp: time.Now(),
		Source:    c.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

//func toEntry(name string, v interface{}) map[string]interface{} {
//	n := strings.ToLower(name)
//	keys := strings.FieldsFunc(n, split)
//	reverse(keys)
//	tmp := make(map[string]interface{})
//	for i, k := range keys {
//		if i == 0 {
//			tmp[k] = v
//			continue
//		}
//
//		tmp = map[string]interface{}{k: tmp}
//	}
//	return tmp
//}

//func reverse(ss []string) {
//	for i := len(ss)/2 - 1; i >= 0; i-- {
//		opp := len(ss) - 1 - i
//		ss[i], ss[opp] = ss[opp], ss[i]
//	}
//}
//
//func split(r rune) bool {
//	return r == '-' || r == '_'
//}

func (c *cliSource) Watch() (source.Watcher, error) {
	return source.NewNoopWatcher()
}

func (c *cliSource) String() string {
	return "cli"
}

func (c *cliSource) setValue(input map[string]interface{}, v interface{}, keys ...string) {
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
		c.setValue(tmpMap, v, keys[1:]...)
	}
}

// NewSource returns a config source for integrating parsed flags from a stack/cli.Context.
// Hyphens are delimiters for nesting, and all keys are lowercased. The assumption is that
// command line flags have already been parsed.
//
// Example:
//      cli.StringFlag{Name: "db-host"},
//
//
//      {
//          "database": {
//              "host": "localhost"
//          }
//      }
func NewSource(app *cli.App, opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)

	var ctx *cli.Context

	c, ok := options.Context.Value(contextKey{}).(*cli.Context)
	if ok {
		ctx = c
	}

	// no context
	if ctx == nil {
		flags := app.Flags

		// create flagset
		set := flag.NewFlagSet(app.Name, flag.ContinueOnError)

		// apply flags to set
		for _, f := range flags {
			f.Apply(set)
		}

		// parse flags
		set.SetOutput(ioutil.Discard)
		_ = set.Parse(os.Args[1:])

		// normalise flags
		_ = normalizeFlags(app.Flags, set)

		// create context
		ctx = cli.NewContext(app, set, nil)
	}

	return &cliSource{
		ctx:  ctx,
		opts: options,
	}
}

// WithContext returns a new source with the context specified.
// The assumption is that Context is retrieved within an app.Action function.
func WithContext(ctx *cli.Context, opts ...source.Option) source.Source {
	return &cliSource{
		ctx:  ctx,
		opts: source.NewOptions(opts...),
	}
}
