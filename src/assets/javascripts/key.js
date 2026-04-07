'use strict';

// Keyboard shortcuts for yarr (HTMX version)
// No dependency on Vue — operates directly on DOM + HTMX

(function() {

  function scrollContent(direction) {
    var scroll = document.getElementById('article-content-scroll');
    if (!scroll) return;
    var height = scroll.getBoundingClientRect().height;
    var padding = 40;
    var newpos = scroll.scrollTop + (height - padding) * direction;
    if (typeof scroll.scrollTo === 'function') {
      scroll.scrollTo({ top: newpos, left: 0, behavior: 'smooth' });
    } else {
      scroll.scrollTop = newpos;
    }
  }

  function getFeedRadios() {
    return Array.from(document.querySelectorAll('#feed-list-content input[name=feed]'));
  }

  function navigateFeed(offset) {
    var radios = getFeedRadios().filter(function(r) {
      // Skip hidden feeds
      var label = r.closest('.selectgroup');
      return label && !label.classList.contains('d-none') && label.style.display !== 'none';
    });
    if (radios.length === 0) return;

    var current = document.querySelector('#feed-list-content input[name=feed]:checked');
    var idx = current ? radios.indexOf(current) : -1;
    var next = idx + offset;

    if (next < 0 || next >= radios.length) return;

    radios[next].checked = true;
    radios[next].dispatchEvent(new Event('change', { bubbles: true }));

    // Update app state
    document.getElementById('app').classList.add('feed-selected');
    document.getElementById('app').classList.remove('item-selected');

    var label = radios[next].closest('.selectgroup');
    if (label) label.scrollIntoView({ block: 'nearest' });
  }

  var shortcuts = {
    'o': function() {
      var link = document.querySelector('#item-content-inner a[rel="noopener noreferrer"]');
      if (link) window.open(link.href, '_blank', 'noopener,noreferrer');
    },
    'i': function() {
      var btn = document.getElementById('btn-readability');
      if (btn) btn.click();
    },
    'r': function() {
      var btn = document.querySelector('#item-content-inner .toolbar-item[title*="Unread"]');
      if (btn) btn.click();
    },
    'R': function() {
      var btn = document.getElementById('btn-mark-read');
      if (btn && btn.style.display !== 'none') btn.click();
    },
    's': function() {
      var btn = document.querySelector('#item-content-inner .toolbar-item[title*="Star"]');
      if (btn) btn.click();
    },
    '/': function() {
      var search = document.getElementById('searchbar');
      if (search) search.focus();
    },
    'j': function() { if (window.yarr && yarr.selectArticleAtOffset) yarr.selectArticleAtOffset(1); },
    'k': function() { if (window.yarr && yarr.selectArticleAtOffset) yarr.selectArticleAtOffset(-1); },
    'l': function() { navigateFeed(1); },
    'h': function() { navigateFeed(-1); },
    'f': function() { scrollContent(1); },
    'b': function() { scrollContent(-1); },
    'q': function() {
      document.getElementById('app').classList.remove('item-selected');
      document.getElementById('col-item-content').innerHTML = '';
    },
    '1': function() {
      var btn = document.querySelector('[data-filter="unread"]');
      if (btn) btn.click();
    },
    '2': function() {
      var btn = document.querySelector('[data-filter="starred"]');
      if (btn) btn.click();
    },
    '3': function() {
      var btn = document.querySelector('[data-filter=""]');
      if (btn) btn.click();
    },
  };

  // Code-based bindings for non-QWERTY layouts
  var codeBindings = {
    'KeyO': shortcuts['o'],
    'KeyI': shortcuts['i'],
    'KeyS': shortcuts['s'],
    'Slash': shortcuts['/'],
    'KeyJ': shortcuts['j'],
    'KeyK': shortcuts['k'],
    'KeyL': shortcuts['l'],
    'KeyH': shortcuts['h'],
    'KeyF': shortcuts['f'],
    'KeyB': shortcuts['b'],
    'KeyQ': shortcuts['q'],
    'Digit1': shortcuts['1'],
    'Digit2': shortcuts['2'],
    'Digit3': shortcuts['3'],
  };

  function isTextInput(el) {
    var tag = el.tagName.toLowerCase();
    if (tag === 'textarea') return true;
    if (tag === 'input') {
      var type = (el.getAttribute('type') || 'text').toLowerCase();
      var blocked = ['button','checkbox','color','file','hidden','image','radio','range','reset','submit'];
      return blocked.indexOf(type) === -1;
    }
    return el.isContentEditable;
  }

  document.addEventListener('keydown', function(e) {
    // Skip in text inputs or with modifiers
    if (isTextInput(e.target) || e.metaKey || e.ctrlKey || e.altKey) return;
    // Skip in dialogs
    if (e.target.closest('dialog')) return;

    var fn = shortcuts[e.key] || codeBindings[e.code];
    if (fn) {
      e.preventDefault();
      fn();
    }
  });

})();
