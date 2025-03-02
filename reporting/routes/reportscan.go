package routes

import (
	"github.com/gin-gonic/gin"
)

type Host struct {
	IpAddress  string
	MacAddress string
	Dynamic    bool
}

func ReportScan(c *gin.Context) {
	var hosts []Host
	err := c.BindJSON(&hosts)
	if err != nil {
		// TODO: Handle errors
	}
}
