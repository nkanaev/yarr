<template>
  <dialog ref="dialog" @close="$emit('hide')">
    <slot></slot>
  </dialog>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: ["open"],
  emits: ["hide"],
  watch: {
    open(open) {
      const el = this.$refs.dialog as HTMLDialogElement;
      if (open) el.showModal();
      else el.close();
    },
  },
});
</script>
