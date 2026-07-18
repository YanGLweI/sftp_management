package models

import (
	"sftpbackend/dao"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Options struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// ! 分页获取所有通讯录列表
func (c *Contact) GetContactList(page int, limit int, name string) (contacts []Contact, total int64, err error) {
	// 根据页码和每页数量计算偏移量
	offset := (page - 1) * limit
	// 构建查询条件
	query := dao.DB.Offset(offset).Limit(limit)
	// 先构建一个用于统计总数的查询条件副本，不添加分页相关设置
	countQuery := dao.DB.Model(&Contact{})

	// 若果有用户名参数，则添加模糊查询条件
	if name != "" {
		likePattern := "%" + name + "%"
		query = query.Where("name LIKE ?", likePattern)
		countQuery = countQuery.Where("name LIKE ?", likePattern)
	}
	// 查询满足条件的所有用户
	if err = query.Find(&contacts).Error; err != nil {
		return nil, 0, err
	}
	// 统计满足条件的用户总数（不分页）
	if err = countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

// ! 查询通讯录是否存在，通过姓名
func (c *Contact) ExistOrNot(name string) (err error) {
	var contact Contact
	err = dao.DB.Where("name = ?", name).First(&contact).Error
	return
}

// ! 查询通讯录是否存在，通过ID
func (c *Contact) ExistOrNotByID(id uint) (err error) {
	var contact Contact
	err = dao.DB.Where("id = ?", id).First(&contact).Error
	return
}

// ! 添加一个联系人
func (c *Contact) AddContact() (err error) {
	err = dao.DB.Create(&c).Error
	return
}

// ! 更新一个联系人
func (c *Contact) UpdateContact() (err error) {
	// 根据ID更新联系人信息
	err = dao.DB.Model(&Contact{}).Where("id = ?", c.ID).Updates(Contact{Name: c.Name, Email: c.Email}).Error
	return
}

// ! 删除一个联系人
func (c *Contact) DeleteContact(id uint) (err error) {
	err = dao.DB.Delete(&Contact{}, id).Error
	return
}

// ! 批量删除联系人
func (c *Contact) DeleteContacts(ids []uint) (err error) {
	err = dao.DB.Delete(&Contact{}, ids).Error
	return
}

func (c *Contact) GetContactoptions() (options []Options, err error) {
	var contacts []Contact
	err = dao.DB.Select("name", "email").Find(&contacts).Error
	// 通过新的结构体整理返回的数据
	for _, contact := range contacts {
		options = append(options, Options{
			Value: contact.Email,
			Label: contact.Name,
		})
	}
	return options, err
}
