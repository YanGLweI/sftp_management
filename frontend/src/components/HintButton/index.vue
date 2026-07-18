<template>
  <!-- 可以巧妙的利用a标签实现按钮提示的功能 -->
  <a :title="title" style="margin: 5px;">
    <el-button  v-bind="$attrs" v-on="$listeners"></el-button>
  </a>
</template>

<script>
export default {
  name:'HintButton',
  props:['title'],
  mounted(){
    // $attrs属于组件的一个属性,可以获取到父组件传递过来的props
    // console.log(this.$attrs)
    // $listeners属于组件的一个属性,可以获取到父组件传递过来的自定义事件
    // console.log(this.$listeners)
  }
}
</script>

<style>

</style>