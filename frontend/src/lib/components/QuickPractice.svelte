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
      return errors.has(index) ? "text-red-600 drop-shadow-sm" : "text-orange-950 drop-shadow-sm";
    if (index === currentIndex)
      return "bg-orange-300/50 border-b-4 border-orange-500 animate-pulse text-amber-800";
    return "text-zinc-600";
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<div
  class="bg-white/60 backdrop-blur-md border border-orange-300 rounded-2xl p-6 sm:p-8 mb-8 shadow-xl shadow-orange-900/5"
>
  <div class="mb-4">
    <h3 class="text-xl font-bold text-amber-600 drop-shadow-sm m-0">Latihan Dasar</h3>
  </div>

  <div
    class="bg-white/40 border border-orange-300 rounded-xl p-6 sm:p-8 font-mono text-3xl sm:text-5xl flex flex-wrap justify-center gap-1 sm:gap-3 shadow-inner"
  >
    {#each sequence as char, i}
      <span
        class="transition-all duration-100 rounded w-10 h-12 sm:w-14 sm:h-16 flex items-center justify-center {getCharClass(
          i,
        )}">{char === " " ? "␣" : char}</span
      >
    {/each}
  </div>

  <p class="text-orange-800 text-sm mt-6 text-center font-medium m-0">
    Ketik huruf di atas secara berurutan. Otomatis berulang!
  </p>
</div>
