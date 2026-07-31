package data_api

import (
	"clover_server/common/res"
	"runtime"

	"github.com/gin-gonic/gin"
)

type ComputerDataResponse struct {
	CPUPercent float64 `json:"cpuPercent"`
	MemPercent float64 `json:"memPercent"`
	DiskPercent float64 `json:"diskPercent"`
	GoRoutines int     `json:"goRoutines"`
}

func (DataApi) ComputerDataView(c *gin.Context) {
	res.OkWithData(ComputerDataResponse{
		CPUPercent:  0,
		MemPercent:  0,
		DiskPercent: 0,
		GoRoutines:  runtime.NumGoroutine(),
	}, c)
}
