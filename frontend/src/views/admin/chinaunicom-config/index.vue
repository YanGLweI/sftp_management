<template>
  <div>
    <SftpModuleConfigForm module-name="chinaunicom" :show-dual-auth="true" />
  </div>
</template>

<script>
import SftpModuleConfigForm from '@/views/admin/sftp-module-config/SftpModuleConfigForm'

/**
 * ChinaUnicomConfig 中国联通模块配置页面
 * 配置 SFTP 登录页"中国联通" tab 的登录方式（本地/LDAP）、可登录角色与双控开关
 */
export default {
  name: 'ChinaUnicomConfig',
  components: { SftpModuleConfigForm }
}
</script>

<style scoped>
</style>
