import * as Blockly from 'blockly/core';

// Custom blocks for kids: focus on icons/arrows, minimal text

export function defineCustomBlocks() {
  Blockly.Blocks['move_forward'] = {
    init: function() {
      this.jsonInit({
        "type": "move_forward",
        "message0": "⬆️ Maju",
        "previousStatement": null,
        "nextStatement": null,
        "colour": 230,
        "tooltip": "Maju satu langkah",
        "helpUrl": ""
      });
    }
  };

  Blockly.Blocks['turn_left'] = {
    init: function() {
      this.jsonInit({
        "type": "turn_left",
        "message0": "↩️ Kiri",
        "previousStatement": null,
        "nextStatement": null,
        "colour": 20,
        "tooltip": "Belok ke kiri",
        "helpUrl": ""
      });
    }
  };

  Blockly.Blocks['turn_right'] = {
    init: function() {
      this.jsonInit({
        "type": "turn_right",
        "message0": "↪️ Kanan",
        "previousStatement": null,
        "nextStatement": null,
        "colour": 160,
        "tooltip": "Belok ke kanan",
        "helpUrl": ""
      });
    }
  };

  Blockly.Blocks['repeat_n'] = {
    init: function() {
      this.jsonInit({
        "type": "repeat_n",
        "message0": "🔁 Ulangi %1 kali",
        "args0": [
          {
            "type": "field_number",
            "name": "TIMES",
            "value": 2,
            "min": 1,
            "max": 10,
            "precision": 1
          }
        ],
        "message1": "Lakukan: %1",
        "args1": [
          {
            "type": "input_statement",
            "name": "DO"
          }
        ],
        "previousStatement": null,
        "nextStatement": null,
        "colour": 120,
        "tooltip": "Ulangi perintah",
        "helpUrl": ""
      });
    }
  };
}

export function getToolbox(allowedBlocks: string[]) {
  return {
    "kind": "flyoutToolbox",
    "contents": allowedBlocks.map(block => ({
      "kind": "block",
      "type": block
    }))
  };
}
