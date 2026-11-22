import { useCallback, useEffect, useLayoutEffect, useRef } from "preact/hooks";
import { useSignal } from "@preact/signals";
import type { Item } from "../api-types";
import api from "../api";
import { useI18n } from "../i18n";
import type { AppSettings } from "../models/settings";
import type { AppRegistry } from "../models/registry";
import { dateTimeString, to } from "../utils";
import Dropdown from "./Dropdown";
import Icon from "./Icon";
import { showToast } from "./Toast";
import Tooltip from "./Tooltip";

export interface ColumnItemViewHandle {
  toggleReadability: () => void;
}

interface ColumnItemProps {
  item: Item | null;
  registry: AppRegistry;
  settings: AppSettings;
  handleRef: { current: ColumnItemViewHandle | null };
  onToggleStarred: (item: Item) => void;
  onToggleRead: (item: Item) => void;
  onNavigate: (dir: number) => void;
  canPrev: boolean;
  canNext: boolean;
  onClose: () => void;
}

export default function ColumnItemView({
  item,
  registry,
  settings,
  handleRef,
  onToggleStarred,
  onToggleRead,
  onNavigate,
  canPrev,
  canNext,
  onClose,
}: ColumnItemProps) {
  const { t } = useI18n();
  const isFullscreen = useSignal(false);
  const readabilityContent = useSignal("");
  const loadingReadability = useSignal(false);
  const contentRef = useRef<HTMLDivElement>(null);

  // latest article id, used to drop stale readability responses
  const itemIdRef = useRef<number | null>(null);
  itemIdRef.current = item?.id ?? null;

  // reset on close; survives article-to-article navigation like before
  useEffect(() => {
    if (!item) isFullscreen.value = false;
  }, [item]);

  // clear readability when the article changes; same id keeps content on star toggle
  useEffect(() => {
    readabilityContent.value = "";
    loadingReadability.value = false;
  }, [item?.id]);

  // start each article at the top
  useLayoutEffect(() => {
    if (contentRef.current) contentRef.current.scrollTop = 0;
  }, [item?.id]);

  const toggleReadability = useCallback(async () => {
    if (readabilityContent.value) {
      readabilityContent.value = "";
      return;
    }
    if (!item?.link) return;
    const targetId = item.id;
    loadingReadability.value = true;
    const [err, data] = await to(api.crawl(item.link));
    if (itemIdRef.current !== targetId) return;
    loadingReadability.value = false;
    if (err) showToast({ title: "fail_readability", description: err }, { level: "fail" });
    else readabilityContent.value = data?.content || "";
  }, [item]);

  useEffect(() => {
    handleRef.current = { toggleReadability };
    return () => {
      handleRef.current = null;
    };
  }, [handleRef, toggleReadability]);

  if (!item)
    return (
      <div id="col-item" className="vh-100 d-flex flex-column w-100" style={{ minWidth: 0 }} />
    );

  const feedId = item.feed_id;
  const feedTitle = registry.feedsById.value[feedId]?.title || "";

  const themeFont = settings.values.theme_font.value;
  const themeSize = settings.values.theme_size.value;
  const language = settings.values.language.value;

  const contentImages = (item.media_links || []).filter(l => l.type === "image");
  const contentAudios = (item.media_links || []).filter(l => l.type === "audio");
  const contentVideos = (item.media_links || []).filter(l => l.type === "video");
  const htmlContent = readabilityContent.value || item.content || "";

  const fontClass =
    themeFont === "" ? "font-sans-serif" : themeFont === "serif" ? "font-serif" : "font-monospace";

  return (
    <div
      id="col-item"
      className={`vh-100 d-flex flex-column w-100 ${isFullscreen.value ? "is-fullscreen" : ""}`}
      style={{ minWidth: 0 }}>
      <div className="px-2 py-1 d-flex gap-1 align-items-center">
        <Tooltip label={t("mark_starred")}>
          <button
            type="button"
            className="c-button-pill"
            aria-label={t("mark_starred")}
            onClick={() => onToggleStarred(item)}>
            <Icon name={item.status === "starred" ? "star-full" : "star"} />
          </button>
        </Tooltip>
        <Tooltip label={t("mark_unread")}>
          <button
            type="button"
            className="c-button-pill"
            aria-label={t("mark_unread")}
            onClick={() => onToggleRead(item)}>
            <Icon name={item.status === "unread" ? "circle-full" : "circle"} />
          </button>
        </Tooltip>

        <Dropdown
          toggleClass="px-2"
          drop="center"
          title={t("appearance")}
          button={<Icon name="sliders" />}>
          <button
            type="button"
            className="c-dropdown-item w-100 text-start font-sans-serif"
            aria-pressed={themeFont === ""}
            onClick={e => {
              e.stopPropagation();
              settings.update({ theme_font: "" });
            }}>
            {t("sans_serif")}
          </button>
          <button
            type="button"
            className="c-dropdown-item w-100 text-start font-serif"
            aria-pressed={themeFont === "serif"}
            onClick={e => {
              e.stopPropagation();
              settings.update({ theme_font: "serif" });
            }}>
            {t("serif")}
          </button>
          <button
            type="button"
            className="c-dropdown-item w-100 text-start font-monospace"
            aria-pressed={themeFont === "monospace"}
            onClick={e => {
              e.stopPropagation();
              settings.update({ theme_font: "monospace" });
            }}>
            {t("monospace")}
          </button>
          <div className="d-flex text-center">
            <button
              type="button"
              className="c-dropdown-item flex-fill"
              style={{ fontSize: "0.8rem" }}
              onClick={e => {
                e.stopPropagation();
                settings.update({ theme_size: +(themeSize - 0.1).toFixed(1) });
              }}>
              A
            </button>
            <button
              type="button"
              className="c-dropdown-item flex-fill"
              style={{ fontSize: "1.2rem" }}
              onClick={e => {
                e.stopPropagation();
                settings.update({ theme_size: +(themeSize + 0.1).toFixed(1) });
              }}>
              A
            </button>
          </div>
        </Dropdown>

        <Tooltip label={t("read_here")}>
          <button
            type="button"
            className="c-button-pill"
            aria-label={t("read_here")}
            aria-pressed={!!readabilityContent.value}
            onClick={toggleReadability}>
            <Icon className={loadingReadability.value ? "is-loading" : ""} name="book-open" />
          </button>
        </Tooltip>

        <Tooltip label={t("open_link")}>
          <a
            className="c-button-pill"
            href={item.link}
            rel="noopener noreferrer"
            target="_blank"
            referrerPolicy="no-referrer"
            aria-label={t("open_link")}>
            <Icon name="external-link" />
          </a>
        </Tooltip>

        <div className="flex-grow-1" />

        <Tooltip label={t("previous_article")}>
          <button
            type="button"
            className="c-button-pill"
            aria-label={t("previous_article")}
            onClick={() => onNavigate(-1)}
            disabled={!canPrev}>
            <Icon name="chevron-left" />
          </button>
        </Tooltip>
        <Tooltip label={t("next_article")}>
          <button
            type="button"
            className="c-button-pill"
            aria-label={t("next_article")}
            onClick={() => onNavigate(1)}
            disabled={!canNext}>
            <Icon name="chevron-right" />
          </button>
        </Tooltip>
        <Tooltip label={isFullscreen.value ? t("fullscreen_leave") : t("fullscreen_enter")}>
          <button
            type="button"
            className="c-button-pill d-none d-lg-flex"
            aria-label={isFullscreen.value ? t("fullscreen_leave") : t("fullscreen_enter")}
            onClick={() => (isFullscreen.value = !isFullscreen.value)}>
            <Icon name={isFullscreen.value ? "fullscreen-exit" : "fullscreen"} />
          </button>
        </Tooltip>
        <Tooltip label={t("close_article")}>
          <button
            type="button"
            className="c-button-pill"
            aria-label={t("close_article")}
            onClick={onClose}>
            <Icon name="x" />
          </button>
        </Tooltip>
      </div>

      <div
        ref={contentRef}
        className={`content px-4 pt-3 pb-5 border-top overflow-auto ${fontClass}`}
        style={{ fontSize: `${themeSize}rem` }}>
        <div className="content-wrapper">
          <h1>
            <b>{item.title || t("untitled")}</b>
          </h1>
          <div className="opacity-50">
            <div>
              <span
                className="cursor-pointer"
                onClick={() => feedId && settings.update({ feed: `feed:${feedId}` })}>
                {feedTitle}
              </span>
            </div>
            <time>{dateTimeString(new Date(item.date), language)}</time>
          </div>
          <hr />
          {!readabilityContent && (
            <div>
              {contentImages.length > 0 && (
                <div>
                  {contentImages.map((media, idx) => (
                    <figure key={idx}>
                      <img src={media.url} loading="lazy" />
                      {media.description && <figcaption>{media.description}</figcaption>}
                    </figure>
                  ))}
                </div>
              )}
              {contentAudios.map((media, idx) => (
                <audio key={idx} className="w-100" controls src={media.url} />
              ))}
              {contentVideos.map((media, idx) => (
                <video key={idx} className="w-100" controls src={media.url} />
              ))}
            </div>
          )}
          <div dangerouslySetInnerHTML={{ __html: htmlContent }} />
        </div>
      </div>
    </div>
  );
}
