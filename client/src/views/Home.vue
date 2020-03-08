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
        <button type="submit">OGP!!!</button>
      </div>
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

  get apiHost() {
    return process.env.VUE_APP_API_HOST || "";
  }

  async onSubmit() {
    try {
      this.loading = true;
      const response = await this.axios.post(`${this.apiHost}/api/image`, {
        words: this.words
      });
      this.$router.push({ name: "Page", params: { id: response.data.id } });
    } catch (e) {
      alert(e);
      console.error(e);
    }
    this.loading = false;
  }
}
</script>
