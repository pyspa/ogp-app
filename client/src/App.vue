<template>
  <div id="app">
    <header>
      <h1 class="title"><router-link to="/">OGP.app</router-link></h1>
    </header>
    <router-view />
    <site-footer />
  </div>
</template>

<script>
import Vue from "vue";
import Component from "vue-class-component";

import SiteFooter from "@/components/Footer.vue";

@Component({
  components: {
    SiteFooter
  }
})
export default class Page extends Vue {
  get imagePath() {
    return `/image/${this.$route.params.id}.png`;
  }
}
</script>

<style lang="scss">
body {
  background-color: #fbfbfb;
}
#app {
  font-family: -apple-system, BlinkMacSystemFont, "Helvetica Neue", Arial,
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin: 16px auto 0;
  display: flex;
  justify-content: center;
  align-items: top;
  flex-direction: column;
  text-align: center;
}

.title {
  font-family: "Lato", sans-serif;
  font-weight: 900;
  font-size: 48px;
  color: #35495e;
  letter-spacing: 0.5px;
  a {
    color: #2c3e50;
    text-decoration: none;
    &:hover {
      color: saturate(#2c3e50, 90%);
    }
  }
}
</style>
