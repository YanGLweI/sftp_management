package controller

import (
	"fmt"
	"net/http"
	"sftpbackend/dao"
	"sftpbackend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 超级管理员角色名（不可删除、不可修改菜单权限）
const superAdminRoleName = "超级管理员"

// GetAllMenus 返回所有预定义的菜单树（前端路由表）
func GetAllMenus(c *gin.Context) {
	menus := []gin.H{
		{"routeName": "Dashboard", "menuTitle": "首页", "icon": "dashboard"},
		{"routeName": "Sftp", "menuTitle": "传输管理", "icon": "el-icon-sort", "children": []gin.H{
			{"routeName": "SftpUser", "menuTitle": "账号管理"},
			{"routeName": "Contacts", "menuTitle": "通讯邮箱"},
		}},
		{"routeName": "Log", "menuTitle": "日志管理", "icon": "el-icon-tickets", "children": []gin.H{
			{"routeName": "PlatformLog", "menuTitle": "平台日志"},
			{"routeName": "SftpLog", "menuTitle": "SFTP日志"},
		}},
		{"routeName": "System", "menuTitle": "系统安全", "icon": "el-icon-lock", "children": []gin.H{
			{"routeName": "SystemUpdate", "menuTitle": "系统更新"},
			{"routeName": "Antivirus", "menuTitle": "病毒管理"},
			{"routeName": "SystemHardening", "menuTitle": "系统加固"},
		}},
		{"routeName": "Settings", "menuTitle": "平台设置", "icon": "el-icon-setting", "children": []gin.H{
			{"routeName": "RoleManagement", "menuTitle": "角色管理"},
			{"routeName": "LocalUserManagement", "menuTitle": "本地账号"},
			{"routeName": "PasswordPolicy", "menuTitle": "密码策略"},
		}},
		{"routeName": "SftpModuleManagement", "menuTitle": "SFTP 管理", "icon": "el-icon-cpu", "children": []gin.H{
			{"routeName": "HotLabelConfig", "menuTitle": "标签上传配置"},
			{"routeName": "ChinaUnicomConfig", "menuTitle": "中国联通配置"},
		}},
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    menus,
	})
}

// GetRoleList 获取角色列表
func GetRoleList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.Query("name")

	var roles []models.Role
	var total int64

	query := dao.DB.Model(&models.Role{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	query.Count(&total)

	if err := query.Offset((page - 1) * limit).Limit(limit).
		Preload("Menus").Preload("LDAPGroups").
		Order("created_at DESC").Find(&roles).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询角色列表失败: " + err.Error()})
		return
	}

	// 补充关联安全组数量
	type RoleResp struct {
		models.Role
		LDAPGroupCount int `json:"ldapGroupCount"`
	}
	var respList []RoleResp
	for _, role := range roles {
		respList = append(respList, RoleResp{
			Role:           role,
			LDAPGroupCount: len(role.LDAPGroups),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":     respList,
			"total":    total,
			"page":     page,
			"limit":    limit,
		},
	})
}

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Menus       []string `json:"menus"`       // 路由名称列表
		LDAPGroups  []gin.H  `json:"ldapGroups"`  // [{group_dn: "...", group_name: "..."}]
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请求数据格式错误"})
		return
	}

	role := models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	// 添加菜单
	for _, menuName := range req.Menus {
		role.Menus = append(role.Menus, models.RoleMenu{
			RouteName: menuName,
			MenuTitle: menuName,
		})
	}

	// 添加LDAP安全组
	for _, group := range req.LDAPGroups {
		dn, _ := group["group_dn"].(string)
		name, _ := group["group_name"].(string)
		if dn != "" {
			role.LDAPGroups = append(role.LDAPGroups, models.RoleLDAPGroup{
				GroupDN:   dn,
				GroupName: name,
			})
		}
	}

	if err := dao.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建角色失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建角色成功",
		"data":    role,
	})
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的角色ID"})
		return
	}

	var req struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Menus       []string `json:"menus"`
		LDAPGroups  []gin.H  `json:"ldapGroups"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请求数据格式错误"})
		return
	}

	role, err := models.GetRoleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "角色不存在"})
		return
	}

	// 超级管理员角色：禁止修改名称、描述、菜单权限，仅允许更新LDAP安全组
	if role.Name == superAdminRoleName {
		// 仅重建LDAP安全组关联
		dao.DB.Where("role_id = ?", role.ID).Delete(&models.RoleLDAPGroup{})
		for _, group := range req.LDAPGroups {
			dn, _ := group["group_dn"].(string)
			name, _ := group["group_name"].(string)
			if dn != "" {
				dao.DB.Create(&models.RoleLDAPGroup{
					RoleID:    role.ID,
					GroupDN:   dn,
					GroupName: name,
				})
			}
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新角色成功"})
		return
	}

	// 更新基本信息
	dao.DB.Model(role).Updates(map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
	})

	// 重建菜单关联
	dao.DB.Where("role_id = ?", role.ID).Delete(&models.RoleMenu{})
	for _, menuName := range req.Menus {
		dao.DB.Create(&models.RoleMenu{
			RoleID:    role.ID,
			RouteName: menuName,
			MenuTitle: menuName,
		})
	}

	// 重建LDAP安全组关联
	dao.DB.Where("role_id = ?", role.ID).Delete(&models.RoleLDAPGroup{})
	for _, group := range req.LDAPGroups {
		dn, _ := group["group_dn"].(string)
		name, _ := group["group_name"].(string)
		if dn != "" {
			dao.DB.Create(&models.RoleLDAPGroup{
				RoleID:    role.ID,
				GroupDN:   dn,
				GroupName: name,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新角色成功"})
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的角色ID"})
		return
	}

	// 禁止删除超级管理员角色
	role, err := models.GetRoleByID(uint(id))
	if err == nil && role.Name == superAdminRoleName {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "不能删除超级管理员角色"})
		return
	}

	// 检查是否有用户绑定该角色
	var userCount int64
	dao.DB.Model(&models.LocalUser{}).Where("role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": fmt.Sprintf("该角色下还有%d个用户，无法删除", userCount)})
		return
	}

	if err := dao.DB.Delete(&models.Role{}, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除角色成功"})
}

// GetRoleDetail 获取角色详情
func GetRoleDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的角色ID"})
		return
	}

	role, err := models.GetRoleByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 400, "message": "角色不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    role,
	})
}

// GetRoleSelect 获取角色下拉列表
func GetRoleSelect(c *gin.Context) {
	roles, err := models.GetRoleSelect()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询角色列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    roles,
	})
}