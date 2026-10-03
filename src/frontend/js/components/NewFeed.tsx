import { useState, useEffect, useRef } from "preact/hooks";
import api, { NetworkError, HTTPError } from "../api";
import { to } from "../utils";
import { useI18n } from "../i18n";
import { showToast } from "./Toast";
import type { Folder, Feed, FeedLink, FeedCreateData } from "../api-types";

interface NewFeedProps {
  folders: Folder[];
  folderId: number | null;
  onCreated: (feed: Feed) => void;
  onFolderCreated: (folder: Folder) => void;
}

export default function NewFeed({ folders, folderId, onCreated, onFolderCreated }: NewFeedProps) {
  const { t } = useI18n();
  const [selectedFolder, setSelectedFolder] = useState<number | null>(folderId);
  const [feedNewChoice, setFeedNewChoice] = useState<FeedLink[]>([]);
  const [feedNewChoiceSelected, setFeedNewChoiceSelected] = useState("");
  const [loading, setLoading] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    setSelectedFolder(folderId);
  }, [folderId]);

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
      folder_id: selectedFolder,
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

  async function createNewFeedFolder() {
    const title = prompt(t("prompt_folder_name"));
    if (!title) return;

    const [folderErr, result] = await to(api.folders.create({ title }));
    if (folderErr) {
      showToast(
        { title: t("fail_save_folder"), description: errDescription(folderErr) },
        { level: "fail", closeable: false },
      );
      return;
    }

    onFolderCreated(result);
    setSelectedFolder(result.id);
  }

  function resetFeedChoice() {
    setFeedNewChoice([]);
    setFeedNewChoiceSelected("");
  }

  return (
    <div className="d-flex flex-column">
      <p className="cursor-default mb-3">
        <b>{t("new_feed")}</b>
      </p>
      <form onSubmit={createFeed} className="d-flex flex-column">
        <label htmlFor="feed-url" className="mb-2">
          {t("url")}
        </label>
        <input
          ref={inputRef}
          id="feed-url"
          name="url"
          type="url"
          className="c-input"
          required
          autoComplete="off"
          readOnly={feedNewChoice.length > 0}
          placeholder="https://example.com/feed"
        />

        <label htmlFor="feed-folder" className="mb-2 mt-3">
          {t("folder")}
          <a
            href="#"
            className="float-end text-decoration-none"
            onClick={e => {
              e.preventDefault();
              createNewFeedFolder();
            }}>
            {t("new_folder")}
          </a>
        </label>
        <select
          className="c-input"
          id="feed-folder"
          name="folder_id"
          value={selectedFolder === null ? "" : selectedFolder}
          onChange={e => {
            const v = (e.target as HTMLSelectElement).value;
            setSelectedFolder(v === "" ? null : Number(v));
          }}>
          <option value="">---</option>
          {folders.map(folder => (
            <option key={folder.id} value={folder.id}>
              {folder.title}
            </option>
          ))}
        </select>

        {feedNewChoice.length > 0 && (
          <div className="mt-3">
            <p className="mb-2">
              {t("multiple_feeds_found")}
              <a
                href="#"
                className="float-end text-decoration-none"
                onClick={e => {
                  e.preventDefault();
                  resetFeedChoice();
                }}>
                {t("cancel")}
              </a>
            </p>
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

        <button className="c-button mt-3" disabled={loading} type="submit">
          {loading ? <span className="c-spinner" /> : <span>{t("add")}</span>}
        </button>
      </form>
    </div>
  );
}
