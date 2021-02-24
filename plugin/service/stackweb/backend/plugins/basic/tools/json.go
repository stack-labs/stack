package tools

import (
	"encoding/json"
	"net/http"

	"github.com/stack-labs/stack/plugin/service/stackweb/plugins"
)

func WriteJsonData(w http.ResponseWriter, data interface{}) {
	rsp := &plugins.Rsp{
		Data:    data,
		Success: true,
	}

	b, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func WriteError(w http.ResponseWriter, err error) {
	rsp := &plugins.Rsp{
		Error:   err.Error(),
		Success: false,
	}

	b, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
