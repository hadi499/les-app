<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { EditorView, basicSetup } from 'codemirror';
  import { EditorState } from '@codemirror/state';
  import { html } from '@codemirror/lang-html';
  import { css } from '@codemirror/lang-css';
  import { javascript } from '@codemirror/lang-javascript';
  import { python } from '@codemirror/lang-python';

  let { 
    value = $bindable(''), 
    language = 'html' 
  } = $props<{
    value?: string;
    language?: 'html' | 'css' | 'js' | 'python';
  }>();

  let editorContainer: HTMLElement;
  let view: EditorView;

  function getLanguageExtension(lang: string) {
    switch (lang) {
      case 'html': return html();
      case 'css': return css();
      case 'js': return javascript();
      case 'python': return python();
      default: return html();
    }
  }

  onMount(() => {
    const updateListener = EditorView.updateListener.of((update) => {
      if (update.docChanged) {
        value = update.state.doc.toString();
      }
    });

    const state = EditorState.create({
      doc: value,
      extensions: [
        basicSetup,
        getLanguageExtension(language),
        updateListener,
        EditorView.theme({
          "&": { height: "100%", fontSize: "14px" },
          ".cm-scroller": { overflow: "auto" }
        })
      ]
    });

    view = new EditorView({
      state,
      parent: editorContainer
    });
  });

  $effect(() => {
    if (view && value !== undefined) {
      const currentValue = view.state.doc.toString();
      if (currentValue !== value) {
        view.dispatch({
          changes: {
            from: 0,
            to: view.state.doc.length,
            insert: value
          }
        });
      }
    }
  });

  onDestroy(() => {
    if (view) {
      view.destroy();
    }
  });
</script>

<div bind:this={editorContainer} class="h-full w-full border border-gray-300 rounded overflow-hidden">
</div>

<style>
  div :global(.cm-editor) {
    height: 100%;
  }
</style>
