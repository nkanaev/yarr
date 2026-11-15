import { useState, useEffect, useRef } from "preact/hooks";
import api, { NetworkError, HTTPError } from "../api";
import { to } from "../utils";
import { useI18n } from "../i18n";
import { showToast } from "./Toast";
import type { Feed, FeedLink, FeedCreateData } from "../api-types";

interface NewFeedProps {
  folderId: number | null;
  onCancel: () => void;
  onCreated: (feed: Feed) => void;
}

export default function NewFeed({ folderId, onCancel, onCreated }: NewFeedProps) {
  const { t } = useI18n();
  const [feedNewChoice, setFeedNewChoice] = useState<FeedLink[]>([]);
  const [feedNewChoiceSelected, setFeedNewChoiceSelected] = useState("");
  const [loading, setLoading] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  function errDescription(err: unknown): string | undefined {
    if (err instanceof HTTPError)
      return t("error_server", { code: err.status, text: err.statusText });
    if (err instanceof NetworkError) return t("error_network");
    return undefined;
  }

  async function createFeed(event: Event) {
    event.preventDefault();
    const form = event.target as HTMLFormElement;
    const urlInput = form.querySelector("input[name=url]") as HTMLInputElement;

    const data: FeedCreateData = {
      url: urlInput.value,
      title: "",
      folder_id: folderId,
    };

    if (feedNewChoiceSelected) {
      const choice = feedNewChoice.find(c => c.url === feedNewChoiceSelected);
      data.url = feedNewChoiceSelected;
      data.title = choice?.title || "";
    }

    setLoading(true);
    const [err, result] = await to(api.feeds.create(data));
    setLoading(false);

    if (err) {
      showToast(
        { title: t("fail_save_feed"), description: errDescription(err) },
        { level: "fail", closeable: false },
      );
      return;
    }

    if (result.status === "success") {
      onCreated(result.feed);
    } else if (result.status === "multiple") {
      setFeedNewChoice(result.choice);
      setFeedNewChoiceSelected(result.choice[0].url);
    } else {
      alert("No feeds found at the given url.");
    }
  }

  return (
    <div className="d-flex flex-column">
      <form onSubmit={createFeed} className="d-flex flex-column">
        <label htmlFor="feed-url" className="mb-2 user-select-none">
          {t("new_feed_url")}
        </label>
        <input
          ref={inputRef}
          id="feed-url"
          name="url"
          type="url"
          className="c-input"
          required
          autoComplete="off"
          placeholder="https://example.com/feed"
          onInput={() => {
            if (feedNewChoice.length > 0) {
              setFeedNewChoice([]);
              setFeedNewChoiceSelected("");
            }
          }}
        />

        {feedNewChoice.length > 0 && (
          <div className="mt-3">
            <p className="mb-2">{t("multiple_feeds_found")}</p>
            <div className="d-flex flex-column gap-1">
              {feedNewChoice.map(choice => (
                <div
                  key={choice.url}
                  className="c-listitem d-flex flex-column user-select-none"
                  role="radio"
                  aria-checked={feedNewChoiceSelected === choice.url}
                  onClick={() => setFeedNewChoiceSelected(choice.url)}>
                  <div className="text-truncate">{choice.title}</div>
                  <div className={`text-truncate ${choice.title ? "opacity-50" : ""}`}>
                    {choice.url}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        <div className="d-flex gap-2 mt-3">
          <button type="button" className="c-button is-link w-50" onClick={onCancel}>
            {t("cancel")}
          </button>
          <button type="submit" className="c-button w-50" disabled={loading}>
            {loading ? <span className="c-spinner" /> : <span>{t("add")}</span>}
          </button>
        </div>
      </form>
    </div>
  );
}
