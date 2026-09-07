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
        ? `${baseClass} text-red-600 drop-shadow-sm`
        : `${baseClass} text-orange-950 drop-shadow-sm`;
    }
    if (index === currentIndex) {
      return `${baseClass} bg-orange-300/50 border-b-[3px] border-orange-500 animate-blink text-amber-800`;
    }
    return `${baseClass} text-orange-700`;
  }
</script>

<div class="w-full">
  <div
    class="bg-white/60 backdrop-blur-md rounded-2xl p-4 md:p-6 font-mono text-[20px] md:text-[30px] leading-loose tracking-[2px] border border-orange-300 shadow-xl shadow-orange-900/5 min-h-[120px] break-all"
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
