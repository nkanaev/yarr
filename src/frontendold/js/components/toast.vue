<template>
  <TransitionGroup
    tag="div"
    class="position-fixed bottom-0 end-0 pb-3 pe-3 d-flex flex-column align-items-end gap-2 z-4 pe-none"
    name="toast">
    <div
      class="c-toast d-flex align-items-start gap-2"
      v-for="t in list"
      :key="t.id"
      :class="{ 'text-bg-danger text-light': t.opts?.level === 'fail' }">
      <v-icon v-if="t.opts?.level === 'fail'" name="alert-circle" />
      <div class="flex-grow-1">
        <div>{{ t.message.title }}</div>
        <div class="opacity-50" v-if="t.message.description">
          {{ t.message.description }}
        </div>
      </div>
      <button class="c-button-link" @click="remove(t.id)" v-if="t.opts?.closeable === true">
        <v-icon name="x" />
      </button>
    </div>
  </TransitionGroup>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import icon from "./icon.vue";

type Level = "info" | "fail";
type Options = { level?: Level; closeable?: boolean; err?: unknown };
type Toast = { id: string; message: Message; opts?: Options };
type Message = { title: string; description?: string };

const getUUID = () =>
  crypto?.randomUUID?.() ?? Date.now().toString(36) + Math.random().toString(36).slice(2);

export default defineComponent({
  components: { "v-icon": icon },
  data() {
    return { list: [] as Toast[] };
  },
  methods: {
    addToast(message: Message, opts?: Options) {
      var t: Toast = { id: getUUID(), message, opts: opts };
      this.list.push(t);
      setTimeout(() => this.remove(t.id), 3000);
    },
    remove(id: string) {
      this.list = this.list.filter(x => x.id !== id);
    },
  },
});
</script>
