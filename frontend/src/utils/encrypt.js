import JSEncrypt from 'jsencrypt'
import axios from 'axios'

let publicKeyCache = null
let loading = false

export async function loadPublicKey() {
  if (publicKeyCache) return publicKeyCache
  if (loading) {
    // 等待公钥加载完成
    return new Promise((resolve, reject) => {
      const timer = setInterval(() => {
        if (publicKeyCache) {
          clearInterval(timer)
          resolve(publicKeyCache)
        } else if (!loading) {
          clearInterval(timer)
          reject(new Error('公钥加载超时'))
        }
      }, 100)
    })
  }
  
  loading = true
  try {
    const res = await axios.get('/dev-api/rsa/public-key')
    if (res.data.code === 200) {
      publicKeyCache = res.data.data
      return publicKeyCache
    }
    throw new Error(res.data.message || '获取公钥失败')
  } finally {
    loading = false
  }
}

// rsaEncrypt 改为异步函数以支持动态获取公钥
export async function rsaEncrypt(txt) {
  const pubKey = await loadPublicKey()
  const encryptor = new JSEncrypt()
  encryptor.setPublicKey(pubKey) // 设置公钥
  return encryptor.encrypt(txt)  // 执行加密
}