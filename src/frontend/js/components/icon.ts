import icons from "../icons";
import { defineComponent } from "vue";

export default defineComponent({
  props: { name: { type: String, required: true } },
  template: '<span class="icon" v-html="content"></span>',
  computed: {
    content: function () {
      return (icons as Record<string, string>)[this.name] || "";
    },
  },
});
