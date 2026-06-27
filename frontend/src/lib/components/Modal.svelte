<script lang="ts">
  import type { Snippet } from "svelte";

  let {
    show = false,
    onclose,
    title = "",
    maxWidth = "max-w-2xl",
    children,
  }: {
    show: boolean;
    onclose: () => void;
    title?: string;
    maxWidth?: string;
    children: Snippet;
  } = $props();
</script>

{#if show}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4"
  >
    <div
      class="bg-slate-50 border border-slate-200 rounded-xl shadow-2xl {maxWidth} w-full max-h-[90vh] overflow-y-auto"
    >
      {#if title}
        <div
          class="flex items-center justify-between p-4 border-b border-slate-200"
        >
          <h2 class="text-lg font-semibold text-slate-900 drop-shadow-sm">
            {title}
          </h2>
          <button
            onclick={onclose}
            class="w-7 h-7 flex items-center justify-center rounded bg-red-500 hover:bg-red-500 cursor-pointer text-white transition-colors"
            >&times;</button
          >
        </div>
      {/if}
      <div class="p-4 text-slate-800">
        {@render children()}
      </div>
    </div>
  </div>
{/if}
