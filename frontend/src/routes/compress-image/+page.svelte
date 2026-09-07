<script lang="ts">
  import { fade, slide } from "svelte/transition";

  let selectedFile = $state<File | null>(null);
  let originalPreviewUrl = $state<string | null>(null);
  let originalSize = $state<number>(0);
  
  let compressedBlob = $state<Blob | null>(null);
  let compressedPreviewUrl = $state<string | null>(null);
  let compressedSize = $state<number>(0);
  
  let isCompressing = $state(false);
  let quality = $state(0.7);
  let maxWidth = $state(1920);

  let isDragging = $state(false);

  function handleDragOver(e: DragEvent) {
    e.preventDefault();
    isDragging = true;
  }

  function handleDragLeave() {
    isDragging = false;
  }

  function handleDrop(e: DragEvent) {
    e.preventDefault();
    isDragging = false;
    
    if (e.dataTransfer && e.dataTransfer.files.length > 0) {
      processSelectedFile(e.dataTransfer.files[0]);
    }
  }

  function handleFileChange(event: Event) {
    const target = event.target as HTMLInputElement;
    if (target.files && target.files.length > 0) {
      processSelectedFile(target.files[0]);
    }
  }

  function processSelectedFile(file: File) {
      if (!file.type.startsWith("image/")) {
        alert("Pilih file gambar yang valid!");
        return;
      }
      selectedFile = file;
      originalSize = file.size;
      
      if (originalPreviewUrl) URL.revokeObjectURL(originalPreviewUrl);
      originalPreviewUrl = URL.createObjectURL(file);
      
      // Reset compressed
      compressedBlob = null;
      if (compressedPreviewUrl) URL.revokeObjectURL(compressedPreviewUrl);
      compressedPreviewUrl = null;
      compressedSize = 0;
  }

  async function compressImage() {
    if (!selectedFile) return;
    
    isCompressing = true;
    
    try {
      const result = await new Promise<{blob: Blob, url: string, size: number}>((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(selectedFile!);
        reader.onload = event => {
            const img = new Image();
            img.src = event.target!.result as string;
            img.onload = () => {
                let width = img.width;
                let height = img.height;

                if (width > maxWidth) {
                    height = Math.round((height * maxWidth) / width);
                    width = maxWidth;
                }

                const canvas = document.createElement('canvas');
                canvas.width = width;
                canvas.height = height;

                const ctx = canvas.getContext('2d');
                if (!ctx) {
                  reject("Could not get 2d context");
                  return;
                }
                ctx.drawImage(img, 0, 0, width, height);

                canvas.toBlob(blob => {
                    if (blob) {
                      resolve({
                          blob: blob,
                          url: URL.createObjectURL(blob),
                          size: blob.size
                      });
                    } else {
                      reject("Compression failed");
                    }
                }, 'image/jpeg', quality);
            };
            img.onerror = error => reject(error);
        };
        reader.onerror = error => reject(error);
      });
      
      compressedBlob = result.blob;
      compressedPreviewUrl = result.url;
      compressedSize = result.size;
    } catch (e) {
      console.error(e);
      alert("Gagal mengkompres gambar.");
    } finally {
      isCompressing = false;
    }
  }

  function downloadImage() {
    if (!compressedBlob || !compressedPreviewUrl) return;
    
    const a = document.createElement("a");
    a.href = compressedPreviewUrl;
    
    const originalName = selectedFile?.name || "image.jpg";
    const nameParts = originalName.split('.');
    nameParts.pop(); // remove extension
    const baseName = nameParts.join('.');
    
    a.download = `${baseName}_compressed.jpg`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
  }

  function formatBytes(bytes: number, decimals = 2) {
      if (!+bytes) return '0 Bytes';
      const k = 1024;
      const dm = decimals < 0 ? 0 : decimals;
      const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
  }
  
  $effect(() => {
    return () => {
      if (originalPreviewUrl) URL.revokeObjectURL(originalPreviewUrl);
      if (compressedPreviewUrl) URL.revokeObjectURL(compressedPreviewUrl);
    };
  });
</script>

