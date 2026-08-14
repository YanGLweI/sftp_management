<template>
  <div>
    <SftpModuleConfigForm module-name="hotlabel" :show-dual-auth="false" />
  </div>
</template>

<script>
import SftpModuleConfigForm from '@/views/admin/sftp-module-config/SftpModuleConfigForm'

/**
 * HotLabelConfig 标签上传模块配置页面
 * 配置 SFTP 登录页"标签上传" tab 的登录方式（本地/LDAP）与可登录角色
 */
export default {
  name: 'HotLabelConfig',
  components: { SftpModuleConfigForm }
}
</script>

<style scoped>
  
</style>
