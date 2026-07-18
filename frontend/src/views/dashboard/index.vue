<template>
  <div>
    <Card/>
    <Sale/>
    <Observe/>
  </div>
</template>

<script>
import Card from './Card'
import Sale from './Sale'
import Observe from './Observe'
export default {
  name: 'Dashboard',
  components:{Card,Sale,Observe},
  mounted(){
    this.$store.dispatch('getData')
    const h = this.$createElement;
    this.$notify.info({
      title: 'Sorry~',
      message: h('i', { style: 'color: teal'}, '首页开发中,仅供参考...'),
      // position: 'bottom-right'
    });
  }
}
</script>

<style scoped>

</style>
