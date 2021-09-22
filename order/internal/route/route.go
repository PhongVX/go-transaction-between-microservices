package route

import (
	"github.com/PhongVX/micro-protos/order"
	"github.com/PhongVX/micro-protos/transaction"
	"net/http"

	"github.com/gorilla/mux"
	"order/pkg/http/middleware"
	"order/pkg/http/router"
)

func NewRouter(orderC order.OrderClient, txC transaction.TransactionClient) (http.Handler, error) {
	r := mux.NewRouter()
	orderXHandler, err := newOrderXHandler(orderC, txC)
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

