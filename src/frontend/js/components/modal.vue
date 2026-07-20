<template>
  <div class="modal custom-modal" tabindex="-1" v-if="$props.open">
    <div class="modal-dialog">
      <div class="modal-content" ref="content">
        <div class="modal-body">
          <slot v-if="$props.open"></slot>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: ["open"],
  data() {
    return { opening: false };
  },
  watch: {
    open(newVal) {
      if (newVal) {
        this.opening = true;
        document.addEventListener("click", this.handleClick);
      } else {
        document.removeEventListener("click", this.handleClick);
      }
    },
  },
  methods: {
    handleClick(e: Event) {
      const target = e.target as HTMLElement;
      if (this.opening) {
        this.opening = false;
        return;
      }
      if (target.closest(".modal-content") == null) this.$emit("hide");
    },
  },
});
</script>
