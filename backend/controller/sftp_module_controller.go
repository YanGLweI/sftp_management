package controller

import (
	"encoding/json"
	"net/http"
	"sftpbackend/models"

	"github.com/gin-gonic/gin"
)

// GetModuleConfig 获取指定模块的配置
func GetModuleConfig(c *gin.Context) {
	moduleName := c.Param("name")
	if moduleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "模块名称不能为空",
		})
		return
	}

	config, err := models.GetSFTPModuleConfig(moduleName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "模块配置不存在：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    config,
	})
}

// UpdateModuleConfig 更新模块配置
func UpdateModuleConfig(c *gin.Context) {
	moduleName := c.Param("name")
	if moduleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "模块名称不能为空",
		})
		return
	}

	var req struct {
		LoginType      string `json:"loginType" binding:"required"` // local or ldap
		EnabledRoles   []uint `json:"enabledRoles"`                 // 允许登录的角色 ID 列表
		DualAuthEnabled bool  `json:"dualAuthEnabled"`              // 是否启用双控（仅中国联通）
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求数据格式错误",
		})
		return
	}

	// 验证登录类型
	if req.LoginType != models.LoginTypeLocal && req.LoginType != models.LoginTypeLDAP {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "登录类型必须为 local 或 ldap",
		})
		return
	}

	// 获取现有配置
	config, err := models.GetSFTPModuleConfig(moduleName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "模块配置不存在",
		})
		return
	}

	// 更新配置
	config.LoginType = req.LoginType
	if roleBytes, marshalErr := json.Marshal(req.EnabledRoles); marshalErr == nil {
		config.EnabledRoles = string(roleBytes)
	}

	// 仅中国联通模块可以配置双控开关
	if moduleName == models.ModuleNameChinaUnicom {
		config.DualAuthEnabled = req.DualAuthEnabled
	} else {
		// 标签上传模块不支持双控
		config.DualAuthEnabled = false
	}

	if err := models.UpdateSFTPModuleConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新配置失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    config,
	})
}

// GetAllModuleConfigs 获取所有模块配置
func GetAllModuleConfigs(c *gin.Context) {
	configs, err := models.GetAllSFTPModuleConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取配置失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    configs,
	})
}

// GetPublicModuleConfigs 获取所有模块的公共配置（无需登录，供 /file SFTP 登录页渲染表单）
// 仅返回模块名、登录方式、双控开关，不返回角色白名单等敏感信息
func GetPublicModuleConfigs(c *gin.Context) {
	configs, err := models.GetAllSFTPModuleConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取配置失败：" + err.Error(),
		})
		return
	}

	publicConfigs := make([]gin.H, 0, len(configs))
	for _, cfg := range configs {
		publicConfigs = append(publicConfigs, gin.H{
			"moduleName":      cfg.ModuleName,
			"loginType":       cfg.LoginType,
			"dualAuthEnabled": cfg.DualAuthEnabled,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    publicConfigs,
	})
}

// CheckRolePermission 检查角色 ID 是否在模块 enabled_roles 白名单中（本地账号登录用）
func CheckRolePermission(roleID uint, moduleName string) bool {
	config, err := models.GetSFTPModuleConfig(moduleName)
	if err != nil {
		// 配置不存在时，默认拒绝访问
		return false
	}

	// 解析 enabled_roles JSON 数组
	var enabledRoles []uint
	if err := json.Unmarshal([]byte(config.EnabledRoles), &enabledRoles); err != nil {
		return false
	}

	// 如果没有配置任何角色，默认拒绝访问
	if len(enabledRoles) == 0 {
		return false
	}

	// 检查当前角色是否在允许列表中
	for _, role := range enabledRoles {
		if role == roleID {
			return true
		}
	}

	return false
}

// CheckLDAPRolePermission 检查 LDAP 用户的安全组是否能匹配到模块允许登录的角色
// 逻辑：LDAP 用户的安全组（memberOf）→ RoleLDAPGroup 表匹配角色 → 角色是否在模块 enabled_roles 白名单中
func CheckLDAPRolePermission(userGroups []string, moduleName string) bool {
	if len(userGroups) == 0 {
		return false
	}

	config, err := models.GetSFTPModuleConfig(moduleName)
	if err != nil {
		// 配置不存在时，默认拒绝访问
		return false
	}

	// 解析 enabled_roles JSON 数组
	var enabledRoles []uint
	if err := json.Unmarshal([]byte(config.EnabledRoles), &enabledRoles); err != nil {
		return false
	}

	// 如果没有配置任何角色，默认拒绝访问
	if len(enabledRoles) == 0 {
		return false
	}

	// 查询所有角色的 LDAP 安全组关联
	var roleLinks []models.RoleLDAPGroup
	if err := models.GetAllRoleLDAPGroups(&roleLinks); err != nil {
		return false
	}

	// 用户安全组 → 匹配角色 → 检查角色是否在白名单中
	enabledSet := make(map[uint]bool, len(enabledRoles))
	for _, id := range enabledRoles {
		enabledSet[id] = true
	}

	for _, link := range roleLinks {
		for _, groupDN := range userGroups {
			if link.GroupDN == groupDN && enabledSet[link.RoleID] {
				return true
			}
		}
	}

	return false
}
