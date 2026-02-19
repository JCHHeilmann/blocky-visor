<script lang="ts">
  interface Props {
    value: string;
    onchange: (value: string) => void;
    readonly?: boolean;
  }

  let { value, onchange, readonly = false }: Props = $props();

  let textarea: HTMLTextAreaElement | undefined = $state();
  let pre: HTMLPreElement | undefined = $state();

  function syncScroll() {
    if (textarea && pre) {
      pre.scrollTop = textarea.scrollTop;
      pre.scrollLeft = textarea.scrollLeft;
    }
  }

  function handleInput(e: Event) {
    const target = e.target as HTMLTextAreaElement;
    onchange(target.value);
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Tab") {
      e.preventDefault();
      const target = e.target as HTMLTextAreaElement;
      const start = target.selectionStart;
      const end = target.selectionEnd;
      const newValue = value.slice(0, start) + "  " + value.slice(end);
      onchange(newValue);
      requestAnimationFrame(() => {
        target.selectionStart = target.selectionEnd = start + 2;
      });
    }
  }

  function escapeHtml(str: string): string {
    return str
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;");
  }

  function highlightValue(raw: string): string {
    const trimmed = raw.trim();
    if (!trimmed || trimmed === "") return escapeHtml(raw);

    // Inline comment after value
    const commentIdx = raw.indexOf(" #");
    if (commentIdx > 0) {
      const before = raw.slice(0, commentIdx);
      const comment = raw.slice(commentIdx);
      return highlightValue(before) + `<span class="yl-comment">${escapeHtml(comment)}</span>`;
    }

    // Quoted strings
    if (/^\s*["']/.test(raw)) {
      return `<span class="yl-string">${escapeHtml(raw)}</span>`;
    }

    // Booleans
    if (/^\s*(true|false|yes|no|on|off)$/i.test(raw)) {
      return `<span class="yl-bool">${escapeHtml(raw)}</span>`;
    }

    // Null
    if (/^\s*(null|~)$/i.test(raw)) {
      return `<span class="yl-null">${escapeHtml(raw)}</span>`;
    }

    // Numbers
    if (/^\s*-?(\d+\.?\d*|\.\d+)([eE][+-]?\d+)?$/.test(trimmed)) {
      return `<span class="yl-number">${escapeHtml(raw)}</span>`;
    }

    // Duration-like values (e.g. 5m, 1h, 30s)
    if (/^\s*\d+[smhd]$/i.test(trimmed)) {
      return `<span class="yl-number">${escapeHtml(raw)}</span>`;
    }

    return `<span class="yl-value">${escapeHtml(raw)}</span>`;
  }

  function highlightYaml(code: string): string {
    return code
      .split("\n")
      .map((line) => {
        // Full-line comment
        if (/^\s*#/.test(line)) {
          return `<span class="yl-comment">${escapeHtml(line)}</span>`;
        }

        // Document markers
        if (/^---\s*$/.test(line) || /^\.\.\.\s*$/.test(line)) {
          return `<span class="yl-doc">${escapeHtml(line)}</span>`;
        }

        // Key: value pairs
        const kvMatch = line.match(/^(\s*)([\w][\w.\-/]*)(:\s*)(.*)/);
        if (kvMatch) {
          const [, indent, key, colon, rest] = kvMatch;
          const highlighted =
            escapeHtml(indent) +
            `<span class="yl-key">${escapeHtml(key)}</span>` +
            `<span class="yl-colon">${escapeHtml(colon)}</span>` +
            (rest ? highlightValue(rest) : "");
          return highlighted;
        }

        // List items with key: value
        const listKvMatch = line.match(/^(\s*)(- )(\s*)([\w][\w.\-/]*)(:\s*)(.*)/);
        if (listKvMatch) {
          const [, indent, dash, space, key, colon, rest] = listKvMatch;
          return (
            escapeHtml(indent) +
            `<span class="yl-dash">${escapeHtml(dash)}</span>` +
            escapeHtml(space) +
            `<span class="yl-key">${escapeHtml(key)}</span>` +
            `<span class="yl-colon">${escapeHtml(colon)}</span>` +
            (rest ? highlightValue(rest) : "")
          );
        }

        // Plain list items
        const listMatch = line.match(/^(\s*)(- )(.*)/);
        if (listMatch) {
          const [, indent, dash, rest] = listMatch;
          return (
            escapeHtml(indent) +
            `<span class="yl-dash">${escapeHtml(dash)}</span>` +
            highlightValue(rest)
          );
        }

        return escapeHtml(line);
      })
      .join("\n");
  }

  let highlighted = $derived(highlightYaml(value));
</script>

<div class="yaml-editor">
  <pre
    bind:this={pre}
    class="yaml-highlight"
    aria-hidden="true"
  >{@html highlighted + "\n"}</pre>
  <textarea
    bind:this={textarea}
    {value}
    oninput={handleInput}
    onscroll={syncScroll}
    onkeydown={handleKeydown}
    {readonly}
    spellcheck="false"
    autocomplete="off"
    class="yaml-input"
  ></textarea>
</div>

<style>
  .yaml-editor {
    position: relative;
    width: 100%;
    height: 100%;
    min-height: 0;
    overflow: hidden;
    border-radius: 0.5rem;
    border: 1px solid var(--color-surface-border);
    background: var(--color-surface-bg);
  }

  .yaml-highlight,
  .yaml-input {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    margin: 0;
    padding: 1rem;
    font-family: var(--font-mono);
    font-size: 0.8125rem;
    line-height: 1.6;
    tab-size: 2;
    white-space: pre;
    overflow: auto;
    box-sizing: border-box;
  }

  .yaml-highlight {
    color: var(--color-text-primary);
    pointer-events: none;
    z-index: 1;
  }

  .yaml-input {
    color: transparent;
    caret-color: var(--color-text-primary);
    background: transparent;
    border: none;
    outline: none;
    resize: none;
    z-index: 2;
    -webkit-text-fill-color: transparent;
  }

  .yaml-input::selection {
    background: oklch(0.5 0.13 195 / 0.3);
    -webkit-text-fill-color: transparent;
  }

  .yaml-input:focus {
    outline: none;
  }

  /* Syntax colors - dark mode defaults */
  :global(.yl-comment) {
    color: var(--color-gray-500);
    font-style: italic;
  }

  :global(.yl-key) {
    color: oklch(0.75 0.15 195);
  }

  :global(.yl-colon) {
    color: var(--color-gray-500);
  }

  :global(.yl-string) {
    color: oklch(0.72 0.14 150);
  }

  :global(.yl-bool) {
    color: oklch(0.72 0.15 310);
  }

  :global(.yl-number) {
    color: oklch(0.75 0.14 60);
  }

  :global(.yl-null) {
    color: var(--color-gray-500);
    font-style: italic;
  }

  :global(.yl-value) {
    color: var(--color-text-primary);
  }

  :global(.yl-dash) {
    color: var(--color-gray-400);
  }

  :global(.yl-doc) {
    color: var(--color-gray-400);
  }

  /* Light mode overrides */
  :root:not(.dark) :global(.yl-key) {
    color: oklch(0.45 0.18 195);
  }

  :root:not(.dark) :global(.yl-string) {
    color: oklch(0.42 0.16 150);
  }

  :root:not(.dark) :global(.yl-bool) {
    color: oklch(0.48 0.18 310);
  }

  :root:not(.dark) :global(.yl-number) {
    color: oklch(0.48 0.16 60);
  }
</style>
