'use strict';

// yarr - vanilla JS app layer for HTMX-driven UI
window.yarr = window.yarr || {};

(function(yarr) {
  var TITLE = document.title;
  var settings = yarr.settings || {};
  var _confirmCb = null;
  var _promptCb = null;
  var _pollTimer = null;

  // --- API helper ---
  function api(method, endpoint, data) {
    var opts = {
      method: method,
      headers: {'Content-Type': 'application/json', 'x-requested-by': 'yarr'}
    };
    if (data) opts.body = JSON.stringify(data);
    return fetch(endpoint, opts);
  }

  // --- Toast ---
  yarr.toast = function(message, type) {
    var container = document.getElementById('toast-container');
    if (!container) return;
    var el = document.createElement('div');
    el.className = 'toast' + (type === 'error' ? ' toast-error' : '');
    el.textContent = message;
    container.appendChild(el);
    setTimeout(function() {
      el.classList.add('toast-out');
      setTimeout(function() { el.remove(); }, 200);
    }, 3000);
  };

  // --- Theme ---
  var themeColors = { light: '#fff', sepia: '#f4f0e5', night: '#0e0e0e' };

  yarr.setTheme = function(name) {
    document.body.setAttribute('data-theme', name);
    document.querySelector("meta[name='theme-color']").content = themeColors[name] || '#fff';
    // Update theme dots
    document.querySelectorAll('[data-theme-dot]').forEach(function(dot) {
      dot.classList.toggle('active', dot.getAttribute('data-theme-dot') === name);
    });
    settings.theme_name = name;
    api('put', './api/settings', { theme_name: name });
  };

  // --- Filter ---
  yarr.setFilter = function(btn, filter) {
    // Update active state on filter buttons
    btn.closest('.toolbar').querySelectorAll('.toolbar-item[data-filter]').forEach(function(b) {
      b.classList.toggle('active', b.getAttribute('data-filter') === filter);
    });
    // Update hidden filter input
    var hidden = document.getElementById('item-filter-status');
    if (hidden) hidden.value = filter;
    // Show/hide mark-read button
    var markRead = document.getElementById('btn-mark-read');
    if (markRead) markRead.style.display = filter === 'unread' ? '' : 'none';
    // Clear selection
    document.getElementById('app').classList.remove('item-selected');
    document.getElementById('col-item-content').innerHTML = '';
    // Persist
    settings.filter = filter;
    api('put', './api/settings', { filter: filter });
  };

  // --- Sort ---
  yarr.setSortOrder = function(newestFirst) {
    settings.sort_newest_first = newestFirst;
    api('put', './api/settings', { sort_newest_first: newestFirst }).then(function() {
      // Reload items
      htmx.ajax('GET', './partials/items', { target: '#item-list-content', swap: 'innerHTML' });
    });
  };

  // --- Refresh Rate ---
  yarr.setRefreshRate = function(val) {
    val = parseInt(val, 10);
    settings.refresh_rate = val;
    api('put', './api/settings', { refresh_rate: val });
  };

  // --- Column Resize ---
  function initDrag(handleId, minW, maxW, settingKey) {
    var handle = document.getElementById(handleId);
    if (!handle) return;
    var startX, initW, col;
    handle.addEventListener('mousedown', function(e) {
      startX = e.clientX;
      col = handle.parentElement;
      initW = col.offsetWidth;
      var onMove = function(e) {
        var w = Math.min(Math.max(minW, initW + e.clientX - startX), maxW);
        col.style.width = w + 'px';
      };
      var onUp = function() {
        document.removeEventListener('mousemove', onMove);
        document.removeEventListener('mouseup', onUp);
        var w = parseInt(col.style.width, 10);
        if (w) {
          var update = {};
          update[settingKey] = w;
          api('put', './api/settings', update);
        }
      };
      document.addEventListener('mousemove', onMove);
      document.addEventListener('mouseup', onUp);
    });
  }

  // --- Dropdown ---
  yarr.toggleDropdown = function(id) {
    var el = document.getElementById(id);
    var menu = el.querySelector('.dropdown-menu');
    var isOpen = menu.classList.contains('show');
    // Close all dropdowns first
    document.querySelectorAll('.dropdown-menu.show').forEach(function(m) { m.classList.remove('show'); });
    if (!isOpen) {
      menu.classList.add('show');
      // Close on click outside
      setTimeout(function() {
        var handler = function(e) {
          if (!el.contains(e.target) || e.target.closest('.dropdown-item')) {
            menu.classList.remove('show');
            document.removeEventListener('click', handler);
          }
        };
        document.addEventListener('click', handler);
      }, 0);
    }
  };

  yarr.closeDropdown = function(id) {
    var el = document.getElementById(id);
    if (el) el.querySelector('.dropdown-menu').classList.remove('show');
  };

  // --- Dialogs ---
  yarr.showDialog = function(id) {
    var d = document.getElementById(id);
    if (d) d.showModal();
  };

  yarr.confirm = function(message, callback) {
    document.getElementById('confirm-message').textContent = message;
    _confirmCb = callback;
    document.getElementById('confirm-dialog').showModal();
  };

  yarr.confirmAction = function() {
    if (_confirmCb) { _confirmCb(); _confirmCb = null; }
  };

  yarr.prompt = function(message, defaultVal, callback) {
    document.getElementById('prompt-message').textContent = message;
    document.getElementById('prompt-input').value = defaultVal || '';
    _promptCb = callback;
    document.getElementById('prompt-dialog').showModal();
    setTimeout(function() { document.getElementById('prompt-input').focus(); }, 50);
  };

  yarr.promptAction = function(val) {
    if (_promptCb) { _promptCb(val); _promptCb = null; }
  };

  // --- Feed Management ---
  yarr.createFeed = function(event) {
    event.preventDefault();
    var form = event.target;
    var url = form.querySelector('[name=url]').value;
    var folderId = form.querySelector('[name=folder_id]').value || null;

    var btn = document.getElementById('add-feed-btn');
    btn.classList.add('loading');

    api('post', './api/feeds', { url: url, folder_id: folderId ? parseInt(folderId) : null })
      .then(function(r) { return r.json(); })
      .then(function(result) {
        btn.classList.remove('loading');
        if (result.status === 'success') {
          document.getElementById('feed-dialog').close();
          yarr.toast('Feed added');
          // Reload feed list
          htmx.ajax('GET', './partials/feeds', { target: '#feed-list-content', swap: 'outerHTML' });
        } else if (result.status === 'multiple') {
          // Show choices
          var choiceDiv = document.getElementById('feed-choice');
          var html = '<p class="mb-2">Multiple feeds found. Choose one:</p>';
          result.choice.forEach(function(c) {
            html += '<label class="selectgroup"><input type="radio" name="feedToAdd" value="' + c.url + '"' +
                    (c.url === result.choice[0].url ? ' checked' : '') + '>' +
                    '<div class="selectgroup-label"><div class="text-truncate">' + (c.title || '') +
                    '</div><div class="text-truncate light">' + c.url + '</div></div></label>';
          });
          choiceDiv.innerHTML = html;
          choiceDiv.style.display = '';
          form.querySelector('[name=url]').readOnly = true;
        } else {
          yarr.toast('No feeds found at the given URL', 'error');
        }
      })
      .catch(function() {
        btn.classList.remove('loading');
        yarr.toast('Error adding feed', 'error');
      });
  };

  yarr.createFolderInline = function(event) {
    event.preventDefault();
    yarr.prompt('Enter folder name:', '', function(title) {
      if (!title) return;
      api('post', './api/folders', { title: title })
        .then(function(r) { return r.json(); })
        .then(function(folder) {
          var select = document.getElementById('feed-folder');
          var opt = document.createElement('option');
          opt.value = folder.id;
          opt.textContent = folder.title;
          opt.selected = true;
          select.appendChild(opt);
        });
    });
  };

  // --- Folder Expand/Collapse ---
  yarr.toggleFolder = function(id, expand) {
    api('put', './api/folders/' + id, { is_expanded: expand }).then(function() {
      htmx.ajax('GET', './partials/feeds', { target: '#feed-list-content', swap: 'outerHTML' });
    });
  };

  // --- Feed/Folder Context Menu ---
  function getSelectedFeed() {
    var checked = document.querySelector('#feed-list-content input[name=feed]:checked');
    return checked ? checked.value : '';
  }

  yarr.showFeedMenu = function() {
    var sel = getSelectedFeed();
    if (!sel) return;
    // Load menu content via fetch, then show dropdown
    fetch('./partials/feed-menu?sel=' + encodeURIComponent(sel), {
      headers: { 'HX-Request': 'true' }
    }).then(function(r) { return r.text(); }).then(function(html) {
      document.getElementById('feed-menu-content').innerHTML = html;
      yarr.toggleDropdown('feed-context-menu');
    });
  };

  yarr.renameFeed = function(id, currentTitle) {
    yarr.closeDropdown('feed-context-menu');
    yarr.prompt('Enter new title:', currentTitle, function(title) {
      if (!title) return;
      api('put', './api/feeds/' + id, { title: title }).then(function() {
        yarr.toast('Feed renamed');
        htmx.ajax('GET', './partials/feeds', { target: '#feed-list-content', swap: 'outerHTML' });
      });
    });
  };

  yarr.changeFeedLink = function(id, currentLink) {
    yarr.closeDropdown('feed-context-menu');
    yarr.prompt('Enter feed link:', currentLink, function(link) {
      if (!link) return;
      api('put', './api/feeds/' + id, { feed_link: link }).then(function() {
        yarr.toast('Feed link updated');
      });
    });
  };

  yarr.deleteFeed = function(id, title) {
    yarr.closeDropdown('feed-context-menu');
    yarr.confirm('Are you sure you want to delete ' + title + '?', function() {
      api('delete', './api/feeds/' + id).then(function() {
        yarr.toast('Feed deleted');
        // Deselect and reload
        var allRadio = document.querySelector('#feed-list-content input[name=feed][value=""]');
        if (allRadio) { allRadio.checked = true; allRadio.dispatchEvent(new Event('change', {bubbles: true})); }
        htmx.ajax('GET', './partials/feeds', { target: '#feed-list-content', swap: 'outerHTML' });
      });
    });
  };

  yarr.archiveFeed = function(id, archive) {
    yarr.closeDropdown('feed-context-menu');
    api('put', './api/feeds/' + id, { archived: archive }).then(function() {
      yarr.toast(archive ? 'Feed archived' : 'Feed unarchived');
      htmx.ajax('GET', './partials/feeds', { target: '#feed-list-content', swap: 'outerHTML' });
    });
  };

  yarr.moveFeed = function(feedId, folderId) {
    yarr.closeDropdown('feed-context-menu');
    api('put', './api/feeds/' + feedId, { folder_id: folderId }).then(function() {
      yarr.toast('Feed moved');
      htmx.ajax('GET', './partials/feeds', { target: '#feed-list-content', swap: 'outerHTML' });
    });
  };

  yarr.moveFeedToNewFolder = function(feedId) {
    yarr.closeDropdown('feed-context-menu');
    yarr.prompt('Enter folder name:', '', function(title) {
      if (!title) return;
      api('post', './api/folders', { title: title })
        .then(function(r) { return r.json(); })
        .then(function(folder) {
          return api('put', './api/feeds/' + feedId, { folder_id: folder.id });
        })
        .then(function() {
          yarr.toast('Feed moved to new folder');
          htmx.ajax('GET', './partials/feeds', { target: '#feed-list-content', swap: 'outerHTML' });
        });
    });
  };

  yarr.renameFolder = function(id, currentTitle) {
    yarr.closeDropdown('feed-context-menu');
    yarr.prompt('Enter new title:', currentTitle, function(title) {
      if (!title) return;
      api('put', './api/folders/' + id, { title: title }).then(function() {
        yarr.toast('Folder renamed');
        htmx.ajax('GET', './partials/feeds', { target: '#feed-list-content', swap: 'outerHTML' });
      });
    });
  };

  yarr.deleteFolder = function(id, title) {
    yarr.closeDropdown('feed-context-menu');
    yarr.confirm('Are you sure you want to delete ' + title + '?', function() {
      api('delete', './api/folders/' + id).then(function() {
        yarr.toast('Folder deleted');
        var allRadio = document.querySelector('#feed-list-content input[name=feed][value=""]');
        if (allRadio) { allRadio.checked = true; allRadio.dispatchEvent(new Event('change', {bubbles: true})); }
        htmx.ajax('GET', './partials/feeds', { target: '#feed-list-content', swap: 'outerHTML' });
      });
    });
  };

  // --- OPML ---
  yarr.importOPML = function(input) {
    var form = document.getElementById('opml-import-form');
    fetch('./opml/import', {
      method: 'post',
      headers: { 'x-requested-by': 'yarr' },
      body: new FormData(form)
    }).then(function() {
      input.value = '';
      yarr.toast('OPML imported');
      htmx.ajax('GET', './partials/feeds', { target: '#feed-list-content', swap: 'outerHTML' });
      yarr.pollStatus();
    });
  };

  // --- Status Polling ---
  yarr.pollStatus = function() {
    if (_pollTimer) return;
    var statusEl = document.getElementById('feed-status');
    var statusText = document.getElementById('feed-status-text');
    statusEl.style.display = '';
    var poll = function() {
      api('get', './api/status').then(function(r) { return r.json(); }).then(function(data) {
        if (data.running > 0) {
          statusText.textContent = 'Refreshing (' + data.running + ' left)';
          _pollTimer = setTimeout(poll, 500);
        } else {
          statusEl.style.display = 'none';
          _pollTimer = null;
          // Reload feed list to get updated counts
          htmx.ajax('GET', './partials/feeds', { target: '#feed-list-content', swap: 'outerHTML' });
        }
      });
    };
    poll();
  };

  // --- Density / Triage Mode ---
  yarr.toggleDensity = function() {
    var list = document.getElementById('item-list-scroll');
    if (list) list.classList.toggle('density-compact');
    var btn = document.getElementById('btn-density-toggle');
    if (btn) btn.classList.toggle('active');
  };

  // --- Reader Mode ---
  yarr.toggleReaderMode = function() {
    document.getElementById('app').classList.toggle('reader-mode');
  };

  // --- Topics View ---
  var _topicsActive = false;
  var _topicsLoaded = false;

  yarr.toggleTopicsView = function() {
    _topicsActive = !_topicsActive;
    var feedScroll = document.getElementById('feed-list-scroll');
    var topicsView = document.getElementById('topics-view');
    var btn = document.getElementById('btn-topics');

    if (_topicsActive) {
      feedScroll.style.display = 'none';
      topicsView.style.display = '';
      btn.classList.add('active');
      if (!_topicsLoaded) {
        yarr.loadTopics();
      }
    } else {
      feedScroll.style.display = '';
      topicsView.style.display = 'none';
      btn.classList.remove('active');
    }
  };

  yarr.loadTopics = function() {
    var container = document.getElementById('topics-view');
    container.innerHTML = '<div class="empty-state"><span class="icon loading"></span><p>Loading topics...</p></div>';

    // Fetch clusters, tags, and health in parallel
    Promise.all([
      fetch('./api/ai/clusters').then(function(r) { return r.json(); }).catch(function() { return { clusters: [] }; }),
      fetch('./api/ai/tags').then(function(r) { return r.json(); }).catch(function() { return []; }),
      fetch('./api/ai/health').then(function(r) { return r.json(); }).catch(function() { return {}; })
    ]).then(function(results) {
      var clusterData = results[0];
      var tagsData = results[1];
      var health = results[2];

      // Toolbar with action buttons
      var html = '<div class="topics-toolbar">'
        + '<button class="btn btn-default" onclick="yarr.reindexArticles(this)" title="Embed all articles into the AI index">'
        + 'Reindex' + (health.chroma_docs ? ' (' + health.chroma_docs + ')' : '') + '</button>'
        + '<button class="btn btn-default" onclick="yarr.rebuildTopics(this)" title="Run clustering to discover topics">'
        + 'Rebuild Topics</button>'
        + '</div>';

      // Clusters section
      if (clusterData.clusters && clusterData.clusters.length > 0) {
        html += '<div class="topic-section-header">Topics</div>';
        clusterData.clusters.forEach(function(c) {
          html += '<div class="topic-entry" data-topic-type="cluster" data-topic-tag="' + escapeAttr(c.label) + '" onclick="yarr.selectTopic(this)">'
            + '<span class="flex-fill text-truncate">' + escapeHtml(c.label) + '</span>'
            + '<span class="topic-badge">' + c.article_count + '</span>'
            + '</div>';
        });
      }

      // Tags section — only show tags that differ from cluster labels
      var clusterLabels = new Set((clusterData.clusters || []).map(function(c) { return c.label; }));
      var uniqueTags = (tagsData || []).filter(function(t) { return !clusterLabels.has(t.tag); });
      if (uniqueTags.length > 0) {
        html += '<div class="topic-section-header mt-3">Tags</div>';
        uniqueTags.slice(0, 30).forEach(function(t) {
          html += '<div class="topic-entry" data-topic-type="tag" data-topic-tag="' + escapeAttr(t.tag) + '" onclick="yarr.selectTopic(this)">'
            + '<span class="flex-fill text-truncate">' + escapeHtml(t.tag) + '</span>'
            + '<span class="topic-badge">' + t.article_count + '</span>'
            + '</div>';
        });
      }

      if (!clusterData.clusters?.length && !(tagsData && tagsData.length)) {
        html += '<div class="empty-state"><p>No topics yet. Reindex articles and rebuild topics.</p></div>';
      }

      container.innerHTML = html;
      _topicsLoaded = true;
    });
  };

  // --- AI Progress Bar ---
  var _aiPollTimer = null;
  var _aiActiveBtn = null;

  function showAiStatus(text) {
    var bar = document.getElementById('ai-status');
    var txt = document.getElementById('ai-status-text');
    txt.textContent = text;
    bar.style.display = '';
  }

  function hideAiStatus() {
    document.getElementById('ai-status').style.display = 'none';
    if (_aiPollTimer) { clearInterval(_aiPollTimer); _aiPollTimer = null; }
    if (_aiActiveBtn) {
      _aiActiveBtn.classList.remove('loading');
      _aiActiveBtn.disabled = false;
      _aiActiveBtn = null;
    }
  }

  function startTaskPoll() {
    if (_aiPollTimer) return;
    _aiPollTimer = setInterval(function() {
      fetch('./api/ai/task-status').then(function(r) { return r.json(); }).then(function(task) {
        if (task.type) {
          showAiStatus(task.detail || 'Processing...');
        } else {
          var lastDetail = document.getElementById('ai-status-text').textContent;
          hideAiStatus();
          if (lastDetail.indexOf('Complete') >= 0) {
            yarr.toast(lastDetail);
          } else {
            yarr.toast('AI task complete');
          }
          _topicsLoaded = false;
          yarr.loadTopics();
        }
      }).catch(function() {});
    }, 2000);
  }

  // Check for running AI task on page load (survives reload)
  yarr.checkTaskStatus = function() {
    fetch('./api/ai/task-status').then(function(r) { return r.json(); }).then(function(task) {
      if (task.type) {
        showAiStatus(task.detail || 'Processing...');
        startTaskPoll();
      }
    }).catch(function() {});
  };

  yarr.reindexArticles = function(btn) {
    btn.classList.add('loading');
    btn.disabled = true;
    _aiActiveBtn = btn;
    showAiStatus('Indexing: starting...');
    api('post', './api/ai/reindex').then(function() {
      startTaskPoll();
    });
  };

  yarr.rebuildTopics = function(btn) {
    btn.classList.add('loading');
    btn.disabled = true;
    _aiActiveBtn = btn;
    showAiStatus('Clustering: starting...');
    api('post', './api/ai/recluster').then(function() {
      startTaskPoll();
    });
  };

  yarr.selectTopic = function(el) {
    // Deselect all
    document.querySelectorAll('#topics-view .topic-entry.active').forEach(function(e) { e.classList.remove('active'); });
    el.classList.add('active');

    var tag = el.getAttribute('data-topic-tag');
    document.getElementById('app').classList.add('feed-selected');

    // Fetch articles for this tag
    fetch('./api/ai/articles?tag=' + encodeURIComponent(tag))
      .then(function(r) { return r.json(); })
      .then(function(articles) {
        var itemList = document.getElementById('item-list-content');
        if (!articles || articles.length === 0) {
          itemList.innerHTML = '<div class="empty-state flex-grow-1"><p>No articles for this topic</p></div>';
          return;
        }
        var html = '<div class="p-2 overflow-auto scroll-touch flex-grow-1" id="topic-articles-scroll" style="min-height:0">';
        articles.forEach(function(a) {
          html += '<div class="selectgroup topic-article-row" data-item-id="' + (a.id || 0) + '" data-item-url="' + escapeAttr(a.url || '') + '">'
            + '<div class="selectgroup-label d-flex flex-column cursor-pointer">'
            + '<div style="line-height:100%;opacity:.7;margin-bottom:.1rem" class="d-flex align-items-center">'
            + '<small class="flex-fill text-truncate me-1">' + escapeHtml((a.feed_name || '') + (a.folder ? ' · ' + a.folder : '')) + '</small>'
            + '<small class="flex-shrink-0">' + escapeHtml(dateRepr(a.published)) + '</small>'
            + '</div>'
            + '<div>' + escapeHtml(a.title || 'untitled') + '</div>'
            + '</div></div>';
        });
        html += '</div>';
        itemList.innerHTML = html;
      })
      .catch(function() {
        document.getElementById('item-list-content').innerHTML = '<div class="empty-state flex-grow-1"><p>Could not load articles</p></div>';
      });
  };

  var _months = ['January','February','March','April','May','June','July','August','September','October','November','December'];
  function dateRepr(dateStr) {
    if (!dateStr) return '';
    var d = new Date(dateStr);
    if (isNaN(d)) return dateStr;
    var sec = (Date.now() - d.getTime()) / 1000;
    var neg = sec < 0;
    if (neg) sec = -sec;
    var out;
    if (sec < 2700) out = Math.round(sec / 60) + 'm';
    else if (sec < 86400) out = Math.round(sec / 3600) + 'h';
    else if (sec < 604800) out = Math.round(sec / 86400) + 'd';
    else out = _months[d.getMonth()] + ' ' + d.getDate() + ', ' + d.getFullYear();
    return neg ? '-' + out : out;
  }

  function escapeHtml(s) {
    var d = document.createElement('div');
    d.textContent = s;
    return d.innerHTML;
  }

  function escapeAttr(s) {
    return s.replace(/&/g, '&amp;').replace(/'/g, '&#39;').replace(/"/g, '&quot;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  }

  // --- Search ---
  function initSearch() {
    var searchbar = document.getElementById('searchbar');
    var clearBtn = document.getElementById('search-clear');
    if (!searchbar || !clearBtn) return;
    searchbar.addEventListener('input', function() {
      clearBtn.style.display = this.value ? '' : 'none';
    });
  }

  // --- Reading Progress ---
  function initReadingProgress() {
    document.addEventListener('scroll', function(e) {
      var scroll = e.target;
      if (!scroll.id || scroll.id !== 'article-content-scroll') return;
      var bar = document.getElementById('reading-progress');
      if (!bar) return;
      var pct = scroll.scrollTop / (scroll.scrollHeight - scroll.clientHeight) * 100;
      bar.style.width = Math.min(100, Math.max(0, pct)) + '%';
    }, true);
  }

  // --- AI Chat Panel ---
  var _chatHistory = [];
  var _chatEventSource = null;

  yarr.openChat = function() {
    var panel = document.getElementById('chat-panel');
    panel.style.display = 'flex';
    document.getElementById('chat-input').focus();
  };

  yarr.closeChat = function() {
    document.getElementById('chat-panel').style.display = 'none';
    if (_chatEventSource) { _chatEventSource.close(); _chatEventSource = null; }
  };

  yarr.sendChat = function(event) {
    event.preventDefault();
    var input = document.getElementById('chat-input');
    var query = input.value.trim();
    if (!query) return;
    input.value = '';

    var messages = document.getElementById('chat-messages');

    // Add user message
    var userMsg = document.createElement('div');
    userMsg.className = 'chat-msg chat-msg-user';
    userMsg.textContent = query;
    messages.appendChild(userMsg);

    // Add assistant placeholder
    var assistantMsg = document.createElement('div');
    assistantMsg.className = 'chat-msg chat-msg-assistant streaming-cursor';
    messages.appendChild(assistantMsg);
    messages.scrollTop = messages.scrollHeight;

    // Stream response via SSE
    var body = JSON.stringify({ query: query, history: _chatHistory });
    fetch('./api/ai/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Accept': 'text/event-stream' },
      body: body
    }).then(function(response) {
      var reader = response.body.getReader();
      var decoder = new TextDecoder();
      var buffer = '';
      var fullResponse = '';

      function processChunk(result) {
        if (result.done) {
          assistantMsg.classList.remove('streaming-cursor');
          _chatHistory.push({ role: 'user', content: query });
          _chatHistory.push({ role: 'assistant', content: fullResponse });
          // Keep last 6 messages (3 turns)
          if (_chatHistory.length > 6) _chatHistory = _chatHistory.slice(-6);
          return;
        }
        buffer += decoder.decode(result.value, { stream: true });
        var lines = buffer.split('\n');
        buffer = lines.pop();

        for (var i = 0; i < lines.length; i++) {
          var line = lines[i];
          if (line.startsWith('event: sources')) {
            // Next data line has sources
            continue;
          }
          if (line.startsWith('data: ')) {
            var data = line.slice(6);
            if (data === '[DONE]') continue;
            // Check for sources JSON
            try {
              var parsed = JSON.parse(data);
              if (parsed.sources) {
                var srcDiv = document.createElement('div');
                srcDiv.className = 'chat-msg-sources';
                srcDiv.innerHTML = parsed.sources.map(function(s, idx) {
                  return '<div>[' + (idx+1) + '] <a href="' + s.url + '" target="_blank" rel="noopener">' + (s.title || s.url) + '</a></div>';
                }).join('');
                assistantMsg.appendChild(srcDiv);
                continue;
              }
              if (parsed.error) {
                assistantMsg.textContent += '\n[Error: ' + parsed.error + ']';
                continue;
              }
            } catch(e) {}
            // Regular text token
            fullResponse += data;
            assistantMsg.firstChild
              ? assistantMsg.firstChild.textContent = fullResponse
              : assistantMsg.textContent = fullResponse;
          }
        }
        messages.scrollTop = messages.scrollHeight;
        return reader.read().then(processChunk);
      }

      return reader.read().then(processChunk);
    }).catch(function(err) {
      assistantMsg.classList.remove('streaming-cursor');
      assistantMsg.textContent = 'Error: Could not connect to AI service. Is it running?';
    });
  };

  // --- AI Briefing Panel ---
  var _briefingEventSource = null;

  yarr.openBriefing = function() {
    document.getElementById('briefing-panel').style.display = 'flex';
    yarr.loadBriefing();
  };

  yarr.closeBriefing = function() {
    document.getElementById('briefing-panel').style.display = 'none';
    if (_briefingEventSource) { _briefingEventSource.close(); _briefingEventSource = null; }
  };

  yarr.loadBriefing = function() {
    var since = document.getElementById('briefing-since').value;
    var content = document.getElementById('briefing-content');
    content.innerHTML = '<div class="streaming-cursor" id="briefing-text"></div>';

    var fullText = '';
    fetch('./api/ai/briefing?since=' + since, {
      headers: { 'Accept': 'text/event-stream' }
    }).then(function(response) {
      var reader = response.body.getReader();
      var decoder = new TextDecoder();
      var buffer = '';
      var textEl = document.getElementById('briefing-text');

      function processChunk(result) {
        if (result.done) {
          textEl.classList.remove('streaming-cursor');
          // Simple markdown rendering: paragraphs
          textEl.innerHTML = fullText
            .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
            .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
            .replace(/\[([\d,\s]+)\]/g, '<sup>[$1]</sup>')
            .replace(/\n\n/g, '</p><p>')
            .replace(/^/, '<p>').replace(/$/, '</p>')
            .replace(/###\s*(.*?)<\/p>/g, '<h3>$1</h3>')
            .replace(/##\s*(.*?)<\/p>/g, '<h2>$1</h2>');
          return;
        }
        buffer += decoder.decode(result.value, { stream: true });
        var lines = buffer.split('\n');
        buffer = lines.pop();

        for (var i = 0; i < lines.length; i++) {
          var line = lines[i];
          if (line.startsWith('data: ')) {
            var data = line.slice(6);
            if (data === '[DONE]') continue;
            fullText += data;
            textEl.textContent = fullText;
          }
        }
        content.scrollTop = content.scrollHeight;
        return reader.read().then(processChunk);
      }

      return reader.read().then(processChunk);
    }).catch(function(err) {
      content.innerHTML = '<div class="empty-state"><p>Could not connect to AI service. Is it running?</p></div>';
    });
  };

  // --- Close dropdowns on outside click ---
  document.addEventListener('click', function(e) {
    if (!e.target.closest('.dropdown')) {
      document.querySelectorAll('.dropdown-menu.show').forEach(function(m) { m.classList.remove('show'); });
    }
  });

  // --- Close dialogs on backdrop click ---
  document.querySelectorAll('dialog').forEach(function(d) {
    d.addEventListener('click', function(e) {
      if (e.target === d) d.close();
    });
  });

  // --- Delegated click handler for topic article rows ---
  document.addEventListener('click', function(e) {
    var row = e.target.closest('.topic-article-row');
    if (!row) return;
    var itemId = parseInt(row.getAttribute('data-item-id'), 10);
    var url = row.getAttribute('data-item-url');
    if (itemId && itemId > 0) {
      htmx.ajax('GET', './partials/items/' + itemId, { target: '#col-item-content', swap: 'innerHTML' });
      document.getElementById('app').classList.add('item-selected');
    } else if (url) {
      window.open(url, '_blank', 'noopener,noreferrer');
    }
  });

  // --- Init ---
  document.addEventListener('DOMContentLoaded', function() {
    initDrag('drag-feed', 200, 700, 'feed_list_width');
    initDrag('drag-item', 200, 700, 'item_list_width');
    initSearch();
    initReadingProgress();

    // Check if feeds are refreshing on load
    if (document.getElementById('feed-status').style.display !== 'none') {
      yarr.pollStatus();
    }

    // Check if an AI task is running (survives page reload)
    yarr.checkTaskStatus();
  });

})(window.yarr);
