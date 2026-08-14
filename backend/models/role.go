package models

import (
	"sftpbackend/dao"

	"gorm.io/gorm"
)

// Role 角色
type Role struct {
	gorm.Model
	Name        string           `json:"name" gorm:"column:name;type:varchar(64);uniqueIndex;not null"`
	Description string           `json:"description" gorm:"column:description;type:varchar(255)"`
	Menus       []RoleMenu       `json:"menus" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	LDAPGroups  []RoleLDAPGroup  `json:"ldapGroups" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}

// RoleMenu 角色可访问的菜单路由
type RoleMenu struct {
	gorm.Model
	RoleID    uint   `json:"roleId" gorm:"column:role_id;index;not null"`
	RouteName string `json:"routeName" gorm:"column:route_name;type:varchar(128);not null"`
	MenuTitle string `json:"menuTitle" gorm:"column:menu_title;type:varchar(64)"`
	ParentID  *uint  `json:"parentId" gorm:"column:parent_id"`
}

// RoleLDAPGroup 角色关联的LDAP安全组
type RoleLDAPGroup struct {
	gorm.Model
	RoleID    uint   `json:"roleId" gorm:"column:role_id;index;not null"`
	GroupDN   string `json:"groupDN" gorm:"column:group_dn;type:varchar(512);not null"`
	GroupName string `json:"groupName" gorm:"column:group_name;type:varchar(128)"`
}

// GetRoleByID 通过ID获取角色（含关联的菜单和安全组）
func GetRoleByID(id uint) (*Role, error) {
	var role Role
	err := dao.DB.Preload("Menus").Preload("LDAPGroups").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByName 通过名称获取角色
func GetRoleByName(name string) (*Role, error) {
	var role Role
	err := dao.DB.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByLDAPGroupDN 通过LDAP安全组DN查找关联的角色
func GetRoleByLDAPGroupDN(groupDN string) (*Role, error) {
	var link RoleLDAPGroup
	err := dao.DB.Where("group_dn = ?", groupDN).First(&link).Error
	if err != nil {
		return nil, err
	}
	return GetRoleByID(link.RoleID)
}

// GetRoleSelect 获取角色下拉列表
func GetRoleSelect() ([]Role, error) {
	var roles []Role
	err := dao.DB.Select("id, name, description").Order("name ASC").Find(&roles).Error
	return roles, err
}