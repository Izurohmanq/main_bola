package routes

import (
	"order-service/clients"
	controllers "order-service/controllers/http"
	routes "order-service/routes/order"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	group      *gin.RouterGroup
	controller controllers.IControllerRegistry
	client     clients.IClientRegistry
}

type IRouteRegistry interface {
	Serve()
}

func NewRouteRegistry(group *gin.RouterGroup, controller controllers.IControllerRegistry, client clients.IClientRegistry) IRouteRegistry {
	return &Registry{group: group, controller: controller, client: client}
}

func (r *Registry) Serve() {
	r.orderRoute().Run()
}

func (r *Registry) orderRoute() routes.IOrderRoute {
	return routes.NewOrderRoute(r.group, r.controller, r.client)
}
