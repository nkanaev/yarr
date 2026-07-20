<template>
  <div class="login-page">
    <form @submit.prevent="login">
      <div class="logo" v-html="logo"></div>
      <div class="text-danger text-center my-3" v-if="hasError">{{ $t("login_error") }}</div>
      <div class="form-group">
        <label for="username">{{ $t("username") }}</label>
        <input
          name="username"
          class="form-control"
          id="username"
          autocomplete="off"
          required
          autofocus />
      </div>
      <div class="form-group">
        <label for="password">{{ $t("password") }}</label>
        <input name="password" class="form-control" id="password" type="password" required />
      </div>
      <button class="btn btn-block btn-default" type="submit">{{ $t("login") }}</button>
    </form>
  </div>
</template>

<script lang="ts">
import icons from "../icons";
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    onLogin: { type: Function, required: true },
  },
  data() {
    return {
      logo: icons.anchor,
      hasError: false,
    };
  },
  created() {
    this.$setLang(window.app.settings.language);
  },
  methods: {
    login(event: Event) {
      event.preventDefault();
      var data = new FormData(event.target as HTMLFormElement);
      fetch("./login", { method: "POST", body: data }).then(res => {
        if (res.ok) {
          this.onLogin();
        } else {
          this.hasError = true;
        }
      });
    },
  },
});
</script>
