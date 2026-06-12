<script lang="ts">
	let { stars, maxStars = 3 }: { stars: number; maxStars?: number } = $props();
</script>

<div class="stars">
	{#each Array(maxStars) as _, i}
		<span class="star" class:filled={i < stars}>
			{i < stars ? '⭐' : '☆'}
		</span>
	{/each}
</div>

<style>
	.stars {
		display: flex;
		gap: 4px;
	}

	.star {
		font-size: 20px;
		transition: transform 0.2s ease;
	}

	.star.filled {
		animation: pop 0.3s ease;
	}

	@keyframes pop {
		0% { transform: scale(0); }
		50% { transform: scale(1.3); }
		100% { transform: scale(1); }
	}
</style>
