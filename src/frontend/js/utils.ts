export function scrollto(target: Element, scroll: Element) {
  const padding = 10;
  const targetRect = target.getBoundingClientRect();
  const scrollRect = scroll.getBoundingClientRect();

  const relativeOffset = targetRect.y - scrollRect.y;
  const absoluteOffset = relativeOffset + scroll.scrollTop;

  if (
    padding <= relativeOffset &&
    relativeOffset + targetRect.height <= scrollRect.height - padding
  ) {
    return;
  }

  let newPos = scroll.scrollTop;
  if (relativeOffset < padding) {
    newPos = absoluteOffset - padding;
  } else {
    newPos = absoluteOffset - scrollRect.height + targetRect.height + padding;
  }
  scroll.scrollTop = Math.round(newPos);
}

export function debounce<T extends (...args: any[]) => any>(fn: T, timeout: number) {
  let timerId: ReturnType<typeof setTimeout> | null = null;
  return function (...args: Parameters<T>): void {
    if (timerId) clearTimeout(timerId);
    timerId = setTimeout(() => {
      fn(...args);
    }, timeout);
  };
}

export type RelUnit = "minute" | "hour" | "day";

export function dateRepr(d: Date, locale?: string): string {
  let sec = (new Date().getTime() - d.getTime()) / 1000;
  const neg = sec < 0;
  let out = "";

  sec = Math.abs(sec);
  if (sec < 2700) {
    out = formatRel(Math.round(sec / 60), "minute", locale);
  } else if (sec < 86400) {
    out = formatRel(Math.round(sec / 3600), "hour", locale);
  } else if (sec < 604800) {
    out = formatRel(Math.round(sec / 86400), "day", locale);
  } else {
    return dateString(d, locale);
  }

  if (neg) return "-" + out;
  return out;
}

export function dateString(d: Date, locale?: string): string {
  if (locale === "en") locale = "en-GB";
  return d.toLocaleDateString(locale, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

export function dateTimeString(d: Date, locale: string): string {
  if (locale === "en") locale = "en-GB";
  return d.toLocaleDateString(locale, {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function formatRel(n: number, unit: RelUnit, locale?: string): string {
  return new Intl.NumberFormat(locale, { style: "unit", unit }).format(n);
}

export function relRepaintDelay(d: Date): number | null {
  const sec = (Date.now() - d.getTime()) / 1000;
  if (isNaN(sec)) return null;

  if (sec >= 0) {
    let unit: number;
    if (sec < 2700) unit = 60;
    else if (sec < 86400) unit = 3600;
    else if (sec < 604800) unit = 86400;
    else return null;

    const next = (Math.round(sec / unit) + 0.5) * unit;
    const delay = next - sec;
    return Math.max(delay * 1000, 1000);
  } else {
    if (Math.abs(sec) >= 604800) return null;
    return 60 * 1000;
  }
}

export async function to<T, E = Error>(
  promise: Promise<T>,
): Promise<[E, undefined] | [undefined, T]> {
  try {
    const result = await promise;
    return [undefined, result];
  } catch (err) {
    return [err as E, undefined];
  }
}