<svelte:head>
  <title>Compress Image - Les App</title>
  <meta name="description" content="Kecilkan ukuran gambar secara gratis langsung di browser. 100% aman dan privat — gambar tidak pernah diunggah ke server manapun." />
  <link rel="canonical" href="https://lesbalonggarut.my.id/compress-image" />

  <!-- Open Graph -->
  <meta property="og:type" content="website" />
  <meta property="og:url" content="https://lesbalonggarut.my.id/compress-image" />
  <meta property="og:title" content="Compress Image Online - Gratis & Privat" />
  <meta property="og:description" content="Kecilkan ukuran gambar secara gratis, cepat, dan privat langsung di browser tanpa upload ke server." />
  <meta property="og:site_name" content="Les Balongarut" />
  <meta property="og:locale" content="id_ID" />

  <!-- Twitter Card -->
  <meta name="twitter:card" content="summary" />
  <meta name="twitter:title" content="Compress Image Online - Gratis & Privat" />
  <meta name="twitter:description" content="Kecilkan ukuran gambar secara gratis dan privat langsung di browser." />
</svelte:head>

<div class="min-h-screen bg-slate-50 pt-24 pb-12">
  <div class="max-w-5xl mx-auto p-4 md:p-6 lg:p-8 animate-in fade-in duration-500">
  <div class="mb-8">
    <h1 class="text-2xl md:text-3xl font-bold text-slate-900 tracking-tight">Compress Image</h1>
    <p class="text-slate-500 mt-2 text-base leading-relaxed">Kecilkan ukuran gambar Anda dengan cepat dan mudah langsung dari browser. <strong>100% Aman & Privat</strong> — gambar Anda tidak pernah diunggah ke server manapun, sehingga privasi Anda sepenuhnya terjaga.</p>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 lg:gap-8">
    
    <!-- Controls Section -->
    <div class="lg:col-span-1 space-y-6">
      <!-- File Upload Card -->
      <div class="bg-white rounded-2xl shadow-sm border border-slate-200 p-5 md:p-6 transition-shadow hover:shadow-md">
        <h3 class="block text-sm font-semibold text-slate-700 mb-3">Pilih Gambar</h3>
        
        <label 
          for="image-upload" 
          ondragover={handleDragOver}
          ondragleave={handleDragLeave}
          ondrop={handleDrop}
          class="flex flex-col items-center justify-center w-full h-32 md:h-40 border-2 border-dashed rounded-xl cursor-pointer transition-all group {isDragging ? 'border-blue-500 bg-blue-50 scale-[1.02]' : 'border-slate-300 bg-slate-50 hover:bg-slate-100 hover:border-blue-400'}"
        >
            <div class="flex flex-col items-center justify-center pt-5 pb-6">
                <svg class="w-8 h-8 md:w-10 md:h-10 mb-3 text-slate-400 group-hover:text-blue-500 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"></path></svg>
                <p class="mb-2 text-sm text-slate-500 text-center px-4"><span class="font-semibold text-blue-600">Klik untuk memilih</span> atau drag dan drop</p>
                <p class="text-xs text-slate-400">PNG, JPG, JPEG, WEBP</p>
            </div>
            <input id="image-upload" type="file" accept="image/*" class="hidden" onchange={handleFileChange} />
        </label>
      </div>

      <!-- Settings Card -->
      {#if selectedFile}
        <div in:slide class="bg-white rounded-2xl shadow-sm border border-slate-200 p-5 md:p-6 transition-shadow hover:shadow-md">
          <h3 class="text-base font-semibold text-slate-800 mb-4 flex items-center gap-2">
            <svg class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"></path></svg>
            Pengaturan Kompresi
          </h3>
          
          <div class="space-y-5">
            <div>
              <div class="flex justify-between mb-1">
                <label for="quality" class="text-sm font-medium text-slate-700">Kualitas Gambar</label>
                <span class="text-sm text-blue-600 font-semibold">{Math.round(quality * 100)}%</span>
              </div>
              <input type="range" id="quality" min="0.1" max="1" step="0.05" bind:value={quality} class="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-blue-600">
              <div class="flex justify-between text-xs text-slate-400 mt-1">
                <span>Ukuran Kecil</span>
                <span>Kualitas Tinggi</span>
              </div>
            </div>

            <div>
              <label for="maxWidth" class="block text-sm font-medium text-slate-700 mb-1">Lebar Maksimal (px)</label>
              <input type="number" id="maxWidth" bind:value={maxWidth} class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
            </div>

            <button 
              onclick={compressImage}
              disabled={isCompressing}
              class="w-full py-3 px-4 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-semibold rounded-xl shadow-sm shadow-blue-600/20 transition-all flex justify-center items-center gap-2"
            >
              {#if isCompressing}
                <svg class="animate-spin -ml-1 mr-2 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Memproses...
              {:else}
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
                Mulai Kompres
              {/if}
            </button>
          </div>
        </div>
      {/if}
    </div>

    <!-- Preview Section -->
    <div class="lg:col-span-2 space-y-6">
      {#if !selectedFile}
        <div class="h-full min-h-[300px] border-2 border-dashed border-slate-200 rounded-2xl flex flex-col items-center justify-center text-slate-400 bg-white/50 p-8 text-center">
          <svg class="w-16 h-16 mb-4 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
          <p class="text-lg font-medium text-slate-500">Belum ada gambar yang dipilih</p>
          <p class="text-sm mt-1">Pilih gambar dari panel di sebelah kiri untuk melihat preview.</p>
        </div>
      {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6" in:fade>
          
          <!-- Original Preview -->
          <div class="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden flex flex-col">
            <div class="p-4 border-b border-slate-100 bg-slate-50/50 flex justify-between items-center">
              <h3 class="font-semibold text-slate-800">Original</h3>
              <span class="px-2.5 py-1 bg-slate-200 text-slate-700 text-xs font-bold rounded-lg">{formatBytes(originalSize)}</span>
            </div>
            <div class="p-4 bg-slate-100/50 flex-1 flex items-center justify-center relative min-h-[200px]">
              <div class="absolute inset-0 pattern-checkered opacity-5"></div>
              {#if originalPreviewUrl}
                <img src={originalPreviewUrl} alt="Original" class="max-w-full max-h-[300px] object-contain relative z-10 rounded shadow-sm border border-slate-200/50 bg-white" />
              {/if}
            </div>
          </div>

          <!-- Compressed Preview -->
          <div class="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden flex flex-col">
            <div class="p-4 border-b border-slate-100 bg-slate-50/50 flex justify-between items-center">
              <h3 class="font-semibold text-slate-800">Hasil Kompresi</h3>
              {#if compressedSize}
                <div class="flex items-center gap-2">
                  <span class="text-xs font-bold text-emerald-600 bg-emerald-100 px-2 py-0.5 rounded-md">
                    -{Math.round((1 - compressedSize / originalSize) * 100)}%
                  </span>
                  <span class="px-2.5 py-1 bg-blue-100 text-blue-700 text-xs font-bold rounded-lg">{formatBytes(compressedSize)}</span>
                </div>
              {:else}
                <span class="text-xs text-slate-400 font-medium">Menunggu kompresi...</span>
              {/if}
            </div>
            <div class="p-4 bg-slate-100/50 flex-1 flex items-center justify-center relative min-h-[200px]">
              <div class="absolute inset-0 pattern-checkered opacity-5"></div>
              {#if compressedPreviewUrl}
                <img src={compressedPreviewUrl} alt="Compressed" class="max-w-full max-h-[300px] object-contain relative z-10 rounded shadow-sm border border-slate-200/50 bg-white" />
              {:else if isCompressing}
                 <div class="flex flex-col items-center text-slate-400">
                    <svg class="animate-spin mb-3 h-8 w-8 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span class="text-sm font-medium">Memproses gambar...</span>
                 </div>
              {:else}
                <div class="text-center text-slate-400">
                  <svg class="w-12 h-12 mx-auto mb-2 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
                  <span class="text-sm font-medium block">Hasil akan muncul di sini</span>
                </div>
              {/if}
            </div>
            
            {#if compressedPreviewUrl}
              <div class="p-4 bg-white border-t border-slate-100" in:slide>
                <button 
                  onclick={downloadImage}
                  class="w-full py-2.5 px-4 bg-emerald-500 hover:bg-emerald-600 text-white font-semibold rounded-xl shadow-sm shadow-emerald-500/20 transition-all flex justify-center items-center gap-2"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path></svg>
                  Download Gambar
                </button>
              </div>
            {/if}
          </div>
          
        </div>
      {/if}
    </div>
  </div>
</div>
</div>

<style>
  .pattern-checkered {
    background-image: linear-gradient(45deg, #ccc 25%, transparent 25%), linear-gradient(-45deg, #ccc 25%, transparent 25%), linear-gradient(45deg, transparent 75%, #ccc 75%), linear-gradient(-45deg, transparent 75%, #ccc 75%);
    background-size: 20px 20px;
    background-position: 0 0, 0 10px, 10px -10px, -10px 0px;
  }
</style>
