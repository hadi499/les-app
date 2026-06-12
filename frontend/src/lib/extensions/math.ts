import katex from 'katex';

export function renderMathContent(html: string): string {
  let out = html;

  out = out.replace(/\$\$([\s\S]*?)\$\$/g, (_m, latex) => {
    try {
      const rendered = katex.renderToString(latex.trim(), { throwOnError: false, displayMode: true });
      return `<div style="text-align:center;margin:4px 0">${rendered}</div>`;
    } catch {
      return _m;
    }
  });

  out = out.replace(/\$([^\$]+?)\$/g, (_m, latex) => {
    try {
      const rendered = katex.renderToString(latex.trim(), { throwOnError: false, displayMode: false });
      return rendered;
    } catch {
      return _m;
    }
  });

  return out;
}
