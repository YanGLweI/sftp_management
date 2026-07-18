//?平台属性管理模块API
import request from '@/utils/request.js'

//! 获取一级分类的接口 /admin/product/getCategory1    get 
export const reqCategory1List = () => request({url:'/admin/product/getCategory1',method:'get'})

//! 获取二级分类的接口 /admin/product/getCategory2/{category1Id}    get 
// * 通过一级分类id 获取二级分类列表
export const reqCategory2List = (category1Id) => request({url:`/admin/product/getCategory2/${category1Id}`,method:'get'})

//! 获取三级分类的接口  /admin/product/getCategory3/{category2Id}    get
//* 通过二级分类id 获取三级分类列表
export const reqCategory3List = (category2Id) => request({url:`/admin/product/getCategory3/${category2Id}`,method:'get'})

//!  通过分类id 获取平台属性数据 /admin/product/attrInfoList/{category1Id}/{category2Id}/{category3Id}  get
export const reqAttrList = (category1Id,category2Id,category3Id) => request({url:`/admin/product/attrInfoList/${category1Id}/${category2Id}/${category3Id}`,method:'get'})

//! 添加或者修改属性与属性值的接口 /admin/product/saveAttrInfo   post  
/* 
  {
    "attrName": "string",  属性名
    "attrValueList": [     属性值数组,属性值可以是多个
      {
        "attrId": 0,            属性的id
        "valueName": "string"   属性值
      }
    ],
    "categoryId": 0,          category3Id  
    "categoryLevel": 3,       
  }
*/
export const reqAddOrUpdateAttr = (data) => request({url:'/admin/product/saveAttrInfo',method:'post',data})

//! 删除属性的接口  /admin/product/deleteAttr/{attrId}  delete
export const  reqDeleteAttr = (attrId) => request({url:`/admin/product/deleteAttr/${attrId}`,method:'delete'})
