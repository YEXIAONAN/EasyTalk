import { Marked } from 'marked'
import hljs from 'highlight.js/lib/common'
import DOMPurify from 'dompurify'
import 'highlight.js/styles/github-dark.css'

const marked = new Marked({
  gfm: true,
  breaks: true,
  renderer: {
    code({ text, lang }) {
      const language = lang && hljs.getLanguage(lang) ? lang : ''
      let highlighted: string
      try {
        highlighted = language
          ? hljs.highlight(text, { language }).value
          : hljs.highlightAuto(text).value
      } catch {
        highlighted = escapeHtml(text)
      }
      const label = language || 'text'
      return (
        `<div class="code-block">` +
        `<div class="code-header"><span class="code-lang">${escapeHtml(label)}</span>` +
        `<button type="button" class="copy-btn">${escapeHtml(copyLabel)}</button></div>` +
        `<pre><code class="hljs">${highlighted}</code></pre></div>`
      )
    },
  },
})

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

// The copy button label is localized at render time. It is stored in a
// module-level variable because the marked renderer is configured once; since
// parsing is synchronous the override is safe.
let copyLabel = 'Copy'

export function renderMarkdown(content: string, label?: string): string {
  if (label !== undefined) copyLabel = label
  const raw = marked.parse(content, { async: false }) as string
  return DOMPurify.sanitize(raw)
}