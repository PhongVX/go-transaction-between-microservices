package orderx

import (
	"net/http"
	"order/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:    "/api/v1/order",
			Method:  http.MethodPost,
			Handler: h.InsertOrder,
		},
	}
}
