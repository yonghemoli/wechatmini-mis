package soulartifactapi

import (
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/soulartifact"

	"github.com/gin-gonic/gin"
)

func GetSoulArtifactAnalysis(c *gin.Context) {
	data, err := soulartifact.GetSoulArtifactAnalysis()
	if err != nil {
		response.Fail(c, "\u83b7\u53d6\u672c\u547d\u7269\u5206\u6790\u5931\u8d25: "+err.Error())
		return
	}
	response.OK(c, data)
}
