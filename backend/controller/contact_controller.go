package controller

import (
	"net/http"
	"sftpbackend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ! 获取所有用户列表
func GetContactList(c *gin.Context) {
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

	// 获取查询参数中的用户名
	name := c.Query("name")
	var contact models.Contact
	if contacts, total, err := contact.GetContactList(page, limit, name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"contacts": contacts,
				"total":    total,
			},
		})
	}
}

// ! 添加一个联系人
func AddContact(c *gin.Context) {
	var contact models.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 查询通讯录是否已经存在
	if err := contact.ExistOrNot(contact.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "联系人已存在",
		})
		return
	} else {
		if err := contact.AddContact(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "添加联系人成功",
			})
			return
		}
	}
}

// ! 更新一个联系人
func UpdateContact(c *gin.Context) {
	var contact models.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 查询通讯录是否已经存在
	if err := contact.ExistOrNotByID(contact.ID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "联系人不存在",
		})
		return
	} else {
		if err := contact.UpdateContact(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "更新联系人成功",
			})
			return
		}
	}
}

// ! 删除一个联系人
func DeleteContact(c *gin.Context) {
	// 获取路径参数中的联系人ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的联系人ID参数",
		})
		return
	}

	var contact models.Contact
	// 查询通讯录是否已经存在
	if err := contact.ExistOrNotByID(uint(id)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "联系人不存在",
		})
		return
	} else {
		if err := contact.DeleteContact(uint(id)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "删除联系人成功",
			})
			return
		}
	}
}

// ! 批量删除联系人
func DeleteContacts(c *gin.Context) {
	var ids []uint
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	var contact models.Contact
	if err := contact.DeleteContacts(ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "批量删除联系人成功",
		})
		return
	}
}

// ! 获取所有联系人姓名和邮箱
func GetContactoptions(c *gin.Context) {
	var contact models.Contact
	if options, err := contact.GetContactoptions(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"options": options,
				"message": "获取选项成功",
			},
		})
	}
}
