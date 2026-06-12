let message = $state('');
let type = $state<'success' | 'error'>('success');
let visible = $state(false);
let timeout: ReturnType<typeof setTimeout> | null = null;

export const toast = {
  get message() { return message; },
  get type() { return type; },
  get visible() { return visible; },

  show(msg: string, t: 'success' | 'error' = 'success') {
    if (timeout) clearTimeout(timeout);
    message = msg;
    type = t;
    visible = true;
    timeout = setTimeout(() => { visible = false; }, 3000);
  },

  success(msg: string) { this.show(msg, 'success'); },
  error(msg: string) { this.show(msg, 'error'); },
};
