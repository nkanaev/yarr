export interface KeyActions {
  openItemLink: () => void;
  toggleReadability: () => void;
  toggleItemRead: () => void;
  markAllRead: () => void;
  toggleItemStarred: () => void;
  nextItem: () => void;
  previousItem: () => void;
  nextFeed: () => void;
  previousFeed: () => void;
  closeItem: () => void;
  showShortcuts: () => void;
  showAll: () => void;
  showUnread: () => void;
  showStarred: () => void;
}

export function setupKeybindings(actionsRef: { current: KeyActions }) {
  const scrollContent = (direction: number) => {
    const padding = 40;
    const scroll = document.querySelector(".content");
    if (!scroll) return;

    const height = scroll.getBoundingClientRect().height;
    const newpos = scroll.scrollTop + (height - padding) * direction;

    if (typeof scroll.scrollTo === "function") {
      scroll.scrollTo({ top: newpos, left: 0, behavior: "smooth" });
    } else {
      scroll.scrollTop = newpos;
    }
  };

  const focusSearch = () => {
    document.getElementById("searchbar")?.focus();
  };

  const shortcutFunctions = {
    openItemLink: () => actionsRef.current.openItemLink(),
    toggleReadability: () => actionsRef.current.toggleReadability(),
    toggleItemRead: () => actionsRef.current.toggleItemRead(),
    markAllRead: () => actionsRef.current.markAllRead(),
    toggleItemStarred: () => actionsRef.current.toggleItemStarred(),
    focusSearch,
    nextItem: () => actionsRef.current.nextItem(),
    previousItem: () => actionsRef.current.previousItem(),
    nextFeed: () => actionsRef.current.nextFeed(),
    previousFeed: () => actionsRef.current.previousFeed(),
    scrollForward: () => scrollContent(+1),
    scrollBackward: () => scrollContent(-1),
    closeItem: () => actionsRef.current.closeItem(),
    showShortcuts: () => actionsRef.current.showShortcuts(),
    showAll: () => actionsRef.current.showAll(),
    showUnread: () => actionsRef.current.showUnread(),
    showStarred: () => actionsRef.current.showStarred(),
  };

  const keybindings: Record<string, () => void> = {
    o: shortcutFunctions.openItemLink,
    i: shortcutFunctions.toggleReadability,
    r: shortcutFunctions.toggleItemRead,
    R: shortcutFunctions.markAllRead,
    s: shortcutFunctions.toggleItemStarred,
    "/": shortcutFunctions.focusSearch,
    j: shortcutFunctions.nextItem,
    k: shortcutFunctions.previousItem,
    l: shortcutFunctions.nextFeed,
    h: shortcutFunctions.previousFeed,
    f: shortcutFunctions.scrollForward,
    b: shortcutFunctions.scrollBackward,
    q: shortcutFunctions.closeItem,
    "?": shortcutFunctions.showShortcuts,
    "1": shortcutFunctions.showUnread,
    "2": shortcutFunctions.showStarred,
    "3": shortcutFunctions.showAll,
  };

  const codebindings: Record<string, () => void> = {
    KeyO: shortcutFunctions.openItemLink,
    KeyI: shortcutFunctions.toggleReadability,
    KeyS: shortcutFunctions.toggleItemStarred,
    Slash: shortcutFunctions.focusSearch,
    KeyJ: shortcutFunctions.nextItem,
    KeyK: shortcutFunctions.previousItem,
    KeyL: shortcutFunctions.nextFeed,
    KeyH: shortcutFunctions.previousFeed,
    KeyF: shortcutFunctions.scrollForward,
    KeyB: shortcutFunctions.scrollBackward,
    KeyQ: shortcutFunctions.closeItem,
    Digit1: shortcutFunctions.showUnread,
    Digit2: shortcutFunctions.showStarred,
    Digit3: shortcutFunctions.showAll,
  };

  function isTextBox(element: Element) {
    const tagName = element.tagName.toLowerCase();
    // Input elements that aren't text
    const inputBlocklist = [
      "button",
      "checkbox",
      "color",
      "file",
      "hidden",
      "image",
      "radio",
      "range",
      "reset",
      "submit",
    ];

    return (
      tagName === "textarea" ||
      (tagName === "input" &&
        inputBlocklist.indexOf(element.getAttribute("type")?.toLowerCase() || "") === -1)
    );
  }

  const handler = (event: KeyboardEvent) => {
    if (isTextBox(event.target as Element) || event.metaKey || event.ctrlKey || event.altKey) {
      return;
    }
    const keybindFunction = keybindings[event.key] || codebindings[event.code];
    if (keybindFunction) {
      event.preventDefault();
      keybindFunction();
    }
  };

  document.addEventListener("keydown", handler);
  return () => document.removeEventListener("keydown", handler);
}
