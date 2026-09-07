<script lang="ts">
	import { KEY_MAP, FINGER_COLORS } from '$lib/data/lessons';

	let { activeKey = '', pressedKey = '', highlightKeys = [] }: {
		activeKey?: string;
		pressedKey?: string;
		highlightKeys?: string[];
	} = $props();

	const rows = [
		[
			{ key: '`', width: 1 },
			{ key: '1', width: 1 },
			{ key: '2', width: 1 },
			{ key: '3', width: 1 },
			{ key: '4', width: 1 },
			{ key: '5', width: 1 },
			{ key: '6', width: 1 },
			{ key: '7', width: 1 },
			{ key: '8', width: 1 },
			{ key: '9', width: 1 },
			{ key: '0', width: 1 },
			{ key: '-', width: 1 },
			{ key: '=', width: 1 },
			{ key: 'Backspace', width: 2, label: '⌫' }
		],
		[
			{ key: 'Tab', width: 1.5, label: 'Tab' },
			{ key: 'q', width: 1 },
			{ key: 'w', width: 1 },
			{ key: 'e', width: 1 },
			{ key: 'r', width: 1 },
			{ key: 't', width: 1 },
			{ key: 'y', width: 1 },
			{ key: 'u', width: 1 },
			{ key: 'i', width: 1 },
			{ key: 'o', width: 1 },
			{ key: 'p', width: 1 },
			{ key: '[', width: 1 },
			{ key: ']', width: 1 },
			{ key: '\\', width: 1.5 }
		],
		[
			{ key: 'CapsLock', width: 1.75, label: 'Caps' },
			{ key: 'a', width: 1 },
			{ key: 's', width: 1 },
			{ key: 'd', width: 1 },
			{ key: 'f', width: 1 },
			{ key: 'g', width: 1 },
			{ key: 'h', width: 1 },
			{ key: 'j', width: 1 },
			{ key: 'k', width: 1 },
			{ key: 'l', width: 1 },
			{ key: ';', width: 1 },
			{ key: "'", width: 1 },
			{ key: 'Enter', width: 2.25, label: 'Enter' }
		],
		[
			{ key: 'ShiftLeft', width: 2.25, label: 'Shift' },
			{ key: 'z', width: 1 },
			{ key: 'x', width: 1 },
			{ key: 'c', width: 1 },
			{ key: 'v', width: 1 },
			{ key: 'b', width: 1 },
			{ key: 'n', width: 1 },
			{ key: 'm', width: 1 },
			{ key: ',', width: 1 },
			{ key: '.', width: 1 },
			{ key: '/', width: 1 },
			{ key: 'ShiftRight', width: 2.75, label: 'Shift' }
		],
		[
			{ key: 'CtrlLeft', width: 1.25, label: 'Ctrl' },
			{ key: 'WinLeft', width: 1.25, label: '⊞' },
			{ key: 'AltLeft', width: 1.25, label: 'Alt' },
			{ key: ' ', width: 6.25, label: 'Space' },
			{ key: 'AltRight', width: 1.25, label: 'Alt' },
			{ key: 'WinRight', width: 1.25, label: '⊞' },
			{ key: 'Menu', width: 1.25, label: '☰' },
			{ key: 'CtrlRight', width: 1.25, label: 'Ctrl' }
		]
	];

	function getKeyColor(key: string): string {
		const info = KEY_MAP[key];
		if (!info) return '#374151';
		return FINGER_COLORS[info.finger];
	}

	function isActive(key: string): boolean {
		return activeKey.toLowerCase() === key;
	}

	function isPressed(key: string): boolean {
		return pressedKey.toLowerCase() === key;
	}

	function isHighlighted(key: string): boolean {
		return highlightKeys.includes(key);
	}

	function getKeyLabel(item: { key: string; label?: string }): string {
		if (item.label) return item.label;
		if (item.key === ';') return ';';
		if (item.key === ',') return ',';
		if (item.key === '.') return '.';
		if (item.key === '/') return '/';
		if (item.key === '`') return '`';
		if (item.key === "'") return "'";
		if (item.key === '[') return '[';
		if (item.key === ']') return ']';
		if (item.key === '\\') return '\\';
		if (item.key === '-') return '-';
		if (item.key === '=') return '=';
		return item.key.toUpperCase();
	}

	function isSpecialKey(key: string): boolean {
		return ['Backspace', 'Tab', 'CapsLock', 'Enter', 'ShiftLeft', 'ShiftRight', 'CtrlLeft', 'CtrlRight', 'AltLeft', 'AltRight', 'WinLeft', 'WinRight', 'Menu'].includes(key);
	}
</script>

<div class="keyboard">
	{#each rows as row}
		<div class="keyboard-row">
			{#each row as item}
				{@const color = getKeyColor(item.key)}
				{@const active = isActive(item.key)}
				{@const pressed = isPressed(item.key)}
				{@const highlighted = isHighlighted(item.key)}
				{@const special = isSpecialKey(item.key)}
				<button
					class="key"
					class:active
					class:pressed
					class:highlighted
					class:special
					style="--key-color: {color}; --key-color-light: {color}33; --key-width: {item.width}"
					disabled
				>
					<span class="key-label">{getKeyLabel(item)}</span>
					{#if highlighted && !special}
						<span class="key-dot"></span>
					{/if}
				</button>
			{/each}
		</div>
	{/each}
</div>

<style>
	.keyboard {
		display: flex;
		flex-direction: column;
		gap: 6px;
		padding: 16px;
		background: rgba(255, 255, 255, 0.6);
		border-radius: 16px;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
		backdrop-filter: blur(8px);
		border: 1px solid rgba(251, 146, 60, 0.3);
	}

	.keyboard-row {
		display: flex;
		gap: 6px;
		justify-content: flex-start;
	}

	.key {
		position: relative;
		width: calc(48px * var(--key-width, 1));
		height: 48px;
		border-radius: 10px;
		border: 2px solid var(--key-color);
		background: var(--key-color-light);
		color: #431407;
		text-shadow: 0 1px 1px rgba(0,0,0,0.1);
		font-size: 14px;
		font-weight: 700;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.1s ease;
		cursor: default;
		user-select: none;
	}

	.key.special {
		background: rgba(255, 255, 255, 0.6);
		border-color: rgba(251, 146, 60, 0.3);
		color: #9a3412;
		font-size: 12px;
	}

	.key.active {
		background: var(--key-color);
		color: #431407;
		transform: scale(1.1);
		box-shadow: 0 0 20px var(--key-color);
		animation: pulse 0.6s ease infinite;
	}

	.key.pressed {
		transform: scale(0.95);
		background: var(--key-color);
		color: #431407;
		box-shadow: 0 0 12px var(--key-color);
	}

	.key.highlighted {
		border-width: 3px;
		box-shadow: 0 0 8px var(--key-color);
	}

	.key-dot {
		position: absolute;
		bottom: 4px;
		left: 50%;
		transform: translateX(-50%);
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--key-color);
	}

	.key-label {
		pointer-events: none;
	}

	@keyframes pulse {
		0%, 100% { box-shadow: 0 0 20px var(--key-color); }
		50% { box-shadow: 0 0 30px var(--key-color); }
	}

	@media (max-width: 768px) {
		.key {
			width: calc(32px * var(--key-width, 1));
			height: 36px;
			font-size: 11px;
			border-radius: 6px;
		}

		.keyboard {
			padding: 8px;
			gap: 4px;
		}

		.keyboard-row {
			gap: 3px;
		}
	}
</style>
