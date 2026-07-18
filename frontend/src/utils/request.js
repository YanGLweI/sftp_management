import axios from 'axios'
import { MessageBox, Message } from 'element-ui'
import store from '@/store'
import { getToken } from '@/utils/auth'

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 600000 // request timeout 10分钟
})

// !请求拦截器
service.interceptors.request.use(
  config => {
    // do something before request is sent

    if (store.getters.token) {
      // let each request carry token
      // ['X-Token'] is a custom headers key
      // please modify it according to the actual situation
      config.headers['token'] = getToken()
    }
    // 新增SFTP-Token头
    config.headers['X-SFTP-Token'] = sessionStorage.getItem("sftp_token")
    return config
  },
  error => {
    // do something with request error
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// !响应拦截器
service.interceptors.response.use(
  /**
   * If you want to get http information such as headers or status
   * Please return  response => response
  */

  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  response => {
    // 如果是blob类型，直接返回response，不进行code检查
    if (response.config.responseType === 'blob') {
      return response;
    }
    const res = response.data

    //! 响应失败的处理,真实服务器返回的code可能是 20000 或 200
    if (res.code !== 20000 && res.code != 200) {
      Message({
        message: res.message || 'Error',
        type: 'error',
        duration: 5 * 1000
      })

      // 50008: Illegal token; 50012: Other clients logged in; 50014: Token expired;
      if (res.code === 50008 || res.code === 50012 || res.code === 50014) {
        // to re-login
        MessageBox.confirm('你的登录已过期，你可以取消以停留在该页面，或者再次登录。', '登录过期', {
          confirmButtonText: '重新登录',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          store.dispatch('user/resetToken').then(() => {
            location.reload()
          })
        })
      }
      return Promise.reject(new Error(res.message || 'Error'))
    } else {
      // !响应成功
      return res
    }
  },
  // error => {
  //   console.log('err' + error) // for debug
  //   Message({
  //     message: error.message,
  //     type: 'error',
  //     duration: 5 * 1000
  //   })
  //   return Promise.reject(error)
  // }
  async error => {
    console.log('错误详情:', error) // 调试用
    
    // 初始化错误消息
    let errorMessage = '网络错误或服务器异常'
    
    // 处理有响应的错误（如500、404等状态码）
    if (error.response) {
      try {
        // 尝试解析响应数据
        const responseData = error.response.data
        
        // 如果是JSON格式且包含message字段
        if (typeof responseData === 'object' && responseData.message) {
          errorMessage = responseData.message
        } 
        // 如果是文本格式
        else if (typeof responseData === 'string') {
          errorMessage = responseData
        }
        // 如果是blob类型（某些错误可能仍返回blob）
        else if (responseData instanceof Blob) {
          const text = await responseData.text()
          const jsonData = JSON.parse(text)
          errorMessage = jsonData.message || '服务器返回错误'
        }
      } catch (e) {
        // 解析失败时使用HTTP状态信息
        errorMessage = `请求失败 (${error.response.status}): ${error.response.statusText}`
      }
    } 
    // 处理无响应的错误（如网络中断）
    else if (error.message) {
      errorMessage = error.message.includes('timeout') 
        ? '请求超时，请稍后重试' 
        : error.message
    }

    // 显示错误消息
    Message({
      message: errorMessage,
      type: 'error',
      duration: 5 * 1000
    })
    
    return Promise.reject(error)
  }
)

export default service
