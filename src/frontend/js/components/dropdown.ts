import { defineComponent } from "vue";

export default defineComponent({
  inheritAttrs: false,
  props: {
    toggleClass: {
      type: String,
      required: true,
    },
    drop: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
  },
  data() {
    return { open: false };
  },
  template: `
    <div class="dropdown" :class="$attrs.class">
      <button ref="btn" @click="toggle" :class="btnToggleClass" :title="$props.title"><slot name="button"></slot></button>
      <div ref="menu" class="dropdown-menu" :class="{show: open}"><slot v-if="open"></slot></div>
    </div>
  `,
  computed: {
    btnToggleClass() {
      var c = this.$props.toggleClass || "";
      c += " dropdown-toggle dropdown-toggle-no-caret";
      c += this.open ? " show" : "";
      return c.trim();
    },
  },
  methods: {
    toggle() {
      this.open ? this.hide() : this.show();
    },
    show() {
      this.open = true;
      const menu = this.$refs.menu as HTMLElement;
      const btn = this.$refs.btn as HTMLElement;
      menu.style.top = btn.offsetHeight + "px";
      var drop = this.$props.drop;

      if (drop === "right") {
        menu.style.left = "auto";
        menu.style.right = "0";
      } else if (drop === "center") {
        this.$nextTick(() => {
          const b = this.$refs.btn as HTMLElement;
          const m = this.$refs.menu as HTMLElement;
          m.style.left =
            "-" +
            (m.getBoundingClientRect().width -
              b.getBoundingClientRect().width) /
              2 +
            "px";
        });
      }

      document.addEventListener("click", this.clickHandler);
    },
    hide() {
      this.open = false;
      document.removeEventListener("click", this.clickHandler);
    },
    clickHandler(e: MouseEvent) {
      const target = e.target as HTMLElement;
      var dropdown = target.closest(".dropdown");
      if (dropdown == null || dropdown != this.$el) return this.hide();
      if (target.closest(".dropdown-item") != null) return this.hide();
    },
  },
});
