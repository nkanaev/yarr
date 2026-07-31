<template>
  <div class="d-flex flex-column gap-1">
    <div
      class="c-listitem d-flex user-select-none gap-2"
      role="radio"
      :aria-checked="modelValue === ''"
      @click="$emit('update:modelValue', '')">
      <div class="flex-shrink-0">
        <v-icon class="flex-shrink-0" name="layers" />
      </div>
      <div class="flex-grow-1 min-w-0 text-truncate">
        {{
          { unread: $t("all_unread"), starred: $t("all_starred"), "": $t("all_feeds") }[
            filterSelected
          ]
        }}
      </div>
      <div class="flex-shrink-0" v-if="filterSelected">
        <span class="ps-2 text-end opacity-50">{{ stats.total[filterSelected] }}</span>
      </div>
    </div>
    <template
      v-for="node in tree"
      :key="node.type === 'folder' ? 'folder:' + node.folder.id : 'feed:' + node.feed.id">
      <template v-if="node.type === 'folder'">
        <div
          class="c-listitem d-flex user-select-none gap-2"
          role="radio"
          :aria-checked="modelValue === 'folder:' + node.folder.id"
          @click="$emit('update:modelValue', 'folder:' + node.folder.id)">
          <div class="flex-shrink-0">
            <div @click.stop="$emit('toggle-folder', node.folder)" class="p-2 m-n2">
              <v-icon name="chevron-right" v-if="!node.folder.is_expanded" />
              <v-icon name="chevron-down" v-else-if="node.folder.is_expanded" />
            </div>
          </div>
          <div class="flex-grow-1 min-w-0 text-truncate">
            {{ node.folder.title }}
          </div>
          <div class="flex-shrink-0" v-if="filterSelected">
            <span class="ps-2 text-end opacity-50">{{
              stats.folders[node.folder.id]?.[filterSelected]
            }}</span>
          </div>
        </div>
        <div
          v-if="node.folder.is_expanded && node.feeds.length > 0"
          class="d-flex flex-column gap-1 ps-3">
          <div
            v-for="feedNode in node.feeds"
            :key="'feed:' + feedNode.feed.id"
            class="c-listitem d-flex user-select-none gap-2"
            role="radio"
            :aria-checked="modelValue === 'feed:' + feedNode.feed.id"
            @click="$emit('update:modelValue', 'feed:' + feedNode.feed.id)">
            <div class="flex-shrink-0">
              <v-icon class="flex-shrink-0" name="rss" v-if="!feedNode.feed.icon" />
              <span class="c-icon" v-else>
                <img :src="feedNode.feed.icon" alt="" loading="lazy" />
              </span>
            </div>
            <div class="flex-grow-1 min-w-0 text-truncate">
              {{ feedNode.feed.title }}
            </div>
            <div class="flex-shrink-0">
              <span class="ps-2 text-end opacity-50" v-if="filterSelected">{{
                stats.feeds[feedNode.feed.id]?.[filterSelected]
              }}</span>
              <v-icon
                class="flex-shrink-0"
                :title="feedErrors[feedNode.feed.id]"
                v-if="!filterSelected && feedErrors[feedNode.feed.id]"
                name="alert-circle" />
            </div>
          </div>
        </div>
      </template>
      <template v-else-if="node.type === 'feed'">
        <div
          class="c-listitem d-flex user-select-none gap-2"
          role="radio"
          :aria-checked="modelValue === 'feed:' + node.feed.id"
          @click="$emit('update:modelValue', 'feed:' + node.feed.id)">
          <div class="flex-shrink-0">
            <v-icon class="flex-shrink-0" name="rss" v-if="!node.feed.icon" />
            <span class="c-icon" v-else>
              <img :src="node.feed.icon" alt="" loading="lazy" />
            </span>
          </div>
          <div class="flex-grow-1 min-w-0 text-truncate">
            {{ node.feed.title }}
          </div>
          <div class="flex-shrink-0">
            <span class="ps-2 text-end opacity-50" v-if="filterSelected">{{
              stats.feeds[node.feed.id]?.[filterSelected]
            }}</span>
            <v-icon
              class="flex-shrink-0"
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
    filterSelected: { type: String as PropType<"" | "unread" | "starred">, required: true },
    stats: {
      type: Object as PropType<{
        feeds: Record<number, { unread: number; starred: number }>;
        folders: Record<number, { unread: number; starred: number }>;
        total: { unread: number; starred: number };
      }>,
      required: true,
    },
    feedErrors: { type: Object as PropType<Record<number, string>>, required: true },
  },
  emits: ["update:modelValue", "toggle-folder"],
});
</script>
