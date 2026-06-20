package socialnetwork

import (
	"yonghemolimis/src/apps/api/response"
	uc "yonghemolimis/src/usecase/socialnetwork"

	"github.com/gin-gonic/gin"
)

func GetSocialNetwork(c *gin.Context) {
	data, err := uc.GetSocialNetworkData()
	if err != nil {
		response.Fail(c, "获取社交网络数据失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
