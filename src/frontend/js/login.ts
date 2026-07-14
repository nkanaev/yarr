import template from "./templates/login.html" with { type: "text" };
import icons from "./icons";
import { defineComponent } from "vue";

export default defineComponent({
  template: template,
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
      fetch("./login", { method: "POST", body: data }).then((res) => {
        if (res.ok) {
          // TODO: reload settings instead of refreshing the page
          document.location.assign("./");
        } else {
          this.hasError = true;
        }
      });
    },
  },
});
