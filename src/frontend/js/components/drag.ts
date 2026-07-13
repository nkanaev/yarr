import { defineComponent } from 'vue'

export default defineComponent({
  props: ['width'],
  template: '<div class="drag"></div>',
  mounted: function() {
    var self = this
    let startX = 0
    let initW = 0
    var onMouseMove = function(e: MouseEvent) {
      var offset = e.clientX - startX
      var newWidth = initW + offset
      self.$emit('resize', newWidth)
    }
    var onMouseUp = function(e: MouseEvent) {
      document.removeEventListener('mousemove', onMouseMove)
      document.removeEventListener('mouseup', onMouseUp)
    }
    this.$el.addEventListener('mousedown', function(e: MouseEvent) {
      startX = e.clientX
      initW = self.width
      document.addEventListener('mousemove', onMouseMove)
      document.addEventListener('mouseup', onMouseUp)
    })
  },
})
