/// <reference lib="webworker" />

let pyodide: any;

self.onmessage = async (event) => {
  const { id, type, code } = event.data;

  if (type === 'INIT') {
    try {
      if (!pyodide) {
        // Module workers (like Vite's default) don't support importScripts.
        // We must use dynamic import to load Pyodide's ESM version.
        const pyodideModule = await import('https://cdn.jsdelivr.net/pyodide/v0.25.0/full/pyodide.mjs');
        pyodide = await pyodideModule.loadPyodide();
        
        pyodide.setStdout({
          batched: (msg: string) => {
            self.postMessage({ id, type: 'STDOUT', msg });
          }
        });
      }
      self.postMessage({ id, type: 'INIT_DONE' });
    } catch (e: any) {
      self.postMessage({ id, type: 'ERROR', error: e.message || e.toString() });
    }
  } else if (type === 'RUN') {
    try {
      if (!pyodide) throw new Error("Pyodide not initialized");
      await pyodide.runPythonAsync(code);
      self.postMessage({ id, type: 'RUN_DONE' });
    } catch (e: any) {
      self.postMessage({ id, type: 'ERROR', error: e.message || e.toString() });
    }
  }
};
