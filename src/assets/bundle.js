(() => {
  // node_modules/vue/dist/vue.esm.js
  var emptyObject = Object.freeze({});
  var isArray = Array.isArray;
  function isUndef(v) {
    return v === void 0 || v === null;
  }
  function isDef(v) {
    return v !== void 0 && v !== null;
  }
  function isTrue(v) {
    return v === true;
  }
  function isFalse(v) {
    return v === false;
  }
  function isPrimitive(value) {
    return typeof value === "string" || typeof value === "number" || // $flow-disable-line
    typeof value === "symbol" || typeof value === "boolean";
  }
  function isFunction(value) {
    return typeof value === "function";
  }
  function isObject(obj) {
    return obj !== null && typeof obj === "object";
  }
  var _toString = Object.prototype.toString;
  function toRawType(value) {
    return _toString.call(value).slice(8, -1);
  }
  function isPlainObject(obj) {
    return _toString.call(obj) === "[object Object]";
  }
  function isRegExp(v) {
    return _toString.call(v) === "[object RegExp]";
  }
  function isValidArrayIndex(val) {
    var n = parseFloat(String(val));
    return n >= 0 && Math.floor(n) === n && isFinite(val);
  }
  function isPromise(val) {
    return isDef(val) && typeof val.then === "function" && typeof val.catch === "function";
  }
  function toString(val) {
    return val == null ? "" : Array.isArray(val) || isPlainObject(val) && val.toString === _toString ? JSON.stringify(val, replacer, 2) : String(val);
  }
  function replacer(_key, val) {
    if (val && val.__v_isRef) {
      return val.value;
    }
    return val;
  }
  function toNumber(val) {
    var n = parseFloat(val);
    return isNaN(n) ? val : n;
  }
  function makeMap(str2, expectsLowerCase) {
    var map = /* @__PURE__ */ Object.create(null);
    var list = str2.split(",");
    for (var i = 0; i < list.length; i++) {
      map[list[i]] = true;
    }
    return expectsLowerCase ? function(val) {
      return map[val.toLowerCase()];
    } : function(val) {
      return map[val];
    };
  }
  var isBuiltInTag = makeMap("slot,component", true);
  var isReservedAttribute = makeMap("key,ref,slot,slot-scope,is");
  function remove$2(arr, item) {
    var len2 = arr.length;
    if (len2) {
      if (item === arr[len2 - 1]) {
        arr.length = len2 - 1;
        return;
      }
      var index2 = arr.indexOf(item);
      if (index2 > -1) {
        return arr.splice(index2, 1);
      }
    }
  }
  var hasOwnProperty = Object.prototype.hasOwnProperty;
  function hasOwn(obj, key) {
    return hasOwnProperty.call(obj, key);
  }
  function cached(fn) {
    var cache2 = /* @__PURE__ */ Object.create(null);
    return function cachedFn(str2) {
      var hit = cache2[str2];
      return hit || (cache2[str2] = fn(str2));
    };
  }
  var camelizeRE = /-(\w)/g;
  var camelize = cached(function(str2) {
    return str2.replace(camelizeRE, function(_, c) {
      return c ? c.toUpperCase() : "";
    });
  });
  var capitalize = cached(function(str2) {
    return str2.charAt(0).toUpperCase() + str2.slice(1);
  });
  var hyphenateRE = /\B([A-Z])/g;
  var hyphenate = cached(function(str2) {
    return str2.replace(hyphenateRE, "-$1").toLowerCase();
  });
  function polyfillBind(fn, ctx) {
    function boundFn(a) {
      var l = arguments.length;
      return l ? l > 1 ? fn.apply(ctx, arguments) : fn.call(ctx, a) : fn.call(ctx);
    }
    boundFn._length = fn.length;
    return boundFn;
  }
  function nativeBind(fn, ctx) {
    return fn.bind(ctx);
  }
  var bind$1 = Function.prototype.bind ? nativeBind : polyfillBind;
  function toArray(list, start) {
    start = start || 0;
    var i = list.length - start;
    var ret = new Array(i);
    while (i--) {
      ret[i] = list[i + start];
    }
    return ret;
  }
  function extend(to, _from) {
    for (var key in _from) {
      to[key] = _from[key];
    }
    return to;
  }
  function toObject(arr) {
    var res = {};
    for (var i = 0; i < arr.length; i++) {
      if (arr[i]) {
        extend(res, arr[i]);
      }
    }
    return res;
  }
  function noop(a, b, c) {
  }
  var no = function(a, b, c) {
    return false;
  };
  var identity = function(_) {
    return _;
  };
  function genStaticKeys$1(modules2) {
    return modules2.reduce(function(keys, m) {
      return keys.concat(m.staticKeys || []);
    }, []).join(",");
  }
  function looseEqual(a, b) {
    if (a === b)
      return true;
    var isObjectA = isObject(a);
    var isObjectB = isObject(b);
    if (isObjectA && isObjectB) {
      try {
        var isArrayA = Array.isArray(a);
        var isArrayB = Array.isArray(b);
        if (isArrayA && isArrayB) {
          return a.length === b.length && a.every(function(e, i) {
            return looseEqual(e, b[i]);
          });
        } else if (a instanceof Date && b instanceof Date) {
          return a.getTime() === b.getTime();
        } else if (!isArrayA && !isArrayB) {
          var keysA = Object.keys(a);
          var keysB = Object.keys(b);
          return keysA.length === keysB.length && keysA.every(function(key) {
            return looseEqual(a[key], b[key]);
          });
        } else {
          return false;
        }
      } catch (e) {
        return false;
      }
    } else if (!isObjectA && !isObjectB) {
      return String(a) === String(b);
    } else {
      return false;
    }
  }
  function looseIndexOf(arr, val) {
    for (var i = 0; i < arr.length; i++) {
      if (looseEqual(arr[i], val))
        return i;
    }
    return -1;
  }
  function once(fn) {
    var called = false;
    return function() {
      if (!called) {
        called = true;
        fn.apply(this, arguments);
      }
    };
  }
  function hasChanged(x, y) {
    if (x === y) {
      return x === 0 && 1 / x !== 1 / y;
    } else {
      return x === x || y === y;
    }
  }
  var SSR_ATTR = "data-server-rendered";
  var ASSET_TYPES = ["component", "directive", "filter"];
  var LIFECYCLE_HOOKS = [
    "beforeCreate",
    "created",
    "beforeMount",
    "mounted",
    "beforeUpdate",
    "updated",
    "beforeDestroy",
    "destroyed",
    "activated",
    "deactivated",
    "errorCaptured",
    "serverPrefetch",
    "renderTracked",
    "renderTriggered"
  ];
  var config = {
    /**
     * Option merge strategies (used in core/util/options)
     */
    // $flow-disable-line
    optionMergeStrategies: /* @__PURE__ */ Object.create(null),
    /**
     * Whether to suppress warnings.
     */
    silent: false,
    /**
     * Show production mode tip message on boot?
     */
    productionTip: true,
    /**
     * Whether to enable devtools
     */
    devtools: true,
    /**
     * Whether to record perf
     */
    performance: false,
    /**
     * Error handler for watcher errors
     */
    errorHandler: null,
    /**
     * Warn handler for watcher warns
     */
    warnHandler: null,
    /**
     * Ignore certain custom elements
     */
    ignoredElements: [],
    /**
     * Custom user key aliases for v-on
     */
    // $flow-disable-line
    keyCodes: /* @__PURE__ */ Object.create(null),
    /**
     * Check if a tag is reserved so that it cannot be registered as a
     * component. This is platform-dependent and may be overwritten.
     */
    isReservedTag: no,
    /**
     * Check if an attribute is reserved so that it cannot be used as a component
     * prop. This is platform-dependent and may be overwritten.
     */
    isReservedAttr: no,
    /**
     * Check if a tag is an unknown element.
     * Platform-dependent.
     */
    isUnknownElement: no,
    /**
     * Get the namespace of an element
     */
    getTagNamespace: noop,
    /**
     * Parse the real tag name for the specific platform.
     */
    parsePlatformTagName: identity,
    /**
     * Check if an attribute must be bound using property, e.g. value
     * Platform-dependent.
     */
    mustUseProp: no,
    /**
     * Perform updates asynchronously. Intended to be used by Vue Test Utils
     * This will significantly reduce performance if set to false.
     */
    async: true,
    /**
     * Exposed for legacy reasons
     */
    _lifecycleHooks: LIFECYCLE_HOOKS
  };
  var unicodeRegExp = /a-zA-Z\u00B7\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u037D\u037F-\u1FFF\u200C-\u200D\u203F-\u2040\u2070-\u218F\u2C00-\u2FEF\u3001-\uD7FF\uF900-\uFDCF\uFDF0-\uFFFD/;
  function isReserved(str2) {
    var c = (str2 + "").charCodeAt(0);
    return c === 36 || c === 95;
  }
  function def(obj, key, val, enumerable) {
    Object.defineProperty(obj, key, {
      value: val,
      enumerable: !!enumerable,
      writable: true,
      configurable: true
    });
  }
  var bailRE = new RegExp("[^".concat(unicodeRegExp.source, ".$_\\d]"));
  function parsePath(path) {
    if (bailRE.test(path)) {
      return;
    }
    var segments = path.split(".");
    return function(obj) {
      for (var i = 0; i < segments.length; i++) {
        if (!obj)
          return;
        obj = obj[segments[i]];
      }
      return obj;
    };
  }
  var hasProto = "__proto__" in {};
  var inBrowser = typeof window !== "undefined";
  var UA = inBrowser && window.navigator.userAgent.toLowerCase();
  var isIE = UA && /msie|trident/.test(UA);
  var isIE9 = UA && UA.indexOf("msie 9.0") > 0;
  var isEdge = UA && UA.indexOf("edge/") > 0;
  UA && UA.indexOf("android") > 0;
  var isIOS = UA && /iphone|ipad|ipod|ios/.test(UA);
  UA && /chrome\/\d+/.test(UA) && !isEdge;
  UA && /phantomjs/.test(UA);
  var isFF = UA && UA.match(/firefox\/(\d+)/);
  var nativeWatch = {}.watch;
  var supportsPassive = false;
  if (inBrowser) {
    try {
      opts = {};
      Object.defineProperty(opts, "passive", {
        get: function() {
          supportsPassive = true;
        }
      });
      window.addEventListener("test-passive", null, opts);
    } catch (e) {
    }
  }
  var opts;
  var _isServer;
  var isServerRendering = function() {
    if (_isServer === void 0) {
      if (!inBrowser && typeof global !== "undefined") {
        _isServer = global["process"] && global["process"].env.VUE_ENV === "server";
      } else {
        _isServer = false;
      }
    }
    return _isServer;
  };
  var devtools = inBrowser && window.__VUE_DEVTOOLS_GLOBAL_HOOK__;
  function isNative(Ctor) {
    return typeof Ctor === "function" && /native code/.test(Ctor.toString());
  }
  var hasSymbol = typeof Symbol !== "undefined" && isNative(Symbol) && typeof Reflect !== "undefined" && isNative(Reflect.ownKeys);
  var _Set;
  if (typeof Set !== "undefined" && isNative(Set)) {
    _Set = Set;
  } else {
    _Set = /** @class */
    (function() {
      function Set2() {
        this.set = /* @__PURE__ */ Object.create(null);
      }
      Set2.prototype.has = function(key) {
        return this.set[key] === true;
      };
      Set2.prototype.add = function(key) {
        this.set[key] = true;
      };
      Set2.prototype.clear = function() {
        this.set = /* @__PURE__ */ Object.create(null);
      };
      return Set2;
    })();
  }
  var currentInstance = null;
  function setCurrentInstance(vm3) {
    if (vm3 === void 0) {
      vm3 = null;
    }
    if (!vm3)
      currentInstance && currentInstance._scope.off();
    currentInstance = vm3;
    vm3 && vm3._scope.on();
  }
  var VNode = (
    /** @class */
    (function() {
      function VNode2(tag, data, children, text2, elm, context, componentOptions, asyncFactory) {
        this.tag = tag;
        this.data = data;
        this.children = children;
        this.text = text2;
        this.elm = elm;
        this.ns = void 0;
        this.context = context;
        this.fnContext = void 0;
        this.fnOptions = void 0;
        this.fnScopeId = void 0;
        this.key = data && data.key;
        this.componentOptions = componentOptions;
        this.componentInstance = void 0;
        this.parent = void 0;
        this.raw = false;
        this.isStatic = false;
        this.isRootInsert = true;
        this.isComment = false;
        this.isCloned = false;
        this.isOnce = false;
        this.asyncFactory = asyncFactory;
        this.asyncMeta = void 0;
        this.isAsyncPlaceholder = false;
      }
      Object.defineProperty(VNode2.prototype, "child", {
        // DEPRECATED: alias for componentInstance for backwards compat.
        /* istanbul ignore next */
        get: function() {
          return this.componentInstance;
        },
        enumerable: false,
        configurable: true
      });
      return VNode2;
    })()
  );
  var createEmptyVNode = function(text2) {
    if (text2 === void 0) {
      text2 = "";
    }
    var node = new VNode();
    node.text = text2;
    node.isComment = true;
    return node;
  };
  function createTextVNode(val) {
    return new VNode(void 0, void 0, void 0, String(val));
  }
  function cloneVNode(vnode) {
    var cloned = new VNode(
      vnode.tag,
      vnode.data,
      // #7975
      // clone children array to avoid mutating original in case of cloning
      // a child.
      vnode.children && vnode.children.slice(),
      vnode.text,
      vnode.elm,
      vnode.context,
      vnode.componentOptions,
      vnode.asyncFactory
    );
    cloned.ns = vnode.ns;
    cloned.isStatic = vnode.isStatic;
    cloned.key = vnode.key;
    cloned.isComment = vnode.isComment;
    cloned.fnContext = vnode.fnContext;
    cloned.fnOptions = vnode.fnOptions;
    cloned.fnScopeId = vnode.fnScopeId;
    cloned.asyncMeta = vnode.asyncMeta;
    cloned.isCloned = true;
    return cloned;
  }
  var initProxy;
  if (true) {
    allowedGlobals_1 = makeMap(
      "Infinity,undefined,NaN,isFinite,isNaN,parseFloat,parseInt,decodeURI,decodeURIComponent,encodeURI,encodeURIComponent,Math,Number,Date,Array,Object,Boolean,String,RegExp,Map,Set,JSON,Intl,BigInt,require"
      // for Webpack/Browserify
    );
    warnNonPresent_1 = function(target2, key) {
      warn$2('Property or method "'.concat(key, '" is not defined on the instance but ') + "referenced during render. Make sure that this property is reactive, either in the data option, or for class-based components, by initializing the property. See: https://v2.vuejs.org/v2/guide/reactivity.html#Declaring-Reactive-Properties.", target2);
    };
    warnReservedPrefix_1 = function(target2, key) {
      warn$2('Property "'.concat(key, '" must be accessed with "$data.').concat(key, '" because ') + 'properties starting with "$" or "_" are not proxied in the Vue instance to prevent conflicts with Vue internals. See: https://v2.vuejs.org/v2/api/#data', target2);
    };
    hasProxy_1 = typeof Proxy !== "undefined" && isNative(Proxy);
    if (hasProxy_1) {
      isBuiltInModifier_1 = makeMap("stop,prevent,self,ctrl,shift,alt,meta,exact");
      config.keyCodes = new Proxy(config.keyCodes, {
        set: function(target2, key, value) {
          if (isBuiltInModifier_1(key)) {
            warn$2("Avoid overwriting built-in modifier in config.keyCodes: .".concat(key));
            return false;
          } else {
            target2[key] = value;
            return true;
          }
        }
      });
    }
    hasHandler_1 = {
      has: function(target2, key) {
        var has2 = key in target2;
        var isAllowed = allowedGlobals_1(key) || typeof key === "string" && key.charAt(0) === "_" && !(key in target2.$data);
        if (!has2 && !isAllowed) {
          if (key in target2.$data)
            warnReservedPrefix_1(target2, key);
          else
            warnNonPresent_1(target2, key);
        }
        return has2 || !isAllowed;
      }
    };
    getHandler_1 = {
      get: function(target2, key) {
        if (typeof key === "string" && !(key in target2)) {
          if (key in target2.$data)
            warnReservedPrefix_1(target2, key);
          else
            warnNonPresent_1(target2, key);
        }
        return target2[key];
      }
    };
    initProxy = function initProxy2(vm3) {
      if (hasProxy_1) {
        var options = vm3.$options;
        var handlers = options.render && options.render._withStripped ? getHandler_1 : hasHandler_1;
        vm3._renderProxy = new Proxy(vm3, handlers);
      } else {
        vm3._renderProxy = vm3;
      }
    };
  }
  var allowedGlobals_1;
  var warnNonPresent_1;
  var warnReservedPrefix_1;
  var hasProxy_1;
  var isBuiltInModifier_1;
  var hasHandler_1;
  var getHandler_1;
  var __assign = function() {
    __assign = Object.assign || function __assign2(t) {
      for (var s, i = 1, n = arguments.length; i < n; i++) {
        s = arguments[i];
        for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
      }
      return t;
    };
    return __assign.apply(this, arguments);
  };
  var uid$2 = 0;
  var pendingCleanupDeps = [];
  var cleanupDeps = function() {
    for (var i = 0; i < pendingCleanupDeps.length; i++) {
      var dep = pendingCleanupDeps[i];
      dep.subs = dep.subs.filter(function(s) {
        return s;
      });
      dep._pending = false;
    }
    pendingCleanupDeps.length = 0;
  };
  var Dep = (
    /** @class */
    (function() {
      function Dep2() {
        this._pending = false;
        this.id = uid$2++;
        this.subs = [];
      }
      Dep2.prototype.addSub = function(sub) {
        this.subs.push(sub);
      };
      Dep2.prototype.removeSub = function(sub) {
        this.subs[this.subs.indexOf(sub)] = null;
        if (!this._pending) {
          this._pending = true;
          pendingCleanupDeps.push(this);
        }
      };
      Dep2.prototype.depend = function(info) {
        if (Dep2.target) {
          Dep2.target.addDep(this);
          if (info && Dep2.target.onTrack) {
            Dep2.target.onTrack(__assign({ effect: Dep2.target }, info));
          }
        }
      };
      Dep2.prototype.notify = function(info) {
        var subs = this.subs.filter(function(s) {
          return s;
        });
        if (!config.async) {
          subs.sort(function(a, b) {
            return a.id - b.id;
          });
        }
        for (var i = 0, l = subs.length; i < l; i++) {
          var sub = subs[i];
          if (info) {
            sub.onTrigger && sub.onTrigger(__assign({ effect: subs[i] }, info));
          }
          sub.update();
        }
      };
      return Dep2;
    })()
  );
  Dep.target = null;
  var targetStack = [];
  function pushTarget(target2) {
    targetStack.push(target2);
    Dep.target = target2;
  }
  function popTarget() {
    targetStack.pop();
    Dep.target = targetStack[targetStack.length - 1];
  }
  var arrayProto = Array.prototype;
  var arrayMethods = Object.create(arrayProto);
  var methodsToPatch = [
    "push",
    "pop",
    "shift",
    "unshift",
    "splice",
    "sort",
    "reverse"
  ];
  methodsToPatch.forEach(function(method) {
    var original = arrayProto[method];
    def(arrayMethods, method, function mutator() {
      var args = [];
      for (var _i = 0; _i < arguments.length; _i++) {
        args[_i] = arguments[_i];
      }
      var result = original.apply(this, args);
      var ob = this.__ob__;
      var inserted;
      switch (method) {
        case "push":
        case "unshift":
          inserted = args;
          break;
        case "splice":
          inserted = args.slice(2);
          break;
      }
      if (inserted)
        ob.observeArray(inserted);
      if (true) {
        ob.dep.notify({
          type: "array mutation",
          target: this,
          key: method
        });
      } else {
        ob.dep.notify();
      }
      return result;
    });
  });
  var arrayKeys = Object.getOwnPropertyNames(arrayMethods);
  var NO_INITIAL_VALUE = {};
  var shouldObserve = true;
  function toggleObserving(value) {
    shouldObserve = value;
  }
  var mockDep = {
    notify: noop,
    depend: noop,
    addSub: noop,
    removeSub: noop
  };
  var Observer = (
    /** @class */
    (function() {
      function Observer2(value, shallow, mock) {
        if (shallow === void 0) {
          shallow = false;
        }
        if (mock === void 0) {
          mock = false;
        }
        this.value = value;
        this.shallow = shallow;
        this.mock = mock;
        this.dep = mock ? mockDep : new Dep();
        this.vmCount = 0;
        def(value, "__ob__", this);
        if (isArray(value)) {
          if (!mock) {
            if (hasProto) {
              value.__proto__ = arrayMethods;
            } else {
              for (var i = 0, l = arrayKeys.length; i < l; i++) {
                var key = arrayKeys[i];
                def(value, key, arrayMethods[key]);
              }
            }
          }
          if (!shallow) {
            this.observeArray(value);
          }
        } else {
          var keys = Object.keys(value);
          for (var i = 0; i < keys.length; i++) {
            var key = keys[i];
            defineReactive(value, key, NO_INITIAL_VALUE, void 0, shallow, mock);
          }
        }
      }
      Observer2.prototype.observeArray = function(value) {
        for (var i = 0, l = value.length; i < l; i++) {
          observe(value[i], false, this.mock);
        }
      };
      return Observer2;
    })()
  );
  function observe(value, shallow, ssrMockReactivity) {
    if (value && hasOwn(value, "__ob__") && value.__ob__ instanceof Observer) {
      return value.__ob__;
    }
    if (shouldObserve && (ssrMockReactivity || !isServerRendering()) && (isArray(value) || isPlainObject(value)) && Object.isExtensible(value) && !value.__v_skip && !isRef(value) && !(value instanceof VNode)) {
      return new Observer(value, shallow, ssrMockReactivity);
    }
  }
  function defineReactive(obj, key, val, customSetter, shallow, mock, observeEvenIfShallow) {
    if (observeEvenIfShallow === void 0) {
      observeEvenIfShallow = false;
    }
    var dep = new Dep();
    var property = Object.getOwnPropertyDescriptor(obj, key);
    if (property && property.configurable === false) {
      return;
    }
    var getter = property && property.get;
    var setter = property && property.set;
    if ((!getter || setter) && (val === NO_INITIAL_VALUE || arguments.length === 2)) {
      val = obj[key];
    }
    var childOb = shallow ? val && val.__ob__ : observe(val, false, mock);
    Object.defineProperty(obj, key, {
      enumerable: true,
      configurable: true,
      get: function reactiveGetter() {
        var value = getter ? getter.call(obj) : val;
        if (Dep.target) {
          if (true) {
            dep.depend({
              target: obj,
              type: "get",
              key
            });
          } else {
            dep.depend();
          }
          if (childOb) {
            childOb.dep.depend();
            if (isArray(value)) {
              dependArray(value);
            }
          }
        }
        return isRef(value) && !shallow ? value.value : value;
      },
      set: function reactiveSetter(newVal) {
        var value = getter ? getter.call(obj) : val;
        if (!hasChanged(value, newVal)) {
          return;
        }
        if (customSetter) {
          customSetter();
        }
        if (setter) {
          setter.call(obj, newVal);
        } else if (getter) {
          return;
        } else if (!shallow && isRef(value) && !isRef(newVal)) {
          value.value = newVal;
          return;
        } else {
          val = newVal;
        }
        childOb = shallow ? newVal && newVal.__ob__ : observe(newVal, false, mock);
        if (true) {
          dep.notify({
            type: "set",
            target: obj,
            key,
            newValue: newVal,
            oldValue: value
          });
        } else {
          dep.notify();
        }
      }
    });
    return dep;
  }
  function set(target2, key, val) {
    if (isUndef(target2) || isPrimitive(target2)) {
      warn$2("Cannot set reactive property on undefined, null, or primitive value: ".concat(target2));
    }
    if (isReadonly(target2)) {
      warn$2('Set operation on key "'.concat(key, '" failed: target is readonly.'));
      return;
    }
    var ob = target2.__ob__;
    if (isArray(target2) && isValidArrayIndex(key)) {
      target2.length = Math.max(target2.length, key);
      target2.splice(key, 1, val);
      if (ob && !ob.shallow && ob.mock) {
        observe(val, false, true);
      }
      return val;
    }
    if (key in target2 && !(key in Object.prototype)) {
      target2[key] = val;
      return val;
    }
    if (target2._isVue || ob && ob.vmCount) {
      warn$2("Avoid adding reactive properties to a Vue instance or its root $data at runtime - declare it upfront in the data option.");
      return val;
    }
    if (!ob) {
      target2[key] = val;
      return val;
    }
    defineReactive(ob.value, key, val, void 0, ob.shallow, ob.mock);
    if (true) {
      ob.dep.notify({
        type: "add",
        target: target2,
        key,
        newValue: val,
        oldValue: void 0
      });
    } else {
      ob.dep.notify();
    }
    return val;
  }
  function del(target2, key) {
    if (isUndef(target2) || isPrimitive(target2)) {
      warn$2("Cannot delete reactive property on undefined, null, or primitive value: ".concat(target2));
    }
    if (isArray(target2) && isValidArrayIndex(key)) {
      target2.splice(key, 1);
      return;
    }
    var ob = target2.__ob__;
    if (target2._isVue || ob && ob.vmCount) {
      warn$2("Avoid deleting properties on a Vue instance or its root $data - just set it to null.");
      return;
    }
    if (isReadonly(target2)) {
      warn$2('Delete operation on key "'.concat(key, '" failed: target is readonly.'));
      return;
    }
    if (!hasOwn(target2, key)) {
      return;
    }
    delete target2[key];
    if (!ob) {
      return;
    }
    if (true) {
      ob.dep.notify({
        type: "delete",
        target: target2,
        key
      });
    } else {
      ob.dep.notify();
    }
  }
  function dependArray(value) {
    for (var e = void 0, i = 0, l = value.length; i < l; i++) {
      e = value[i];
      if (e && e.__ob__) {
        e.__ob__.dep.depend();
      }
      if (isArray(e)) {
        dependArray(e);
      }
    }
  }
  function shallowReactive(target2) {
    makeReactive(target2, true);
    def(target2, "__v_isShallow", true);
    return target2;
  }
  function makeReactive(target2, shallow) {
    if (!isReadonly(target2)) {
      if (true) {
        if (isArray(target2)) {
          warn$2("Avoid using Array as root value for ".concat(shallow ? "shallowReactive()" : "reactive()", " as it cannot be tracked in watch() or watchEffect(). Use ").concat(shallow ? "shallowRef()" : "ref()", " instead. This is a Vue-2-only limitation."));
        }
        var existingOb = target2 && target2.__ob__;
        if (existingOb && existingOb.shallow !== shallow) {
          warn$2("Target is already a ".concat(existingOb.shallow ? "" : "non-", "shallow reactive object, and cannot be converted to ").concat(shallow ? "" : "non-", "shallow."));
        }
      }
      var ob = observe(
        target2,
        shallow,
        isServerRendering()
        /* ssr mock reactivity */
      );
      if (!ob) {
        if (target2 == null || isPrimitive(target2)) {
          warn$2("value cannot be made reactive: ".concat(String(target2)));
        }
        if (isCollectionType(target2)) {
          warn$2("Vue 2 does not support reactive collection types such as Map or Set.");
        }
      }
    }
  }
  function isReadonly(value) {
    return !!(value && value.__v_isReadonly);
  }
  function isCollectionType(value) {
    var type = toRawType(value);
    return type === "Map" || type === "WeakMap" || type === "Set" || type === "WeakSet";
  }
  function isRef(r) {
    return !!(r && r.__v_isRef === true);
  }
  function proxyWithRefUnwrap(target2, source, key) {
    Object.defineProperty(target2, key, {
      enumerable: true,
      configurable: true,
      get: function() {
        var val = source[key];
        if (isRef(val)) {
          return val.value;
        } else {
          var ob = val && val.__ob__;
          if (ob)
            ob.dep.depend();
          return val;
        }
      },
      set: function(value) {
        var oldValue = source[key];
        if (isRef(oldValue) && !isRef(value)) {
          oldValue.value = value;
        } else {
          source[key] = value;
        }
      }
    });
  }
  var mark;
  var measure;
  if (true) {
    perf_1 = inBrowser && window.performance;
    if (perf_1 && // @ts-ignore
    perf_1.mark && // @ts-ignore
    perf_1.measure && // @ts-ignore
    perf_1.clearMarks && // @ts-ignore
    perf_1.clearMeasures) {
      mark = function(tag) {
        return perf_1.mark(tag);
      };
      measure = function(name, startTag, endTag2) {
        perf_1.measure(name, startTag, endTag2);
        perf_1.clearMarks(startTag);
        perf_1.clearMarks(endTag2);
      };
    }
  }
  var perf_1;
  var normalizeEvent = cached(function(name) {
    var passive = name.charAt(0) === "&";
    name = passive ? name.slice(1) : name;
    var once2 = name.charAt(0) === "~";
    name = once2 ? name.slice(1) : name;
    var capture = name.charAt(0) === "!";
    name = capture ? name.slice(1) : name;
    return {
      name,
      once: once2,
      capture,
      passive
    };
  });
  function createFnInvoker(fns, vm3) {
    function invoker() {
      var fns2 = invoker.fns;
      if (isArray(fns2)) {
        var cloned = fns2.slice();
        for (var i = 0; i < cloned.length; i++) {
          invokeWithErrorHandling(cloned[i], null, arguments, vm3, "v-on handler");
        }
      } else {
        return invokeWithErrorHandling(fns2, null, arguments, vm3, "v-on handler");
      }
    }
    invoker.fns = fns;
    return invoker;
  }
  function updateListeners(on2, oldOn, add2, remove2, createOnceHandler2, vm3) {
    var name, cur, old, event;
    for (name in on2) {
      cur = on2[name];
      old = oldOn[name];
      event = normalizeEvent(name);
      if (isUndef(cur)) {
        warn$2('Invalid handler for event "'.concat(event.name, '": got ') + String(cur), vm3);
      } else if (isUndef(old)) {
        if (isUndef(cur.fns)) {
          cur = on2[name] = createFnInvoker(cur, vm3);
        }
        if (isTrue(event.once)) {
          cur = on2[name] = createOnceHandler2(event.name, cur, event.capture);
        }
        add2(event.name, cur, event.capture, event.passive, event.params);
      } else if (cur !== old) {
        old.fns = cur;
        on2[name] = old;
      }
    }
    for (name in oldOn) {
      if (isUndef(on2[name])) {
        event = normalizeEvent(name);
        remove2(event.name, oldOn[name], event.capture);
      }
    }
  }
  function mergeVNodeHook(def2, hookKey, hook) {
    if (def2 instanceof VNode) {
      def2 = def2.data.hook || (def2.data.hook = {});
    }
    var invoker;
    var oldHook = def2[hookKey];
    function wrappedHook() {
      hook.apply(this, arguments);
      remove$2(invoker.fns, wrappedHook);
    }
    if (isUndef(oldHook)) {
      invoker = createFnInvoker([wrappedHook]);
    } else {
      if (isDef(oldHook.fns) && isTrue(oldHook.merged)) {
        invoker = oldHook;
        invoker.fns.push(wrappedHook);
      } else {
        invoker = createFnInvoker([oldHook, wrappedHook]);
      }
    }
    invoker.merged = true;
    def2[hookKey] = invoker;
  }
  function extractPropsFromVNodeData(data, Ctor, tag) {
    var propOptions = Ctor.options.props;
    if (isUndef(propOptions)) {
      return;
    }
    var res = {};
    var attrs2 = data.attrs, props2 = data.props;
    if (isDef(attrs2) || isDef(props2)) {
      for (var key in propOptions) {
        var altKey = hyphenate(key);
        if (true) {
          var keyInLowerCase = key.toLowerCase();
          if (key !== keyInLowerCase && attrs2 && hasOwn(attrs2, keyInLowerCase)) {
            tip('Prop "'.concat(keyInLowerCase, '" is passed to component ') + "".concat(formatComponentName(
              // @ts-expect-error tag is string
              tag || Ctor
            ), ", but the declared prop name is") + ' "'.concat(key, '". ') + "Note that HTML attributes are case-insensitive and camelCased props need to use their kebab-case equivalents when using in-DOM " + 'templates. You should probably use "'.concat(altKey, '" instead of "').concat(key, '".'));
          }
        }
        checkProp(res, props2, key, altKey, true) || checkProp(res, attrs2, key, altKey, false);
      }
    }
    return res;
  }
  function checkProp(res, hash2, key, altKey, preserve) {
    if (isDef(hash2)) {
      if (hasOwn(hash2, key)) {
        res[key] = hash2[key];
        if (!preserve) {
          delete hash2[key];
        }
        return true;
      } else if (hasOwn(hash2, altKey)) {
        res[key] = hash2[altKey];
        if (!preserve) {
          delete hash2[altKey];
        }
        return true;
      }
    }
    return false;
  }
  function simpleNormalizeChildren(children) {
    for (var i = 0; i < children.length; i++) {
      if (isArray(children[i])) {
        return Array.prototype.concat.apply([], children);
      }
    }
    return children;
  }
  function normalizeChildren(children) {
    return isPrimitive(children) ? [createTextVNode(children)] : isArray(children) ? normalizeArrayChildren(children) : void 0;
  }
  function isTextNode(node) {
    return isDef(node) && isDef(node.text) && isFalse(node.isComment);
  }
  function normalizeArrayChildren(children, nestedIndex) {
    var res = [];
    var i, c, lastIndex, last;
    for (i = 0; i < children.length; i++) {
      c = children[i];
      if (isUndef(c) || typeof c === "boolean")
        continue;
      lastIndex = res.length - 1;
      last = res[lastIndex];
      if (isArray(c)) {
        if (c.length > 0) {
          c = normalizeArrayChildren(c, "".concat(nestedIndex || "", "_").concat(i));
          if (isTextNode(c[0]) && isTextNode(last)) {
            res[lastIndex] = createTextVNode(last.text + c[0].text);
            c.shift();
          }
          res.push.apply(res, c);
        }
      } else if (isPrimitive(c)) {
        if (isTextNode(last)) {
          res[lastIndex] = createTextVNode(last.text + c);
        } else if (c !== "") {
          res.push(createTextVNode(c));
        }
      } else {
        if (isTextNode(c) && isTextNode(last)) {
          res[lastIndex] = createTextVNode(last.text + c.text);
        } else {
          if (isTrue(children._isVList) && isDef(c.tag) && isUndef(c.key) && isDef(nestedIndex)) {
            c.key = "__vlist".concat(nestedIndex, "_").concat(i, "__");
          }
          res.push(c);
        }
      }
    }
    return res;
  }
  var SIMPLE_NORMALIZE = 1;
  var ALWAYS_NORMALIZE = 2;
  function createElement$1(context, tag, data, children, normalizationType, alwaysNormalize) {
    if (isArray(data) || isPrimitive(data)) {
      normalizationType = children;
      children = data;
      data = void 0;
    }
    if (isTrue(alwaysNormalize)) {
      normalizationType = ALWAYS_NORMALIZE;
    }
    return _createElement(context, tag, data, children, normalizationType);
  }
  function _createElement(context, tag, data, children, normalizationType) {
    if (isDef(data) && isDef(data.__ob__)) {
      warn$2("Avoid using observed data object as vnode data: ".concat(JSON.stringify(data), "\n") + "Always create fresh vnode data objects in each render!", context);
      return createEmptyVNode();
    }
    if (isDef(data) && isDef(data.is)) {
      tag = data.is;
    }
    if (!tag) {
      return createEmptyVNode();
    }
    if (isDef(data) && isDef(data.key) && !isPrimitive(data.key)) {
      warn$2("Avoid using non-primitive value as key, use string/number value instead.", context);
    }
    if (isArray(children) && isFunction(children[0])) {
      data = data || {};
      data.scopedSlots = { default: children[0] };
      children.length = 0;
    }
    if (normalizationType === ALWAYS_NORMALIZE) {
      children = normalizeChildren(children);
    } else if (normalizationType === SIMPLE_NORMALIZE) {
      children = simpleNormalizeChildren(children);
    }
    var vnode, ns;
    if (typeof tag === "string") {
      var Ctor = void 0;
      ns = context.$vnode && context.$vnode.ns || config.getTagNamespace(tag);
      if (config.isReservedTag(tag)) {
        if (isDef(data) && isDef(data.nativeOn) && data.tag !== "component") {
          warn$2("The .native modifier for v-on is only valid on components but it was used on <".concat(tag, ">."), context);
        }
        vnode = new VNode(config.parsePlatformTagName(tag), data, children, void 0, void 0, context);
      } else if ((!data || !data.pre) && isDef(Ctor = resolveAsset(context.$options, "components", tag))) {
        vnode = createComponent(Ctor, data, context, children, tag);
      } else {
        vnode = new VNode(tag, data, children, void 0, void 0, context);
      }
    } else {
      vnode = createComponent(tag, data, context, children);
    }
    if (isArray(vnode)) {
      return vnode;
    } else if (isDef(vnode)) {
      if (isDef(ns))
        applyNS(vnode, ns);
      if (isDef(data))
        registerDeepBindings(data);
      return vnode;
    } else {
      return createEmptyVNode();
    }
  }
  function applyNS(vnode, ns, force) {
    vnode.ns = ns;
    if (vnode.tag === "foreignObject") {
      ns = void 0;
      force = true;
    }
    if (isDef(vnode.children)) {
      for (var i = 0, l = vnode.children.length; i < l; i++) {
        var child = vnode.children[i];
        if (isDef(child.tag) && (isUndef(child.ns) || isTrue(force) && child.tag !== "svg")) {
          applyNS(child, ns, force);
        }
      }
    }
  }
  function registerDeepBindings(data) {
    if (isObject(data.style)) {
      traverse(data.style);
    }
    if (isObject(data.class)) {
      traverse(data.class);
    }
  }
  function renderList(val, render) {
    var ret = null, i, l, keys, key;
    if (isArray(val) || typeof val === "string") {
      ret = new Array(val.length);
      for (i = 0, l = val.length; i < l; i++) {
        ret[i] = render(val[i], i);
      }
    } else if (typeof val === "number") {
      ret = new Array(val);
      for (i = 0; i < val; i++) {
        ret[i] = render(i + 1, i);
      }
    } else if (isObject(val)) {
      if (hasSymbol && val[Symbol.iterator]) {
        ret = [];
        var iterator = val[Symbol.iterator]();
        var result = iterator.next();
        while (!result.done) {
          ret.push(render(result.value, ret.length));
          result = iterator.next();
        }
      } else {
        keys = Object.keys(val);
        ret = new Array(keys.length);
        for (i = 0, l = keys.length; i < l; i++) {
          key = keys[i];
          ret[i] = render(val[key], key, i);
        }
      }
    }
    if (!isDef(ret)) {
      ret = [];
    }
    ret._isVList = true;
    return ret;
  }
  function renderSlot(name, fallbackRender, props2, bindObject) {
    var scopedSlotFn = this.$scopedSlots[name];
    var nodes;
    if (scopedSlotFn) {
      props2 = props2 || {};
      if (bindObject) {
        if (!isObject(bindObject)) {
          warn$2("slot v-bind without argument expects an Object", this);
        }
        props2 = extend(extend({}, bindObject), props2);
      }
      nodes = scopedSlotFn(props2) || (isFunction(fallbackRender) ? fallbackRender() : fallbackRender);
    } else {
      nodes = this.$slots[name] || (isFunction(fallbackRender) ? fallbackRender() : fallbackRender);
    }
    var target2 = props2 && props2.slot;
    if (target2) {
      return this.$createElement("template", { slot: target2 }, nodes);
    } else {
      return nodes;
    }
  }
  function resolveFilter(id) {
    return resolveAsset(this.$options, "filters", id, true) || identity;
  }
  function isKeyNotMatch(expect, actual) {
    if (isArray(expect)) {
      return expect.indexOf(actual) === -1;
    } else {
      return expect !== actual;
    }
  }
  function checkKeyCodes(eventKeyCode, key, builtInKeyCode, eventKeyName, builtInKeyName) {
    var mappedKeyCode = config.keyCodes[key] || builtInKeyCode;
    if (builtInKeyName && eventKeyName && !config.keyCodes[key]) {
      return isKeyNotMatch(builtInKeyName, eventKeyName);
    } else if (mappedKeyCode) {
      return isKeyNotMatch(mappedKeyCode, eventKeyCode);
    } else if (eventKeyName) {
      return hyphenate(eventKeyName) !== key;
    }
    return eventKeyCode === void 0;
  }
  function bindObjectProps(data, tag, value, asProp, isSync) {
    if (value) {
      if (!isObject(value)) {
        warn$2("v-bind without argument expects an Object or Array value", this);
      } else {
        if (isArray(value)) {
          value = toObject(value);
        }
        var hash2 = void 0;
        var _loop_1 = function(key2) {
          if (key2 === "class" || key2 === "style" || isReservedAttribute(key2)) {
            hash2 = data;
          } else {
            var type = data.attrs && data.attrs.type;
            hash2 = asProp || config.mustUseProp(tag, type, key2) ? data.domProps || (data.domProps = {}) : data.attrs || (data.attrs = {});
          }
          var camelizedKey = camelize(key2);
          var hyphenatedKey = hyphenate(key2);
          if (!(camelizedKey in hash2) && !(hyphenatedKey in hash2)) {
            hash2[key2] = value[key2];
            if (isSync) {
              var on2 = data.on || (data.on = {});
              on2["update:".concat(key2)] = function($event) {
                value[key2] = $event;
              };
            }
          }
        };
        for (var key in value) {
          _loop_1(key);
        }
      }
    }
    return data;
  }
  function renderStatic(index2, isInFor) {
    var cached2 = this._staticTrees || (this._staticTrees = []);
    var tree = cached2[index2];
    if (tree && !isInFor) {
      return tree;
    }
    tree = cached2[index2] = this.$options.staticRenderFns[index2].call(
      this._renderProxy,
      this._c,
      this
      // for render fns generated for functional component templates
    );
    markStatic$1(tree, "__static__".concat(index2), false);
    return tree;
  }
  function markOnce(tree, index2, key) {
    markStatic$1(tree, "__once__".concat(index2).concat(key ? "_".concat(key) : ""), true);
    return tree;
  }
  function markStatic$1(tree, key, isOnce) {
    if (isArray(tree)) {
      for (var i = 0; i < tree.length; i++) {
        if (tree[i] && typeof tree[i] !== "string") {
          markStaticNode(tree[i], "".concat(key, "_").concat(i), isOnce);
        }
      }
    } else {
      markStaticNode(tree, key, isOnce);
    }
  }
  function markStaticNode(node, key, isOnce) {
    node.isStatic = true;
    node.key = key;
    node.isOnce = isOnce;
  }
  function bindObjectListeners(data, value) {
    if (value) {
      if (!isPlainObject(value)) {
        warn$2("v-on without argument expects an Object value", this);
      } else {
        var on2 = data.on = data.on ? extend({}, data.on) : {};
        for (var key in value) {
          var existing = on2[key];
          var ours = value[key];
          on2[key] = existing ? [].concat(existing, ours) : ours;
        }
      }
    }
    return data;
  }
  function resolveScopedSlots(fns, res, hasDynamicKeys, contentHashKey) {
    res = res || { $stable: !hasDynamicKeys };
    for (var i = 0; i < fns.length; i++) {
      var slot = fns[i];
      if (isArray(slot)) {
        resolveScopedSlots(slot, res, hasDynamicKeys);
      } else if (slot) {
        if (slot.proxy) {
          slot.fn.proxy = true;
        }
        res[slot.key] = slot.fn;
      }
    }
    if (contentHashKey) {
      res.$key = contentHashKey;
    }
    return res;
  }
  function bindDynamicKeys(baseObj, values2) {
    for (var i = 0; i < values2.length; i += 2) {
      var key = values2[i];
      if (typeof key === "string" && key) {
        baseObj[values2[i]] = values2[i + 1];
      } else if (key !== "" && key !== null) {
        warn$2("Invalid value for dynamic directive argument (expected string or null): ".concat(key), this);
      }
    }
    return baseObj;
  }
  function prependModifier(value, symbol) {
    return typeof value === "string" ? symbol + value : value;
  }
  function installRenderHelpers(target2) {
    target2._o = markOnce;
    target2._n = toNumber;
    target2._s = toString;
    target2._l = renderList;
    target2._t = renderSlot;
    target2._q = looseEqual;
    target2._i = looseIndexOf;
    target2._m = renderStatic;
    target2._f = resolveFilter;
    target2._k = checkKeyCodes;
    target2._b = bindObjectProps;
    target2._v = createTextVNode;
    target2._e = createEmptyVNode;
    target2._u = resolveScopedSlots;
    target2._g = bindObjectListeners;
    target2._d = bindDynamicKeys;
    target2._p = prependModifier;
  }
  function resolveSlots(children, context) {
    if (!children || !children.length) {
      return {};
    }
    var slots = {};
    for (var i = 0, l = children.length; i < l; i++) {
      var child = children[i];
      var data = child.data;
      if (data && data.attrs && data.attrs.slot) {
        delete data.attrs.slot;
      }
      if ((child.context === context || child.fnContext === context) && data && data.slot != null) {
        var name_1 = data.slot;
        var slot = slots[name_1] || (slots[name_1] = []);
        if (child.tag === "template") {
          slot.push.apply(slot, child.children || []);
        } else {
          slot.push(child);
        }
      } else {
        (slots.default || (slots.default = [])).push(child);
      }
    }
    for (var name_2 in slots) {
      if (slots[name_2].every(isWhitespace)) {
        delete slots[name_2];
      }
    }
    return slots;
  }
  function isWhitespace(node) {
    return node.isComment && !node.asyncFactory || node.text === " ";
  }
  function isAsyncPlaceholder(node) {
    return node.isComment && node.asyncFactory;
  }
  function normalizeScopedSlots(ownerVm, scopedSlots, normalSlots, prevScopedSlots) {
    var res;
    var hasNormalSlots = Object.keys(normalSlots).length > 0;
    var isStable = scopedSlots ? !!scopedSlots.$stable : !hasNormalSlots;
    var key = scopedSlots && scopedSlots.$key;
    if (!scopedSlots) {
      res = {};
    } else if (scopedSlots._normalized) {
      return scopedSlots._normalized;
    } else if (isStable && prevScopedSlots && prevScopedSlots !== emptyObject && key === prevScopedSlots.$key && !hasNormalSlots && !prevScopedSlots.$hasNormal) {
      return prevScopedSlots;
    } else {
      res = {};
      for (var key_1 in scopedSlots) {
        if (scopedSlots[key_1] && key_1[0] !== "$") {
          res[key_1] = normalizeScopedSlot(ownerVm, normalSlots, key_1, scopedSlots[key_1]);
        }
      }
    }
    for (var key_2 in normalSlots) {
      if (!(key_2 in res)) {
        res[key_2] = proxyNormalSlot(normalSlots, key_2);
      }
    }
    if (scopedSlots && Object.isExtensible(scopedSlots)) {
      scopedSlots._normalized = res;
    }
    def(res, "$stable", isStable);
    def(res, "$key", key);
    def(res, "$hasNormal", hasNormalSlots);
    return res;
  }
  function normalizeScopedSlot(vm3, normalSlots, key, fn) {
    var normalized = function() {
      var cur = currentInstance;
      setCurrentInstance(vm3);
      var res = arguments.length ? fn.apply(null, arguments) : fn({});
      res = res && typeof res === "object" && !isArray(res) ? [res] : normalizeChildren(res);
      var vnode = res && res[0];
      setCurrentInstance(cur);
      return res && (!vnode || res.length === 1 && vnode.isComment && !isAsyncPlaceholder(vnode)) ? void 0 : res;
    };
    if (fn.proxy) {
      Object.defineProperty(normalSlots, key, {
        get: normalized,
        enumerable: true,
        configurable: true
      });
    }
    return normalized;
  }
  function proxyNormalSlot(slots, key) {
    return function() {
      return slots[key];
    };
  }
  function initSetup(vm3) {
    var options = vm3.$options;
    var setup = options.setup;
    if (setup) {
      var ctx = vm3._setupContext = createSetupContext(vm3);
      setCurrentInstance(vm3);
      pushTarget();
      var setupResult = invokeWithErrorHandling(setup, null, [vm3._props || shallowReactive({}), ctx], vm3, "setup");
      popTarget();
      setCurrentInstance();
      if (isFunction(setupResult)) {
        options.render = setupResult;
      } else if (isObject(setupResult)) {
        if (setupResult instanceof VNode) {
          warn$2("setup() should not return VNodes directly - return a render function instead.");
        }
        vm3._setupState = setupResult;
        if (!setupResult.__sfc) {
          for (var key in setupResult) {
            if (!isReserved(key)) {
              proxyWithRefUnwrap(vm3, setupResult, key);
            } else if (true) {
              warn$2("Avoid using variables that start with _ or $ in setup().");
            }
          }
        } else {
          var proxy2 = vm3._setupProxy = {};
          for (var key in setupResult) {
            if (key !== "__sfc") {
              proxyWithRefUnwrap(proxy2, setupResult, key);
            }
          }
        }
      } else if (setupResult !== void 0) {
        warn$2("setup() should return an object. Received: ".concat(setupResult === null ? "null" : typeof setupResult));
      }
    }
  }
  function createSetupContext(vm3) {
    var exposeCalled = false;
    return {
      get attrs() {
        if (!vm3._attrsProxy) {
          var proxy2 = vm3._attrsProxy = {};
          def(proxy2, "_v_attr_proxy", true);
          syncSetupProxy(proxy2, vm3.$attrs, emptyObject, vm3, "$attrs");
        }
        return vm3._attrsProxy;
      },
      get listeners() {
        if (!vm3._listenersProxy) {
          var proxy2 = vm3._listenersProxy = {};
          syncSetupProxy(proxy2, vm3.$listeners, emptyObject, vm3, "$listeners");
        }
        return vm3._listenersProxy;
      },
      get slots() {
        return initSlotsProxy(vm3);
      },
      emit: bind$1(vm3.$emit, vm3),
      expose: function(exposed) {
        if (true) {
          if (exposeCalled) {
            warn$2("expose() should be called only once per setup().", vm3);
          }
          exposeCalled = true;
        }
        if (exposed) {
          Object.keys(exposed).forEach(function(key) {
            return proxyWithRefUnwrap(vm3, exposed, key);
          });
        }
      }
    };
  }
  function syncSetupProxy(to, from, prev, instance, type) {
    var changed = false;
    for (var key in from) {
      if (!(key in to)) {
        changed = true;
        defineProxyAttr(to, key, instance, type);
      } else if (from[key] !== prev[key]) {
        changed = true;
      }
    }
    for (var key in to) {
      if (!(key in from)) {
        changed = true;
        delete to[key];
      }
    }
    return changed;
  }
  function defineProxyAttr(proxy2, key, instance, type) {
    Object.defineProperty(proxy2, key, {
      enumerable: true,
      configurable: true,
      get: function() {
        return instance[type][key];
      }
    });
  }
  function initSlotsProxy(vm3) {
    if (!vm3._slotsProxy) {
      syncSetupSlots(vm3._slotsProxy = {}, vm3.$scopedSlots);
    }
    return vm3._slotsProxy;
  }
  function syncSetupSlots(to, from) {
    for (var key in from) {
      to[key] = from[key];
    }
    for (var key in to) {
      if (!(key in from)) {
        delete to[key];
      }
    }
  }
  function initRender(vm3) {
    vm3._vnode = null;
    vm3._staticTrees = null;
    var options = vm3.$options;
    var parentVnode = vm3.$vnode = options._parentVnode;
    var renderContext = parentVnode && parentVnode.context;
    vm3.$slots = resolveSlots(options._renderChildren, renderContext);
    vm3.$scopedSlots = parentVnode ? normalizeScopedSlots(vm3.$parent, parentVnode.data.scopedSlots, vm3.$slots) : emptyObject;
    vm3._c = function(a, b, c, d) {
      return createElement$1(vm3, a, b, c, d, false);
    };
    vm3.$createElement = function(a, b, c, d) {
      return createElement$1(vm3, a, b, c, d, true);
    };
    var parentData = parentVnode && parentVnode.data;
    if (true) {
      defineReactive(vm3, "$attrs", parentData && parentData.attrs || emptyObject, function() {
        !isUpdatingChildComponent && warn$2("$attrs is readonly.", vm3);
      }, true);
      defineReactive(vm3, "$listeners", options._parentListeners || emptyObject, function() {
        !isUpdatingChildComponent && warn$2("$listeners is readonly.", vm3);
      }, true);
    } else {
      defineReactive(vm3, "$attrs", parentData && parentData.attrs || emptyObject, null, true);
      defineReactive(vm3, "$listeners", options._parentListeners || emptyObject, null, true);
    }
  }
  var currentRenderingInstance = null;
  function renderMixin(Vue2) {
    installRenderHelpers(Vue2.prototype);
    Vue2.prototype.$nextTick = function(fn) {
      return nextTick(fn, this);
    };
    Vue2.prototype._render = function() {
      var vm3 = this;
      var _a2 = vm3.$options, render = _a2.render, _parentVnode = _a2._parentVnode;
      if (_parentVnode && vm3._isMounted) {
        vm3.$scopedSlots = normalizeScopedSlots(vm3.$parent, _parentVnode.data.scopedSlots, vm3.$slots, vm3.$scopedSlots);
        if (vm3._slotsProxy) {
          syncSetupSlots(vm3._slotsProxy, vm3.$scopedSlots);
        }
      }
      vm3.$vnode = _parentVnode;
      var prevInst = currentInstance;
      var prevRenderInst = currentRenderingInstance;
      var vnode;
      try {
        setCurrentInstance(vm3);
        currentRenderingInstance = vm3;
        vnode = render.call(vm3._renderProxy, vm3.$createElement);
      } catch (e) {
        handleError(e, vm3, "render");
        if (vm3.$options.renderError) {
          try {
            vnode = vm3.$options.renderError.call(vm3._renderProxy, vm3.$createElement, e);
          } catch (e2) {
            handleError(e2, vm3, "renderError");
            vnode = vm3._vnode;
          }
        } else {
          vnode = vm3._vnode;
        }
      } finally {
        currentRenderingInstance = prevRenderInst;
        setCurrentInstance(prevInst);
      }
      if (isArray(vnode) && vnode.length === 1) {
        vnode = vnode[0];
      }
      if (!(vnode instanceof VNode)) {
        if (isArray(vnode)) {
          warn$2("Multiple root nodes returned from render function. Render function should return a single root node.", vm3);
        }
        vnode = createEmptyVNode();
      }
      vnode.parent = _parentVnode;
      return vnode;
    };
  }
  function ensureCtor(comp, base) {
    if (comp.__esModule || hasSymbol && comp[Symbol.toStringTag] === "Module") {
      comp = comp.default;
    }
    return isObject(comp) ? base.extend(comp) : comp;
  }
  function createAsyncPlaceholder(factory, data, context, children, tag) {
    var node = createEmptyVNode();
    node.asyncFactory = factory;
    node.asyncMeta = { data, context, children, tag };
    return node;
  }
  function resolveAsyncComponent(factory, baseCtor) {
    if (isTrue(factory.error) && isDef(factory.errorComp)) {
      return factory.errorComp;
    }
    if (isDef(factory.resolved)) {
      return factory.resolved;
    }
    var owner = currentRenderingInstance;
    if (owner && isDef(factory.owners) && factory.owners.indexOf(owner) === -1) {
      factory.owners.push(owner);
    }
    if (isTrue(factory.loading) && isDef(factory.loadingComp)) {
      return factory.loadingComp;
    }
    if (owner && !isDef(factory.owners)) {
      var owners_1 = factory.owners = [owner];
      var sync_1 = true;
      var timerLoading_1 = null;
      var timerTimeout_1 = null;
      owner.$on("hook:destroyed", function() {
        return remove$2(owners_1, owner);
      });
      var forceRender_1 = function(renderCompleted) {
        for (var i = 0, l = owners_1.length; i < l; i++) {
          owners_1[i].$forceUpdate();
        }
        if (renderCompleted) {
          owners_1.length = 0;
          if (timerLoading_1 !== null) {
            clearTimeout(timerLoading_1);
            timerLoading_1 = null;
          }
          if (timerTimeout_1 !== null) {
            clearTimeout(timerTimeout_1);
            timerTimeout_1 = null;
          }
        }
      };
      var resolve = once(function(res) {
        factory.resolved = ensureCtor(res, baseCtor);
        if (!sync_1) {
          forceRender_1(true);
        } else {
          owners_1.length = 0;
        }
      });
      var reject_1 = once(function(reason) {
        warn$2("Failed to resolve async component: ".concat(String(factory)) + (reason ? "\nReason: ".concat(reason) : ""));
        if (isDef(factory.errorComp)) {
          factory.error = true;
          forceRender_1(true);
        }
      });
      var res_1 = factory(resolve, reject_1);
      if (isObject(res_1)) {
        if (isPromise(res_1)) {
          if (isUndef(factory.resolved)) {
            res_1.then(resolve, reject_1);
          }
        } else if (isPromise(res_1.component)) {
          res_1.component.then(resolve, reject_1);
          if (isDef(res_1.error)) {
            factory.errorComp = ensureCtor(res_1.error, baseCtor);
          }
          if (isDef(res_1.loading)) {
            factory.loadingComp = ensureCtor(res_1.loading, baseCtor);
            if (res_1.delay === 0) {
              factory.loading = true;
            } else {
              timerLoading_1 = setTimeout(function() {
                timerLoading_1 = null;
                if (isUndef(factory.resolved) && isUndef(factory.error)) {
                  factory.loading = true;
                  forceRender_1(false);
                }
              }, res_1.delay || 200);
            }
          }
          if (isDef(res_1.timeout)) {
            timerTimeout_1 = setTimeout(function() {
              timerTimeout_1 = null;
              if (isUndef(factory.resolved)) {
                reject_1(true ? "timeout (".concat(res_1.timeout, "ms)") : null);
              }
            }, res_1.timeout);
          }
        }
      }
      sync_1 = false;
      return factory.loading ? factory.loadingComp : factory.resolved;
    }
  }
  function getFirstComponentChild(children) {
    if (isArray(children)) {
      for (var i = 0; i < children.length; i++) {
        var c = children[i];
        if (isDef(c) && (isDef(c.componentOptions) || isAsyncPlaceholder(c))) {
          return c;
        }
      }
    }
  }
  function initEvents(vm3) {
    vm3._events = /* @__PURE__ */ Object.create(null);
    vm3._hasHookEvent = false;
    var listeners = vm3.$options._parentListeners;
    if (listeners) {
      updateComponentListeners(vm3, listeners);
    }
  }
  var target$1;
  function add$1(event, fn) {
    target$1.$on(event, fn);
  }
  function remove$1(event, fn) {
    target$1.$off(event, fn);
  }
  function createOnceHandler$1(event, fn) {
    var _target = target$1;
    return function onceHandler() {
      var res = fn.apply(null, arguments);
      if (res !== null) {
        _target.$off(event, onceHandler);
      }
    };
  }
  function updateComponentListeners(vm3, listeners, oldListeners) {
    target$1 = vm3;
    updateListeners(listeners, oldListeners || {}, add$1, remove$1, createOnceHandler$1, vm3);
    target$1 = void 0;
  }
  function eventsMixin(Vue2) {
    var hookRE = /^hook:/;
    Vue2.prototype.$on = function(event, fn) {
      var vm3 = this;
      if (isArray(event)) {
        for (var i = 0, l = event.length; i < l; i++) {
          vm3.$on(event[i], fn);
        }
      } else {
        (vm3._events[event] || (vm3._events[event] = [])).push(fn);
        if (hookRE.test(event)) {
          vm3._hasHookEvent = true;
        }
      }
      return vm3;
    };
    Vue2.prototype.$once = function(event, fn) {
      var vm3 = this;
      function on2() {
        vm3.$off(event, on2);
        fn.apply(vm3, arguments);
      }
      on2.fn = fn;
      vm3.$on(event, on2);
      return vm3;
    };
    Vue2.prototype.$off = function(event, fn) {
      var vm3 = this;
      if (!arguments.length) {
        vm3._events = /* @__PURE__ */ Object.create(null);
        return vm3;
      }
      if (isArray(event)) {
        for (var i_1 = 0, l = event.length; i_1 < l; i_1++) {
          vm3.$off(event[i_1], fn);
        }
        return vm3;
      }
      var cbs = vm3._events[event];
      if (!cbs) {
        return vm3;
      }
      if (!fn) {
        vm3._events[event] = null;
        return vm3;
      }
      var cb;
      var i = cbs.length;
      while (i--) {
        cb = cbs[i];
        if (cb === fn || cb.fn === fn) {
          cbs.splice(i, 1);
          break;
        }
      }
      return vm3;
    };
    Vue2.prototype.$emit = function(event) {
      var vm3 = this;
      if (true) {
        var lowerCaseEvent = event.toLowerCase();
        if (lowerCaseEvent !== event && vm3._events[lowerCaseEvent]) {
          tip('Event "'.concat(lowerCaseEvent, '" is emitted in component ') + "".concat(formatComponentName(vm3), ' but the handler is registered for "').concat(event, '". ') + "Note that HTML attributes are case-insensitive and you cannot use v-on to listen to camelCase events when using in-DOM templates. " + 'You should probably use "'.concat(hyphenate(event), '" instead of "').concat(event, '".'));
        }
      }
      var cbs = vm3._events[event];
      if (cbs) {
        cbs = cbs.length > 1 ? toArray(cbs) : cbs;
        var args = toArray(arguments, 1);
        var info = 'event handler for "'.concat(event, '"');
        for (var i = 0, l = cbs.length; i < l; i++) {
          invokeWithErrorHandling(cbs[i], vm3, args, vm3, info);
        }
      }
      return vm3;
    };
  }
  var activeEffectScope;
  var EffectScope = (
    /** @class */
    (function() {
      function EffectScope2(detached) {
        if (detached === void 0) {
          detached = false;
        }
        this.detached = detached;
        this.active = true;
        this.effects = [];
        this.cleanups = [];
        this.parent = activeEffectScope;
        if (!detached && activeEffectScope) {
          this.index = (activeEffectScope.scopes || (activeEffectScope.scopes = [])).push(this) - 1;
        }
      }
      EffectScope2.prototype.run = function(fn) {
        if (this.active) {
          var currentEffectScope = activeEffectScope;
          try {
            activeEffectScope = this;
            return fn();
          } finally {
            activeEffectScope = currentEffectScope;
          }
        } else if (true) {
          warn$2("cannot run an inactive effect scope.");
        }
      };
      EffectScope2.prototype.on = function() {
        activeEffectScope = this;
      };
      EffectScope2.prototype.off = function() {
        activeEffectScope = this.parent;
      };
      EffectScope2.prototype.stop = function(fromParent) {
        if (this.active) {
          var i = void 0, l = void 0;
          for (i = 0, l = this.effects.length; i < l; i++) {
            this.effects[i].teardown();
          }
          for (i = 0, l = this.cleanups.length; i < l; i++) {
            this.cleanups[i]();
          }
          if (this.scopes) {
            for (i = 0, l = this.scopes.length; i < l; i++) {
              this.scopes[i].stop(true);
            }
          }
          if (!this.detached && this.parent && !fromParent) {
            var last = this.parent.scopes.pop();
            if (last && last !== this) {
              this.parent.scopes[this.index] = last;
              last.index = this.index;
            }
          }
          this.parent = void 0;
          this.active = false;
        }
      };
      return EffectScope2;
    })()
  );
  function recordEffectScope(effect, scope) {
    if (scope === void 0) {
      scope = activeEffectScope;
    }
    if (scope && scope.active) {
      scope.effects.push(effect);
    }
  }
  function getCurrentScope() {
    return activeEffectScope;
  }
  var activeInstance = null;
  var isUpdatingChildComponent = false;
  function setActiveInstance(vm3) {
    var prevActiveInstance = activeInstance;
    activeInstance = vm3;
    return function() {
      activeInstance = prevActiveInstance;
    };
  }
  function initLifecycle(vm3) {
    var options = vm3.$options;
    var parent = options.parent;
    if (parent && !options.abstract) {
      while (parent.$options.abstract && parent.$parent) {
        parent = parent.$parent;
      }
      parent.$children.push(vm3);
    }
    vm3.$parent = parent;
    vm3.$root = parent ? parent.$root : vm3;
    vm3.$children = [];
    vm3.$refs = {};
    vm3._provided = parent ? parent._provided : /* @__PURE__ */ Object.create(null);
    vm3._watcher = null;
    vm3._inactive = null;
    vm3._directInactive = false;
    vm3._isMounted = false;
    vm3._isDestroyed = false;
    vm3._isBeingDestroyed = false;
  }
  function lifecycleMixin(Vue2) {
    Vue2.prototype._update = function(vnode, hydrating) {
      var vm3 = this;
      var prevEl = vm3.$el;
      var prevVnode = vm3._vnode;
      var restoreActiveInstance = setActiveInstance(vm3);
      vm3._vnode = vnode;
      if (!prevVnode) {
        vm3.$el = vm3.__patch__(
          vm3.$el,
          vnode,
          hydrating,
          false
          /* removeOnly */
        );
      } else {
        vm3.$el = vm3.__patch__(prevVnode, vnode);
      }
      restoreActiveInstance();
      if (prevEl) {
        prevEl.__vue__ = null;
      }
      if (vm3.$el) {
        vm3.$el.__vue__ = vm3;
      }
      var wrapper = vm3;
      while (wrapper && wrapper.$vnode && wrapper.$parent && wrapper.$vnode === wrapper.$parent._vnode) {
        wrapper.$parent.$el = wrapper.$el;
        wrapper = wrapper.$parent;
      }
    };
    Vue2.prototype.$forceUpdate = function() {
      var vm3 = this;
      if (vm3._watcher) {
        vm3._watcher.update();
      }
    };
    Vue2.prototype.$destroy = function() {
      var vm3 = this;
      if (vm3._isBeingDestroyed) {
        return;
      }
      callHook$1(vm3, "beforeDestroy");
      vm3._isBeingDestroyed = true;
      var parent = vm3.$parent;
      if (parent && !parent._isBeingDestroyed && !vm3.$options.abstract) {
        remove$2(parent.$children, vm3);
      }
      vm3._scope.stop();
      if (vm3._data.__ob__) {
        vm3._data.__ob__.vmCount--;
      }
      vm3._isDestroyed = true;
      vm3.__patch__(vm3._vnode, null);
      callHook$1(vm3, "destroyed");
      vm3.$off();
      if (vm3.$el) {
        vm3.$el.__vue__ = null;
      }
      if (vm3.$vnode) {
        vm3.$vnode.parent = null;
      }
    };
  }
  function mountComponent(vm3, el, hydrating) {
    vm3.$el = el;
    if (!vm3.$options.render) {
      vm3.$options.render = createEmptyVNode;
      if (true) {
        if (vm3.$options.template && vm3.$options.template.charAt(0) !== "#" || vm3.$options.el || el) {
          warn$2("You are using the runtime-only build of Vue where the template compiler is not available. Either pre-compile the templates into render functions, or use the compiler-included build.", vm3);
        } else {
          warn$2("Failed to mount component: template or render function not defined.", vm3);
        }
      }
    }
    callHook$1(vm3, "beforeMount");
    var updateComponent;
    if (config.performance && mark) {
      updateComponent = function() {
        var name = vm3._name;
        var id = vm3._uid;
        var startTag = "vue-perf-start:".concat(id);
        var endTag2 = "vue-perf-end:".concat(id);
        mark(startTag);
        var vnode = vm3._render();
        mark(endTag2);
        measure("vue ".concat(name, " render"), startTag, endTag2);
        mark(startTag);
        vm3._update(vnode, hydrating);
        mark(endTag2);
        measure("vue ".concat(name, " patch"), startTag, endTag2);
      };
    } else {
      updateComponent = function() {
        vm3._update(vm3._render(), hydrating);
      };
    }
    var watcherOptions = {
      before: function() {
        if (vm3._isMounted && !vm3._isDestroyed) {
          callHook$1(vm3, "beforeUpdate");
        }
      }
    };
    if (true) {
      watcherOptions.onTrack = function(e) {
        return callHook$1(vm3, "renderTracked", [e]);
      };
      watcherOptions.onTrigger = function(e) {
        return callHook$1(vm3, "renderTriggered", [e]);
      };
    }
    new Watcher(
      vm3,
      updateComponent,
      noop,
      watcherOptions,
      true
      /* isRenderWatcher */
    );
    hydrating = false;
    var preWatchers = vm3._preWatchers;
    if (preWatchers) {
      for (var i = 0; i < preWatchers.length; i++) {
        preWatchers[i].run();
      }
    }
    if (vm3.$vnode == null) {
      vm3._isMounted = true;
      callHook$1(vm3, "mounted");
    }
    return vm3;
  }
  function updateChildComponent(vm3, propsData, listeners, parentVnode, renderChildren) {
    if (true) {
      isUpdatingChildComponent = true;
    }
    var newScopedSlots = parentVnode.data.scopedSlots;
    var oldScopedSlots = vm3.$scopedSlots;
    var hasDynamicScopedSlot = !!(newScopedSlots && !newScopedSlots.$stable || oldScopedSlots !== emptyObject && !oldScopedSlots.$stable || newScopedSlots && vm3.$scopedSlots.$key !== newScopedSlots.$key || !newScopedSlots && vm3.$scopedSlots.$key);
    var needsForceUpdate = !!(renderChildren || // has new static slots
    vm3.$options._renderChildren || // has old static slots
    hasDynamicScopedSlot);
    var prevVNode = vm3.$vnode;
    vm3.$options._parentVnode = parentVnode;
    vm3.$vnode = parentVnode;
    if (vm3._vnode) {
      vm3._vnode.parent = parentVnode;
    }
    vm3.$options._renderChildren = renderChildren;
    var attrs2 = parentVnode.data.attrs || emptyObject;
    if (vm3._attrsProxy) {
      if (syncSetupProxy(vm3._attrsProxy, attrs2, prevVNode.data && prevVNode.data.attrs || emptyObject, vm3, "$attrs")) {
        needsForceUpdate = true;
      }
    }
    vm3.$attrs = attrs2;
    listeners = listeners || emptyObject;
    var prevListeners = vm3.$options._parentListeners;
    if (vm3._listenersProxy) {
      syncSetupProxy(vm3._listenersProxy, listeners, prevListeners || emptyObject, vm3, "$listeners");
    }
    vm3.$listeners = vm3.$options._parentListeners = listeners;
    updateComponentListeners(vm3, listeners, prevListeners);
    if (propsData && vm3.$options.props) {
      toggleObserving(false);
      var props2 = vm3._props;
      var propKeys = vm3.$options._propKeys || [];
      for (var i = 0; i < propKeys.length; i++) {
        var key = propKeys[i];
        var propOptions = vm3.$options.props;
        props2[key] = validateProp(key, propOptions, propsData, vm3);
      }
      toggleObserving(true);
      vm3.$options.propsData = propsData;
    }
    if (needsForceUpdate) {
      vm3.$slots = resolveSlots(renderChildren, parentVnode.context);
      vm3.$forceUpdate();
    }
    if (true) {
      isUpdatingChildComponent = false;
    }
  }
  function isInInactiveTree(vm3) {
    while (vm3 && (vm3 = vm3.$parent)) {
      if (vm3._inactive)
        return true;
    }
    return false;
  }
  function activateChildComponent(vm3, direct) {
    if (direct) {
      vm3._directInactive = false;
      if (isInInactiveTree(vm3)) {
        return;
      }
    } else if (vm3._directInactive) {
      return;
    }
    if (vm3._inactive || vm3._inactive === null) {
      vm3._inactive = false;
      for (var i = 0; i < vm3.$children.length; i++) {
        activateChildComponent(vm3.$children[i]);
      }
      callHook$1(vm3, "activated");
    }
  }
  function deactivateChildComponent(vm3, direct) {
    if (direct) {
      vm3._directInactive = true;
      if (isInInactiveTree(vm3)) {
        return;
      }
    }
    if (!vm3._inactive) {
      vm3._inactive = true;
      for (var i = 0; i < vm3.$children.length; i++) {
        deactivateChildComponent(vm3.$children[i]);
      }
      callHook$1(vm3, "deactivated");
    }
  }
  function callHook$1(vm3, hook, args, setContext) {
    if (setContext === void 0) {
      setContext = true;
    }
    pushTarget();
    var prevInst = currentInstance;
    var prevScope = getCurrentScope();
    setContext && setCurrentInstance(vm3);
    var handlers = vm3.$options[hook];
    var info = "".concat(hook, " hook");
    if (handlers) {
      for (var i = 0, j = handlers.length; i < j; i++) {
        invokeWithErrorHandling(handlers[i], vm3, args || null, vm3, info);
      }
    }
    if (vm3._hasHookEvent) {
      vm3.$emit("hook:" + hook);
    }
    if (setContext) {
      setCurrentInstance(prevInst);
      prevScope && prevScope.on();
    }
    popTarget();
  }
  var MAX_UPDATE_COUNT = 100;
  var queue = [];
  var activatedChildren = [];
  var has = {};
  var circular = {};
  var waiting = false;
  var flushing = false;
  var index$1 = 0;
  function resetSchedulerState() {
    index$1 = queue.length = activatedChildren.length = 0;
    has = {};
    if (true) {
      circular = {};
    }
    waiting = flushing = false;
  }
  var currentFlushTimestamp = 0;
  var getNow = Date.now;
  if (inBrowser && !isIE) {
    performance_1 = window.performance;
    if (performance_1 && typeof performance_1.now === "function" && getNow() > document.createEvent("Event").timeStamp) {
      getNow = function() {
        return performance_1.now();
      };
    }
  }
  var performance_1;
  var sortCompareFn = function(a, b) {
    if (a.post) {
      if (!b.post)
        return 1;
    } else if (b.post) {
      return -1;
    }
    return a.id - b.id;
  };
  function flushSchedulerQueue() {
    currentFlushTimestamp = getNow();
    flushing = true;
    var watcher, id;
    queue.sort(sortCompareFn);
    for (index$1 = 0; index$1 < queue.length; index$1++) {
      watcher = queue[index$1];
      if (watcher.before) {
        watcher.before();
      }
      id = watcher.id;
      has[id] = null;
      watcher.run();
      if (has[id] != null) {
        circular[id] = (circular[id] || 0) + 1;
        if (circular[id] > MAX_UPDATE_COUNT) {
          warn$2("You may have an infinite update loop " + (watcher.user ? 'in watcher with expression "'.concat(watcher.expression, '"') : "in a component render function."), watcher.vm);
          break;
        }
      }
    }
    var activatedQueue = activatedChildren.slice();
    var updatedQueue = queue.slice();
    resetSchedulerState();
    callActivatedHooks(activatedQueue);
    callUpdatedHooks(updatedQueue);
    cleanupDeps();
    if (devtools && config.devtools) {
      devtools.emit("flush");
    }
  }
  function callUpdatedHooks(queue2) {
    var i = queue2.length;
    while (i--) {
      var watcher = queue2[i];
      var vm3 = watcher.vm;
      if (vm3 && vm3._watcher === watcher && vm3._isMounted && !vm3._isDestroyed) {
        callHook$1(vm3, "updated");
      }
    }
  }
  function queueActivatedComponent(vm3) {
    vm3._inactive = false;
    activatedChildren.push(vm3);
  }
  function callActivatedHooks(queue2) {
    for (var i = 0; i < queue2.length; i++) {
      queue2[i]._inactive = true;
      activateChildComponent(
        queue2[i],
        true
        /* true */
      );
    }
  }
  function queueWatcher(watcher) {
    var id = watcher.id;
    if (has[id] != null) {
      return;
    }
    if (watcher === Dep.target && watcher.noRecurse) {
      return;
    }
    has[id] = true;
    if (!flushing) {
      queue.push(watcher);
    } else {
      var i = queue.length - 1;
      while (i > index$1 && queue[i].id > watcher.id) {
        i--;
      }
      queue.splice(i + 1, 0, watcher);
    }
    if (!waiting) {
      waiting = true;
      if (!config.async) {
        flushSchedulerQueue();
        return;
      }
      nextTick(flushSchedulerQueue);
    }
  }
  var WATCHER = "watcher";
  var WATCHER_CB = "".concat(WATCHER, " callback");
  var WATCHER_GETTER = "".concat(WATCHER, " getter");
  var WATCHER_CLEANUP = "".concat(WATCHER, " cleanup");
  function resolveProvided(vm3) {
    var existing = vm3._provided;
    var parentProvides = vm3.$parent && vm3.$parent._provided;
    if (parentProvides === existing) {
      return vm3._provided = Object.create(parentProvides);
    } else {
      return existing;
    }
  }
  function handleError(err, vm3, info) {
    pushTarget();
    try {
      if (vm3) {
        var cur = vm3;
        while (cur = cur.$parent) {
          var hooks2 = cur.$options.errorCaptured;
          if (hooks2) {
            for (var i = 0; i < hooks2.length; i++) {
              try {
                var capture = hooks2[i].call(cur, err, vm3, info) === false;
                if (capture)
                  return;
              } catch (e) {
                globalHandleError(e, cur, "errorCaptured hook");
              }
            }
          }
        }
      }
      globalHandleError(err, vm3, info);
    } finally {
      popTarget();
    }
  }
  function invokeWithErrorHandling(handler, context, args, vm3, info) {
    var res;
    try {
      res = args ? handler.apply(context, args) : handler.call(context);
      if (res && !res._isVue && isPromise(res) && !res._handled) {
        res.catch(function(e) {
          return handleError(e, vm3, info + " (Promise/async)");
        });
        res._handled = true;
      }
    } catch (e) {
      handleError(e, vm3, info);
    }
    return res;
  }
  function globalHandleError(err, vm3, info) {
    if (config.errorHandler) {
      try {
        return config.errorHandler.call(null, err, vm3, info);
      } catch (e) {
        if (e !== err) {
          logError(e, null, "config.errorHandler");
        }
      }
    }
    logError(err, vm3, info);
  }
  function logError(err, vm3, info) {
    if (true) {
      warn$2("Error in ".concat(info, ': "').concat(err.toString(), '"'), vm3);
    }
    if (inBrowser && typeof console !== "undefined") {
      console.error(err);
    } else {
      throw err;
    }
  }
  var isUsingMicroTask = false;
  var callbacks = [];
  var pending = false;
  function flushCallbacks() {
    pending = false;
    var copies = callbacks.slice(0);
    callbacks.length = 0;
    for (var i = 0; i < copies.length; i++) {
      copies[i]();
    }
  }
  var timerFunc;
  if (typeof Promise !== "undefined" && isNative(Promise)) {
    p_1 = Promise.resolve();
    timerFunc = function() {
      p_1.then(flushCallbacks);
      if (isIOS)
        setTimeout(noop);
    };
    isUsingMicroTask = true;
  } else if (!isIE && typeof MutationObserver !== "undefined" && (isNative(MutationObserver) || // PhantomJS and iOS 7.x
  MutationObserver.toString() === "[object MutationObserverConstructor]")) {
    counter_1 = 1;
    observer = new MutationObserver(flushCallbacks);
    textNode_1 = document.createTextNode(String(counter_1));
    observer.observe(textNode_1, {
      characterData: true
    });
    timerFunc = function() {
      counter_1 = (counter_1 + 1) % 2;
      textNode_1.data = String(counter_1);
    };
    isUsingMicroTask = true;
  } else if (typeof setImmediate !== "undefined" && isNative(setImmediate)) {
    timerFunc = function() {
      setImmediate(flushCallbacks);
    };
  } else {
    timerFunc = function() {
      setTimeout(flushCallbacks, 0);
    };
  }
  var p_1;
  var counter_1;
  var observer;
  var textNode_1;
  function nextTick(cb, ctx) {
    var _resolve;
    callbacks.push(function() {
      if (cb) {
        try {
          cb.call(ctx);
        } catch (e) {
          handleError(e, ctx, "nextTick");
        }
      } else if (_resolve) {
        _resolve(ctx);
      }
    });
    if (!pending) {
      pending = true;
      timerFunc();
    }
    if (!cb && typeof Promise !== "undefined") {
      return new Promise(function(resolve) {
        _resolve = resolve;
      });
    }
  }
  function createLifeCycle(hookName) {
    return function(fn, target2) {
      if (target2 === void 0) {
        target2 = currentInstance;
      }
      if (!target2) {
        warn$2("".concat(formatName(hookName), " is called when there is no active component instance to be ") + "associated with. Lifecycle injection APIs can only be used during execution of setup().");
        return;
      }
      return injectHook(target2, hookName, fn);
    };
  }
  function formatName(name) {
    if (name === "beforeDestroy") {
      name = "beforeUnmount";
    } else if (name === "destroyed") {
      name = "unmounted";
    }
    return "on".concat(name[0].toUpperCase() + name.slice(1));
  }
  function injectHook(instance, hookName, fn) {
    var options = instance.$options;
    options[hookName] = mergeLifecycleHook(options[hookName], fn);
  }
  var onBeforeMount = createLifeCycle("beforeMount");
  var onMounted = createLifeCycle("mounted");
  var onBeforeUpdate = createLifeCycle("beforeUpdate");
  var onUpdated = createLifeCycle("updated");
  var onBeforeUnmount = createLifeCycle("beforeDestroy");
  var onUnmounted = createLifeCycle("destroyed");
  var onActivated = createLifeCycle("activated");
  var onDeactivated = createLifeCycle("deactivated");
  var onServerPrefetch = createLifeCycle("serverPrefetch");
  var onRenderTracked = createLifeCycle("renderTracked");
  var onRenderTriggered = createLifeCycle("renderTriggered");
  var injectErrorCapturedHook = createLifeCycle("errorCaptured");
  var version = "2.7.16";
  var seenObjects = new _Set();
  function traverse(val) {
    _traverse(val, seenObjects);
    seenObjects.clear();
    return val;
  }
  function _traverse(val, seen) {
    var i, keys;
    var isA = isArray(val);
    if (!isA && !isObject(val) || val.__v_skip || Object.isFrozen(val) || val instanceof VNode) {
      return;
    }
    if (val.__ob__) {
      var depId = val.__ob__.dep.id;
      if (seen.has(depId)) {
        return;
      }
      seen.add(depId);
    }
    if (isA) {
      i = val.length;
      while (i--)
        _traverse(val[i], seen);
    } else if (isRef(val)) {
      _traverse(val.value, seen);
    } else {
      keys = Object.keys(val);
      i = keys.length;
      while (i--)
        _traverse(val[keys[i]], seen);
    }
  }
  var uid$1 = 0;
  var Watcher = (
    /** @class */
    (function() {
      function Watcher2(vm3, expOrFn, cb, options, isRenderWatcher) {
        recordEffectScope(
          this,
          // if the active effect scope is manually created (not a component scope),
          // prioritize it
          activeEffectScope && !activeEffectScope._vm ? activeEffectScope : vm3 ? vm3._scope : void 0
        );
        if ((this.vm = vm3) && isRenderWatcher) {
          vm3._watcher = this;
        }
        if (options) {
          this.deep = !!options.deep;
          this.user = !!options.user;
          this.lazy = !!options.lazy;
          this.sync = !!options.sync;
          this.before = options.before;
          if (true) {
            this.onTrack = options.onTrack;
            this.onTrigger = options.onTrigger;
          }
        } else {
          this.deep = this.user = this.lazy = this.sync = false;
        }
        this.cb = cb;
        this.id = ++uid$1;
        this.active = true;
        this.post = false;
        this.dirty = this.lazy;
        this.deps = [];
        this.newDeps = [];
        this.depIds = new _Set();
        this.newDepIds = new _Set();
        this.expression = true ? expOrFn.toString() : "";
        if (isFunction(expOrFn)) {
          this.getter = expOrFn;
        } else {
          this.getter = parsePath(expOrFn);
          if (!this.getter) {
            this.getter = noop;
            warn$2('Failed watching path: "'.concat(expOrFn, '" ') + "Watcher only accepts simple dot-delimited paths. For full control, use a function instead.", vm3);
          }
        }
        this.value = this.lazy ? void 0 : this.get();
      }
      Watcher2.prototype.get = function() {
        pushTarget(this);
        var value;
        var vm3 = this.vm;
        try {
          value = this.getter.call(vm3, vm3);
        } catch (e) {
          if (this.user) {
            handleError(e, vm3, 'getter for watcher "'.concat(this.expression, '"'));
          } else {
            throw e;
          }
        } finally {
          if (this.deep) {
            traverse(value);
          }
          popTarget();
          this.cleanupDeps();
        }
        return value;
      };
      Watcher2.prototype.addDep = function(dep) {
        var id = dep.id;
        if (!this.newDepIds.has(id)) {
          this.newDepIds.add(id);
          this.newDeps.push(dep);
          if (!this.depIds.has(id)) {
            dep.addSub(this);
          }
        }
      };
      Watcher2.prototype.cleanupDeps = function() {
        var i = this.deps.length;
        while (i--) {
          var dep = this.deps[i];
          if (!this.newDepIds.has(dep.id)) {
            dep.removeSub(this);
          }
        }
        var tmp = this.depIds;
        this.depIds = this.newDepIds;
        this.newDepIds = tmp;
        this.newDepIds.clear();
        tmp = this.deps;
        this.deps = this.newDeps;
        this.newDeps = tmp;
        this.newDeps.length = 0;
      };
      Watcher2.prototype.update = function() {
        if (this.lazy) {
          this.dirty = true;
        } else if (this.sync) {
          this.run();
        } else {
          queueWatcher(this);
        }
      };
      Watcher2.prototype.run = function() {
        if (this.active) {
          var value = this.get();
          if (value !== this.value || // Deep watchers and watchers on Object/Arrays should fire even
          // when the value is the same, because the value may
          // have mutated.
          isObject(value) || this.deep) {
            var oldValue = this.value;
            this.value = value;
            if (this.user) {
              var info = 'callback for watcher "'.concat(this.expression, '"');
              invokeWithErrorHandling(this.cb, this.vm, [value, oldValue], this.vm, info);
            } else {
              this.cb.call(this.vm, value, oldValue);
            }
          }
        }
      };
      Watcher2.prototype.evaluate = function() {
        this.value = this.get();
        this.dirty = false;
      };
      Watcher2.prototype.depend = function() {
        var i = this.deps.length;
        while (i--) {
          this.deps[i].depend();
        }
      };
      Watcher2.prototype.teardown = function() {
        if (this.vm && !this.vm._isBeingDestroyed) {
          remove$2(this.vm._scope.effects, this);
        }
        if (this.active) {
          var i = this.deps.length;
          while (i--) {
            this.deps[i].removeSub(this);
          }
          this.active = false;
          if (this.onStop) {
            this.onStop();
          }
        }
      };
      return Watcher2;
    })()
  );
  var sharedPropertyDefinition = {
    enumerable: true,
    configurable: true,
    get: noop,
    set: noop
  };
  function proxy(target2, sourceKey, key) {
    sharedPropertyDefinition.get = function proxyGetter() {
      return this[sourceKey][key];
    };
    sharedPropertyDefinition.set = function proxySetter(val) {
      this[sourceKey][key] = val;
    };
    Object.defineProperty(target2, key, sharedPropertyDefinition);
  }
  function initState(vm3) {
    var opts2 = vm3.$options;
    if (opts2.props)
      initProps$1(vm3, opts2.props);
    initSetup(vm3);
    if (opts2.methods)
      initMethods(vm3, opts2.methods);
    if (opts2.data) {
      initData(vm3);
    } else {
      var ob = observe(vm3._data = {});
      ob && ob.vmCount++;
    }
    if (opts2.computed)
      initComputed$1(vm3, opts2.computed);
    if (opts2.watch && opts2.watch !== nativeWatch) {
      initWatch(vm3, opts2.watch);
    }
  }
  function initProps$1(vm3, propsOptions) {
    var propsData = vm3.$options.propsData || {};
    var props2 = vm3._props = shallowReactive({});
    var keys = vm3.$options._propKeys = [];
    var isRoot = !vm3.$parent;
    if (!isRoot) {
      toggleObserving(false);
    }
    var _loop_1 = function(key2) {
      keys.push(key2);
      var value = validateProp(key2, propsOptions, propsData, vm3);
      if (true) {
        var hyphenatedKey = hyphenate(key2);
        if (isReservedAttribute(hyphenatedKey) || config.isReservedAttr(hyphenatedKey)) {
          warn$2('"'.concat(hyphenatedKey, '" is a reserved attribute and cannot be used as component prop.'), vm3);
        }
        defineReactive(
          props2,
          key2,
          value,
          function() {
            if (!isRoot && !isUpdatingChildComponent) {
              warn$2("Avoid mutating a prop directly since the value will be overwritten whenever the parent component re-renders. Instead, use a data or computed property based on the prop's " + 'value. Prop being mutated: "'.concat(key2, '"'), vm3);
            }
          },
          true
          /* shallow */
        );
      } else {
        defineReactive(
          props2,
          key2,
          value,
          void 0,
          true
          /* shallow */
        );
      }
      if (!(key2 in vm3)) {
        proxy(vm3, "_props", key2);
      }
    };
    for (var key in propsOptions) {
      _loop_1(key);
    }
    toggleObserving(true);
  }
  function initData(vm3) {
    var data = vm3.$options.data;
    data = vm3._data = isFunction(data) ? getData(data, vm3) : data || {};
    if (!isPlainObject(data)) {
      data = {};
      warn$2("data functions should return an object:\nhttps://v2.vuejs.org/v2/guide/components.html#data-Must-Be-a-Function", vm3);
    }
    var keys = Object.keys(data);
    var props2 = vm3.$options.props;
    var methods = vm3.$options.methods;
    var i = keys.length;
    while (i--) {
      var key = keys[i];
      if (true) {
        if (methods && hasOwn(methods, key)) {
          warn$2('Method "'.concat(key, '" has already been defined as a data property.'), vm3);
        }
      }
      if (props2 && hasOwn(props2, key)) {
        warn$2('The data property "'.concat(key, '" is already declared as a prop. ') + "Use prop default value instead.", vm3);
      } else if (!isReserved(key)) {
        proxy(vm3, "_data", key);
      }
    }
    var ob = observe(data);
    ob && ob.vmCount++;
  }
  function getData(data, vm3) {
    pushTarget();
    try {
      return data.call(vm3, vm3);
    } catch (e) {
      handleError(e, vm3, "data()");
      return {};
    } finally {
      popTarget();
    }
  }
  var computedWatcherOptions = { lazy: true };
  function initComputed$1(vm3, computed) {
    var watchers = vm3._computedWatchers = /* @__PURE__ */ Object.create(null);
    var isSSR = isServerRendering();
    for (var key in computed) {
      var userDef = computed[key];
      var getter = isFunction(userDef) ? userDef : userDef.get;
      if (getter == null) {
        warn$2('Getter is missing for computed property "'.concat(key, '".'), vm3);
      }
      if (!isSSR) {
        watchers[key] = new Watcher(vm3, getter || noop, noop, computedWatcherOptions);
      }
      if (!(key in vm3)) {
        defineComputed(vm3, key, userDef);
      } else if (true) {
        if (key in vm3.$data) {
          warn$2('The computed property "'.concat(key, '" is already defined in data.'), vm3);
        } else if (vm3.$options.props && key in vm3.$options.props) {
          warn$2('The computed property "'.concat(key, '" is already defined as a prop.'), vm3);
        } else if (vm3.$options.methods && key in vm3.$options.methods) {
          warn$2('The computed property "'.concat(key, '" is already defined as a method.'), vm3);
        }
      }
    }
  }
  function defineComputed(target2, key, userDef) {
    var shouldCache = !isServerRendering();
    if (isFunction(userDef)) {
      sharedPropertyDefinition.get = shouldCache ? createComputedGetter(key) : createGetterInvoker(userDef);
      sharedPropertyDefinition.set = noop;
    } else {
      sharedPropertyDefinition.get = userDef.get ? shouldCache && userDef.cache !== false ? createComputedGetter(key) : createGetterInvoker(userDef.get) : noop;
      sharedPropertyDefinition.set = userDef.set || noop;
    }
    if (sharedPropertyDefinition.set === noop) {
      sharedPropertyDefinition.set = function() {
        warn$2('Computed property "'.concat(key, '" was assigned to but it has no setter.'), this);
      };
    }
    Object.defineProperty(target2, key, sharedPropertyDefinition);
  }
  function createComputedGetter(key) {
    return function computedGetter() {
      var watcher = this._computedWatchers && this._computedWatchers[key];
      if (watcher) {
        if (watcher.dirty) {
          watcher.evaluate();
        }
        if (Dep.target) {
          if (Dep.target.onTrack) {
            Dep.target.onTrack({
              effect: Dep.target,
              target: this,
              type: "get",
              key
            });
          }
          watcher.depend();
        }
        return watcher.value;
      }
    };
  }
  function createGetterInvoker(fn) {
    return function computedGetter() {
      return fn.call(this, this);
    };
  }
  function initMethods(vm3, methods) {
    var props2 = vm3.$options.props;
    for (var key in methods) {
      if (true) {
        if (typeof methods[key] !== "function") {
          warn$2('Method "'.concat(key, '" has type "').concat(typeof methods[key], '" in the component definition. ') + "Did you reference the function correctly?", vm3);
        }
        if (props2 && hasOwn(props2, key)) {
          warn$2('Method "'.concat(key, '" has already been defined as a prop.'), vm3);
        }
        if (key in vm3 && isReserved(key)) {
          warn$2('Method "'.concat(key, '" conflicts with an existing Vue instance method. ') + "Avoid defining component methods that start with _ or $.");
        }
      }
      vm3[key] = typeof methods[key] !== "function" ? noop : bind$1(methods[key], vm3);
    }
  }
  function initWatch(vm3, watch) {
    for (var key in watch) {
      var handler = watch[key];
      if (isArray(handler)) {
        for (var i = 0; i < handler.length; i++) {
          createWatcher(vm3, key, handler[i]);
        }
      } else {
        createWatcher(vm3, key, handler);
      }
    }
  }
  function createWatcher(vm3, expOrFn, handler, options) {
    if (isPlainObject(handler)) {
      options = handler;
      handler = handler.handler;
    }
    if (typeof handler === "string") {
      handler = vm3[handler];
    }
    return vm3.$watch(expOrFn, handler, options);
  }
  function stateMixin(Vue2) {
    var dataDef = {};
    dataDef.get = function() {
      return this._data;
    };
    var propsDef = {};
    propsDef.get = function() {
      return this._props;
    };
    if (true) {
      dataDef.set = function() {
        warn$2("Avoid replacing instance root $data. Use nested data properties instead.", this);
      };
      propsDef.set = function() {
        warn$2("$props is readonly.", this);
      };
    }
    Object.defineProperty(Vue2.prototype, "$data", dataDef);
    Object.defineProperty(Vue2.prototype, "$props", propsDef);
    Vue2.prototype.$set = set;
    Vue2.prototype.$delete = del;
    Vue2.prototype.$watch = function(expOrFn, cb, options) {
      var vm3 = this;
      if (isPlainObject(cb)) {
        return createWatcher(vm3, expOrFn, cb, options);
      }
      options = options || {};
      options.user = true;
      var watcher = new Watcher(vm3, expOrFn, cb, options);
      if (options.immediate) {
        var info = 'callback for immediate watcher "'.concat(watcher.expression, '"');
        pushTarget();
        invokeWithErrorHandling(cb, vm3, [watcher.value], vm3, info);
        popTarget();
      }
      return function unwatchFn() {
        watcher.teardown();
      };
    };
  }
  function initProvide(vm3) {
    var provideOption = vm3.$options.provide;
    if (provideOption) {
      var provided = isFunction(provideOption) ? provideOption.call(vm3) : provideOption;
      if (!isObject(provided)) {
        return;
      }
      var source = resolveProvided(vm3);
      var keys = hasSymbol ? Reflect.ownKeys(provided) : Object.keys(provided);
      for (var i = 0; i < keys.length; i++) {
        var key = keys[i];
        Object.defineProperty(source, key, Object.getOwnPropertyDescriptor(provided, key));
      }
    }
  }
  function initInjections(vm3) {
    var result = resolveInject(vm3.$options.inject, vm3);
    if (result) {
      toggleObserving(false);
      Object.keys(result).forEach(function(key) {
        if (true) {
          defineReactive(vm3, key, result[key], function() {
            warn$2("Avoid mutating an injected value directly since the changes will be overwritten whenever the provided component re-renders. " + 'injection being mutated: "'.concat(key, '"'), vm3);
          });
        } else {
          defineReactive(vm3, key, result[key]);
        }
      });
      toggleObserving(true);
    }
  }
  function resolveInject(inject, vm3) {
    if (inject) {
      var result = /* @__PURE__ */ Object.create(null);
      var keys = hasSymbol ? Reflect.ownKeys(inject) : Object.keys(inject);
      for (var i = 0; i < keys.length; i++) {
        var key = keys[i];
        if (key === "__ob__")
          continue;
        var provideKey = inject[key].from;
        if (provideKey in vm3._provided) {
          result[key] = vm3._provided[provideKey];
        } else if ("default" in inject[key]) {
          var provideDefault = inject[key].default;
          result[key] = isFunction(provideDefault) ? provideDefault.call(vm3) : provideDefault;
        } else if (true) {
          warn$2('Injection "'.concat(key, '" not found'), vm3);
        }
      }
      return result;
    }
  }
  var uid = 0;
  function initMixin$1(Vue2) {
    Vue2.prototype._init = function(options) {
      var vm3 = this;
      vm3._uid = uid++;
      var startTag, endTag2;
      if (config.performance && mark) {
        startTag = "vue-perf-start:".concat(vm3._uid);
        endTag2 = "vue-perf-end:".concat(vm3._uid);
        mark(startTag);
      }
      vm3._isVue = true;
      vm3.__v_skip = true;
      vm3._scope = new EffectScope(
        true
        /* detached */
      );
      vm3._scope.parent = void 0;
      vm3._scope._vm = true;
      if (options && options._isComponent) {
        initInternalComponent(vm3, options);
      } else {
        vm3.$options = mergeOptions(resolveConstructorOptions(vm3.constructor), options || {}, vm3);
      }
      if (true) {
        initProxy(vm3);
      } else {
        vm3._renderProxy = vm3;
      }
      vm3._self = vm3;
      initLifecycle(vm3);
      initEvents(vm3);
      initRender(vm3);
      callHook$1(
        vm3,
        "beforeCreate",
        void 0,
        false
        /* setContext */
      );
      initInjections(vm3);
      initState(vm3);
      initProvide(vm3);
      callHook$1(vm3, "created");
      if (config.performance && mark) {
        vm3._name = formatComponentName(vm3, false);
        mark(endTag2);
        measure("vue ".concat(vm3._name, " init"), startTag, endTag2);
      }
      if (vm3.$options.el) {
        vm3.$mount(vm3.$options.el);
      }
    };
  }
  function initInternalComponent(vm3, options) {
    var opts2 = vm3.$options = Object.create(vm3.constructor.options);
    var parentVnode = options._parentVnode;
    opts2.parent = options.parent;
    opts2._parentVnode = parentVnode;
    var vnodeComponentOptions = parentVnode.componentOptions;
    opts2.propsData = vnodeComponentOptions.propsData;
    opts2._parentListeners = vnodeComponentOptions.listeners;
    opts2._renderChildren = vnodeComponentOptions.children;
    opts2._componentTag = vnodeComponentOptions.tag;
    if (options.render) {
      opts2.render = options.render;
      opts2.staticRenderFns = options.staticRenderFns;
    }
  }
  function resolveConstructorOptions(Ctor) {
    var options = Ctor.options;
    if (Ctor.super) {
      var superOptions = resolveConstructorOptions(Ctor.super);
      var cachedSuperOptions = Ctor.superOptions;
      if (superOptions !== cachedSuperOptions) {
        Ctor.superOptions = superOptions;
        var modifiedOptions = resolveModifiedOptions(Ctor);
        if (modifiedOptions) {
          extend(Ctor.extendOptions, modifiedOptions);
        }
        options = Ctor.options = mergeOptions(superOptions, Ctor.extendOptions);
        if (options.name) {
          options.components[options.name] = Ctor;
        }
      }
    }
    return options;
  }
  function resolveModifiedOptions(Ctor) {
    var modified;
    var latest = Ctor.options;
    var sealed = Ctor.sealedOptions;
    for (var key in latest) {
      if (latest[key] !== sealed[key]) {
        if (!modified)
          modified = {};
        modified[key] = latest[key];
      }
    }
    return modified;
  }
  function FunctionalRenderContext(data, props2, children, parent, Ctor) {
    var _this = this;
    var options = Ctor.options;
    var contextVm;
    if (hasOwn(parent, "_uid")) {
      contextVm = Object.create(parent);
      contextVm._original = parent;
    } else {
      contextVm = parent;
      parent = parent._original;
    }
    var isCompiled = isTrue(options._compiled);
    var needNormalization = !isCompiled;
    this.data = data;
    this.props = props2;
    this.children = children;
    this.parent = parent;
    this.listeners = data.on || emptyObject;
    this.injections = resolveInject(options.inject, parent);
    this.slots = function() {
      if (!_this.$slots) {
        normalizeScopedSlots(parent, data.scopedSlots, _this.$slots = resolveSlots(children, parent));
      }
      return _this.$slots;
    };
    Object.defineProperty(this, "scopedSlots", {
      enumerable: true,
      get: function() {
        return normalizeScopedSlots(parent, data.scopedSlots, this.slots());
      }
    });
    if (isCompiled) {
      this.$options = options;
      this.$slots = this.slots();
      this.$scopedSlots = normalizeScopedSlots(parent, data.scopedSlots, this.$slots);
    }
    if (options._scopeId) {
      this._c = function(a, b, c, d) {
        var vnode = createElement$1(contextVm, a, b, c, d, needNormalization);
        if (vnode && !isArray(vnode)) {
          vnode.fnScopeId = options._scopeId;
          vnode.fnContext = parent;
        }
        return vnode;
      };
    } else {
      this._c = function(a, b, c, d) {
        return createElement$1(contextVm, a, b, c, d, needNormalization);
      };
    }
  }
  installRenderHelpers(FunctionalRenderContext.prototype);
  function createFunctionalComponent(Ctor, propsData, data, contextVm, children) {
    var options = Ctor.options;
    var props2 = {};
    var propOptions = options.props;
    if (isDef(propOptions)) {
      for (var key in propOptions) {
        props2[key] = validateProp(key, propOptions, propsData || emptyObject);
      }
    } else {
      if (isDef(data.attrs))
        mergeProps(props2, data.attrs);
      if (isDef(data.props))
        mergeProps(props2, data.props);
    }
    var renderContext = new FunctionalRenderContext(data, props2, children, contextVm, Ctor);
    var vnode = options.render.call(null, renderContext._c, renderContext);
    if (vnode instanceof VNode) {
      return cloneAndMarkFunctionalResult(vnode, data, renderContext.parent, options, renderContext);
    } else if (isArray(vnode)) {
      var vnodes = normalizeChildren(vnode) || [];
      var res = new Array(vnodes.length);
      for (var i = 0; i < vnodes.length; i++) {
        res[i] = cloneAndMarkFunctionalResult(vnodes[i], data, renderContext.parent, options, renderContext);
      }
      return res;
    }
  }
  function cloneAndMarkFunctionalResult(vnode, data, contextVm, options, renderContext) {
    var clone = cloneVNode(vnode);
    clone.fnContext = contextVm;
    clone.fnOptions = options;
    if (true) {
      (clone.devtoolsMeta = clone.devtoolsMeta || {}).renderContext = renderContext;
    }
    if (data.slot) {
      (clone.data || (clone.data = {})).slot = data.slot;
    }
    return clone;
  }
  function mergeProps(to, from) {
    for (var key in from) {
      to[camelize(key)] = from[key];
    }
  }
  function getComponentName(options) {
    return options.name || options.__name || options._componentTag;
  }
  var componentVNodeHooks = {
    init: function(vnode, hydrating) {
      if (vnode.componentInstance && !vnode.componentInstance._isDestroyed && vnode.data.keepAlive) {
        var mountedNode = vnode;
        componentVNodeHooks.prepatch(mountedNode, mountedNode);
      } else {
        var child = vnode.componentInstance = createComponentInstanceForVnode(vnode, activeInstance);
        child.$mount(hydrating ? vnode.elm : void 0, hydrating);
      }
    },
    prepatch: function(oldVnode, vnode) {
      var options = vnode.componentOptions;
      var child = vnode.componentInstance = oldVnode.componentInstance;
      updateChildComponent(
        child,
        options.propsData,
        // updated props
        options.listeners,
        // updated listeners
        vnode,
        // new parent vnode
        options.children
        // new children
      );
    },
    insert: function(vnode) {
      var context = vnode.context, componentInstance = vnode.componentInstance;
      if (!componentInstance._isMounted) {
        componentInstance._isMounted = true;
        callHook$1(componentInstance, "mounted");
      }
      if (vnode.data.keepAlive) {
        if (context._isMounted) {
          queueActivatedComponent(componentInstance);
        } else {
          activateChildComponent(
            componentInstance,
            true
            /* direct */
          );
        }
      }
    },
    destroy: function(vnode) {
      var componentInstance = vnode.componentInstance;
      if (!componentInstance._isDestroyed) {
        if (!vnode.data.keepAlive) {
          componentInstance.$destroy();
        } else {
          deactivateChildComponent(
            componentInstance,
            true
            /* direct */
          );
        }
      }
    }
  };
  var hooksToMerge = Object.keys(componentVNodeHooks);
  function createComponent(Ctor, data, context, children, tag) {
    if (isUndef(Ctor)) {
      return;
    }
    var baseCtor = context.$options._base;
    if (isObject(Ctor)) {
      Ctor = baseCtor.extend(Ctor);
    }
    if (typeof Ctor !== "function") {
      if (true) {
        warn$2("Invalid Component definition: ".concat(String(Ctor)), context);
      }
      return;
    }
    var asyncFactory;
    if (isUndef(Ctor.cid)) {
      asyncFactory = Ctor;
      Ctor = resolveAsyncComponent(asyncFactory, baseCtor);
      if (Ctor === void 0) {
        return createAsyncPlaceholder(asyncFactory, data, context, children, tag);
      }
    }
    data = data || {};
    resolveConstructorOptions(Ctor);
    if (isDef(data.model)) {
      transformModel(Ctor.options, data);
    }
    var propsData = extractPropsFromVNodeData(data, Ctor, tag);
    if (isTrue(Ctor.options.functional)) {
      return createFunctionalComponent(Ctor, propsData, data, context, children);
    }
    var listeners = data.on;
    data.on = data.nativeOn;
    if (isTrue(Ctor.options.abstract)) {
      var slot = data.slot;
      data = {};
      if (slot) {
        data.slot = slot;
      }
    }
    installComponentHooks(data);
    var name = getComponentName(Ctor.options) || tag;
    var vnode = new VNode(
      // @ts-expect-error
      "vue-component-".concat(Ctor.cid).concat(name ? "-".concat(name) : ""),
      data,
      void 0,
      void 0,
      void 0,
      context,
      // @ts-expect-error
      { Ctor, propsData, listeners, tag, children },
      asyncFactory
    );
    return vnode;
  }
  function createComponentInstanceForVnode(vnode, parent) {
    var options = {
      _isComponent: true,
      _parentVnode: vnode,
      parent
    };
    var inlineTemplate = vnode.data.inlineTemplate;
    if (isDef(inlineTemplate)) {
      options.render = inlineTemplate.render;
      options.staticRenderFns = inlineTemplate.staticRenderFns;
    }
    return new vnode.componentOptions.Ctor(options);
  }
  function installComponentHooks(data) {
    var hooks2 = data.hook || (data.hook = {});
    for (var i = 0; i < hooksToMerge.length; i++) {
      var key = hooksToMerge[i];
      var existing = hooks2[key];
      var toMerge = componentVNodeHooks[key];
      if (existing !== toMerge && !(existing && existing._merged)) {
        hooks2[key] = existing ? mergeHook(toMerge, existing) : toMerge;
      }
    }
  }
  function mergeHook(f1, f2) {
    var merged = function(a, b) {
      f1(a, b);
      f2(a, b);
    };
    merged._merged = true;
    return merged;
  }
  function transformModel(options, data) {
    var prop = options.model && options.model.prop || "value";
    var event = options.model && options.model.event || "input";
    (data.attrs || (data.attrs = {}))[prop] = data.model.value;
    var on2 = data.on || (data.on = {});
    var existing = on2[event];
    var callback = data.model.callback;
    if (isDef(existing)) {
      if (isArray(existing) ? existing.indexOf(callback) === -1 : existing !== callback) {
        on2[event] = [callback].concat(existing);
      }
    } else {
      on2[event] = callback;
    }
  }
  var warn$2 = noop;
  var tip = noop;
  var generateComponentTrace;
  var formatComponentName;
  if (true) {
    hasConsole_1 = typeof console !== "undefined";
    classifyRE_1 = /(?:^|[-_])(\w)/g;
    classify_1 = function(str2) {
      return str2.replace(classifyRE_1, function(c) {
        return c.toUpperCase();
      }).replace(/[-_]/g, "");
    };
    warn$2 = function(msg, vm3) {
      if (vm3 === void 0) {
        vm3 = currentInstance;
      }
      var trace = vm3 ? generateComponentTrace(vm3) : "";
      if (config.warnHandler) {
        config.warnHandler.call(null, msg, vm3, trace);
      } else if (hasConsole_1 && !config.silent) {
        console.error("[Vue warn]: ".concat(msg).concat(trace));
      }
    };
    tip = function(msg, vm3) {
      if (hasConsole_1 && !config.silent) {
        console.warn("[Vue tip]: ".concat(msg) + (vm3 ? generateComponentTrace(vm3) : ""));
      }
    };
    formatComponentName = function(vm3, includeFile) {
      if (vm3.$root === vm3) {
        return "<Root>";
      }
      var options = isFunction(vm3) && vm3.cid != null ? vm3.options : vm3._isVue ? vm3.$options || vm3.constructor.options : vm3;
      var name = getComponentName(options);
      var file = options.__file;
      if (!name && file) {
        var match2 = file.match(/([^/\\]+)\.vue$/);
        name = match2 && match2[1];
      }
      return (name ? "<".concat(classify_1(name), ">") : "<Anonymous>") + (file && includeFile !== false ? " at ".concat(file) : "");
    };
    repeat_1 = function(str2, n) {
      var res = "";
      while (n) {
        if (n % 2 === 1)
          res += str2;
        if (n > 1)
          str2 += str2;
        n >>= 1;
      }
      return res;
    };
    generateComponentTrace = function(vm3) {
      if (vm3._isVue && vm3.$parent) {
        var tree = [];
        var currentRecursiveSequence = 0;
        while (vm3) {
          if (tree.length > 0) {
            var last = tree[tree.length - 1];
            if (last.constructor === vm3.constructor) {
              currentRecursiveSequence++;
              vm3 = vm3.$parent;
              continue;
            } else if (currentRecursiveSequence > 0) {
              tree[tree.length - 1] = [last, currentRecursiveSequence];
              currentRecursiveSequence = 0;
            }
          }
          tree.push(vm3);
          vm3 = vm3.$parent;
        }
        return "\n\nfound in\n\n" + tree.map(function(vm4, i) {
          return "".concat(i === 0 ? "---> " : repeat_1(" ", 5 + i * 2)).concat(isArray(vm4) ? "".concat(formatComponentName(vm4[0]), "... (").concat(vm4[1], " recursive calls)") : formatComponentName(vm4));
        }).join("\n");
      } else {
        return "\n\n(found in ".concat(formatComponentName(vm3), ")");
      }
    };
  }
  var hasConsole_1;
  var classifyRE_1;
  var classify_1;
  var repeat_1;
  var strats = config.optionMergeStrategies;
  if (true) {
    strats.el = strats.propsData = function(parent, child, vm3, key) {
      if (!vm3) {
        warn$2('option "'.concat(key, '" can only be used during instance ') + "creation with the `new` keyword.");
      }
      return defaultStrat(parent, child);
    };
  }
  function mergeData(to, from, recursive) {
    if (recursive === void 0) {
      recursive = true;
    }
    if (!from)
      return to;
    var key, toVal, fromVal;
    var keys = hasSymbol ? Reflect.ownKeys(from) : Object.keys(from);
    for (var i = 0; i < keys.length; i++) {
      key = keys[i];
      if (key === "__ob__")
        continue;
      toVal = to[key];
      fromVal = from[key];
      if (!recursive || !hasOwn(to, key)) {
        set(to, key, fromVal);
      } else if (toVal !== fromVal && isPlainObject(toVal) && isPlainObject(fromVal)) {
        mergeData(toVal, fromVal);
      }
    }
    return to;
  }
  function mergeDataOrFn(parentVal, childVal, vm3) {
    if (!vm3) {
      if (!childVal) {
        return parentVal;
      }
      if (!parentVal) {
        return childVal;
      }
      return function mergedDataFn() {
        return mergeData(isFunction(childVal) ? childVal.call(this, this) : childVal, isFunction(parentVal) ? parentVal.call(this, this) : parentVal);
      };
    } else {
      return function mergedInstanceDataFn() {
        var instanceData = isFunction(childVal) ? childVal.call(vm3, vm3) : childVal;
        var defaultData = isFunction(parentVal) ? parentVal.call(vm3, vm3) : parentVal;
        if (instanceData) {
          return mergeData(instanceData, defaultData);
        } else {
          return defaultData;
        }
      };
    }
  }
  strats.data = function(parentVal, childVal, vm3) {
    if (!vm3) {
      if (childVal && typeof childVal !== "function") {
        warn$2('The "data" option should be a function that returns a per-instance value in component definitions.', vm3);
        return parentVal;
      }
      return mergeDataOrFn(parentVal, childVal);
    }
    return mergeDataOrFn(parentVal, childVal, vm3);
  };
  function mergeLifecycleHook(parentVal, childVal) {
    var res = childVal ? parentVal ? parentVal.concat(childVal) : isArray(childVal) ? childVal : [childVal] : parentVal;
    return res ? dedupeHooks(res) : res;
  }
  function dedupeHooks(hooks2) {
    var res = [];
    for (var i = 0; i < hooks2.length; i++) {
      if (res.indexOf(hooks2[i]) === -1) {
        res.push(hooks2[i]);
      }
    }
    return res;
  }
  LIFECYCLE_HOOKS.forEach(function(hook) {
    strats[hook] = mergeLifecycleHook;
  });
  function mergeAssets(parentVal, childVal, vm3, key) {
    var res = Object.create(parentVal || null);
    if (childVal) {
      assertObjectType(key, childVal, vm3);
      return extend(res, childVal);
    } else {
      return res;
    }
  }
  ASSET_TYPES.forEach(function(type) {
    strats[type + "s"] = mergeAssets;
  });
  strats.watch = function(parentVal, childVal, vm3, key) {
    if (parentVal === nativeWatch)
      parentVal = void 0;
    if (childVal === nativeWatch)
      childVal = void 0;
    if (!childVal)
      return Object.create(parentVal || null);
    if (true) {
      assertObjectType(key, childVal, vm3);
    }
    if (!parentVal)
      return childVal;
    var ret = {};
    extend(ret, parentVal);
    for (var key_1 in childVal) {
      var parent_1 = ret[key_1];
      var child = childVal[key_1];
      if (parent_1 && !isArray(parent_1)) {
        parent_1 = [parent_1];
      }
      ret[key_1] = parent_1 ? parent_1.concat(child) : isArray(child) ? child : [child];
    }
    return ret;
  };
  strats.props = strats.methods = strats.inject = strats.computed = function(parentVal, childVal, vm3, key) {
    if (childVal && true) {
      assertObjectType(key, childVal, vm3);
    }
    if (!parentVal)
      return childVal;
    var ret = /* @__PURE__ */ Object.create(null);
    extend(ret, parentVal);
    if (childVal)
      extend(ret, childVal);
    return ret;
  };
  strats.provide = function(parentVal, childVal) {
    if (!parentVal)
      return childVal;
    return function() {
      var ret = /* @__PURE__ */ Object.create(null);
      mergeData(ret, isFunction(parentVal) ? parentVal.call(this) : parentVal);
      if (childVal) {
        mergeData(
          ret,
          isFunction(childVal) ? childVal.call(this) : childVal,
          false
          // non-recursive
        );
      }
      return ret;
    };
  };
  var defaultStrat = function(parentVal, childVal) {
    return childVal === void 0 ? parentVal : childVal;
  };
  function checkComponents(options) {
    for (var key in options.components) {
      validateComponentName(key);
    }
  }
  function validateComponentName(name) {
    if (!new RegExp("^[a-zA-Z][\\-\\.0-9_".concat(unicodeRegExp.source, "]*$")).test(name)) {
      warn$2('Invalid component name: "' + name + '". Component names should conform to valid custom element name in html5 specification.');
    }
    if (isBuiltInTag(name) || config.isReservedTag(name)) {
      warn$2("Do not use built-in or reserved HTML elements as component id: " + name);
    }
  }
  function normalizeProps(options, vm3) {
    var props2 = options.props;
    if (!props2)
      return;
    var res = {};
    var i, val, name;
    if (isArray(props2)) {
      i = props2.length;
      while (i--) {
        val = props2[i];
        if (typeof val === "string") {
          name = camelize(val);
          res[name] = { type: null };
        } else if (true) {
          warn$2("props must be strings when using array syntax.");
        }
      }
    } else if (isPlainObject(props2)) {
      for (var key in props2) {
        val = props2[key];
        name = camelize(key);
        res[name] = isPlainObject(val) ? val : { type: val };
      }
    } else if (true) {
      warn$2('Invalid value for option "props": expected an Array or an Object, ' + "but got ".concat(toRawType(props2), "."), vm3);
    }
    options.props = res;
  }
  function normalizeInject(options, vm3) {
    var inject = options.inject;
    if (!inject)
      return;
    var normalized = options.inject = {};
    if (isArray(inject)) {
      for (var i = 0; i < inject.length; i++) {
        normalized[inject[i]] = { from: inject[i] };
      }
    } else if (isPlainObject(inject)) {
      for (var key in inject) {
        var val = inject[key];
        normalized[key] = isPlainObject(val) ? extend({ from: key }, val) : { from: val };
      }
    } else if (true) {
      warn$2('Invalid value for option "inject": expected an Array or an Object, ' + "but got ".concat(toRawType(inject), "."), vm3);
    }
  }
  function normalizeDirectives$1(options) {
    var dirs = options.directives;
    if (dirs) {
      for (var key in dirs) {
        var def2 = dirs[key];
        if (isFunction(def2)) {
          dirs[key] = { bind: def2, update: def2 };
        }
      }
    }
  }
  function assertObjectType(name, value, vm3) {
    if (!isPlainObject(value)) {
      warn$2('Invalid value for option "'.concat(name, '": expected an Object, ') + "but got ".concat(toRawType(value), "."), vm3);
    }
  }
  function mergeOptions(parent, child, vm3) {
    if (true) {
      checkComponents(child);
    }
    if (isFunction(child)) {
      child = child.options;
    }
    normalizeProps(child, vm3);
    normalizeInject(child, vm3);
    normalizeDirectives$1(child);
    if (!child._base) {
      if (child.extends) {
        parent = mergeOptions(parent, child.extends, vm3);
      }
      if (child.mixins) {
        for (var i = 0, l = child.mixins.length; i < l; i++) {
          parent = mergeOptions(parent, child.mixins[i], vm3);
        }
      }
    }
    var options = {};
    var key;
    for (key in parent) {
      mergeField(key);
    }
    for (key in child) {
      if (!hasOwn(parent, key)) {
        mergeField(key);
      }
    }
    function mergeField(key2) {
      var strat = strats[key2] || defaultStrat;
      options[key2] = strat(parent[key2], child[key2], vm3, key2);
    }
    return options;
  }
  function resolveAsset(options, type, id, warnMissing) {
    if (typeof id !== "string") {
      return;
    }
    var assets = options[type];
    if (hasOwn(assets, id))
      return assets[id];
    var camelizedId = camelize(id);
    if (hasOwn(assets, camelizedId))
      return assets[camelizedId];
    var PascalCaseId = capitalize(camelizedId);
    if (hasOwn(assets, PascalCaseId))
      return assets[PascalCaseId];
    var res = assets[id] || assets[camelizedId] || assets[PascalCaseId];
    if (warnMissing && !res) {
      warn$2("Failed to resolve " + type.slice(0, -1) + ": " + id);
    }
    return res;
  }
  function validateProp(key, propOptions, propsData, vm3) {
    var prop = propOptions[key];
    var absent = !hasOwn(propsData, key);
    var value = propsData[key];
    var booleanIndex = getTypeIndex(Boolean, prop.type);
    if (booleanIndex > -1) {
      if (absent && !hasOwn(prop, "default")) {
        value = false;
      } else if (value === "" || value === hyphenate(key)) {
        var stringIndex = getTypeIndex(String, prop.type);
        if (stringIndex < 0 || booleanIndex < stringIndex) {
          value = true;
        }
      }
    }
    if (value === void 0) {
      value = getPropDefaultValue(vm3, prop, key);
      var prevShouldObserve = shouldObserve;
      toggleObserving(true);
      observe(value);
      toggleObserving(prevShouldObserve);
    }
    if (true) {
      assertProp(prop, key, value, vm3, absent);
    }
    return value;
  }
  function getPropDefaultValue(vm3, prop, key) {
    if (!hasOwn(prop, "default")) {
      return void 0;
    }
    var def2 = prop.default;
    if (isObject(def2)) {
      warn$2('Invalid default value for prop "' + key + '": Props with type Object/Array must use a factory function to return the default value.', vm3);
    }
    if (vm3 && vm3.$options.propsData && vm3.$options.propsData[key] === void 0 && vm3._props[key] !== void 0) {
      return vm3._props[key];
    }
    return isFunction(def2) && getType(prop.type) !== "Function" ? def2.call(vm3) : def2;
  }
  function assertProp(prop, name, value, vm3, absent) {
    if (prop.required && absent) {
      warn$2('Missing required prop: "' + name + '"', vm3);
      return;
    }
    if (value == null && !prop.required) {
      return;
    }
    var type = prop.type;
    var valid = !type || type === true;
    var expectedTypes = [];
    if (type) {
      if (!isArray(type)) {
        type = [type];
      }
      for (var i = 0; i < type.length && !valid; i++) {
        var assertedType = assertType(value, type[i], vm3);
        expectedTypes.push(assertedType.expectedType || "");
        valid = assertedType.valid;
      }
    }
    var haveExpectedTypes = expectedTypes.some(function(t) {
      return t;
    });
    if (!valid && haveExpectedTypes) {
      warn$2(getInvalidTypeMessage(name, value, expectedTypes), vm3);
      return;
    }
    var validator = prop.validator;
    if (validator) {
      if (!validator(value)) {
        warn$2('Invalid prop: custom validator check failed for prop "' + name + '".', vm3);
      }
    }
  }
  var simpleCheckRE = /^(String|Number|Boolean|Function|Symbol|BigInt)$/;
  function assertType(value, type, vm3) {
    var valid;
    var expectedType = getType(type);
    if (simpleCheckRE.test(expectedType)) {
      var t = typeof value;
      valid = t === expectedType.toLowerCase();
      if (!valid && t === "object") {
        valid = value instanceof type;
      }
    } else if (expectedType === "Object") {
      valid = isPlainObject(value);
    } else if (expectedType === "Array") {
      valid = isArray(value);
    } else {
      try {
        valid = value instanceof type;
      } catch (e) {
        warn$2('Invalid prop type: "' + String(type) + '" is not a constructor', vm3);
        valid = false;
      }
    }
    return {
      valid,
      expectedType
    };
  }
  var functionTypeCheckRE = /^\s*function (\w+)/;
  function getType(fn) {
    var match2 = fn && fn.toString().match(functionTypeCheckRE);
    return match2 ? match2[1] : "";
  }
  function isSameType(a, b) {
    return getType(a) === getType(b);
  }
  function getTypeIndex(type, expectedTypes) {
    if (!isArray(expectedTypes)) {
      return isSameType(expectedTypes, type) ? 0 : -1;
    }
    for (var i = 0, len2 = expectedTypes.length; i < len2; i++) {
      if (isSameType(expectedTypes[i], type)) {
        return i;
      }
    }
    return -1;
  }
  function getInvalidTypeMessage(name, value, expectedTypes) {
    var message = 'Invalid prop: type check failed for prop "'.concat(name, '".') + " Expected ".concat(expectedTypes.map(capitalize).join(", "));
    var expectedType = expectedTypes[0];
    var receivedType = toRawType(value);
    if (expectedTypes.length === 1 && isExplicable(expectedType) && isExplicable(typeof value) && !isBoolean(expectedType, receivedType)) {
      message += " with value ".concat(styleValue(value, expectedType));
    }
    message += ", got ".concat(receivedType, " ");
    if (isExplicable(receivedType)) {
      message += "with value ".concat(styleValue(value, receivedType), ".");
    }
    return message;
  }
  function styleValue(value, type) {
    if (type === "String") {
      return '"'.concat(value, '"');
    } else if (type === "Number") {
      return "".concat(Number(value));
    } else {
      return "".concat(value);
    }
  }
  var EXPLICABLE_TYPES = ["string", "number", "boolean"];
  function isExplicable(value) {
    return EXPLICABLE_TYPES.some(function(elem) {
      return value.toLowerCase() === elem;
    });
  }
  function isBoolean() {
    var args = [];
    for (var _i = 0; _i < arguments.length; _i++) {
      args[_i] = arguments[_i];
    }
    return args.some(function(elem) {
      return elem.toLowerCase() === "boolean";
    });
  }
  function Vue(options) {
    if (!(this instanceof Vue)) {
      warn$2("Vue is a constructor and should be called with the `new` keyword");
    }
    this._init(options);
  }
  initMixin$1(Vue);
  stateMixin(Vue);
  eventsMixin(Vue);
  lifecycleMixin(Vue);
  renderMixin(Vue);
  function initUse(Vue2) {
    Vue2.use = function(plugin) {
      var installedPlugins = this._installedPlugins || (this._installedPlugins = []);
      if (installedPlugins.indexOf(plugin) > -1) {
        return this;
      }
      var args = toArray(arguments, 1);
      args.unshift(this);
      if (isFunction(plugin.install)) {
        plugin.install.apply(plugin, args);
      } else if (isFunction(plugin)) {
        plugin.apply(null, args);
      }
      installedPlugins.push(plugin);
      return this;
    };
  }
  function initMixin(Vue2) {
    Vue2.mixin = function(mixin) {
      this.options = mergeOptions(this.options, mixin);
      return this;
    };
  }
  function initExtend(Vue2) {
    Vue2.cid = 0;
    var cid = 1;
    Vue2.extend = function(extendOptions) {
      extendOptions = extendOptions || {};
      var Super = this;
      var SuperId = Super.cid;
      var cachedCtors = extendOptions._Ctor || (extendOptions._Ctor = {});
      if (cachedCtors[SuperId]) {
        return cachedCtors[SuperId];
      }
      var name = getComponentName(extendOptions) || getComponentName(Super.options);
      if (name) {
        validateComponentName(name);
      }
      var Sub = function VueComponent(options) {
        this._init(options);
      };
      Sub.prototype = Object.create(Super.prototype);
      Sub.prototype.constructor = Sub;
      Sub.cid = cid++;
      Sub.options = mergeOptions(Super.options, extendOptions);
      Sub["super"] = Super;
      if (Sub.options.props) {
        initProps(Sub);
      }
      if (Sub.options.computed) {
        initComputed(Sub);
      }
      Sub.extend = Super.extend;
      Sub.mixin = Super.mixin;
      Sub.use = Super.use;
      ASSET_TYPES.forEach(function(type) {
        Sub[type] = Super[type];
      });
      if (name) {
        Sub.options.components[name] = Sub;
      }
      Sub.superOptions = Super.options;
      Sub.extendOptions = extendOptions;
      Sub.sealedOptions = extend({}, Sub.options);
      cachedCtors[SuperId] = Sub;
      return Sub;
    };
  }
  function initProps(Comp) {
    var props2 = Comp.options.props;
    for (var key in props2) {
      proxy(Comp.prototype, "_props", key);
    }
  }
  function initComputed(Comp) {
    var computed = Comp.options.computed;
    for (var key in computed) {
      defineComputed(Comp.prototype, key, computed[key]);
    }
  }
  function initAssetRegisters(Vue2) {
    ASSET_TYPES.forEach(function(type) {
      Vue2[type] = function(id, definition) {
        if (!definition) {
          return this.options[type + "s"][id];
        } else {
          if (type === "component") {
            validateComponentName(id);
          }
          if (type === "component" && isPlainObject(definition)) {
            definition.name = definition.name || id;
            definition = this.options._base.extend(definition);
          }
          if (type === "directive" && isFunction(definition)) {
            definition = { bind: definition, update: definition };
          }
          this.options[type + "s"][id] = definition;
          return definition;
        }
      };
    });
  }
  function _getComponentName(opts2) {
    return opts2 && (getComponentName(opts2.Ctor.options) || opts2.tag);
  }
  function matches(pattern, name) {
    if (isArray(pattern)) {
      return pattern.indexOf(name) > -1;
    } else if (typeof pattern === "string") {
      return pattern.split(",").indexOf(name) > -1;
    } else if (isRegExp(pattern)) {
      return pattern.test(name);
    }
    return false;
  }
  function pruneCache(keepAliveInstance, filter) {
    var cache2 = keepAliveInstance.cache, keys = keepAliveInstance.keys, _vnode = keepAliveInstance._vnode, $vnode = keepAliveInstance.$vnode;
    for (var key in cache2) {
      var entry = cache2[key];
      if (entry) {
        var name_1 = entry.name;
        if (name_1 && !filter(name_1)) {
          pruneCacheEntry(cache2, key, keys, _vnode);
        }
      }
    }
    $vnode.componentOptions.children = void 0;
  }
  function pruneCacheEntry(cache2, key, keys, current) {
    var entry = cache2[key];
    if (entry && (!current || entry.tag !== current.tag)) {
      entry.componentInstance.$destroy();
    }
    cache2[key] = null;
    remove$2(keys, key);
  }
  var patternTypes = [String, RegExp, Array];
  var KeepAlive = {
    name: "keep-alive",
    abstract: true,
    props: {
      include: patternTypes,
      exclude: patternTypes,
      max: [String, Number]
    },
    methods: {
      cacheVNode: function() {
        var _a2 = this, cache2 = _a2.cache, keys = _a2.keys, vnodeToCache = _a2.vnodeToCache, keyToCache = _a2.keyToCache;
        if (vnodeToCache) {
          var tag = vnodeToCache.tag, componentInstance = vnodeToCache.componentInstance, componentOptions = vnodeToCache.componentOptions;
          cache2[keyToCache] = {
            name: _getComponentName(componentOptions),
            tag,
            componentInstance
          };
          keys.push(keyToCache);
          if (this.max && keys.length > parseInt(this.max)) {
            pruneCacheEntry(cache2, keys[0], keys, this._vnode);
          }
          this.vnodeToCache = null;
        }
      }
    },
    created: function() {
      this.cache = /* @__PURE__ */ Object.create(null);
      this.keys = [];
    },
    destroyed: function() {
      for (var key in this.cache) {
        pruneCacheEntry(this.cache, key, this.keys);
      }
    },
    mounted: function() {
      var _this = this;
      this.cacheVNode();
      this.$watch("include", function(val) {
        pruneCache(_this, function(name) {
          return matches(val, name);
        });
      });
      this.$watch("exclude", function(val) {
        pruneCache(_this, function(name) {
          return !matches(val, name);
        });
      });
    },
    updated: function() {
      this.cacheVNode();
    },
    render: function() {
      var slot = this.$slots.default;
      var vnode = getFirstComponentChild(slot);
      var componentOptions = vnode && vnode.componentOptions;
      if (componentOptions) {
        var name_2 = _getComponentName(componentOptions);
        var _a2 = this, include = _a2.include, exclude = _a2.exclude;
        if (
          // not included
          include && (!name_2 || !matches(include, name_2)) || // excluded
          exclude && name_2 && matches(exclude, name_2)
        ) {
          return vnode;
        }
        var _b = this, cache2 = _b.cache, keys = _b.keys;
        var key = vnode.key == null ? (
          // same constructor may get registered as different local components
          // so cid alone is not enough (#3269)
          componentOptions.Ctor.cid + (componentOptions.tag ? "::".concat(componentOptions.tag) : "")
        ) : vnode.key;
        if (cache2[key]) {
          vnode.componentInstance = cache2[key].componentInstance;
          remove$2(keys, key);
          keys.push(key);
        } else {
          this.vnodeToCache = vnode;
          this.keyToCache = key;
        }
        vnode.data.keepAlive = true;
      }
      return vnode || slot && slot[0];
    }
  };
  var builtInComponents = {
    KeepAlive
  };
  function initGlobalAPI(Vue2) {
    var configDef = {};
    configDef.get = function() {
      return config;
    };
    if (true) {
      configDef.set = function() {
        warn$2("Do not replace the Vue.config object, set individual fields instead.");
      };
    }
    Object.defineProperty(Vue2, "config", configDef);
    Vue2.util = {
      warn: warn$2,
      extend,
      mergeOptions,
      defineReactive
    };
    Vue2.set = set;
    Vue2.delete = del;
    Vue2.nextTick = nextTick;
    Vue2.observable = function(obj) {
      observe(obj);
      return obj;
    };
    Vue2.options = /* @__PURE__ */ Object.create(null);
    ASSET_TYPES.forEach(function(type) {
      Vue2.options[type + "s"] = /* @__PURE__ */ Object.create(null);
    });
    Vue2.options._base = Vue2;
    extend(Vue2.options.components, builtInComponents);
    initUse(Vue2);
    initMixin(Vue2);
    initExtend(Vue2);
    initAssetRegisters(Vue2);
  }
  initGlobalAPI(Vue);
  Object.defineProperty(Vue.prototype, "$isServer", {
    get: isServerRendering
  });
  Object.defineProperty(Vue.prototype, "$ssrContext", {
    get: function() {
      return this.$vnode && this.$vnode.ssrContext;
    }
  });
  Object.defineProperty(Vue, "FunctionalRenderContext", {
    value: FunctionalRenderContext
  });
  Vue.version = version;
  var isReservedAttr = makeMap("style,class");
  var acceptValue = makeMap("input,textarea,option,select,progress");
  var mustUseProp = function(tag, type, attr) {
    return attr === "value" && acceptValue(tag) && type !== "button" || attr === "selected" && tag === "option" || attr === "checked" && tag === "input" || attr === "muted" && tag === "video";
  };
  var isEnumeratedAttr = makeMap("contenteditable,draggable,spellcheck");
  var isValidContentEditableValue = makeMap("events,caret,typing,plaintext-only");
  var convertEnumeratedValue = function(key, value) {
    return isFalsyAttrValue(value) || value === "false" ? "false" : (
      // allow arbitrary string value for contenteditable
      key === "contenteditable" && isValidContentEditableValue(value) ? value : "true"
    );
  };
  var isBooleanAttr = makeMap("allowfullscreen,async,autofocus,autoplay,checked,compact,controls,declare,default,defaultchecked,defaultmuted,defaultselected,defer,disabled,enabled,formnovalidate,hidden,indeterminate,inert,ismap,itemscope,loop,multiple,muted,nohref,noresize,noshade,novalidate,nowrap,open,pauseonexit,readonly,required,reversed,scoped,seamless,selected,sortable,truespeed,typemustmatch,visible");
  var xlinkNS = "http://www.w3.org/1999/xlink";
  var isXlink = function(name) {
    return name.charAt(5) === ":" && name.slice(0, 5) === "xlink";
  };
  var getXlinkProp = function(name) {
    return isXlink(name) ? name.slice(6, name.length) : "";
  };
  var isFalsyAttrValue = function(val) {
    return val == null || val === false;
  };
  function genClassForVnode(vnode) {
    var data = vnode.data;
    var parentNode2 = vnode;
    var childNode = vnode;
    while (isDef(childNode.componentInstance)) {
      childNode = childNode.componentInstance._vnode;
      if (childNode && childNode.data) {
        data = mergeClassData(childNode.data, data);
      }
    }
    while (isDef(parentNode2 = parentNode2.parent)) {
      if (parentNode2 && parentNode2.data) {
        data = mergeClassData(data, parentNode2.data);
      }
    }
    return renderClass(data.staticClass, data.class);
  }
  function mergeClassData(child, parent) {
    return {
      staticClass: concat(child.staticClass, parent.staticClass),
      class: isDef(child.class) ? [child.class, parent.class] : parent.class
    };
  }
  function renderClass(staticClass, dynamicClass) {
    if (isDef(staticClass) || isDef(dynamicClass)) {
      return concat(staticClass, stringifyClass(dynamicClass));
    }
    return "";
  }
  function concat(a, b) {
    return a ? b ? a + " " + b : a : b || "";
  }
  function stringifyClass(value) {
    if (Array.isArray(value)) {
      return stringifyArray(value);
    }
    if (isObject(value)) {
      return stringifyObject(value);
    }
    if (typeof value === "string") {
      return value;
    }
    return "";
  }
  function stringifyArray(value) {
    var res = "";
    var stringified;
    for (var i = 0, l = value.length; i < l; i++) {
      if (isDef(stringified = stringifyClass(value[i])) && stringified !== "") {
        if (res)
          res += " ";
        res += stringified;
      }
    }
    return res;
  }
  function stringifyObject(value) {
    var res = "";
    for (var key in value) {
      if (value[key]) {
        if (res)
          res += " ";
        res += key;
      }
    }
    return res;
  }
  var namespaceMap = {
    svg: "http://www.w3.org/2000/svg",
    math: "http://www.w3.org/1998/Math/MathML"
  };
  var isHTMLTag = makeMap("html,body,base,head,link,meta,style,title,address,article,aside,footer,header,h1,h2,h3,h4,h5,h6,hgroup,nav,section,div,dd,dl,dt,figcaption,figure,picture,hr,img,li,main,ol,p,pre,ul,a,b,abbr,bdi,bdo,br,cite,code,data,dfn,em,i,kbd,mark,q,rp,rt,rtc,ruby,s,samp,small,span,strong,sub,sup,time,u,var,wbr,area,audio,map,track,video,embed,object,param,source,canvas,script,noscript,del,ins,caption,col,colgroup,table,thead,tbody,td,th,tr,button,datalist,fieldset,form,input,label,legend,meter,optgroup,option,output,progress,select,textarea,details,dialog,menu,menuitem,summary,content,element,shadow,template,blockquote,iframe,tfoot");
  var isSVG = makeMap("svg,animate,circle,clippath,cursor,defs,desc,ellipse,filter,font-face,foreignobject,g,glyph,image,line,marker,mask,missing-glyph,path,pattern,polygon,polyline,rect,switch,symbol,text,textpath,tspan,use,view", true);
  var isPreTag = function(tag) {
    return tag === "pre";
  };
  var isReservedTag = function(tag) {
    return isHTMLTag(tag) || isSVG(tag);
  };
  function getTagNamespace(tag) {
    if (isSVG(tag)) {
      return "svg";
    }
    if (tag === "math") {
      return "math";
    }
  }
  var unknownElementCache = /* @__PURE__ */ Object.create(null);
  function isUnknownElement(tag) {
    if (!inBrowser) {
      return true;
    }
    if (isReservedTag(tag)) {
      return false;
    }
    tag = tag.toLowerCase();
    if (unknownElementCache[tag] != null) {
      return unknownElementCache[tag];
    }
    var el = document.createElement(tag);
    if (tag.indexOf("-") > -1) {
      return unknownElementCache[tag] = el.constructor === window.HTMLUnknownElement || el.constructor === window.HTMLElement;
    } else {
      return unknownElementCache[tag] = /HTMLUnknownElement/.test(el.toString());
    }
  }
  var isTextInputType = makeMap("text,number,password,search,email,tel,url");
  function query(el) {
    if (typeof el === "string") {
      var selected = document.querySelector(el);
      if (!selected) {
        warn$2("Cannot find element: " + el);
        return document.createElement("div");
      }
      return selected;
    } else {
      return el;
    }
  }
  function createElement(tagName2, vnode) {
    var elm = document.createElement(tagName2);
    if (tagName2 !== "select") {
      return elm;
    }
    if (vnode.data && vnode.data.attrs && vnode.data.attrs.multiple !== void 0) {
      elm.setAttribute("multiple", "multiple");
    }
    return elm;
  }
  function createElementNS(namespace, tagName2) {
    return document.createElementNS(namespaceMap[namespace], tagName2);
  }
  function createTextNode(text2) {
    return document.createTextNode(text2);
  }
  function createComment(text2) {
    return document.createComment(text2);
  }
  function insertBefore(parentNode2, newNode, referenceNode) {
    parentNode2.insertBefore(newNode, referenceNode);
  }
  function removeChild(node, child) {
    node.removeChild(child);
  }
  function appendChild(node, child) {
    node.appendChild(child);
  }
  function parentNode(node) {
    return node.parentNode;
  }
  function nextSibling(node) {
    return node.nextSibling;
  }
  function tagName(node) {
    return node.tagName;
  }
  function setTextContent(node, text2) {
    node.textContent = text2;
  }
  function setStyleScope(node, scopeId) {
    node.setAttribute(scopeId, "");
  }
  var nodeOps = /* @__PURE__ */ Object.freeze({
    __proto__: null,
    createElement,
    createElementNS,
    createTextNode,
    createComment,
    insertBefore,
    removeChild,
    appendChild,
    parentNode,
    nextSibling,
    tagName,
    setTextContent,
    setStyleScope
  });
  var ref = {
    create: function(_, vnode) {
      registerRef(vnode);
    },
    update: function(oldVnode, vnode) {
      if (oldVnode.data.ref !== vnode.data.ref) {
        registerRef(oldVnode, true);
        registerRef(vnode);
      }
    },
    destroy: function(vnode) {
      registerRef(vnode, true);
    }
  };
  function registerRef(vnode, isRemoval) {
    var ref2 = vnode.data.ref;
    if (!isDef(ref2))
      return;
    var vm3 = vnode.context;
    var refValue = vnode.componentInstance || vnode.elm;
    var value = isRemoval ? null : refValue;
    var $refsValue = isRemoval ? void 0 : refValue;
    if (isFunction(ref2)) {
      invokeWithErrorHandling(ref2, vm3, [value], vm3, "template ref function");
      return;
    }
    var isFor = vnode.data.refInFor;
    var _isString = typeof ref2 === "string" || typeof ref2 === "number";
    var _isRef = isRef(ref2);
    var refs = vm3.$refs;
    if (_isString || _isRef) {
      if (isFor) {
        var existing = _isString ? refs[ref2] : ref2.value;
        if (isRemoval) {
          isArray(existing) && remove$2(existing, refValue);
        } else {
          if (!isArray(existing)) {
            if (_isString) {
              refs[ref2] = [refValue];
              setSetupRef(vm3, ref2, refs[ref2]);
            } else {
              ref2.value = [refValue];
            }
          } else if (!existing.includes(refValue)) {
            existing.push(refValue);
          }
        }
      } else if (_isString) {
        if (isRemoval && refs[ref2] !== refValue) {
          return;
        }
        refs[ref2] = $refsValue;
        setSetupRef(vm3, ref2, value);
      } else if (_isRef) {
        if (isRemoval && ref2.value !== refValue) {
          return;
        }
        ref2.value = value;
      } else if (true) {
        warn$2("Invalid template ref type: ".concat(typeof ref2));
      }
    }
  }
  function setSetupRef(_a2, key, val) {
    var _setupState = _a2._setupState;
    if (_setupState && hasOwn(_setupState, key)) {
      if (isRef(_setupState[key])) {
        _setupState[key].value = val;
      } else {
        _setupState[key] = val;
      }
    }
  }
  var emptyNode = new VNode("", {}, []);
  var hooks = ["create", "activate", "update", "remove", "destroy"];
  function sameVnode(a, b) {
    return a.key === b.key && a.asyncFactory === b.asyncFactory && (a.tag === b.tag && a.isComment === b.isComment && isDef(a.data) === isDef(b.data) && sameInputType(a, b) || isTrue(a.isAsyncPlaceholder) && isUndef(b.asyncFactory.error));
  }
  function sameInputType(a, b) {
    if (a.tag !== "input")
      return true;
    var i;
    var typeA = isDef(i = a.data) && isDef(i = i.attrs) && i.type;
    var typeB = isDef(i = b.data) && isDef(i = i.attrs) && i.type;
    return typeA === typeB || isTextInputType(typeA) && isTextInputType(typeB);
  }
  function createKeyToOldIdx(children, beginIdx, endIdx) {
    var i, key;
    var map = {};
    for (i = beginIdx; i <= endIdx; ++i) {
      key = children[i].key;
      if (isDef(key))
        map[key] = i;
    }
    return map;
  }
  function createPatchFunction(backend) {
    var i, j;
    var cbs = {};
    var modules2 = backend.modules, nodeOps2 = backend.nodeOps;
    for (i = 0; i < hooks.length; ++i) {
      cbs[hooks[i]] = [];
      for (j = 0; j < modules2.length; ++j) {
        if (isDef(modules2[j][hooks[i]])) {
          cbs[hooks[i]].push(modules2[j][hooks[i]]);
        }
      }
    }
    function emptyNodeAt(elm) {
      return new VNode(nodeOps2.tagName(elm).toLowerCase(), {}, [], void 0, elm);
    }
    function createRmCb(childElm, listeners) {
      function remove2() {
        if (--remove2.listeners === 0) {
          removeNode(childElm);
        }
      }
      remove2.listeners = listeners;
      return remove2;
    }
    function removeNode(el) {
      var parent = nodeOps2.parentNode(el);
      if (isDef(parent)) {
        nodeOps2.removeChild(parent, el);
      }
    }
    function isUnknownElement2(vnode, inVPre) {
      return !inVPre && !vnode.ns && !(config.ignoredElements.length && config.ignoredElements.some(function(ignore) {
        return isRegExp(ignore) ? ignore.test(vnode.tag) : ignore === vnode.tag;
      })) && config.isUnknownElement(vnode.tag);
    }
    var creatingElmInVPre = 0;
    function createElm(vnode, insertedVnodeQueue, parentElm, refElm, nested, ownerArray, index2) {
      if (isDef(vnode.elm) && isDef(ownerArray)) {
        vnode = ownerArray[index2] = cloneVNode(vnode);
      }
      vnode.isRootInsert = !nested;
      if (createComponent2(vnode, insertedVnodeQueue, parentElm, refElm)) {
        return;
      }
      var data = vnode.data;
      var children = vnode.children;
      var tag = vnode.tag;
      if (isDef(tag)) {
        if (true) {
          if (data && data.pre) {
            creatingElmInVPre++;
          }
          if (isUnknownElement2(vnode, creatingElmInVPre)) {
            warn$2("Unknown custom element: <" + tag + '> - did you register the component correctly? For recursive components, make sure to provide the "name" option.', vnode.context);
          }
        }
        vnode.elm = vnode.ns ? nodeOps2.createElementNS(vnode.ns, tag) : nodeOps2.createElement(tag, vnode);
        setScope(vnode);
        createChildren(vnode, children, insertedVnodeQueue);
        if (isDef(data)) {
          invokeCreateHooks(vnode, insertedVnodeQueue);
        }
        insert(parentElm, vnode.elm, refElm);
        if (data && data.pre) {
          creatingElmInVPre--;
        }
      } else if (isTrue(vnode.isComment)) {
        vnode.elm = nodeOps2.createComment(vnode.text);
        insert(parentElm, vnode.elm, refElm);
      } else {
        vnode.elm = nodeOps2.createTextNode(vnode.text);
        insert(parentElm, vnode.elm, refElm);
      }
    }
    function createComponent2(vnode, insertedVnodeQueue, parentElm, refElm) {
      var i2 = vnode.data;
      if (isDef(i2)) {
        var isReactivated = isDef(vnode.componentInstance) && i2.keepAlive;
        if (isDef(i2 = i2.hook) && isDef(i2 = i2.init)) {
          i2(
            vnode,
            false
            /* hydrating */
          );
        }
        if (isDef(vnode.componentInstance)) {
          initComponent(vnode, insertedVnodeQueue);
          insert(parentElm, vnode.elm, refElm);
          if (isTrue(isReactivated)) {
            reactivateComponent(vnode, insertedVnodeQueue, parentElm, refElm);
          }
          return true;
        }
      }
    }
    function initComponent(vnode, insertedVnodeQueue) {
      if (isDef(vnode.data.pendingInsert)) {
        insertedVnodeQueue.push.apply(insertedVnodeQueue, vnode.data.pendingInsert);
        vnode.data.pendingInsert = null;
      }
      vnode.elm = vnode.componentInstance.$el;
      if (isPatchable(vnode)) {
        invokeCreateHooks(vnode, insertedVnodeQueue);
        setScope(vnode);
      } else {
        registerRef(vnode);
        insertedVnodeQueue.push(vnode);
      }
    }
    function reactivateComponent(vnode, insertedVnodeQueue, parentElm, refElm) {
      var i2;
      var innerNode = vnode;
      while (innerNode.componentInstance) {
        innerNode = innerNode.componentInstance._vnode;
        if (isDef(i2 = innerNode.data) && isDef(i2 = i2.transition)) {
          for (i2 = 0; i2 < cbs.activate.length; ++i2) {
            cbs.activate[i2](emptyNode, innerNode);
          }
          insertedVnodeQueue.push(innerNode);
          break;
        }
      }
      insert(parentElm, vnode.elm, refElm);
    }
    function insert(parent, elm, ref2) {
      if (isDef(parent)) {
        if (isDef(ref2)) {
          if (nodeOps2.parentNode(ref2) === parent) {
            nodeOps2.insertBefore(parent, elm, ref2);
          }
        } else {
          nodeOps2.appendChild(parent, elm);
        }
      }
    }
    function createChildren(vnode, children, insertedVnodeQueue) {
      if (isArray(children)) {
        if (true) {
          checkDuplicateKeys(children);
        }
        for (var i_1 = 0; i_1 < children.length; ++i_1) {
          createElm(children[i_1], insertedVnodeQueue, vnode.elm, null, true, children, i_1);
        }
      } else if (isPrimitive(vnode.text)) {
        nodeOps2.appendChild(vnode.elm, nodeOps2.createTextNode(String(vnode.text)));
      }
    }
    function isPatchable(vnode) {
      while (vnode.componentInstance) {
        vnode = vnode.componentInstance._vnode;
      }
      return isDef(vnode.tag);
    }
    function invokeCreateHooks(vnode, insertedVnodeQueue) {
      for (var i_2 = 0; i_2 < cbs.create.length; ++i_2) {
        cbs.create[i_2](emptyNode, vnode);
      }
      i = vnode.data.hook;
      if (isDef(i)) {
        if (isDef(i.create))
          i.create(emptyNode, vnode);
        if (isDef(i.insert))
          insertedVnodeQueue.push(vnode);
      }
    }
    function setScope(vnode) {
      var i2;
      if (isDef(i2 = vnode.fnScopeId)) {
        nodeOps2.setStyleScope(vnode.elm, i2);
      } else {
        var ancestor = vnode;
        while (ancestor) {
          if (isDef(i2 = ancestor.context) && isDef(i2 = i2.$options._scopeId)) {
            nodeOps2.setStyleScope(vnode.elm, i2);
          }
          ancestor = ancestor.parent;
        }
      }
      if (isDef(i2 = activeInstance) && i2 !== vnode.context && i2 !== vnode.fnContext && isDef(i2 = i2.$options._scopeId)) {
        nodeOps2.setStyleScope(vnode.elm, i2);
      }
    }
    function addVnodes(parentElm, refElm, vnodes, startIdx, endIdx, insertedVnodeQueue) {
      for (; startIdx <= endIdx; ++startIdx) {
        createElm(vnodes[startIdx], insertedVnodeQueue, parentElm, refElm, false, vnodes, startIdx);
      }
    }
    function invokeDestroyHook(vnode) {
      var i2, j2;
      var data = vnode.data;
      if (isDef(data)) {
        if (isDef(i2 = data.hook) && isDef(i2 = i2.destroy))
          i2(vnode);
        for (i2 = 0; i2 < cbs.destroy.length; ++i2)
          cbs.destroy[i2](vnode);
      }
      if (isDef(i2 = vnode.children)) {
        for (j2 = 0; j2 < vnode.children.length; ++j2) {
          invokeDestroyHook(vnode.children[j2]);
        }
      }
    }
    function removeVnodes(vnodes, startIdx, endIdx) {
      for (; startIdx <= endIdx; ++startIdx) {
        var ch = vnodes[startIdx];
        if (isDef(ch)) {
          if (isDef(ch.tag)) {
            removeAndInvokeRemoveHook(ch);
            invokeDestroyHook(ch);
          } else {
            removeNode(ch.elm);
          }
        }
      }
    }
    function removeAndInvokeRemoveHook(vnode, rm) {
      if (isDef(rm) || isDef(vnode.data)) {
        var i_3;
        var listeners = cbs.remove.length + 1;
        if (isDef(rm)) {
          rm.listeners += listeners;
        } else {
          rm = createRmCb(vnode.elm, listeners);
        }
        if (isDef(i_3 = vnode.componentInstance) && isDef(i_3 = i_3._vnode) && isDef(i_3.data)) {
          removeAndInvokeRemoveHook(i_3, rm);
        }
        for (i_3 = 0; i_3 < cbs.remove.length; ++i_3) {
          cbs.remove[i_3](vnode, rm);
        }
        if (isDef(i_3 = vnode.data.hook) && isDef(i_3 = i_3.remove)) {
          i_3(vnode, rm);
        } else {
          rm();
        }
      } else {
        removeNode(vnode.elm);
      }
    }
    function updateChildren(parentElm, oldCh, newCh, insertedVnodeQueue, removeOnly) {
      var oldStartIdx = 0;
      var newStartIdx = 0;
      var oldEndIdx = oldCh.length - 1;
      var oldStartVnode = oldCh[0];
      var oldEndVnode = oldCh[oldEndIdx];
      var newEndIdx = newCh.length - 1;
      var newStartVnode = newCh[0];
      var newEndVnode = newCh[newEndIdx];
      var oldKeyToIdx, idxInOld, vnodeToMove, refElm;
      var canMove = !removeOnly;
      if (true) {
        checkDuplicateKeys(newCh);
      }
      while (oldStartIdx <= oldEndIdx && newStartIdx <= newEndIdx) {
        if (isUndef(oldStartVnode)) {
          oldStartVnode = oldCh[++oldStartIdx];
        } else if (isUndef(oldEndVnode)) {
          oldEndVnode = oldCh[--oldEndIdx];
        } else if (sameVnode(oldStartVnode, newStartVnode)) {
          patchVnode(oldStartVnode, newStartVnode, insertedVnodeQueue, newCh, newStartIdx);
          oldStartVnode = oldCh[++oldStartIdx];
          newStartVnode = newCh[++newStartIdx];
        } else if (sameVnode(oldEndVnode, newEndVnode)) {
          patchVnode(oldEndVnode, newEndVnode, insertedVnodeQueue, newCh, newEndIdx);
          oldEndVnode = oldCh[--oldEndIdx];
          newEndVnode = newCh[--newEndIdx];
        } else if (sameVnode(oldStartVnode, newEndVnode)) {
          patchVnode(oldStartVnode, newEndVnode, insertedVnodeQueue, newCh, newEndIdx);
          canMove && nodeOps2.insertBefore(parentElm, oldStartVnode.elm, nodeOps2.nextSibling(oldEndVnode.elm));
          oldStartVnode = oldCh[++oldStartIdx];
          newEndVnode = newCh[--newEndIdx];
        } else if (sameVnode(oldEndVnode, newStartVnode)) {
          patchVnode(oldEndVnode, newStartVnode, insertedVnodeQueue, newCh, newStartIdx);
          canMove && nodeOps2.insertBefore(parentElm, oldEndVnode.elm, oldStartVnode.elm);
          oldEndVnode = oldCh[--oldEndIdx];
          newStartVnode = newCh[++newStartIdx];
        } else {
          if (isUndef(oldKeyToIdx))
            oldKeyToIdx = createKeyToOldIdx(oldCh, oldStartIdx, oldEndIdx);
          idxInOld = isDef(newStartVnode.key) ? oldKeyToIdx[newStartVnode.key] : findIdxInOld(newStartVnode, oldCh, oldStartIdx, oldEndIdx);
          if (isUndef(idxInOld)) {
            createElm(newStartVnode, insertedVnodeQueue, parentElm, oldStartVnode.elm, false, newCh, newStartIdx);
          } else {
            vnodeToMove = oldCh[idxInOld];
            if (sameVnode(vnodeToMove, newStartVnode)) {
              patchVnode(vnodeToMove, newStartVnode, insertedVnodeQueue, newCh, newStartIdx);
              oldCh[idxInOld] = void 0;
              canMove && nodeOps2.insertBefore(parentElm, vnodeToMove.elm, oldStartVnode.elm);
            } else {
              createElm(newStartVnode, insertedVnodeQueue, parentElm, oldStartVnode.elm, false, newCh, newStartIdx);
            }
          }
          newStartVnode = newCh[++newStartIdx];
        }
      }
      if (oldStartIdx > oldEndIdx) {
        refElm = isUndef(newCh[newEndIdx + 1]) ? null : newCh[newEndIdx + 1].elm;
        addVnodes(parentElm, refElm, newCh, newStartIdx, newEndIdx, insertedVnodeQueue);
      } else if (newStartIdx > newEndIdx) {
        removeVnodes(oldCh, oldStartIdx, oldEndIdx);
      }
    }
    function checkDuplicateKeys(children) {
      var seenKeys = {};
      for (var i_4 = 0; i_4 < children.length; i_4++) {
        var vnode = children[i_4];
        var key = vnode.key;
        if (isDef(key)) {
          if (seenKeys[key]) {
            warn$2("Duplicate keys detected: '".concat(key, "'. This may cause an update error."), vnode.context);
          } else {
            seenKeys[key] = true;
          }
        }
      }
    }
    function findIdxInOld(node, oldCh, start, end) {
      for (var i_5 = start; i_5 < end; i_5++) {
        var c = oldCh[i_5];
        if (isDef(c) && sameVnode(node, c))
          return i_5;
      }
    }
    function patchVnode(oldVnode, vnode, insertedVnodeQueue, ownerArray, index2, removeOnly) {
      if (oldVnode === vnode) {
        return;
      }
      if (isDef(vnode.elm) && isDef(ownerArray)) {
        vnode = ownerArray[index2] = cloneVNode(vnode);
      }
      var elm = vnode.elm = oldVnode.elm;
      if (isTrue(oldVnode.isAsyncPlaceholder)) {
        if (isDef(vnode.asyncFactory.resolved)) {
          hydrate(oldVnode.elm, vnode, insertedVnodeQueue);
        } else {
          vnode.isAsyncPlaceholder = true;
        }
        return;
      }
      if (isTrue(vnode.isStatic) && isTrue(oldVnode.isStatic) && vnode.key === oldVnode.key && (isTrue(vnode.isCloned) || isTrue(vnode.isOnce))) {
        vnode.componentInstance = oldVnode.componentInstance;
        return;
      }
      var i2;
      var data = vnode.data;
      if (isDef(data) && isDef(i2 = data.hook) && isDef(i2 = i2.prepatch)) {
        i2(oldVnode, vnode);
      }
      var oldCh = oldVnode.children;
      var ch = vnode.children;
      if (isDef(data) && isPatchable(vnode)) {
        for (i2 = 0; i2 < cbs.update.length; ++i2)
          cbs.update[i2](oldVnode, vnode);
        if (isDef(i2 = data.hook) && isDef(i2 = i2.update))
          i2(oldVnode, vnode);
      }
      if (isUndef(vnode.text)) {
        if (isDef(oldCh) && isDef(ch)) {
          if (oldCh !== ch)
            updateChildren(elm, oldCh, ch, insertedVnodeQueue, removeOnly);
        } else if (isDef(ch)) {
          if (true) {
            checkDuplicateKeys(ch);
          }
          if (isDef(oldVnode.text))
            nodeOps2.setTextContent(elm, "");
          addVnodes(elm, null, ch, 0, ch.length - 1, insertedVnodeQueue);
        } else if (isDef(oldCh)) {
          removeVnodes(oldCh, 0, oldCh.length - 1);
        } else if (isDef(oldVnode.text)) {
          nodeOps2.setTextContent(elm, "");
        }
      } else if (oldVnode.text !== vnode.text) {
        nodeOps2.setTextContent(elm, vnode.text);
      }
      if (isDef(data)) {
        if (isDef(i2 = data.hook) && isDef(i2 = i2.postpatch))
          i2(oldVnode, vnode);
      }
    }
    function invokeInsertHook(vnode, queue2, initial) {
      if (isTrue(initial) && isDef(vnode.parent)) {
        vnode.parent.data.pendingInsert = queue2;
      } else {
        for (var i_6 = 0; i_6 < queue2.length; ++i_6) {
          queue2[i_6].data.hook.insert(queue2[i_6]);
        }
      }
    }
    var hydrationBailed = false;
    var isRenderedModule = makeMap("attrs,class,staticClass,staticStyle,key");
    function hydrate(elm, vnode, insertedVnodeQueue, inVPre) {
      var i2;
      var tag = vnode.tag, data = vnode.data, children = vnode.children;
      inVPre = inVPre || data && data.pre;
      vnode.elm = elm;
      if (isTrue(vnode.isComment) && isDef(vnode.asyncFactory)) {
        vnode.isAsyncPlaceholder = true;
        return true;
      }
      if (true) {
        if (!assertNodeMatch(elm, vnode, inVPre)) {
          return false;
        }
      }
      if (isDef(data)) {
        if (isDef(i2 = data.hook) && isDef(i2 = i2.init))
          i2(
            vnode,
            true
            /* hydrating */
          );
        if (isDef(i2 = vnode.componentInstance)) {
          initComponent(vnode, insertedVnodeQueue);
          return true;
        }
      }
      if (isDef(tag)) {
        if (isDef(children)) {
          if (!elm.hasChildNodes()) {
            createChildren(vnode, children, insertedVnodeQueue);
          } else {
            if (isDef(i2 = data) && isDef(i2 = i2.domProps) && isDef(i2 = i2.innerHTML)) {
              if (i2 !== elm.innerHTML) {
                if (typeof console !== "undefined" && !hydrationBailed) {
                  hydrationBailed = true;
                  console.warn("Parent: ", elm);
                  console.warn("server innerHTML: ", i2);
                  console.warn("client innerHTML: ", elm.innerHTML);
                }
                return false;
              }
            } else {
              var childrenMatch = true;
              var childNode = elm.firstChild;
              for (var i_7 = 0; i_7 < children.length; i_7++) {
                if (!childNode || !hydrate(childNode, children[i_7], insertedVnodeQueue, inVPre)) {
                  childrenMatch = false;
                  break;
                }
                childNode = childNode.nextSibling;
              }
              if (!childrenMatch || childNode) {
                if (typeof console !== "undefined" && !hydrationBailed) {
                  hydrationBailed = true;
                  console.warn("Parent: ", elm);
                  console.warn("Mismatching childNodes vs. VNodes: ", elm.childNodes, children);
                }
                return false;
              }
            }
          }
        }
        if (isDef(data)) {
          var fullInvoke = false;
          for (var key in data) {
            if (!isRenderedModule(key)) {
              fullInvoke = true;
              invokeCreateHooks(vnode, insertedVnodeQueue);
              break;
            }
          }
          if (!fullInvoke && data["class"]) {
            traverse(data["class"]);
          }
        }
      } else if (elm.data !== vnode.text) {
        elm.data = vnode.text;
      }
      return true;
    }
    function assertNodeMatch(node, vnode, inVPre) {
      if (isDef(vnode.tag)) {
        return vnode.tag.indexOf("vue-component") === 0 || !isUnknownElement2(vnode, inVPre) && vnode.tag.toLowerCase() === (node.tagName && node.tagName.toLowerCase());
      } else {
        return node.nodeType === (vnode.isComment ? 8 : 3);
      }
    }
    return function patch2(oldVnode, vnode, hydrating, removeOnly) {
      if (isUndef(vnode)) {
        if (isDef(oldVnode))
          invokeDestroyHook(oldVnode);
        return;
      }
      var isInitialPatch = false;
      var insertedVnodeQueue = [];
      if (isUndef(oldVnode)) {
        isInitialPatch = true;
        createElm(vnode, insertedVnodeQueue);
      } else {
        var isRealElement = isDef(oldVnode.nodeType);
        if (!isRealElement && sameVnode(oldVnode, vnode)) {
          patchVnode(oldVnode, vnode, insertedVnodeQueue, null, null, removeOnly);
        } else {
          if (isRealElement) {
            if (oldVnode.nodeType === 1 && oldVnode.hasAttribute(SSR_ATTR)) {
              oldVnode.removeAttribute(SSR_ATTR);
              hydrating = true;
            }
            if (isTrue(hydrating)) {
              if (hydrate(oldVnode, vnode, insertedVnodeQueue)) {
                invokeInsertHook(vnode, insertedVnodeQueue, true);
                return oldVnode;
              } else if (true) {
                warn$2("The client-side rendered virtual DOM tree is not matching server-rendered content. This is likely caused by incorrect HTML markup, for example nesting block-level elements inside <p>, or missing <tbody>. Bailing hydration and performing full client-side render.");
              }
            }
            oldVnode = emptyNodeAt(oldVnode);
          }
          var oldElm = oldVnode.elm;
          var parentElm = nodeOps2.parentNode(oldElm);
          createElm(
            vnode,
            insertedVnodeQueue,
            // extremely rare edge case: do not insert if old element is in a
            // leaving transition. Only happens when combining transition +
            // keep-alive + HOCs. (#4590)
            oldElm._leaveCb ? null : parentElm,
            nodeOps2.nextSibling(oldElm)
          );
          if (isDef(vnode.parent)) {
            var ancestor = vnode.parent;
            var patchable = isPatchable(vnode);
            while (ancestor) {
              for (var i_8 = 0; i_8 < cbs.destroy.length; ++i_8) {
                cbs.destroy[i_8](ancestor);
              }
              ancestor.elm = vnode.elm;
              if (patchable) {
                for (var i_9 = 0; i_9 < cbs.create.length; ++i_9) {
                  cbs.create[i_9](emptyNode, ancestor);
                }
                var insert_1 = ancestor.data.hook.insert;
                if (insert_1.merged) {
                  var cloned = insert_1.fns.slice(1);
                  for (var i_10 = 0; i_10 < cloned.length; i_10++) {
                    cloned[i_10]();
                  }
                }
              } else {
                registerRef(ancestor);
              }
              ancestor = ancestor.parent;
            }
          }
          if (isDef(parentElm)) {
            removeVnodes([oldVnode], 0, 0);
          } else if (isDef(oldVnode.tag)) {
            invokeDestroyHook(oldVnode);
          }
        }
      }
      invokeInsertHook(vnode, insertedVnodeQueue, isInitialPatch);
      return vnode.elm;
    };
  }
  var directives$1 = {
    create: updateDirectives,
    update: updateDirectives,
    destroy: function unbindDirectives(vnode) {
      updateDirectives(vnode, emptyNode);
    }
  };
  function updateDirectives(oldVnode, vnode) {
    if (oldVnode.data.directives || vnode.data.directives) {
      _update(oldVnode, vnode);
    }
  }
  function _update(oldVnode, vnode) {
    var isCreate = oldVnode === emptyNode;
    var isDestroy = vnode === emptyNode;
    var oldDirs = normalizeDirectives(oldVnode.data.directives, oldVnode.context);
    var newDirs = normalizeDirectives(vnode.data.directives, vnode.context);
    var dirsWithInsert = [];
    var dirsWithPostpatch = [];
    var key, oldDir, dir;
    for (key in newDirs) {
      oldDir = oldDirs[key];
      dir = newDirs[key];
      if (!oldDir) {
        callHook(dir, "bind", vnode, oldVnode);
        if (dir.def && dir.def.inserted) {
          dirsWithInsert.push(dir);
        }
      } else {
        dir.oldValue = oldDir.value;
        dir.oldArg = oldDir.arg;
        callHook(dir, "update", vnode, oldVnode);
        if (dir.def && dir.def.componentUpdated) {
          dirsWithPostpatch.push(dir);
        }
      }
    }
    if (dirsWithInsert.length) {
      var callInsert = function() {
        for (var i = 0; i < dirsWithInsert.length; i++) {
          callHook(dirsWithInsert[i], "inserted", vnode, oldVnode);
        }
      };
      if (isCreate) {
        mergeVNodeHook(vnode, "insert", callInsert);
      } else {
        callInsert();
      }
    }
    if (dirsWithPostpatch.length) {
      mergeVNodeHook(vnode, "postpatch", function() {
        for (var i = 0; i < dirsWithPostpatch.length; i++) {
          callHook(dirsWithPostpatch[i], "componentUpdated", vnode, oldVnode);
        }
      });
    }
    if (!isCreate) {
      for (key in oldDirs) {
        if (!newDirs[key]) {
          callHook(oldDirs[key], "unbind", oldVnode, oldVnode, isDestroy);
        }
      }
    }
  }
  var emptyModifiers = /* @__PURE__ */ Object.create(null);
  function normalizeDirectives(dirs, vm3) {
    var res = /* @__PURE__ */ Object.create(null);
    if (!dirs) {
      return res;
    }
    var i, dir;
    for (i = 0; i < dirs.length; i++) {
      dir = dirs[i];
      if (!dir.modifiers) {
        dir.modifiers = emptyModifiers;
      }
      res[getRawDirName(dir)] = dir;
      if (vm3._setupState && vm3._setupState.__sfc) {
        var setupDef = dir.def || resolveAsset(vm3, "_setupState", "v-" + dir.name);
        if (typeof setupDef === "function") {
          dir.def = {
            bind: setupDef,
            update: setupDef
          };
        } else {
          dir.def = setupDef;
        }
      }
      dir.def = dir.def || resolveAsset(vm3.$options, "directives", dir.name, true);
    }
    return res;
  }
  function getRawDirName(dir) {
    return dir.rawName || "".concat(dir.name, ".").concat(Object.keys(dir.modifiers || {}).join("."));
  }
  function callHook(dir, hook, vnode, oldVnode, isDestroy) {
    var fn = dir.def && dir.def[hook];
    if (fn) {
      try {
        fn(vnode.elm, dir, vnode, oldVnode, isDestroy);
      } catch (e) {
        handleError(e, vnode.context, "directive ".concat(dir.name, " ").concat(hook, " hook"));
      }
    }
  }
  var baseModules = [ref, directives$1];
  function updateAttrs(oldVnode, vnode) {
    var opts2 = vnode.componentOptions;
    if (isDef(opts2) && opts2.Ctor.options.inheritAttrs === false) {
      return;
    }
    if (isUndef(oldVnode.data.attrs) && isUndef(vnode.data.attrs)) {
      return;
    }
    var key, cur, old;
    var elm = vnode.elm;
    var oldAttrs = oldVnode.data.attrs || {};
    var attrs2 = vnode.data.attrs || {};
    if (isDef(attrs2.__ob__) || isTrue(attrs2._v_attr_proxy)) {
      attrs2 = vnode.data.attrs = extend({}, attrs2);
    }
    for (key in attrs2) {
      cur = attrs2[key];
      old = oldAttrs[key];
      if (old !== cur) {
        setAttr(elm, key, cur, vnode.data.pre);
      }
    }
    if ((isIE || isEdge) && attrs2.value !== oldAttrs.value) {
      setAttr(elm, "value", attrs2.value);
    }
    for (key in oldAttrs) {
      if (isUndef(attrs2[key])) {
        if (isXlink(key)) {
          elm.removeAttributeNS(xlinkNS, getXlinkProp(key));
        } else if (!isEnumeratedAttr(key)) {
          elm.removeAttribute(key);
        }
      }
    }
  }
  function setAttr(el, key, value, isInPre) {
    if (isInPre || el.tagName.indexOf("-") > -1) {
      baseSetAttr(el, key, value);
    } else if (isBooleanAttr(key)) {
      if (isFalsyAttrValue(value)) {
        el.removeAttribute(key);
      } else {
        value = key === "allowfullscreen" && el.tagName === "EMBED" ? "true" : key;
        el.setAttribute(key, value);
      }
    } else if (isEnumeratedAttr(key)) {
      el.setAttribute(key, convertEnumeratedValue(key, value));
    } else if (isXlink(key)) {
      if (isFalsyAttrValue(value)) {
        el.removeAttributeNS(xlinkNS, getXlinkProp(key));
      } else {
        el.setAttributeNS(xlinkNS, key, value);
      }
    } else {
      baseSetAttr(el, key, value);
    }
  }
  function baseSetAttr(el, key, value) {
    if (isFalsyAttrValue(value)) {
      el.removeAttribute(key);
    } else {
      if (isIE && !isIE9 && el.tagName === "TEXTAREA" && key === "placeholder" && value !== "" && !el.__ieph) {
        var blocker_1 = function(e) {
          e.stopImmediatePropagation();
          el.removeEventListener("input", blocker_1);
        };
        el.addEventListener("input", blocker_1);
        el.__ieph = true;
      }
      el.setAttribute(key, value);
    }
  }
  var attrs = {
    create: updateAttrs,
    update: updateAttrs
  };
  function updateClass(oldVnode, vnode) {
    var el = vnode.elm;
    var data = vnode.data;
    var oldData = oldVnode.data;
    if (isUndef(data.staticClass) && isUndef(data.class) && (isUndef(oldData) || isUndef(oldData.staticClass) && isUndef(oldData.class))) {
      return;
    }
    var cls = genClassForVnode(vnode);
    var transitionClass = el._transitionClasses;
    if (isDef(transitionClass)) {
      cls = concat(cls, stringifyClass(transitionClass));
    }
    if (cls !== el._prevClass) {
      el.setAttribute("class", cls);
      el._prevClass = cls;
    }
  }
  var klass$1 = {
    create: updateClass,
    update: updateClass
  };
  var validDivisionCharRE = /[\w).+\-_$\]]/;
  function parseFilters(exp) {
    var inSingle = false;
    var inDouble = false;
    var inTemplateString = false;
    var inRegex = false;
    var curly = 0;
    var square = 0;
    var paren = 0;
    var lastFilterIndex = 0;
    var c, prev, i, expression, filters;
    for (i = 0; i < exp.length; i++) {
      prev = c;
      c = exp.charCodeAt(i);
      if (inSingle) {
        if (c === 39 && prev !== 92)
          inSingle = false;
      } else if (inDouble) {
        if (c === 34 && prev !== 92)
          inDouble = false;
      } else if (inTemplateString) {
        if (c === 96 && prev !== 92)
          inTemplateString = false;
      } else if (inRegex) {
        if (c === 47 && prev !== 92)
          inRegex = false;
      } else if (c === 124 && // pipe
      exp.charCodeAt(i + 1) !== 124 && exp.charCodeAt(i - 1) !== 124 && !curly && !square && !paren) {
        if (expression === void 0) {
          lastFilterIndex = i + 1;
          expression = exp.slice(0, i).trim();
        } else {
          pushFilter();
        }
      } else {
        switch (c) {
          case 34:
            inDouble = true;
            break;
          // "
          case 39:
            inSingle = true;
            break;
          // '
          case 96:
            inTemplateString = true;
            break;
          // `
          case 40:
            paren++;
            break;
          // (
          case 41:
            paren--;
            break;
          // )
          case 91:
            square++;
            break;
          // [
          case 93:
            square--;
            break;
          // ]
          case 123:
            curly++;
            break;
          // {
          case 125:
            curly--;
            break;
        }
        if (c === 47) {
          var j = i - 1;
          var p = void 0;
          for (; j >= 0; j--) {
            p = exp.charAt(j);
            if (p !== " ")
              break;
          }
          if (!p || !validDivisionCharRE.test(p)) {
            inRegex = true;
          }
        }
      }
    }
    if (expression === void 0) {
      expression = exp.slice(0, i).trim();
    } else if (lastFilterIndex !== 0) {
      pushFilter();
    }
    function pushFilter() {
      (filters || (filters = [])).push(exp.slice(lastFilterIndex, i).trim());
      lastFilterIndex = i + 1;
    }
    if (filters) {
      for (i = 0; i < filters.length; i++) {
        expression = wrapFilter(expression, filters[i]);
      }
    }
    return expression;
  }
  function wrapFilter(exp, filter) {
    var i = filter.indexOf("(");
    if (i < 0) {
      return '_f("'.concat(filter, '")(').concat(exp, ")");
    } else {
      var name_1 = filter.slice(0, i);
      var args = filter.slice(i + 1);
      return '_f("'.concat(name_1, '")(').concat(exp).concat(args !== ")" ? "," + args : args);
    }
  }
  function baseWarn(msg, range2) {
    console.error("[Vue compiler]: ".concat(msg));
  }
  function pluckModuleFunction(modules2, key) {
    return modules2 ? modules2.map(function(m) {
      return m[key];
    }).filter(function(_) {
      return _;
    }) : [];
  }
  function addProp(el, name, value, range2, dynamic) {
    (el.props || (el.props = [])).push(rangeSetItem({ name, value, dynamic }, range2));
    el.plain = false;
  }
  function addAttr(el, name, value, range2, dynamic) {
    var attrs2 = dynamic ? el.dynamicAttrs || (el.dynamicAttrs = []) : el.attrs || (el.attrs = []);
    attrs2.push(rangeSetItem({ name, value, dynamic }, range2));
    el.plain = false;
  }
  function addRawAttr(el, name, value, range2) {
    el.attrsMap[name] = value;
    el.attrsList.push(rangeSetItem({ name, value }, range2));
  }
  function addDirective(el, name, rawName, value, arg, isDynamicArg, modifiers, range2) {
    (el.directives || (el.directives = [])).push(rangeSetItem({
      name,
      rawName,
      value,
      arg,
      isDynamicArg,
      modifiers
    }, range2));
    el.plain = false;
  }
  function prependModifierMarker(symbol, name, dynamic) {
    return dynamic ? "_p(".concat(name, ',"').concat(symbol, '")') : symbol + name;
  }
  function addHandler(el, name, value, modifiers, important, warn2, range2, dynamic) {
    modifiers = modifiers || emptyObject;
    if (warn2 && modifiers.prevent && modifiers.passive) {
      warn2("passive and prevent can't be used together. Passive handler can't prevent default event.", range2);
    }
    if (modifiers.right) {
      if (dynamic) {
        name = "(".concat(name, ")==='click'?'contextmenu':(").concat(name, ")");
      } else if (name === "click") {
        name = "contextmenu";
        delete modifiers.right;
      }
    } else if (modifiers.middle) {
      if (dynamic) {
        name = "(".concat(name, ")==='click'?'mouseup':(").concat(name, ")");
      } else if (name === "click") {
        name = "mouseup";
      }
    }
    if (modifiers.capture) {
      delete modifiers.capture;
      name = prependModifierMarker("!", name, dynamic);
    }
    if (modifiers.once) {
      delete modifiers.once;
      name = prependModifierMarker("~", name, dynamic);
    }
    if (modifiers.passive) {
      delete modifiers.passive;
      name = prependModifierMarker("&", name, dynamic);
    }
    var events2;
    if (modifiers.native) {
      delete modifiers.native;
      events2 = el.nativeEvents || (el.nativeEvents = {});
    } else {
      events2 = el.events || (el.events = {});
    }
    var newHandler = rangeSetItem({ value: value.trim(), dynamic }, range2);
    if (modifiers !== emptyObject) {
      newHandler.modifiers = modifiers;
    }
    var handlers = events2[name];
    if (Array.isArray(handlers)) {
      important ? handlers.unshift(newHandler) : handlers.push(newHandler);
    } else if (handlers) {
      events2[name] = important ? [newHandler, handlers] : [handlers, newHandler];
    } else {
      events2[name] = newHandler;
    }
    el.plain = false;
  }
  function getRawBindingAttr(el, name) {
    return el.rawAttrsMap[":" + name] || el.rawAttrsMap["v-bind:" + name] || el.rawAttrsMap[name];
  }
  function getBindingAttr(el, name, getStatic) {
    var dynamicValue = getAndRemoveAttr(el, ":" + name) || getAndRemoveAttr(el, "v-bind:" + name);
    if (dynamicValue != null) {
      return parseFilters(dynamicValue);
    } else if (getStatic !== false) {
      var staticValue = getAndRemoveAttr(el, name);
      if (staticValue != null) {
        return JSON.stringify(staticValue);
      }
    }
  }
  function getAndRemoveAttr(el, name, removeFromMap) {
    var val;
    if ((val = el.attrsMap[name]) != null) {
      var list = el.attrsList;
      for (var i = 0, l = list.length; i < l; i++) {
        if (list[i].name === name) {
          list.splice(i, 1);
          break;
        }
      }
    }
    if (removeFromMap) {
      delete el.attrsMap[name];
    }
    return val;
  }
  function getAndRemoveAttrByRegex(el, name) {
    var list = el.attrsList;
    for (var i = 0, l = list.length; i < l; i++) {
      var attr = list[i];
      if (name.test(attr.name)) {
        list.splice(i, 1);
        return attr;
      }
    }
  }
  function rangeSetItem(item, range2) {
    if (range2) {
      if (range2.start != null) {
        item.start = range2.start;
      }
      if (range2.end != null) {
        item.end = range2.end;
      }
    }
    return item;
  }
  function genComponentModel(el, value, modifiers) {
    var _a2 = modifiers || {}, number = _a2.number, trim = _a2.trim;
    var baseValueExpression = "$$v";
    var valueExpression = baseValueExpression;
    if (trim) {
      valueExpression = "(typeof ".concat(baseValueExpression, " === 'string'") + "? ".concat(baseValueExpression, ".trim()") + ": ".concat(baseValueExpression, ")");
    }
    if (number) {
      valueExpression = "_n(".concat(valueExpression, ")");
    }
    var assignment = genAssignmentCode(value, valueExpression);
    el.model = {
      value: "(".concat(value, ")"),
      expression: JSON.stringify(value),
      callback: "function (".concat(baseValueExpression, ") {").concat(assignment, "}")
    };
  }
  function genAssignmentCode(value, assignment) {
    var res = parseModel(value);
    if (res.key === null) {
      return "".concat(value, "=").concat(assignment);
    } else {
      return "$set(".concat(res.exp, ", ").concat(res.key, ", ").concat(assignment, ")");
    }
  }
  var len;
  var str;
  var chr;
  var index;
  var expressionPos;
  var expressionEndPos;
  function parseModel(val) {
    val = val.trim();
    len = val.length;
    if (val.indexOf("[") < 0 || val.lastIndexOf("]") < len - 1) {
      index = val.lastIndexOf(".");
      if (index > -1) {
        return {
          exp: val.slice(0, index),
          key: '"' + val.slice(index + 1) + '"'
        };
      } else {
        return {
          exp: val,
          key: null
        };
      }
    }
    str = val;
    index = expressionPos = expressionEndPos = 0;
    while (!eof()) {
      chr = next();
      if (isStringStart(chr)) {
        parseString(chr);
      } else if (chr === 91) {
        parseBracket(chr);
      }
    }
    return {
      exp: val.slice(0, expressionPos),
      key: val.slice(expressionPos + 1, expressionEndPos)
    };
  }
  function next() {
    return str.charCodeAt(++index);
  }
  function eof() {
    return index >= len;
  }
  function isStringStart(chr2) {
    return chr2 === 34 || chr2 === 39;
  }
  function parseBracket(chr2) {
    var inBracket = 1;
    expressionPos = index;
    while (!eof()) {
      chr2 = next();
      if (isStringStart(chr2)) {
        parseString(chr2);
        continue;
      }
      if (chr2 === 91)
        inBracket++;
      if (chr2 === 93)
        inBracket--;
      if (inBracket === 0) {
        expressionEndPos = index;
        break;
      }
    }
  }
  function parseString(chr2) {
    var stringQuote = chr2;
    while (!eof()) {
      chr2 = next();
      if (chr2 === stringQuote) {
        break;
      }
    }
  }
  var warn$1;
  var RANGE_TOKEN = "__r";
  var CHECKBOX_RADIO_TOKEN = "__c";
  function model$1(el, dir, _warn) {
    warn$1 = _warn;
    var value = dir.value;
    var modifiers = dir.modifiers;
    var tag = el.tag;
    var type = el.attrsMap.type;
    if (true) {
      if (tag === "input" && type === "file") {
        warn$1("<".concat(el.tag, ' v-model="').concat(value, '" type="file">:\n') + "File inputs are read only. Use a v-on:change listener instead.", el.rawAttrsMap["v-model"]);
      }
    }
    if (el.component) {
      genComponentModel(el, value, modifiers);
      return false;
    } else if (tag === "select") {
      genSelect(el, value, modifiers);
    } else if (tag === "input" && type === "checkbox") {
      genCheckboxModel(el, value, modifiers);
    } else if (tag === "input" && type === "radio") {
      genRadioModel(el, value, modifiers);
    } else if (tag === "input" || tag === "textarea") {
      genDefaultModel(el, value, modifiers);
    } else if (!config.isReservedTag(tag)) {
      genComponentModel(el, value, modifiers);
      return false;
    } else if (true) {
      warn$1("<".concat(el.tag, ' v-model="').concat(value, '">: ') + "v-model is not supported on this element type. If you are working with contenteditable, it's recommended to wrap a library dedicated for that purpose inside a custom component.", el.rawAttrsMap["v-model"]);
    }
    return true;
  }
  function genCheckboxModel(el, value, modifiers) {
    var number = modifiers && modifiers.number;
    var valueBinding = getBindingAttr(el, "value") || "null";
    var trueValueBinding = getBindingAttr(el, "true-value") || "true";
    var falseValueBinding = getBindingAttr(el, "false-value") || "false";
    addProp(el, "checked", "Array.isArray(".concat(value, ")") + "?_i(".concat(value, ",").concat(valueBinding, ")>-1") + (trueValueBinding === "true" ? ":(".concat(value, ")") : ":_q(".concat(value, ",").concat(trueValueBinding, ")")));
    addHandler(el, "change", "var $$a=".concat(value, ",") + "$$el=$event.target," + "$$c=$$el.checked?(".concat(trueValueBinding, "):(").concat(falseValueBinding, ");") + "if(Array.isArray($$a)){" + "var $$v=".concat(number ? "_n(" + valueBinding + ")" : valueBinding, ",") + "$$i=_i($$a,$$v);" + "if($$el.checked){$$i<0&&(".concat(genAssignmentCode(value, "$$a.concat([$$v])"), ")}") + "else{$$i>-1&&(".concat(genAssignmentCode(value, "$$a.slice(0,$$i).concat($$a.slice($$i+1))"), ")}") + "}else{".concat(genAssignmentCode(value, "$$c"), "}"), null, true);
  }
  function genRadioModel(el, value, modifiers) {
    var number = modifiers && modifiers.number;
    var valueBinding = getBindingAttr(el, "value") || "null";
    valueBinding = number ? "_n(".concat(valueBinding, ")") : valueBinding;
    addProp(el, "checked", "_q(".concat(value, ",").concat(valueBinding, ")"));
    addHandler(el, "change", genAssignmentCode(value, valueBinding), null, true);
  }
  function genSelect(el, value, modifiers) {
    var number = modifiers && modifiers.number;
    var selectedVal = 'Array.prototype.filter.call($event.target.options,function(o){return o.selected}).map(function(o){var val = "_value" in o ? o._value : o.value;' + "return ".concat(number ? "_n(val)" : "val", "})");
    var assignment = "$event.target.multiple ? $$selectedVal : $$selectedVal[0]";
    var code = "var $$selectedVal = ".concat(selectedVal, ";");
    code = "".concat(code, " ").concat(genAssignmentCode(value, assignment));
    addHandler(el, "change", code, null, true);
  }
  function genDefaultModel(el, value, modifiers) {
    var type = el.attrsMap.type;
    if (true) {
      var value_1 = el.attrsMap["v-bind:value"] || el.attrsMap[":value"];
      var typeBinding = el.attrsMap["v-bind:type"] || el.attrsMap[":type"];
      if (value_1 && !typeBinding) {
        var binding = el.attrsMap["v-bind:value"] ? "v-bind:value" : ":value";
        warn$1("".concat(binding, '="').concat(value_1, '" conflicts with v-model on the same element ') + "because the latter already expands to a value binding internally", el.rawAttrsMap[binding]);
      }
    }
    var _a2 = modifiers || {}, lazy = _a2.lazy, number = _a2.number, trim = _a2.trim;
    var needCompositionGuard = !lazy && type !== "range";
    var event = lazy ? "change" : type === "range" ? RANGE_TOKEN : "input";
    var valueExpression = "$event.target.value";
    if (trim) {
      valueExpression = "$event.target.value.trim()";
    }
    if (number) {
      valueExpression = "_n(".concat(valueExpression, ")");
    }
    var code = genAssignmentCode(value, valueExpression);
    if (needCompositionGuard) {
      code = "if($event.target.composing)return;".concat(code);
    }
    addProp(el, "value", "(".concat(value, ")"));
    addHandler(el, event, code, null, true);
    if (trim || number) {
      addHandler(el, "blur", "$forceUpdate()");
    }
  }
  function normalizeEvents(on2) {
    if (isDef(on2[RANGE_TOKEN])) {
      var event_1 = isIE ? "change" : "input";
      on2[event_1] = [].concat(on2[RANGE_TOKEN], on2[event_1] || []);
      delete on2[RANGE_TOKEN];
    }
    if (isDef(on2[CHECKBOX_RADIO_TOKEN])) {
      on2.change = [].concat(on2[CHECKBOX_RADIO_TOKEN], on2.change || []);
      delete on2[CHECKBOX_RADIO_TOKEN];
    }
  }
  var target;
  function createOnceHandler(event, handler, capture) {
    var _target = target;
    return function onceHandler() {
      var res = handler.apply(null, arguments);
      if (res !== null) {
        remove(event, onceHandler, capture, _target);
      }
    };
  }
  var useMicrotaskFix = isUsingMicroTask && !(isFF && Number(isFF[1]) <= 53);
  function add(name, handler, capture, passive) {
    if (useMicrotaskFix) {
      var attachedTimestamp_1 = currentFlushTimestamp;
      var original_1 = handler;
      handler = original_1._wrapper = function(e) {
        if (
          // no bubbling, should always fire.
          // this is just a safety net in case event.timeStamp is unreliable in
          // certain weird environments...
          e.target === e.currentTarget || // event is fired after handler attachment
          e.timeStamp >= attachedTimestamp_1 || // bail for environments that have buggy event.timeStamp implementations
          // #9462 iOS 9 bug: event.timeStamp is 0 after history.pushState
          // #9681 QtWebEngine event.timeStamp is negative value
          e.timeStamp <= 0 || // #9448 bail if event is fired in another document in a multi-page
          // electron/nw.js app, since event.timeStamp will be using a different
          // starting reference
          e.target.ownerDocument !== document
        ) {
          return original_1.apply(this, arguments);
        }
      };
    }
    target.addEventListener(name, handler, supportsPassive ? { capture, passive } : capture);
  }
  function remove(name, handler, capture, _target) {
    (_target || target).removeEventListener(
      name,
      //@ts-expect-error
      handler._wrapper || handler,
      capture
    );
  }
  function updateDOMListeners(oldVnode, vnode) {
    if (isUndef(oldVnode.data.on) && isUndef(vnode.data.on)) {
      return;
    }
    var on2 = vnode.data.on || {};
    var oldOn = oldVnode.data.on || {};
    target = vnode.elm || oldVnode.elm;
    normalizeEvents(on2);
    updateListeners(on2, oldOn, add, remove, createOnceHandler, vnode.context);
    target = void 0;
  }
  var events = {
    create: updateDOMListeners,
    update: updateDOMListeners,
    // @ts-expect-error emptyNode has actually data
    destroy: function(vnode) {
      return updateDOMListeners(vnode, emptyNode);
    }
  };
  var svgContainer;
  function updateDOMProps(oldVnode, vnode) {
    if (isUndef(oldVnode.data.domProps) && isUndef(vnode.data.domProps)) {
      return;
    }
    var key, cur;
    var elm = vnode.elm;
    var oldProps = oldVnode.data.domProps || {};
    var props2 = vnode.data.domProps || {};
    if (isDef(props2.__ob__) || isTrue(props2._v_attr_proxy)) {
      props2 = vnode.data.domProps = extend({}, props2);
    }
    for (key in oldProps) {
      if (!(key in props2)) {
        elm[key] = "";
      }
    }
    for (key in props2) {
      cur = props2[key];
      if (key === "textContent" || key === "innerHTML") {
        if (vnode.children)
          vnode.children.length = 0;
        if (cur === oldProps[key])
          continue;
        if (elm.childNodes.length === 1) {
          elm.removeChild(elm.childNodes[0]);
        }
      }
      if (key === "value" && elm.tagName !== "PROGRESS") {
        elm._value = cur;
        var strCur = isUndef(cur) ? "" : String(cur);
        if (shouldUpdateValue(elm, strCur)) {
          elm.value = strCur;
        }
      } else if (key === "innerHTML" && isSVG(elm.tagName) && isUndef(elm.innerHTML)) {
        svgContainer = svgContainer || document.createElement("div");
        svgContainer.innerHTML = "<svg>".concat(cur, "</svg>");
        var svg = svgContainer.firstChild;
        while (elm.firstChild) {
          elm.removeChild(elm.firstChild);
        }
        while (svg.firstChild) {
          elm.appendChild(svg.firstChild);
        }
      } else if (
        // skip the update if old and new VDOM state is the same.
        // `value` is handled separately because the DOM value may be temporarily
        // out of sync with VDOM state due to focus, composition and modifiers.
        // This  #4521 by skipping the unnecessary `checked` update.
        cur !== oldProps[key]
      ) {
        try {
          elm[key] = cur;
        } catch (e) {
        }
      }
    }
  }
  function shouldUpdateValue(elm, checkVal) {
    return (
      //@ts-expect-error
      !elm.composing && (elm.tagName === "OPTION" || isNotInFocusAndDirty(elm, checkVal) || isDirtyWithModifiers(elm, checkVal))
    );
  }
  function isNotInFocusAndDirty(elm, checkVal) {
    var notInFocus = true;
    try {
      notInFocus = document.activeElement !== elm;
    } catch (e) {
    }
    return notInFocus && elm.value !== checkVal;
  }
  function isDirtyWithModifiers(elm, newVal) {
    var value = elm.value;
    var modifiers = elm._vModifiers;
    if (isDef(modifiers)) {
      if (modifiers.number) {
        return toNumber(value) !== toNumber(newVal);
      }
      if (modifiers.trim) {
        return value.trim() !== newVal.trim();
      }
    }
    return value !== newVal;
  }
  var domProps = {
    create: updateDOMProps,
    update: updateDOMProps
  };
  var parseStyleText = cached(function(cssText) {
    var res = {};
    var listDelimiter = /;(?![^(]*\))/g;
    var propertyDelimiter = /:(.+)/;
    cssText.split(listDelimiter).forEach(function(item) {
      if (item) {
        var tmp = item.split(propertyDelimiter);
        tmp.length > 1 && (res[tmp[0].trim()] = tmp[1].trim());
      }
    });
    return res;
  });
  function normalizeStyleData(data) {
    var style2 = normalizeStyleBinding(data.style);
    return data.staticStyle ? extend(data.staticStyle, style2) : style2;
  }
  function normalizeStyleBinding(bindingStyle) {
    if (Array.isArray(bindingStyle)) {
      return toObject(bindingStyle);
    }
    if (typeof bindingStyle === "string") {
      return parseStyleText(bindingStyle);
    }
    return bindingStyle;
  }
  function getStyle(vnode, checkChild) {
    var res = {};
    var styleData;
    if (checkChild) {
      var childNode = vnode;
      while (childNode.componentInstance) {
        childNode = childNode.componentInstance._vnode;
        if (childNode && childNode.data && (styleData = normalizeStyleData(childNode.data))) {
          extend(res, styleData);
        }
      }
    }
    if (styleData = normalizeStyleData(vnode.data)) {
      extend(res, styleData);
    }
    var parentNode2 = vnode;
    while (parentNode2 = parentNode2.parent) {
      if (parentNode2.data && (styleData = normalizeStyleData(parentNode2.data))) {
        extend(res, styleData);
      }
    }
    return res;
  }
  var cssVarRE = /^--/;
  var importantRE = /\s*!important$/;
  var setProp = function(el, name, val) {
    if (cssVarRE.test(name)) {
      el.style.setProperty(name, val);
    } else if (importantRE.test(val)) {
      el.style.setProperty(hyphenate(name), val.replace(importantRE, ""), "important");
    } else {
      var normalizedName = normalize(name);
      if (Array.isArray(val)) {
        for (var i = 0, len2 = val.length; i < len2; i++) {
          el.style[normalizedName] = val[i];
        }
      } else {
        el.style[normalizedName] = val;
      }
    }
  };
  var vendorNames = ["Webkit", "Moz", "ms"];
  var emptyStyle;
  var normalize = cached(function(prop) {
    emptyStyle = emptyStyle || document.createElement("div").style;
    prop = camelize(prop);
    if (prop !== "filter" && prop in emptyStyle) {
      return prop;
    }
    var capName = prop.charAt(0).toUpperCase() + prop.slice(1);
    for (var i = 0; i < vendorNames.length; i++) {
      var name_1 = vendorNames[i] + capName;
      if (name_1 in emptyStyle) {
        return name_1;
      }
    }
  });
  function updateStyle(oldVnode, vnode) {
    var data = vnode.data;
    var oldData = oldVnode.data;
    if (isUndef(data.staticStyle) && isUndef(data.style) && isUndef(oldData.staticStyle) && isUndef(oldData.style)) {
      return;
    }
    var cur, name;
    var el = vnode.elm;
    var oldStaticStyle = oldData.staticStyle;
    var oldStyleBinding = oldData.normalizedStyle || oldData.style || {};
    var oldStyle = oldStaticStyle || oldStyleBinding;
    var style2 = normalizeStyleBinding(vnode.data.style) || {};
    vnode.data.normalizedStyle = isDef(style2.__ob__) ? extend({}, style2) : style2;
    var newStyle = getStyle(vnode, true);
    for (name in oldStyle) {
      if (isUndef(newStyle[name])) {
        setProp(el, name, "");
      }
    }
    for (name in newStyle) {
      cur = newStyle[name];
      setProp(el, name, cur == null ? "" : cur);
    }
  }
  var style$1 = {
    create: updateStyle,
    update: updateStyle
  };
  var whitespaceRE$1 = /\s+/;
  function addClass(el, cls) {
    if (!cls || !(cls = cls.trim())) {
      return;
    }
    if (el.classList) {
      if (cls.indexOf(" ") > -1) {
        cls.split(whitespaceRE$1).forEach(function(c) {
          return el.classList.add(c);
        });
      } else {
        el.classList.add(cls);
      }
    } else {
      var cur = " ".concat(el.getAttribute("class") || "", " ");
      if (cur.indexOf(" " + cls + " ") < 0) {
        el.setAttribute("class", (cur + cls).trim());
      }
    }
  }
  function removeClass(el, cls) {
    if (!cls || !(cls = cls.trim())) {
      return;
    }
    if (el.classList) {
      if (cls.indexOf(" ") > -1) {
        cls.split(whitespaceRE$1).forEach(function(c) {
          return el.classList.remove(c);
        });
      } else {
        el.classList.remove(cls);
      }
      if (!el.classList.length) {
        el.removeAttribute("class");
      }
    } else {
      var cur = " ".concat(el.getAttribute("class") || "", " ");
      var tar = " " + cls + " ";
      while (cur.indexOf(tar) >= 0) {
        cur = cur.replace(tar, " ");
      }
      cur = cur.trim();
      if (cur) {
        el.setAttribute("class", cur);
      } else {
        el.removeAttribute("class");
      }
    }
  }
  function resolveTransition(def2) {
    if (!def2) {
      return;
    }
    if (typeof def2 === "object") {
      var res = {};
      if (def2.css !== false) {
        extend(res, autoCssTransition(def2.name || "v"));
      }
      extend(res, def2);
      return res;
    } else if (typeof def2 === "string") {
      return autoCssTransition(def2);
    }
  }
  var autoCssTransition = cached(function(name) {
    return {
      enterClass: "".concat(name, "-enter"),
      enterToClass: "".concat(name, "-enter-to"),
      enterActiveClass: "".concat(name, "-enter-active"),
      leaveClass: "".concat(name, "-leave"),
      leaveToClass: "".concat(name, "-leave-to"),
      leaveActiveClass: "".concat(name, "-leave-active")
    };
  });
  var hasTransition = inBrowser && !isIE9;
  var TRANSITION = "transition";
  var ANIMATION = "animation";
  var transitionProp = "transition";
  var transitionEndEvent = "transitionend";
  var animationProp = "animation";
  var animationEndEvent = "animationend";
  if (hasTransition) {
    if (window.ontransitionend === void 0 && window.onwebkittransitionend !== void 0) {
      transitionProp = "WebkitTransition";
      transitionEndEvent = "webkitTransitionEnd";
    }
    if (window.onanimationend === void 0 && window.onwebkitanimationend !== void 0) {
      animationProp = "WebkitAnimation";
      animationEndEvent = "webkitAnimationEnd";
    }
  }
  var raf = inBrowser ? window.requestAnimationFrame ? window.requestAnimationFrame.bind(window) : setTimeout : (
    /* istanbul ignore next */
    function(fn) {
      return fn();
    }
  );
  function nextFrame(fn) {
    raf(function() {
      raf(fn);
    });
  }
  function addTransitionClass(el, cls) {
    var transitionClasses = el._transitionClasses || (el._transitionClasses = []);
    if (transitionClasses.indexOf(cls) < 0) {
      transitionClasses.push(cls);
      addClass(el, cls);
    }
  }
  function removeTransitionClass(el, cls) {
    if (el._transitionClasses) {
      remove$2(el._transitionClasses, cls);
    }
    removeClass(el, cls);
  }
  function whenTransitionEnds(el, expectedType, cb) {
    var _a2 = getTransitionInfo(el, expectedType), type = _a2.type, timeout = _a2.timeout, propCount = _a2.propCount;
    if (!type)
      return cb();
    var event = type === TRANSITION ? transitionEndEvent : animationEndEvent;
    var ended = 0;
    var end = function() {
      el.removeEventListener(event, onEnd);
      cb();
    };
    var onEnd = function(e) {
      if (e.target === el) {
        if (++ended >= propCount) {
          end();
        }
      }
    };
    setTimeout(function() {
      if (ended < propCount) {
        end();
      }
    }, timeout + 1);
    el.addEventListener(event, onEnd);
  }
  var transformRE = /\b(transform|all)(,|$)/;
  function getTransitionInfo(el, expectedType) {
    var styles = window.getComputedStyle(el);
    var transitionDelays = (styles[transitionProp + "Delay"] || "").split(", ");
    var transitionDurations = (styles[transitionProp + "Duration"] || "").split(", ");
    var transitionTimeout = getTimeout(transitionDelays, transitionDurations);
    var animationDelays = (styles[animationProp + "Delay"] || "").split(", ");
    var animationDurations = (styles[animationProp + "Duration"] || "").split(", ");
    var animationTimeout = getTimeout(animationDelays, animationDurations);
    var type;
    var timeout = 0;
    var propCount = 0;
    if (expectedType === TRANSITION) {
      if (transitionTimeout > 0) {
        type = TRANSITION;
        timeout = transitionTimeout;
        propCount = transitionDurations.length;
      }
    } else if (expectedType === ANIMATION) {
      if (animationTimeout > 0) {
        type = ANIMATION;
        timeout = animationTimeout;
        propCount = animationDurations.length;
      }
    } else {
      timeout = Math.max(transitionTimeout, animationTimeout);
      type = timeout > 0 ? transitionTimeout > animationTimeout ? TRANSITION : ANIMATION : null;
      propCount = type ? type === TRANSITION ? transitionDurations.length : animationDurations.length : 0;
    }
    var hasTransform = type === TRANSITION && transformRE.test(styles[transitionProp + "Property"]);
    return {
      type,
      timeout,
      propCount,
      hasTransform
    };
  }
  function getTimeout(delays, durations) {
    while (delays.length < durations.length) {
      delays = delays.concat(delays);
    }
    return Math.max.apply(null, durations.map(function(d, i) {
      return toMs(d) + toMs(delays[i]);
    }));
  }
  function toMs(s) {
    return Number(s.slice(0, -1).replace(",", ".")) * 1e3;
  }
  function enter(vnode, toggleDisplay) {
    var el = vnode.elm;
    if (isDef(el._leaveCb)) {
      el._leaveCb.cancelled = true;
      el._leaveCb();
    }
    var data = resolveTransition(vnode.data.transition);
    if (isUndef(data)) {
      return;
    }
    if (isDef(el._enterCb) || el.nodeType !== 1) {
      return;
    }
    var css = data.css, type = data.type, enterClass = data.enterClass, enterToClass = data.enterToClass, enterActiveClass = data.enterActiveClass, appearClass = data.appearClass, appearToClass = data.appearToClass, appearActiveClass = data.appearActiveClass, beforeEnter = data.beforeEnter, enter2 = data.enter, afterEnter = data.afterEnter, enterCancelled = data.enterCancelled, beforeAppear = data.beforeAppear, appear = data.appear, afterAppear = data.afterAppear, appearCancelled = data.appearCancelled, duration = data.duration;
    var context = activeInstance;
    var transitionNode = activeInstance.$vnode;
    while (transitionNode && transitionNode.parent) {
      context = transitionNode.context;
      transitionNode = transitionNode.parent;
    }
    var isAppear = !context._isMounted || !vnode.isRootInsert;
    if (isAppear && !appear && appear !== "") {
      return;
    }
    var startClass = isAppear && appearClass ? appearClass : enterClass;
    var activeClass = isAppear && appearActiveClass ? appearActiveClass : enterActiveClass;
    var toClass = isAppear && appearToClass ? appearToClass : enterToClass;
    var beforeEnterHook = isAppear ? beforeAppear || beforeEnter : beforeEnter;
    var enterHook = isAppear ? isFunction(appear) ? appear : enter2 : enter2;
    var afterEnterHook = isAppear ? afterAppear || afterEnter : afterEnter;
    var enterCancelledHook = isAppear ? appearCancelled || enterCancelled : enterCancelled;
    var explicitEnterDuration = toNumber(isObject(duration) ? duration.enter : duration);
    if (explicitEnterDuration != null) {
      checkDuration(explicitEnterDuration, "enter", vnode);
    }
    var expectsCSS = css !== false && !isIE9;
    var userWantsControl = getHookArgumentsLength(enterHook);
    var cb = el._enterCb = once(function() {
      if (expectsCSS) {
        removeTransitionClass(el, toClass);
        removeTransitionClass(el, activeClass);
      }
      if (cb.cancelled) {
        if (expectsCSS) {
          removeTransitionClass(el, startClass);
        }
        enterCancelledHook && enterCancelledHook(el);
      } else {
        afterEnterHook && afterEnterHook(el);
      }
      el._enterCb = null;
    });
    if (!vnode.data.show) {
      mergeVNodeHook(vnode, "insert", function() {
        var parent = el.parentNode;
        var pendingNode = parent && parent._pending && parent._pending[vnode.key];
        if (pendingNode && pendingNode.tag === vnode.tag && pendingNode.elm._leaveCb) {
          pendingNode.elm._leaveCb();
        }
        enterHook && enterHook(el, cb);
      });
    }
    beforeEnterHook && beforeEnterHook(el);
    if (expectsCSS) {
      addTransitionClass(el, startClass);
      addTransitionClass(el, activeClass);
      nextFrame(function() {
        removeTransitionClass(el, startClass);
        if (!cb.cancelled) {
          addTransitionClass(el, toClass);
          if (!userWantsControl) {
            if (isValidDuration(explicitEnterDuration)) {
              setTimeout(cb, explicitEnterDuration);
            } else {
              whenTransitionEnds(el, type, cb);
            }
          }
        }
      });
    }
    if (vnode.data.show) {
      toggleDisplay && toggleDisplay();
      enterHook && enterHook(el, cb);
    }
    if (!expectsCSS && !userWantsControl) {
      cb();
    }
  }
  function leave(vnode, rm) {
    var el = vnode.elm;
    if (isDef(el._enterCb)) {
      el._enterCb.cancelled = true;
      el._enterCb();
    }
    var data = resolveTransition(vnode.data.transition);
    if (isUndef(data) || el.nodeType !== 1) {
      return rm();
    }
    if (isDef(el._leaveCb)) {
      return;
    }
    var css = data.css, type = data.type, leaveClass = data.leaveClass, leaveToClass = data.leaveToClass, leaveActiveClass = data.leaveActiveClass, beforeLeave = data.beforeLeave, leave2 = data.leave, afterLeave = data.afterLeave, leaveCancelled = data.leaveCancelled, delayLeave = data.delayLeave, duration = data.duration;
    var expectsCSS = css !== false && !isIE9;
    var userWantsControl = getHookArgumentsLength(leave2);
    var explicitLeaveDuration = toNumber(isObject(duration) ? duration.leave : duration);
    if (isDef(explicitLeaveDuration)) {
      checkDuration(explicitLeaveDuration, "leave", vnode);
    }
    var cb = el._leaveCb = once(function() {
      if (el.parentNode && el.parentNode._pending) {
        el.parentNode._pending[vnode.key] = null;
      }
      if (expectsCSS) {
        removeTransitionClass(el, leaveToClass);
        removeTransitionClass(el, leaveActiveClass);
      }
      if (cb.cancelled) {
        if (expectsCSS) {
          removeTransitionClass(el, leaveClass);
        }
        leaveCancelled && leaveCancelled(el);
      } else {
        rm();
        afterLeave && afterLeave(el);
      }
      el._leaveCb = null;
    });
    if (delayLeave) {
      delayLeave(performLeave);
    } else {
      performLeave();
    }
    function performLeave() {
      if (cb.cancelled) {
        return;
      }
      if (!vnode.data.show && el.parentNode) {
        (el.parentNode._pending || (el.parentNode._pending = {}))[vnode.key] = vnode;
      }
      beforeLeave && beforeLeave(el);
      if (expectsCSS) {
        addTransitionClass(el, leaveClass);
        addTransitionClass(el, leaveActiveClass);
        nextFrame(function() {
          removeTransitionClass(el, leaveClass);
          if (!cb.cancelled) {
            addTransitionClass(el, leaveToClass);
            if (!userWantsControl) {
              if (isValidDuration(explicitLeaveDuration)) {
                setTimeout(cb, explicitLeaveDuration);
              } else {
                whenTransitionEnds(el, type, cb);
              }
            }
          }
        });
      }
      leave2 && leave2(el, cb);
      if (!expectsCSS && !userWantsControl) {
        cb();
      }
    }
  }
  function checkDuration(val, name, vnode) {
    if (typeof val !== "number") {
      warn$2("<transition> explicit ".concat(name, " duration is not a valid number - ") + "got ".concat(JSON.stringify(val), "."), vnode.context);
    } else if (isNaN(val)) {
      warn$2("<transition> explicit ".concat(name, " duration is NaN - ") + "the duration expression might be incorrect.", vnode.context);
    }
  }
  function isValidDuration(val) {
    return typeof val === "number" && !isNaN(val);
  }
  function getHookArgumentsLength(fn) {
    if (isUndef(fn)) {
      return false;
    }
    var invokerFns = fn.fns;
    if (isDef(invokerFns)) {
      return getHookArgumentsLength(Array.isArray(invokerFns) ? invokerFns[0] : invokerFns);
    } else {
      return (fn._length || fn.length) > 1;
    }
  }
  function _enter(_, vnode) {
    if (vnode.data.show !== true) {
      enter(vnode);
    }
  }
  var transition = inBrowser ? {
    create: _enter,
    activate: _enter,
    remove: function(vnode, rm) {
      if (vnode.data.show !== true) {
        leave(vnode, rm);
      } else {
        rm();
      }
    }
  } : {};
  var platformModules = [attrs, klass$1, events, domProps, style$1, transition];
  var modules$1 = platformModules.concat(baseModules);
  var patch = createPatchFunction({ nodeOps, modules: modules$1 });
  if (isIE9) {
    document.addEventListener("selectionchange", function() {
      var el = document.activeElement;
      if (el && el.vmodel) {
        trigger(el, "input");
      }
    });
  }
  var directive = {
    inserted: function(el, binding, vnode, oldVnode) {
      if (vnode.tag === "select") {
        if (oldVnode.elm && !oldVnode.elm._vOptions) {
          mergeVNodeHook(vnode, "postpatch", function() {
            directive.componentUpdated(el, binding, vnode);
          });
        } else {
          setSelected(el, binding, vnode.context);
        }
        el._vOptions = [].map.call(el.options, getValue);
      } else if (vnode.tag === "textarea" || isTextInputType(el.type)) {
        el._vModifiers = binding.modifiers;
        if (!binding.modifiers.lazy) {
          el.addEventListener("compositionstart", onCompositionStart);
          el.addEventListener("compositionend", onCompositionEnd);
          el.addEventListener("change", onCompositionEnd);
          if (isIE9) {
            el.vmodel = true;
          }
        }
      }
    },
    componentUpdated: function(el, binding, vnode) {
      if (vnode.tag === "select") {
        setSelected(el, binding, vnode.context);
        var prevOptions_1 = el._vOptions;
        var curOptions_1 = el._vOptions = [].map.call(el.options, getValue);
        if (curOptions_1.some(function(o, i) {
          return !looseEqual(o, prevOptions_1[i]);
        })) {
          var needReset = el.multiple ? binding.value.some(function(v) {
            return hasNoMatchingOption(v, curOptions_1);
          }) : binding.value !== binding.oldValue && hasNoMatchingOption(binding.value, curOptions_1);
          if (needReset) {
            trigger(el, "change");
          }
        }
      }
    }
  };
  function setSelected(el, binding, vm3) {
    actuallySetSelected(el, binding, vm3);
    if (isIE || isEdge) {
      setTimeout(function() {
        actuallySetSelected(el, binding, vm3);
      }, 0);
    }
  }
  function actuallySetSelected(el, binding, vm3) {
    var value = binding.value;
    var isMultiple = el.multiple;
    if (isMultiple && !Array.isArray(value)) {
      warn$2('<select multiple v-model="'.concat(binding.expression, '"> ') + "expects an Array value for its binding, but got ".concat(Object.prototype.toString.call(value).slice(8, -1)), vm3);
      return;
    }
    var selected, option;
    for (var i = 0, l = el.options.length; i < l; i++) {
      option = el.options[i];
      if (isMultiple) {
        selected = looseIndexOf(value, getValue(option)) > -1;
        if (option.selected !== selected) {
          option.selected = selected;
        }
      } else {
        if (looseEqual(getValue(option), value)) {
          if (el.selectedIndex !== i) {
            el.selectedIndex = i;
          }
          return;
        }
      }
    }
    if (!isMultiple) {
      el.selectedIndex = -1;
    }
  }
  function hasNoMatchingOption(value, options) {
    return options.every(function(o) {
      return !looseEqual(o, value);
    });
  }
  function getValue(option) {
    return "_value" in option ? option._value : option.value;
  }
  function onCompositionStart(e) {
    e.target.composing = true;
  }
  function onCompositionEnd(e) {
    if (!e.target.composing)
      return;
    e.target.composing = false;
    trigger(e.target, "input");
  }
  function trigger(el, type) {
    var e = document.createEvent("HTMLEvents");
    e.initEvent(type, true, true);
    el.dispatchEvent(e);
  }
  function locateNode(vnode) {
    return vnode.componentInstance && (!vnode.data || !vnode.data.transition) ? locateNode(vnode.componentInstance._vnode) : vnode;
  }
  var show = {
    bind: function(el, _a2, vnode) {
      var value = _a2.value;
      vnode = locateNode(vnode);
      var transition2 = vnode.data && vnode.data.transition;
      var originalDisplay = el.__vOriginalDisplay = el.style.display === "none" ? "" : el.style.display;
      if (value && transition2) {
        vnode.data.show = true;
        enter(vnode, function() {
          el.style.display = originalDisplay;
        });
      } else {
        el.style.display = value ? originalDisplay : "none";
      }
    },
    update: function(el, _a2, vnode) {
      var value = _a2.value, oldValue = _a2.oldValue;
      if (!value === !oldValue)
        return;
      vnode = locateNode(vnode);
      var transition2 = vnode.data && vnode.data.transition;
      if (transition2) {
        vnode.data.show = true;
        if (value) {
          enter(vnode, function() {
            el.style.display = el.__vOriginalDisplay;
          });
        } else {
          leave(vnode, function() {
            el.style.display = "none";
          });
        }
      } else {
        el.style.display = value ? el.__vOriginalDisplay : "none";
      }
    },
    unbind: function(el, binding, vnode, oldVnode, isDestroy) {
      if (!isDestroy) {
        el.style.display = el.__vOriginalDisplay;
      }
    }
  };
  var platformDirectives = {
    model: directive,
    show
  };
  var transitionProps = {
    name: String,
    appear: Boolean,
    css: Boolean,
    mode: String,
    type: String,
    enterClass: String,
    leaveClass: String,
    enterToClass: String,
    leaveToClass: String,
    enterActiveClass: String,
    leaveActiveClass: String,
    appearClass: String,
    appearActiveClass: String,
    appearToClass: String,
    duration: [Number, String, Object]
  };
  function getRealChild(vnode) {
    var compOptions = vnode && vnode.componentOptions;
    if (compOptions && compOptions.Ctor.options.abstract) {
      return getRealChild(getFirstComponentChild(compOptions.children));
    } else {
      return vnode;
    }
  }
  function extractTransitionData(comp) {
    var data = {};
    var options = comp.$options;
    for (var key in options.propsData) {
      data[key] = comp[key];
    }
    var listeners = options._parentListeners;
    for (var key in listeners) {
      data[camelize(key)] = listeners[key];
    }
    return data;
  }
  function placeholder(h, rawChild) {
    if (/\d-keep-alive$/.test(rawChild.tag)) {
      return h("keep-alive", {
        props: rawChild.componentOptions.propsData
      });
    }
  }
  function hasParentTransition(vnode) {
    while (vnode = vnode.parent) {
      if (vnode.data.transition) {
        return true;
      }
    }
  }
  function isSameChild(child, oldChild) {
    return oldChild.key === child.key && oldChild.tag === child.tag;
  }
  var isNotTextNode = function(c) {
    return c.tag || isAsyncPlaceholder(c);
  };
  var isVShowDirective = function(d) {
    return d.name === "show";
  };
  var Transition = {
    name: "transition",
    props: transitionProps,
    abstract: true,
    render: function(h) {
      var _this = this;
      var children = this.$slots.default;
      if (!children) {
        return;
      }
      children = children.filter(isNotTextNode);
      if (!children.length) {
        return;
      }
      if (children.length > 1) {
        warn$2("<transition> can only be used on a single element. Use <transition-group> for lists.", this.$parent);
      }
      var mode = this.mode;
      if (mode && mode !== "in-out" && mode !== "out-in") {
        warn$2("invalid <transition> mode: " + mode, this.$parent);
      }
      var rawChild = children[0];
      if (hasParentTransition(this.$vnode)) {
        return rawChild;
      }
      var child = getRealChild(rawChild);
      if (!child) {
        return rawChild;
      }
      if (this._leaving) {
        return placeholder(h, rawChild);
      }
      var id = "__transition-".concat(this._uid, "-");
      child.key = child.key == null ? child.isComment ? id + "comment" : id + child.tag : isPrimitive(child.key) ? String(child.key).indexOf(id) === 0 ? child.key : id + child.key : child.key;
      var data = (child.data || (child.data = {})).transition = extractTransitionData(this);
      var oldRawChild = this._vnode;
      var oldChild = getRealChild(oldRawChild);
      if (child.data.directives && child.data.directives.some(isVShowDirective)) {
        child.data.show = true;
      }
      if (oldChild && oldChild.data && !isSameChild(child, oldChild) && !isAsyncPlaceholder(oldChild) && // #6687 component root is a comment node
      !(oldChild.componentInstance && oldChild.componentInstance._vnode.isComment)) {
        var oldData = oldChild.data.transition = extend({}, data);
        if (mode === "out-in") {
          this._leaving = true;
          mergeVNodeHook(oldData, "afterLeave", function() {
            _this._leaving = false;
            _this.$forceUpdate();
          });
          return placeholder(h, rawChild);
        } else if (mode === "in-out") {
          if (isAsyncPlaceholder(child)) {
            return oldRawChild;
          }
          var delayedLeave_1;
          var performLeave = function() {
            delayedLeave_1();
          };
          mergeVNodeHook(data, "afterEnter", performLeave);
          mergeVNodeHook(data, "enterCancelled", performLeave);
          mergeVNodeHook(oldData, "delayLeave", function(leave2) {
            delayedLeave_1 = leave2;
          });
        }
      }
      return rawChild;
    }
  };
  var props = extend({
    tag: String,
    moveClass: String
  }, transitionProps);
  delete props.mode;
  var TransitionGroup = {
    props,
    beforeMount: function() {
      var _this = this;
      var update = this._update;
      this._update = function(vnode, hydrating) {
        var restoreActiveInstance = setActiveInstance(_this);
        _this.__patch__(
          _this._vnode,
          _this.kept,
          false,
          // hydrating
          true
          // removeOnly (!important, avoids unnecessary moves)
        );
        _this._vnode = _this.kept;
        restoreActiveInstance();
        update.call(_this, vnode, hydrating);
      };
    },
    render: function(h) {
      var tag = this.tag || this.$vnode.data.tag || "span";
      var map = /* @__PURE__ */ Object.create(null);
      var prevChildren = this.prevChildren = this.children;
      var rawChildren = this.$slots.default || [];
      var children = this.children = [];
      var transitionData = extractTransitionData(this);
      for (var i = 0; i < rawChildren.length; i++) {
        var c = rawChildren[i];
        if (c.tag) {
          if (c.key != null && String(c.key).indexOf("__vlist") !== 0) {
            children.push(c);
            map[c.key] = c;
            (c.data || (c.data = {})).transition = transitionData;
          } else if (true) {
            var opts2 = c.componentOptions;
            var name_1 = opts2 ? getComponentName(opts2.Ctor.options) || opts2.tag || "" : c.tag;
            warn$2("<transition-group> children must be keyed: <".concat(name_1, ">"));
          }
        }
      }
      if (prevChildren) {
        var kept = [];
        var removed = [];
        for (var i = 0; i < prevChildren.length; i++) {
          var c = prevChildren[i];
          c.data.transition = transitionData;
          c.data.pos = c.elm.getBoundingClientRect();
          if (map[c.key]) {
            kept.push(c);
          } else {
            removed.push(c);
          }
        }
        this.kept = h(tag, null, kept);
        this.removed = removed;
      }
      return h(tag, null, children);
    },
    updated: function() {
      var children = this.prevChildren;
      var moveClass = this.moveClass || (this.name || "v") + "-move";
      if (!children.length || !this.hasMove(children[0].elm, moveClass)) {
        return;
      }
      children.forEach(callPendingCbs);
      children.forEach(recordPosition);
      children.forEach(applyTranslation);
      this._reflow = document.body.offsetHeight;
      children.forEach(function(c) {
        if (c.data.moved) {
          var el_1 = c.elm;
          var s = el_1.style;
          addTransitionClass(el_1, moveClass);
          s.transform = s.WebkitTransform = s.transitionDuration = "";
          el_1.addEventListener(transitionEndEvent, el_1._moveCb = function cb(e) {
            if (e && e.target !== el_1) {
              return;
            }
            if (!e || /transform$/.test(e.propertyName)) {
              el_1.removeEventListener(transitionEndEvent, cb);
              el_1._moveCb = null;
              removeTransitionClass(el_1, moveClass);
            }
          });
        }
      });
    },
    methods: {
      hasMove: function(el, moveClass) {
        if (!hasTransition) {
          return false;
        }
        if (this._hasMove) {
          return this._hasMove;
        }
        var clone = el.cloneNode();
        if (el._transitionClasses) {
          el._transitionClasses.forEach(function(cls) {
            removeClass(clone, cls);
          });
        }
        addClass(clone, moveClass);
        clone.style.display = "none";
        this.$el.appendChild(clone);
        var info = getTransitionInfo(clone);
        this.$el.removeChild(clone);
        return this._hasMove = info.hasTransform;
      }
    }
  };
  function callPendingCbs(c) {
    if (c.elm._moveCb) {
      c.elm._moveCb();
    }
    if (c.elm._enterCb) {
      c.elm._enterCb();
    }
  }
  function recordPosition(c) {
    c.data.newPos = c.elm.getBoundingClientRect();
  }
  function applyTranslation(c) {
    var oldPos = c.data.pos;
    var newPos = c.data.newPos;
    var dx = oldPos.left - newPos.left;
    var dy = oldPos.top - newPos.top;
    if (dx || dy) {
      c.data.moved = true;
      var s = c.elm.style;
      s.transform = s.WebkitTransform = "translate(".concat(dx, "px,").concat(dy, "px)");
      s.transitionDuration = "0s";
    }
  }
  var platformComponents = {
    Transition,
    TransitionGroup
  };
  Vue.config.mustUseProp = mustUseProp;
  Vue.config.isReservedTag = isReservedTag;
  Vue.config.isReservedAttr = isReservedAttr;
  Vue.config.getTagNamespace = getTagNamespace;
  Vue.config.isUnknownElement = isUnknownElement;
  extend(Vue.options.directives, platformDirectives);
  extend(Vue.options.components, platformComponents);
  Vue.prototype.__patch__ = inBrowser ? patch : noop;
  Vue.prototype.$mount = function(el, hydrating) {
    el = el && inBrowser ? query(el) : void 0;
    return mountComponent(this, el, hydrating);
  };
  if (inBrowser) {
    setTimeout(function() {
      if (config.devtools) {
        if (devtools) {
          devtools.emit("init", Vue);
        } else if (true) {
          console[console.info ? "info" : "log"]("Download the Vue Devtools extension for a better development experience:\nhttps://github.com/vuejs/vue-devtools");
        }
      }
      if (config.productionTip !== false && typeof console !== "undefined") {
        console[console.info ? "info" : "log"]("You are running Vue in development mode.\nMake sure to turn on production mode when deploying for production.\nSee more tips at https://vuejs.org/guide/deployment.html");
      }
    }, 0);
  }
  var defaultTagRE = /\{\{((?:.|\r?\n)+?)\}\}/g;
  var regexEscapeRE = /[-.*+?^${}()|[\]\/\\]/g;
  var buildRegex = cached(function(delimiters2) {
    var open = delimiters2[0].replace(regexEscapeRE, "\\$&");
    var close = delimiters2[1].replace(regexEscapeRE, "\\$&");
    return new RegExp(open + "((?:.|\\n)+?)" + close, "g");
  });
  function parseText(text2, delimiters2) {
    var tagRE = delimiters2 ? buildRegex(delimiters2) : defaultTagRE;
    if (!tagRE.test(text2)) {
      return;
    }
    var tokens = [];
    var rawTokens = [];
    var lastIndex = tagRE.lastIndex = 0;
    var match2, index2, tokenValue;
    while (match2 = tagRE.exec(text2)) {
      index2 = match2.index;
      if (index2 > lastIndex) {
        rawTokens.push(tokenValue = text2.slice(lastIndex, index2));
        tokens.push(JSON.stringify(tokenValue));
      }
      var exp = parseFilters(match2[1].trim());
      tokens.push("_s(".concat(exp, ")"));
      rawTokens.push({ "@binding": exp });
      lastIndex = index2 + match2[0].length;
    }
    if (lastIndex < text2.length) {
      rawTokens.push(tokenValue = text2.slice(lastIndex));
      tokens.push(JSON.stringify(tokenValue));
    }
    return {
      expression: tokens.join("+"),
      tokens: rawTokens
    };
  }
  function transformNode$1(el, options) {
    var warn2 = options.warn || baseWarn;
    var staticClass = getAndRemoveAttr(el, "class");
    if (staticClass) {
      var res = parseText(staticClass, options.delimiters);
      if (res) {
        warn2('class="'.concat(staticClass, '": ') + 'Interpolation inside attributes has been removed. Use v-bind or the colon shorthand instead. For example, instead of <div class="{{ val }}">, use <div :class="val">.', el.rawAttrsMap["class"]);
      }
    }
    if (staticClass) {
      el.staticClass = JSON.stringify(staticClass.replace(/\s+/g, " ").trim());
    }
    var classBinding = getBindingAttr(
      el,
      "class",
      false
      /* getStatic */
    );
    if (classBinding) {
      el.classBinding = classBinding;
    }
  }
  function genData$2(el) {
    var data = "";
    if (el.staticClass) {
      data += "staticClass:".concat(el.staticClass, ",");
    }
    if (el.classBinding) {
      data += "class:".concat(el.classBinding, ",");
    }
    return data;
  }
  var klass = {
    staticKeys: ["staticClass"],
    transformNode: transformNode$1,
    genData: genData$2
  };
  function transformNode(el, options) {
    var warn2 = options.warn || baseWarn;
    var staticStyle = getAndRemoveAttr(el, "style");
    if (staticStyle) {
      if (true) {
        var res = parseText(staticStyle, options.delimiters);
        if (res) {
          warn2('style="'.concat(staticStyle, '": ') + 'Interpolation inside attributes has been removed. Use v-bind or the colon shorthand instead. For example, instead of <div style="{{ val }}">, use <div :style="val">.', el.rawAttrsMap["style"]);
        }
      }
      el.staticStyle = JSON.stringify(parseStyleText(staticStyle));
    }
    var styleBinding = getBindingAttr(
      el,
      "style",
      false
      /* getStatic */
    );
    if (styleBinding) {
      el.styleBinding = styleBinding;
    }
  }
  function genData$1(el) {
    var data = "";
    if (el.staticStyle) {
      data += "staticStyle:".concat(el.staticStyle, ",");
    }
    if (el.styleBinding) {
      data += "style:(".concat(el.styleBinding, "),");
    }
    return data;
  }
  var style = {
    staticKeys: ["staticStyle"],
    transformNode,
    genData: genData$1
  };
  var decoder;
  var he = {
    decode: function(html2) {
      decoder = decoder || document.createElement("div");
      decoder.innerHTML = html2;
      return decoder.textContent;
    }
  };
  var isUnaryTag = makeMap("area,base,br,col,embed,frame,hr,img,input,isindex,keygen,link,meta,param,source,track,wbr");
  var canBeLeftOpenTag = makeMap("colgroup,dd,dt,li,options,p,td,tfoot,th,thead,tr,source");
  var isNonPhrasingTag = makeMap("address,article,aside,base,blockquote,body,caption,col,colgroup,dd,details,dialog,div,dl,dt,fieldset,figcaption,figure,footer,form,h1,h2,h3,h4,h5,h6,head,header,hgroup,hr,html,legend,li,menuitem,meta,optgroup,option,param,rp,rt,source,style,summary,tbody,td,tfoot,th,thead,title,tr,track");
  var attribute = /^\s*([^\s"'<>\/=]+)(?:\s*(=)\s*(?:"([^"]*)"+|'([^']*)'+|([^\s"'=<>`]+)))?/;
  var dynamicArgAttribute = /^\s*((?:v-[\w-]+:|@|:|#)\[[^=]+?\][^\s"'<>\/=]*)(?:\s*(=)\s*(?:"([^"]*)"+|'([^']*)'+|([^\s"'=<>`]+)))?/;
  var ncname = "[a-zA-Z_][\\-\\.0-9_a-zA-Z".concat(unicodeRegExp.source, "]*");
  var qnameCapture = "((?:".concat(ncname, "\\:)?").concat(ncname, ")");
  var startTagOpen = new RegExp("^<".concat(qnameCapture));
  var startTagClose = /^\s*(\/?)>/;
  var endTag = new RegExp("^<\\/".concat(qnameCapture, "[^>]*>"));
  var doctype = /^<!DOCTYPE [^>]+>/i;
  var comment = /^<!\--/;
  var conditionalComment = /^<!\[/;
  var isPlainTextElement = makeMap("script,style,textarea", true);
  var reCache = {};
  var decodingMap = {
    "&lt;": "<",
    "&gt;": ">",
    "&quot;": '"',
    "&amp;": "&",
    "&#10;": "\n",
    "&#9;": "	",
    "&#39;": "'"
  };
  var encodedAttr = /&(?:lt|gt|quot|amp|#39);/g;
  var encodedAttrWithNewLines = /&(?:lt|gt|quot|amp|#39|#10|#9);/g;
  var isIgnoreNewlineTag = makeMap("pre,textarea", true);
  var shouldIgnoreFirstNewline = function(tag, html2) {
    return tag && isIgnoreNewlineTag(tag) && html2[0] === "\n";
  };
  function decodeAttr(value, shouldDecodeNewlines2) {
    var re = shouldDecodeNewlines2 ? encodedAttrWithNewLines : encodedAttr;
    return value.replace(re, function(match2) {
      return decodingMap[match2];
    });
  }
  function parseHTML(html2, options) {
    var stack = [];
    var expectHTML = options.expectHTML;
    var isUnaryTag2 = options.isUnaryTag || no;
    var canBeLeftOpenTag2 = options.canBeLeftOpenTag || no;
    var index2 = 0;
    var last, lastTag;
    var _loop_1 = function() {
      last = html2;
      if (!lastTag || !isPlainTextElement(lastTag)) {
        var textEnd = html2.indexOf("<");
        if (textEnd === 0) {
          if (comment.test(html2)) {
            var commentEnd = html2.indexOf("-->");
            if (commentEnd >= 0) {
              if (options.shouldKeepComment && options.comment) {
                options.comment(html2.substring(4, commentEnd), index2, index2 + commentEnd + 3);
              }
              advance(commentEnd + 3);
              return "continue";
            }
          }
          if (conditionalComment.test(html2)) {
            var conditionalEnd = html2.indexOf("]>");
            if (conditionalEnd >= 0) {
              advance(conditionalEnd + 2);
              return "continue";
            }
          }
          var doctypeMatch = html2.match(doctype);
          if (doctypeMatch) {
            advance(doctypeMatch[0].length);
            return "continue";
          }
          var endTagMatch = html2.match(endTag);
          if (endTagMatch) {
            var curIndex = index2;
            advance(endTagMatch[0].length);
            parseEndTag(endTagMatch[1], curIndex, index2);
            return "continue";
          }
          var startTagMatch = parseStartTag();
          if (startTagMatch) {
            handleStartTag(startTagMatch);
            if (shouldIgnoreFirstNewline(startTagMatch.tagName, html2)) {
              advance(1);
            }
            return "continue";
          }
        }
        var text2 = void 0, rest = void 0, next2 = void 0;
        if (textEnd >= 0) {
          rest = html2.slice(textEnd);
          while (!endTag.test(rest) && !startTagOpen.test(rest) && !comment.test(rest) && !conditionalComment.test(rest)) {
            next2 = rest.indexOf("<", 1);
            if (next2 < 0)
              break;
            textEnd += next2;
            rest = html2.slice(textEnd);
          }
          text2 = html2.substring(0, textEnd);
        }
        if (textEnd < 0) {
          text2 = html2;
        }
        if (text2) {
          advance(text2.length);
        }
        if (options.chars && text2) {
          options.chars(text2, index2 - text2.length, index2);
        }
      } else {
        var endTagLength_1 = 0;
        var stackedTag_1 = lastTag.toLowerCase();
        var reStackedTag = reCache[stackedTag_1] || (reCache[stackedTag_1] = new RegExp("([\\s\\S]*?)(</" + stackedTag_1 + "[^>]*>)", "i"));
        var rest = html2.replace(reStackedTag, function(all, text3, endTag2) {
          endTagLength_1 = endTag2.length;
          if (!isPlainTextElement(stackedTag_1) && stackedTag_1 !== "noscript") {
            text3 = text3.replace(/<!\--([\s\S]*?)-->/g, "$1").replace(/<!\[CDATA\[([\s\S]*?)]]>/g, "$1");
          }
          if (shouldIgnoreFirstNewline(stackedTag_1, text3)) {
            text3 = text3.slice(1);
          }
          if (options.chars) {
            options.chars(text3);
          }
          return "";
        });
        index2 += html2.length - rest.length;
        html2 = rest;
        parseEndTag(stackedTag_1, index2 - endTagLength_1, index2);
      }
      if (html2 === last) {
        options.chars && options.chars(html2);
        if (!stack.length && options.warn) {
          options.warn('Mal-formatted tag at end of template: "'.concat(html2, '"'), {
            start: index2 + html2.length
          });
        }
        return "break";
      }
    };
    while (html2) {
      var state_1 = _loop_1();
      if (state_1 === "break")
        break;
    }
    parseEndTag();
    function advance(n) {
      index2 += n;
      html2 = html2.substring(n);
    }
    function parseStartTag() {
      var start = html2.match(startTagOpen);
      if (start) {
        var match2 = {
          tagName: start[1],
          attrs: [],
          start: index2
        };
        advance(start[0].length);
        var end = void 0, attr = void 0;
        while (!(end = html2.match(startTagClose)) && (attr = html2.match(dynamicArgAttribute) || html2.match(attribute))) {
          attr.start = index2;
          advance(attr[0].length);
          attr.end = index2;
          match2.attrs.push(attr);
        }
        if (end) {
          match2.unarySlash = end[1];
          advance(end[0].length);
          match2.end = index2;
          return match2;
        }
      }
    }
    function handleStartTag(match2) {
      var tagName2 = match2.tagName;
      var unarySlash = match2.unarySlash;
      if (expectHTML) {
        if (lastTag === "p" && isNonPhrasingTag(tagName2)) {
          parseEndTag(lastTag);
        }
        if (canBeLeftOpenTag2(tagName2) && lastTag === tagName2) {
          parseEndTag(tagName2);
        }
      }
      var unary = isUnaryTag2(tagName2) || !!unarySlash;
      var l = match2.attrs.length;
      var attrs2 = new Array(l);
      for (var i = 0; i < l; i++) {
        var args = match2.attrs[i];
        var value = args[3] || args[4] || args[5] || "";
        var shouldDecodeNewlines2 = tagName2 === "a" && args[1] === "href" ? options.shouldDecodeNewlinesForHref : options.shouldDecodeNewlines;
        attrs2[i] = {
          name: args[1],
          value: decodeAttr(value, shouldDecodeNewlines2)
        };
        if (options.outputSourceRange) {
          attrs2[i].start = args.start + args[0].match(/^\s*/).length;
          attrs2[i].end = args.end;
        }
      }
      if (!unary) {
        stack.push({
          tag: tagName2,
          lowerCasedTag: tagName2.toLowerCase(),
          attrs: attrs2,
          start: match2.start,
          end: match2.end
        });
        lastTag = tagName2;
      }
      if (options.start) {
        options.start(tagName2, attrs2, unary, match2.start, match2.end);
      }
    }
    function parseEndTag(tagName2, start, end) {
      var pos, lowerCasedTagName;
      if (start == null)
        start = index2;
      if (end == null)
        end = index2;
      if (tagName2) {
        lowerCasedTagName = tagName2.toLowerCase();
        for (pos = stack.length - 1; pos >= 0; pos--) {
          if (stack[pos].lowerCasedTag === lowerCasedTagName) {
            break;
          }
        }
      } else {
        pos = 0;
      }
      if (pos >= 0) {
        for (var i = stack.length - 1; i >= pos; i--) {
          if ((i > pos || !tagName2) && options.warn) {
            options.warn("tag <".concat(stack[i].tag, "> has no matching end tag."), {
              start: stack[i].start,
              end: stack[i].end
            });
          }
          if (options.end) {
            options.end(stack[i].tag, start, end);
          }
        }
        stack.length = pos;
        lastTag = pos && stack[pos - 1].tag;
      } else if (lowerCasedTagName === "br") {
        if (options.start) {
          options.start(tagName2, [], true, start, end);
        }
      } else if (lowerCasedTagName === "p") {
        if (options.start) {
          options.start(tagName2, [], false, start, end);
        }
        if (options.end) {
          options.end(tagName2, start, end);
        }
      }
    }
  }
  var onRE = /^@|^v-on:/;
  var dirRE = /^v-|^@|^:|^#/;
  var forAliasRE = /([\s\S]*?)\s+(?:in|of)\s+([\s\S]*)/;
  var forIteratorRE = /,([^,\}\]]*)(?:,([^,\}\]]*))?$/;
  var stripParensRE = /^\(|\)$/g;
  var dynamicArgRE = /^\[.*\]$/;
  var argRE = /:(.*)$/;
  var bindRE = /^:|^\.|^v-bind:/;
  var modifierRE = /\.[^.\]]+(?=[^\]]*$)/g;
  var slotRE = /^v-slot(:|$)|^#/;
  var lineBreakRE = /[\r\n]/;
  var whitespaceRE = /[ \f\t\r\n]+/g;
  var invalidAttributeRE = /[\s"'<>\/=]/;
  var decodeHTMLCached = cached(he.decode);
  var emptySlotScopeToken = "_empty_";
  var warn;
  var delimiters;
  var transforms;
  var preTransforms;
  var postTransforms;
  var platformIsPreTag;
  var platformMustUseProp;
  var platformGetTagNamespace;
  var maybeComponent;
  function createASTElement(tag, attrs2, parent) {
    return {
      type: 1,
      tag,
      attrsList: attrs2,
      attrsMap: makeAttrsMap(attrs2),
      rawAttrsMap: {},
      parent,
      children: []
    };
  }
  function parse(template, options) {
    warn = options.warn || baseWarn;
    platformIsPreTag = options.isPreTag || no;
    platformMustUseProp = options.mustUseProp || no;
    platformGetTagNamespace = options.getTagNamespace || no;
    var isReservedTag2 = options.isReservedTag || no;
    maybeComponent = function(el) {
      return !!(el.component || el.attrsMap[":is"] || el.attrsMap["v-bind:is"] || !(el.attrsMap.is ? isReservedTag2(el.attrsMap.is) : isReservedTag2(el.tag)));
    };
    transforms = pluckModuleFunction(options.modules, "transformNode");
    preTransforms = pluckModuleFunction(options.modules, "preTransformNode");
    postTransforms = pluckModuleFunction(options.modules, "postTransformNode");
    delimiters = options.delimiters;
    var stack = [];
    var preserveWhitespace = options.preserveWhitespace !== false;
    var whitespaceOption = options.whitespace;
    var root;
    var currentParent;
    var inVPre = false;
    var inPre = false;
    var warned = false;
    function warnOnce(msg, range2) {
      if (!warned) {
        warned = true;
        warn(msg, range2);
      }
    }
    function closeElement(element) {
      trimEndingWhitespace(element);
      if (!inVPre && !element.processed) {
        element = processElement(element, options);
      }
      if (!stack.length && element !== root) {
        if (root.if && (element.elseif || element.else)) {
          if (true) {
            checkRootConstraints(element);
          }
          addIfCondition(root, {
            exp: element.elseif,
            block: element
          });
        } else if (true) {
          warnOnce("Component template should contain exactly one root element. If you are using v-if on multiple elements, use v-else-if to chain them instead.", { start: element.start });
        }
      }
      if (currentParent && !element.forbidden) {
        if (element.elseif || element.else) {
          processIfConditions(element, currentParent);
        } else {
          if (element.slotScope) {
            var name_1 = element.slotTarget || '"default"';
            (currentParent.scopedSlots || (currentParent.scopedSlots = {}))[name_1] = element;
          }
          currentParent.children.push(element);
          element.parent = currentParent;
        }
      }
      element.children = element.children.filter(function(c) {
        return !c.slotScope;
      });
      trimEndingWhitespace(element);
      if (element.pre) {
        inVPre = false;
      }
      if (platformIsPreTag(element.tag)) {
        inPre = false;
      }
      for (var i = 0; i < postTransforms.length; i++) {
        postTransforms[i](element, options);
      }
    }
    function trimEndingWhitespace(el) {
      if (!inPre) {
        var lastNode = void 0;
        while ((lastNode = el.children[el.children.length - 1]) && lastNode.type === 3 && lastNode.text === " ") {
          el.children.pop();
        }
      }
    }
    function checkRootConstraints(el) {
      if (el.tag === "slot" || el.tag === "template") {
        warnOnce("Cannot use <".concat(el.tag, "> as component root element because it may ") + "contain multiple nodes.", { start: el.start });
      }
      if (el.attrsMap.hasOwnProperty("v-for")) {
        warnOnce("Cannot use v-for on stateful component root element because it renders multiple elements.", el.rawAttrsMap["v-for"]);
      }
    }
    parseHTML(template, {
      warn,
      expectHTML: options.expectHTML,
      isUnaryTag: options.isUnaryTag,
      canBeLeftOpenTag: options.canBeLeftOpenTag,
      shouldDecodeNewlines: options.shouldDecodeNewlines,
      shouldDecodeNewlinesForHref: options.shouldDecodeNewlinesForHref,
      shouldKeepComment: options.comments,
      outputSourceRange: options.outputSourceRange,
      start: function(tag, attrs2, unary, start, end) {
        var ns = currentParent && currentParent.ns || platformGetTagNamespace(tag);
        if (isIE && ns === "svg") {
          attrs2 = guardIESVGBug(attrs2);
        }
        var element = createASTElement(tag, attrs2, currentParent);
        if (ns) {
          element.ns = ns;
        }
        if (true) {
          if (options.outputSourceRange) {
            element.start = start;
            element.end = end;
            element.rawAttrsMap = element.attrsList.reduce(function(cumulated, attr) {
              cumulated[attr.name] = attr;
              return cumulated;
            }, {});
          }
          attrs2.forEach(function(attr) {
            if (invalidAttributeRE.test(attr.name)) {
              warn("Invalid dynamic argument expression: attribute names cannot contain spaces, quotes, <, >, / or =.", options.outputSourceRange ? {
                start: attr.start + attr.name.indexOf("["),
                end: attr.start + attr.name.length
              } : void 0);
            }
          });
        }
        if (isForbiddenTag(element) && !isServerRendering()) {
          element.forbidden = true;
          warn("Templates should only be responsible for mapping the state to the UI. Avoid placing tags with side-effects in your templates, such as " + "<".concat(tag, ">") + ", as they will not be parsed.", { start: element.start });
        }
        for (var i = 0; i < preTransforms.length; i++) {
          element = preTransforms[i](element, options) || element;
        }
        if (!inVPre) {
          processPre(element);
          if (element.pre) {
            inVPre = true;
          }
        }
        if (platformIsPreTag(element.tag)) {
          inPre = true;
        }
        if (inVPre) {
          processRawAttrs(element);
        } else if (!element.processed) {
          processFor(element);
          processIf(element);
          processOnce(element);
        }
        if (!root) {
          root = element;
          if (true) {
            checkRootConstraints(root);
          }
        }
        if (!unary) {
          currentParent = element;
          stack.push(element);
        } else {
          closeElement(element);
        }
      },
      end: function(tag, start, end) {
        var element = stack[stack.length - 1];
        stack.length -= 1;
        currentParent = stack[stack.length - 1];
        if (options.outputSourceRange) {
          element.end = end;
        }
        closeElement(element);
      },
      chars: function(text2, start, end) {
        if (!currentParent) {
          if (true) {
            if (text2 === template) {
              warnOnce("Component template requires a root element, rather than just text.", { start });
            } else if (text2 = text2.trim()) {
              warnOnce('text "'.concat(text2, '" outside root element will be ignored.'), {
                start
              });
            }
          }
          return;
        }
        if (isIE && currentParent.tag === "textarea" && currentParent.attrsMap.placeholder === text2) {
          return;
        }
        var children = currentParent.children;
        if (inPre || text2.trim()) {
          text2 = isTextTag(currentParent) ? text2 : decodeHTMLCached(text2);
        } else if (!children.length) {
          text2 = "";
        } else if (whitespaceOption) {
          if (whitespaceOption === "condense") {
            text2 = lineBreakRE.test(text2) ? "" : " ";
          } else {
            text2 = " ";
          }
        } else {
          text2 = preserveWhitespace ? " " : "";
        }
        if (text2) {
          if (!inPre && whitespaceOption === "condense") {
            text2 = text2.replace(whitespaceRE, " ");
          }
          var res = void 0;
          var child = void 0;
          if (!inVPre && text2 !== " " && (res = parseText(text2, delimiters))) {
            child = {
              type: 2,
              expression: res.expression,
              tokens: res.tokens,
              text: text2
            };
          } else if (text2 !== " " || !children.length || children[children.length - 1].text !== " ") {
            child = {
              type: 3,
              text: text2
            };
          }
          if (child) {
            if (options.outputSourceRange) {
              child.start = start;
              child.end = end;
            }
            children.push(child);
          }
        }
      },
      comment: function(text2, start, end) {
        if (currentParent) {
          var child = {
            type: 3,
            text: text2,
            isComment: true
          };
          if (options.outputSourceRange) {
            child.start = start;
            child.end = end;
          }
          currentParent.children.push(child);
        }
      }
    });
    return root;
  }
  function processPre(el) {
    if (getAndRemoveAttr(el, "v-pre") != null) {
      el.pre = true;
    }
  }
  function processRawAttrs(el) {
    var list = el.attrsList;
    var len2 = list.length;
    if (len2) {
      var attrs2 = el.attrs = new Array(len2);
      for (var i = 0; i < len2; i++) {
        attrs2[i] = {
          name: list[i].name,
          value: JSON.stringify(list[i].value)
        };
        if (list[i].start != null) {
          attrs2[i].start = list[i].start;
          attrs2[i].end = list[i].end;
        }
      }
    } else if (!el.pre) {
      el.plain = true;
    }
  }
  function processElement(element, options) {
    processKey(element);
    element.plain = !element.key && !element.scopedSlots && !element.attrsList.length;
    processRef(element);
    processSlotContent(element);
    processSlotOutlet(element);
    processComponent(element);
    for (var i = 0; i < transforms.length; i++) {
      element = transforms[i](element, options) || element;
    }
    processAttrs(element);
    return element;
  }
  function processKey(el) {
    var exp = getBindingAttr(el, "key");
    if (exp) {
      if (true) {
        if (el.tag === "template") {
          warn("<template> cannot be keyed. Place the key on real elements instead.", getRawBindingAttr(el, "key"));
        }
        if (el.for) {
          var iterator = el.iterator2 || el.iterator1;
          var parent_1 = el.parent;
          if (iterator && iterator === exp && parent_1 && parent_1.tag === "transition-group") {
            warn(
              "Do not use v-for index as key on <transition-group> children, this is the same as not using keys.",
              getRawBindingAttr(el, "key"),
              true
              /* tip */
            );
          }
        }
      }
      el.key = exp;
    }
  }
  function processRef(el) {
    var ref2 = getBindingAttr(el, "ref");
    if (ref2) {
      el.ref = ref2;
      el.refInFor = checkInFor(el);
    }
  }
  function processFor(el) {
    var exp;
    if (exp = getAndRemoveAttr(el, "v-for")) {
      var res = parseFor(exp);
      if (res) {
        extend(el, res);
      } else if (true) {
        warn("Invalid v-for expression: ".concat(exp), el.rawAttrsMap["v-for"]);
      }
    }
  }
  function parseFor(exp) {
    var inMatch = exp.match(forAliasRE);
    if (!inMatch)
      return;
    var res = {};
    res.for = inMatch[2].trim();
    var alias = inMatch[1].trim().replace(stripParensRE, "");
    var iteratorMatch = alias.match(forIteratorRE);
    if (iteratorMatch) {
      res.alias = alias.replace(forIteratorRE, "").trim();
      res.iterator1 = iteratorMatch[1].trim();
      if (iteratorMatch[2]) {
        res.iterator2 = iteratorMatch[2].trim();
      }
    } else {
      res.alias = alias;
    }
    return res;
  }
  function processIf(el) {
    var exp = getAndRemoveAttr(el, "v-if");
    if (exp) {
      el.if = exp;
      addIfCondition(el, {
        exp,
        block: el
      });
    } else {
      if (getAndRemoveAttr(el, "v-else") != null) {
        el.else = true;
      }
      var elseif = getAndRemoveAttr(el, "v-else-if");
      if (elseif) {
        el.elseif = elseif;
      }
    }
  }
  function processIfConditions(el, parent) {
    var prev = findPrevElement(parent.children);
    if (prev && prev.if) {
      addIfCondition(prev, {
        exp: el.elseif,
        block: el
      });
    } else if (true) {
      warn("v-".concat(el.elseif ? 'else-if="' + el.elseif + '"' : "else", " ") + "used on element <".concat(el.tag, "> without corresponding v-if."), el.rawAttrsMap[el.elseif ? "v-else-if" : "v-else"]);
    }
  }
  function findPrevElement(children) {
    var i = children.length;
    while (i--) {
      if (children[i].type === 1) {
        return children[i];
      } else {
        if (children[i].text !== " ") {
          warn('text "'.concat(children[i].text.trim(), '" between v-if and v-else(-if) ') + "will be ignored.", children[i]);
        }
        children.pop();
      }
    }
  }
  function addIfCondition(el, condition) {
    if (!el.ifConditions) {
      el.ifConditions = [];
    }
    el.ifConditions.push(condition);
  }
  function processOnce(el) {
    var once2 = getAndRemoveAttr(el, "v-once");
    if (once2 != null) {
      el.once = true;
    }
  }
  function processSlotContent(el) {
    var slotScope;
    if (el.tag === "template") {
      slotScope = getAndRemoveAttr(el, "scope");
      if (slotScope) {
        warn('the "scope" attribute for scoped slots have been deprecated and replaced by "slot-scope" since 2.5. The new "slot-scope" attribute can also be used on plain elements in addition to <template> to denote scoped slots.', el.rawAttrsMap["scope"], true);
      }
      el.slotScope = slotScope || getAndRemoveAttr(el, "slot-scope");
    } else if (slotScope = getAndRemoveAttr(el, "slot-scope")) {
      if (el.attrsMap["v-for"]) {
        warn("Ambiguous combined usage of slot-scope and v-for on <".concat(el.tag, "> ") + "(v-for takes higher priority). Use a wrapper <template> for the scoped slot to make it clearer.", el.rawAttrsMap["slot-scope"], true);
      }
      el.slotScope = slotScope;
    }
    var slotTarget = getBindingAttr(el, "slot");
    if (slotTarget) {
      el.slotTarget = slotTarget === '""' ? '"default"' : slotTarget;
      el.slotTargetDynamic = !!(el.attrsMap[":slot"] || el.attrsMap["v-bind:slot"]);
      if (el.tag !== "template" && !el.slotScope) {
        addAttr(el, "slot", slotTarget, getRawBindingAttr(el, "slot"));
      }
    }
    {
      if (el.tag === "template") {
        var slotBinding = getAndRemoveAttrByRegex(el, slotRE);
        if (slotBinding) {
          if (true) {
            if (el.slotTarget || el.slotScope) {
              warn("Unexpected mixed usage of different slot syntaxes.", el);
            }
            if (el.parent && !maybeComponent(el.parent)) {
              warn("<template v-slot> can only appear at the root level inside the receiving component", el);
            }
          }
          var _a2 = getSlotName(slotBinding), name_2 = _a2.name, dynamic = _a2.dynamic;
          el.slotTarget = name_2;
          el.slotTargetDynamic = dynamic;
          el.slotScope = slotBinding.value || emptySlotScopeToken;
        }
      } else {
        var slotBinding = getAndRemoveAttrByRegex(el, slotRE);
        if (slotBinding) {
          if (true) {
            if (!maybeComponent(el)) {
              warn("v-slot can only be used on components or <template>.", slotBinding);
            }
            if (el.slotScope || el.slotTarget) {
              warn("Unexpected mixed usage of different slot syntaxes.", el);
            }
            if (el.scopedSlots) {
              warn("To avoid scope ambiguity, the default slot should also use <template> syntax when there are other named slots.", slotBinding);
            }
          }
          var slots = el.scopedSlots || (el.scopedSlots = {});
          var _b = getSlotName(slotBinding), name_3 = _b.name, dynamic = _b.dynamic;
          var slotContainer_1 = slots[name_3] = createASTElement("template", [], el);
          slotContainer_1.slotTarget = name_3;
          slotContainer_1.slotTargetDynamic = dynamic;
          slotContainer_1.children = el.children.filter(function(c) {
            if (!c.slotScope) {
              c.parent = slotContainer_1;
              return true;
            }
          });
          slotContainer_1.slotScope = slotBinding.value || emptySlotScopeToken;
          el.children = [];
          el.plain = false;
        }
      }
    }
  }
  function getSlotName(binding) {
    var name = binding.name.replace(slotRE, "");
    if (!name) {
      if (binding.name[0] !== "#") {
        name = "default";
      } else if (true) {
        warn("v-slot shorthand syntax requires a slot name.", binding);
      }
    }
    return dynamicArgRE.test(name) ? (
      // dynamic [name]
      { name: name.slice(1, -1), dynamic: true }
    ) : (
      // static name
      { name: '"'.concat(name, '"'), dynamic: false }
    );
  }
  function processSlotOutlet(el) {
    if (el.tag === "slot") {
      el.slotName = getBindingAttr(el, "name");
      if (el.key) {
        warn("`key` does not work on <slot> because slots are abstract outlets and can possibly expand into multiple elements. Use the key on a wrapping element instead.", getRawBindingAttr(el, "key"));
      }
    }
  }
  function processComponent(el) {
    var binding;
    if (binding = getBindingAttr(el, "is")) {
      el.component = binding;
    }
    if (getAndRemoveAttr(el, "inline-template") != null) {
      el.inlineTemplate = true;
    }
  }
  function processAttrs(el) {
    var list = el.attrsList;
    var i, l, name, rawName, value, modifiers, syncGen, isDynamic;
    for (i = 0, l = list.length; i < l; i++) {
      name = rawName = list[i].name;
      value = list[i].value;
      if (dirRE.test(name)) {
        el.hasBindings = true;
        modifiers = parseModifiers(name.replace(dirRE, ""));
        if (modifiers) {
          name = name.replace(modifierRE, "");
        }
        if (bindRE.test(name)) {
          name = name.replace(bindRE, "");
          value = parseFilters(value);
          isDynamic = dynamicArgRE.test(name);
          if (isDynamic) {
            name = name.slice(1, -1);
          }
          if (value.trim().length === 0) {
            warn('The value for a v-bind expression cannot be empty. Found in "v-bind:'.concat(name, '"'));
          }
          if (modifiers) {
            if (modifiers.prop && !isDynamic) {
              name = camelize(name);
              if (name === "innerHtml")
                name = "innerHTML";
            }
            if (modifiers.camel && !isDynamic) {
              name = camelize(name);
            }
            if (modifiers.sync) {
              syncGen = genAssignmentCode(value, "$event");
              if (!isDynamic) {
                addHandler(el, "update:".concat(camelize(name)), syncGen, null, false, warn, list[i]);
                if (hyphenate(name) !== camelize(name)) {
                  addHandler(el, "update:".concat(hyphenate(name)), syncGen, null, false, warn, list[i]);
                }
              } else {
                addHandler(
                  el,
                  '"update:"+('.concat(name, ")"),
                  syncGen,
                  null,
                  false,
                  warn,
                  list[i],
                  true
                  // dynamic
                );
              }
            }
          }
          if (modifiers && modifiers.prop || !el.component && platformMustUseProp(el.tag, el.attrsMap.type, name)) {
            addProp(el, name, value, list[i], isDynamic);
          } else {
            addAttr(el, name, value, list[i], isDynamic);
          }
        } else if (onRE.test(name)) {
          name = name.replace(onRE, "");
          isDynamic = dynamicArgRE.test(name);
          if (isDynamic) {
            name = name.slice(1, -1);
          }
          addHandler(el, name, value, modifiers, false, warn, list[i], isDynamic);
        } else {
          name = name.replace(dirRE, "");
          var argMatch = name.match(argRE);
          var arg = argMatch && argMatch[1];
          isDynamic = false;
          if (arg) {
            name = name.slice(0, -(arg.length + 1));
            if (dynamicArgRE.test(arg)) {
              arg = arg.slice(1, -1);
              isDynamic = true;
            }
          }
          addDirective(el, name, rawName, value, arg, isDynamic, modifiers, list[i]);
          if (name === "model") {
            checkForAliasModel(el, value);
          }
        }
      } else {
        if (true) {
          var res = parseText(value, delimiters);
          if (res) {
            warn("".concat(name, '="').concat(value, '": ') + 'Interpolation inside attributes has been removed. Use v-bind or the colon shorthand instead. For example, instead of <div id="{{ val }}">, use <div :id="val">.', list[i]);
          }
        }
        addAttr(el, name, JSON.stringify(value), list[i]);
        if (!el.component && name === "muted" && platformMustUseProp(el.tag, el.attrsMap.type, name)) {
          addProp(el, name, "true", list[i]);
        }
      }
    }
  }
  function checkInFor(el) {
    var parent = el;
    while (parent) {
      if (parent.for !== void 0) {
        return true;
      }
      parent = parent.parent;
    }
    return false;
  }
  function parseModifiers(name) {
    var match2 = name.match(modifierRE);
    if (match2) {
      var ret_1 = {};
      match2.forEach(function(m) {
        ret_1[m.slice(1)] = true;
      });
      return ret_1;
    }
  }
  function makeAttrsMap(attrs2) {
    var map = {};
    for (var i = 0, l = attrs2.length; i < l; i++) {
      if (map[attrs2[i].name] && !isIE && !isEdge) {
        warn("duplicate attribute: " + attrs2[i].name, attrs2[i]);
      }
      map[attrs2[i].name] = attrs2[i].value;
    }
    return map;
  }
  function isTextTag(el) {
    return el.tag === "script" || el.tag === "style";
  }
  function isForbiddenTag(el) {
    return el.tag === "style" || el.tag === "script" && (!el.attrsMap.type || el.attrsMap.type === "text/javascript");
  }
  var ieNSBug = /^xmlns:NS\d+/;
  var ieNSPrefix = /^NS\d+:/;
  function guardIESVGBug(attrs2) {
    var res = [];
    for (var i = 0; i < attrs2.length; i++) {
      var attr = attrs2[i];
      if (!ieNSBug.test(attr.name)) {
        attr.name = attr.name.replace(ieNSPrefix, "");
        res.push(attr);
      }
    }
    return res;
  }
  function checkForAliasModel(el, value) {
    var _el = el;
    while (_el) {
      if (_el.for && _el.alias === value) {
        warn("<".concat(el.tag, ' v-model="').concat(value, '">: ') + "You are binding v-model directly to a v-for iteration alias. This will not be able to modify the v-for source array because writing to the alias is like modifying a function local variable. Consider using an array of objects and use v-model on an object property instead.", el.rawAttrsMap["v-model"]);
      }
      _el = _el.parent;
    }
  }
  function preTransformNode(el, options) {
    if (el.tag === "input") {
      var map = el.attrsMap;
      if (!map["v-model"]) {
        return;
      }
      var typeBinding = void 0;
      if (map[":type"] || map["v-bind:type"]) {
        typeBinding = getBindingAttr(el, "type");
      }
      if (!map.type && !typeBinding && map["v-bind"]) {
        typeBinding = "(".concat(map["v-bind"], ").type");
      }
      if (typeBinding) {
        var ifCondition = getAndRemoveAttr(el, "v-if", true);
        var ifConditionExtra = ifCondition ? "&&(".concat(ifCondition, ")") : "";
        var hasElse = getAndRemoveAttr(el, "v-else", true) != null;
        var elseIfCondition = getAndRemoveAttr(el, "v-else-if", true);
        var branch0 = cloneASTElement(el);
        processFor(branch0);
        addRawAttr(branch0, "type", "checkbox");
        processElement(branch0, options);
        branch0.processed = true;
        branch0.if = "(".concat(typeBinding, ")==='checkbox'") + ifConditionExtra;
        addIfCondition(branch0, {
          exp: branch0.if,
          block: branch0
        });
        var branch1 = cloneASTElement(el);
        getAndRemoveAttr(branch1, "v-for", true);
        addRawAttr(branch1, "type", "radio");
        processElement(branch1, options);
        addIfCondition(branch0, {
          exp: "(".concat(typeBinding, ")==='radio'") + ifConditionExtra,
          block: branch1
        });
        var branch2 = cloneASTElement(el);
        getAndRemoveAttr(branch2, "v-for", true);
        addRawAttr(branch2, ":type", typeBinding);
        processElement(branch2, options);
        addIfCondition(branch0, {
          exp: ifCondition,
          block: branch2
        });
        if (hasElse) {
          branch0.else = true;
        } else if (elseIfCondition) {
          branch0.elseif = elseIfCondition;
        }
        return branch0;
      }
    }
  }
  function cloneASTElement(el) {
    return createASTElement(el.tag, el.attrsList.slice(), el.parent);
  }
  var model = {
    preTransformNode
  };
  var modules = [klass, style, model];
  function text(el, dir) {
    if (dir.value) {
      addProp(el, "textContent", "_s(".concat(dir.value, ")"), dir);
    }
  }
  function html(el, dir) {
    if (dir.value) {
      addProp(el, "innerHTML", "_s(".concat(dir.value, ")"), dir);
    }
  }
  var directives = {
    model: model$1,
    text,
    html
  };
  var baseOptions = {
    expectHTML: true,
    modules,
    directives,
    isPreTag,
    isUnaryTag,
    mustUseProp,
    canBeLeftOpenTag,
    isReservedTag,
    getTagNamespace,
    staticKeys: genStaticKeys$1(modules)
  };
  var isStaticKey;
  var isPlatformReservedTag;
  var genStaticKeysCached = cached(genStaticKeys);
  function optimize(root, options) {
    if (!root)
      return;
    isStaticKey = genStaticKeysCached(options.staticKeys || "");
    isPlatformReservedTag = options.isReservedTag || no;
    markStatic(root);
    markStaticRoots(root, false);
  }
  function genStaticKeys(keys) {
    return makeMap("type,tag,attrsList,attrsMap,plain,parent,children,attrs,start,end,rawAttrsMap" + (keys ? "," + keys : ""));
  }
  function markStatic(node) {
    node.static = isStatic(node);
    if (node.type === 1) {
      if (!isPlatformReservedTag(node.tag) && node.tag !== "slot" && node.attrsMap["inline-template"] == null) {
        return;
      }
      for (var i = 0, l = node.children.length; i < l; i++) {
        var child = node.children[i];
        markStatic(child);
        if (!child.static) {
          node.static = false;
        }
      }
      if (node.ifConditions) {
        for (var i = 1, l = node.ifConditions.length; i < l; i++) {
          var block = node.ifConditions[i].block;
          markStatic(block);
          if (!block.static) {
            node.static = false;
          }
        }
      }
    }
  }
  function markStaticRoots(node, isInFor) {
    if (node.type === 1) {
      if (node.static || node.once) {
        node.staticInFor = isInFor;
      }
      if (node.static && node.children.length && !(node.children.length === 1 && node.children[0].type === 3)) {
        node.staticRoot = true;
        return;
      } else {
        node.staticRoot = false;
      }
      if (node.children) {
        for (var i = 0, l = node.children.length; i < l; i++) {
          markStaticRoots(node.children[i], isInFor || !!node.for);
        }
      }
      if (node.ifConditions) {
        for (var i = 1, l = node.ifConditions.length; i < l; i++) {
          markStaticRoots(node.ifConditions[i].block, isInFor);
        }
      }
    }
  }
  function isStatic(node) {
    if (node.type === 2) {
      return false;
    }
    if (node.type === 3) {
      return true;
    }
    return !!(node.pre || !node.hasBindings && // no dynamic bindings
    !node.if && !node.for && // not v-if or v-for or v-else
    !isBuiltInTag(node.tag) && // not a built-in
    isPlatformReservedTag(node.tag) && // not a component
    !isDirectChildOfTemplateFor(node) && Object.keys(node).every(isStaticKey));
  }
  function isDirectChildOfTemplateFor(node) {
    while (node.parent) {
      node = node.parent;
      if (node.tag !== "template") {
        return false;
      }
      if (node.for) {
        return true;
      }
    }
    return false;
  }
  var fnExpRE = /^([\w$_]+|\([^)]*?\))\s*=>|^function(?:\s+[\w$]+)?\s*\(/;
  var fnInvokeRE = /\([^)]*?\);*$/;
  var simplePathRE = /^[A-Za-z_$][\w$]*(?:\.[A-Za-z_$][\w$]*|\['[^']*?']|\["[^"]*?"]|\[\d+]|\[[A-Za-z_$][\w$]*])*$/;
  var keyCodes = {
    esc: 27,
    tab: 9,
    enter: 13,
    space: 32,
    up: 38,
    left: 37,
    right: 39,
    down: 40,
    delete: [8, 46]
  };
  var keyNames = {
    // #7880: IE11 and Edge use `Esc` for Escape key name.
    esc: ["Esc", "Escape"],
    tab: "Tab",
    enter: "Enter",
    // #9112: IE11 uses `Spacebar` for Space key name.
    space: [" ", "Spacebar"],
    // #7806: IE11 uses key names without `Arrow` prefix for arrow keys.
    up: ["Up", "ArrowUp"],
    left: ["Left", "ArrowLeft"],
    right: ["Right", "ArrowRight"],
    down: ["Down", "ArrowDown"],
    // #9112: IE11 uses `Del` for Delete key name.
    delete: ["Backspace", "Delete", "Del"]
  };
  var genGuard = function(condition) {
    return "if(".concat(condition, ")return null;");
  };
  var modifierCode = {
    stop: "$event.stopPropagation();",
    prevent: "$event.preventDefault();",
    self: genGuard("$event.target !== $event.currentTarget"),
    ctrl: genGuard("!$event.ctrlKey"),
    shift: genGuard("!$event.shiftKey"),
    alt: genGuard("!$event.altKey"),
    meta: genGuard("!$event.metaKey"),
    left: genGuard("'button' in $event && $event.button !== 0"),
    middle: genGuard("'button' in $event && $event.button !== 1"),
    right: genGuard("'button' in $event && $event.button !== 2")
  };
  function genHandlers(events2, isNative2) {
    var prefix = isNative2 ? "nativeOn:" : "on:";
    var staticHandlers = "";
    var dynamicHandlers = "";
    for (var name_1 in events2) {
      var handlerCode = genHandler(events2[name_1]);
      if (events2[name_1] && events2[name_1].dynamic) {
        dynamicHandlers += "".concat(name_1, ",").concat(handlerCode, ",");
      } else {
        staticHandlers += '"'.concat(name_1, '":').concat(handlerCode, ",");
      }
    }
    staticHandlers = "{".concat(staticHandlers.slice(0, -1), "}");
    if (dynamicHandlers) {
      return prefix + "_d(".concat(staticHandlers, ",[").concat(dynamicHandlers.slice(0, -1), "])");
    } else {
      return prefix + staticHandlers;
    }
  }
  function genHandler(handler) {
    if (!handler) {
      return "function(){}";
    }
    if (Array.isArray(handler)) {
      return "[".concat(handler.map(function(handler2) {
        return genHandler(handler2);
      }).join(","), "]");
    }
    var isMethodPath = simplePathRE.test(handler.value);
    var isFunctionExpression = fnExpRE.test(handler.value);
    var isFunctionInvocation = simplePathRE.test(handler.value.replace(fnInvokeRE, ""));
    if (!handler.modifiers) {
      if (isMethodPath || isFunctionExpression) {
        return handler.value;
      }
      return "function($event){".concat(isFunctionInvocation ? "return ".concat(handler.value) : handler.value, "}");
    } else {
      var code = "";
      var genModifierCode = "";
      var keys = [];
      var _loop_1 = function(key2) {
        if (modifierCode[key2]) {
          genModifierCode += modifierCode[key2];
          if (keyCodes[key2]) {
            keys.push(key2);
          }
        } else if (key2 === "exact") {
          var modifiers_1 = handler.modifiers;
          genModifierCode += genGuard(["ctrl", "shift", "alt", "meta"].filter(function(keyModifier) {
            return !modifiers_1[keyModifier];
          }).map(function(keyModifier) {
            return "$event.".concat(keyModifier, "Key");
          }).join("||"));
        } else {
          keys.push(key2);
        }
      };
      for (var key in handler.modifiers) {
        _loop_1(key);
      }
      if (keys.length) {
        code += genKeyFilter(keys);
      }
      if (genModifierCode) {
        code += genModifierCode;
      }
      var handlerCode = isMethodPath ? "return ".concat(handler.value, ".apply(null, arguments)") : isFunctionExpression ? "return (".concat(handler.value, ").apply(null, arguments)") : isFunctionInvocation ? "return ".concat(handler.value) : handler.value;
      return "function($event){".concat(code).concat(handlerCode, "}");
    }
  }
  function genKeyFilter(keys) {
    return (
      // make sure the key filters only apply to KeyboardEvents
      // #9441: can't use 'keyCode' in $event because Chrome autofill fires fake
      // key events that do not have keyCode property...
      "if(!$event.type.indexOf('key')&&" + "".concat(keys.map(genFilterCode).join("&&"), ")return null;")
    );
  }
  function genFilterCode(key) {
    var keyVal = parseInt(key, 10);
    if (keyVal) {
      return "$event.keyCode!==".concat(keyVal);
    }
    var keyCode = keyCodes[key];
    var keyName = keyNames[key];
    return "_k($event.keyCode," + "".concat(JSON.stringify(key), ",") + "".concat(JSON.stringify(keyCode), ",") + "$event.key," + "".concat(JSON.stringify(keyName)) + ")";
  }
  function on(el, dir) {
    if (dir.modifiers) {
      warn$2("v-on without argument does not support modifiers.");
    }
    el.wrapListeners = function(code) {
      return "_g(".concat(code, ",").concat(dir.value, ")");
    };
  }
  function bind(el, dir) {
    el.wrapData = function(code) {
      return "_b(".concat(code, ",'").concat(el.tag, "',").concat(dir.value, ",").concat(dir.modifiers && dir.modifiers.prop ? "true" : "false").concat(dir.modifiers && dir.modifiers.sync ? ",true" : "", ")");
    };
  }
  var baseDirectives = {
    on,
    bind,
    cloak: noop
  };
  var CodegenState = (
    /** @class */
    /* @__PURE__ */ (function() {
      function CodegenState2(options) {
        this.options = options;
        this.warn = options.warn || baseWarn;
        this.transforms = pluckModuleFunction(options.modules, "transformCode");
        this.dataGenFns = pluckModuleFunction(options.modules, "genData");
        this.directives = extend(extend({}, baseDirectives), options.directives);
        var isReservedTag2 = options.isReservedTag || no;
        this.maybeComponent = function(el) {
          return !!el.component || !isReservedTag2(el.tag);
        };
        this.onceId = 0;
        this.staticRenderFns = [];
        this.pre = false;
      }
      return CodegenState2;
    })()
  );
  function generate(ast, options) {
    var state = new CodegenState(options);
    var code = ast ? ast.tag === "script" ? "null" : genElement(ast, state) : '_c("div")';
    return {
      render: "with(this){return ".concat(code, "}"),
      staticRenderFns: state.staticRenderFns
    };
  }
  function genElement(el, state) {
    if (el.parent) {
      el.pre = el.pre || el.parent.pre;
    }
    if (el.staticRoot && !el.staticProcessed) {
      return genStatic(el, state);
    } else if (el.once && !el.onceProcessed) {
      return genOnce(el, state);
    } else if (el.for && !el.forProcessed) {
      return genFor(el, state);
    } else if (el.if && !el.ifProcessed) {
      return genIf(el, state);
    } else if (el.tag === "template" && !el.slotTarget && !state.pre) {
      return genChildren(el, state) || "void 0";
    } else if (el.tag === "slot") {
      return genSlot(el, state);
    } else {
      var code = void 0;
      if (el.component) {
        code = genComponent(el.component, el, state);
      } else {
        var data = void 0;
        var maybeComponent2 = state.maybeComponent(el);
        if (!el.plain || el.pre && maybeComponent2) {
          data = genData(el, state);
        }
        var tag = void 0;
        var bindings = state.options.bindings;
        if (maybeComponent2 && bindings && bindings.__isScriptSetup !== false) {
          tag = checkBindingType(bindings, el.tag);
        }
        if (!tag)
          tag = "'".concat(el.tag, "'");
        var children = el.inlineTemplate ? null : genChildren(el, state, true);
        code = "_c(".concat(tag).concat(
          data ? ",".concat(data) : ""
          // data
        ).concat(
          children ? ",".concat(children) : "",
          ")"
        );
      }
      for (var i = 0; i < state.transforms.length; i++) {
        code = state.transforms[i](el, code);
      }
      return code;
    }
  }
  function checkBindingType(bindings, key) {
    var camelName = camelize(key);
    var PascalName = capitalize(camelName);
    var checkType = function(type) {
      if (bindings[key] === type) {
        return key;
      }
      if (bindings[camelName] === type) {
        return camelName;
      }
      if (bindings[PascalName] === type) {
        return PascalName;
      }
    };
    var fromConst = checkType(
      "setup-const"
      /* BindingTypes.SETUP_CONST */
    ) || checkType(
      "setup-reactive-const"
      /* BindingTypes.SETUP_REACTIVE_CONST */
    );
    if (fromConst) {
      return fromConst;
    }
    var fromMaybeRef = checkType(
      "setup-let"
      /* BindingTypes.SETUP_LET */
    ) || checkType(
      "setup-ref"
      /* BindingTypes.SETUP_REF */
    ) || checkType(
      "setup-maybe-ref"
      /* BindingTypes.SETUP_MAYBE_REF */
    );
    if (fromMaybeRef) {
      return fromMaybeRef;
    }
  }
  function genStatic(el, state) {
    el.staticProcessed = true;
    var originalPreState = state.pre;
    if (el.pre) {
      state.pre = el.pre;
    }
    state.staticRenderFns.push("with(this){return ".concat(genElement(el, state), "}"));
    state.pre = originalPreState;
    return "_m(".concat(state.staticRenderFns.length - 1).concat(el.staticInFor ? ",true" : "", ")");
  }
  function genOnce(el, state) {
    el.onceProcessed = true;
    if (el.if && !el.ifProcessed) {
      return genIf(el, state);
    } else if (el.staticInFor) {
      var key = "";
      var parent_1 = el.parent;
      while (parent_1) {
        if (parent_1.for) {
          key = parent_1.key;
          break;
        }
        parent_1 = parent_1.parent;
      }
      if (!key) {
        state.warn("v-once can only be used inside v-for that is keyed. ", el.rawAttrsMap["v-once"]);
        return genElement(el, state);
      }
      return "_o(".concat(genElement(el, state), ",").concat(state.onceId++, ",").concat(key, ")");
    } else {
      return genStatic(el, state);
    }
  }
  function genIf(el, state, altGen, altEmpty) {
    el.ifProcessed = true;
    return genIfConditions(el.ifConditions.slice(), state, altGen, altEmpty);
  }
  function genIfConditions(conditions, state, altGen, altEmpty) {
    if (!conditions.length) {
      return altEmpty || "_e()";
    }
    var condition = conditions.shift();
    if (condition.exp) {
      return "(".concat(condition.exp, ")?").concat(genTernaryExp(condition.block), ":").concat(genIfConditions(conditions, state, altGen, altEmpty));
    } else {
      return "".concat(genTernaryExp(condition.block));
    }
    function genTernaryExp(el) {
      return altGen ? altGen(el, state) : el.once ? genOnce(el, state) : genElement(el, state);
    }
  }
  function genFor(el, state, altGen, altHelper) {
    var exp = el.for;
    var alias = el.alias;
    var iterator1 = el.iterator1 ? ",".concat(el.iterator1) : "";
    var iterator2 = el.iterator2 ? ",".concat(el.iterator2) : "";
    if (state.maybeComponent(el) && el.tag !== "slot" && el.tag !== "template" && !el.key) {
      state.warn(
        "<".concat(el.tag, ' v-for="').concat(alias, " in ").concat(exp, '">: component lists rendered with ') + "v-for should have explicit keys. See https://v2.vuejs.org/v2/guide/list.html#key for more info.",
        el.rawAttrsMap["v-for"],
        true
        /* tip */
      );
    }
    el.forProcessed = true;
    return "".concat(altHelper || "_l", "((").concat(exp, "),") + "function(".concat(alias).concat(iterator1).concat(iterator2, "){") + "return ".concat((altGen || genElement)(el, state)) + "})";
  }
  function genData(el, state) {
    var data = "{";
    var dirs = genDirectives(el, state);
    if (dirs)
      data += dirs + ",";
    if (el.key) {
      data += "key:".concat(el.key, ",");
    }
    if (el.ref) {
      data += "ref:".concat(el.ref, ",");
    }
    if (el.refInFor) {
      data += "refInFor:true,";
    }
    if (el.pre) {
      data += "pre:true,";
    }
    if (el.component) {
      data += 'tag:"'.concat(el.tag, '",');
    }
    for (var i = 0; i < state.dataGenFns.length; i++) {
      data += state.dataGenFns[i](el);
    }
    if (el.attrs) {
      data += "attrs:".concat(genProps(el.attrs), ",");
    }
    if (el.props) {
      data += "domProps:".concat(genProps(el.props), ",");
    }
    if (el.events) {
      data += "".concat(genHandlers(el.events, false), ",");
    }
    if (el.nativeEvents) {
      data += "".concat(genHandlers(el.nativeEvents, true), ",");
    }
    if (el.slotTarget && !el.slotScope) {
      data += "slot:".concat(el.slotTarget, ",");
    }
    if (el.scopedSlots) {
      data += "".concat(genScopedSlots(el, el.scopedSlots, state), ",");
    }
    if (el.model) {
      data += "model:{value:".concat(el.model.value, ",callback:").concat(el.model.callback, ",expression:").concat(el.model.expression, "},");
    }
    if (el.inlineTemplate) {
      var inlineTemplate = genInlineTemplate(el, state);
      if (inlineTemplate) {
        data += "".concat(inlineTemplate, ",");
      }
    }
    data = data.replace(/,$/, "") + "}";
    if (el.dynamicAttrs) {
      data = "_b(".concat(data, ',"').concat(el.tag, '",').concat(genProps(el.dynamicAttrs), ")");
    }
    if (el.wrapData) {
      data = el.wrapData(data);
    }
    if (el.wrapListeners) {
      data = el.wrapListeners(data);
    }
    return data;
  }
  function genDirectives(el, state) {
    var dirs = el.directives;
    if (!dirs)
      return;
    var res = "directives:[";
    var hasRuntime = false;
    var i, l, dir, needRuntime;
    for (i = 0, l = dirs.length; i < l; i++) {
      dir = dirs[i];
      needRuntime = true;
      var gen = state.directives[dir.name];
      if (gen) {
        needRuntime = !!gen(el, dir, state.warn);
      }
      if (needRuntime) {
        hasRuntime = true;
        res += '{name:"'.concat(dir.name, '",rawName:"').concat(dir.rawName, '"').concat(dir.value ? ",value:(".concat(dir.value, "),expression:").concat(JSON.stringify(dir.value)) : "").concat(dir.arg ? ",arg:".concat(dir.isDynamicArg ? dir.arg : '"'.concat(dir.arg, '"')) : "").concat(dir.modifiers ? ",modifiers:".concat(JSON.stringify(dir.modifiers)) : "", "},");
      }
    }
    if (hasRuntime) {
      return res.slice(0, -1) + "]";
    }
  }
  function genInlineTemplate(el, state) {
    var ast = el.children[0];
    if (el.children.length !== 1 || ast.type !== 1) {
      state.warn("Inline-template components must have exactly one child element.", { start: el.start });
    }
    if (ast && ast.type === 1) {
      var inlineRenderFns = generate(ast, state.options);
      return "inlineTemplate:{render:function(){".concat(inlineRenderFns.render, "},staticRenderFns:[").concat(inlineRenderFns.staticRenderFns.map(function(code) {
        return "function(){".concat(code, "}");
      }).join(","), "]}");
    }
  }
  function genScopedSlots(el, slots, state) {
    var needsForceUpdate = el.for || Object.keys(slots).some(function(key) {
      var slot = slots[key];
      return slot.slotTargetDynamic || slot.if || slot.for || containsSlotChild(slot);
    });
    var needsKey = !!el.if;
    if (!needsForceUpdate) {
      var parent_2 = el.parent;
      while (parent_2) {
        if (parent_2.slotScope && parent_2.slotScope !== emptySlotScopeToken || parent_2.for) {
          needsForceUpdate = true;
          break;
        }
        if (parent_2.if) {
          needsKey = true;
        }
        parent_2 = parent_2.parent;
      }
    }
    var generatedSlots = Object.keys(slots).map(function(key) {
      return genScopedSlot(slots[key], state);
    }).join(",");
    return "scopedSlots:_u([".concat(generatedSlots, "]").concat(needsForceUpdate ? ",null,true" : "").concat(!needsForceUpdate && needsKey ? ",null,false,".concat(hash(generatedSlots)) : "", ")");
  }
  function hash(str2) {
    var hash2 = 5381;
    var i = str2.length;
    while (i) {
      hash2 = hash2 * 33 ^ str2.charCodeAt(--i);
    }
    return hash2 >>> 0;
  }
  function containsSlotChild(el) {
    if (el.type === 1) {
      if (el.tag === "slot") {
        return true;
      }
      return el.children.some(containsSlotChild);
    }
    return false;
  }
  function genScopedSlot(el, state) {
    var isLegacySyntax = el.attrsMap["slot-scope"];
    if (el.if && !el.ifProcessed && !isLegacySyntax) {
      return genIf(el, state, genScopedSlot, "null");
    }
    if (el.for && !el.forProcessed) {
      return genFor(el, state, genScopedSlot);
    }
    var slotScope = el.slotScope === emptySlotScopeToken ? "" : String(el.slotScope);
    var fn = "function(".concat(slotScope, "){") + "return ".concat(el.tag === "template" ? el.if && isLegacySyntax ? "(".concat(el.if, ")?").concat(genChildren(el, state) || "undefined", ":undefined") : genChildren(el, state) || "undefined" : genElement(el, state), "}");
    var reverseProxy = slotScope ? "" : ",proxy:true";
    return "{key:".concat(el.slotTarget || '"default"', ",fn:").concat(fn).concat(reverseProxy, "}");
  }
  function genChildren(el, state, checkSkip, altGenElement, altGenNode) {
    var children = el.children;
    if (children.length) {
      var el_1 = children[0];
      if (children.length === 1 && el_1.for && el_1.tag !== "template" && el_1.tag !== "slot") {
        var normalizationType_1 = checkSkip ? state.maybeComponent(el_1) ? ",1" : ",0" : "";
        return "".concat((altGenElement || genElement)(el_1, state)).concat(normalizationType_1);
      }
      var normalizationType = checkSkip ? getNormalizationType(children, state.maybeComponent) : 0;
      var gen_1 = altGenNode || genNode;
      return "[".concat(children.map(function(c) {
        return gen_1(c, state);
      }).join(","), "]").concat(normalizationType ? ",".concat(normalizationType) : "");
    }
  }
  function getNormalizationType(children, maybeComponent2) {
    var res = 0;
    for (var i = 0; i < children.length; i++) {
      var el = children[i];
      if (el.type !== 1) {
        continue;
      }
      if (needsNormalization(el) || el.ifConditions && el.ifConditions.some(function(c) {
        return needsNormalization(c.block);
      })) {
        res = 2;
        break;
      }
      if (maybeComponent2(el) || el.ifConditions && el.ifConditions.some(function(c) {
        return maybeComponent2(c.block);
      })) {
        res = 1;
      }
    }
    return res;
  }
  function needsNormalization(el) {
    return el.for !== void 0 || el.tag === "template" || el.tag === "slot";
  }
  function genNode(node, state) {
    if (node.type === 1) {
      return genElement(node, state);
    } else if (node.type === 3 && node.isComment) {
      return genComment(node);
    } else {
      return genText(node);
    }
  }
  function genText(text2) {
    return "_v(".concat(text2.type === 2 ? text2.expression : transformSpecialNewlines(JSON.stringify(text2.text)), ")");
  }
  function genComment(comment2) {
    return "_e(".concat(JSON.stringify(comment2.text), ")");
  }
  function genSlot(el, state) {
    var slotName = el.slotName || '"default"';
    var children = genChildren(el, state);
    var res = "_t(".concat(slotName).concat(children ? ",function(){return ".concat(children, "}") : "");
    var attrs2 = el.attrs || el.dynamicAttrs ? genProps((el.attrs || []).concat(el.dynamicAttrs || []).map(function(attr) {
      return {
        // slot props are camelized
        name: camelize(attr.name),
        value: attr.value,
        dynamic: attr.dynamic
      };
    })) : null;
    var bind2 = el.attrsMap["v-bind"];
    if ((attrs2 || bind2) && !children) {
      res += ",null";
    }
    if (attrs2) {
      res += ",".concat(attrs2);
    }
    if (bind2) {
      res += "".concat(attrs2 ? "" : ",null", ",").concat(bind2);
    }
    return res + ")";
  }
  function genComponent(componentName, el, state) {
    var children = el.inlineTemplate ? null : genChildren(el, state, true);
    return "_c(".concat(componentName, ",").concat(genData(el, state)).concat(children ? ",".concat(children) : "", ")");
  }
  function genProps(props2) {
    var staticProps = "";
    var dynamicProps = "";
    for (var i = 0; i < props2.length; i++) {
      var prop = props2[i];
      var value = transformSpecialNewlines(prop.value);
      if (prop.dynamic) {
        dynamicProps += "".concat(prop.name, ",").concat(value, ",");
      } else {
        staticProps += '"'.concat(prop.name, '":').concat(value, ",");
      }
    }
    staticProps = "{".concat(staticProps.slice(0, -1), "}");
    if (dynamicProps) {
      return "_d(".concat(staticProps, ",[").concat(dynamicProps.slice(0, -1), "])");
    } else {
      return staticProps;
    }
  }
  function transformSpecialNewlines(text2) {
    return text2.replace(/\u2028/g, "\\u2028").replace(/\u2029/g, "\\u2029");
  }
  var prohibitedKeywordRE = new RegExp("\\b" + "do,if,for,let,new,try,var,case,else,with,await,break,catch,class,const,super,throw,while,yield,delete,export,import,return,switch,default,extends,finally,continue,debugger,function,arguments".split(",").join("\\b|\\b") + "\\b");
  var unaryOperatorsRE = new RegExp("\\b" + "delete,typeof,void".split(",").join("\\s*\\([^\\)]*\\)|\\b") + "\\s*\\([^\\)]*\\)");
  var stripStringRE = /'(?:[^'\\]|\\.)*'|"(?:[^"\\]|\\.)*"|`(?:[^`\\]|\\.)*\$\{|\}(?:[^`\\]|\\.)*`|`(?:[^`\\]|\\.)*`/g;
  function detectErrors(ast, warn2) {
    if (ast) {
      checkNode(ast, warn2);
    }
  }
  function checkNode(node, warn2) {
    if (node.type === 1) {
      for (var name_1 in node.attrsMap) {
        if (dirRE.test(name_1)) {
          var value = node.attrsMap[name_1];
          if (value) {
            var range2 = node.rawAttrsMap[name_1];
            if (name_1 === "v-for") {
              checkFor(node, 'v-for="'.concat(value, '"'), warn2, range2);
            } else if (name_1 === "v-slot" || name_1[0] === "#") {
              checkFunctionParameterExpression(value, "".concat(name_1, '="').concat(value, '"'), warn2, range2);
            } else if (onRE.test(name_1)) {
              checkEvent(value, "".concat(name_1, '="').concat(value, '"'), warn2, range2);
            } else {
              checkExpression(value, "".concat(name_1, '="').concat(value, '"'), warn2, range2);
            }
          }
        }
      }
      if (node.children) {
        for (var i = 0; i < node.children.length; i++) {
          checkNode(node.children[i], warn2);
        }
      }
    } else if (node.type === 2) {
      checkExpression(node.expression, node.text, warn2, node);
    }
  }
  function checkEvent(exp, text2, warn2, range2) {
    var stripped = exp.replace(stripStringRE, "");
    var keywordMatch = stripped.match(unaryOperatorsRE);
    if (keywordMatch && stripped.charAt(keywordMatch.index - 1) !== "$") {
      warn2("avoid using JavaScript unary operator as property name: " + '"'.concat(keywordMatch[0], '" in expression ').concat(text2.trim()), range2);
    }
    checkExpression(exp, text2, warn2, range2);
  }
  function checkFor(node, text2, warn2, range2) {
    checkExpression(node.for || "", text2, warn2, range2);
    checkIdentifier(node.alias, "v-for alias", text2, warn2, range2);
    checkIdentifier(node.iterator1, "v-for iterator", text2, warn2, range2);
    checkIdentifier(node.iterator2, "v-for iterator", text2, warn2, range2);
  }
  function checkIdentifier(ident, type, text2, warn2, range2) {
    if (typeof ident === "string") {
      try {
        new Function("var ".concat(ident, "=_"));
      } catch (e) {
        warn2("invalid ".concat(type, ' "').concat(ident, '" in expression: ').concat(text2.trim()), range2);
      }
    }
  }
  function checkExpression(exp, text2, warn2, range2) {
    try {
      new Function("return ".concat(exp));
    } catch (e) {
      var keywordMatch = exp.replace(stripStringRE, "").match(prohibitedKeywordRE);
      if (keywordMatch) {
        warn2("avoid using JavaScript keyword as property name: " + '"'.concat(keywordMatch[0], '"\n  Raw expression: ').concat(text2.trim()), range2);
      } else {
        warn2("invalid expression: ".concat(e.message, " in\n\n") + "    ".concat(exp, "\n\n") + "  Raw expression: ".concat(text2.trim(), "\n"), range2);
      }
    }
  }
  function checkFunctionParameterExpression(exp, text2, warn2, range2) {
    try {
      new Function(exp, "");
    } catch (e) {
      warn2("invalid function parameter expression: ".concat(e.message, " in\n\n") + "    ".concat(exp, "\n\n") + "  Raw expression: ".concat(text2.trim(), "\n"), range2);
    }
  }
  var range = 2;
  function generateCodeFrame(source, start, end) {
    if (start === void 0) {
      start = 0;
    }
    if (end === void 0) {
      end = source.length;
    }
    var lines = source.split(/\r?\n/);
    var count = 0;
    var res = [];
    for (var i = 0; i < lines.length; i++) {
      count += lines[i].length + 1;
      if (count >= start) {
        for (var j = i - range; j <= i + range || end > count; j++) {
          if (j < 0 || j >= lines.length)
            continue;
          res.push("".concat(j + 1).concat(repeat(" ", 3 - String(j + 1).length), "|  ").concat(lines[j]));
          var lineLength = lines[j].length;
          if (j === i) {
            var pad = start - (count - lineLength) + 1;
            var length_1 = end > count ? lineLength - pad : end - start;
            res.push("   |  " + repeat(" ", pad) + repeat("^", length_1));
          } else if (j > i) {
            if (end > count) {
              var length_2 = Math.min(end - count, lineLength);
              res.push("   |  " + repeat("^", length_2));
            }
            count += lineLength + 1;
          }
        }
        break;
      }
    }
    return res.join("\n");
  }
  function repeat(str2, n) {
    var result = "";
    if (n > 0) {
      while (true) {
        if (n & 1)
          result += str2;
        n >>>= 1;
        if (n <= 0)
          break;
        str2 += str2;
      }
    }
    return result;
  }
  function createFunction(code, errors) {
    try {
      return new Function(code);
    } catch (err) {
      errors.push({ err, code });
      return noop;
    }
  }
  function createCompileToFunctionFn(compile) {
    var cache2 = /* @__PURE__ */ Object.create(null);
    return function compileToFunctions2(template, options, vm3) {
      options = extend({}, options);
      var warn2 = options.warn || warn$2;
      delete options.warn;
      if (true) {
        try {
          new Function("return 1");
        } catch (e) {
          if (e.toString().match(/unsafe-eval|CSP/)) {
            warn2("It seems you are using the standalone build of Vue.js in an environment with Content Security Policy that prohibits unsafe-eval. The template compiler cannot work in this environment. Consider relaxing the policy to allow unsafe-eval or pre-compiling your templates into render functions.");
          }
        }
      }
      var key = options.delimiters ? String(options.delimiters) + template : template;
      if (cache2[key]) {
        return cache2[key];
      }
      var compiled = compile(template, options);
      if (true) {
        if (compiled.errors && compiled.errors.length) {
          if (options.outputSourceRange) {
            compiled.errors.forEach(function(e) {
              warn2("Error compiling template:\n\n".concat(e.msg, "\n\n") + generateCodeFrame(template, e.start, e.end), vm3);
            });
          } else {
            warn2("Error compiling template:\n\n".concat(template, "\n\n") + compiled.errors.map(function(e) {
              return "- ".concat(e);
            }).join("\n") + "\n", vm3);
          }
        }
        if (compiled.tips && compiled.tips.length) {
          if (options.outputSourceRange) {
            compiled.tips.forEach(function(e) {
              return tip(e.msg, vm3);
            });
          } else {
            compiled.tips.forEach(function(msg) {
              return tip(msg, vm3);
            });
          }
        }
      }
      var res = {};
      var fnGenErrors = [];
      res.render = createFunction(compiled.render, fnGenErrors);
      res.staticRenderFns = compiled.staticRenderFns.map(function(code) {
        return createFunction(code, fnGenErrors);
      });
      if (true) {
        if ((!compiled.errors || !compiled.errors.length) && fnGenErrors.length) {
          warn2("Failed to generate render function:\n\n" + fnGenErrors.map(function(_a2) {
            var err = _a2.err, code = _a2.code;
            return "".concat(err.toString(), " in\n\n").concat(code, "\n");
          }).join("\n"), vm3);
        }
      }
      return cache2[key] = res;
    };
  }
  function createCompilerCreator(baseCompile2) {
    return function createCompiler2(baseOptions2) {
      function compile(template, options) {
        var finalOptions = Object.create(baseOptions2);
        var errors = [];
        var tips = [];
        var warn2 = function(msg, range2, tip2) {
          (tip2 ? tips : errors).push(msg);
        };
        if (options) {
          if (options.outputSourceRange) {
            var leadingSpaceLength_1 = template.match(/^\s*/)[0].length;
            warn2 = function(msg, range2, tip2) {
              var data = typeof msg === "string" ? { msg } : msg;
              if (range2) {
                if (range2.start != null) {
                  data.start = range2.start + leadingSpaceLength_1;
                }
                if (range2.end != null) {
                  data.end = range2.end + leadingSpaceLength_1;
                }
              }
              (tip2 ? tips : errors).push(data);
            };
          }
          if (options.modules) {
            finalOptions.modules = (baseOptions2.modules || []).concat(options.modules);
          }
          if (options.directives) {
            finalOptions.directives = extend(Object.create(baseOptions2.directives || null), options.directives);
          }
          for (var key in options) {
            if (key !== "modules" && key !== "directives") {
              finalOptions[key] = options[key];
            }
          }
        }
        finalOptions.warn = warn2;
        var compiled = baseCompile2(template.trim(), finalOptions);
        if (true) {
          detectErrors(compiled.ast, warn2);
        }
        compiled.errors = errors;
        compiled.tips = tips;
        return compiled;
      }
      return {
        compile,
        compileToFunctions: createCompileToFunctionFn(compile)
      };
    };
  }
  var createCompiler = createCompilerCreator(function baseCompile(template, options) {
    var ast = parse(template.trim(), options);
    if (options.optimize !== false) {
      optimize(ast, options);
    }
    var code = generate(ast, options);
    return {
      ast,
      render: code.render,
      staticRenderFns: code.staticRenderFns
    };
  });
  var _a = createCompiler(baseOptions);
  var compileToFunctions = _a.compileToFunctions;
  var div;
  function getShouldDecode(href) {
    div = div || document.createElement("div");
    div.innerHTML = href ? '<a href="\n"/>' : '<div a="\n"/>';
    return div.innerHTML.indexOf("&#10;") > 0;
  }
  var shouldDecodeNewlines = inBrowser ? getShouldDecode(false) : false;
  var shouldDecodeNewlinesForHref = inBrowser ? getShouldDecode(true) : false;
  var idToTemplate = cached(function(id) {
    var el = query(id);
    return el && el.innerHTML;
  });
  var mount = Vue.prototype.$mount;
  Vue.prototype.$mount = function(el, hydrating) {
    el = el && query(el);
    if (el === document.body || el === document.documentElement) {
      warn$2("Do not mount Vue to <html> or <body> - mount to normal elements instead.");
      return this;
    }
    var options = this.$options;
    if (!options.render) {
      var template = options.template;
      if (template) {
        if (typeof template === "string") {
          if (template.charAt(0) === "#") {
            template = idToTemplate(template);
            if (!template) {
              warn$2("Template element not found or is empty: ".concat(options.template), this);
            }
          }
        } else if (template.nodeType) {
          template = template.innerHTML;
        } else {
          if (true) {
            warn$2("invalid template option:" + template, this);
          }
          return this;
        }
      } else if (el) {
        template = getOuterHTML(el);
      }
      if (template) {
        if (config.performance && mark) {
          mark("compile");
        }
        var _a2 = compileToFunctions(template, {
          outputSourceRange: true,
          shouldDecodeNewlines,
          shouldDecodeNewlinesForHref,
          delimiters: options.delimiters,
          comments: options.comments
        }, this), render = _a2.render, staticRenderFns = _a2.staticRenderFns;
        options.render = render;
        options.staticRenderFns = staticRenderFns;
        if (config.performance && mark) {
          mark("compile end");
          measure("vue ".concat(this._name, " compile"), "compile", "compile end");
        }
      }
    }
    return mount.call(this, el, hydrating);
  };
  function getOuterHTML(el) {
    if (el.outerHTML) {
      return el.outerHTML;
    } else {
      var container = document.createElement("div");
      container.appendChild(el.cloneNode(true));
      return container.innerHTML;
    }
  }
  Vue.compile = compileToFunctions;

  // node_modules/@fluent/bundle/esm/types.js
  var FluentType = class {
    /**
     * Create a `FluentType` instance.
     *
     * @param value The JavaScript value to wrap.
     */
    constructor(value) {
      this.value = value;
    }
    /**
     * Unwrap the raw value stored by this `FluentType`.
     */
    valueOf() {
      return this.value;
    }
  };
  var FluentNone = class extends FluentType {
    /**
     * Create an instance of `FluentNone` with an optional fallback value.
     * @param value The fallback value of this `FluentNone`.
     */
    constructor(value = "???") {
      super(value);
    }
    /**
     * Format this `FluentNone` to the fallback string.
     */
    toString(scope) {
      return `{${this.value}}`;
    }
  };
  var FluentNumber = class extends FluentType {
    /**
     * Create an instance of `FluentNumber` with options to the
     * `Intl.NumberFormat` constructor.
     *
     * @param value The number value of this `FluentNumber`.
     * @param opts Options which will be passed to `Intl.NumberFormat`.
     */
    constructor(value, opts2 = {}) {
      super(value);
      this.opts = opts2;
    }
    /**
     * Format this `FluentNumber` to a string.
     */
    toString(scope) {
      if (scope) {
        try {
          const nf = scope.memoizeIntlObject(Intl.NumberFormat, this.opts);
          return nf.format(this.value);
        } catch (err) {
          scope.reportError(err);
        }
      }
      return this.value.toString(10);
    }
  };
  var FluentDateTime = class _FluentDateTime extends FluentType {
    static supportsValue(value) {
      if (typeof value === "number")
        return true;
      if (value instanceof Date)
        return true;
      if (value instanceof FluentType)
        return _FluentDateTime.supportsValue(value.valueOf());
      if ("Temporal" in globalThis) {
        const _Temporal = globalThis.Temporal;
        if (value instanceof _Temporal.Instant || value instanceof _Temporal.PlainDateTime || value instanceof _Temporal.PlainDate || value instanceof _Temporal.PlainMonthDay || value instanceof _Temporal.PlainTime || value instanceof _Temporal.PlainYearMonth) {
          return true;
        }
      }
      return false;
    }
    /**
     * Create an instance of `FluentDateTime` with options to the
     * `Intl.DateTimeFormat` constructor.
     *
     * @param value The number value of this `FluentDateTime`, in milliseconds.
     * @param opts Options which will be passed to `Intl.DateTimeFormat`.
     */
    constructor(value, opts2 = {}) {
      if (value instanceof _FluentDateTime) {
        opts2 = { ...value.opts, ...opts2 };
        value = value.value;
      } else if (value instanceof FluentType) {
        value = value.valueOf();
      }
      if (typeof value === "object" && "calendarId" in value && opts2.calendar === void 0) {
        opts2 = { ...opts2, calendar: value.calendarId };
      }
      super(value);
      this.opts = opts2;
    }
    [Symbol.toPrimitive](hint) {
      return hint === "string" ? this.toString() : this.toNumber();
    }
    /**
     * Convert this `FluentDateTime` to a number.
     * Note that this isn't always possible due to the nature of Temporal objects.
     * In such cases, a TypeError will be thrown.
     */
    toNumber() {
      const value = this.value;
      if (typeof value === "number")
        return value;
      if (value instanceof Date)
        return value.getTime();
      if ("epochMilliseconds" in value) {
        return value.epochMilliseconds;
      }
      if ("toZonedDateTime" in value) {
        return value.toZonedDateTime("UTC").epochMilliseconds;
      }
      throw new TypeError("Unwrapping a non-number value as a number");
    }
    /**
     * Format this `FluentDateTime` to a string.
     */
    toString(scope) {
      if (scope) {
        try {
          const dtf = scope.memoizeIntlObject(Intl.DateTimeFormat, this.opts);
          return dtf.format(this.value);
        } catch (err) {
          scope.reportError(err);
        }
      }
      if (typeof this.value === "number" || this.value instanceof Date) {
        return new Date(this.value).toISOString();
      }
      return this.value.toString();
    }
  };

  // node_modules/@fluent/bundle/esm/resolver.js
  var MAX_PLACEABLES = 100;
  var FSI = "\u2068";
  var PDI = "\u2069";
  function match(scope, selector, key) {
    if (key === selector) {
      return true;
    }
    if (key instanceof FluentNumber && selector instanceof FluentNumber && key.value === selector.value) {
      return true;
    }
    if (selector instanceof FluentNumber && typeof key === "string") {
      let category = scope.memoizeIntlObject(Intl.PluralRules, selector.opts).select(selector.value);
      if (key === category) {
        return true;
      }
    }
    return false;
  }
  function getDefault(scope, variants, star) {
    if (variants[star]) {
      return resolvePattern(scope, variants[star].value);
    }
    scope.reportError(new RangeError("No default"));
    return new FluentNone();
  }
  function getArguments(scope, args) {
    const positional = [];
    const named = /* @__PURE__ */ Object.create(null);
    for (const arg of args) {
      if (arg.type === "narg") {
        named[arg.name] = resolveExpression(scope, arg.value);
      } else {
        positional.push(resolveExpression(scope, arg));
      }
    }
    return { positional, named };
  }
  function resolveExpression(scope, expr) {
    switch (expr.type) {
      case "str":
        return expr.value;
      case "num":
        return new FluentNumber(expr.value, {
          minimumFractionDigits: expr.precision
        });
      case "var":
        return resolveVariableReference(scope, expr);
      case "mesg":
        return resolveMessageReference(scope, expr);
      case "term":
        return resolveTermReference(scope, expr);
      case "func":
        return resolveFunctionReference(scope, expr);
      case "select":
        return resolveSelectExpression(scope, expr);
      default:
        return new FluentNone();
    }
  }
  function resolveVariableReference(scope, { name }) {
    let arg;
    if (scope.params) {
      if (Object.prototype.hasOwnProperty.call(scope.params, name)) {
        arg = scope.params[name];
      } else {
        return new FluentNone(`$${name}`);
      }
    } else if (scope.args && Object.prototype.hasOwnProperty.call(scope.args, name)) {
      arg = scope.args[name];
    } else {
      scope.reportError(new ReferenceError(`Unknown variable: $${name}`));
      return new FluentNone(`$${name}`);
    }
    if (arg instanceof FluentType) {
      return arg;
    }
    switch (typeof arg) {
      case "string":
        return arg;
      case "number":
        return new FluentNumber(arg);
      case "object":
        if (FluentDateTime.supportsValue(arg)) {
          return new FluentDateTime(arg);
        }
      // eslint-disable-next-line no-fallthrough
      default:
        scope.reportError(new TypeError(`Variable type not supported: $${name}, ${typeof arg}`));
        return new FluentNone(`$${name}`);
    }
  }
  function resolveMessageReference(scope, { name, attr }) {
    const message = scope.bundle._messages.get(name);
    if (!message) {
      scope.reportError(new ReferenceError(`Unknown message: ${name}`));
      return new FluentNone(name);
    }
    if (attr) {
      const attribute2 = message.attributes[attr];
      if (attribute2) {
        return resolvePattern(scope, attribute2);
      }
      scope.reportError(new ReferenceError(`Unknown attribute: ${attr}`));
      return new FluentNone(`${name}.${attr}`);
    }
    if (message.value) {
      return resolvePattern(scope, message.value);
    }
    scope.reportError(new ReferenceError(`No value: ${name}`));
    return new FluentNone(name);
  }
  function resolveTermReference(scope, { name, attr, args }) {
    const id = `-${name}`;
    const term = scope.bundle._terms.get(id);
    if (!term) {
      scope.reportError(new ReferenceError(`Unknown term: ${id}`));
      return new FluentNone(id);
    }
    if (attr) {
      const attribute2 = term.attributes[attr];
      if (attribute2) {
        scope.params = getArguments(scope, args).named;
        const resolved2 = resolvePattern(scope, attribute2);
        scope.params = null;
        return resolved2;
      }
      scope.reportError(new ReferenceError(`Unknown attribute: ${attr}`));
      return new FluentNone(`${id}.${attr}`);
    }
    scope.params = getArguments(scope, args).named;
    const resolved = resolvePattern(scope, term.value);
    scope.params = null;
    return resolved;
  }
  function resolveFunctionReference(scope, { name, args }) {
    let func = scope.bundle._functions[name];
    if (!func) {
      scope.reportError(new ReferenceError(`Unknown function: ${name}()`));
      return new FluentNone(`${name}()`);
    }
    if (typeof func !== "function") {
      scope.reportError(new TypeError(`Function ${name}() is not callable`));
      return new FluentNone(`${name}()`);
    }
    try {
      let resolved = getArguments(scope, args);
      return func(resolved.positional, resolved.named);
    } catch (err) {
      scope.reportError(err);
      return new FluentNone(`${name}()`);
    }
  }
  function resolveSelectExpression(scope, { selector, variants, star }) {
    let sel = resolveExpression(scope, selector);
    if (sel instanceof FluentNone) {
      return getDefault(scope, variants, star);
    }
    for (const variant of variants) {
      const key = resolveExpression(scope, variant.key);
      if (match(scope, sel, key)) {
        return resolvePattern(scope, variant.value);
      }
    }
    return getDefault(scope, variants, star);
  }
  function resolveComplexPattern(scope, ptn) {
    if (scope.dirty.has(ptn)) {
      scope.reportError(new RangeError("Cyclic reference"));
      return new FluentNone();
    }
    scope.dirty.add(ptn);
    const result = [];
    const useIsolating = scope.bundle._useIsolating && ptn.length > 1;
    for (const elem of ptn) {
      if (typeof elem === "string") {
        result.push(scope.bundle._transform(elem));
        continue;
      }
      scope.placeables++;
      if (scope.placeables > MAX_PLACEABLES) {
        scope.dirty.delete(ptn);
        throw new RangeError(`Too many placeables expanded: ${scope.placeables}, max allowed is ${MAX_PLACEABLES}`);
      }
      if (useIsolating) {
        result.push(FSI);
      }
      result.push(resolveExpression(scope, elem).toString(scope));
      if (useIsolating) {
        result.push(PDI);
      }
    }
    scope.dirty.delete(ptn);
    return result.join("");
  }
  function resolvePattern(scope, value) {
    if (typeof value === "string") {
      return scope.bundle._transform(value);
    }
    return resolveComplexPattern(scope, value);
  }

  // node_modules/@fluent/bundle/esm/scope.js
  var Scope = class {
    constructor(bundle, errors, args) {
      this.dirty = /* @__PURE__ */ new WeakSet();
      this.params = null;
      this.placeables = 0;
      this.bundle = bundle;
      this.errors = errors;
      this.args = args;
    }
    reportError(error) {
      if (!this.errors || !(error instanceof Error)) {
        throw error;
      }
      this.errors.push(error);
    }
    memoizeIntlObject(ctor, opts2) {
      let cache2 = this.bundle._intls.get(ctor);
      if (!cache2) {
        cache2 = {};
        this.bundle._intls.set(ctor, cache2);
      }
      let id = JSON.stringify(opts2);
      if (!cache2[id]) {
        cache2[id] = new ctor(this.bundle.locales, opts2);
      }
      return cache2[id];
    }
  };

  // node_modules/@fluent/bundle/esm/builtins.js
  function values(opts2, allowed) {
    const unwrapped = /* @__PURE__ */ Object.create(null);
    for (const [name, opt] of Object.entries(opts2)) {
      if (allowed.includes(name)) {
        unwrapped[name] = opt.valueOf();
      }
    }
    return unwrapped;
  }
  var NUMBER_ALLOWED = [
    "unitDisplay",
    "currencyDisplay",
    "useGrouping",
    "minimumIntegerDigits",
    "minimumFractionDigits",
    "maximumFractionDigits",
    "minimumSignificantDigits",
    "maximumSignificantDigits"
  ];
  function NUMBER(args, opts2) {
    let arg = args[0];
    if (arg instanceof FluentNone) {
      return new FluentNone(`NUMBER(${arg.valueOf()})`);
    }
    if (arg instanceof FluentNumber) {
      return new FluentNumber(arg.valueOf(), {
        ...arg.opts,
        ...values(opts2, NUMBER_ALLOWED)
      });
    }
    if (arg instanceof FluentDateTime) {
      return new FluentNumber(arg.toNumber(), {
        ...values(opts2, NUMBER_ALLOWED)
      });
    }
    throw new TypeError("Invalid argument to NUMBER");
  }
  var DATETIME_ALLOWED = [
    "dateStyle",
    "timeStyle",
    "fractionalSecondDigits",
    "dayPeriod",
    "hour12",
    "weekday",
    "era",
    "year",
    "month",
    "day",
    "hour",
    "minute",
    "second",
    "timeZoneName"
  ];
  function DATETIME(args, opts2) {
    let arg = args[0];
    if (arg instanceof FluentNone) {
      return new FluentNone(`DATETIME(${arg.valueOf()})`);
    }
    if (arg instanceof FluentDateTime || arg instanceof FluentNumber) {
      return new FluentDateTime(arg, values(opts2, DATETIME_ALLOWED));
    }
    throw new TypeError("Invalid argument to DATETIME");
  }

  // node_modules/@fluent/bundle/esm/memoizer.js
  var cache = /* @__PURE__ */ new Map();
  function getMemoizerForLocale(locales) {
    const stringLocale = Array.isArray(locales) ? locales.join(" ") : locales;
    let memoizer = cache.get(stringLocale);
    if (memoizer === void 0) {
      memoizer = /* @__PURE__ */ new Map();
      cache.set(stringLocale, memoizer);
    }
    return memoizer;
  }

  // node_modules/@fluent/bundle/esm/bundle.js
  var FluentBundle = class {
    /**
     * Create an instance of `FluentBundle`.
     *
     * @example
     * ```js
     * let bundle = new FluentBundle(["en-US", "en"]);
     *
     * let bundle = new FluentBundle(locales, {useIsolating: false});
     *
     * let bundle = new FluentBundle(locales, {
     *   useIsolating: true,
     *   functions: {
     *     NODE_ENV: () => process.env.NODE_ENV
     *   }
     * });
     * ```
     *
     * @param locales - Used to instantiate `Intl` formatters used by translations.
     * @param options - Optional configuration for the bundle.
     */
    constructor(locales, { functions, useIsolating = true, transform = (v) => v } = {}) {
      this._terms = /* @__PURE__ */ new Map();
      this._messages = /* @__PURE__ */ new Map();
      this.locales = Array.isArray(locales) ? locales : [locales];
      this._functions = {
        NUMBER,
        DATETIME,
        ...functions
      };
      this._useIsolating = useIsolating;
      this._transform = transform;
      this._intls = getMemoizerForLocale(locales);
    }
    /**
     * Check if a message is present in the bundle.
     *
     * @param id - The identifier of the message to check.
     */
    hasMessage(id) {
      return this._messages.has(id);
    }
    /**
     * Return a raw unformatted message object from the bundle.
     *
     * Raw messages are `{value, attributes}` shapes containing translation units
     * called `Patterns`. `Patterns` are implementation-specific; they should be
     * treated as black boxes and formatted with `FluentBundle.formatPattern`.
     *
     * @param id - The identifier of the message to check.
     */
    getMessage(id) {
      return this._messages.get(id);
    }
    /**
     * Add a translation resource to the bundle.
     *
     * @example
     * ```js
     * let res = new FluentResource("foo = Foo");
     * bundle.addResource(res);
     * bundle.getMessage("foo");
     * // → {value: .., attributes: {..}}
     * ```
     *
     * @param res
     * @param options
     */
    addResource(res, { allowOverrides = false } = {}) {
      const errors = [];
      for (let i = 0; i < res.body.length; i++) {
        let entry = res.body[i];
        if (entry.id.startsWith("-")) {
          if (allowOverrides === false && this._terms.has(entry.id)) {
            errors.push(new Error(`Attempt to override an existing term: "${entry.id}"`));
            continue;
          }
          this._terms.set(entry.id, entry);
        } else {
          if (allowOverrides === false && this._messages.has(entry.id)) {
            errors.push(new Error(`Attempt to override an existing message: "${entry.id}"`));
            continue;
          }
          this._messages.set(entry.id, entry);
        }
      }
      return errors;
    }
    /**
     * Format a `Pattern` to a string.
     *
     * Format a raw `Pattern` into a string. `args` will be used to resolve
     * references to variables passed as arguments to the translation.
     *
     * In case of errors `formatPattern` will try to salvage as much of the
     * translation as possible and will still return a string. For performance
     * reasons, the encountered errors are not returned but instead are appended
     * to the `errors` array passed as the third argument.
     *
     * If `errors` is omitted, the first encountered error will be thrown.
     *
     * @example
     * ```js
     * let errors = [];
     * bundle.addResource(
     *     new FluentResource("hello = Hello, {$name}!"));
     *
     * let hello = bundle.getMessage("hello");
     * if (hello.value) {
     *     bundle.formatPattern(hello.value, {name: "Jane"}, errors);
     *     // Returns "Hello, Jane!" and `errors` is empty.
     *
     *     bundle.formatPattern(hello.value, undefined, errors);
     *     // Returns "Hello, {$name}!" and `errors` is now:
     *     // [<ReferenceError: Unknown variable: name>]
     * }
     * ```
     */
    formatPattern(pattern, args = null, errors = null) {
      if (typeof pattern === "string") {
        return this._transform(pattern);
      }
      let scope = new Scope(this, errors, args);
      try {
        let value = resolveComplexPattern(scope, pattern);
        return value.toString(scope);
      } catch (err) {
        if (scope.errors && err instanceof Error) {
          scope.errors.push(err);
          return new FluentNone().toString(scope);
        }
        throw err;
      }
    }
  };

  // node_modules/@fluent/bundle/esm/resource.js
  var RE_MESSAGE_START = /^(-?[a-zA-Z][\w-]*) *= */gm;
  var RE_ATTRIBUTE_START = /\.([a-zA-Z][\w-]*) *= */y;
  var RE_VARIANT_START = /\*?\[/y;
  var RE_NUMBER_LITERAL = /(-?[0-9]+(?:\.([0-9]+))?)/y;
  var RE_IDENTIFIER = /([a-zA-Z][\w-]*)/y;
  var RE_REFERENCE = /([$-])?([a-zA-Z][\w-]*)(?:\.([a-zA-Z][\w-]*))?/y;
  var RE_FUNCTION_NAME = /^[A-Z][A-Z0-9_-]*$/;
  var RE_TEXT_RUN = /([^{}\n\r]+)/y;
  var RE_STRING_RUN = /([^\\"\n\r]*)/y;
  var RE_STRING_ESCAPE = /\\([\\"])/y;
  var RE_UNICODE_ESCAPE = /\\u([a-fA-F0-9]{4})|\\U([a-fA-F0-9]{6})/y;
  var RE_LEADING_NEWLINES = /^\n+/;
  var RE_TRAILING_SPACES = / +$/;
  var RE_BLANK_LINES = / *\r?\n/g;
  var RE_INDENT = /( *)$/;
  var TOKEN_BRACE_OPEN = /{\s*/y;
  var TOKEN_BRACE_CLOSE = /\s*}/y;
  var TOKEN_BRACKET_OPEN = /\[\s*/y;
  var TOKEN_BRACKET_CLOSE = /\s*] */y;
  var TOKEN_PAREN_OPEN = /\s*\(\s*/y;
  var TOKEN_ARROW = /\s*->\s*/y;
  var TOKEN_COLON = /\s*:\s*/y;
  var TOKEN_COMMA = /\s*,?\s*/y;
  var TOKEN_BLANK = /\s+/y;
  var FluentResource = class {
    constructor(source) {
      this.body = [];
      RE_MESSAGE_START.lastIndex = 0;
      let cursor = 0;
      while (true) {
        let next2 = RE_MESSAGE_START.exec(source);
        if (next2 === null) {
          break;
        }
        cursor = RE_MESSAGE_START.lastIndex;
        try {
          this.body.push(parseMessage(next2[1]));
        } catch (err) {
          if (err instanceof SyntaxError) {
            continue;
          }
          throw err;
        }
      }
      function test(re) {
        re.lastIndex = cursor;
        return re.test(source);
      }
      function consumeChar(char, errorClass) {
        if (source[cursor] === char) {
          cursor++;
          return true;
        }
        if (errorClass) {
          throw new errorClass(`Expected ${char}`);
        }
        return false;
      }
      function consumeToken(re, errorClass) {
        if (test(re)) {
          cursor = re.lastIndex;
          return true;
        }
        if (errorClass) {
          throw new errorClass(`Expected ${re.toString()}`);
        }
        return false;
      }
      function match2(re) {
        re.lastIndex = cursor;
        let result = re.exec(source);
        if (result === null) {
          throw new SyntaxError(`Expected ${re.toString()}`);
        }
        cursor = re.lastIndex;
        return result;
      }
      function match1(re) {
        return match2(re)[1];
      }
      function parseMessage(id) {
        let value = parsePattern();
        let attributes = parseAttributes();
        if (value === null && Object.keys(attributes).length === 0) {
          throw new SyntaxError("Expected message value or attributes");
        }
        return { id, value, attributes };
      }
      function parseAttributes() {
        let attrs2 = /* @__PURE__ */ Object.create(null);
        while (test(RE_ATTRIBUTE_START)) {
          let name = match1(RE_ATTRIBUTE_START);
          let value = parsePattern();
          if (value === null) {
            throw new SyntaxError("Expected attribute value");
          }
          attrs2[name] = value;
        }
        return attrs2;
      }
      function parsePattern() {
        let first;
        if (test(RE_TEXT_RUN)) {
          first = match1(RE_TEXT_RUN);
        }
        if (source[cursor] === "{" || source[cursor] === "}") {
          return parsePatternElements(first ? [first] : [], Infinity);
        }
        let indent = parseIndent();
        if (indent) {
          if (first) {
            return parsePatternElements([first, indent], indent.length);
          }
          indent.value = trim(indent.value, RE_LEADING_NEWLINES);
          return parsePatternElements([indent], indent.length);
        }
        if (first) {
          return trim(first, RE_TRAILING_SPACES);
        }
        return null;
      }
      function parsePatternElements(elements = [], commonIndent) {
        while (true) {
          if (test(RE_TEXT_RUN)) {
            elements.push(match1(RE_TEXT_RUN));
            continue;
          }
          if (source[cursor] === "{") {
            elements.push(parsePlaceable());
            continue;
          }
          if (source[cursor] === "}") {
            throw new SyntaxError("Unbalanced closing brace");
          }
          let indent = parseIndent();
          if (indent) {
            elements.push(indent);
            commonIndent = Math.min(commonIndent, indent.length);
            continue;
          }
          break;
        }
        let lastIndex = elements.length - 1;
        let lastElement = elements[lastIndex];
        if (typeof lastElement === "string") {
          elements[lastIndex] = trim(lastElement, RE_TRAILING_SPACES);
        }
        let baked = [];
        for (let element of elements) {
          if (element instanceof Indent) {
            element = element.value.slice(0, element.value.length - commonIndent);
          }
          if (element) {
            baked.push(element);
          }
        }
        return baked;
      }
      function parsePlaceable() {
        consumeToken(TOKEN_BRACE_OPEN, SyntaxError);
        let selector = parseInlineExpression();
        if (consumeToken(TOKEN_BRACE_CLOSE)) {
          return selector;
        }
        if (consumeToken(TOKEN_ARROW)) {
          let variants = parseVariants();
          consumeToken(TOKEN_BRACE_CLOSE, SyntaxError);
          return {
            type: "select",
            selector,
            ...variants
          };
        }
        throw new SyntaxError("Unclosed placeable");
      }
      function parseInlineExpression() {
        if (source[cursor] === "{") {
          return parsePlaceable();
        }
        if (test(RE_REFERENCE)) {
          let [, sigil, name, attr = null] = match2(RE_REFERENCE);
          if (sigil === "$") {
            return { type: "var", name };
          }
          if (consumeToken(TOKEN_PAREN_OPEN)) {
            let args = parseArguments();
            if (sigil === "-") {
              return { type: "term", name, attr, args };
            }
            if (RE_FUNCTION_NAME.test(name)) {
              return { type: "func", name, args };
            }
            throw new SyntaxError("Function names must be all upper-case");
          }
          if (sigil === "-") {
            return {
              type: "term",
              name,
              attr,
              args: []
            };
          }
          return { type: "mesg", name, attr };
        }
        return parseLiteral();
      }
      function parseArguments() {
        let args = [];
        while (true) {
          switch (source[cursor]) {
            case ")":
              cursor++;
              return args;
            case void 0:
              throw new SyntaxError("Unclosed argument list");
          }
          args.push(parseArgument());
          consumeToken(TOKEN_COMMA);
        }
      }
      function parseArgument() {
        let expr = parseInlineExpression();
        if (expr.type !== "mesg") {
          return expr;
        }
        if (consumeToken(TOKEN_COLON)) {
          return {
            type: "narg",
            name: expr.name,
            value: parseLiteral()
          };
        }
        return expr;
      }
      function parseVariants() {
        let variants = [];
        let count = 0;
        let star;
        while (test(RE_VARIANT_START)) {
          if (consumeChar("*")) {
            star = count;
          }
          let key = parseVariantKey();
          let value = parsePattern();
          if (value === null) {
            throw new SyntaxError("Expected variant value");
          }
          variants[count++] = { key, value };
        }
        if (count === 0) {
          return null;
        }
        if (star === void 0) {
          throw new SyntaxError("Expected default variant");
        }
        return { variants, star };
      }
      function parseVariantKey() {
        consumeToken(TOKEN_BRACKET_OPEN, SyntaxError);
        let key;
        if (test(RE_NUMBER_LITERAL)) {
          key = parseNumberLiteral();
        } else {
          key = {
            type: "str",
            value: match1(RE_IDENTIFIER)
          };
        }
        consumeToken(TOKEN_BRACKET_CLOSE, SyntaxError);
        return key;
      }
      function parseLiteral() {
        if (test(RE_NUMBER_LITERAL)) {
          return parseNumberLiteral();
        }
        if (source[cursor] === '"') {
          return parseStringLiteral();
        }
        throw new SyntaxError("Invalid expression");
      }
      function parseNumberLiteral() {
        let [, value, fraction = ""] = match2(RE_NUMBER_LITERAL);
        let precision = fraction.length;
        return {
          type: "num",
          value: parseFloat(value),
          precision
        };
      }
      function parseStringLiteral() {
        consumeChar('"', SyntaxError);
        let value = "";
        while (true) {
          value += match1(RE_STRING_RUN);
          if (source[cursor] === "\\") {
            value += parseEscapeSequence();
            continue;
          }
          if (consumeChar('"')) {
            return { type: "str", value };
          }
          throw new SyntaxError("Unclosed string literal");
        }
      }
      function parseEscapeSequence() {
        if (test(RE_STRING_ESCAPE)) {
          return match1(RE_STRING_ESCAPE);
        }
        if (test(RE_UNICODE_ESCAPE)) {
          let [, codepoint4, codepoint6] = match2(RE_UNICODE_ESCAPE);
          let codepoint = parseInt(codepoint4 || codepoint6, 16);
          return codepoint <= 55295 || 57344 <= codepoint ? (
            // It's a Unicode scalar value.
            String.fromCodePoint(codepoint)
          ) : (
            // Lonely surrogates can cause trouble when the parsing result is
            // saved using UTF-8. Use U+FFFD REPLACEMENT CHARACTER instead.
            "\uFFFD"
          );
        }
        throw new SyntaxError("Unknown escape sequence");
      }
      function parseIndent() {
        let start = cursor;
        consumeToken(TOKEN_BLANK);
        switch (source[cursor]) {
          case ".":
          case "[":
          case "*":
          case "}":
          case void 0:
            return false;
          case "{":
            return makeIndent(source.slice(start, cursor));
        }
        if (source[cursor - 1] === " ") {
          return makeIndent(source.slice(start, cursor));
        }
        return false;
      }
      function trim(text2, re) {
        return text2.replace(re, "");
      }
      function makeIndent(blank) {
        let value = blank.replace(RE_BLANK_LINES, "\n");
        let length = RE_INDENT.exec(blank)[1].length;
        return new Indent(value, length);
      }
    }
  };
  var Indent = class {
    constructor(value, length) {
      this.value = value;
      this.length = length;
    }
  };

  // src/assets/javascripts/i18n.ts
  var translations = {
    "unread": {
      "en": "Unread",
      "de": "Ungelesene",
      "fr": "Non lus",
      "es": "No le\xEDdos",
      "ja": "\u672A\u8AAD",
      "pt": "N\xE3o lidos",
      "zh": "\u672A\u8BFB",
      "ru": "\u041D\u0435\u043F\u0440\u043E\u0447\u0438\u0442\u0430\u043D\u043D\u044B\u0435"
    },
    "starred": {
      "en": "Starred",
      "de": "Markierte",
      "fr": "Favoris",
      "es": "Destacados",
      "ja": "\u30B9\u30BF\u30FC\u4ED8\u304D",
      "pt": "Favoritos",
      "zh": "\u661F\u6807",
      "ru": "\u0418\u0437\u0431\u0440\u0430\u043D\u043D\u044B\u0435"
    },
    "all": {
      "en": "All",
      "de": "Alle",
      "fr": "Tout",
      "es": "Todo",
      "ja": "\u3059\u3079\u3066",
      "pt": "Tudo",
      "zh": "\u5168\u90E8",
      "ru": "\u0412\u0441\u0435"
    },
    "settings": {
      "en": "Settings",
      "de": "Einstellungen",
      "fr": "Param\xE8tres",
      "es": "Ajustes",
      "ja": "\u8A2D\u5B9A",
      "pt": "Configura\xE7\xF5es",
      "zh": "\u8BBE\u7F6E",
      "ru": "\u041D\u0430\u0441\u0442\u0440\u043E\u0439\u043A\u0438"
    },
    "new_feed": {
      "en": "New Feed",
      "de": "Neuer Feed",
      "fr": "Nouveau flux",
      "es": "Nueva fuente",
      "ja": "\u65B0\u898F\u30D5\u30A3\u30FC\u30C9",
      "pt": "Novo feed",
      "zh": "\u65B0\u5EFA\u8BA2\u9605",
      "ru": "\u041D\u043E\u0432\u0430\u044F \u043B\u0435\u043D\u0442\u0430"
    },
    "refresh_feeds": {
      "en": "Refresh Feeds",
      "de": "Feeds aktualisieren",
      "fr": "Actualiser les flux",
      "es": "Actualizar fuentes",
      "ja": "\u30D5\u30A3\u30FC\u30C9\u3092\u66F4\u65B0",
      "pt": "Atualizar feeds",
      "zh": "\u5237\u65B0\u8BA2\u9605",
      "ru": "\u041E\u0431\u043D\u043E\u0432\u0438\u0442\u044C \u043B\u0435\u043D\u0442\u044B"
    },
    "theme": {
      "en": "Theme",
      "de": "Design",
      "fr": "Th\xE8me",
      "es": "Tema",
      "ja": "\u30C6\u30FC\u30DE",
      "pt": "Tema",
      "zh": "\u4E3B\u9898",
      "ru": "\u0422\u0435\u043C\u0430"
    },
    "auto_refresh": {
      "en": "Auto Refresh",
      "de": "Automatisch aktualisieren",
      "fr": "Actualisation automatique",
      "es": "Actualizaci\xF3n autom\xE1tica",
      "ja": "\u81EA\u52D5\u66F4\u65B0",
      "pt": "Atualiza\xE7\xE3o autom\xE1tica",
      "zh": "\u81EA\u52A8\u5237\u65B0",
      "ru": "\u0410\u0432\u0442\u043E\u043E\u0431\u043D\u043E\u0432\u043B\u0435\u043D\u0438\u0435"
    },
    "show_first": {
      "en": "Show first",
      "de": "Zuerst anzeigen",
      "fr": "Afficher d'abord",
      "es": "Mostrar primero",
      "ja": "\u8868\u793A\u9806",
      "pt": "Mostrar primeiro",
      "zh": "\u4F18\u5148\u663E\u793A",
      "ru": "\u0421\u043D\u0430\u0447\u0430\u043B\u0430"
    },
    "new": {
      "en": "New",
      "de": "Neue",
      "fr": "R\xE9cents",
      "es": "Nuevos",
      "ja": "\u65B0\u3057\u3044\u9806",
      "pt": "Novos",
      "zh": "\u6700\u65B0",
      "ru": "\u041D\u043E\u0432\u044B\u0435"
    },
    "old": {
      "en": "Old",
      "de": "Alte",
      "fr": "Anciens",
      "es": "Antiguos",
      "ja": "\u53E4\u3044\u9806",
      "pt": "Antigos",
      "zh": "\u6700\u65E7",
      "ru": "\u0421\u0442\u0430\u0440\u044B\u0435"
    },
    "subscriptions": {
      "en": "Subscriptions",
      "de": "Abonnements",
      "fr": "Abonnements",
      "es": "Suscripciones",
      "ja": "\u8CFC\u8AAD\u7BA1\u7406",
      "pt": "Assinaturas",
      "zh": "\u8BA2\u9605\u7BA1\u7406",
      "ru": "\u041F\u043E\u0434\u043F\u0438\u0441\u043A\u0438"
    },
    "import": {
      "en": "Import",
      "de": "Importieren",
      "fr": "Importer",
      "es": "Importar",
      "ja": "\u30A4\u30F3\u30DD\u30FC\u30C8",
      "pt": "Importar",
      "zh": "\u5BFC\u5165",
      "ru": "\u0418\u043C\u043F\u043E\u0440\u0442"
    },
    "export": {
      "en": "Export",
      "de": "Exportieren",
      "fr": "Exporter",
      "es": "Exportar",
      "ja": "\u30A8\u30AF\u30B9\u30DD\u30FC\u30C8",
      "pt": "Exportar",
      "zh": "\u5BFC\u51FA",
      "ru": "\u042D\u043A\u0441\u043F\u043E\u0440\u0442"
    },
    "shortcuts": {
      "en": "Shortcuts",
      "de": "Tastenk\xFCrzel",
      "fr": "Raccourcis",
      "es": "Atajos",
      "ja": "\u30B7\u30E7\u30FC\u30C8\u30AB\u30C3\u30C8",
      "pt": "Atalhos",
      "zh": "\u5FEB\u6377\u952E",
      "ru": "\u0413\u043E\u0440\u044F\u0447\u0438\u0435 \u043A\u043B\u0430\u0432\u0438\u0448\u0438"
    },
    "log_out": {
      "en": "Log out",
      "de": "Abmelden",
      "fr": "D\xE9connexion",
      "es": "Cerrar sesi\xF3n",
      "ja": "\u30ED\u30B0\u30A2\u30A6\u30C8",
      "pt": "Sair",
      "zh": "\u767B\u51FA",
      "ru": "\u0412\u044B\u0439\u0442\u0438"
    },
    "all_unread": {
      "en": "All Unread",
      "de": "Alle ungelesenen",
      "fr": "Tous les non lus",
      "es": "Todos los no le\xEDdos",
      "ja": "\u3059\u3079\u3066\u306E\u672A\u8AAD",
      "pt": "Todos os n\xE3o lidos",
      "zh": "\u5168\u90E8\u672A\u8BFB",
      "ru": "\u0412\u0441\u0435 \u043D\u0435\u043F\u0440\u043E\u0447\u0438\u0442\u0430\u043D\u043D\u044B\u0435"
    },
    "all_starred": {
      "en": "All Starred",
      "de": "Alle markierten",
      "fr": "Tous les favoris",
      "es": "Todos los destacados",
      "ja": "\u3059\u3079\u3066\u306E\u30B9\u30BF\u30FC\u4ED8\u304D",
      "pt": "Todos os favoritos",
      "zh": "\u5168\u90E8\u661F\u6807",
      "ru": "\u0412\u0441\u0435 \u0438\u0437\u0431\u0440\u0430\u043D\u043D\u044B\u0435"
    },
    "all_feeds": {
      "en": "All Feeds",
      "de": "Alle Feeds",
      "fr": "Tous les flux",
      "es": "Todas las fuentes",
      "ja": "\u3059\u3079\u3066\u306E\u30D5\u30A3\u30FC\u30C9",
      "pt": "Todos os feeds",
      "zh": "\u5168\u90E8\u8BA2\u9605",
      "ru": "\u0412\u0441\u0435 \u043B\u0435\u043D\u0442\u044B"
    },
    "refreshing_progress": {
      "en": "Refreshing ({ $count } left)",
      "de": "Aktualisiere ({ $count } \xFCbrig)",
      "fr": "Actualisation ({ $count } restantes)",
      "es": "Actualizando ({ $count } restantes)",
      "ja": "\u66F4\u65B0\u4E2D\uFF08\u6B8B\u308A{ $count }\uFF09",
      "pt": "Atualizando ({ $count } restantes)",
      "zh": "\u6B63\u5728\u5237\u65B0\uFF08\u5269\u4F59{ $count }\uFF09",
      "ru": "\u041E\u0431\u043D\u043E\u0432\u043B\u0435\u043D\u0438\u0435: \u043E\u0441\u0442\u0430\u043B\u043E\u0441\u044C { $count }"
    },
    "show_feeds": {
      "en": "Show Feeds",
      "de": "Feeds anzeigen",
      "fr": "Afficher les flux",
      "es": "Mostrar fuentes",
      "ja": "\u30D5\u30A3\u30FC\u30C9\u3092\u8868\u793A",
      "pt": "Mostrar feeds",
      "zh": "\u663E\u793A\u8BA2\u9605",
      "ru": "\u041F\u043E\u043A\u0430\u0437\u0430\u0442\u044C \u043B\u0435\u043D\u0442\u044B"
    },
    "mark_all_read": {
      "en": "Mark All Read",
      "de": "Alle als gelesen markieren",
      "fr": "Tout marquer comme lu",
      "es": "Marcar todo como le\xEDdo",
      "ja": "\u3059\u3079\u3066\u65E2\u8AAD\u306B\u3059\u308B",
      "pt": "Marcar todos como lidos",
      "zh": "\u5168\u90E8\u6807\u8BB0\u4E3A\u5DF2\u8BFB",
      "ru": "\u041E\u0442\u043C\u0435\u0442\u0438\u0442\u044C \u0432\u0441\u0435 \u043A\u0430\u043A \u043F\u0440\u043E\u0447\u0438\u0442\u0430\u043D\u043D\u044B\u0435"
    },
    "feed_settings": {
      "en": "Feed Settings",
      "de": "Feed-Einstellungen",
      "fr": "Param\xE8tres du flux",
      "es": "Ajustes de fuente",
      "ja": "\u30D5\u30A3\u30FC\u30C9\u8A2D\u5B9A",
      "pt": "Configura\xE7\xF5es do feed",
      "zh": "\u8BA2\u9605\u8BBE\u7F6E",
      "ru": "\u041D\u0430\u0441\u0442\u0440\u043E\u0439\u043A\u0438 \u043B\u0435\u043D\u0442\u044B"
    },
    "folder_settings": {
      "en": "Folder Settings",
      "de": "Ordner-Einstellungen",
      "fr": "Param\xE8tres du dossier",
      "es": "Ajustes de carpeta",
      "ja": "\u30D5\u30A9\u30EB\u30C0\u8A2D\u5B9A",
      "pt": "Configura\xE7\xF5es da pasta",
      "zh": "\u6587\u4EF6\u5939\u8BBE\u7F6E",
      "ru": "\u041D\u0430\u0441\u0442\u0440\u043E\u0439\u043A\u0438 \u043F\u0430\u043F\u043A\u0438"
    },
    "website": {
      "en": "Website",
      "de": "Webseite",
      "fr": "Site web",
      "es": "Sitio web",
      "ja": "\u30A6\u30A7\u30D6\u30B5\u30A4\u30C8",
      "pt": "Site",
      "zh": "\u7F51\u7AD9",
      "ru": "\u0421\u0430\u0439\u0442"
    },
    "feed_link": {
      "en": "Feed Link",
      "de": "Feed-Link",
      "fr": "Lien du flux",
      "es": "Enlace de la fuente",
      "ja": "\u30D5\u30A3\u30FC\u30C9\u30EA\u30F3\u30AF",
      "pt": "Link do feed",
      "zh": "\u8BA2\u9605\u94FE\u63A5",
      "ru": "\u0421\u0441\u044B\u043B\u043A\u0430 \u043D\u0430 \u043B\u0435\u043D\u0442\u0443"
    },
    "rename": {
      "en": "Rename",
      "de": "Umbenennen",
      "fr": "Renommer",
      "es": "Renombrar",
      "ja": "\u540D\u524D\u5909\u66F4",
      "pt": "Renomear",
      "zh": "\u91CD\u547D\u540D",
      "ru": "\u041F\u0435\u0440\u0435\u0438\u043C\u0435\u043D\u043E\u0432\u0430\u0442\u044C"
    },
    "change_link": {
      "en": "Change Link",
      "de": "Link \xE4ndern",
      "fr": "Changer le lien",
      "es": "Cambiar enlace",
      "ja": "\u30EA\u30F3\u30AF\u5909\u66F4",
      "pt": "Alterar link",
      "zh": "\u4FEE\u6539\u94FE\u63A5",
      "ru": "\u0418\u0437\u043C\u0435\u043D\u0438\u0442\u044C \u0441\u0441\u044B\u043B\u043A\u0443"
    },
    "move_to": {
      "en": "Move to...",
      "de": "Verschieben nach...",
      "fr": "D\xE9placer vers...",
      "es": "Mover a...",
      "ja": "\u79FB\u52D5...",
      "pt": "Mover para...",
      "zh": "\u79FB\u52A8\u5230...",
      "ru": "\u041F\u0435\u0440\u0435\u043C\u0435\u0441\u0442\u0438\u0442\u044C \u0432..."
    },
    "new_folder": {
      "en": "new folder",
      "de": "neuer Ordner",
      "fr": "nouveau dossier",
      "es": "nueva carpeta",
      "ja": "\u65B0\u898F\u30D5\u30A9\u30EB\u30C0",
      "pt": "nova pasta",
      "zh": "\u65B0\u5EFA\u6587\u4EF6\u5939",
      "ru": "\u043D\u043E\u0432\u0430\u044F \u043F\u0430\u043F\u043A\u0430"
    },
    "delete": {
      "en": "Delete",
      "de": "L\xF6schen",
      "fr": "Supprimer",
      "es": "Eliminar",
      "ja": "\u524A\u9664",
      "pt": "Excluir",
      "zh": "\u5220\u9664",
      "ru": "\u0423\u0434\u0430\u043B\u0438\u0442\u044C"
    },
    "mark_starred": {
      "en": "Mark Starred",
      "de": "Als markiert kennzeichnen",
      "fr": "Marquer comme favori",
      "es": "Marcar como destacado",
      "ja": "\u30B9\u30BF\u30FC\u3092\u4ED8\u3051\u308B",
      "pt": "Marcar como favorito",
      "zh": "\u6807\u8BB0\u661F\u6807",
      "ru": "\u041F\u043E\u043C\u0435\u0442\u0438\u0442\u044C \u0438\u0437\u0431\u0440\u0430\u043D\u043D\u044B\u043C"
    },
    "mark_unread": {
      "en": "Mark Unread",
      "de": "Als ungelesen kennzeichnen",
      "fr": "Marquer comme non lu",
      "es": "Marcar como no le\xEDdo",
      "ja": "\u672A\u8AAD\u306B\u3059\u308B",
      "pt": "Marcar como n\xE3o lido",
      "zh": "\u6807\u8BB0\u672A\u8BFB",
      "ru": "\u041F\u043E\u043C\u0435\u0442\u0438\u0442\u044C \u043D\u0435\u043F\u0440\u043E\u0447\u0438\u0442\u0430\u043D\u043D\u044B\u043C"
    },
    "appearance": {
      "en": "Appearance",
      "de": "Darstellung",
      "fr": "Apparence",
      "es": "Apariencia",
      "ja": "\u8868\u793A\u8A2D\u5B9A",
      "pt": "Apar\xEAncia",
      "zh": "\u5916\u89C2",
      "ru": "\u0412\u043D\u0435\u0448\u043D\u0438\u0439 \u0432\u0438\u0434"
    },
    "read_here": {
      "en": "Read Here",
      "de": "Hier lesen",
      "fr": "Lire ici",
      "es": "Leer aqu\xED",
      "ja": "\u3053\u3053\u3067\u8AAD\u3080",
      "pt": "Ler aqui",
      "zh": "\u5728\u6B64\u9605\u8BFB",
      "ru": "\u0427\u0438\u0442\u0430\u0442\u044C \u0437\u0434\u0435\u0441\u044C"
    },
    "open_link": {
      "en": "Open Link",
      "de": "Link \xF6ffnen",
      "fr": "Ouvrir le lien",
      "es": "Abrir enlace",
      "ja": "\u30EA\u30F3\u30AF\u3092\u958B\u304F",
      "pt": "Abrir link",
      "zh": "\u6253\u5F00\u94FE\u63A5",
      "ru": "\u041E\u0442\u043A\u0440\u044B\u0442\u044C \u0441\u0441\u044B\u043B\u043A\u0443"
    },
    "previous_article": {
      "en": "Previous Article",
      "de": "Vorheriger Artikel",
      "fr": "Article pr\xE9c\xE9dent",
      "es": "Art\xEDculo anterior",
      "ja": "\u524D\u306E\u8A18\u4E8B",
      "pt": "Artigo anterior",
      "zh": "\u4E0A\u4E00\u7BC7",
      "ru": "\u041F\u0440\u0435\u0434\u044B\u0434\u0443\u0449\u0430\u044F \u0441\u0442\u0430\u0442\u044C\u044F"
    },
    "next_article": {
      "en": "Next Article",
      "de": "N\xE4chster Artikel",
      "fr": "Article suivant",
      "es": "Art\xEDculo siguiente",
      "ja": "\u6B21\u306E\u8A18\u4E8B",
      "pt": "Pr\xF3ximo artigo",
      "zh": "\u4E0B\u4E00\u7BC7",
      "ru": "\u0421\u043B\u0435\u0434\u0443\u044E\u0449\u0430\u044F \u0441\u0442\u0430\u0442\u044C\u044F"
    },
    "close_article": {
      "en": "Close Article",
      "de": "Artikel schlie\xDFen",
      "fr": "Fermer l'article",
      "es": "Cerrar art\xEDculo",
      "ja": "\u8A18\u4E8B\u3092\u9589\u3058\u308B",
      "pt": "Fechar artigo",
      "zh": "\u5173\u95ED\u6587\u7AE0",
      "ru": "\u0417\u0430\u043A\u0440\u044B\u0442\u044C \u0441\u0442\u0430\u0442\u044C\u044E"
    },
    "untitled": {
      "en": "untitled",
      "de": "unbenannt",
      "fr": "sans titre",
      "es": "sin t\xEDtulo",
      "ja": "\u7121\u984C",
      "pt": "sem t\xEDtulo",
      "zh": "\u65E0\u6807\u9898",
      "ru": "\u0431\u0435\u0437 \u043D\u0430\u0437\u0432\u0430\u043D\u0438\u044F"
    },
    "sans_serif": {
      "en": "sans-serif",
      "de": "serifenlos",
      "fr": "sans empattement",
      "es": "sans-serif",
      "ja": "\u30B4\u30B7\u30C3\u30AF\u4F53",
      "pt": "sem serifa",
      "zh": "\u65E0\u886C\u7EBF",
      "ru": "sans-serif"
    },
    "serif": {
      "en": "serif",
      "de": "Serife",
      "fr": "empattement",
      "es": "serifa",
      "ja": "\u660E\u671D\u4F53",
      "pt": "com serifa",
      "zh": "\u886C\u7EBF",
      "ru": "serif"
    },
    "monospace": {
      "en": "monospace",
      "de": "monospace",
      "fr": "monospace",
      "es": "monoespacio",
      "ja": "\u7B49\u5E45",
      "pt": "monoespa\xE7ada",
      "zh": "\u7B49\u5BBD",
      "ru": "monospace"
    },
    "url": {
      "en": "URL",
      "de": "URL",
      "fr": "URL",
      "es": "URL",
      "ja": "URL",
      "pt": "URL",
      "zh": "\u7F51\u5740",
      "ru": "URL"
    },
    "folder": {
      "en": "Folder",
      "de": "Ordner",
      "fr": "Dossier",
      "es": "Carpeta",
      "ja": "\u30D5\u30A9\u30EB\u30C0",
      "pt": "Pasta",
      "zh": "\u6587\u4EF6\u5939",
      "ru": "\u041F\u0430\u043F\u043A\u0430"
    },
    "add": {
      "en": "Add",
      "de": "Hinzuf\xFCgen",
      "fr": "Ajouter",
      "es": "A\xF1adir",
      "ja": "\u8FFD\u52A0",
      "pt": "Adicionar",
      "zh": "\u6DFB\u52A0",
      "ru": "\u0414\u043E\u0431\u0430\u0432\u0438\u0442\u044C"
    },
    "keyboard_shortcuts": {
      "en": "Keyboard Shortcuts",
      "de": "Tastenk\xFCrzel",
      "fr": "Raccourcis clavier",
      "es": "Atajos de teclado",
      "ja": "\u30AD\u30FC\u30DC\u30FC\u30C9\u30B7\u30E7\u30FC\u30C8\u30AB\u30C3\u30C8",
      "pt": "Atalhos do teclado",
      "zh": "\u952E\u76D8\u5FEB\u6377\u952E",
      "ru": "\u0413\u043E\u0440\u044F\u0447\u0438\u0435 \u043A\u043B\u0430\u0432\u0438\u0448\u0438"
    },
    "multiple_feeds_found": {
      "en": "Multiple feeds found. Choose one below:",
      "de": "Mehrere Feeds gefunden. Bitte w\xE4hlen Sie einen aus:",
      "fr": "Plusieurs flux trouv\xE9s. Choisissez-en un ci-dessous :",
      "es": "M\xFAltiples fuentes encontradas. Elija una:",
      "ja": "\u8907\u6570\u306E\u30D5\u30A3\u30FC\u30C9\u304C\u898B\u3064\u304B\u308A\u307E\u3057\u305F\u3002\u4EE5\u4E0B\u304B\u3089\u9078\u629E\u3057\u3066\u304F\u3060\u3055\u3044\uFF1A",
      "pt": "M\xFAltiplos feeds encontrados. Escolha um abaixo:",
      "zh": "\u627E\u5230\u591A\u4E2A\u8BA2\u9605\u6E90\uFF0C\u8BF7\u9009\u62E9\u4E00\u4E2A\uFF1A",
      "ru": "\u041D\u0430\u0439\u0434\u0435\u043D\u043E \u043D\u0435\u0441\u043A\u043E\u043B\u044C\u043A\u043E \u043B\u0435\u043D\u0442. \u0412\u044B\u0431\u0435\u0440\u0438\u0442\u0435 \u043E\u0434\u043D\u0443:"
    },
    "cancel": {
      "en": "cancel",
      "de": "abbrechen",
      "fr": "annuler",
      "es": "cancelar",
      "ja": "\u30AD\u30E3\u30F3\u30BB\u30EB",
      "pt": "cancelar",
      "zh": "\u53D6\u6D88",
      "ru": "\u043E\u0442\u043C\u0435\u043D\u0430"
    },
    "kb_show_filters": {
      "en": "show unread / starred / all feeds",
      "de": "ungelesene / markierte / alle Feeds anzeigen",
      "fr": "afficher les flux non lus / favoris / tous",
      "es": "mostrar fuentes no le\xEDdas / destacadas / todas",
      "ja": "\u672A\u8AAD/\u30B9\u30BF\u30FC\u4ED8\u304D/\u3059\u3079\u3066\u306E\u30D5\u30A3\u30FC\u30C9\u3092\u8868\u793A",
      "pt": "mostrar feeds n\xE3o lidos / favoritos / todos",
      "zh": "\u663E\u793A\u672A\u8BFB/\u661F\u6807/\u5168\u90E8\u8BA2\u9605",
      "ru": "\u043F\u043E\u043A\u0430\u0437\u0430\u0442\u044C \u043D\u0435\u043F\u0440\u043E\u0447\u0438\u0442\u0430\u043D\u043D\u044B\u0435 / \u0438\u0437\u0431\u0440\u0430\u043D\u043D\u044B\u0435 / \u0432\u0441\u0435 \u043B\u0435\u043D\u0442\u044B"
    },
    "kb_focus_search": {
      "en": "focus the search bar",
      "de": "Suchleiste fokussieren",
      "fr": "focus sur la barre de recherche",
      "es": "enfocar la barra de b\xFAsqueda",
      "ja": "\u691C\u7D22\u30D0\u30FC\u306B\u30D5\u30A9\u30FC\u30AB\u30B9",
      "pt": "focar na barra de pesquisa",
      "zh": "\u805A\u7126\u641C\u7D22\u680F",
      "ru": "\u0444\u043E\u043A\u0443\u0441 \u043D\u0430 \u0441\u0442\u0440\u043E\u043A\u0443 \u043F\u043E\u0438\u0441\u043A\u0430"
    },
    "kb_next_prev_article": {
      "en": "next / prev article",
      "de": "n\xE4chster / vorheriger Artikel",
      "fr": "article suivant / pr\xE9c\xE9dent",
      "es": "art\xEDculo siguiente / anterior",
      "ja": "\u6B21\u306E/\u524D\u306E\u8A18\u4E8B",
      "pt": "pr\xF3ximo / artigo anterior",
      "zh": "\u4E0B\u4E00\u7BC7/\u4E0A\u4E00\u7BC7\u6587\u7AE0",
      "ru": "\u0441\u043B\u0435\u0434\u0443\u044E\u0449\u0430\u044F / \u043F\u0440\u0435\u0434\u044B\u0434\u0443\u0449\u0430\u044F \u0441\u0442\u0430\u0442\u044C\u044F"
    },
    "kb_next_prev_feed": {
      "en": "next / prev feed",
      "de": "n\xE4chster / vorheriger Feed",
      "fr": "flux suivant / pr\xE9c\xE9dent",
      "es": "fuente siguiente / anterior",
      "ja": "\u6B21\u306E/\u524D\u306E\u30D5\u30A3\u30FC\u30C9",
      "pt": "pr\xF3ximo / feed anterior",
      "zh": "\u4E0B\u4E00\u4E2A/\u4E0A\u4E00\u4E2A\u8BA2\u9605",
      "ru": "\u0441\u043B\u0435\u0434\u0443\u044E\u0449\u0430\u044F / \u043F\u0440\u0435\u0434\u044B\u0434\u0443\u0449\u0430\u044F \u043B\u0435\u043D\u0442\u0430"
    },
    "kb_close_article": {
      "en": "close article",
      "de": "Artikel schlie\xDFen",
      "fr": "fermer l'article",
      "es": "cerrar art\xEDculo",
      "ja": "\u8A18\u4E8B\u3092\u9589\u3058\u308B",
      "pt": "fechar artigo",
      "zh": "\u5173\u95ED\u6587\u7AE0",
      "ru": "\u0437\u0430\u043A\u0440\u044B\u0442\u044C \u0441\u0442\u0430\u0442\u044C\u044E"
    },
    "kb_mark_all_read": {
      "en": "mark all read",
      "de": "alle als gelesen markieren",
      "fr": "tout marquer comme lu",
      "es": "marcar todo como le\xEDdo",
      "ja": "\u3059\u3079\u3066\u65E2\u8AAD\u306B\u3059\u308B",
      "pt": "marcar todos como lidos",
      "zh": "\u5168\u90E8\u6807\u8BB0\u4E3A\u5DF2\u8BFB",
      "ru": "\u043E\u0442\u043C\u0435\u0442\u0438\u0442\u044C \u0432\u0441\u0435 \u043A\u0430\u043A \u043F\u0440\u043E\u0447\u0438\u0442\u0430\u043D\u043D\u044B\u0435"
    },
    "kb_mark_read": {
      "en": "mark read / unread",
      "de": "als gelesen / ungelesen markieren",
      "fr": "marquer comme lu / non lu",
      "es": "marcar como le\xEDdo / no le\xEDdo",
      "ja": "\u65E2\u8AAD/\u672A\u8AAD\u3092\u5207\u308A\u66FF\u3048",
      "pt": "marcar como lido / n\xE3o lido",
      "zh": "\u6807\u8BB0\u5DF2\u8BFB/\u672A\u8BFB",
      "ru": "\u043E\u0442\u043C\u0435\u0442\u0438\u0442\u044C \u043A\u0430\u043A \u043F\u0440\u043E\u0447\u0438\u0442\u0430\u043D\u043D\u043E\u0435 / \u043D\u0435\u043F\u0440\u043E\u0447\u0438\u0442\u0430\u043D\u043D\u043E\u0435"
    },
    "kb_mark_starred": {
      "en": "mark starred / unstarred",
      "de": "als markiert / nicht markiert kennzeichnen",
      "fr": "marquer comme favori / non favori",
      "es": "marcar como destacado / no destacado",
      "ja": "\u30B9\u30BF\u30FC\u3092\u4ED8\u3051\u308B/\u5916\u3059",
      "pt": "marcar como favorito / n\xE3o favorito",
      "zh": "\u6807\u8BB0\u661F\u6807/\u53D6\u6D88\u661F\u6807",
      "ru": "\u043F\u043E\u043C\u0435\u0442\u0438\u0442\u044C \u0438\u0437\u0431\u0440\u0430\u043D\u043D\u044B\u043C / \u0443\u0431\u0440\u0430\u0442\u044C \u0438\u0437 \u0438\u0437\u0431\u0440\u0430\u043D\u043D\u043E\u0433\u043E"
    },
    "kb_open_link": {
      "en": "open link",
      "de": "Link \xF6ffnen",
      "fr": "ouvrir le lien",
      "es": "abrir enlace",
      "ja": "\u30EA\u30F3\u30AF\u3092\u958B\u304F",
      "pt": "abrir link",
      "zh": "\u6253\u5F00\u94FE\u63A5",
      "ru": "\u043E\u0442\u043A\u0440\u044B\u0442\u044C \u0441\u0441\u044B\u043B\u043A\u0443"
    },
    "kb_read_here": {
      "en": "read here",
      "de": "hier lesen",
      "fr": "lire ici",
      "es": "leer aqu\xED",
      "ja": "\u3053\u3053\u3067\u8AAD\u3080",
      "pt": "ler aqui",
      "zh": "\u5728\u6B64\u9605\u8BFB",
      "ru": "\u0447\u0438\u0442\u0430\u0442\u044C \u0437\u0434\u0435\u0441\u044C"
    },
    "kb_scroll_content": {
      "en": "scroll content forward / backward",
      "de": "Inhalt vorw\xE4rts / r\xFCckw\xE4rts scrollen",
      "fr": "faire d\xE9filer le contenu avant / arri\xE8re",
      "es": "desplazar contenido hacia adelante / atr\xE1s",
      "ja": "\u30B3\u30F3\u30C6\u30F3\u30C4\u3092\u524D/\u5F8C\u306B\u30B9\u30AF\u30ED\u30FC\u30EB",
      "pt": "rolar conte\xFAdo para frente / tr\xE1s",
      "zh": "\u5411\u524D/\u5411\u540E\u6EDA\u52A8\u5185\u5BB9",
      "ru": "\u043F\u0440\u043E\u043A\u0440\u0443\u0442\u043A\u0430 \u0432\u043F\u0435\u0440\u0435\u0434 / \u043D\u0430\u0437\u0430\u0434"
    },
    "prompt_folder_name": {
      "en": "Enter folder name:",
      "de": "Ordnernamen eingeben:",
      "fr": "Entrez le nom du dossier :",
      "es": "Introduzca el nombre de la carpeta:",
      "ja": "\u30D5\u30A9\u30EB\u30C0\u540D\u3092\u5165\u529B\u3057\u3066\u304F\u3060\u3055\u3044\uFF1A",
      "pt": "Digite o nome da pasta:",
      "zh": "\u8BF7\u8F93\u5165\u6587\u4EF6\u5939\u540D\u79F0\uFF1A",
      "ru": "\u0412\u0432\u0435\u0434\u0438\u0442\u0435 \u0438\u043C\u044F \u043F\u0430\u043F\u043A\u0438:"
    },
    "prompt_new_title": {
      "en": "Enter new title",
      "de": "Neuen Titel eingeben",
      "fr": "Entrez un nouveau titre",
      "es": "Introduzca un nuevo t\xEDtulo",
      "ja": "\u65B0\u3057\u3044\u30BF\u30A4\u30C8\u30EB\u3092\u5165\u529B\u3057\u3066\u304F\u3060\u3055\u3044",
      "pt": "Digite o novo t\xEDtulo",
      "zh": "\u8BF7\u8F93\u5165\u65B0\u6807\u9898",
      "ru": "\u0412\u0432\u0435\u0434\u0438\u0442\u0435 \u043D\u043E\u0432\u044B\u0439 \u0437\u0430\u0433\u043E\u043B\u043E\u0432\u043E\u043A"
    },
    "prompt_feed_link": {
      "en": "Enter feed link",
      "de": "Feed-Link eingeben",
      "fr": "Entrez le lien du flux",
      "es": "Introduzca el enlace de la fuente",
      "ja": "\u30D5\u30A3\u30FC\u30C9\u30EA\u30F3\u30AF\u3092\u5165\u529B\u3057\u3066\u304F\u3060\u3055\u3044",
      "pt": "Digite o link do feed",
      "zh": "\u8BF7\u8F93\u5165\u8BA2\u9605\u94FE\u63A5",
      "ru": "\u0412\u0432\u0435\u0434\u0438\u0442\u0435 \u0441\u0441\u044B\u043B\u043A\u0443 \u043D\u0430 \u043B\u0435\u043D\u0442\u0443"
    },
    "confirm_delete": {
      "en": "Are you sure you want to delete { $name }?",
      "de": "M\xF6chten Sie { $name } wirklich l\xF6schen?",
      "fr": "Voulez-vous vraiment supprimer { $name } ?",
      "es": "\xBFEst\xE1 seguro de que quiere eliminar { $name }?",
      "ja": "{ $name }\u3092\u524A\u9664\u3057\u3066\u3082\u3088\u308D\u3057\u3044\u3067\u3059\u304B\uFF1F",
      "pt": "Tem certeza que deseja excluir { $name }?",
      "zh": "\u786E\u5B9A\u8981\u5220\u9664{ $name }\uFF1F",
      "ru": "\u0412\u044B \u0443\u0432\u0435\u0440\u0435\u043D\u044B, \u0447\u0442\u043E \u0445\u043E\u0442\u0438\u0442\u0435 \u0443\u0434\u0430\u043B\u0438\u0442\u044C { $name }?"
    },
    "alert_no_feeds": {
      "en": "No feeds found at the given url.",
      "de": "Keine Feeds unter der angegebenen URL gefunden.",
      "fr": "Aucun flux trouv\xE9 \xE0 cette URL.",
      "es": "No se encontraron fuentes en la URL proporcionada.",
      "ja": "\u6307\u5B9A\u3055\u308C\u305FURL\u306B\u30D5\u30A3\u30FC\u30C9\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3067\u3057\u305F\u3002",
      "pt": "Nenhum feed encontrado no URL fornecido.",
      "zh": "\u5728\u6307\u5B9A\u7684\u7F51\u5740\u672A\u627E\u5230\u8BA2\u9605\u6E90\u3002",
      "ru": "\u041B\u0435\u043D\u0442 \u043F\u043E \u0434\u0430\u043D\u043D\u043E\u043C\u0443 \u0430\u0434\u0440\u0435\u0441\u0443 \u043D\u0435 \u043D\u0430\u0439\u0434\u0435\u043D\u043E."
    },
    "login": {
      "en": "Login",
      "de": "Anmelden",
      "fr": "Connexion",
      "es": "Iniciar sesi\xF3n",
      "ja": "\u30ED\u30B0\u30A4\u30F3",
      "pt": "Entrar",
      "zh": "\u767B\u5F55",
      "ru": "\u0412\u0445\u043E\u0434"
    },
    "login_error": {
      "en": "Invalid username or password",
      "de": "Ung\xFCltiger Benutzername oder Passwort",
      "fr": "Nom d'utilisateur ou mot de passe invalide",
      "es": "Nombre de usuario o contrase\xF1a inv\xE1lidos",
      "ja": "\u30E6\u30FC\u30B6\u30FC\u540D\u307E\u305F\u306F\u30D1\u30B9\u30EF\u30FC\u30C9\u304C\u7121\u52B9\u3067\u3059",
      "pt": "Nome de usu\xE1rio ou senha inv\xE1lidos",
      "zh": "\u7528\u6237\u540D\u6216\u5BC6\u7801\u9519\u8BEF",
      "ru": "\u041D\u0435\u0432\u0435\u0440\u043D\u043E\u0435 \u0438\u043C\u044F \u043F\u043E\u043B\u044C\u0437\u043E\u0432\u0430\u0442\u0435\u043B\u044F \u0438\u043B\u0438 \u043F\u0430\u0440\u043E\u043B\u044C"
    },
    "username": {
      "en": "Username",
      "de": "Benutzername",
      "fr": "Nom d'utilisateur",
      "es": "Nombre de usuario",
      "ja": "\u30E6\u30FC\u30B6\u30FC\u540D",
      "pt": "Nome de usu\xE1rio",
      "zh": "\u7528\u6237\u540D",
      "ru": "\u0418\u043C\u044F \u043F\u043E\u043B\u044C\u0437\u043E\u0432\u0430\u0442\u0435\u043B\u044F"
    },
    "password": {
      "en": "Password",
      "de": "Passwort",
      "fr": "Mot de passe",
      "es": "Contrase\xF1a",
      "ja": "\u30D1\u30B9\u30EF\u30FC\u30C9",
      "pt": "Senha",
      "zh": "\u5BC6\u7801",
      "ru": "\u041F\u0430\u0440\u043E\u043B\u044C"
    }
  };
  function ftlFrom(lang) {
    return Object.entries(translations).map(([key, langs]) => `${key} = ${langs[lang]}`).join("\n");
  }
  var i18n_default = {
    install(Vue2) {
      let bundle = null;
      Vue2.prototype.$setLang = function(lang) {
        const ftl = ftlFrom(lang);
        const resource = new FluentResource(ftl);
        bundle = new FluentBundle(lang);
        bundle.addResource(resource);
      };
      Vue2.prototype.$t = function(code, args) {
        if (!bundle) return;
        const msg = bundle.getMessage(code);
        if (!msg || !msg.value) return;
        return bundle.formatPattern(msg.value, args);
      };
    }
  };

  // src/assets/javascripts/api.ts
  var xfetch = function(resource, init) {
    init = init || {};
    if (["post", "put", "delete"].indexOf(init.method) !== -1) {
      init["headers"] = init["headers"] || {};
      init["headers"]["x-requested-by"] = "yarr";
    }
    return fetch(resource, init);
  };
  var api = function(method, endpoint, data) {
    var headers = { "Content-Type": "application/json" };
    return xfetch(endpoint, {
      method,
      headers,
      body: JSON.stringify(data)
    });
  };
  var json = function(res) {
    return res.json();
  };
  var param = function(query2) {
    if (!query2) return "";
    return "?" + Object.keys(query2).map(function(key) {
      return encodeURIComponent(key) + "=" + encodeURIComponent(query2[key]);
    }).join("&");
  };
  var api_default = {
    feeds: {
      list: function() {
        return api("get", "./api/feeds").then(json);
      },
      create: function(data) {
        return api("post", "./api/feeds", data).then(json);
      },
      update: function(id, data) {
        return api("put", "./api/feeds/" + id, data);
      },
      delete: function(id) {
        return api("delete", "./api/feeds/" + id);
      },
      list_items: function(id) {
        return api("get", "./api/feeds/" + id + "/items").then(json);
      },
      refresh: function() {
        return api("post", "./api/feeds/refresh");
      },
      list_errors: function() {
        return api("get", "./api/feeds/errors").then(json);
      }
    },
    folders: {
      list: function() {
        return api("get", "./api/folders").then(json);
      },
      create: function(data) {
        return api("post", "./api/folders", data).then(json);
      },
      update: function(id, data) {
        return api("put", "./api/folders/" + id, data);
      },
      delete: function(id) {
        return api("delete", "./api/folders/" + id);
      },
      list_items: function(id) {
        return api("get", "./api/folders/" + id + "/items").then(json);
      }
    },
    items: {
      get: function(id) {
        return api("get", "./api/items/" + id).then(json);
      },
      list: function(query2) {
        return api("get", "./api/items" + param(query2)).then(json);
      },
      update: function(id, data) {
        return api("put", "./api/items/" + id, data);
      },
      mark_read: function(query2) {
        return api("put", "./api/items" + param(query2));
      }
    },
    settings: {
      get: function() {
        return api("get", "./api/settings").then(json);
      },
      update: function(data) {
        return api("put", "./api/settings", data);
      }
    },
    status: function() {
      return api("get", "./api/status").then(json);
    },
    upload_opml: function(form) {
      return xfetch("./opml/import", {
        method: "post",
        body: new FormData(form)
      });
    },
    logout: function() {
      return api("post", "./logout");
    },
    crawl: function(url) {
      return api("get", "./page?url=" + encodeURIComponent(url)).then(json);
    }
  };

  // src/assets/javascripts/templates/index.html
  var templates_default = `<div class="d-flex" :class="{'feed-selected': feedSelected !== null, 'item-selected': itemSelected !== null}">
    <!-- feed list -->
    <div id="col-feed-list" class="vh-100 position-relative d-flex flex-column border-right flex-shrink-0" :style="{width: feedListWidth+'px'}">
        <drag :width="feedListWidth" @resize="resizeFeedList"></drag>
        <div class="p-2 toolbar d-flex align-items-center">
            <v-icon class="mx-2" name="anchor" />
            <div class="flex-grow-1"></div>
            <button class="toolbar-item ml-1"
                    :class="{active: filterSelected == 'unread'}"
                    :aria-pressed="filterSelected == 'unread'"
                    :title="$t('unread')"
                    @click="filterSelected = 'unread'">
                <v-icon name="circle-full" />
            </button>
            <button class="toolbar-item mx-1"
                    :class="{active: filterSelected == 'starred'}"
                    :aria-pressed="filterSelected == 'starred'"
                    :title="$t('starred')"
                    @click="filterSelected = 'starred'">
                <v-icon name="star-full" />
            </button>
            <button class="toolbar-item mr-1"
                    :class="{active: filterSelected == ''}"
                    :aria-pressed="filterSelected == ''"
                    :title="$t('all')"
                    @click="filterSelected = ''">
                <v-icon name="assorted" />
            </button>
            <div class="flex-grow-1"></div>
            <dropdown class="settings-dropdown" toggle-class="btn btn-link toolbar-item px-2" ref="menuDropdown" drop="right" :title="$t('settings')">
                <template v-slot:button>
                    <v-icon name="more-horizontal" />
                </template>

                <button class="dropdown-item" @click="showSettings('create')">
                    <v-icon class="mr-1" name="plus" />
                    {{ $t('new_feed') }}
                </button>
                <div class="dropdown-divider"></div>
                <button class="dropdown-item" @click="fetchAllFeeds()">
                    <v-icon class="mr-1" name="rotate-cw" />
                    {{ $t('refresh_feeds') }}
                </button>

                <div class="dropdown-divider"></div>

                <header class="dropdown-header" role="heading" aria-level="2">{{ $t('theme') }}</header>
                <div class="row text-center m-0">
                    <button class="btn btn-link col-4 px-0 rounded-0"
                            :class="'theme-'+t"
                            :title="t"
                            :aria-label="t"
                            :aria-pressed="theme.name == t"
                            @click.stop="theme.name = t"
                            v-for="t in ['light', 'sepia', 'night']">
                        <v-icon v-if="theme.name == t" name="check" />
                    </button>
                </div>

                <div class="dropdown-divider"></div>

                <header class="dropdown-header" role="heading" aria-level="2">{{ $t('auto_refresh') }}</header>
                <div class="row text-center m-0">
                    <button class="dropdown-item col-4 px-0"
                            @click.stop="changeRefreshRate(-1)"
                            :disabled="!refreshRate">
                        <v-icon name="chevron-down" />
                    </button>
                    <div class="col-4 d-flex align-items-center justify-content-center">{{ refreshRateTitle }}</div>
                    <button class="dropdown-item col-4 px-0"
                            @click.stop="changeRefreshRate(1)" :disabled="refreshRate === refreshRateOptions.at(-1).value">
                        <v-icon name="chevron-up" />
                    </button>
                </div>

                <div class="dropdown-divider"></div>

                <header class="dropdown-header" role="heading" aria-level="2">{{ $t('show_first') }}</header>
                <div class="d-flex text-center">
                    <button class="dropdown-item px-0" :aria-pressed="itemSortNewestFirst"  :class="{active: itemSortNewestFirst}"  @click.stop="itemSortNewestFirst=true">{{ $t('new') }}</button>
                    <button class="dropdown-item px-0" :aria-pressed="!itemSortNewestFirst" :class="{active: !itemSortNewestFirst}" @click.stop="itemSortNewestFirst=false">{{ $t('old') }}</button>
                </div>
                <div class="dropdown-divider"></div>
                <header class="dropdown-header" role="heading" aria-level="2">{{ $t('subscriptions') }}</header>
                <form id="opml-import-form" enctype="multipart/form-data" tabindex="-1">
                    <input type="file"
                            id="opml-import"
                            @change="importOPML"
                            name="opml"
                            style="opacity: 0; width: 1px; height: 0; position: absolute; z-index: -1;">
                    <label class="dropdown-item mb-0 cursor-pointer" for="opml-import" @click.stop="">
                        <v-icon class="mr-1" name="download" />
                        {{ $t('import') }}
                    </label>
                </form>
                <a class="dropdown-item" href="./opml/export">
                    <v-icon class="mr-1" name="upload" />
                    {{ $t('export') }}
                </a>
                <div class="dropdown-divider"></div>
                <button class="dropdown-item" @click="showSettings('shortcuts')">
                    <v-icon class="mr-1" name="help-circle" />
                    {{ $t('shortcuts') }}
                </button>
                <div class="dropdown-divider"></div>
                <header class="dropdown-header" role="heading" aria-level="2">A / \u3042 / \u6587</header>
                <div class="container">
                    <div class="row">
                        <button
                            v-for="lang in languages"
                            class="dropdown-item text-center col-3 px-0"
                            :aria-label="lang.name"
                            :title="lang.name"
                            :class="{active: language==lang.code}"
                            @click.stop="changeLanguage(lang.code)">
                                {{ lang.code }}
                        </button>
                    </div>
                </div>
                <div class="dropdown-divider" v-if="authenticated"></div>
                <button class="dropdown-item" v-if="authenticated" @click="logout()">
                    <v-icon class="mr-1" name="log-out" />
                    {{ $t('log_out') }}
                </button>
            </dropdown>
        </div>
        <div id="feed-list-scroll" class="p-2 overflow-auto scroll-touch border-top flex-grow-1">
            <label class="selectgroup">
                <input type="radio" name="feed" value="" v-model="feedSelected">
                <div class="selectgroup-label d-flex align-items-center w-100">
                    <v-icon class="mr-2" name="layers" />
                    <span class="flex-fill text-left text-truncate" v-if="filterSelected=='unread'">{{ $t('all_unread') }}</span>
                    <span class="flex-fill text-left text-truncate" v-if="filterSelected=='starred'">{{ $t('all_starred') }}</span>
                    <span class="flex-fill text-left text-truncate" v-if="filterSelected==''">{{ $t('all_feeds') }}</span>
                    <span class="counter text-right">{{ filteredTotalStats }}</span>
                </div>
            </label>
            <div v-for="folder in foldersWithFeeds">
                <label class="selectgroup mt-1"
                        :class="{'d-none': mustHideFolder(folder)}"
                        v-if="folder.id">
                    <input type="radio" name="feed" :value="'folder:'+folder.id" v-model="feedSelected" v-if="folder.id">
                    <div class="selectgroup-label d-flex align-items-center w-100" v-if="folder.id">
                        <v-icon class="mr-2"
                                :class="{expanded: folder.is_expanded}"
                                @click.prevent="toggleFolderExpanded(folder)"
                                name="chevron-right" />
                        <span class="flex-fill text-left text-truncate">{{ folder.title }}</span>
                        <span class="counter text-right">{{ filteredFolderStats[folder.id] || '' }}</span>
                    </div>
                </label>
                <div v-show="!folder.id || folder.is_expanded" class="mt-1" :class="{'pl-3': folder.id}">
                    <label class="selectgroup"
                            :class="{'d-none': mustHideFeed(feed)}"
                            v-for="feed in folder.feeds">
                        <input type="radio" name="feed" :value="'feed:'+feed.id" v-model="feedSelected">
                        <div class="selectgroup-label d-flex align-items-center w-100">
                            <v-icon class="mr-2" name="rss" v-if="!feed.has_icon" />
                            <span class="icon mr-2" v-else><img :src="'./api/feeds/'+feed.id+'/icon'" alt="" loading="lazy"></span>
                            <span class="flex-fill text-left text-truncate">{{ feed.title }}</span>
                            <span class="counter text-right">{{ filteredFeedStats[feed.id] || '' }}</span>
                            <v-icon class="flex-shrink-0 mx-2"
                                    :title="feed_errors[feed.id]"
                                    v-if="!filterSelected && feed_errors[feed.id]"
                                    name="alert-circle" />
                        </div>
                    </label>
                </div>
            </div>
        </div>
        <div class="p-2 toolbar d-flex align-items-center border-top flex-shrink-0" v-if="loading.feeds">
            <span class="icon loading mx-2"></span>
            <span class="text-truncate cursor-default noselect">{{ $t('refreshing_progress', {count: loading.feeds}) }}</span>
        </div>
    </div>
    <!-- item list -->
    <div id="col-item-list" class="vh-100 position-relative d-flex flex-column border-right flex-shrink-0" :style="{width: itemListWidth+'px'}">
        <drag :width="itemListWidth" @resize="resizeItemList"></drag>
        <div class="px-2 toolbar d-flex align-items-center">
            <button class="toolbar-item mr-2 d-block d-md-none"
                    @click="feedSelected = null"
                    :title="$t('show_feeds')">
                <v-icon name="chevron-left" />
            </button>
            <div class="input-icon flex-grow-1">
                <v-icon name="search" />
                <!-- id used by keybindings -->
                <input id="searchbar" type="" class="d-block toolbar-search" v-model="itemSearch" @keydown.enter="$event.target.blur()">
            </div>
            <button class="toolbar-item ml-2"
                    @click="markItemsRead()"
                    v-if="filterSelected == 'unread'"
                    :title="$t('mark_all_read')">
                <v-icon name="check" />
            </button>


            <button class="btn btn-link toolbar-item px-2 ml-2" v-if="!current.type" disabled>
                <v-icon name="more-horizontal" />
            </button>
            <dropdown class="settings-dropdown"
                        toggle-class="btn btn-link toolbar-item px-2 ml-2"
                        drop="right"
                        :title="$t('feed_settings')"
                        v-if="current.type == 'feed'">
                <template v-slot:button>
                    <v-icon name="more-horizontal" />
                </template>
                <header class="dropdown-header" role="heading" aria-level="2">{{ current.feed.title }}</header>
                <a class="dropdown-item" :href="current.feed.link" rel="noopener noreferrer" target="_blank" referrerpolicy="no-referrer" v-if="current.feed.link">
                    <v-icon class="mr-1" name="globe" />
                    {{ $t('website') }}
                </a>
                <a class="dropdown-item" :href="current.feed.feed_link" rel="noopener noreferrer" target="_blank" referrerpolicy="no-referrer" v-if="current.feed.feed_link">
                    <v-icon class="mr-1" name="rss" />
                    {{ $t('feed_link') }}
                </a>
                <div class="dropdown-divider" v-if="current.feed.link || current.feed.feed_link"></div>
                <button class="dropdown-item" @click="renameFeed(current.feed)">
                    <v-icon class="mr-1" name="edit" />
                    {{ $t('rename') }}
                </button>
                <button class="dropdown-item" @click="updateFeedLink(current.feed)" v-if="current.feed.feed_link">
                    <v-icon class="mr-1" name="edit" />
                    {{ $t('change_link') }}
                </button>
                <div class="dropdown-divider"></div>
                <header class="dropdown-header" role="heading" aria-level="2">{{ $t('move_to') }}</header>
                <button class="dropdown-item"
                    v-if="folder.id != current.feed.folder_id"
                    v-for="folder in folders"
                    @click="moveFeed(current.feed, folder)">
                    <v-icon class="mr-1" name="folder" />
                    {{ folder.title }}
                </button>
                <button class="dropdown-item text-muted" @click="moveFeed(current.feed, null)" v-if="current.feed.folder_id">
                    <v-icon class="mr-1" name="folder-minus" />
                    \u2500\u2500
                </button>
                <button class="dropdown-item text-muted" @click="moveFeedToNewFolder(current.feed)">
                    <v-icon class="mr-1" name="folder-plus" />
                    {{ $t('new_folder') }}
                </button>
                <div class="dropdown-divider"></div>
                <button class="dropdown-item text-danger" @click.prevent="deleteFeed(current.feed)">
                    <v-icon class="mr-1" name="trash" />
                    {{ $t('delete') }}
                </button>
            </dropdown>
            <dropdown class="settings-dropdown"
                        toggle-class="btn btn-link toolbar-item px-2 ml-2"
                        :title="$t('folder_settings')"
                        drop="right"
                        v-if="current.type == 'folder'">
                <template v-slot:button>
                    <v-icon name="more-horizontal" />
                </template>
                <header class="dropdown-header" role="heading" aria-level="2">{{ current.folder.title }}</header>
                <button class="dropdown-item" @click="renameFolder(current.folder)">
                    <v-icon class="mr-1" name="edit" />
                    {{ $t('rename') }}
                </button>
                <div class="dropdown-divider"></div>
                <button class="dropdown-item text-danger" @click="deleteFolder(current.folder)">
                    <v-icon class="mr-1" name="trash" />
                    {{ $t('delete') }}
                </button>
            </dropdown>
        </div>
        <div id="item-list-scroll" class="p-2 overflow-auto scroll-touch border-top flex-grow-1" v-scroll="loadMoreItems" ref="itemlist">
            <label v-for="item in items" :key="item.id"
                    class="selectgroup">
                <input type="radio" name="item" :value="item.id" v-model="itemSelected">
                <div class="selectgroup-label d-flex flex-column">
                    <div style="line-height: 100%; opacity: .7; margin-bottom: .1rem;" class="d-flex align-items-center">
                        <transition name="indicator">
                            <v-icon class="icon-small mr-1" name="circle-full" v-if="item.status=='unread'" />
                            <v-icon class="icon-small mr-1" name="star-full" v-if="item.status=='starred'" />
                        </transition>
                        <small class="flex-fill text-truncate mr-1">
                            {{ (feedsById[item.feed_id] || {}).title }}
                        </small>
                        <small class="flex-shrink-0"><relative-time v-bind:title="formatDate(item.date)" :val="item.date"/></small>
                    </div>
                    <div>{{ item.title || $t('untitled') }}</div>
                </div>
            </label>
            <button class="btn btn-link btn-block loading my-3" v-if="itemsHasMore"></button>
        </div>
        <div class="px-3 py-2 border-top text-danger text-break" v-if="feed_errors[current.feed.id]">
            {{ feed_errors[current.feed.id] }}
        </div>
    </div>
    <!-- item show -->
    <div id="col-item" class="vh-100 d-flex flex-column w-100" style="min-width: 0;">
        <div class="toolbar px-2 d-flex align-items-center" v-if="itemSelectedDetails">
            <button class="toolbar-item"
                    @click="toggleItemStarred(itemSelectedDetails)"
                    :title="$t('mark_starred')">
                <v-icon name="star-full" v-if="itemSelectedDetails.status=='starred'" />
                <v-icon name="star" v-else-if="itemSelectedDetails.status!='starred'" />
            </button>
            <button class="toolbar-item"
                    :title="$t('mark_unread')"
                    @click="toggleItemRead(itemSelectedDetails)">
                <v-icon name="circle-full" v-if="itemSelectedDetails.status=='unread'" />
                <v-icon name="circle" v-if="itemSelectedDetails.status!='unread'" />
            </button>
            <dropdown class="settings-dropdown" toggle-class="toolbar-item px-2" drop="center" :title="$t('appearance')">
                <template v-slot:button>
                    <v-icon name="sliders" />
                </template>

                <button class="dropdown-item" :class="{active: !theme.font}" @click.stop="theme.font = ''">{{ $t('sans_serif') }}</button>
                <button class="dropdown-item font-serif" :class="{active: theme.font == 'serif'}" @click.stop="theme.font = 'serif'">{{ $t('serif') }}</button>
                <button class="dropdown-item font-monospace" :class="{active: theme.font == 'monospace'}" @click.stop="theme.font = 'monospace'">{{ $t('monospace') }}</button>

                <div class="d-flex text-center">
                    <button class="dropdown-item" style="font-size: 0.8rem" @click.stop="incrFont(-1)">A</button>
                    <button class="dropdown-item" style="font-size: 1.2rem" @click.stop="incrFont(1)">A</button>
                </div>
            </dropdown>
            <button class="toolbar-item"
                    :class="{active: itemSelectedReadability}"
                    @click="toggleReadability()"
                    :title="$t('read_here')">
                <v-icon :class="{'icon-loading': loading.readability}" name="book-open" />
            </button>
            <a class="toolbar-item" :href="itemSelectedDetails.link" rel="noopener noreferrer" target="_blank" referrerpolicy="no-referrer" :title="$t('open_link')">
                <v-icon name="external-link" />
            </a>
            <div class="flex-grow-1"></div>
            <button class="toolbar-item" @click="navigateToItem(-1)" :title="$t('previous_article')" :disabled="!items.length || itemSelected == items[0].id">
                <v-icon name="chevron-left" />
            </button>
            <button class="toolbar-item" @click="navigateToItem(+1)" :title="$t('next_article')" :disabled="!items.length || itemSelected == items[items.length - 1].id">
                <v-icon name="chevron-right" />
            </button>
            <button class="toolbar-item" @click="itemSelected=null" :title="$t('close_article')">
                <v-icon name="x" />
            </button>
        </div>
        <div v-if="itemSelectedDetails"
                ref="content"
                class="content px-4 pt-3 pb-5 border-top overflow-auto scroll-touch"
                :class="{'font-serif': theme.font == 'serif', 'font-monospace': theme.font == 'monospace'}"
                :style="{'font-size': theme.size + 'rem'}">
            <div class="content-wrapper">
                <h1><b>{{ itemSelectedDetails.title || $t('untitled') }}</b></h1>
                <div class="text-muted">
                    <div>
                        <span class="cursor-pointer" @click="feedSelected = 'feed:'+(feedsById[itemSelectedDetails.feed_id] || {}).id">
                            {{ (feedsById[itemSelectedDetails.feed_id] || {}).title }}
                        </span>
                    </div>
                    <time>{{ formatDate(itemSelectedDetails.date) }}</time>
                </div>
                <hr>
                <div v-if="!itemSelectedReadability">
                    <div v-if="contentImages.length">
                        <figure v-for="media in contentImages">
                            <img :src="media.url" loading="lazy">
                            <figcaption v-if="media.description">{{ media.description }}</figcaption>
                        </figure>
                    </div>
                    <audio class="w-100" controls v-for="media in contentAudios" :src="media.url"></audio>
                    <video class="w-100" controls v-for="media in contentVideos" :src="media.url"></video>
                </div>
                <div v-html="itemSelectedContent"></div>
            </div>
        </div>
    </div>
    <modal :open="!!settings" @hide="settings = ''">
        <button class="btn btn-link outline-none float-right p-2 mr-n2 mt-n2" style="line-height: 1" @click="settings = ''">
            <v-icon name="x" />
        </button>
        <div v-if="settings=='create'">
            <p class="cursor-default"><b>{{ $t('new_feed') }}</b></p>
            <form action="" @submit.prevent="createFeed(event)" class="mt-4">
                <label for="feed-url">{{ $t('url') }}</label>
                <input id="feed-url" name="url" type="url" class="form-control" required autocomplete="off" :readonly="feedNewChoice.length > 0" placeholder="https://example.com/feed" v-focus>
                <label for="feed-folder" class="mt-3 d-block">
                    {{ $t('folder') }}
                    <a href="#" class="float-right text-decoration-none" @click.prevent="createNewFeedFolder()">{{ $t('new_folder') }}</a>
                </label>
                <select class="form-control" id="feed-folder" name="folder_id" ref="newFeedFolder">
                    <option value="">---</option>
                    <option :value="folder.id" v-for="folder in folders" :selected="folder.id === current.feed.folder_id || folder.id === current.folder.id">{{ folder.title }}</option>
                </select>
                <div class="mt-4" v-if="feedNewChoice.length">
                    <p class="mb-2">
                        {{ $t('multiple_feeds_found') }}
                        <a href="#" class="float-right text-decoration-none" @click.prevent="resetFeedChoice()">{{ $t('cancel') }}</a>
                    </p>
                    <label class="selectgroup" v-for="choice in feedNewChoice">
                        <input type="radio" name="feedToAdd" :value="choice.url" v-model="feedNewChoiceSelected">
                        <div class="selectgroup-label">
                            <div class="text-truncate">{{ choice.title }}</div>
                            <div class="text-truncate" :class="{light: choice.title}">{{ choice.url }}</div>
                        </div>
                    </label>
                </div>
                <button class="btn btn-block btn-default mt-3" :class="{loading: loading.newfeed}" type="submit">{{ $t('add') }}</button>
            </form>
        </div>
        <div v-else-if="settings=='shortcuts'">
            <p class="cursor-default"><b>{{ $t('keyboard_shortcuts') }}</b></p>

            <table class="table table-borderless table-sm table-compact m-0">
                <tr><td><kbd>1</kbd> <kbd>2</kbd> <kbd>3</kbd></td>
                                                        <td>{{ $t('kb_show_filters') }}</td></tr>
                <tr><td><kbd>/</kbd></td>               <td>{{ $t('kb_focus_search') }}</td></tr>

                <tr><td colspan=2>&nbsp;</td></tr>
                <tr><td><kbd>j</kbd> <kbd>k</kbd></td>  <td>{{ $t('kb_next_prev_article') }}</td></tr>
                <tr><td><kbd>l</kbd> <kbd>h</kbd></td>  <td>{{ $t('kb_next_prev_feed') }}</td></tr>
                <tr><td><kbd>q</kbd></td>               <td>{{ $t('kb_close_article') }}</td></tr>

                <tr><td colspan=2>&nbsp;</td></tr>
                <tr><td><kbd>R</kbd></td>               <td>{{ $t('kb_mark_all_read') }}</td></tr>
                <tr><td><kbd>r</kbd></td>               <td>{{ $t('kb_mark_read') }}</td></tr>
                <tr><td><kbd>s</kbd></td>               <td>{{ $t('kb_mark_starred') }}</td></tr>
                <tr><td><kbd>o</kbd></td>               <td>{{ $t('kb_open_link') }}</td></tr>
                <tr><td><kbd>i</kbd></td>               <td>{{ $t('kb_read_here') }}</td> </tr>
                <tr><td><kbd>f</kbd> <kbd>b</kbd></td>  <td>{{ $t('kb_scroll_content') }}</td>
                </tr>
            </table>
        </div>
    </modal>
</div>
`;

  // src/assets/graphicarts/anchor.svg
  var anchor_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-anchor"><circle cx="12" cy="5" r="3"></circle><line x1="12" y1="22" x2="12" y2="8"></line><path d="M5 12H2a10 10 0 0 0 20 0h-3"></path></svg>';

  // src/assets/graphicarts/alert-circle.svg
  var alert_circle_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-alert-circle"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>';

  // src/assets/graphicarts/assorted.svg
  var assorted_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-bar-chart-2"><line x1="4" y1="6" x2="14" y2="6"></line><line x1="4" y1="12" x2="20" y2="12"></line><line x1="4" y1="18" x2="8" y2="18"></line></svg>\n';

  // src/assets/graphicarts/book-open.svg
  var book_open_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-book-open"><path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path></svg>';

  // src/assets/graphicarts/check.svg
  var check_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-check"><polyline points="20 6 9 17 4 12"></polyline></svg>';

  // src/assets/graphicarts/chevron-down.svg
  var chevron_down_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-chevron-down"><polyline points="6 9 12 15 18 9"></polyline></svg>';

  // src/assets/graphicarts/chevron-left.svg
  var chevron_left_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-chevron-left"><polyline points="15 18 9 12 15 6"></polyline></svg>';

  // src/assets/graphicarts/chevron-right.svg
  var chevron_right_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-chevron-right"><polyline points="9 18 15 12 9 6"></polyline></svg>';

  // src/assets/graphicarts/chevron-up.svg
  var chevron_up_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-chevron-up"><polyline points="18 15 12 9 6 15"></polyline></svg>';

  // src/assets/graphicarts/circle.svg
  var circle_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-circle"><circle cx="12" cy="12" r="10"></circle></svg>';

  // src/assets/graphicarts/circle-full.svg
  var circle_full_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-circle"><circle cx="12" cy="12" r="10"></circle></svg>\n';

  // src/assets/graphicarts/download.svg
  var download_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-download"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="7 10 12 15 17 10"></polyline><line x1="12" y1="15" x2="12" y2="3"></line></svg>';

  // src/assets/graphicarts/edit.svg
  var edit_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-edit"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>';

  // src/assets/graphicarts/external-link.svg
  var external_link_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-external-link"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>';

  // src/assets/graphicarts/folder.svg
  var folder_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-folder"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>';

  // src/assets/graphicarts/folder-minus.svg
  var folder_minus_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-folder-minus"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path><line x1="9" y1="14" x2="15" y2="14"></line></svg>';

  // src/assets/graphicarts/folder-plus.svg
  var folder_plus_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-folder-plus"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path><line x1="12" y1="11" x2="12" y2="17"></line><line x1="9" y1="14" x2="15" y2="14"></line></svg>';

  // src/assets/graphicarts/globe.svg
  var globe_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-globe"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>';

  // src/assets/graphicarts/help-circle.svg
  var help_circle_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-help-circle"><circle cx="12" cy="12" r="10"></circle><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>';

  // src/assets/graphicarts/layers.svg
  var layers_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-layers"><polygon points="12 2 2 7 12 12 22 7 12 2"></polygon><polyline points="2 17 12 22 22 17"></polyline><polyline points="2 12 12 17 22 12"></polyline></svg>';

  // src/assets/graphicarts/log-out.svg
  var log_out_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-log-out"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>';

  // src/assets/graphicarts/more-horizontal.svg
  var more_horizontal_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-more-horizontal"><circle cx="12" cy="12" r="1"></circle><circle cx="19" cy="12" r="1"></circle><circle cx="5" cy="12" r="1"></circle></svg>';

  // src/assets/graphicarts/plus.svg
  var plus_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-plus"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>';

  // src/assets/graphicarts/rotate-cw.svg
  var rotate_cw_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-rotate-cw"><polyline points="23 4 23 10 17 10"></polyline><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path></svg>';

  // src/assets/graphicarts/rss.svg
  var rss_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-rss"><path d="M4 11a9 9 0 0 1 9 9"></path><path d="M4 4a16 16 0 0 1 16 16"></path><circle cx="5" cy="19" r="1"></circle></svg>';

  // src/assets/graphicarts/search.svg
  var search_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-search"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>';

  // src/assets/graphicarts/sliders.svg
  var sliders_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-sliders"><line x1="4" y1="21" x2="4" y2="14"></line><line x1="4" y1="10" x2="4" y2="3"></line><line x1="12" y1="21" x2="12" y2="12"></line><line x1="12" y1="8" x2="12" y2="3"></line><line x1="20" y1="21" x2="20" y2="16"></line><line x1="20" y1="12" x2="20" y2="3"></line><line x1="1" y1="14" x2="7" y2="14"></line><line x1="9" y1="8" x2="15" y2="8"></line><line x1="17" y1="16" x2="23" y2="16"></line></svg>';

  // src/assets/graphicarts/star.svg
  var star_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-star"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>';

  // src/assets/graphicarts/star-full.svg
  var star_full_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-star"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>\n';

  // src/assets/graphicarts/trash.svg
  var trash_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-trash"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>';

  // src/assets/graphicarts/upload.svg
  var upload_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-upload"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="17 8 12 3 7 8"></polyline><line x1="12" y1="3" x2="12" y2="15"></line></svg>';

  // src/assets/graphicarts/x.svg
  var x_default = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-x"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>';

  // src/assets/javascripts/icons.ts
  var icons_default = {
    anchor: anchor_default,
    "alert-circle": alert_circle_default,
    assorted: assorted_default,
    "book-open": book_open_default,
    check: check_default,
    "chevron-down": chevron_down_default,
    "chevron-left": chevron_left_default,
    "chevron-right": chevron_right_default,
    "chevron-up": chevron_up_default,
    circle: circle_default,
    "circle-full": circle_full_default,
    download: download_default,
    edit: edit_default,
    "external-link": external_link_default,
    folder: folder_default,
    "folder-minus": folder_minus_default,
    "folder-plus": folder_plus_default,
    globe: globe_default,
    "help-circle": help_circle_default,
    layers: layers_default,
    "log-out": log_out_default,
    "more-horizontal": more_horizontal_default,
    plus: plus_default,
    "rotate-cw": rotate_cw_default,
    rss: rss_default,
    search: search_default,
    sliders: sliders_default,
    star: star_default,
    "star-full": star_full_default,
    trash: trash_default,
    upload: upload_default,
    x: x_default
  };

  // src/assets/javascripts/key.ts
  function setupKeybindings(vm3) {
    var helperFunctions = {
      scrollContent: function(direction) {
        var padding = 40;
        var scroll = document.querySelector(".content");
        if (!scroll) return;
        var height = scroll.getBoundingClientRect().height;
        var newpos = scroll.scrollTop + (height - padding) * direction;
        if (typeof scroll.scrollTo == "function") {
          scroll.scrollTo({ top: newpos, left: 0, behavior: "smooth" });
        } else {
          scroll.scrollTop = newpos;
        }
      }
    };
    var shortcutFunctions = {
      openItemLink: function() {
        if (vm3.itemSelectedDetails && vm3.itemSelectedDetails.link) {
          window.open(vm3.itemSelectedDetails.link, "_blank", "noopener,noreferrer");
        }
      },
      toggleReadability: function() {
        vm3.toggleReadability();
      },
      toggleItemRead: function() {
        if (vm3.itemSelected != null) {
          vm3.toggleItemRead(vm3.itemSelectedDetails);
        }
      },
      markAllRead: function() {
        if (vm3.filterSelected == "unread") {
          vm3.markItemsRead();
        }
      },
      toggleItemStarred: function() {
        if (vm3.itemSelected != null) {
          vm3.toggleItemStarred(vm3.itemSelectedDetails);
        }
      },
      focusSearch: function() {
        document.getElementById("searchbar").focus();
      },
      nextItem() {
        vm3.navigateToItem(1);
      },
      previousItem() {
        vm3.navigateToItem(-1);
      },
      nextFeed() {
        vm3.navigateToFeed(1);
      },
      previousFeed() {
        vm3.navigateToFeed(-1);
      },
      scrollForward: function() {
        helperFunctions.scrollContent(1);
      },
      scrollBackward: function() {
        helperFunctions.scrollContent(-1);
      },
      closeItem: function() {
        vm3.itemSelected = null;
      },
      showAll() {
        vm3.filterSelected = "";
      },
      showUnread() {
        vm3.filterSelected = "unread";
      },
      showStarred() {
        vm3.filterSelected = "starred";
      }
    };
    var keybindings = {
      "o": shortcutFunctions.openItemLink,
      "i": shortcutFunctions.toggleReadability,
      "r": shortcutFunctions.toggleItemRead,
      "R": shortcutFunctions.markAllRead,
      "s": shortcutFunctions.toggleItemStarred,
      "/": shortcutFunctions.focusSearch,
      "j": shortcutFunctions.nextItem,
      "k": shortcutFunctions.previousItem,
      "l": shortcutFunctions.nextFeed,
      "h": shortcutFunctions.previousFeed,
      "f": shortcutFunctions.scrollForward,
      "b": shortcutFunctions.scrollBackward,
      "q": shortcutFunctions.closeItem,
      "1": shortcutFunctions.showUnread,
      "2": shortcutFunctions.showStarred,
      "3": shortcutFunctions.showAll
    };
    var codebindings = {
      "KeyO": shortcutFunctions.openItemLink,
      "KeyI": shortcutFunctions.toggleReadability,
      //"r": shortcutFunctions.toggleItemRead,
      //"KeyR": shortcutFunctions.markAllRead,
      "KeyS": shortcutFunctions.toggleItemStarred,
      "Slash": shortcutFunctions.focusSearch,
      "KeyJ": shortcutFunctions.nextItem,
      "KeyK": shortcutFunctions.previousItem,
      "KeyL": shortcutFunctions.nextFeed,
      "KeyH": shortcutFunctions.previousFeed,
      "KeyF": shortcutFunctions.scrollForward,
      "KeyB": shortcutFunctions.scrollBackward,
      "KeyQ": shortcutFunctions.closeItem,
      "Digit1": shortcutFunctions.showUnread,
      "Digit2": shortcutFunctions.showStarred,
      "Digit3": shortcutFunctions.showAll
    };
    function isTextBox(element) {
      var tagName2 = element.tagName.toLowerCase();
      var inputBlocklist = ["button", "checkbox", "color", "file", "hidden", "image", "radio", "range", "reset", "search", "submit"];
      return tagName2 === "textarea" || tagName2 === "input" && inputBlocklist.indexOf(element.getAttribute("type").toLowerCase()) == -1;
    }
    document.addEventListener("keydown", function(event) {
      if (isTextBox(event.target) || event.metaKey || event.ctrlKey || event.altKey) {
        return;
      }
      var keybindFunction = keybindings[event.key] || codebindings[event.code];
      if (keybindFunction) {
        event.preventDefault();
        keybindFunction();
      }
    });
  }

  // src/assets/javascripts/app.ts
  var app = window.app;
  var vm;
  var TITLE = document.title;
  function scrollto(target2, scroll) {
    var padding = 10;
    var targetRect = target2.getBoundingClientRect();
    var scrollRect = scroll.getBoundingClientRect();
    var relativeOffset = targetRect.y - scrollRect.y;
    var absoluteOffset = relativeOffset + scroll.scrollTop;
    if (padding <= relativeOffset && relativeOffset + targetRect.height <= scrollRect.height - padding) return;
    var newPos = scroll.scrollTop;
    if (relativeOffset < padding) {
      newPos = absoluteOffset - padding;
    } else {
      newPos = absoluteOffset - scrollRect.height + targetRect.height + padding;
    }
    scroll.scrollTop = Math.round(newPos);
  }
  var debounce = function(callback, wait) {
    var timeout;
    return function() {
      var ctx = this, args = arguments;
      clearTimeout(timeout);
      timeout = setTimeout(function() {
        callback.apply(ctx, args);
      }, wait);
    };
  };
  Vue.directive("scroll", {
    inserted: function(el, binding) {
      el.addEventListener("scroll", debounce(function(event) {
        binding.value(event, el);
      }, 200));
    }
  });
  Vue.directive("focus", {
    inserted: function(el) {
      el.focus();
    }
  });
  Vue.component("drag", {
    props: ["width"],
    template: '<div class="drag"></div>',
    mounted: function() {
      var self = this;
      var startX = void 0;
      var initW = void 0;
      var onMouseMove = function(e) {
        var offset = e.clientX - startX;
        var newWidth = initW + offset;
        self.$emit("resize", newWidth);
      };
      var onMouseUp = function(e) {
        document.removeEventListener("mousemove", onMouseMove);
        document.removeEventListener("mouseup", onMouseUp);
      };
      this.$el.addEventListener("mousedown", function(e) {
        startX = e.clientX;
        initW = self.width;
        document.addEventListener("mousemove", onMouseMove);
        document.addEventListener("mouseup", onMouseUp);
      });
    }
  });
  Vue.component("dropdown", {
    props: ["toggle-class", "drop", "title"],
    data: function() {
      return { open: false };
    },
    template: `
    <div class="dropdown" :class="$attrs.class">
      <button ref="btn" @click="toggle" :class="btnToggleClass" :title="$props.title"><slot name="button"></slot></button>
      <div ref="menu" class="dropdown-menu" :class="{show: open}"><slot v-if="open"></slot></div>
    </div>
  `,
    computed: {
      btnToggleClass: function() {
        var c = this.$props.toggleClass || "";
        c += " dropdown-toggle dropdown-toggle-no-caret";
        c += this.open ? " show" : "";
        return c.trim();
      }
    },
    methods: {
      toggle: function(e) {
        this.open ? this.hide() : this.show();
      },
      show: function(e) {
        this.open = true;
        this.$refs.menu.style.top = this.$refs.btn.offsetHeight + "px";
        var drop = this.$props.drop;
        if (drop === "right") {
          this.$refs.menu.style.left = "auto";
          this.$refs.menu.style.right = "0";
        } else if (drop === "center") {
          this.$nextTick(function() {
            var btnWidth = this.$refs.btn.getBoundingClientRect().width;
            var menuWidth = this.$refs.menu.getBoundingClientRect().width;
            this.$refs.menu.style.left = "-" + (menuWidth - btnWidth) / 2 + "px";
          }.bind(this));
        }
        document.addEventListener("click", this.clickHandler);
      },
      hide: function() {
        this.open = false;
        document.removeEventListener("click", this.clickHandler);
      },
      clickHandler: function(e) {
        var dropdown = e.target.closest(".dropdown");
        if (dropdown == null || dropdown != this.$el) return this.hide();
        if (e.target.closest(".dropdown-item") != null) return this.hide();
      }
    }
  });
  Vue.component("modal", {
    props: ["open"],
    template: `
    <div class="modal custom-modal" tabindex="-1" v-if="$props.open">
      <div class="modal-dialog">
        <div class="modal-content" ref="content">
          <div class="modal-body">
            <slot v-if="$props.open"></slot>
          </div>
        </div>
      </div>
    </div>
  `,
    data: function() {
      return { opening: false };
    },
    watch: {
      "open": function(newVal) {
        if (newVal) {
          this.opening = true;
          document.addEventListener("click", this.handleClick);
        } else {
          document.removeEventListener("click", this.handleClick);
        }
      }
    },
    methods: {
      handleClick: function(e) {
        if (this.opening) {
          this.opening = false;
          return;
        }
        if (e.target.closest(".modal-content") == null) this.$emit("hide");
      }
    }
  });
  function dateRepr(d) {
    var sec = ((/* @__PURE__ */ new Date()).getTime() - d.getTime()) / 1e3;
    var neg = sec < 0;
    var out = "";
    sec = Math.abs(sec);
    if (sec < 2700)
      out = Math.round(sec / 60) + "m";
    else if (sec < 86400)
      out = Math.round(sec / 3600) + "h";
    else if (sec < 604800)
      out = Math.round(sec / 86400) + "d";
    else
      out = d.toLocaleDateString(void 0, { year: "numeric", month: "long", day: "numeric" });
    if (neg) return "-" + out;
    return out;
  }
  Vue.component("relative-time", {
    props: ["val"],
    data: function() {
      var d = new Date(this.val);
      return {
        "date": d,
        "formatted": dateRepr(d),
        "interval": null
      };
    },
    template: '<time :datetime="val">{{ formatted }}</time>',
    mounted: function() {
      this.interval = setInterval(function() {
        this.formatted = dateRepr(this.date);
      }.bind(this), 6e5);
    },
    destroyed: function() {
      clearInterval(this.interval);
    }
  });
  Vue.component("v-icon", {
    props: ["name"],
    template: '<span class="icon" v-html="content"></span>',
    computed: {
      content: function() {
        return icons_default[this.name] || "";
      }
    }
  });
  var app_default = {
    template: templates_default,
    created: function() {
      vm = this;
      this.refreshStats().then(this.refreshFeeds.bind(this)).then(this.refreshItems.bind(this, false));
      api_default.feeds.list_errors().then(function(errors) {
        vm.feed_errors = errors;
      });
      this.updateMetaTheme(app.settings.theme_name);
      this.$setLang(app.settings.language);
    },
    mounted: function() {
      setupKeybindings(this);
    },
    data: function() {
      var s = app.settings;
      return {
        "filterSelected": s.filter,
        "folders": [],
        "feeds": [],
        "feedSelected": s.feed,
        "feedListWidth": s.feed_list_width || 300,
        "feedNewChoice": [],
        "feedNewChoiceSelected": "",
        "items": [],
        "itemsHasMore": true,
        "itemSelected": null,
        "itemSelectedDetails": null,
        "itemSelectedReadability": "",
        "itemSearch": "",
        "itemSortNewestFirst": s.sort_newest_first,
        "itemListWidth": s.item_list_width || 300,
        "filteredFeedStats": {},
        "filteredFolderStats": {},
        "filteredTotalStats": null,
        "settings": "",
        "loading": {
          "feeds": 0,
          "newfeed": false,
          "items": false,
          "readability": false
        },
        "fonts": ["", "serif", "monospace"],
        "feedStats": {},
        "theme": {
          "name": s.theme_name,
          "font": s.theme_font,
          "size": s.theme_size
        },
        "themeColors": {
          "night": "#0e0e0e",
          "sepia": "#f4f0e5",
          "light": "#fff"
        },
        "refreshRate": s.refresh_rate,
        "authenticated": app.authenticated,
        "feed_errors": {},
        "refreshRateOptions": [
          { title: "0", value: 0 },
          { title: "10m", value: 10 },
          { title: "30m", value: 30 },
          { title: "1h", value: 60 },
          { title: "2h", value: 120 },
          { title: "4h", value: 240 },
          { title: "12h", value: 720 },
          { title: "24h", value: 1440 }
        ],
        "language": s.language,
        "languages": [
          { code: "en", name: "English" },
          { code: "de", name: "Deutsch" },
          { code: "es", name: "Espa\xF1ol" },
          { code: "fr", name: "Fran\xE7ais" },
          { code: "ja", name: "\u65E5\u672C\u8A9E" },
          { code: "pt", name: "Portugu\xEAs" },
          { code: "ru", name: "\u0420\u0443\u0441\u0441\u043A\u0438\u0439" },
          { code: "zh", name: "\u7B80\u4F53\u4E2D\u6587" }
        ]
      };
    },
    computed: {
      foldersWithFeeds: function() {
        var feedsByFolders = this.feeds.reduce(function(folders2, feed) {
          if (!folders2[feed.folder_id])
            folders2[feed.folder_id] = [feed];
          else
            folders2[feed.folder_id].push(feed);
          return folders2;
        }, {});
        var folders = this.folders.slice().map(function(folder) {
          folder.feeds = feedsByFolders[folder.id];
          return folder;
        });
        folders.push({ id: null, feeds: feedsByFolders[null] });
        return folders;
      },
      feedsById: function() {
        return this.feeds.reduce(function(acc, f) {
          acc[f.id] = f;
          return acc;
        }, {});
      },
      foldersById: function() {
        return this.folders.reduce(function(acc, f) {
          acc[f.id] = f;
          return acc;
        }, {});
      },
      current: function() {
        var parts = (this.feedSelected || "").split(":", 2);
        var type = parts[0];
        var guid = parts[1];
        var folder = {}, feed = {};
        if (type == "feed")
          feed = this.feedsById[guid] || {};
        if (type == "folder")
          folder = this.foldersById[guid] || {};
        return { type, feed, folder };
      },
      itemSelectedContent: function() {
        if (!this.itemSelected) return "";
        if (this.itemSelectedReadability)
          return this.itemSelectedReadability;
        return this.itemSelectedDetails.content || "";
      },
      contentImages: function() {
        if (!this.itemSelectedDetails) return [];
        return (this.itemSelectedDetails.media_links || []).filter((l) => l.type === "image");
      },
      contentAudios: function() {
        if (!this.itemSelectedDetails) return [];
        return (this.itemSelectedDetails.media_links || []).filter((l) => l.type === "audio");
      },
      contentVideos: function() {
        if (!this.itemSelectedDetails) return [];
        return (this.itemSelectedDetails.media_links || []).filter((l) => l.type === "video");
      },
      refreshRateTitle: function() {
        const entry = this.refreshRateOptions.find((o) => o.value === this.refreshRate);
        return entry ? entry.title : "0";
      }
    },
    watch: {
      "theme": {
        deep: true,
        handler: function(theme) {
          this.updateMetaTheme(theme.name);
          document.body.classList.value = "theme-" + theme.name;
          api_default.settings.update({
            theme_name: theme.name,
            theme_font: theme.font,
            theme_size: theme.size
          });
        }
      },
      "feedStats": {
        deep: true,
        handler: debounce(function() {
          var title = TITLE;
          var unreadCount = Object.values(this.feedStats).reduce(function(acc, stat) {
            return acc + stat.unread;
          }, 0);
          if (unreadCount) {
            title += " (" + unreadCount + ")";
          }
          document.title = title;
          this.computeStats();
        }, 500)
      },
      "filterSelected": function(newVal, oldVal) {
        if (oldVal === void 0) return;
        this.itemSelected = null;
        this.items = [];
        this.itemsHasMore = true;
        api_default.settings.update({ filter: newVal }).then(this.refreshItems.bind(this, false));
        this.computeStats();
      },
      "feedSelected": function(newVal, oldVal) {
        if (oldVal === void 0) return;
        this.itemSelected = null;
        this.items = [];
        this.itemsHasMore = true;
        api_default.settings.update({ feed: newVal }).then(this.refreshItems.bind(this, false));
        if (this.$refs.itemlist) this.$refs.itemlist.scrollTop = 0;
      },
      "itemSelected": function(newVal, oldVal) {
        this.itemSelectedReadability = "";
        if (newVal === null) {
          this.itemSelectedDetails = null;
          return;
        }
        if (this.$refs.content) this.$refs.content.scrollTop = 0;
        api_default.items.get(newVal).then(function(item) {
          this.itemSelectedDetails = item;
          if (this.itemSelectedDetails.status == "unread") {
            api_default.items.update(this.itemSelectedDetails.id, { status: "read" }).then(function() {
              this.feedStats[this.itemSelectedDetails.feed_id].unread -= 1;
              var itemInList = this.items.find(function(i) {
                return i.id == item.id;
              });
              if (itemInList) itemInList.status = "read";
              this.itemSelectedDetails.status = "read";
            }.bind(this));
          }
        }.bind(this));
      },
      "itemSearch": debounce(function(newVal) {
        this.refreshItems();
      }, 500),
      "itemSortNewestFirst": function(newVal, oldVal) {
        if (oldVal === void 0) return;
        api_default.settings.update({ sort_newest_first: newVal }).then(vm.refreshItems.bind(this, false));
      },
      "feedListWidth": debounce(function(newVal, oldVal) {
        if (oldVal === void 0) return;
        api_default.settings.update({ feed_list_width: newVal });
      }, 1e3),
      "itemListWidth": debounce(function(newVal, oldVal) {
        if (oldVal === void 0) return;
        api_default.settings.update({ item_list_width: newVal });
      }, 1e3),
      "refreshRate": function(newVal, oldVal) {
        if (oldVal === void 0) return;
        api_default.settings.update({ refresh_rate: newVal });
      }
    },
    methods: {
      updateMetaTheme: function(theme) {
        document.querySelector("meta[name='theme-color']").content = this.themeColors[theme];
      },
      refreshStats: function(loopMode) {
        return api_default.status().then(function(data) {
          if (loopMode && !vm.itemSelected) vm.refreshItems();
          vm.loading.feeds = data.running;
          if (data.running) {
            setTimeout(vm.refreshStats.bind(vm, true), 500);
          }
          vm.feedStats = data.stats.reduce(function(acc, stat) {
            acc[stat.feed_id] = stat;
            return acc;
          }, {});
          api_default.feeds.list_errors().then(function(errors) {
            vm.feed_errors = errors;
          });
        });
      },
      getItemsQuery: function() {
        var query2 = {};
        if (this.feedSelected) {
          var parts = this.feedSelected.split(":", 2);
          var type = parts[0];
          var guid = parts[1];
          if (type == "feed") {
            query2.feed_id = guid;
          } else if (type == "folder") {
            query2.folder_id = guid;
          }
        }
        if (this.filterSelected) {
          query2.status = this.filterSelected;
        }
        if (this.itemSearch) {
          query2.search = this.itemSearch;
        }
        if (!this.itemSortNewestFirst) {
          query2.oldest_first = true;
        }
        return query2;
      },
      refreshFeeds: function() {
        return Promise.all([api_default.folders.list(), api_default.feeds.list()]).then(function(values2) {
          vm.folders = values2[0];
          vm.feeds = values2[1];
        });
      },
      refreshItems: function(loadMore = false) {
        if (this.feedSelected === null) {
          vm.items = [];
          vm.itemsHasMore = false;
          return;
        }
        var query2 = this.getItemsQuery();
        if (loadMore) {
          query2.after = vm.items[vm.items.length - 1].id;
        }
        this.loading.items = true;
        return api_default.items.list(query2).then(function(data) {
          if (loadMore) {
            vm.items = vm.items.concat(data.list);
          } else {
            vm.items = data.list;
          }
          vm.itemsHasMore = data.has_more;
          vm.loading.items = false;
          vm.$nextTick(function() {
            if (vm.itemsHasMore && !vm.loading.items && vm.itemListCloseToBottom()) {
              vm.refreshItems(true);
            }
          });
        });
      },
      itemListCloseToBottom: function() {
        var bottomSpace = 70;
        var scale = (parseFloat(getComputedStyle(document.documentElement).fontSize) || 16) / 16;
        var el = this.$refs.itemlist;
        if (el.scrollHeight === 0) return false;
        var closeToBottom = el.scrollHeight - el.scrollTop - el.offsetHeight < bottomSpace * scale;
        return closeToBottom;
      },
      loadMoreItems: function(event, el) {
        if (!this.itemsHasMore) return;
        if (this.loading.items) return;
        if (this.itemListCloseToBottom()) return this.refreshItems(true);
        if (this.itemSelected && this.itemSelected === this.items[this.items.length - 1].id) return this.refreshItems(true);
      },
      markItemsRead: function() {
        var query2 = this.getItemsQuery();
        api_default.items.mark_read(query2).then(function() {
          vm.items = [];
          vm.itemsPage = { "cur": 1, "num": 1 };
          vm.itemSelected = null;
          vm.itemsHasMore = false;
          vm.refreshStats();
        });
      },
      toggleFolderExpanded: function(folder) {
        folder.is_expanded = !folder.is_expanded;
        api_default.folders.update(folder.id, { is_expanded: folder.is_expanded });
      },
      formatDate: function(datestr) {
        var options = {
          year: "numeric",
          month: "long",
          day: "numeric",
          hour: "2-digit",
          minute: "2-digit"
        };
        return new Date(datestr).toLocaleDateString(void 0, options);
      },
      moveFeed: function(feed, folder) {
        var folder_id = folder ? folder.id : null;
        api_default.feeds.update(feed.id, { folder_id }).then(function() {
          feed.folder_id = folder_id;
          vm.refreshStats();
        });
      },
      moveFeedToNewFolder: function(feed) {
        var title = prompt(this.$t("prompt_folder_name"));
        if (!title) return;
        api_default.folders.create({ "title": title }).then(function(folder) {
          api_default.feeds.update(feed.id, { folder_id: folder.id }).then(function() {
            vm.refreshFeeds().then(function() {
              vm.refreshStats();
            });
          });
        });
      },
      createNewFeedFolder: function() {
        var title = prompt(this.$t("prompt_folder_name"));
        if (!title) return;
        api_default.folders.create({ "title": title }).then(function(result) {
          vm.refreshFeeds().then(function() {
            vm.$nextTick(function() {
              if (vm.$refs.newFeedFolder) {
                vm.$refs.newFeedFolder.value = result.id;
              }
            });
          });
        });
      },
      renameFolder: function(folder) {
        var newTitle = prompt(this.$t("prompt_new_title"), folder.title);
        if (newTitle) {
          api_default.folders.update(folder.id, { title: newTitle }).then(function() {
            folder.title = newTitle;
            this.folders.sort(function(a, b) {
              return a.title.localeCompare(b.title);
            });
          }.bind(this));
        }
      },
      deleteFolder: function(folder) {
        if (confirm(this.$t("confirm_delete", { name: folder.title }))) {
          api_default.folders.delete(folder.id).then(function() {
            vm.feedSelected = null;
            vm.refreshStats();
            vm.refreshFeeds();
          });
        }
      },
      updateFeedLink: function(feed) {
        var newLink = prompt(this.$t("prompt_feed_link"), feed.feed_link);
        if (newLink) {
          api_default.feeds.update(feed.id, { feed_link: newLink }).then(function() {
            feed.feed_link = newLink;
          });
        }
      },
      renameFeed: function(feed) {
        var newTitle = prompt(this.$t("prompt_new_title"), feed.title);
        if (newTitle) {
          api_default.feeds.update(feed.id, { title: newTitle }).then(function() {
            feed.title = newTitle;
          });
        }
      },
      deleteFeed: function(feed) {
        if (confirm(this.$t("confirm_delete", { name: feed.title }))) {
          api_default.feeds.delete(feed.id).then(function() {
            vm.feedSelected = null;
            vm.refreshStats();
            vm.refreshFeeds();
          });
        }
      },
      createFeed: function(event) {
        var form = event.target;
        var data = {
          url: form.querySelector("input[name=url]").value,
          folder_id: parseInt(form.querySelector("select[name=folder_id]").value) || null
        };
        if (this.feedNewChoiceSelected) {
          data.url = this.feedNewChoiceSelected;
        }
        this.loading.newfeed = true;
        api_default.feeds.create(data).then(function(result) {
          if (result.status === "success") {
            vm.refreshFeeds();
            vm.refreshStats();
            vm.settings = "";
            vm.feedSelected = "feed:" + result.feed.id;
          } else if (result.status === "multiple") {
            vm.feedNewChoice = result.choice;
            vm.feedNewChoiceSelected = result.choice[0].url;
          } else {
            alert("No feeds found at the given url.");
          }
          vm.loading.newfeed = false;
        });
      },
      toggleItemStatus: function(item, targetstatus, fallbackstatus) {
        var oldstatus = item.status;
        var newstatus = item.status !== targetstatus ? targetstatus : fallbackstatus;
        var updateStats = function(status, incr) {
          if (status == "unread" || status == "starred") {
            this.feedStats[item.feed_id][status] += incr;
          }
        }.bind(this);
        api_default.items.update(item.id, { status: newstatus }).then(function() {
          updateStats(oldstatus, -1);
          updateStats(newstatus, 1);
          var itemInList = this.items.find(function(i) {
            return i.id == item.id;
          });
          if (itemInList) itemInList.status = newstatus;
          item.status = newstatus;
        }.bind(this));
      },
      toggleItemStarred: function(item) {
        this.toggleItemStatus(item, "starred", "read");
      },
      toggleItemRead: function(item) {
        this.toggleItemStatus(item, "unread", "read");
      },
      importOPML: function(event) {
        var input = event.target;
        var form = document.querySelector("#opml-import-form");
        this.$refs.menuDropdown.hide();
        api_default.upload_opml(form).then(function() {
          input.value = "";
          vm.refreshFeeds();
          vm.refreshStats();
        });
      },
      logout: function() {
        api_default.logout().then(function() {
          document.location.reload();
        });
      },
      toggleReadability: function() {
        if (this.itemSelectedReadability) {
          this.itemSelectedReadability = null;
          return;
        }
        var item = this.itemSelectedDetails;
        if (!item) return;
        if (item.link) {
          this.loading.readability = true;
          api_default.crawl(item.link).then(function(data) {
            vm.itemSelectedReadability = data && data.content;
            vm.loading.readability = false;
          });
        }
      },
      showSettings: function(settings) {
        this.settings = settings;
        if (settings === "create") {
          vm.feedNewChoice = [];
          vm.feedNewChoiceSelected = "";
        }
      },
      resizeFeedList: function(width) {
        this.feedListWidth = Math.min(Math.max(200, width), 700);
      },
      resizeItemList: function(width) {
        this.itemListWidth = Math.min(Math.max(200, width), 700);
      },
      resetFeedChoice: function() {
        this.feedNewChoice = [];
        this.feedNewChoiceSelected = "";
      },
      incrFont: function(x) {
        this.theme.size = +(this.theme.size + 0.1 * x).toFixed(1);
      },
      fetchAllFeeds: function() {
        if (this.loading.feeds) return;
        api_default.feeds.refresh().then(function() {
          vm.refreshStats();
        });
      },
      computeStats: function() {
        var filter = this.filterSelected;
        if (!filter) {
          this.filteredFeedStats = {};
          this.filteredFolderStats = {};
          this.filteredTotalStats = null;
          return;
        }
        var statsFeeds = {}, statsFolders = {}, statsTotal = 0;
        for (var i = 0; i < this.feeds.length; i++) {
          var feed = this.feeds[i];
          if (!this.feedStats[feed.id]) continue;
          var n = vm.feedStats[feed.id][filter] || 0;
          if (!statsFolders[feed.folder_id]) statsFolders[feed.folder_id] = 0;
          statsFeeds[feed.id] = n;
          statsFolders[feed.folder_id] += n;
          statsTotal += n;
        }
        this.filteredFeedStats = statsFeeds;
        this.filteredFolderStats = statsFolders;
        this.filteredTotalStats = statsTotal;
      },
      // navigation helper, navigate relative to selected item
      navigateToItem: function(relativePosition) {
        let vm3 = this;
        if (vm3.itemSelected == null) {
          if (vm3.items.length !== 0) vm3.itemSelected = vm3.items[0].id;
          return;
        }
        var itemPosition = vm3.items.findIndex(function(x) {
          return x.id === vm3.itemSelected;
        });
        if (itemPosition === -1) {
          if (vm3.items.length !== 0) vm3.itemSelected = vm3.items[0].id;
          return;
        }
        var newPosition = itemPosition + relativePosition;
        if (newPosition < 0 || newPosition >= vm3.items.length) return;
        vm3.itemSelected = vm3.items[newPosition].id;
        vm3.$nextTick(function() {
          var scroll = document.querySelector("#item-list-scroll");
          var handle = scroll.querySelector("input[type=radio]:checked");
          var target2 = handle && handle.parentElement;
          if (target2 && scroll) scrollto(target2, scroll);
          vm3.loadMoreItems();
        });
      },
      // navigation helper, navigate relative to selected feed
      navigateToFeed: function(relativePosition) {
        let vm3 = this;
        const navigationList = this.foldersWithFeeds.filter((folder) => !folder.id || !vm3.mustHideFolder(folder)).map((folder) => {
          if (this.mustHideFolder(folder)) return [];
          const folds = folder.id ? [`folder:${folder.id}`] : [];
          const feeds = folder.is_expanded || !folder.id ? (folder.feeds || []).filter((f) => !vm3.mustHideFeed(f)).map((f) => `feed:${f.id}`) : [];
          return folds.concat(feeds);
        }).flat();
        navigationList.unshift("");
        var currentFeedPosition = navigationList.indexOf(vm3.feedSelected);
        if (currentFeedPosition == -1) {
          vm3.feedSelected = "";
          return;
        }
        var newPosition = currentFeedPosition + relativePosition;
        if (newPosition < 0 || newPosition >= navigationList.length) return;
        vm3.feedSelected = navigationList[newPosition];
        vm3.$nextTick(function() {
          var scroll = document.querySelector("#feed-list-scroll");
          var handle = scroll.querySelector("input[type=radio]:checked");
          var target2 = handle && handle.parentElement;
          if (target2 && scroll) scrollto(target2, scroll);
        });
      },
      changeRefreshRate: function(offset) {
        const curIdx = this.refreshRateOptions.findIndex((o) => o.value === this.refreshRate);
        if (curIdx <= 0 && offset < 0) return;
        if (curIdx >= this.refreshRateOptions.length - 1 && offset > 0) return;
        this.refreshRate = this.refreshRateOptions[curIdx + offset].value;
      },
      mustHideFolder: function(folder) {
        return this.filterSelected && !(this.current.folder.id == folder.id || this.current.feed.folder_id == folder.id) && !this.filteredFolderStats[folder.id] && (!this.itemSelectedDetails || (this.feedsById[this.itemSelectedDetails.feed_id] || {}).folder_id != folder.id);
      },
      mustHideFeed: function(feed) {
        return this.filterSelected && !(this.current.feed.id == feed.id) && !this.filteredFeedStats[feed.id] && (!this.itemSelectedDetails || this.itemSelectedDetails.feed_id != feed.id);
      },
      changeLanguage(lang) {
        this.$setLang(lang);
        this.language = lang;
        api_default.settings.update({ language: lang });
      }
    }
  };

  // src/assets/javascripts/templates/login.html
  var login_default = `<div class="login-page">
    <form @submit.prevent="login">
        <img src="./static/graphicarts/anchor.svg" alt="">
        <div class="text-danger text-center my-3" v-if="hasError">{{ $t('login_error') }}</div>
        <div class="form-group">
            <label for="username">{{ $t('username') }}</label>
            <input name="username" class="form-control" id="username" autocomplete="off" required autofocus>
        </div>
        <div class="form-group">
            <label for="password">{{ $t('password') }}</label>
            <input name="password" class="form-control" id="password" type="password" required>
        </div>
        <button class="btn btn-block btn-default" type="submit">{{ $t('login') }}</button>
    </form>
</div>
`;

  // src/assets/javascripts/login.ts
  var login_default2 = {
    template: login_default,
    data: function() {
      return { hasError: false };
    },
    created: function() {
      this.$setLang(window.app.settings.language);
    },
    methods: {
      login: function(event) {
        event.preventDefault();
        var data = new FormData(event.target);
        fetch("./login", { method: "POST", body: data }).then(function(res) {
          if (res.ok) {
            document.location.assign("./");
          } else {
            this.hasError = true;
          }
        }.bind(this));
      }
    }
  };

  // src/assets/javascripts/main.ts
  Vue.use(i18n_default);
  var vm2 = new Vue({
    render: function(h) {
      return h(window.app.authenticated ? app_default : login_default2);
    }
  }).$mount("#app");
})();
/*! Bundled license information:

vue/dist/vue.esm.js:
  (*!
   * Vue.js v2.7.16
   * (c) 2014-2023 Evan You
   * Released under the MIT License.
   *)
*/
