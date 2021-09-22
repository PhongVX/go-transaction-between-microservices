package productx

import (
	"net/http"
	"product/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:    "/api/v1/product",
			Method:  http.MethodPut,
			Handler: h.UpdateProduct,
		},
	}
}
