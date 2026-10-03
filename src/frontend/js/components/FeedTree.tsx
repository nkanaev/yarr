import type { Folder, Feed } from "../api-types";
import { useI18n } from "../i18n";
import Icon from "./Icon";

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

interface FeedTreeProps {
  tree: FeedTreeNode[];
  value: string | null;
  filterSelected: "" | "unread" | "starred";
  stats: {
    feeds: Record<number, { unread: number; starred: number }>;
    folders: Record<number, { unread: number; starred: number }>;
    total: { unread: number; starred: number };
  };
  feedErrors: Record<number, string>;
  onSelect: (val: string) => void;
  onToggleFolder: (folder: Folder) => void;
}

export default function FeedTree({
  tree,
  value,
  filterSelected,
  stats,
  feedErrors,
  onSelect,
  onToggleFolder,
}: FeedTreeProps) {
  const { t } = useI18n();

  const allTitle =
    filterSelected === "unread"
      ? t("all_unread")
      : filterSelected === "starred"
        ? t("all_starred")
        : t("all_feeds");

  return (
    <div className="d-flex flex-column gap-1">
      <div
        className="c-listitem d-flex user-select-none gap-2"
        role="radio"
        aria-checked={value === ""}
        onClick={() => onSelect("")}>
        <div className="flex-shrink-0 d-flex">
          <Icon name="layers" />
        </div>
        <div className="flex-grow-1 min-w-0 text-truncate">{allTitle}</div>
        {filterSelected && (
          <div className="flex-shrink-0">
            <span className="ps-2 text-end opacity-50">{stats.total[filterSelected]}</span>
          </div>
        )}
      </div>

      {tree.map(node => {
        if (node.type === "folder") {
          const folderKey = "folder:" + node.folder.id;
          const isSelected = value === folderKey;
          return (
            <div key={folderKey}>
              <div
                className="c-listitem d-flex user-select-none gap-2"
                role="radio"
                aria-checked={isSelected}
                onClick={() => onSelect(folderKey)}>
                <div className="flex-shrink-0 d-flex">
                  <div
                    onClick={e => {
                      e.stopPropagation();
                      onToggleFolder(node.folder);
                    }}
                    className="p-2 m-n2 d-inline-flex">
                    <Icon name={node.folder.is_expanded ? "chevron-down" : "chevron-right"} />
                  </div>
                </div>
                <div className="flex-grow-1 min-w-0 text-truncate">{node.folder.title}</div>
                {filterSelected && (
                  <div className="flex-shrink-0">
                    <span className="ps-2 text-end opacity-50">
                      {stats.folders[node.folder.id]?.[filterSelected]}
                    </span>
                  </div>
                )}
              </div>
              {node.folder.is_expanded && node.feeds.length > 0 && (
                <div className="d-flex flex-column gap-1 ps-3 pt-1">
                  {node.feeds.map(feedNode => {
                    const feedKey = "feed:" + feedNode.feed.id;
                    return (
                      <div
                        key={feedKey}
                        className="c-listitem d-flex user-select-none gap-2"
                        role="radio"
                        aria-checked={value === feedKey}
                        onClick={() => onSelect(feedKey)}>
                        <div className="flex-shrink-0 d-flex align-items-center">
                          {!feedNode.feed.icon ? (
                            <Icon className="flex-shrink-0" name="rss" />
                          ) : (
                            <span className="c-icon">
                              <img src={feedNode.feed.icon} alt="" loading="lazy" />
                            </span>
                          )}
                        </div>
                        <div className="flex-grow-1 min-w-0 text-truncate">
                          {feedNode.feed.title}
                        </div>
                        <div className="flex-shrink-0 d-flex">
                          {filterSelected && (
                            <span className="ps-2 text-end opacity-50">
                              {stats.feeds[feedNode.feed.id]?.[filterSelected]}
                            </span>
                          )}
                          {!filterSelected && feedErrors[feedNode.feed.id] && (
                            <Icon
                              className="flex-shrink-0"
                              title={feedErrors[feedNode.feed.id]}
                              name="alert-circle"
                            />
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </div>
          );
        }

        const feedKey = "feed:" + node.feed.id;
        return (
          <div
            key={feedKey}
            className="c-listitem d-flex user-select-none gap-2"
            role="radio"
            aria-checked={value === feedKey}
            onClick={() => onSelect(feedKey)}>
            <div className="flex-shrink-0 d-flex">
              {!node.feed.icon ? (
                <Icon className="flex-shrink-0" name="rss" />
              ) : (
                <span className="c-icon">
                  <img src={node.feed.icon} alt="" loading="lazy" />
                </span>
              )}
            </div>
            <div className="flex-grow-1 min-w-0 text-truncate">{node.feed.title}</div>
            <div className="flex-shrink-0 d-flex">
              {filterSelected && (
                <span className="ps-2 text-end opacity-50">
                  {stats.feeds[node.feed.id]?.[filterSelected]}
                </span>
              )}
              {!filterSelected && feedErrors[node.feed.id] && (
                <Icon
                  className="flex-shrink-0"
                  title={feedErrors[node.feed.id]}
                  name="alert-circle"
                />
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}
