<script lang="ts">
	let { mode = 'both' }: { mode?: 'left' | 'right' | 'both' } = $props();

	const leftKeys = ['a', 's', 'd', 'f'];
	const rightKeys = ['j', 'k', 'l', ';'];
	const allKeys = [...leftKeys, ...rightKeys];

	interface FallingLetter {
		id: number;
		letter: string;
		x: number;
		y: number;
		speed: number;
		color: string;
	}

	let letters: FallingLetter[] = $state([]);
	let score = $state(0);
	let lives = $state(5);
	let gameActive = $state(false);
	let nextId = $state(0);
	let spawnInterval: ReturnType<typeof setInterval> | null = null;
	let animationFrame: number | null = null;
	let lastTime = 0;
	let currentDifficulty = $state(1);

	function getSpawnInterval(): number {
		const baseInterval = 1200;
		const reduction = Math.floor(score / 100) * 100;
		return Math.max(400, baseInterval - reduction);
	}

	function getSpeedMultiplier(): number {
		return 1 + Math.floor(score / 100) * 0.1;
	}

	const colors = ['#fbbf24', '#4ade80', '#3b82f6', '#ec4899', '#8b5cf6'];

	function getAvailableKeys(): string[] {
		if (mode === 'left') return leftKeys;
		if (mode === 'right') return rightKeys;
		return allKeys;
	}

	function spawnLetter() {
		const keys = getAvailableKeys();
		const letter = keys[Math.floor(Math.random() * keys.length)];
		const baseSpeed = 0.3 + Math.random() * 0.4;
		const newLetter: FallingLetter = {
			id: nextId++,
			letter,
			x: Math.random() * 80 + 10,
			y: 105,
			speed: baseSpeed * getSpeedMultiplier(),
			color: colors[Math.floor(Math.random() * colors.length)]
		};
		letters = [...letters, newLetter];
	}

	function startGame() {
		score = 0;
		lives = 5;
		letters = [];
		gameActive = true;
		nextId = 0;
		currentDifficulty = 1;

		spawnInterval = setInterval(() => {
			if (letters.length < 6) {
				spawnLetter();
			}
		}, getSpawnInterval());

		lastTime = performance.now();
		animationFrame = requestAnimationFrame(update);
	}

	function update(currentTime: number) {
		if (!gameActive) return;

		const deltaTime = currentTime - lastTime;
		lastTime = currentTime;

		const newDifficulty = Math.floor(score / 100) + 1;
		if (newDifficulty !== currentDifficulty) {
			currentDifficulty = newDifficulty;
			if (spawnInterval) clearInterval(spawnInterval);
			spawnInterval = setInterval(() => {
				if (letters.length < 6) {
					spawnLetter();
				}
			}, getSpawnInterval());
		}

		letters = letters
			.map(l => ({ ...l, y: l.y - l.speed * (deltaTime / 16) }))
			.filter(l => {
				if (l.y < -5) {
					lives--;
					if (lives <= 0) {
						endGame();
					}
					return false;
				}
				return true;
			});

		animationFrame = requestAnimationFrame(update);
	}

	function endGame() {
		gameActive = false;
		if (spawnInterval) clearInterval(spawnInterval);
		if (animationFrame) cancelAnimationFrame(animationFrame);
		letters = [];
	}

	function handleKeydown(e: KeyboardEvent) {
		if (!gameActive) return;

		const key = e.key.toLowerCase();
		const keys = getAvailableKeys();

		if (!keys.includes(key)) return;

		const index = letters.findIndex(l => l.letter === key);
		if (index !== -1) {
			letters = letters.filter((_, i) => i !== index);
			score += 10;
		}
	}

	function handleRestart() {
		startGame();
	}

	$effect(() => {
		return () => {
			if (spawnInterval) clearInterval(spawnInterval);
			if (animationFrame) cancelAnimationFrame(animationFrame);
		};
	});
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="game-container">
	{#if !gameActive}
		<div class="start-screen">
			<h2 class="game-title">🎮 Tangkap Huruf!</h2>
			<p class="game-desc">Ketik huruf yang melayang sebelum menyentuh atas!</p>
			<div class="mode-info">
				<p>Mode: <strong>{mode === 'left' ? 'Home Row Kiri (ASDF)' : mode === 'right' ? 'Home Row Kanan (JKL;)' : 'Home Row Lengkap'}</strong></p>
			</div>
			<button class="start-btn" onclick={startGame}>Mulai Game</button>
		</div>
	{:else}
		<div class="game-hud">
			<div class="hud-item">
				<span class="hud-label">Skor</span>
				<span class="hud-value">{score}</span>
			</div>
			<div class="hud-item">
				<span class="hud-label">Level</span>
				<span class="hud-value">{currentDifficulty}</span>
			</div>
			<div class="hud-item">
				<span class="hud-label">Nyawa</span>
				<span class="hud-value">{'❤️'.repeat(lives)}</span>
			</div>
		</div>

		<div class="game-area">
			{#each letters as letter (letter.id)}
				<div
					class="falling-letter"
					style="left: {letter.x}%; top: {letter.y}%; color: {letter.color}"
				>
					{letter.letter.toUpperCase()}
				</div>
			{/each}
			<div class="ground"></div>
		</div>

		{#if lives <= 0}
			<div class="game-over">
				<h3>Game Over!</h3>
				<p>Skor Akhir: <strong>{score}</strong></p>
				<button class="restart-btn" onclick={handleRestart}>Main Lagi</button>
			</div>
		{/if}
	{/if}
</div>

<style>
	.game-container {
		width: 100%;
		max-width: 900px;
		margin: 0 auto;
		background: linear-gradient(135deg, #ffffff, #f8fafc);
		border-radius: 24px;
		padding: 32px;
		border: 2px solid rgba(251, 146, 60, 0.3);
		min-height: 500px;
		position: relative;
		overflow: hidden;
	}

	.start-screen {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 400px;
		text-align: center;
	}

	.game-title {
		font-size: 48px;
		font-weight: 800;
		margin: 0 0 16px 0;
		background: linear-gradient(135deg, #fbbf24, #f59e0b);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
	}

	.game-desc {
		font-size: 18px;
		color: #9a3412;
		margin: 0 0 24px 0;
	}

	.mode-info {
		background: rgba(255, 255, 255, 0.6);
		padding: 16px 24px;
		border-radius: 12px;
		margin-bottom: 32px;
	}

	.mode-info p {
		margin: 0;
		color: #475569;
	}

	.mode-info strong {
		color: #fbbf24;
	}

	.start-btn {
		padding: 16px 48px;
		font-size: 20px;
		font-weight: 700;
		background: linear-gradient(135deg, #fbbf24, #f59e0b);
		color: #431407;
		border: none;
		border-radius: 12px;
		cursor: pointer;
		transition: all 0.2s;
	}

	.start-btn:hover {
		transform: scale(1.05);
		box-shadow: 0 8px 24px rgba(251, 191, 36, 0.4);
	}

	.game-hud {
		display: flex;
		justify-content: space-between;
		margin-bottom: 24px;
		padding: 16px;
		background: rgba(0, 0, 0, 0.05);
		border-radius: 12px;
	}

	.hud-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 4px;
	}

	.hud-label {
		font-size: 14px;
		color: #9a3412;
		text-transform: uppercase;
		letter-spacing: 1px;
	}

	.hud-value {
		font-size: 32px;
		font-weight: 800;
		color: #fbbf24;
	}

	.game-area {
		position: relative;
		width: 100%;
		height: 400px;
		background: rgba(0, 0, 0, 0.2);
		border-radius: 16px;
		overflow: hidden;
		border: 2px solid rgba(251, 146, 60, 0.3);
	}

	.falling-letter {
		position: absolute;
		font-size: 64px;
		font-weight: 900;
		text-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
		transition: transform 0.1s;
		user-select: none;
		pointer-events: none;
	}

	.ground {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		height: 4px;
		background: linear-gradient(90deg, #ef4444, #f97316, #eab308);
	}

	.game-over {
		position: absolute;
		inset: 0;
		background: rgba(0, 0, 0, 0.9);
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		border-radius: 24px;
		text-align: center;
	}

	.game-over h3 {
		font-size: 48px;
		font-weight: 800;
		margin: 0 0 16px 0;
		color: #ef4444;
	}

	.game-over p {
		font-size: 24px;
		color: #431407;
		margin: 0 0 32px 0;
	}

	.game-over strong {
		color: #fbbf24;
		font-size: 32px;
	}

	.restart-btn {
		padding: 12px 32px;
		font-size: 18px;
		font-weight: 700;
		background: linear-gradient(135deg, #4ade80, #22c55e);
		color: #431407;
		border: none;
		border-radius: 12px;
		cursor: pointer;
		transition: all 0.2s;
	}

	.restart-btn:hover {
		transform: scale(1.05);
		box-shadow: 0 8px 24px rgba(74, 222, 128, 0.4);
	}
</style>
