import { createApp, h } from 'vue'
import i18n from './i18n'
import App from './app'
import Login from './login'

const application = createApp({
  render: function () {
    return h(window.app.authenticated ? App : Login)
  }
})
application.use(i18n)
application.mount('#app')
