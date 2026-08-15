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

// ============ P0: 核心数据统计 API ============

// GetTotalAccessCount 获取累计访问量（登录总次数）
func GetTotalAccessCount(c *gin.Context) {
	count, err := models.GetTotalLoginCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    count,
		"message": "获取累计访问量成功",
	})
}

// GetTotalTransferCount 获取累计传输数（上传 + 下载总次数）
func GetTotalTransferCount(c *gin.Context) {
	count, err := models.GetTotalTransferCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    count,
		"message": "获取累计传输数成功",
	})
}

// GetTodayAccessCount 获取今日访问量（登录次数）
func GetTodayAccessCount(c *gin.Context) {
	count, err := models.GetTodayLoginCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    count,
		"message": "获取今日访问量成功",
	})
}

// GetAccessGrowth 获取访问量增长率（今日 vs 昨日）
type AccessGrowth struct {
	Today      int64   `json:"today"`      // 今日访问量
	Yesterday  int64   `json:"yesterday"`  // 昨日访问量
	GrowthRate float64 `json:"growthRate"` // 同比增长率（百分比数值，如 12.5 表示 +12.5%）
}

func GetAccessGrowth(c *gin.Context) {
	today, err := models.GetTodayLoginCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	yesterday, err := models.GetYesterdayLoginCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	// 计算增长率：昨日为 0 时返回 100（视为全部新增），否则为百分比
	var growthRate float64
	if yesterday == 0 {
		if today > 0 {
			growthRate = 100
		} else {
			growthRate = 0
		}
	} else {
		growthRate = (float64(today) - float64(yesterday)) / float64(yesterday) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": AccessGrowth{
			Today:      today,
			Yesterday:  yesterday,
			GrowthRate: growthRate,
		},
		"message": "获取访问量增长率成功",
	})
}

// GetTodayTransferCount 获取今日传输总量（上传 + 下载次数）
func GetTodayTransferCount(c *gin.Context) {
	total, err := models.GetTodayTransferCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    total,
		"message": "获取今日传输量成功",
	})
}

// GetTransferGrowth 获取传输量增长率（今日 vs 昨日）
type TransferGrowth struct {
	Today      int64   `json:"today"`      // 今日传输量
	Yesterday  int64   `json:"yesterday"`  // 昨日传输量
	GrowthRate float64 `json:"growthRate"` // 同比增长率（百分比数值）
}

func GetTransferGrowth(c *gin.Context) {
	today, err := models.GetTodayTransferCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	yesterday, err := models.GetYesterdayTransferCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	var growthRate float64
	if yesterday == 0 {
		if today > 0 {
			growthRate = 100
		} else {
			growthRate = 0
		}
	} else {
		growthRate = (float64(today) - float64(yesterday)) / float64(yesterday) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": TransferGrowth{
			Today:      today,
			Yesterday:  yesterday,
			GrowthRate: growthRate,
		},
		"message": "获取传输量增长率成功",
	})
}

// GetAuthDistribution 获取认证方式分布（密钥/密码）
func GetAuthDistribution(c *gin.Context) {
	data, err := models.GetAuthMethodDistribution()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    data,
		"message": "获取认证分布成功",
	})
}

// GetActiveUsersTop6 获取活跃用户 Top6
func GetActiveUsersTop6(c *gin.Context) {
	users, err := models.GetActiveUsersTop6()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    users,
		"message": "获取活跃用户成功",
	})
}

// GetTopTransferUsers 获取传输量排行 Top10
func GetTopTransferUsers(c *gin.Context) {
	users, err := models.GetTopTransferUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    users,
		"message": "获取传输量排行成功",
	})
}
