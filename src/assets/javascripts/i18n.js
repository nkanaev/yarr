(function (exports) {
  const translations = {
    "unread": {
      "en": "Unread",
      "zh": "未读"
    },
    "starred": {
      "en": "Starred",
      "zh": "星标"
    },
    "all": {
      "en": "All",
      "zh": "全部"
    },
    "settings": {
      "en": "Settings",
      "zh": "设置"
    },
    "new_feed": {
      "en": "New Feed",
      "zh": "新建订阅"
    },
    "refresh_feeds": {
      "en": "Refresh Feeds",
      "zh": "刷新订阅"
    },
    "theme": {
      "en": "Theme",
      "zh": "主题"
    },
    "auto_refresh": {
      "en": "Auto Refresh",
      "zh": "自动刷新"
    },
    "show_first": {
      "en": "Show first",
      "zh": "优先显示"
    },
    "new": {
      "en": "New",
      "zh": "最新"
    },
    "old": {
      "en": "Old",
      "zh": "最旧"
    },
    "subscriptions": {
      "en": "Subscriptions",
      "zh": "订阅管理"
    },
    "import": {
      "en": "Import",
      "zh": "导入"
    },
    "export": {
      "en": "Export",
      "zh": "导出"
    },
    "shortcuts": {
      "en": "Shortcuts",
      "zh": "快捷键"
    },
    "log_out": {
      "en": "Log out",
      "zh": "登出"
    },
    "all_unread": {
      "en": "All Unread",
      "zh": "全部未读"
    },
    "all_starred": {
      "en": "All Starred",
      "zh": "全部星标"
    },
    "all_feeds": {
      "en": "All Feeds",
      "zh": "全部订阅"
    },
    "refreshing": {
      "en": "Refreshing",
      "zh": "正在刷新"
    },
    "left": {
      "en": "left",
      "zh": "剩余"
    },
    "show_feeds": {
      "en": "Show Feeds",
      "zh": "显示订阅"
    },
    "mark_all_read": {
      "en": "Mark All Read",
      "zh": "全部标记为已读"
    },
    "feed_settings": {
      "en": "Feed Settings",
      "zh": "订阅设置"
    },
    "folder_settings": {
      "en": "Folder Settings",
      "zh": "文件夹设置"
    },
    "website": {
      "en": "Website",
      "zh": "网站"
    },
    "feed_link": {
      "en": "Feed Link",
      "zh": "订阅链接"
    },
    "rename": {
      "en": "Rename",
      "zh": "重命名"
    },
    "change_link": {
      "en": "Change Link",
      "zh": "修改链接"
    },
    "move_to": {
      "en": "Move to...",
      "zh": "移动到..."
    },
    "new_folder": {
      "en": "new folder",
      "zh": "新建文件夹"
    },
    "delete": {
      "en": "Delete",
      "zh": "删除"
    },
    "mark_starred": {
      "en": "Mark Starred",
      "zh": "标记星标"
    },
    "mark_unread": {
      "en": "Mark Unread",
      "zh": "标记未读"
    },
    "appearance": {
      "en": "Appearance",
      "zh": "外观"
    },
    "read_here": {
      "en": "Read Here",
      "zh": "在此阅读"
    },
    "open_link": {
      "en": "Open Link",
      "zh": "打开链接"
    },
    "previous_article": {
      "en": "Previous Article",
      "zh": "上一篇"
    },
    "next_article": {
      "en": "Next Article",
      "zh": "下一篇"
    },
    "close_article": {
      "en": "Close Article",
      "zh": "关闭文章"
    },
    "untitled": {
      "en": "untitled",
      "zh": "无标题"
    },
    "sans_serif": {
      "en": "sans-serif",
      "zh": "无衬线"
    },
    "serif": {
      "en": "serif",
      "zh": "衬线"
    },
    "monospace": {
      "en": "monospace",
      "zh": "等宽"
    },
    "url": {
      "en": "URL",
      "zh": "网址"
    },
    "folder": {
      "en": "Folder",
      "zh": "文件夹"
    },
    "add": {
      "en": "Add",
      "zh": "添加"
    },
    "keyboard_shortcuts": {
      "en": "Keyboard Shortcuts",
      "zh": "键盘快捷键"
    },
    "multiple_feeds_found": {
      "en": "Multiple feeds found. Choose one below:",
      "zh": "找到多个订阅源，请选择一个："
    },
    "cancel": {
      "en": "cancel",
      "zh": "取消"
    },
    "kb_show_filters": {
      "en": "show unread / starred / all feeds",
      "zh": "显示未读/星标/全部订阅"
    },
    "kb_focus_search": {
      "en": "focus the search bar",
      "zh": "聚焦搜索栏"
    },
    "kb_next_prev_article": {
      "en": "next / prev article",
      "zh": "下一篇/上一篇文章"
    },
    "kb_next_prev_feed": {
      "en": "next / prev feed",
      "zh": "下一个/上一个订阅"
    },
    "kb_close_article": {
      "en": "close article",
      "zh": "关闭文章"
    },
    "kb_mark_all_read": {
      "en": "mark all read",
      "zh": "全部标记为已读"
    },
    "kb_mark_read": {
      "en": "mark read / unread",
      "zh": "标记已读/未读"
    },
    "kb_mark_starred": {
      "en": "mark starred / unstarred",
      "zh": "标记星标/取消星标"
    },
    "kb_open_link": {
      "en": "open link",
      "zh": "打开链接"
    },
    "kb_read_here": {
      "en": "read here",
      "zh": "在此阅读"
    },
    "kb_scroll_content": {
      "en": "scroll content forward / backward",
      "zh": "向前/向后滚动内容"
    },
    "prompt_folder_name": {
      "en": "Enter folder name:",
      "zh": "请输入文件夹名称："
    },
    "prompt_new_title": {
      "en": "Enter new title",
      "zh": "请输入新标题"
    },
    "prompt_feed_link": {
      "en": "Enter feed link",
      "zh": "请输入订阅链接"
    },
    "confirm_delete_folder": {
      "en": "Are you sure you want to delete",
      "zh": "确定要删除"
    },
    "confirm_delete_feed": {
      "en": "Are you sure you want to delete",
      "zh": "确定要删除"
    },
    "alert_no_feeds": {
      "en": "No feeds found at the given url.",
      "zh": "在指定的网址未找到订阅源。"
    },
    "login": {
      "en": "Login",
      "zh": "登录"
    },
    "username": {
      "en": "Username",
      "zh": "用户名"
    },
    "password": {
      "en": "Password",
      "zh": "密码"
    },
    "language": {
      "en": "Language",
      "zh": "语言"
    }
  };
  class i18n {
    setLang(lang) {
      this.lang = lang
    }
    $t(code) {
      return translations[code][this.lang]
    }
  }
  exports.i18n = {
    install(Vue, opts) {
      const x = new i18n();
      Vue.prototype.$t = x.$t
      Vue.prototype.$setLang = x.setLang
    }
  }
})(window)
