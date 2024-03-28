package routes

import (
	controller "RESTAURANT-MANAGEMENT/controllers"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/invoices", controller.GetInvoices())
	incomingRoutes.GET("/invoices/:invoices_id", controller.GetInvoices())
	incomingRoutes.POST("/invoices", controller.CreateInvoices())
	incomingRoutes.PATCH("/invoices/:invoices_id", controller.UpdateInvoices())

}
