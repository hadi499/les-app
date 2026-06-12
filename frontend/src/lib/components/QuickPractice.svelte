<script lang="ts">
  const sequence = "asdf jkl;";
  let currentIndex = $state(0);
  let errors = $state(new Set<number>());

  function reset() {
    currentIndex = 0;
    errors = new Set<number>();
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Backspace") {
      if (currentIndex > 0) {
        currentIndex--;
        errors.delete(currentIndex);
      }
      return;
    }

    if (e.key.length > 1) return; // Ignore modifier keys (Shift, Ctrl, etc.)
    
    e.preventDefault();

    const key = e.key.toLowerCase();
    const expected = sequence[currentIndex].toLowerCase();

    if (key !== expected) {
      errors.add(currentIndex);
    }

    currentIndex++;

    if (currentIndex >= sequence.length) {
      reset();
    }
  }

  function getCharClass(index: number): string {
    if (index < currentIndex)
      return errors.has(index) ? "text-red-400 drop-shadow-sm" : "text-white drop-shadow-sm";
    if (index === currentIndex)
      return "bg-amber-400/20 border-b-4 border-amber-400 animate-pulse text-amber-100";
    return "text-zinc-600";
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<div
  class="bg-zinc-900/60 backdrop-blur-md border border-zinc-700/50 rounded-2xl p-6 sm:p-8 mb-8 shadow-xl shadow-black/20"
>
  <div class="mb-4">
    <h3 class="text-xl font-bold text-amber-400 drop-shadow-sm m-0">Latihan Dasar</h3>
  </div>

  <div
    class="bg-zinc-800/50 border border-zinc-700 rounded-xl p-6 sm:p-8 font-mono text-3xl sm:text-5xl flex flex-wrap justify-center gap-1 sm:gap-3 shadow-inner"
  >
    {#each sequence as char, i}
      <span
        class="transition-all duration-100 rounded w-10 h-12 sm:w-14 sm:h-16 flex items-center justify-center {getCharClass(
          i,
        )}">{char === " " ? "␣" : char}</span
      >
    {/each}
  </div>

  <p class="text-blue-200 text-sm mt-6 text-center font-medium m-0">
    Ketik huruf di atas secara berurutan. Otomatis berulang!
  </p>
</div>
