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
  </div>
</template>

<script>
import Vue from "vue";
import Component from "vue-class-component";

@Component
export default class Home extends Vue {
  words = "";

  get apiHost() {
    return process.env.VUE_APP_API_HOST || "";
  }

  onSubmit() {
    this.axios
      .post(`${this.apiHost}/api/image`, {
        words: this.words
      })
      .then(response => {
        this.$router.push({ name: "Page", params: { id: response.data.id } });
      })
      .catch(e => {
        alert(e);
      });
  }
}
</script>
