<template>
  <div class="mx-auto my-2 p-3" style="max-width: 20rem">
    <form @submit.prevent="login" class="d-flex flex-column">
      <div class="login-logo my-5 d-flex justify-content-center" v-html="logo"></div>
      <label for="username" class="mb-2">{{ $t("username") }}</label>
      <input name="username" class="c-input" id="username" autocomplete="off" required autofocus />
      <label for="password" class="mb-2 mt-3">{{ $t("password") }}</label>
      <input name="password" class="c-input" id="password" type="password" required />
      <button class="c-button mt-3" type="submit">{{ $t("login") }}</button>
      <div class="fixed-top p-2 text-center bg-danger text-white" v-if="hasError">
        {{ $t("login_error") }}
      </div>
    </form>
  </div>
</template>

<script lang="ts">
import { Lang } from "../i18n";
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
    this.$t.set(document.documentElement.lang as Lang);
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
