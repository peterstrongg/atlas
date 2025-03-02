package main

import (
	"atlas/reporting/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/api/v1/reportscan", routes.ReportScan)
	r.Run()
}
