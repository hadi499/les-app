<script lang="ts">
  import type { Snippet } from "svelte";

  let {
    show = false,
    onclose,
    title = "",
    children,
  }: {
    show: boolean;
    onclose: () => void;
    title?: string;
    children: Snippet;
  } = $props();
</script>

{#if show}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4"
  >
    <div
      class="bg-orange-50 border border-orange-200 rounded-xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto"
    >
      {#if title}
        <div
          class="flex items-center justify-between p-4 border-b border-orange-200"
        >
          <h2 class="text-lg font-semibold text-orange-950 drop-shadow-sm">
            {title}
          </h2>
          <button
            onclick={onclose}
            class="w-7 h-7 flex items-center justify-center rounded bg-red-500 hover:bg-red-500 cursor-pointer text-white transition-colors"
            >&times;</button
          >
        </div>
      {/if}
      <div class="p-4 text-orange-900">
        {@render children()}
      </div>
    </div>
  </div>
{/if}
