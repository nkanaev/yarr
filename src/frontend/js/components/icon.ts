import icons from '../icons'
import { defineComponent } from 'vue'

export default defineComponent({
  props: ['name'],
  template: '<span class="icon" v-html="content"></span>',
  computed: {
    content: function () { return icons[this.name] || '' }
  }
})
