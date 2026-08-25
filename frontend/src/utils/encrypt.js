import JSEncrypt from 'jsencrypt'
import axios from 'axios'

let publicKeyCache = null
let loadingPromise = null

export async function loadPublicKey() {
  // 检查缓存是否有效且非空
  if (publicKeyCache && typeof publicKeyCache === 'string' && publicKeyCache.length > 0) {
    return publicKeyCache
  }
  
  // 返回现有加载中的 Promise
  if (loadingPromise) return loadingPromise
  
  // 创建新的 loading Promise
  loadingPromise = (async () => {
    try {
      const res = await axios.get('/dev-api/rsa/public-key', { timeout: 5000 })
      if (res.data.code === 200) {
        publicKeyCache = res.data.data
        loadingPromise = null
        return publicKeyCache
      }
      throw new Error(res.data.message || '获取公钥失败')
    } catch (error) {
      loadingPromise = null
      console.error('RSA 公钥加载失败:', error)
      throw error
    }
  })()
  
  return loadingPromise
}

// rsaEncrypt 改为异步函数以支持动态获取公钥
export async function rsaEncrypt(txt) {
  const pubKey = await loadPublicKey()
  const encryptor = new JSEncrypt()
  encryptor.setPublicKey(pubKey) // 设置公钥
  return encryptor.encrypt(txt)  // 执行加密
}