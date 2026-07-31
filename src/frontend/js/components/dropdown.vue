<template>
  <div class="position-relative d-inline-flex flex-column">
    <button class="c-button-pill" :class="toggleClass" @click="toggle">
      <slot name="button"></slot>
    </button>
    <div
      ref="dropdown"
      class="c-dropdown position-absolute top-100 z-1"
      :class="{
        'start-50 translate-middle-x': $props.drop === 'center',
        'end-0': $props.drop === 'right',
      }"
      v-if="open">
      <slot></slot>
    </div>
  </div>
</template>

<script lang="ts">
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
      validator: (value: string) => ["right", "center"].includes(value),
    },
    title: {
      type: String,
      required: true,
    },
  },
  data() {
    return { open: false };
  },
  methods: {
    toggle(event: Event) {
      event.stopPropagation();
      this.open ? this.hide() : this.show();
    },
    show() {
      this.open = true;
      document.addEventListener("click", this.clickHandler);
    },
    hide() {
      this.open = false;
      document.removeEventListener("click", this.clickHandler);
    },
    clickHandler(e: MouseEvent) {
      const dropdown = this.$refs.dropdown as HTMLElement;
      const target = e.target as HTMLElement;
      // Did click happen outside this component's DOM element?
      if (dropdown === null || !dropdown.contains(target)) this.hide();
      // Is the target (or its parent) a clickable option (e.g., button or link)?
      if (target.closest(".c-dropdown-item") !== null) this.hide();
    },
  },
});
</script>
