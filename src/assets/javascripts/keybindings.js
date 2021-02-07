const keybindings = {}

function isTextBox(element) {
  var tagName = element.tagName.toLowerCase();
  // Input elements that aren't text
  const inputBlocklist = ['button','checkbox','color','file','hidden','image','radio','range','reset','search','submit']

  return tagName === 'textarea' ||
    ( tagName === 'input'
      && !inputBlocklist.includes(element.getAttribute('type').toLowerCase())
    )
}

document.addEventListener('keydown',function(event) {
  if(isTextBox(event.target)) {
    return;
  }
  const keybindFunction = keybindings[event.key];
  if(keybindFunction) {
    event.preventDefault();
    keybindFunction();
  }
})
