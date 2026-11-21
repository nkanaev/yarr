import { useState, useEffect } from "preact/hooks";
import { useI18n } from "../i18n";
import type { TranslationKey, TranslateFn } from "../i18n";
import { NetworkError, HTTPError } from "../api";
import Icon from "./Icon";

export type ToastLevel = "info" | "fail";

export interface ToastOptions {
  level?: ToastLevel;
  closeable?: boolean;
  err?: unknown;
}

export interface ToastMessage {
  title: TranslationKey;
  description?: TranslationKey | Error;
}

export interface ToastItem {
  id: string;
  message: ToastMessage;
  opts?: ToastOptions;
}

const getUUID = () =>
  crypto?.randomUUID?.() ?? Date.now().toString(36) + Math.random().toString(36).slice(2);

type Listener = (toasts: ToastItem[]) => void;
let toasts: ToastItem[] = [];
const listeners = new Set<Listener>();

function notify() {
  listeners.forEach(l => l([...toasts]));
}

export function showToast(message: ToastMessage, opts?: ToastOptions) {
  const item: ToastItem = { id: getUUID(), message, opts };
  toasts = [...toasts, item];
  notify();
  setTimeout(() => removeToast(item.id), 3000);
}

export function removeToast(id: string) {
  toasts = toasts.filter(t => t.id !== id);
  notify();
}

function describeToast(description: ToastMessage["description"], t: TranslateFn) {
  if (typeof description === "string") return t(description);
  if (description instanceof HTTPError)
    return t("error_server", { code: description.status, text: description.statusText });
  if (description instanceof NetworkError) return t("error_network");
  return undefined;
}

export default function ToastContainer() {
  const { t } = useI18n();
  const [list, setList] = useState<ToastItem[]>(toasts);

  useEffect(() => {
    listeners.add(setList);
    return () => {
      listeners.delete(setList);
    };
  }, []);

  return (
    <div className="position-fixed bottom-0 end-0 pb-3 pe-3 d-flex flex-column align-items-end gap-2 z-4 pe-none">
      {list.map(toast => {
        const description = describeToast(toast.message.description, t);
        return (
          <div
            key={toast.id}
            className={`c-toast d-flex align-items-start gap-2 ${
              toast.opts?.level === "fail" ? "text-bg-danger text-light" : ""
            }`}>
            {toast.opts?.level === "fail" && <Icon name="alert-circle" />}
            <div className="flex-grow-1">
              <div>{t(toast.message.title)}</div>
              {description && <div className="opacity-50">{description}</div>}
            </div>
            {toast.opts?.closeable === true && (
              <button type="button" className="c-button-link" onClick={() => removeToast(toast.id)}>
                <Icon name="x" />
              </button>
            )}
          </div>
        );
      })}
    </div>
  );
}
