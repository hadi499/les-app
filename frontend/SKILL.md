---
name: live-coding-editor
description: Build or modify an in-browser live coding editor for teaching basic programming (HTML/CSS/JS and Python) in a SvelteKit app. Use this skill whenever the user asks to add a "live coding" feature, an interactive code playground, an in-browser code editor with live preview, or a coding exercise/practice tool for beginners — especially for the lesbalonggarut.my.id educational site. Covers editor setup (CodeMirror 6), sandboxed HTML/CSS/JS preview via iframe, Python execution via Pyodide, and beginner-friendly UX patterns (starter templates, simplified errors, expected-output checking).
---

# Live Coding Editor (Beginner-Focused)

A skill for building an in-browser code editor + live preview feature aimed at people learning to code for the first time. Built for a SvelteKit frontend; no backend execution needed for the languages covered here (HTML/CSS/JS run natively in the browser, Python runs via WebAssembly).

## When to use this

- User wants to add a "coba kode langsung" / live preview / code playground feature to a SvelteKit site
- User wants a coding practice/exercise tool with instant feedback for beginners
- User mentions lesbalonggarut.my.id and wants to add a coding module alongside its existing typing practice and quizzes
- User wants to embed a Python or JS sandbox for teaching purposes

## Design principles for this context

The audience is **absolute beginners**, not developers. This changes defaults from a general-purpose code playground:

1. **Never start from a blank editor.** Always pre-fill a small starter snippet relevant to the lesson.
2. **Simplify errors.** Raw stack traces (`Uncaught TypeError: Cannot read properties of undefined...`) are intimidating. Catch errors and rewrite them in plain Indonesian where possible, e.g. "Ada bagian kode yang belum lengkap, coba cek baris ke-X".
3. **Show output immediately and visually.** Split-pane: code on one side, rendered result or console output on the other, updated live (debounced).
4. **No login/backend required.** Persist progress with `localStorage` so the feature stays free to host and simple to ship.
5. **Prefer CodeMirror 6 over Monaco.** Monaco looks like a full IDE (VS Code), which can be intimidating; CodeMirror is lighter and friendlier for beginners, and has a smaller bundle size.

## Architecture

### HTML / CSS / JS (fully client-side, no backend)

1. Three editor panels (or tabbed single panel): HTML, CSS, JS
2. On change (debounce ~300ms), combine into one HTML document
3. Render inside a sandboxed iframe using `srcdoc`

```svelte
<script>
  let html = $state('');
  let css = $state('');
  let js = $state('');

  let srcdoc = $derived(`
    <html>
      <head><style>${css}</style></head>
      <body>${html}
        <script>${js}<\/script>
      </body>
    </html>
  `);
</script>

<iframe sandbox="allow-scripts" srcdoc={srcdoc} title="preview"></iframe>
```

Important: use `sandbox="allow-scripts"` **without** `allow-same-origin` — this keeps user-written JS isolated from the parent page.

### Python (via Pyodide, fully client-side, no backend)

Pyodide compiles CPython to WebAssembly and runs entirely in the browser — no server execution needed, which keeps hosting costs at zero for this feature.

1. Load Pyodide once (lazy-load on first use, since the WASM bundle is ~10-20MB):
   ```js
   const pyodide = await loadPyodide();
   ```
2. Redirect `print()` output to a console box:
   ```js
   pyodide.setStdout({ batched: (msg) => appendToConsole(msg) });
   ```
3. Run user code inside a try/catch, and simplify the error message before displaying it.
4. Debounce execution the same way as the JS panel, or use an explicit "Run" button for Python (execution is heavier than JS, so on-every-keystroke is not recommended — prefer a Run button over live-as-you-type for this language specifically).

### Beginner UX additions

- **Starter templates**: keyed by exercise/lesson id, pre-filled on load
- **Reset button**: restores the starter template, discarding user changes
- **Expected-output checking** (optional): for structured exercises, compare captured console/DOM output against an expected value and show a simple pass/fail indicator — useful for turning this into a graded exercise rather than a free-form playground
- **Progress persistence**: save last-edited code per exercise to `localStorage`, keyed by exercise id, so users don't lose work on refresh

## What NOT to add for this use case

- No Monaco (too heavy/complex for beginners)
- No backend code execution / Docker sandboxes — unnecessary complexity and cost for HTML/CSS/JS/Python at this level
- No real-time multi-user collaboration — out of scope unless explicitly requested
- No arbitrary npm package imports in the JS sandbox — keep the surface area small and predictable for beginners

## Reference: minimal file layout in a SvelteKit project

```
src/lib/components/LiveEditor/
  Editor.svelte       - CodeMirror wrapper, one instance per language panel
  Preview.svelte       - the sandboxed iframe (HTML/CSS/JS)
  PythonRunner.svelte  - Pyodide loader + run button + console output
  starters.ts          - map of exercise id -> starter code
src/routes/belajar-koding/[exercise]/+page.svelte
  - wires Editor + Preview/PythonRunner together based on the exercise's language
```
