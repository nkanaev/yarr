<template>
  <div>
    <div
      class="selectgroup"
      role="radio"
      :aria-checked="modelValue === ''"
      @click="$emit('update:modelValue', '')">
      <div class="selectgroup-label d-flex align-items-center w-100">
        <v-icon class="mr-2" name="layers" />
        <span class="flex-fill text-left text-truncate" v-if="filterSelected == 'unread'">{{
          $t("all_unread")
        }}</span>
        <span class="flex-fill text-left text-truncate" v-if="filterSelected == 'starred'">{{
          $t("all_starred")
        }}</span>
        <span class="flex-fill text-left text-truncate" v-if="filterSelected == ''">{{
          $t("all_feeds")
        }}</span>
        <span class="counter text-right">{{ filteredTotalStats }}</span>
      </div>
    </div>
    <template
      v-for="node in tree"
      :key="node.type === 'folder' ? 'folder:' + node.folder.id : 'feed:' + node.feed.id">
      <template v-if="node.type === 'folder'">
        <div
          class="selectgroup mt-1"
          role="radio"
          :aria-checked="modelValue === 'folder:' + node.folder.id"
          @click="$emit('update:modelValue', 'folder:' + node.folder.id)">
          <div class="selectgroup-label d-flex align-items-center w-100">
            <div @click.stop="$emit('toggle-folder', node.folder)" class="m-n1 p-1">
              <v-icon
                class="mr-2"
                :class="{ expanded: node.folder.is_expanded }"
                name="chevron-right" />
            </div>
            <span class="flex-fill text-left text-truncate">{{ node.folder.title }}</span>
            <span class="counter text-right">{{ filteredFolderStats[node.folder.id] || "" }}</span>
          </div>
        </div>
        <div v-show="node.folder.is_expanded" class="mt-1 pl-3">
          <div
            class="selectgroup"
            role="radio"
            :aria-checked="modelValue === 'feed:' + feedNode.feed.id"
            @click="$emit('update:modelValue', 'feed:' + feedNode.feed.id)"
            v-for="feedNode in node.feeds">
            <div class="selectgroup-label d-flex align-items-center w-100">
              <v-icon class="mr-2" name="rss" v-if="!feedNode.feed.icon" />
              <span class="icon mr-2" v-else
                ><img :src="feedNode.feed.icon" alt="" loading="lazy"
              /></span>
              <span class="flex-fill text-left text-truncate">{{ feedNode.feed.title }}</span>
              <span class="counter text-right">{{
                filteredFeedStats[feedNode.feed.id] || ""
              }}</span>
              <v-icon
                class="flex-shrink-0 mx-2"
                :title="feedErrors[feedNode.feed.id]"
                v-if="!filterSelected && feedErrors[feedNode.feed.id]"
                name="alert-circle" />
            </div>
          </div>
        </div>
      </template>
      <template v-else-if="node.type === 'feed'">
        <div
          class="selectgroup"
          role="radio"
          :aria-checked="modelValue === 'feed:' + node.feed.id"
          @click="$emit('update:modelValue', 'feed:' + node.feed.id)">
          <div class="selectgroup-label d-flex align-items-center w-100">
            <v-icon class="mr-2" name="rss" v-if="!node.feed.icon" />
            <span class="icon mr-2" v-else
              ><img :src="node.feed.icon" alt="" loading="lazy"
            /></span>
            <span class="flex-fill text-left text-truncate">{{ node.feed.title }}</span>
            <span class="counter text-right">{{ filteredFeedStats[node.feed.id] || "" }}</span>
            <v-icon
              class="flex-shrink-0 mx-2"
              :title="feedErrors[node.feed.id]"
              v-if="!filterSelected && feedErrors[node.feed.id]"
              name="alert-circle" />
          </div>
        </div>
      </template>
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import type { PropType } from "vue";
import type { Folder, Feed } from "../api-types";
import icon from "../components/icon.vue";

export interface TreeFeedNode {
  type: "feed";
  feed: Feed;
}
export interface TreeFolderNode {
  type: "folder";
  folder: Folder;
  feeds: TreeFeedNode[];
}
export type FeedTreeNode = TreeFolderNode | TreeFeedNode;

export default defineComponent({
  components: { "v-icon": icon },
  props: {
    tree: { type: Array as PropType<FeedTreeNode[]>, required: true },
    modelValue: { type: String, required: true },
    filterSelected: { type: String, required: true },
    filteredTotalStats: { type: Number as PropType<number | null>, default: null },
    filteredFeedStats: { type: Object as PropType<Record<number, number>>, required: true },
    filteredFolderStats: { type: Object as PropType<Record<string, number>>, required: true },
    feedErrors: { type: Object as PropType<Record<number, string>>, required: true },
  },
  emits: ["update:modelValue", "toggle-folder"],
  computed: {},
});
</script>
