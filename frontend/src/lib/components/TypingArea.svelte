<script lang="ts">
  let {
    text,
    currentIndex,
    errors,
  }: {
    text: string;
    currentIndex: number;
    errors: Set<number>;
  } = $props();

  function getCharClass(index: number): string {
    const baseClass =
      "transition-all duration-100 ease-in-out rounded py-[2px] px-[1px]";
    if (index < currentIndex) {
      return errors.has(index)
        ? `${baseClass} text-red-400 drop-shadow-sm`
        : `${baseClass} text-white drop-shadow-sm`;
    }
    if (index === currentIndex) {
      return `${baseClass} bg-amber-400/20 border-b-[3px] border-amber-400 animate-blink text-amber-100`;
    }
    return `${baseClass} text-zinc-600`;
  }
</script>

<div class="w-full">
  <div
    class="bg-zinc-900/60 backdrop-blur-md rounded-2xl p-4 md:p-6 font-mono text-[20px] md:text-[30px] leading-loose tracking-[2px] border border-zinc-700 shadow-xl shadow-black/20 min-h-[120px] break-all"
  >
    {#each text as char, i}
      <span class={getCharClass(i)}>{char === " " ? "␣" : char}</span>
    {/each}
  </div>
</div>

<style>
  :global(.animate-blink) {
    animation: blink 1s ease infinite;
  }
  @keyframes blink {
    0%,
    100% {
      border-bottom-color: #fbbf24;
    }
    50% {
      border-bottom-color: transparent;
    }
  }
</style>
