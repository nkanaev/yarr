import type { Ref } from "preact";
import type { Item } from "../api-types";
import { useI18n } from "../i18n";
import type { AppSettings } from "../models/settings";
import { dateTimeString } from "../utils";
import Dropdown from "./Dropdown";
import Icon from "./Icon";
import Tooltip from "./Tooltip";

interface ArticleViewProps {
  item: Item | null;
  feedTitle: string;
  feedId?: number;
  settings: AppSettings;
  readabilityContent: string;
  loadingReadability: boolean;
  onToggleReadability: () => void;
  onToggleStarred: (item: Item) => void;
  onToggleRead: (item: Item) => void;
  onNavigate: (dir: number) => void;
  canPrev: boolean;
  canNext: boolean;
  isFullscreen: boolean;
  onToggleFullscreen: () => void;
  onClose: () => void;
  contentRef: Ref<HTMLDivElement>;
}

export default function ArticleView({
  item,
  feedTitle,
  feedId,
  settings,
  readabilityContent,
  loadingReadability,
  onToggleReadability,
  onToggleStarred,
  onToggleRead,
  onNavigate,
  canPrev,
  canNext,
  isFullscreen,
  onToggleFullscreen,
  onClose,
  contentRef,
}: ArticleViewProps) {
  const { t } = useI18n();

  if (!item) return null;

  const themeFont = settings.values.theme_font.value;
  const themeSize = settings.values.theme_size.value;
  const language = settings.values.language.value;

  const contentImages = (item.media_links || []).filter(l => l.type === "image");
  const contentAudios = (item.media_links || []).filter(l => l.type === "audio");
  const contentVideos = (item.media_links || []).filter(l => l.type === "video");
  const htmlContent = readabilityContent || item.content || "";

  const fontClass =
    themeFont === ""
      ? "font-sans-serif"
      : themeFont === "serif"
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
    </>
  );
}
