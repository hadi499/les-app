<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import * as Blockly from 'blockly/core';
  import * as En from 'blockly/msg/en'; // Base locale
  import 'blockly/blocks';
  
  import { defineCustomBlocks, getToolbox } from './blocks/custom-blocks';

  let { allowedBlocks, onBlocksChanged } = $props<{
    allowedBlocks: string[];
    onBlocksChanged: (blocks: any[]) => void;
  }>();

  let blocklyDiv: HTMLDivElement;
  let workspace: Blockly.WorkspaceSvg;

  onMount(() => {
    Blockly.setLocale(En as any);
    defineCustomBlocks();

    workspace = Blockly.inject(blocklyDiv, {
      toolbox: getToolbox(allowedBlocks),
      scrollbars: true,
      trashcan: true,
      zoom: {
        controls: true,
        wheel: true,
        startScale: 1.0,
        maxScale: 2,
        minScale: 0.5,
        scaleSpeed: 1.2
      }
    });

    workspace.addChangeListener(() => {
      // Extract the sequence of blocks
      const topBlocks = workspace.getTopBlocks(true);
      if (topBlocks.length > 0) {
        const sequence = extractSequence(topBlocks[0]);
        onBlocksChanged(sequence);
      } else {
        onBlocksChanged([]);
      }
    });
  });

  onDestroy(() => {
    if (workspace) {
      workspace.dispose();
    }
  });

  // Helper to extract a sequence of block intents from a start block
  function extractSequence(block: Blockly.Block | null): any[] {
    const seq = [];
    let currentBlock = block;
    while (currentBlock) {
      const type = currentBlock.type;
      
      if (type === 'repeat_n') {
        const times = currentBlock.getFieldValue('TIMES');
        const doBlock = currentBlock.getInputTargetBlock('DO');
        seq.push({
          type: 'repeat_n',
          times: parseInt(times, 10),
          do: extractSequence(doBlock)
        });
      } else {
        seq.push({ type });
      }

      currentBlock = currentBlock.getNextBlock();
    }
    return seq;
  }
</script>

<div bind:this={blocklyDiv} class="absolute inset-0"></div>

<style>
  /* Fix blockly injecting directly to body occasionally */
  :global(.blocklyWidgetDiv) {
    z-index: 99999;
  }
</style>
