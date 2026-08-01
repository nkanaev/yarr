<template>
  <div class="d-flex flex-column">
    <p class="cursor-default mb-3">
      <b>{{ $t("new_feed") }}</b>
    </p>
    <form @submit.prevent="createFeed($event)" class="d-flex flex-column">
      <label for="feed-url" class="mb-2">{{ $t("url") }}</label>
      <input
        id="feed-url"
        name="url"
        type="url"
        class="c-input"
        required
        autocomplete="off"
        :readonly="feedNewChoice.length > 0"
        placeholder="https://example.com/feed"
        v-focus />
      <label for="feed-folder" class="mb-2 mt-3">
        {{ $t("folder") }}
        <a
          href="#"
          class="float-end text-decoration-none"
          @click.prevent="createNewFeedFolder()"
          >{{ $t("new_folder") }}</a
        >
      </label>
      <select class="c-input" id="feed-folder" name="folder_id" v-model="selectedFolder">
        <option :value="null">---</option>
        <option :value="folder.id" v-for="folder in folders">
          {{ folder.title }}
        </option>
      </select>
      <div class="mt-3" v-if="feedNewChoice.length">
        <p class="mb-2">
          {{ $t("multiple_feeds_found") }}
          <a
            href="#"
            class="float-end text-decoration-none"
            @click.prevent="resetFeedChoice()"
            >{{ $t("cancel") }}</a
          >
        </p>
        <div class="d-flex flex-column gap-1">
          <div
            class="c-listitem d-flex flex-column user-select-none"
            role="radio"
            :aria-checked="feedNewChoiceSelected === choice.url"
            @click="feedNewChoiceSelected = choice.url"
            v-for="choice in feedNewChoice">
            <div class="text-truncate">{{ choice.title }}</div>
            <div class="text-truncate" :class="{ 'opacity-50': choice.title }">
              {{ choice.url }}
            </div>
          </div>
        </div>
      </div>
      <button class="c-button mt-3" :disabled="loading.newfeed" type="submit">
        <span class="c-spinner" v-if="loading.newfeed"></span>
        <span v-else>{{ $t("add") }}</span>
      </button>
    </form>
    <v-toast ref="toast" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import type { PropType } from "vue";
import api, { NetworkError, HTTPError } from "../api";
import { to } from "../utils";
import focusDir from "../directives/focus";
import icon from "./icon.vue";
import toast from "./toast.vue";
import type { Folder, Feed, FeedLink, FeedCreateData } from "../api-types";

export default defineComponent({
  components: {
    "v-icon": icon,
    "v-toast": toast,
  },
  directives: {
    focus: focusDir,
  },
  props: {
    folders: { type: Array as PropType<Folder[]>, required: true },
    folderId: { type: Number as PropType<number | null>, default: null },
  },
  emits: ["created", "folder-created"],
  data() {
    return {
      selectedFolder: this.$props.folderId as number | null,
      feedNewChoice: [] as FeedLink[],
      feedNewChoiceSelected: "",
      loading: { newfeed: false },
    };
  },
  watch: {
    folderId(val) {
      this.selectedFolder = val;
    },
  },
  methods: {
    toast(): InstanceType<typeof toast> {
      return this.$refs.toast as InstanceType<typeof toast>;
    },
    async createFeed($event: Event) {
      var form = $event.target as HTMLFormElement;
      var data: FeedCreateData = {
        url: (form.querySelector("input[name=url]") as HTMLInputElement).value,
        folder_id: this.selectedFolder,
      };
      if (this.feedNewChoiceSelected) {
        var choice = this.feedNewChoice.find(c => c.url === this.feedNewChoiceSelected);
        data.url = this.feedNewChoiceSelected;
        if (choice && choice.title_override) data.title_override = choice.title_override;
      }
      this.loading.newfeed = true;
      const [err, result] = await to(api.feeds.create(data));
      this.loading.newfeed = false;
      if (err) {
        this.toast().addToast(
          { title: this.$t("fail_save_feed"), description: this.errDescription(err) },
          { level: "fail", closeable: false },
        );
        return;
      }
      if (result.status === "success") {
        this.$emit("created", result.feed);
      } else if (result.status === "multiple") {
        this.feedNewChoice = result.choice;
        this.feedNewChoiceSelected = result.choice[0].url;
      } else {
        alert("No feeds found at the given url.");
      }
    },
    async createNewFeedFolder() {
      const title = prompt(this.$t("prompt_folder_name"));
      if (!title) return;

      const [folderErr, result] = await to(api.folders.create({ title: title }));
      if (folderErr) {
        this.toast().addToast(
          { title: this.$t("fail_save_folder"), description: this.errDescription(folderErr) },
          { level: "fail", closeable: false },
        );
        return;
      }

      this.$emit("folder-created", result);
      this.selectedFolder = result.id;
    },
    resetFeedChoice() {
      this.feedNewChoice = [];
      this.feedNewChoiceSelected = "";
    },
    errDescription(err: unknown): string | undefined {
      if (err instanceof HTTPError)
        return this.$t("error_server", { code: err.status, text: err.statusText });
      if (err instanceof NetworkError) return this.$t("error_network");
      return undefined;
    },
  },
});
</script>
