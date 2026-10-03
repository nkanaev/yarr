export function scrollto(target: Element, scroll: Element) {
  var padding = 10;
  var targetRect = target.getBoundingClientRect();
  var scrollRect = scroll.getBoundingClientRect();

  // target
  var relativeOffset = targetRect.y - scrollRect.y;
  var absoluteOffset = relativeOffset + scroll.scrollTop;

  if (
    padding <= relativeOffset &&
    relativeOffset + targetRect.height <= scrollRect.height - padding
  )
    return;

  var newPos = scroll.scrollTop;
  if (relativeOffset < padding) {
    newPos = absoluteOffset - padding;
  } else {
    newPos = absoluteOffset - scrollRect.height + targetRect.height + padding;
  }
  scroll.scrollTop = Math.round(newPos);
}

// TODO: replace or rewrite debounce/$debounce.
// `debounce` has type checker issues with Options API.
// `$debounce` initialises callback once, capturing only the first variables/arguments.

export function debounce<T extends (...args: any[]) => any>(fn: T, timeout: number) {
  let timerId: ReturnType<typeof setTimeout> | null = null;
  return function (...args: Parameters<T>): void {
    if (timerId) clearTimeout(timerId);
    timerId = setTimeout(() => {
      fn(...args);
    }, timeout);
  };
}

const debounceCache = new WeakMap();

export const debounceMixin = {
  methods: {
    $debounce(id: string, fn: (...args: any[]) => any, delay = 300) {
      // 'this' inside mixin methods automatically refers to the Vue component instance
      let keywordMap = debounceCache.get(this);

      if (!keywordMap) {
        keywordMap = {};
        debounceCache.set(this, keywordMap);
      }

      if (!keywordMap[id]) {
        keywordMap[id] = debounce(fn, delay);
      }

      keywordMap[id]();
    },
  },
};

export type RelUnit = "minute" | "hour" | "day";

export function dateRepr(d: Date, locale?: string): string {
  var sec = (new Date().getTime() - d.getTime()) / 1000;
  var neg = sec < 0;
  var out = "";

  sec = Math.abs(sec);
  if (sec < 2700)
    // less than 45 minutes
    out = formatRel(Math.round(sec / 60), "minute", locale);
  else if (sec < 86400)
    // less than 24 hours
    out = formatRel(Math.round(sec / 3600), "hour", locale);
  else if (sec < 604800)
    // less than a week
    out = formatRel(Math.round(sec / 86400), "day", locale);
  else
    return dateString(d, locale);

  if (neg) return "-" + out;
  return out;
}

export function dateString(d: Date, locale?: string): string {
  if (locale == "en") locale = "en-GB";  // empire strikes back
  return d.toLocaleDateString(locale, {
    year: "numeric",
    month: "short",
    day: "numeric",
  })
}

export function dateTimeString(d: Date, locale: string): string {
  if (locale == "en") locale = "en-GB";  // empire strikes back
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

// returns ms until relative text rendered by dateRepr would change,
// or null when static (older than a week).
export function relRepaintDelay(d: Date): number | null {
  const sec = (Date.now() - d.getTime()) / 1000;
  if (isNaN(sec)) return null;

  if (sec >= 0) {
    let unit: number;
    if (sec < 2700) {
      // less than 45 minutes
      unit = 60;
    } else if (sec < 86400) {
      // less than 24 hours
      unit = 3600;
    } else if (sec < 604800) {
      // less than a week
      unit = 86400;
    } else {
      // older than a week
      return null;
    }

    const next = (Math.round(sec / unit) + 0.5) * unit;
    const delay = next - sec;
    return Math.max(delay * 1000, 1000);
  } else {
    // future dates
    if (Math.abs(sec) >= 604800) {
      return null;
    }
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
