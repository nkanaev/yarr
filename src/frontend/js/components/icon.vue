<template>
  <span class="c-icon" :class="{ 'is-small': $props.small }" v-html="content"></span>
</template>

<script lang="ts">
import icons from "../icons";
import { defineComponent } from "vue";

export type IconName = keyof typeof icons;

export default defineComponent({
  props: {
    name: { type: String, required: true } as unknown as { type: () => IconName },
    small: { type: Boolean, required: false },
  },
  computed: {
    content() {
      return (icons as Record<string, string>)[this.name] || "";
    },
  },
});
</script>
