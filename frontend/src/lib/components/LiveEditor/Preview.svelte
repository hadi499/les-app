<script lang="ts">
  let { 
    html = '', 
    css = '', 
    js = '' 
  } = $props<{
    html?: string;
    css?: string;
    js?: string;
  }>();

  let srcdoc = $derived(`
    <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8">
        <style>
          body { font-family: sans-serif; }
          ${css}
        </style>
      </head>
      <body>
        ${html}
        <script>
          try {
            ${js}
          } catch (e) {
            console.error(e);
            document.body.innerHTML += '<div style="color: red; padding: 10px; border: 1px solid red; margin-top: 20px;">' + e.toString() + '</div>';
          }
        <\/script>
      </body>
    </html>
  `);
</script>

<div class="h-full w-full border border-gray-300 rounded overflow-hidden bg-white">
  <iframe
    sandbox="allow-scripts"
    {srcdoc}
    title="preview"
    class="w-full h-full border-none"
  ></iframe>
</div>
