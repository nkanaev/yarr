import { dateRepr } from '../utils'

export default {
  props: ['val'],
  data: function() {
    var d = new Date(this.val)
    return {
      'date': d,
      'formatted': dateRepr(d),
      'interval': null,
    }
  },
  template: '<time :datetime="val">{{ formatted }}</time>',
  mounted: function() {
    this.interval = setInterval(function() {
      this.formatted = dateRepr(this.date)
    }.bind(this), 600000)  // every 10 minutes
  },
  unmounted: function() {
    clearInterval(this.interval)
  },
}
