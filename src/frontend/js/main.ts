import { createApp, h } from "vue";
import i18n from "./i18n";
import App from "./pages/App.vue";
import Login from "./pages/Login.vue";

const application = createApp({
  render() {
    return h(window.app.authenticated ? App : Login);
  },
});
application.use(i18n);
application.mount("#app");
