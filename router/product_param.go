package router

import (
	"goshop/api/controller"
	"goshop/api/pkg/core/routerhelper"
	"goshop/api/pkg/middleware"
	
	"github.com/gin-gonic/gin"
)

func init() {
	routerhelper.Use(func(r *gin.Engine) {
		g := routerhelper.NewGroupRouter("product-param", new(controller.ProductParam), r, middleware.Cors())
		g.Get("/index")
	})
}