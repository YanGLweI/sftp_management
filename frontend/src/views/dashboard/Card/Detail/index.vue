<template>
  <div>
    <div class="card-header">
      <span>{{ title }}</span><i class="el-icon-info"></i>
    </div>
    <div class="card-content">{{ count }}</div>
    <div class="card-charts">
      <slot name="charts"></slot>
    </div>
    <div class="card-footer">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script>

export default {
  name:'Detail',
  props:['title','count']
}
</script>

<style>
.card-header{
  /* 设置为弹性容器 */
  display: flex;
  /* 在弹性容器中，将子元素在主轴方向上两端对齐，子元素之间的间距相等。这样可以使子元素均匀分布在容器中，两端不留空白。 */
  justify-content: space-between;
  color: #d9d9d9;
}
.card-content{
  font-size: 30px;
  padding: 10px 10px;
}
.card-charts{
  height: 50px;
  display: flex;
  align-items: flex-end;
  padding-bottom: 5px;
}
.card-footer{
  border-top: 1px solid #eee;
  padding-top: 10px;
}
</style>