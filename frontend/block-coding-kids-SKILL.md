---
name: block-coding-for-kids
description: Build or modify a block-based visual coding module (using Blockly) for teaching computational thinking to children aged 5-10, in a SvelteKit app. Use this skill whenever the user asks to add a kids' coding module, a drag-and-drop block programming feature, a Scratch/ScratchJr-like activity, or a computational-thinking learning tool for young children — especially for the lesbalonggarut.my.id educational site. This is distinct from the text-editor live-coding-editor skill, which targets older/teen learners writing real HTML/CSS/JS/Python. Covers custom icon-based Blockly blocks, a canvas/sprite stage for visual feedback, a safe custom interpreter (not eval), level/curriculum progression, and unplugged-activity considerations for pre-readers.
---

# Block Coding for Kids (Ages 5-10)

A skill for building a drag-and-drop, icon-based block programming module aimed at young children, where the goal is developing computational thinking (sequencing, loops, conditionals, debugging) rather than teaching a real programming language. Built for a SvelteKit frontend.

## When to use this

- User wants to add a kids' coding activity/module to a SvelteKit site, especially for ages roughly 5-10
- User mentions lesbalonggarut.my.id and wants a coding module for young children, separate from the text-based live coding editor for older learners
- User wants a Scratch-like or ScratchJr-like drag-and-drop programming experience
- User wants to build computational thinking / logic skills through a visual/game-like interface

## Design principles for this audience

The audience cannot reliably read fluent sentences (especially ages 5-7). This drives every design decision:

1. **Icons over text.** Blocks should communicate through arrows, colors, and pictograms, not sentences. Avoid Blockly's default verbose English/Indonesian block labels where possible — design custom block visuals.
2. **Visible, animated execution.** Never jump straight to the end result. Run the program step-by-step with visible motion (character moves one step at a time) so the child can connect cause and effect — this is the pedagogical core of the exercise, not a nice-to-have.
3. **Friendly failure, not error messages.** When a level fails (character hits a wall, doesn't reach the goal), highlight the block responsible and give a warm, concrete hint — never a technical error.
4. **Small, progressive steps.** Introduce exactly one new concept per level band (sequencing → loops → conditionals). Don't mix concepts early.
5. **No login, no typing required.** Everything is drag-and-drop; persist via `localStorage`, no backend needed.
6. **Consider unplugged activities first for ages 5-7.** Physical command cards (move forward, turn) that guide a toy or a friend to a target build the same sequencing/debugging intuition before introducing a screen at all. Recommend this as a possible precursor, especially for the younger end of the range.

## Architecture

### Core components (SvelteKit + Blockly)

```
src/lib/components/BlockPlayground/
  BlocklyWorkspace.svelte   - Blockly workspace wrapper, drag-and-drop area
  BlockRunner.svelte        - executes the block program step-by-step
  CanvasStage.svelte        - canvas where the character/sprite moves, the visual feedback surface
  blocks/
    custom-blocks.ts        - custom block definitions (move, turn, repeat, if, start/end)
  levels.ts                 - level list: starting state, target, allowed blocks, validation
```

### Custom blocks (icon-based, not default Blockly text blocks)

Design blocks with `Blockly.Blocks[...]` and custom SVG icons instead of the default text-heavy blocks:
- **Move**: forward / backward / turn left / turn right — arrow icons
- **Repeat**: a loop-shaped block, "repeat N times" shown as a visual counter, not a text field the child must type into (use a stepper/dial instead)
- **If**: a simple condition block with a pictogram (e.g. "if wall ahead")
- **Start / End**: clear visual bookends for the program

### Execution model — use a custom interpreter, not `eval` / Blockly's JS generator directly

For safety and control, walk the Blockly block tree yourself (or use Blockly's generator to produce an intermediate representation, then interpret that) rather than generating and `eval`-ing raw JavaScript. This gives full control over step timing/animation and avoids any code-injection surface, which matters more here than in the teen-focused editor since there's no sandboxing expectation from the child.

Execution flow:
1. Child arranges blocks, taps "Run"
2. Interpreter walks the program one instruction at a time
3. Each instruction animates the character on `CanvasStage` (don't skip straight to end state)
4. On success: celebratory feedback tied to the visual target (e.g. reaching a star)
5. On failure: pause at the failing block, highlight it, show a simple encouraging hint

### Level / curriculum structure

Progressive, one new concept per band:
- **Levels 1-3**: movement blocks only (sequencing)
- **Levels 4-6**: introduce "repeat" (loops)
- **Levels 7-10**: introduce "if" (conditionals)

Each level: one clear visual goal + the minimal block set needed, with automatic validation against the target state.

### Progress persistence

No backend/login required:
```js
localStorage.setItem(`level-${id}-solution`, JSON.stringify(blocklyXml));
```

## What NOT to add for this use case

- No default verbose Blockly text blocks aimed at adult/teen developers
- No raw `eval()` of generated JavaScript — use a controlled interpreter
- No typed input fields for young children where a visual control (stepper, icon picker) can substitute
- No login/account system — keep it zero-friction
- Don't reuse the text-based `live-coding-editor` skill's CodeMirror/Pyodide approach here — this age group needs block-based, not text-based, programming
