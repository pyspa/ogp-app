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
        <button type="submit">OPG!!!</button>
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
  get form() {
    return this.$form.createForm(this, { name: "horizontal_login" });
  }

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

<style lang="scss" scoped>
.form {
  margin: 0 auto;
  max-width: 320px;
  button {
    background-image: linear-gradient(to left, #fff, #f2f2f2);
    border: 1px solid #e3e3e3;
    border-radius: 3px;
    width: 63px;
    height: 34px;
    &:hover {
      border-color: #bdbdbd;
    }
  }
}

.word-input {
  background: #fff;
  border-radius: 3px;
  padding: 12px 8px 12px 18px;
  margin-bottom: 15px;
  position: relative;
  box-shadow: 0 0 5px rgba(16, 14, 23, 0.25);
}

.form input {
  border: 0;
  font-family: 24px;
  font-size: 16px;
  width: 210px;
  margin-right: 5px;
  &:focus {
    outline: none;
  }
}
</style>
