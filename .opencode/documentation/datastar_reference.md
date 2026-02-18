# Datastar Reference Documentation

Fetched from https://data-star.dev/reference on 2026-02-17

## Key Attributes for Side Panel Implementation

### `data-signals`

Patches (adds, updates or removes) one or more signals into the existing signals. Values defined later in the DOM tree override those defined earlier.

```html
<div data-signals:foo="1"></div>
```

Signals can be nested using dot-notation.

```html
<div data-signals:foo.bar="1"></div>
```

The `data-signals` attribute can also be used to patch multiple signals using a set of key-value pairs:

```html
<div data-signals="{foo: {bar: 1, baz: 2}}"></div>
```

**Important**: Setting a signal's value to `null` or `undefined` removes the signal.

```html
<div data-signals="{foo: null}"></div>
```

Keys used in `data-signals:*` are converted to camel case, so the signal name `mySignal` must be written as `data-signals:my-signal` or `data-signals="{mySignal: 1}"`.

Signal names cannot begin with nor contain a double underscore (`__`), due to its use as a modifier delimiter.

**Modifiers**:
- `__case` – Converts the casing of the signal name
- `__ifmissing` – Only patches signals if their keys do not already exist (useful for defaults)

### `data-class`

Adds or removes a class to or from an element based on an expression.

```html
<div data-class:font-bold="$foo == 'strong'"></div>
```

If the expression evaluates to `true`, the class is added to the element; otherwise, it is removed.

The `data-class` attribute can also add or remove multiple classes:

```html
<div data-class="{success: $foo != '', 'font-bold': $foo == 'strong'}"></div>
```

**Modifiers**:
- `__case` – Converts the casing of the class name

### `data-on`

Attaches an event listener to an element, executing an expression whenever the event is triggered.

```html
<button data-on:click="$foo = ''">Reset</button>
```

An `evt` variable that represents the event object is available in the expression.

**Modifiers**:
- `__once` – Only trigger the event listener once
- `__passive` – Do not call `preventDefault` on the event listener
- `__prevent` – Calls `preventDefault` on the event listener
- `__stop` – Calls `stopPropagation` on the event listener
- `__viewtransition` – Wraps the expression in `document.startViewTransition()`
- `__delay` – Delay the event listener (`.500ms`, `.1s`, etc.)
- `__debounce` – Debounce the event listener
- `__throttle` – Throttle the event listener

### `data-text`

Binds the text content of an element to an expression.

```html
<div data-text="$foo"></div>
```

### `data-show`

Shows or hides an element based on whether an expression evaluates to `true` or `false`.

```html
<div data-show="$foo"></div>
```

To prevent flickering before Datastar processes the DOM, add `style="display: none"`:

```html
<div data-show="$foo" style="display: none"></div>
```

## Store Initialization

Stores can be initialized using `data-store` attribute on the `<body>` tag:

```html
<body data-store='{"sidepanel": {"collapsed": true}}'>
```

## Signal Naming and Usage

- Signal names are automatically converted to camelCase
- Use `$signalName` to reference signals in expressions
- Signals beginning with underscore `_` are not included in backend requests by default
- Nested signals can be accessed using dot notation: `$sidepanel.collapsed`

## Attribute Evaluation Order

Elements are evaluated by walking the DOM in a depth-first manner, and attributes are applied in the order they appear in the element. This is important for cases like using `data-indicator` with `data-init`.

## Error Handling

Datastar has built-in error handling. When a data attribute is used incorrectly, an error message is logged to the console with a "More info" link to a context-aware error page.

## Common Patterns for Side Panel

```html
<!-- Store initialization -->
<body data-store='{"sidepanel": {"collapsed": true}}'>

<!-- Menu toggle - only visible when collapsed -->
<button data-on:click="$sidepanel.collapsed = false" 
        data-class:hidden="!$sidepanel.collapsed">
</button>

<!-- Side panel with conditional class -->
<div class="sidepanel" data-class:sidepanel--collapsed="$sidepanel.collapsed">
  <!-- Close button -->
  <button data-on:click="$sidepanel.collapsed = true">✕</button>
  
  <!-- Debug state -->
  <div data-text='$sidepanel.collapsed ? "CLOSED" : "OPEN"'></div>
</div>

<!-- Overlay to close panel -->
<div data-on:click="$sidepanel.collapsed = true"
     data-class:visible="!$sidepanel.collapsed">
</div>
```

## Troubleshooting

1. **Signals not updating**: Ensure Datastar script is loaded
2. **Attributes not working**: Check console for Datastar errors
3. **State not persisting**: Verify `data-store` attribute syntax
4. **CSS not applying**: Check `data-class` expressions evaluate to boolean