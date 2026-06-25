import katex from 'katex';

export function renderMathContent(html: string): string {
  let out = html;

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

  out = out.replace(/\$\$([\s\S]*?)\$\$/g, (_m, latex) => {
    try {
      const cleanLatex = unescapeHtml(latex.trim());
      const rendered = katex.renderToString(cleanLatex, { throwOnError: false, displayMode: true });
      return `<div style="text-align:center;margin:4px 0">${rendered}</div>`;
    } catch {
      return _m;
    }
  });

  out = out.replace(/\$([^\$]+?)\$/g, (_m, latex) => {
    try {
      const cleanLatex = unescapeHtml(latex.trim());
      const rendered = katex.renderToString(cleanLatex, { throwOnError: false, displayMode: false });
      return rendered;
    } catch {
      return _m;
    }
  });

  return out;
}
