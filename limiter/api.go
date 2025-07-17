package limiter

import (
	"encoding/json"
	"io"
	"net/http"
)

type API struct {
	NodeMgr   *NodeManager
	Limiter   *Limiter
	GlobalQPS float64
}

// UpdateNodesHandler allows external systems to update the node list via HTTP POST
func (api *API) UpdateNodesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var nodes []Node
	if err := json.Unmarshal(body, &nodes); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	api.NodeMgr.UpdateNodes(nodes)
	newRate := api.NodeMgr.CalcSelfRate(api.GlobalQPS)
	api.Limiter.UpdateRate(newRate)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
