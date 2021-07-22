package api

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/voxpupuli/webhook-go/parsers"
)

func (r Routes) addModule(rg *gin.RouterGroup) {
	module := rg.Group("/module")

	module.POST("/", deployModule)
}

func deployModule(c *gin.Context) {
	data := parsers.Data{}

	err := data.ParseData(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	cmd := exec.Command("r10k", "module", "deploy", data.ModuleName)

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(stdout))
}
