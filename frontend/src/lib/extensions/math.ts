import katex from 'katex';

export function renderMathContent(html: string): string {
  let out = html;

  // Remove anchor tags but keep their inner text
  out = out.replace(/<a\b[^>]*>(.*?)<\/a>/gi, '$1');

  // Unescape HTML entities inside the latex string before passing to KaTeX
  const unescapeHtml = (text: string) => {
    return text
      .replace(/&amp;/g, '&')
      .replace(/&lt;/g, '<')
      .replace(/&gt;/g, '>')
      .replace(/&quot;/g, '"')
      .replace(/&#39;/g, "'")
      .replace(/&nbsp;/g, ' ');
  };

  // Replace display math block first. Ignore math inside code/pre blocks.
  out = out.replace(/(<code\b[^>]*>[\s\S]*?<\/code>|<pre\b[^>]*>[\s\S]*?<\/pre>)|(\$\$([\s\S]*?)\$\$)/gi, (match, codeBlock, mathGroup, latex) => {
    if (codeBlock) return codeBlock;
    try {
      const cleanLatex = unescapeHtml(latex.trim());
      const rendered = katex.renderToString(cleanLatex, { throwOnError: false, displayMode: true });
      return `<div style="text-align:center;margin:4px 0">${rendered}</div>`;
    } catch {
      return match;
    }
  });

  // Replace inline math. Ignore math inside code/pre blocks, and ensure proper word boundaries.
  out = out.replace(/(<code\b[^>]*>[\s\S]*?<\/code>|<pre\b[^>]*>[\s\S]*?<\/pre>)|((^|[^a-zA-Z0-9_])\$(?!\s)([^$<]*?[^\s$<])\$(?=(?:th|st|nd|rd)?(?:[^a-zA-Z0-9_]|$)))/gi, (match, codeBlock, mathGroup, pre, latex) => {
    if (codeBlock) return codeBlock;
    try {
      const cleanLatex = unescapeHtml(latex.trim());
      const rendered = katex.renderToString(cleanLatex, { throwOnError: false, displayMode: false });
      return pre + rendered;
    } catch {
      return match;
    }
  });

  return out;
}
