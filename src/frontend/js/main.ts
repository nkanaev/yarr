import Vue from 'vue/dist/vue.esm.js'
import i18n from './i18n'
import App from './app'
import Login from './login'

Vue.use(i18n)

var vm = new Vue({
  render: function (h) {
    return h(window.app.authenticated ? App : Login)
  }
}).$mount('#app')
