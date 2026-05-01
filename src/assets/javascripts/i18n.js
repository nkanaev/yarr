(function (exports) {
  const translations = {
    "unread": {
      "en": "Unread",
      "zh": "未读",
      "ru": "Непрочитанные"
    },
    "starred": {
      "en": "Starred",
      "zh": "星标",
      "ru": "Избранные"
    },
    "all": {
      "en": "All",
      "zh": "全部",
      "ru": "Все"
    },
    "settings": {
      "en": "Settings",
      "zh": "设置",
      "ru": "Настройки"
    },
    "new_feed": {
      "en": "New Feed",
      "zh": "新建订阅",
      "ru": "Новая лента"
    },
    "refresh_feeds": {
      "en": "Refresh Feeds",
      "zh": "刷新订阅",
      "ru": "Обновить ленты"
    },
    "theme": {
      "en": "Theme",
      "zh": "主题",
      "ru": "Тема"
    },
    "auto_refresh": {
      "en": "Auto Refresh",
      "zh": "自动刷新",
      "ru": "Автообновление"
    },
    "show_first": {
      "en": "Show first",
      "zh": "优先显示",
      "ru": "Сначала"
    },
    "new": {
      "en": "New",
      "zh": "最新",
      "ru": "Новые"
    },
    "old": {
      "en": "Old",
      "zh": "最旧",
      "ru": "Старые"
    },
    "subscriptions": {
      "en": "Subscriptions",
      "zh": "订阅管理",
      "ru": "Подписки"
    },
    "import": {
      "en": "Import",
      "zh": "导入",
      "ru": "Импорт"
    },
    "export": {
      "en": "Export",
      "zh": "导出",
      "ru": "Экспорт"
    },
    "shortcuts": {
      "en": "Shortcuts",
      "zh": "快捷键",
      "ru": "Горячие клавиши"
    },
    "log_out": {
      "en": "Log out",
      "zh": "登出",
      "ru": "Выйти"
    },
    "all_unread": {
      "en": "All Unread",
      "zh": "全部未读",
      "ru": "Все непрочитанные"
    },
    "all_starred": {
      "en": "All Starred",
      "zh": "全部星标",
      "ru": "Все избранные"
    },
    "all_feeds": {
      "en": "All Feeds",
      "zh": "全部订阅",
      "ru": "Все ленты"
    },
    "refreshing": {
      "en": "Refreshing",
      "zh": "正在刷新",
      "ru": "Обновление"
    },
    "left": {
      "en": "left",
      "zh": "剩余",
      "ru": "осталось"
    },
    "show_feeds": {
      "en": "Show Feeds",
      "zh": "显示订阅",
      "ru": "Показать ленты"
    },
    "mark_all_read": {
      "en": "Mark All Read",
      "zh": "全部标记为已读",
      "ru": "Отметить все как прочитанные"
    },
    "feed_settings": {
      "en": "Feed Settings",
      "zh": "订阅设置",
      "ru": "Настройки ленты"
    },
    "folder_settings": {
      "en": "Folder Settings",
      "zh": "文件夹设置",
      "ru": "Настройки папки"
    },
    "website": {
      "en": "Website",
      "zh": "网站",
      "ru": "Сайт"
    },
    "feed_link": {
      "en": "Feed Link",
      "zh": "订阅链接",
      "ru": "Ссылка на ленту"
    },
    "rename": {
      "en": "Rename",
      "zh": "重命名",
      "ru": "Переименовать"
    },
    "change_link": {
      "en": "Change Link",
      "zh": "修改链接",
      "ru": "Изменить ссылку"
    },
    "move_to": {
      "en": "Move to...",
      "zh": "移动到...",
      "ru": "Переместить в..."
    },
    "new_folder": {
      "en": "new folder",
      "zh": "新建文件夹",
      "ru": "новая папка"
    },
    "delete": {
      "en": "Delete",
      "zh": "删除",
      "ru": "Удалить"
    },
    "mark_starred": {
      "en": "Mark Starred",
      "zh": "标记星标",
      "ru": "Пометить избранным"
    },
    "mark_unread": {
      "en": "Mark Unread",
      "zh": "标记未读",
      "ru": "Пометить непрочитанным"
    },
    "appearance": {
      "en": "Appearance",
      "zh": "外观",
      "ru": "Внешний вид"
    },
    "read_here": {
      "en": "Read Here",
      "zh": "在此阅读",
      "ru": "Читать здесь"
    },
    "open_link": {
      "en": "Open Link",
      "zh": "打开链接",
      "ru": "Открыть ссылку"
    },
    "previous_article": {
      "en": "Previous Article",
      "zh": "上一篇",
      "ru": "Предыдущая статья"
    },
    "next_article": {
      "en": "Next Article",
      "zh": "下一篇",
      "ru": "Следующая статья"
    },
    "close_article": {
      "en": "Close Article",
      "zh": "关闭文章",
      "ru": "Закрыть статью"
    },
    "untitled": {
      "en": "untitled",
      "zh": "无标题",
      "ru": "без названия"
    },
    "sans_serif": {
      "en": "sans-serif",
      "zh": "无衬线",
      "ru": "sans-serif"
    },
    "serif": {
      "en": "serif",
      "zh": "衬线",
      "ru": "serif"
    },
    "monospace": {
      "en": "monospace",
      "zh": "等宽",
      "ru": "monospace"
    },
    "url": {
      "en": "URL",
      "zh": "网址",
      "ru": "URL"
    },
    "folder": {
      "en": "Folder",
      "zh": "文件夹",
      "ru": "Папка"
    },
    "add": {
      "en": "Add",
      "zh": "添加",
      "ru": "Добавить"
    },
    "keyboard_shortcuts": {
      "en": "Keyboard Shortcuts",
      "zh": "键盘快捷键",
      "ru": "Горячие клавиши"
    },
    "multiple_feeds_found": {
      "en": "Multiple feeds found. Choose one below:",
      "zh": "找到多个订阅源，请选择一个：",
      "ru": "Найдено несколько лент. Выберите одну:"
    },
    "cancel": {
      "en": "cancel",
      "zh": "取消",
      "ru": "отмена"
    },
    "kb_show_filters": {
      "en": "show unread / starred / all feeds",
      "zh": "显示未读/星标/全部订阅",
      "ru": "показать непрочитанные / избранные / все ленты"
    },
    "kb_focus_search": {
      "en": "focus the search bar",
      "zh": "聚焦搜索栏",
      "ru": "фокус на строку поиска"
    },
    "kb_next_prev_article": {
      "en": "next / prev article",
      "zh": "下一篇/上一篇文章",
      "ru": "следующая / предыдущая статья"
    },
    "kb_next_prev_feed": {
      "en": "next / prev feed",
      "zh": "下一个/上一个订阅",
      "ru": "следующая / предыдущая лента"
    },
    "kb_close_article": {
      "en": "close article",
      "zh": "关闭文章",
      "ru": "закрыть статью"
    },
    "kb_mark_all_read": {
      "en": "mark all read",
      "zh": "全部标记为已读",
      "ru": "отметить все как прочитанные"
    },
    "kb_mark_read": {
      "en": "mark read / unread",
      "zh": "标记已读/未读",
      "ru": "отметить как прочитанное / непрочитанное"
    },
    "kb_mark_starred": {
      "en": "mark starred / unstarred",
      "zh": "标记星标/取消星标",
      "ru": "пометить избранным / убрать из избранного"
    },
    "kb_open_link": {
      "en": "open link",
      "zh": "打开链接",
      "ru": "открыть ссылку"
    },
    "kb_read_here": {
      "en": "read here",
      "zh": "在此阅读",
      "ru": "читать здесь"
    },
    "kb_scroll_content": {
      "en": "scroll content forward / backward",
      "zh": "向前/向后滚动内容",
      "ru": "прокрутка вперед / назад"
    },
    "prompt_folder_name": {
      "en": "Enter folder name:",
      "zh": "请输入文件夹名称：",
      "ru": "Введите имя папки:"
    },
    "prompt_new_title": {
      "en": "Enter new title",
      "zh": "请输入新标题",
      "ru": "Введите новый заголовок"
    },
    "prompt_feed_link": {
      "en": "Enter feed link",
      "zh": "请输入订阅链接",
      "ru": "Введите ссылку на ленту"
    },
    "confirm_delete_folder": {
      "en": "Are you sure you want to delete",
      "zh": "确定要删除",
      "ru": "Вы уверены, что хотите удалить"
    },
    "confirm_delete_feed": {
      "en": "Are you sure you want to delete",
      "zh": "确定要删除",
      "ru": "Вы уверены, что хотите удалить"
    },
    "alert_no_feeds": {
      "en": "No feeds found at the given url.",
      "zh": "在指定的网址未找到订阅源。",
      "ru": "Лент по данному адресу не найдено."
    },
    "login": {
      "en": "Login",
      "zh": "登录",
      "ru": "Вход"
    },
    "username": {
      "en": "Username",
      "zh": "用户名",
      "ru": "Имя пользователя"
    },
    "password": {
      "en": "Password",
      "zh": "密码",
      "ru": "Пароль"
    },
  };
  class i18n {
    constructor() {
      this.lang = 'en'
    }
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
