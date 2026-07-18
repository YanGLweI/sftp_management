package controller

import (
	"net/http"
	"sftpbackend/models"

	"github.com/gin-gonic/gin"
)

type Dashboard struct {
	AccountCount    int64    `json:"accountCount"`    // 用户总数
	MonthlyNewCount int64    `json:"monthlyNewCount"` // 月新增用户数
	TransXaxis      []string `json:"transXaxis"`      // 传输量x轴数据
	TransFullDay    []string `json:"transFullDay"`    // 传输量7天数据
	AccessXaxis     []string `json:"accessXaxis"`     // 访问量x轴数据
	AccessFullDay   []string `json:"accessFullDay"`   // 访问量7天数据
}

func GetDashboard(c *gin.Context) {
	var u models.SftpUsers
	// 查询用户总数
	accountcount, err := u.GetTotalCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 查询月新增用户数
	monthlynewcount, err := u.GetMonthlyNewCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 获取传输统计数据
	transXaxis, transFullDay := models.GetTransCount()
	// 获取访问统计数据
	accessXaxis, accessFullDay := models.GetAccessCount()

	// 返回用户总数
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": Dashboard{
			AccountCount:    accountcount,
			MonthlyNewCount: monthlynewcount,
			TransXaxis:      transXaxis,
			TransFullDay:    transFullDay,
			AccessXaxis:     accessXaxis,
			AccessFullDay:   accessFullDay,
		},
		"message": "获取首页数据成功",
	})
}
