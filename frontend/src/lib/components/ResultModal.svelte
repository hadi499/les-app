<script lang="ts">
	let { show, wpm, accuracy, stars, onRetry, onNext, hasNext }: {
		show: boolean;
		wpm: number;
		accuracy: number;
		stars: number;
		onRetry: () => void;
		onNext: () => void;
		hasNext: boolean;
	} = $props();

	function getMessage(): string {
		if (stars === 3) return 'LUAR BIASA! 🎉';
		if (stars === 2) return 'Hebat! 👏';
		if (stars === 1) return 'Bagus! Terus latihan! 💪';
		return 'Coba lagi ya! 😊';
	}

	function getStarLabel(): string {
		if (stars === 3) return 'Sempurna!';
		if (stars === 2) return 'Keren!';
		if (stars === 1) return 'Lumayan!';
		return 'Latihan lagi!';
	}
</script>

{#if show}
	<div class="fixed inset-0 bg-black/60 flex items-center justify-center z-[100] backdrop-blur-[4px] animate-fadeIn" role="dialog" aria-modal="true">
		<div class="bg-gradient-to-br from-white/90 to-white/70 border border-orange-300/60 shadow-2xl shadow-orange-900/10 rounded-3xl p-8 sm:p-10 text-center max-w-[420px] w-[90%] animate-slideUp backdrop-blur-md">
			<h2 class="text-2xl sm:text-3xl font-extrabold mb-5 text-orange-950 drop-shadow-sm">{getMessage()}</h2>

			<div class="flex justify-center gap-3 mb-2">
				{#each Array(3) as _, i}
					<span class="text-4xl sm:text-5xl opacity-30 transition-all duration-300 drop-shadow-sm {i < stars ? 'opacity-100 animate-starPop drop-shadow-md' : ''}" style="animation-delay: {i * 0.2}s">
						{i < stars ? '⭐' : '☆'}
					</span>
				{/each}
			</div>
			<p class="text-base text-amber-500 font-bold mb-6 tracking-wide">{getStarLabel()}</p>

			<div class="flex gap-4 sm:gap-6 justify-center mb-8">
				<div class="flex flex-col items-center bg-white/40 px-5 sm:px-6 py-4 rounded-2xl shadow-inner border border-orange-300">
					<span class="text-3xl sm:text-4xl font-black text-amber-500 drop-shadow-sm">{wpm}</span>
					<span class="text-xs text-orange-800 uppercase tracking-[0.15em] mt-1 font-bold">WPM</span>
				</div>
				<div class="flex flex-col items-center bg-white/40 px-5 sm:px-6 py-4 rounded-2xl shadow-inner border border-orange-300">
					<span class="text-3xl sm:text-4xl font-black text-amber-500 drop-shadow-sm">{accuracy}%</span>
					<span class="text-xs text-orange-800 uppercase tracking-[0.15em] mt-1 font-bold">Akurasi</span>
				</div>
			</div>

			<div class="flex gap-3 justify-center mb-4">
				<button class="px-5 sm:px-7 py-3 rounded-xl border-2 border-orange-300 bg-white text-orange-950 text-sm sm:text-base font-bold cursor-pointer transition-all duration-200 hover:bg-orange-50 hover:border-orange-300 hover:-translate-y-0.5 active:translate-y-0 active:scale-95 shadow-sm" onclick={onRetry}>
					🔄 Ulangi
				</button>
				{#if hasNext && stars > 0}
					<button class="px-5 sm:px-7 py-3 rounded-xl border-none bg-gradient-to-r from-amber-400 to-orange-400 text-orange-950 text-sm sm:text-base font-bold cursor-pointer transition-all duration-200 hover:scale-105 hover:shadow-[0_8px_20px_rgba(251,191,36,0.3)] active:scale-95 shadow-md" onclick={onNext}>
						Lanjut ➡️
					</button>
				{/if}
			</div>

			<a href="/mengetik" class="text-sm text-orange-800 no-underline inline-block mt-3 transition-colors duration-200 hover:text-orange-950 font-medium">← Kembali ke Menu</a>
		</div>
	</div>
{/if}

<style>
	.animate-fadeIn {
		animation: fadeIn 0.3s ease forwards;
	}

	.animate-slideUp {
		animation: slideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1) forwards;
	}

	.animate-starPop {
		animation: starPop 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
	}

	@keyframes fadeIn {
		from { opacity: 0; }
		to { opacity: 1; }
	}

	@keyframes slideUp {
		from { transform: translateY(40px); opacity: 0; }
		to { transform: translateY(0); opacity: 1; }
	}

	@keyframes starPop {
		0% { transform: scale(0); }
		50% { transform: scale(1.3); }
		100% { transform: scale(1); }
	}
</style>
