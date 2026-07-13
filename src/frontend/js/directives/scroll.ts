import { debounce } from '../utils'

export default {
  mounted: function(el, binding) {
    el.addEventListener('scroll', debounce(function(event) {
      binding.value(event, el)
    }, 200))
  },
}
