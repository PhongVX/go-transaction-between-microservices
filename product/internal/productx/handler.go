package productx

import (
	"encoding/json"
	"net/http"
	"product/pkg/http/response"
)

func NewHandler(srv ServiceI) *Handler{
	return &Handler{
		srv: srv,
	}
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var o UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	_, err := h.srv.UpdateProduct(r.Context(), o)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, map[string]*int{
		"id": o.Body.ID,
	})
}
