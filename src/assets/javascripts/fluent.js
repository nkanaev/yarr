/** @fluent/bundle@0.19.1 */
(function (global, factory) {
    typeof exports === 'object' && typeof module !== 'undefined' ? factory(exports) :
    typeof define === 'function' && define.amd ? define('@fluent/bundle', ['exports'], factory) :
    (global = typeof globalThis !== 'undefined' ? globalThis : global || self, factory(global.FluentBundle = {}));
})(this, (function (exports) { 'use strict';

    /**
     * The `FluentType` class is the base of Fluent's type system.
     *
     * Fluent types wrap JavaScript values and store additional configuration for
     * them, which can then be used in the `toString` method together with a proper
     * `Intl` formatter.
     */
    class FluentType {
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
    }
    /**
     * A {@link FluentType} representing no correct value.
     */
    class FluentNone extends FluentType {
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
    }
    /**
     * A {@link FluentType} representing a number.
     *
     * A `FluentNumber` instance stores the number value of the number it
     * represents. It may also store an option bag of options which will be passed
     * to `Intl.NumerFormat` when the `FluentNumber` is formatted to a string.
     */
    class FluentNumber extends FluentType {
        /**
         * Create an instance of `FluentNumber` with options to the
         * `Intl.NumberFormat` constructor.
         *
         * @param value The number value of this `FluentNumber`.
         * @param opts Options which will be passed to `Intl.NumberFormat`.
         */
        constructor(value, opts = {}) {
            super(value);
            this.opts = opts;
        }
        /**
         * Format this `FluentNumber` to a string.
         */
        toString(scope) {
            if (scope) {
                try {
                    const nf = scope.memoizeIntlObject(Intl.NumberFormat, this.opts);
                    return nf.format(this.value);
                }
                catch (err) {
                    scope.reportError(err);
                }
            }
            return this.value.toString(10);
        }
    }
    /**
     * A {@link FluentType} representing a date and time.
     *
     * A `FluentDateTime` instance stores a Date object, Temporal object, or a number
     * as a numerical timestamp in milliseconds. It may also store an
     * option bag of options which will be passed to `Intl.DateTimeFormat` when the
     * `FluentDateTime` is formatted to a string.
     */
    class FluentDateTime extends FluentType {
        static supportsValue(value) {
            if (typeof value === "number")
                return true;
            if (value instanceof Date)
                return true;
            if (value instanceof FluentType)
                return FluentDateTime.supportsValue(value.valueOf());
            // Temporary workaround to support environments without Temporal
            if ("Temporal" in globalThis) {
                // for TypeScript, which doesn't know about Temporal yet
                const _Temporal = globalThis.Temporal;
                if (value instanceof _Temporal.Instant ||
                    value instanceof _Temporal.PlainDateTime ||
                    value instanceof _Temporal.PlainDate ||
                    value instanceof _Temporal.PlainMonthDay ||
                    value instanceof _Temporal.PlainTime ||
                    value instanceof _Temporal.PlainYearMonth) {
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
        constructor(value, opts = {}) {
            // unwrap any FluentType value, but only retain the opts from FluentDateTime
            if (value instanceof FluentDateTime) {
                opts = { ...value.opts, ...opts };
                value = value.value;
            }
            else if (value instanceof FluentType) {
                value = value.valueOf();
            }
            // Intl.DateTimeFormat defaults to gregorian calendar, but Temporal defaults to iso8601
            if (typeof value === "object" &&
                "calendarId" in value &&
                opts.calendar === undefined) {
                opts = { ...opts, calendar: value.calendarId };
            }
            super(value);
            this.opts = opts;
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
                }
                catch (err) {
                    scope.reportError(err);
                }
            }
            if (typeof this.value === "number" || this.value instanceof Date) {
                return new Date(this.value).toISOString();
            }
            return this.value.toString();
        }
    }

    /**
     * The role of the Fluent resolver is to format a `Pattern` to an instance of
     * `FluentValue`. For performance reasons, primitive strings are considered
     * such instances, too.
     *
     * Translations can contain references to other messages or variables,
     * conditional logic in form of select expressions, traits which describe their
     * grammatical features, and can use Fluent builtins which make use of the
     * `Intl` formatters to format numbers and dates into the bundle's languages.
     * See the documentation of the Fluent syntax for more information.
     *
     * In case of errors the resolver will try to salvage as much of the
     * translation as possible. In rare situations where the resolver didn't know
     * how to recover from an error it will return an instance of `FluentNone`.
     *
     * All expressions resolve to an instance of `FluentValue`. The caller should
     * use the `toString` method to convert the instance to a native value.
     *
     * Functions in this file pass around an instance of the `Scope` class, which
     * stores the data required for successful resolution and error recovery.
     */
    /**
     * The maximum number of placeables which can be expanded in a single call to
     * `formatPattern`. The limit protects against the Billion Laughs and Quadratic
     * Blowup attacks. See https://msdn.microsoft.com/en-us/magazine/ee335713.aspx.
     */
    const MAX_PLACEABLES = 100;
    /** Unicode bidi isolation characters. */
    const FSI = "\u2068";
    const PDI = "\u2069";
    /** Helper: match a variant key to the given selector. */
    function match(scope, selector, key) {
        if (key === selector) {
            // Both are strings.
            return true;
        }
        // XXX Consider comparing options too, e.g. minimumFractionDigits.
        if (key instanceof FluentNumber &&
            selector instanceof FluentNumber &&
            key.value === selector.value) {
            return true;
        }
        if (selector instanceof FluentNumber && typeof key === "string") {
            let category = scope
                .memoizeIntlObject(Intl.PluralRules, selector.opts)
                .select(selector.value);
            if (key === category) {
                return true;
            }
        }
        return false;
    }
    /** Helper: resolve the default variant from a list of variants. */
    function getDefault(scope, variants, star) {
        if (variants[star]) {
            return resolvePattern(scope, variants[star].value);
        }
        scope.reportError(new RangeError("No default"));
        return new FluentNone();
    }
    /** Helper: resolve arguments to a call expression. */
    function getArguments(scope, args) {
        const positional = [];
        const named = Object.create(null);
        for (const arg of args) {
            if (arg.type === "narg") {
                named[arg.name] = resolveExpression(scope, arg.value);
            }
            else {
                positional.push(resolveExpression(scope, arg));
            }
        }
        return { positional, named };
    }
    /** Resolve an expression to a Fluent type. */
    function resolveExpression(scope, expr) {
        switch (expr.type) {
            case "str":
                return expr.value;
            case "num":
                return new FluentNumber(expr.value, {
                    minimumFractionDigits: expr.precision,
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
    /** Resolve a reference to a variable. */
    function resolveVariableReference(scope, { name }) {
        let arg;
        if (scope.params) {
            // We're inside a TermReference. It's OK to reference undefined parameters.
            if (Object.prototype.hasOwnProperty.call(scope.params, name)) {
                arg = scope.params[name];
            }
            else {
                return new FluentNone(`$${name}`);
            }
        }
        else if (scope.args &&
            Object.prototype.hasOwnProperty.call(scope.args, name)) {
            // We're in the top-level Pattern or inside a MessageReference. Missing
            // variables references produce ReferenceErrors.
            arg = scope.args[name];
        }
        else {
            scope.reportError(new ReferenceError(`Unknown variable: $${name}`));
            return new FluentNone(`$${name}`);
        }
        // Return early if the argument already is an instance of FluentType.
        if (arg instanceof FluentType) {
            return arg;
        }
        // Convert the argument to a Fluent type.
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
    /** Resolve a reference to another message. */
    function resolveMessageReference(scope, { name, attr }) {
        const message = scope.bundle._messages.get(name);
        if (!message) {
            scope.reportError(new ReferenceError(`Unknown message: ${name}`));
            return new FluentNone(name);
        }
        if (attr) {
            const attribute = message.attributes[attr];
            if (attribute) {
                return resolvePattern(scope, attribute);
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
    /** Resolve a call to a Term with key-value arguments. */
    function resolveTermReference(scope, { name, attr, args }) {
        const id = `-${name}`;
        const term = scope.bundle._terms.get(id);
        if (!term) {
            scope.reportError(new ReferenceError(`Unknown term: ${id}`));
            return new FluentNone(id);
        }
        if (attr) {
            const attribute = term.attributes[attr];
            if (attribute) {
                // Every TermReference has its own variables.
                scope.params = getArguments(scope, args).named;
                const resolved = resolvePattern(scope, attribute);
                scope.params = null;
                return resolved;
            }
            scope.reportError(new ReferenceError(`Unknown attribute: ${attr}`));
            return new FluentNone(`${id}.${attr}`);
        }
        scope.params = getArguments(scope, args).named;
        const resolved = resolvePattern(scope, term.value);
        scope.params = null;
        return resolved;
    }
    /** Resolve a call to a Function with positional and key-value arguments. */
    function resolveFunctionReference(scope, { name, args }) {
        // Some functions are built-in. Others may be provided by the runtime via
        // the `FluentBundle` constructor.
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
        }
        catch (err) {
            scope.reportError(err);
            return new FluentNone(`${name}()`);
        }
    }
    /** Resolve a select expression to the member object. */
    function resolveSelectExpression(scope, { selector, variants, star }) {
        let sel = resolveExpression(scope, selector);
        if (sel instanceof FluentNone) {
            return getDefault(scope, variants, star);
        }
        // Match the selector against keys of each variant, in order.
        for (const variant of variants) {
            const key = resolveExpression(scope, variant.key);
            if (match(scope, sel, key)) {
                return resolvePattern(scope, variant.value);
            }
        }
        return getDefault(scope, variants, star);
    }
    /** Resolve a pattern (a complex string with placeables). */
    function resolveComplexPattern(scope, ptn) {
        if (scope.dirty.has(ptn)) {
            scope.reportError(new RangeError("Cyclic reference"));
            return new FluentNone();
        }
        // Tag the pattern as dirty for the purpose of the current resolution.
        scope.dirty.add(ptn);
        const result = [];
        // Wrap interpolations with Directional Isolate Formatting characters
        // only when the pattern has more than one element.
        const useIsolating = scope.bundle._useIsolating && ptn.length > 1;
        for (const elem of ptn) {
            if (typeof elem === "string") {
                result.push(scope.bundle._transform(elem));
                continue;
            }
            scope.placeables++;
            if (scope.placeables > MAX_PLACEABLES) {
                scope.dirty.delete(ptn);
                // This is a fatal error which causes the resolver to instantly bail out
                // on this pattern. The length check protects against excessive memory
                // usage, and throwing protects against eating up the CPU when long
                // placeables are deeply nested.
                throw new RangeError(`Too many placeables expanded: ${scope.placeables}, ` +
                    `max allowed is ${MAX_PLACEABLES}`);
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
    /**
     * Resolve a simple or a complex Pattern to a FluentString
     * (which is really the string primitive).
     */
    function resolvePattern(scope, value) {
        // Resolve a simple pattern.
        if (typeof value === "string") {
            return scope.bundle._transform(value);
        }
        return resolveComplexPattern(scope, value);
    }

    class Scope {
        constructor(bundle, errors, args) {
            /**
             * The Set of patterns already encountered during this resolution.
             * Used to detect and prevent cyclic resolutions.
             * @ignore
             */
            this.dirty = new WeakSet();
            /** A dict of parameters passed to a TermReference. */
            this.params = null;
            /**
             * The running count of placeables resolved so far.
             * Used to detect the Billion Laughs and Quadratic Blowup attacks.
             * @ignore
             */
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
        memoizeIntlObject(ctor, opts) {
            let cache = this.bundle._intls.get(ctor);
            if (!cache) {
                cache = {};
                this.bundle._intls.set(ctor, cache);
            }
            let id = JSON.stringify(opts);
            if (!cache[id]) {
                // @ts-expect-error This is fine.
                // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
                cache[id] = new ctor(this.bundle.locales, opts);
            }
            return cache[id];
        }
    }

    /**
     * @overview
     *
     * The FTL resolver ships with a number of functions built-in.
     *
     * Each function take two arguments:
     *   - args - an array of positional args
     *   - opts - an object of key-value args
     *
     * Arguments to functions are guaranteed to already be instances of
     * `FluentValue`.  Functions must return `FluentValues` as well.
     */
    function values(opts, allowed) {
        const unwrapped = Object.create(null);
        for (const [name, opt] of Object.entries(opts)) {
            if (allowed.includes(name)) {
                unwrapped[name] = opt.valueOf();
            }
        }
        return unwrapped;
    }
    const NUMBER_ALLOWED = [
        "unitDisplay",
        "currencyDisplay",
        "useGrouping",
        "minimumIntegerDigits",
        "minimumFractionDigits",
        "maximumFractionDigits",
        "minimumSignificantDigits",
        "maximumSignificantDigits",
    ];
    /**
     * The implementation of the `NUMBER()` builtin available to translations.
     *
     * Translations may call the `NUMBER()` builtin in order to specify formatting
     * options of a number. For example:
     *
     *     pi = The value of π is {NUMBER($pi, maximumFractionDigits: 2)}.
     *
     * The implementation expects an array of {@link FluentValue | FluentValues} representing the
     * positional arguments, and an object of named {@link FluentValue | FluentValues} representing the
     * named parameters.
     *
     * The following options are recognized:
     *
     *     unitDisplay
     *     currencyDisplay
     *     useGrouping
     *     minimumIntegerDigits
     *     minimumFractionDigits
     *     maximumFractionDigits
     *     minimumSignificantDigits
     *     maximumSignificantDigits
     *
     * Other options are ignored.
     *
     * @param args The positional arguments passed to this `NUMBER()`.
     * @param opts The named argments passed to this `NUMBER()`.
     */
    function NUMBER(args, opts) {
        let arg = args[0];
        if (arg instanceof FluentNone) {
            return new FluentNone(`NUMBER(${arg.valueOf()})`);
        }
        if (arg instanceof FluentNumber) {
            return new FluentNumber(arg.valueOf(), {
                ...arg.opts,
                ...values(opts, NUMBER_ALLOWED),
            });
        }
        if (arg instanceof FluentDateTime) {
            return new FluentNumber(arg.toNumber(), {
                ...values(opts, NUMBER_ALLOWED),
            });
        }
        throw new TypeError("Invalid argument to NUMBER");
    }
    const DATETIME_ALLOWED = [
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
        "timeZoneName",
    ];
    /**
     * The implementation of the `DATETIME()` builtin available to translations.
     *
     * Translations may call the `DATETIME()` builtin in order to specify
     * formatting options of a number. For example:
     *
     *     now = It's {DATETIME($today, month: "long")}.
     *
     * The implementation expects an array of {@link FluentValue | FluentValues} representing the
     * positional arguments, and an object of named {@link FluentValue | FluentValues} representing the
     * named parameters.
     *
     * The following options are recognized:
     *
     *     dateStyle
     *     timeStyle
     *     fractionalSecondDigits
     *     dayPeriod
     *     hour12
     *     weekday
     *     era
     *     year
     *     month
     *     day
     *     hour
     *     minute
     *     second
     *     timeZoneName
     *
     * Other options are ignored.
     *
     * @param args The positional arguments passed to this `DATETIME()`.
     * @param opts The named argments passed to this `DATETIME()`.
     */
    function DATETIME(args, opts) {
        let arg = args[0];
        if (arg instanceof FluentNone) {
            return new FluentNone(`DATETIME(${arg.valueOf()})`);
        }
        if (arg instanceof FluentDateTime || arg instanceof FluentNumber) {
            return new FluentDateTime(arg, values(opts, DATETIME_ALLOWED));
        }
        throw new TypeError("Invalid argument to DATETIME");
    }

    const cache = new Map();
    function getMemoizerForLocale(locales) {
        const stringLocale = Array.isArray(locales) ? locales.join(" ") : locales;
        let memoizer = cache.get(stringLocale);
        if (memoizer === undefined) {
            memoizer = new Map();
            cache.set(stringLocale, memoizer);
        }
        return memoizer;
    }

    /**
     * Message bundles are single-language stores of translation resources. They are
     * responsible for formatting message values and attributes to strings.
     */
    class FluentBundle {
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
        constructor(locales, { functions, useIsolating = true, transform = (v) => v, } = {}) {
            /** @ignore */
            this._terms = new Map();
            /** @ignore */
            this._messages = new Map();
            this.locales = Array.isArray(locales) ? locales : [locales];
            this._functions = {
                NUMBER,
                DATETIME,
                ...functions,
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
        addResource(res, { allowOverrides = false, } = {}) {
            const errors = [];
            for (let i = 0; i < res.body.length; i++) {
                let entry = res.body[i];
                if (entry.id.startsWith("-")) {
                    // Identifiers starting with a dash (-) define terms. Terms are private
                    // and cannot be retrieved from FluentBundle.
                    if (allowOverrides === false && this._terms.has(entry.id)) {
                        errors.push(new Error(`Attempt to override an existing term: "${entry.id}"`));
                        continue;
                    }
                    this._terms.set(entry.id, entry);
                }
                else {
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
            // Resolve a simple pattern without creating a scope. No error handling is
            // required; by definition simple patterns don't have placeables.
            if (typeof pattern === "string") {
                return this._transform(pattern);
            }
            // Resolve a complex pattern.
            let scope = new Scope(this, errors, args);
            try {
                let value = resolveComplexPattern(scope, pattern);
                return value.toString(scope);
            }
            catch (err) {
                if (scope.errors && err instanceof Error) {
                    scope.errors.push(err);
                    return new FluentNone().toString(scope);
                }
                throw err;
            }
        }
    }

    // This regex is used to iterate through the beginnings of messages and terms.
    // With the /m flag, the ^ matches at the beginning of every line.
    const RE_MESSAGE_START = /^(-?[a-zA-Z][\w-]*) *= */gm;
    // Both Attributes and Variants are parsed in while loops. These regexes are
    // used to break out of them.
    const RE_ATTRIBUTE_START = /\.([a-zA-Z][\w-]*) *= */y;
    const RE_VARIANT_START = /\*?\[/y;
    const RE_NUMBER_LITERAL = /(-?[0-9]+(?:\.([0-9]+))?)/y;
    const RE_IDENTIFIER = /([a-zA-Z][\w-]*)/y;
    const RE_REFERENCE = /([$-])?([a-zA-Z][\w-]*)(?:\.([a-zA-Z][\w-]*))?/y;
    const RE_FUNCTION_NAME = /^[A-Z][A-Z0-9_-]*$/;
    // A "run" is a sequence of text or string literal characters which don't
    // require any special handling. For TextElements such special characters are: {
    // (starts a placeable), and line breaks which require additional logic to check
    // if the next line is indented. For StringLiterals they are: \ (starts an
    // escape sequence), " (ends the literal), and line breaks which are not allowed
    // in StringLiterals. Note that string runs may be empty; text runs may not.
    const RE_TEXT_RUN = /([^{}\n\r]+)/y;
    const RE_STRING_RUN = /([^\\"\n\r]*)/y;
    // Escape sequences.
    const RE_STRING_ESCAPE = /\\([\\"])/y;
    const RE_UNICODE_ESCAPE = /\\u([a-fA-F0-9]{4})|\\U([a-fA-F0-9]{6})/y;
    // Used for trimming TextElements and indents.
    const RE_LEADING_NEWLINES = /^\n+/;
    const RE_TRAILING_SPACES = / +$/;
    // Used in makeIndent to strip spaces from blank lines and normalize CRLF to LF.
    const RE_BLANK_LINES = / *\r?\n/g;
    // Used in makeIndent to measure the indentation.
    const RE_INDENT = /( *)$/;
    // Common tokens.
    const TOKEN_BRACE_OPEN = /{\s*/y;
    const TOKEN_BRACE_CLOSE = /\s*}/y;
    const TOKEN_BRACKET_OPEN = /\[\s*/y;
    const TOKEN_BRACKET_CLOSE = /\s*] */y;
    const TOKEN_PAREN_OPEN = /\s*\(\s*/y;
    const TOKEN_ARROW = /\s*->\s*/y;
    const TOKEN_COLON = /\s*:\s*/y;
    // Note the optional comma. As a deviation from the Fluent EBNF, the parser
    // doesn't enforce commas between call arguments.
    const TOKEN_COMMA = /\s*,?\s*/y;
    const TOKEN_BLANK = /\s+/y;
    /**
     * Fluent Resource is a structure storing parsed localization entries.
     */
    class FluentResource {
        constructor(source) {
            this.body = [];
            RE_MESSAGE_START.lastIndex = 0;
            let cursor = 0;
            // Iterate over the beginnings of messages and terms to efficiently skip
            // comments and recover from errors.
            while (true) {
                let next = RE_MESSAGE_START.exec(source);
                if (next === null) {
                    break;
                }
                cursor = RE_MESSAGE_START.lastIndex;
                try {
                    this.body.push(parseMessage(next[1]));
                }
                catch (err) {
                    if (err instanceof SyntaxError) {
                        // Don't report any Fluent syntax errors. Skip directly to the
                        // beginning of the next message or term.
                        continue;
                    }
                    throw err;
                }
            }
            // The parser implementation is inlined below for performance reasons,
            // as well as for convenience of accessing `source` and `cursor`.
            // The parser focuses on minimizing the number of false negatives at the
            // expense of increasing the risk of false positives. In other words, it
            // aims at parsing valid Fluent messages with a success rate of 100%, but it
            // may also parse a few invalid messages which the reference parser would
            // reject. The parser doesn't perform any validation and may produce entries
            // which wouldn't make sense in the real world. For best results users are
            // advised to validate translations with the fluent-syntax parser
            // pre-runtime.
            // The parser makes an extensive use of sticky regexes which can be anchored
            // to any offset of the source string without slicing it. Errors are thrown
            // to bail out of parsing of ill-formed messages.
            function test(re) {
                re.lastIndex = cursor;
                return re.test(source);
            }
            // Advance the cursor by the char if it matches. May be used as a predicate
            // (was the match found?) or, if errorClass is passed, as an assertion.
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
            // Advance the cursor by the token if it matches. May be used as a predicate
            // (was the match found?) or, if errorClass is passed, as an assertion.
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
            // Execute a regex, advance the cursor, and return all capture groups.
            function match(re) {
                re.lastIndex = cursor;
                let result = re.exec(source);
                if (result === null) {
                    throw new SyntaxError(`Expected ${re.toString()}`);
                }
                cursor = re.lastIndex;
                return result;
            }
            // Execute a regex, advance the cursor, and return the capture group.
            function match1(re) {
                return match(re)[1];
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
                let attrs = Object.create(null);
                while (test(RE_ATTRIBUTE_START)) {
                    let name = match1(RE_ATTRIBUTE_START);
                    let value = parsePattern();
                    if (value === null) {
                        throw new SyntaxError("Expected attribute value");
                    }
                    attrs[name] = value;
                }
                return attrs;
            }
            function parsePattern() {
                let first;
                // First try to parse any simple text on the same line as the id.
                if (test(RE_TEXT_RUN)) {
                    first = match1(RE_TEXT_RUN);
                }
                // If there's a placeable on the first line, parse a complex pattern.
                if (source[cursor] === "{" || source[cursor] === "}") {
                    // Re-use the text parsed above, if possible.
                    return parsePatternElements(first ? [first] : [], Infinity);
                }
                // RE_TEXT_VALUE stops at newlines. Only continue parsing the pattern if
                // what comes after the newline is indented.
                let indent = parseIndent();
                if (indent) {
                    if (first) {
                        // If there's text on the first line, the blank block is part of the
                        // translation content in its entirety.
                        return parsePatternElements([first, indent], indent.length);
                    }
                    // Otherwise, we're dealing with a block pattern, i.e. a pattern which
                    // starts on a new line. Discrad the leading newlines but keep the
                    // inline indent; it will be used by the dedentation logic.
                    indent.value = trim(indent.value, RE_LEADING_NEWLINES);
                    return parsePatternElements([indent], indent.length);
                }
                if (first) {
                    // It was just a simple inline text after all.
                    return trim(first, RE_TRAILING_SPACES);
                }
                return null;
            }
            // Parse a complex pattern as an array of elements.
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
                // Trim the trailing spaces in the last element if it's a TextElement.
                if (typeof lastElement === "string") {
                    elements[lastIndex] = trim(lastElement, RE_TRAILING_SPACES);
                }
                let baked = [];
                for (let element of elements) {
                    if (element instanceof Indent) {
                        // Dedent indented lines by the maximum common indent.
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
                        ...variants,
                    };
                }
                throw new SyntaxError("Unclosed placeable");
            }
            function parseInlineExpression() {
                if (source[cursor] === "{") {
                    // It's a nested placeable.
                    return parsePlaceable();
                }
                if (test(RE_REFERENCE)) {
                    let [, sigil, name, attr = null] = match(RE_REFERENCE);
                    if (sigil === "$") {
                        return { type: "var", name };
                    }
                    if (consumeToken(TOKEN_PAREN_OPEN)) {
                        let args = parseArguments();
                        if (sigil === "-") {
                            // A parameterized term: -term(...).
                            return { type: "term", name, attr, args };
                        }
                        if (RE_FUNCTION_NAME.test(name)) {
                            return { type: "func", name, args };
                        }
                        throw new SyntaxError("Function names must be all upper-case");
                    }
                    if (sigil === "-") {
                        // A non-parameterized term: -term.
                        return {
                            type: "term",
                            name,
                            attr,
                            args: [],
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
                        case ")": // End of the argument list.
                            cursor++;
                            return args;
                        case undefined: // EOF
                            throw new SyntaxError("Unclosed argument list");
                    }
                    args.push(parseArgument());
                    // Commas between arguments are treated as whitespace.
                    consumeToken(TOKEN_COMMA);
                }
            }
            function parseArgument() {
                let expr = parseInlineExpression();
                if (expr.type !== "mesg") {
                    return expr;
                }
                if (consumeToken(TOKEN_COLON)) {
                    // The reference is the beginning of a named argument.
                    return {
                        type: "narg",
                        name: expr.name,
                        value: parseLiteral(),
                    };
                }
                // It's a regular message reference.
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
                if (star === undefined) {
                    throw new SyntaxError("Expected default variant");
                }
                return { variants, star };
            }
            function parseVariantKey() {
                consumeToken(TOKEN_BRACKET_OPEN, SyntaxError);
                let key;
                if (test(RE_NUMBER_LITERAL)) {
                    key = parseNumberLiteral();
                }
                else {
                    key = {
                        type: "str",
                        value: match1(RE_IDENTIFIER),
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
                let [, value, fraction = ""] = match(RE_NUMBER_LITERAL);
                let precision = fraction.length;
                return {
                    type: "num",
                    value: parseFloat(value),
                    precision,
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
                    // We've reached an EOL of EOF.
                    throw new SyntaxError("Unclosed string literal");
                }
            }
            // Unescape known escape sequences.
            function parseEscapeSequence() {
                if (test(RE_STRING_ESCAPE)) {
                    return match1(RE_STRING_ESCAPE);
                }
                if (test(RE_UNICODE_ESCAPE)) {
                    let [, codepoint4, codepoint6] = match(RE_UNICODE_ESCAPE);
                    let codepoint = parseInt(codepoint4 || codepoint6, 16);
                    return codepoint <= 0xd7ff || 0xe000 <= codepoint
                        ? // It's a Unicode scalar value.
                            String.fromCodePoint(codepoint)
                        : // Lonely surrogates can cause trouble when the parsing result is
                            // saved using UTF-8. Use U+FFFD REPLACEMENT CHARACTER instead.
                            "�";
                }
                throw new SyntaxError("Unknown escape sequence");
            }
            // Parse blank space. Return it if it looks like indent before a pattern
            // line. Skip it othwerwise.
            function parseIndent() {
                let start = cursor;
                consumeToken(TOKEN_BLANK);
                // Check the first non-blank character after the indent.
                switch (source[cursor]) {
                    case ".":
                    case "[":
                    case "*":
                    case "}":
                    case undefined: // EOF
                        // A special character. End the Pattern.
                        return false;
                    case "{":
                        // Placeables don't require indentation (in EBNF: block-placeable).
                        // Continue the Pattern.
                        return makeIndent(source.slice(start, cursor));
                }
                // If the first character on the line is not one of the special characters
                // listed above, it's a regular text character. Check if there's at least
                // one space of indent before it.
                if (source[cursor - 1] === " ") {
                    // It's an indented text character (in EBNF: indented-char). Continue
                    // the Pattern.
                    return makeIndent(source.slice(start, cursor));
                }
                // A not-indented text character is likely the identifier of the next
                // message. End the Pattern.
                return false;
            }
            // Trim blanks in text according to the given regex.
            function trim(text, re) {
                return text.replace(re, "");
            }
            // Normalize a blank block and extract the indent details.
            function makeIndent(blank) {
                let value = blank.replace(RE_BLANK_LINES, "\n");
                let length = RE_INDENT.exec(blank)[1].length;
                return new Indent(value, length);
            }
        }
    }
    class Indent {
        constructor(value, length) {
            this.value = value;
            this.length = length;
        }
    }

    exports.FluentBundle = FluentBundle;
    exports.FluentDateTime = FluentDateTime;
    exports.FluentNone = FluentNone;
    exports.FluentNumber = FluentNumber;
    exports.FluentResource = FluentResource;
    exports.FluentType = FluentType;

}));
