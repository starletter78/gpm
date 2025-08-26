package file

import (
	"gpm/app/service/search"
	"gpm/common/res"
	"gpm/global"

	"github.com/gin-gonic/gin"
)

func (File) FileTreeView(c *gin.Context) {
	filePath := global.Config.Log.Dir
	tree, err := search.BuildFileTree(filePath)
	if err != nil {
		res.FailWithError(c, err)
		return
	}
	res.SuccessWithData(c, tree)
}
