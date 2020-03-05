<template>
  <div class="image">
    <img :src="imagePath" />
  </div>
</template>

<script>
import Vue from "vue";
import Component from "vue-class-component";

@Component
export default class Page extends Vue {
  get imagePath() {
    return `/image/${this.$route.params.id}.png`;
  }
}
</script>

<style lang="scss" scoped>
.image {
  img {
    box-shadow: 0 0 5px rgba(16, 14, 23, 0.25);
  }
}
@media screen and (max-width: 480px) {
  .image {
    img {
      max-width: 320px;
    }
  }
}
</style>
