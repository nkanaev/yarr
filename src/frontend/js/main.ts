import { createApp, ref, defineComponent, h, nextTick, onMounted } from "vue";
import i18n from "./i18n";
import App from "./pages/App.vue";
import Login from "./pages/Login.vue";
import api from "./api";
import { setupKeybindings } from "./key";

const ready = ref(window.app.authenticated);

const Root = defineComponent({
  setup() {
    const appRef = ref<InstanceType<typeof App>>();

    const onLogin = () => {
      api.settings.get().then(settings => {
        window.app.settings = settings;
        window.app.authenticated = true;
        ready.value = true;
        nextTick(() => appRef.value && setupKeybindings(appRef.value));
      });
    };

    onMounted(() => {
      if (appRef.value) setupKeybindings(appRef.value);
    });

    return () =>
      ready.value
        ? h(App, { ref: appRef })
        : h(Login, { onLogin });
  },
});

const application = createApp(Root);
application.use(i18n);
application.mount("#app");
