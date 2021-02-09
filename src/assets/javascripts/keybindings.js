const helperFunctions = {
  // navigation helper, navigate relative to selected item
  navigateToItem: function(relativePosition) {
    if(vm.itemSelected == null){
      if(vm.items.length !== 0) {
        // if no item is selected, select first
        vm.itemSelected = vm.items[0].id
      }
      return
    }
    const itemPosition = vm.items.findIndex(x=>x.id==vm.itemSelected)
    if(itemPosition == -1){
      // Item not found error
      return
    }
    const newPosition = itemPosition+relativePosition
    if(newPosition < 0 || newPosition >= vm.items.length){
      return
    }
    vm.itemSelected = vm.items[newPosition].id
  },
  // navigation helper, navigate relative to selected feed
  navigateToFeed: function(relativePosition) {
    // create a list with feed and folders guids, ignore feeds in collapsed folders
    // Example result with folder 2 collapsed:
    // ['','folder:1','feed:1','feed:2','folder:2', 'folder:3','feed:3']
    // The empty string is the "All Feeds" option
    const navigationList = [''].concat(vm.foldersWithFeeds.map(
      folder =>
        folder.is_expanded
        ? ['folder:'+folder.id].concat(folder.feeds.map(feed=>'feed:'+feed.id))
        : 'folder:'+folder.id
    ).flat())
    const currentFeedPosition = navigationList.indexOf(vm.feedSelected)
    if(currentFeedPosition== -1){
      // feed not found error
      return
    }
    const newPosition = currentFeedPosition+relativePosition
    if(newPosition < 0 || newPosition >= navigationList.length){
      return
    }
    vm.feedSelected = navigationList[newPosition];
  }
}
const shortcutFunctions = {
  toggleItemRead: function() {
    if(vm.itemSelected != null) {
      vm.toggleItemRead(vm.itemSelectedDetails)
    }
  },
  markAllRead: function() {
    // same condition as 'Mark all read button'
    if(vm.filterSelected == 'unread'){
      vm.markItemsRead()
    }
  },
  toggleItemStarred: function() {
    if(vm.itemSelected != null) {
      vm.toggleItemStarred(vm.itemSelectedDetails)
    }
  },
  focusSearch: function() {
    document.getElementById("searchbar").focus()
  },
  nextItem(){
    helperFunctions.navigateToItem(+1)
  },
  previousItem() {
    helperFunctions.navigateToItem(-1)
  },
  nextFeed(){
    helperFunctions.navigateToFeed(+1)
  },
  previousFeed() {
    helperFunctions.navigateToFeed(-1)
  },
  showAll() {
    vm.filterSelected = ''
  },
  showUnread() {
    vm.filterSelected = 'unread'
  },
  showStarred() {
    vm.filterSelected = 'starred'
  },
}

const keybindings = {
  "r": shortcutFunctions.toggleItemRead,
  "R": shortcutFunctions.markAllRead,
  "s": shortcutFunctions.toggleItemStarred,
  "?": shortcutFunctions.focusSearch,
  "j": shortcutFunctions.nextItem,
  "k": shortcutFunctions.previousItem,
  "l": shortcutFunctions.nextFeed,
  "h": shortcutFunctions.previousFeed,
  "A": shortcutFunctions.showAll,
  "U": shortcutFunctions.showUnread,
  "S": shortcutFunctions.showStarred,
}

function isTextBox(element) {
  var tagName = element.tagName.toLowerCase()
  // Input elements that aren't text
  const inputBlocklist = ['button','checkbox','color','file','hidden','image','radio','range','reset','search','submit']

  return tagName === 'textarea' ||
    ( tagName === 'input'
      && !inputBlocklist.includes(element.getAttribute('type').toLowerCase())
    )
}

document.addEventListener('keydown',function(event) {
  // Ignore while focused on text or
  // when using modifier keys (to not clash with browser behaviour)
  if(isTextBox(event.target) || event.metaKey || event.ctrlKey) {
    return
  }
  const keybindFunction = keybindings[event.key]
  if(keybindFunction) {
    event.preventDefault()
    keybindFunction()
  }
})
