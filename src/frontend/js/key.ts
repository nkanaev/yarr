export function setupKeybindings(vm) {
  var helperFunctions = {
    scrollContent(direction: number) {
      var padding = 40;
      var scroll = document.querySelector(".content");
      if (!scroll) return;

      var height = scroll.getBoundingClientRect().height;
      var newpos = scroll.scrollTop + (height - padding) * direction;

      if (typeof scroll.scrollTo == "function") {
        scroll.scrollTo({ top: newpos, left: 0, behavior: "smooth" });
      } else {
        scroll.scrollTop = newpos;
      }
    },
  };
  var shortcutFunctions = {
    openItemLink() {
      if (vm.itemSelectedDetails && vm.itemSelectedDetails.link) {
        window.open(
          vm.itemSelectedDetails.link,
          "_blank",
          "noopener,noreferrer",
        );
      }
    },
    toggleReadability() {
      vm.toggleReadability();
    },
    toggleItemRead() {
      if (vm.itemSelected != null) {
        vm.toggleItemRead(vm.itemSelectedDetails);
      }
    },
    markAllRead() {
      // same condition as 'Mark all read button'
      if (vm.filterSelected == "unread") {
        vm.markItemsRead();
      }
    },
    toggleItemStarred() {
      if (vm.itemSelected != null) {
        vm.toggleItemStarred(vm.itemSelectedDetails);
      }
    },
    focusSearch() {
      document.getElementById("searchbar").focus();
    },
    nextItem() {
      vm.navigateToItem(+1);
    },
    previousItem() {
      vm.navigateToItem(-1);
    },
    nextFeed() {
      vm.navigateToFeed(+1);
    },
    previousFeed() {
      vm.navigateToFeed(-1);
    },
    scrollForward() {
      helperFunctions.scrollContent(+1);
    },
    scrollBackward() {
      helperFunctions.scrollContent(-1);
    },
    closeItem() {
      vm.itemSelected = null;
    },
    showAll() {
      vm.filterSelected = "";
    },
    showUnread() {
      vm.filterSelected = "unread";
    },
    showStarred() {
      vm.filterSelected = "starred";
    },
  };

  // If you edit, make sure you update the help modal
  var keybindings = {
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
    "1": shortcutFunctions.showUnread,
    "2": shortcutFunctions.showStarred,
    "3": shortcutFunctions.showAll,
  };

  var codebindings = {
    KeyO: shortcutFunctions.openItemLink,
    KeyI: shortcutFunctions.toggleReadability,
    //"r": shortcutFunctions.toggleItemRead,
    //"KeyR": shortcutFunctions.markAllRead,
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
    var tagName = element.tagName.toLowerCase();
    // Input elements that aren't text
    var inputBlocklist = [
      "button",
      "checkbox",
      "color",
      "file",
      "hidden",
      "image",
      "radio",
      "range",
      "reset",
      "search",
      "submit",
    ];

    return (
      tagName === "textarea" ||
      (tagName === "input" &&
        inputBlocklist.indexOf(element.getAttribute("type").toLowerCase()) ==
          -1)
    );
  }

  document.addEventListener("keydown", function (event) {
    // Ignore while focused on text or
    // when using modifier keys (to not clash with browser behaviour)
    if (
      isTextBox(event.target as Element) ||
      event.metaKey ||
      event.ctrlKey ||
      event.altKey
    ) {
      return;
    }
    var keybindFunction = keybindings[event.key] || codebindings[event.code];
    if (keybindFunction) {
      event.preventDefault();
      keybindFunction();
    }
  });
}
