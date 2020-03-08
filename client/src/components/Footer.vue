<template>
  <div class="footer">
    <p>&copy; 2020 Powered by pyspa.</p>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Footer extends Vue {}
</script>

<style scoped lang="scss">
.footer {
  display: flex;
  justify-content: center;
  align-items: flex-end;
  margin-top: 16px;
}
</style>
