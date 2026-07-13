import template from './templates/login.html' with {type: 'text'}
import icons from './icons'
import { defineComponent } from 'vue'

export default defineComponent({
  template: template,
  data: function () {
    return {
      logo: icons.anchor,
      hasError: false,
    }
  },
  created: function () {
    this.$setLang(window.app.settings.language)
  },
  methods: {
    login: function (event: Event) {
      event.preventDefault()
      var data = new FormData(event.target)
      fetch('./login', { method: 'POST', body: data }).then(function (res) {
        if (res.ok) {
          // TODO: 
          document.location.assign('./')
        } else {
          this.hasError = true
        }
        }.bind(this))
    }
  }
})
