import icons from '../icons'

export default {
  props: ['name'],
  template: '<span class="icon" v-html="content"></span>',
  computed: {
    content: function () { return icons[this.name] || '' }
  }
}
