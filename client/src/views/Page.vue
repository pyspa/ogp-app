<template>
  <div>
    <div class="form">
      <div class="word-input noshadow">
        <input
          type="text"
          id="url"
          v-model="shareUrl"
          @click.prevent="onClickUrl"
        /><button @click.prevent="copyShareUrl">
          Copy
        </button>
      </div>
    </div>
    <div class="image">
      <img :src="imagePath" />
    </div>
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
  get shareUrl() {
    const url = new URL(document.URL);
    return `${url.origin}/p/${this.$route.params.id}`;
  }
  onClickUrl() {
    document.getElementById("url").select();
  }
  copyShareUrl() {
    document.getElementById("url").select();
    document.execCommand("copy");
    alert("success!");
  }
}
</script>

<style lang="scss" scoped>
.image {
  margin-top: 20px;
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
