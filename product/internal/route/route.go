package route

import (
	"github.com/PhongVX/micro-protos/product"
	"github.com/PhongVX/micro-protos/transaction"
	"net/http"

	"github.com/gorilla/mux"
	"product/pkg/http/middleware"
	"product/pkg/http/router"
)

func NewRouter(productC product.ProductClient, txC transaction.TransactionClient) (http.Handler, error) {
	r := mux.NewRouter()
	orderXHandler, err := newOrderXHandler(productC, txC)
	if err != nil {
		return nil, err
	}

	routes := []router.Route{}
	routes = append(routes, orderXHandler.Routes()...)

	//Routes
	for _, rt := range routes {
		var h http.Handler
		h = rt.Handler
		for i := len(rt.Middlewares) - 1; i >= 0; i-- {
			h = rt.Middlewares[i](h)
		}
		r.Path(rt.Path).Methods(rt.Method).Handler(h)
	}

	return middleware.CORS(r), nil
}

