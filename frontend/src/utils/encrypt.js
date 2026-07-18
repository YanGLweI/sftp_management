import JSEncrypt from 'jsencrypt'

// 从后端获取的公钥，或是写死在配置文件中
const PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA4BvuaUCd62R3mH4X/AT7
7FWZQI1R6IPqLpn3fGQVIqUDKgqW3KiPlY1LYI5JhYz8O5krzA9kuEkGA6E0VqYs
qoYRU6nI0IG/QbkmQihxZGk0pHoQ5yD2L3Ib3AfBJLgq/9Tlrhyr3+X6T480+NXH
R7HL4V26sdv5cid+0Hdnd2sZqQwrDzAdYumnxW3haWQYbu6rty/sYNMfAU05C/Rh
PRICjBuluYryQn77RSscaxLKWuvq4nAFpMM3DkT8gXtHNZGqzp+iBJL5A9f/u5jc
tirRFxiJvHATJtdEBB8HEfkUj3DX1WqiCJoLFRGyKjRs5WlOLMdE12PuBq5xzc1i
dxpAeVN8Gx6HTT+c3Ok17zIuKxDzO7APjILS7O2x+p9d1VmU/gVuApZIPXhst1sb
txG1f4LtS677OIsymt2YIlIWu9SE55R/+zXEqWiMcPOWf6O0S6Kuu/C6pXsRbZlM
5qVcEdJVFPR6yTMVBvYiGXnm7eJq4locRxzzl1h5XYqzryhpZTky8h0XSLeQ2K3T
K8M4uDOObwozXDJiS+dSxfRgXE7RZvVZgErqK0XN04dOtmOEnfNH4kxBU5REk2FA
tRPzz5KQQzk4HQI/7/GIpGBPgU/HaUaMVGcKpRUdSoI5ivOjEgPvD2eFYYQ+glwX
nBcHL1zRKmqaNwuVBj+Y7IMCAwEAAQ==
-----END PUBLIC KEY-----
`

export function rsaEncrypt(txt) {
  const encryptor = new JSEncrypt()
  encryptor.setPublicKey(PUBLIC_KEY) // 设置公钥
  return encryptor.encrypt(txt)      // 执行加密
}