export async function compressImageFile(file: File, quality = 0.7, maxWidth = 1920): Promise<File> {
  if (!file.type.startsWith("image/")) {
    return file;
  }
  
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
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
              resolve(file);
              return;
            }
            ctx.drawImage(img, 0, 0, width, height);

            canvas.toBlob(blob => {
                if (blob) {
                  // Only return compressed if it's actually smaller
                  if (blob.size >= file.size) {
                    resolve(file);
                  } else {
                    const compressedFile = new File([blob], file.name, { type: 'image/jpeg' });
                    resolve(compressedFile);
                  }
                } else {
                  resolve(file);
                }
            }, 'image/jpeg', quality);
        };
        img.onerror = () => resolve(file);
    };
    reader.onerror = () => resolve(file);
  });
}
