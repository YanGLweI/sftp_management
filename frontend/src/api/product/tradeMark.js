//? 这个模块主要获取的是品牌数据的模块
import request from '@/utils/request'

//! 获取品牌列表接口 /admin/product/baseTrademark/{page}/{limit}   get
//? page:当前页码，limit:每页记录数
export const reqTradeMarkList = (page,limit)=>request({url:`/admin/product/baseTrademark/${page}/${limit}`,method:'get'})

//! 处理添加品牌的接口  /admin/product/baseTrademark/save     post  需要2个参数 品牌LOGO:logoUrl 品牌名称:tmName
//! 修改品牌的接口     /admin/product/baseTrademark/update    put  需要3个参数 品牌LOGO:logoUrl 品牌名称:tmName 品牌id:id
//? 区别在于是否带id
//* 可以封装成一个函数
export const reqAddOrUpdateTradeMark = (tradeMark)=>{
  // 如果携带id,代表是修改
  if(tradeMark.id){
    return request({url:'/admin/product/baseTrademark/update ',data:tradeMark,method:'put'})
  }else{
    // 新增品牌
    return request({url:'/admin/product/baseTrademark/save ',data:tradeMark,method:'post'})
  }
}

//! 删除品牌的接口 /admin/product/baseTrademark/remove/{id}  delete
export const reqDeleteTradeMark = (id)=>request({url:`/admin/product/baseTrademark/remove/${id}`,method:'delete'})
