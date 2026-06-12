<script lang="ts">
	import type { Lesson, LessonProgress } from '$lib/types';
	import Stars from './Stars.svelte';

	let { lesson, progress, locked = false }: {
		lesson: Lesson;
		progress: LessonProgress;
		locked?: boolean;
	} = $props();
</script>

<a
	href={locked ? undefined : `/practice/${lesson.id}`}
	class="lesson-card"
	class:locked
	class:completed={progress.completed}
>
	<div class="card-icon">{locked ? '🔒' : lesson.icon}</div>
	<div class="card-content">
		<h3 class="card-title">{lesson.title}</h3>
		<p class="card-desc">{lesson.description}</p>
		{#if !locked}
			<div class="card-footer">
				<Stars stars={progress.stars} />
				{#if progress.bestWPM > 0}
					<span class="card-wpm">{progress.bestWPM} WPM</span>
				{/if}
			</div>
		{:else}
			<p class="card-lock-text">Butuh {lesson.unlockScore} bintang</p>
		{/if}
	</div>
</a>

<style>
	.lesson-card {
		display: flex;
		gap: 16px;
		padding: 20px;
		background: rgba(255, 255, 255, 0.05);
		border: 2px solid rgba(255, 255, 255, 0.1);
		border-radius: 16px;
		text-decoration: none;
		color: #f8fafc;
		transition: all 0.3s ease;
		cursor: pointer;
	}

	.lesson-card:not(.locked):hover {
		background: rgba(255, 255, 255, 0.1);
		border-color: #fbbf24;
		transform: translateY(-2px);
		box-shadow: 0 8px 24px rgba(251, 191, 36, 0.2);
	}

	.lesson-card.locked {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.lesson-card.completed {
		border-color: rgba(74, 222, 128, 0.3);
	}

	.card-icon {
		font-size: 40px;
		flex-shrink: 0;
		width: 60px;
		height: 60px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(255, 255, 255, 0.05);
		border-radius: 12px;
	}

	.card-content {
		flex: 1;
		min-width: 0;
	}

	.card-title {
		font-size: 18px;
		font-weight: 700;
		margin: 0 0 4px 0;
	}

	.card-desc {
		font-size: 13px;
		color: #64748b;
		margin: 0 0 8px 0;
	}

	.card-footer {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.card-wpm {
		font-size: 12px;
		color: #fbbf24;
		font-weight: 600;
	}

	.card-lock-text {
		font-size: 12px;
		color: #64748b;
		margin: 0;
	}
</style>
