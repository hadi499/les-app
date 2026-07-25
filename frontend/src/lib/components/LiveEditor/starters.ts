export type Language = 'html' | 'css' | 'js' | 'python';

export interface Exercise {
  id: string;
  title: string;
  type: 'web' | 'python'; // web uses html/css/js, python uses python
  description: string;
  starter: {
    html?: string;
    css?: string;
    js?: string;
    python?: string;
  };
}

export const exercises: Record<string, Exercise> = {
  'hello-html': {
    id: 'hello-html',
    title: 'HTML, CSS, JS',
    type: 'web',
    description: 'HTML, CSS, dan JavaScript(JS) adalah tiga pilar utama pembuat halaman web: HTML bertindak sebagai struktur atau kerangka yang menentukan isi seperti teks dan gambar, CSS mengatur tampilan visual seperti warna, tata letak, dan gaya agar web terlihat menarik, sedangkan JavaScript memberikan interaktivitas yang membuat web hidup dan dapat merespons aksi pengguna.',
    starter: {
      html: ``,
      css: ``,
      js: ``
    }
  },
  'hello-python': {
    id: 'hello-python',
    title: 'Python',
    type: 'python',
    description: 'Python merupakan bahasa pemrograman multi-guna yang banyak digunakan untuk pengembangan web backend, analisis data, kecerdasan buatan, hingga otomasi skrip. Keunggulan utamanya terletak pada sintaksis yang bersih dan menyerupai bahasa Inggris sehingga sangat ramah bagi pemula.',
    starter: {
      python: ``
    }
  }
};
