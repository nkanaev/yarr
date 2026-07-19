<template>
  <time :datetime="val">{{ formatted }}</time>
</template>

<script lang="ts">
import { dateRepr } from "../utils";
import { defineComponent } from "vue";

export default defineComponent({
  props: ["val"],
  data() {
    var d = new Date(this.val);
    return {
      date: d,
      formatted: dateRepr(d),
      interval: undefined as number | undefined,
    };
  },
  mounted() {
    this.interval = setInterval(() => {
      this.formatted = dateRepr(this.date);
    }, 600000); // every 10 minutes
  },
  unmounted() {
    clearInterval(this.interval);
  },
});
</script>
