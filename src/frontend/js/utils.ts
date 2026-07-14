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

export function debounce<A extends unknown[]>(
  callback: (...args: A) => void,
  wait: number,
): (...args: A) => void {
  let timeout: ReturnType<typeof setTimeout> | undefined;
  return function (this: unknown, ...args: A) {
    const ctx = this;
    clearTimeout(timeout);
    timeout = setTimeout(function () {
      callback.apply(ctx, args);
    }, wait);
  };
}

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
