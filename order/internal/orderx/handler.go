package orderx

import (
	"encoding/json"
	"net/http"
	"order/pkg/http/response"
)

func NewHandler(srv ServiceI) *Handler{
	return &Handler{
		srv: srv,
	}
}

func (h *Handler) InsertOrder(w http.ResponseWriter, r *http.Request) {
	var o OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	id, err := h.srv.InsertOrder(r.Context(), o)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{
		"id": *id,
	})
}
