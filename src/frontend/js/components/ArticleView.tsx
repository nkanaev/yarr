import type { Ref } from "preact";
import type { Item } from "../api-types";
import { useI18n } from "../i18n";
import { dateTimeString } from "../utils";
import Dropdown from "./Dropdown";
import Icon from "./Icon";
import Tooltip from "./Tooltip";

export type ThemeFont = "" | "serif" | "monospace";

interface ArticleViewProps {
  item: Item | null;
  feedTitle: string;
  feedId?: number;
  theme: { font: ThemeFont; size: number };
  onFontChange: (font: ThemeFont) => void;
  onFontSizeChange: (delta: number) => void;
  readabilityContent: string;
  loadingReadability: boolean;
  onToggleReadability: () => void;
  onToggleStarred: (item: Item) => void;
  onToggleRead: (item: Item) => void;
  onSelectFeed: (feedId: number) => void;
  onNavigate: (dir: number) => void;
  canPrev: boolean;
  canNext: boolean;
  isFullscreen: boolean;
  onToggleFullscreen: () => void;
  onClose: () => void;
  contentRef: Ref<HTMLDivElement>;
  language: string;
}

export default function ArticleView({
  item,
  feedTitle,
  feedId,
  theme,
  onFontChange,
  onFontSizeChange,
  readabilityContent,
  loadingReadability,
  onToggleReadability,
  onToggleStarred,
  onToggleRead,
  onSelectFeed,
  onNavigate,
  canPrev,
  canNext,
  isFullscreen,
  onToggleFullscreen,
  onClose,
  contentRef,
  language,
}: ArticleViewProps) {
  const { t } = useI18n();

  if (!item) return null;

  const contentImages = (item.media_links || []).filter(l => l.type === "image");
  const contentAudios = (item.media_links || []).filter(l => l.type === "audio");
  const contentVideos = (item.media_links || []).filter(l => l.type === "video");
  const htmlContent = readabilityContent || item.content || "";

  const fontClass =
    theme.font === ""
      ? "font-sans-serif"
      : theme.font === "serif"
        ? "font-serif"
        : "font-monospace";

  return (
    <>
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
            aria-pressed={theme.font === ""}
            onClick={e => {
              e.stopPropagation();
              onFontChange("");
            }}>
            {t("sans_serif")}
          </button>
          <button
            type="button"
            className="c-dropdown-item w-100 text-start font-serif"
            aria-pressed={theme.font === "serif"}
            onClick={e => {
              e.stopPropagation();
              onFontChange("serif");
            }}>
            {t("serif")}
          </button>
          <button
            type="button"
            className="c-dropdown-item w-100 text-start font-monospace"
            aria-pressed={theme.font === "monospace"}
            onClick={e => {
              e.stopPropagation();
              onFontChange("monospace");
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
                onFontSizeChange(-1);
              }}>
              A
            </button>
            <button
              type="button"
              className="c-dropdown-item flex-fill"
              style={{ fontSize: "1.2rem" }}
              onClick={e => {
                e.stopPropagation();
                onFontSizeChange(1);
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
            aria-pressed={!!readabilityContent}
            onClick={onToggleReadability}>
            <Icon className={loadingReadability ? "is-loading" : ""} name="book-open" />
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
        <Tooltip label={isFullscreen ? t("fullscreen_leave") : t("fullscreen_enter")}>
          <button
            type="button"
            className="c-button-pill d-none d-lg-flex"
            aria-label={isFullscreen ? t("fullscreen_leave") : t("fullscreen_enter")}
            onClick={onToggleFullscreen}>
            <Icon name={isFullscreen ? "fullscreen-exit" : "fullscreen"} />
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
        style={{ fontSize: `${theme.size}rem` }}>
        <div className="content-wrapper">
          <h1>
            <b>{item.title || t("untitled")}</b>
          </h1>
          <div className="opacity-50">
            <div>
              <span className="cursor-pointer" onClick={() => feedId && onSelectFeed(feedId)}>
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
    </>
  );
}
