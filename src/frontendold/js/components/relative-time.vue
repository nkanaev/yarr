<template>
  <time :datetime="val" :title="title">{{ formatted }}</time>
</template>

<script lang="ts">
import { dateRepr, dateTimeString, relRepaintDelay } from "../utils";
import { defineComponent } from "vue";

export default defineComponent({
  props: ["val", "locale"],
  data() {
    return {
      date: new Date(this.val),
      formatted: "" as string,
      timer: undefined as number | undefined,
    };
  },
  computed: {
    title(): string {
      return dateTimeString(this.date, this.locale);
    },
  },
  created() {
    this.repaint();
  },
  unmounted() {
    window.clearTimeout(this.timer);
  },
  methods: {
    repaint() {
      this.formatted = dateRepr(this.date, this.locale);
      window.clearTimeout(this.timer);
      const delay = relRepaintDelay(this.date);
      if (delay !== null) {
        this.timer = window.setTimeout(() => this.repaint(), delay);
      }
    },
  },
  watch: {
    locale() {
      this.repaint();
    },
  },
});
</script>
