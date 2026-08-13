package controller

import (
	"net/http"
	"sftpbackend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ! GetSftpLogList 获取SFTP登录与操作日志列表
func GetSftpLogList(c *gin.Context) {
	// 获取路径参数中的页码和每页记录数
	pageStr := c.Param("page")
	limitStr := c.Param("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的页码参数",
		})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的每页记录数参数",
		})
		return
	}
	// 获取查询参数中的时间与用户名
	date := c.Query("datetime")
	username := c.Query("username")
	var log models.SftpLog
	if logs, total, err := log.GetSftpLogList(page, limit, date, username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"logs":  logs,
				"total": total,
			},
		})
	}
}
