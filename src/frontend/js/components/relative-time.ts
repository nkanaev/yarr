import { dateRepr } from '../utils'
import { defineComponent } from 'vue'

export default defineComponent({
  props: ['val'],
  data: function() {
    var d = new Date(this.val)
    return {
      'date': d,
      'formatted': dateRepr(d),
      'interval': undefined as number | undefined,
    }
  },
  template: '<time :datetime="val">{{ formatted }}</time>',
  mounted: function() {
    this.interval = setInterval(() => {
      this.formatted = dateRepr(this.date)
    }, 600000)  // every 10 minutes
  },
  unmounted: function() {
    clearInterval(this.interval)
  },
})
