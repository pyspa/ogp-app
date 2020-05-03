<template>
  <div class="home">
    <form @submit.prevent="onSubmit" class="form">
      <div class="word-input">
        <input
          type="text"
          v-model="words"
          placeholder="Please input your words."
          required
        />
      </div>
      <div>
        <label>space ogp: </label>
        <input type="checkbox" v-model="isSpace" />
      </div>
      <button type="submit">OGP!!!</button>
    </form>
    <loading v-if="loading" />
  </div>
</template>

<script>
import Vue from "vue";
import Component from "vue-class-component";

import Loading from "@/components/Loading";

@Component({
  components: {
    Loading
  }
})
export default class Home extends Vue {
  loading = false;
  words = "";
  isSpace = false;

  get apiHost() {
    return process.env.VUE_APP_API_HOST || "";
  }

  async onSubmit() {
    try {
      this.loading = true;
      const response = await this.axios.post(`${this.apiHost}/api/image`, {
        words: this.words,
        isSpace: this.isSpace
      });
      this.$router.push({ name: "Page", params: { id: response.data.id } });
    } catch (e) {
      alert(e);
    } finally {
      this.loading = false;
    }
  }
}
</script>
