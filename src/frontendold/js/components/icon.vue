<template>
  <span class="c-icon" :class="{ 'is-small': $props.small }" v-html="content"></span>
</template>

<script lang="ts">
import icons from "../icons";
import { defineComponent, PropType } from "vue";

export type IconName = keyof typeof icons;

export default defineComponent({
  props: {
    name: { type: String as PropType<IconName>, required: true },
    small: { type: Boolean, required: false },
  },
  computed: {
    content() {
      return (icons as Record<string, string>)[this.name] || "";
    },
  },
});
</script>
