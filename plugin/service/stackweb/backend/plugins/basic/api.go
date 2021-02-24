package basic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/stack-labs/stack/client/selector"
	"github.com/stack-labs/stack/pkg/metadata"
	"github.com/stack-labs/stack/plugin/service/stackweb/plugins/basic/tools"
	"github.com/stack-labs/stack/registry"
	"github.com/stack-labs/stack/service"
)

type api struct {
	service service.Options
}

type rpcRequest struct {
	Service  string
	Endpoint string
	Method   string
	Address  string
	URL      string
	timeout  int
	Request  interface{}
}

// serviceAPIDetail is the service api detail
type serviceAPIDetail struct {
	Name      string               `json:"name,omitempty"`
	Endpoints []*registry.Endpoint `json:"endpoints,omitempty"`
}

func (a *api) webServices(w http.ResponseWriter, r *http.Request) {
	services, err := a.service.Registry.ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	webServices := make([]*registry.Service, 0)
	for _, s := range services {
		for _, webN := range WebNamespacePrefix {
			if strings.Index(s.Name, webN) == 0 && len(strings.TrimPrefix(s.Name, webN)) > 0 {
				s.Name = strings.Replace(s.Name, webN+".", "", 1)
				webServices = append(webServices, s)
			}
		}
	}

	sort.Sort(tools.SortedServices{Services: services})

	tools.WriteJsonData(w, webServices)

	return
}

func (a *api) services(w http.ResponseWriter, r *http.Request) {
	services, err := a.service.Registry.ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	for _, service := range services {
		ss, err := a.service.Registry.GetService(service.Name)
		if err != nil {
			continue
		}
		if len(ss) == 0 {
			continue
		}

		for _, s := range ss {
			service.Nodes = append(service.Nodes, s.Nodes...)
			service.Endpoints = s.Endpoints
		}

	}

	sort.Sort(tools.SortedServices{Services: services})

	tools.WriteJsonData(w, services)
	return
}

func (a *api) microServices(w http.ResponseWriter, r *http.Request) {
	services, err := a.service.Registry.ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	ret := make([]*registry.Service, 0)

	for _, srv := range services {
		temp, err := a.service.Registry.GetService(srv.Name)
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}

		for _, s := range temp {
			for _, n := range s.Nodes {
				if n.Metadata["registry"] != "" {
					ret = append(ret, s)
					break
				}
			}
		}
	}

	sort.Sort(tools.SortedServices{Services: ret})

	tools.WriteJsonData(w, ret)
	return
}

func (a *api) serviceDetails(w http.ResponseWriter, r *http.Request) {
	services, err := a.service.Registry.ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	sort.Sort(tools.SortedServices{Services: services})

	serviceDetails := make([]*serviceAPIDetail, 0)
	for _, service := range services {
		s, err := a.service.Registry.GetService(service.Name)
		if err != nil {
			continue
		}
		if len(s) == 0 {
			continue
		}

		serviceDetails = append(serviceDetails, &serviceAPIDetail{
			Name:      service.Name,
			Endpoints: s[0].Endpoints,
		})
	}

	tools.WriteJsonData(w, serviceDetails)
	return
}

func (a *api) handler(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("service")

	if len(serviceName) > 0 {
		s, err := a.service.Registry.GetService(serviceName)
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}

		if len(s) == 0 {
			tools.WriteError(w, fmt.Errorf("Service Is Not found %s: ", serviceName))
			return
		}

		tools.WriteJsonData(w, s)
		return
	}

	return
}

func (a *api) apiGatewayServices(w http.ResponseWriter, r *http.Request) {
	services, err := a.service.Registry.ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	ret := make([]*registry.Service, 0)
	for _, service := range services {
		_, _ = a.service.Selector.Next(service.Name, func(options *selector.SelectOptions) {
			filter := func(services []*registry.Service) []*registry.Service {
				for _, s := range services {
					for _, gwN := range GatewayNamespaces {
						if s.Name == gwN {
							ret = append(ret, s)
							break
						}
					}
				}
				return ret
			}

			options.Filters = append(options.Filters, filter)
		})
	}

	tools.WriteJsonData(w, ret)
	return
}

func (a *api) rpc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rpcReq := &rpcRequest{}

	d := json.NewDecoder(r.Body)
	d.UseNumber()

	if err := d.Decode(&rpcReq); err != nil {
		tools.WriteError(w, fmt.Errorf("rpc decode err %s: ", err))
		return
	}

	if len(rpcReq.Endpoint) == 0 {
		rpcReq.Endpoint = rpcReq.Method
	}

	rpcReq.timeout, _ = strconv.Atoi(r.Header.Get("Timeout"))
	rpcReq.URL = r.URL.Path

	a.rpcCall(w, requestToContext(r), rpcReq)
}

func (a *api) health(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	rpcReq := &rpcRequest{
		Service:  r.URL.Query().Get("service"),
		Endpoint: "Debug.Health",
		Request:  "{}",
		URL:      r.URL.Path,
		Address:  r.URL.Query().Get("address"),
	}

	a.rpcCall(w, requestToContext(r), rpcReq)
}

func (a *api) stats(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	rpcReq := &rpcRequest{
		Service:  r.URL.Query().Get("service"),
		Endpoint: "Debug.Stats",
		Request:  "{}",
		URL:      r.URL.Path,
		Address:  r.URL.Query().Get("address"),
	}

	a.rpcCall(w, requestToContext(r), rpcReq)
	return
}

func requestToContext(r *http.Request) context.Context {
	ctx := context.Background()
	md := make(metadata.Metadata)
	for k, v := range r.Header {
		md[k] = strings.Join(v, ",")
	}
	return metadata.NewContext(ctx, md)
}
