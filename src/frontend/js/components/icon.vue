<template>
  <span class="icon" v-html="content"></span>
</template>

<script lang="ts">
import icons from "../icons";
import { defineComponent } from "vue";

export default defineComponent({
  props: { name: { type: String, required: true } },
  computed: {
    content() {
      return (icons as Record<string, string>)[this.name] || "";
    },
  },
});
</script>
