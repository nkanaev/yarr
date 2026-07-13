import type { Directive } from 'vue'
import { debounce } from '../utils'

export default {
  mounted: function(el, binding) {
    el.addEventListener('scroll', debounce(function(event) {
      binding.value(event, el)
    }, 200))
  },
  } satisfies Directive<HTMLElement, (event: Event, el: HTMLElement) => void>
