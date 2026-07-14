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
      var data = new FormData(event.target);
      fetch("./login", { method: "POST", body: data }).then(
        function (res) {
          if (res.ok) {
            // TODO:
            document.location.assign("./");
          } else {
            this.hasError = true;
          }
        }.bind(this),
      );
    },
  },
});
