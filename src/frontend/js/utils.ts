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

export function dateRepr(d: Date): string {
  var sec = (new Date().getTime() - d.getTime()) / 1000;
  var neg = sec < 0;
  var out = "";

  sec = Math.abs(sec);
  if (sec < 2700)
    // less than 45 minutes
    out = Math.round(sec / 60) + "m";
  else if (sec < 86400)
    // less than 24 hours
    out = Math.round(sec / 3600) + "h";
  else if (sec < 604800)
    // less than a week
    out = Math.round(sec / 86400) + "d";
  else
    out = d.toLocaleDateString(undefined, {
      year: "numeric",
      month: "long",
      day: "numeric",
    });

  if (neg) return "-" + out;
  return out;
}

async function to<T, E = Error>(promise: Promise<T>): Promise<[E, undefined] | [undefined, T]> {
  try {
    const result = await promise;
    return [undefined, result];
  } catch (err) {
    return [err as E, undefined];
  }
}
