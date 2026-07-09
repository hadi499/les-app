<script lang="ts">
  import { onMount } from "svelte";
  import mathCard2 from "$lib/assets/card2.png";
  import mathCard3 from "$lib/assets/card3.png";
  import mathCard4 from "$lib/assets/card4.png";

  let isAuthenticated = $state(false);
  let authChecked = $state(false);

  onMount(async () => {
    try {
      const res = await fetch(`/me`, {
        credentials: "include",
      });
      if (res.ok) {
        const data = (await res.json()) as { authenticated: boolean };
        if (data.authenticated) {
          isAuthenticated = true;
        }
      }
    } catch (e) {
      // Ignore error if no session
    } finally {
      authChecked = true;
    }
  });
</script>

<svelte:head>
  <title>Panduan & Metode Belajar - Les Balongarut</title>
</svelte:head>

<div class="min-h-screen bg-slate-50 pt-24 pb-16 font-sans">
  <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 space-y-16">
    <!-- Header Section -->
    <section
      class="text-center space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-700"
    >
      <h1
        class="text-3xl md:text-5xl font-extrabold text-slate-900 tracking-tight"
      >
        Panduan Les <span class="text-blue-600">Balongarut</span>
      </h1>
      <p class="text-lg text-slate-600 max-w-2xl mx-auto leading-relaxed">
        Tempat belajar modern dengan metode <span
          class="font-semibold text-blue-600">Card Memory</span
        > dan latihan soal terstruktur untuk hasil yang optimal.
      </p>
    </section>

    <!-- Aturan Dasar Section -->
    <section
      class="bg-blue-50/50 border border-blue-100 rounded-3xl p-8 shadow-sm animate-in fade-in slide-in-from-bottom-6 duration-700 delay-75"
    >
      <div class="space-y-6 max-w-4xl mx-auto">
        <h2 class="text-2xl md:text-3xl font-bold text-slate-800 text-center">
          Aturan Dasar
        </h2>
        <div class="grid gap-4 md:grid-cols-2">
          <div class="bg-white rounded-2xl p-6 shadow-sm border border-slate-100 flex flex-col gap-2">
            <div class="flex items-center gap-4">
              <div
                class="shrink-0 w-10 h-10 rounded-full bg-red-50 text-red-600 flex items-center justify-center font-bold text-md"
              >
                1
              </div>
              <h3 class="font-bold text-slate-800 text-md">
                Syarat Kelas Komputer
              </h3>
            </div>
            <p class="text-slate-600 text-md leading-relaxed pl-14 text-left">
              Penguasaan perkalian dasar adalah syarat wajib sebelum siswa
              dapat mengikuti pelajaran komputer.
            </p>
          </div>

          <div class="bg-white rounded-2xl p-6 shadow-sm border border-slate-100 flex flex-col gap-2">
            <div class="flex items-center gap-4">
              <div
                class="shrink-0 w-10 h-10 rounded-full bg-blue-50 text-blue-600 flex items-center justify-center font-bold text-md"
              >
                2
              </div>
              <h3 class="font-bold text-slate-800 text-md">
                Kebijakan Mengerjakan PR
              </h3>
            </div>
            <p class="text-slate-600 text-md leading-relaxed pl-14 text-left">
              Fokus utama bimbingan adalah materi les. PR sekolah dapat
              dibahas setelah sesi utama selesai, di mana guru hanya akan
              memberikan panduan penyelesaiannya.
            </p>
          </div>
        </div>
      </div>
    </section>
    <!-- Metode Card Memory Section -->
    <section
      class="bg-white rounded-3xl p-8 md:p-12 shadow-xl shadow-blue-900/5 border border-slate-100 animate-in fade-in slide-in-from-bottom-8 duration-700 delay-150"
    >
      <div class="grid md:grid-cols-2 gap-12 items-center">
        <div class="space-y-6">
          <div
            class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-50 text-blue-600 font-medium text-sm"
          >
            <svg
              class="w-4 h-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 10V3L4 14h7v7l9-11h-7z"
              /></svg
            >
            Metode Unggulan
          </div>
          <h2
            class="text-2xl md:text-3xl font-bold text-slate-900 text-center sm:text-left"
          >
            Metode Card Memory
          </h2>
          <p
            class="text-slate-600 leading-relaxed text-lg text-center sm:text-left"
          >
            Mata pelajaran dipecah menjadi bagian kecil-kecil sehingga tidak
            membebani anak-anak.
          </p>
          <ul class="space-y-4">
            <li class="flex flex-col gap-1">
              <div class="flex items-center gap-4">
                <div
                  class="w-8 h-8 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center font-bold shrink-0"
                >
                  1
                </div>
                <h4 class="font-bold text-slate-900 text-lg">
                  Baca & Hafalkan
                </h4>
              </div>
              <p class="text-slate-700 pl-12">
                Baca card memory, hafalkan selama 5 menit per card lalu tulis.
              </p>
            </li>
            <li class="flex flex-col gap-2">
              <div class="flex items-center gap-4">
                <div
                  class="w-8 h-8 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center font-bold shrink-0"
                >
                  2
                </div>
                <h4 class="font-bold text-slate-900 text-lg">
                  Latihan Teratur
                </h4>
              </div>
              <p class="text-slate-700 pl-12">
                Soal latihan 2 atau 3 kali perminggu untuk memastikan anak bisa
                menerapkan apa yang dihafal.
              </p>
            </li>
          </ul>
        </div>

        <div
          class="bg-slate-50 rounded-2xl p-6 border border-slate-100 relative overflow-hidden group"
        >
          <div
            class="absolute top-0 right-0 -mt-4 -mr-4 w-24 h-24 bg-blue-100 rounded-full blur-2xl opacity-50"
          ></div>
          <h3
            class="text-md md:text-xl font-semibold text-slate-900 mb-4 text-center sm:text-left"
          >
            Mengapa Harus Card Memory?
          </h3>
          <div class="space-y-4">
            <div
              class="bg-white p-4 rounded-xl shadow-sm border border-slate-100 group-hover:-translate-y-1 transition-transform duration-300"
            >
              <p class="text-slate-700">
                Anak zaman sekarang <span class="font-semibold text-red-500"
                  >tidak suka membaca</span
                >
                . Baca buku sedikit sudah ngantuk.
                <span
                  class="text-2xl ml-1 align-middle inline-block animate-pulse"
                  >😴</span
                >
              </p>
            </div>
            <div
              class="bg-white p-4 rounded-xl shadow-sm border border-slate-100 group-hover:-translate-y-1 transition-transform duration-300 delay-75"
            >
              <p class="text-slate-700">
                Walaupun membaca materi, kita <span
                  class="font-semibold text-amber-500">tidak tahu pasti</span
                > apakah materi tersebut benar-benar diserap dan dipahami secara
                optimal atau tidak.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Contoh Card Memory Section -->
    <section
      class="space-y-8 animate-in fade-in slide-in-from-bottom-8 duration-700 delay-300"
    >
      <div class="text-center">
        <h2 class="text-2xl md:text-3xl font-bold text-slate-900">
          Contoh Card Memory
        </h2>
        <p class="text-slate-600 mt-2">
          Visual yang menarik membantu otak mengingat lebih cepat.
        </p>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-3 gap-6">
        <div
          class="bg-transparent overflow-visible group transition-all duration-300 hover:-translate-y-2"
        >
          <img
            src={mathCard2}
            alt="Card Memory"
            class="w-[70%] sm:w-full mx-auto block h-auto object-contain group-hover:scale-105 transition-transform duration-500"
          />
        </div>
        <div
          class="bg-transparent overflow-visible group transition-all duration-300 hover:-translate-y-2 delay-75"
        >
          <img
            src={mathCard3}
            alt="Card Memory"
            class="w-[70%] sm:w-full mx-auto block h-auto object-contain group-hover:scale-105 transition-transform duration-500"
          />
        </div>
        <div
          class="bg-transparent overflow-visible group transition-all duration-300 hover:-translate-y-2 delay-150"
        >
          <img
            src={mathCard4}
            alt="Card Memory"
            class="w-[70%] sm:w-full mx-auto block h-auto object-contain group-hover:scale-105 transition-transform duration-500"
          />
        </div>
      </div>
    </section>

    <!-- Panduan Aplikasi Section -->
    <section
      class="grid md:grid-cols-2 gap-8 animate-in fade-in slide-in-from-bottom-8 duration-700 delay-500"
    >
      <div
        class="bg-white rounded-3xl p-8 shadow-xl shadow-blue-900/5 border border-slate-100 flex flex-col"
      >
        <h2
          class="text-xl md:text-2xl font-bold text-slate-900 mb-6 flex items-center gap-3"
        >
          <div class="p-2 bg-blue-100 text-blue-600 rounded-lg">
            <svg
              class="w-6 h-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"
              /></svg
            >
          </div>
          Panduan Login
        </h2>
        <ul class="space-y-4 mb-8">
          <li class="flex gap-3 items-start">
            <div class="mt-1 w-2 h-2 rounded-full bg-blue-500 shrink-0"></div>
            <p class="text-slate-700">
              Klik menu <strong class="text-slate-900">Masuk Portal</strong> di navigasi
              atas.
            </p>
          </li>
          <li class="flex gap-3 items-start">
            <div class="mt-1 w-2 h-2 rounded-full bg-blue-500 shrink-0"></div>
            <p class="text-slate-700">
              Masukkan <strong class="text-slate-900">username</strong> dan
              <strong class="text-slate-900">password</strong> yang telah diberikan
              oleh admin.
            </p>
          </li>
        </ul>

        <hr class="my-auto border-slate-100" />

        <h2
          class="text-xl md:text-2xl font-bold text-slate-900 mt-8 mb-6 flex items-center gap-3"
        >
          <div class="p-2 bg-amber-100 text-amber-600 rounded-lg">
            <svg
              class="w-6 h-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
              /></svg
            >
          </div>
          Panduan Ganti Password
        </h2>
        <ul class="space-y-4">
          <li class="flex gap-3 items-start">
            <div class="mt-1 w-2 h-2 rounded-full bg-amber-500 shrink-0"></div>
            <p class="text-slate-700">
              Masuk ke <strong class="text-slate-900">Dashboard</strong>.
            </p>
          </li>
          <li class="flex gap-3 items-start">
            <div class="mt-1 w-2 h-2 rounded-full bg-amber-500 shrink-0"></div>
            <p class="text-slate-700 leading-relaxed">
              Klik icon <span
                class="inline-flex items-center justify-center p-1 bg-slate-100 rounded border border-slate-200 mx-1 align-middle"
                ><svg
                  class="w-4 h-4 text-slate-700"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  ><path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 6h16M4 12h16M4 18h16"
                  /></svg
                ></span
              > di pojok kanan atas untuk membuka menu.
            </p>
          </li>
          <li class="flex gap-3 items-start">
            <div class="mt-1 w-2 h-2 rounded-full bg-amber-500 shrink-0"></div>
            <p class="text-slate-700">
              Klik <strong class="text-slate-900">tanda panah</strong> di pojok bawah
              kiri menu profil.
            </p>
          </li>
          <li class="flex gap-3 items-start">
            <div class="mt-1 w-2 h-2 rounded-full bg-amber-500 shrink-0"></div>
            <p class="text-slate-700">
              Pilih opsi <strong class="text-slate-900">Ganti Password</strong>.
            </p>
          </li>
        </ul>
      </div>

      <div
        class="bg-white rounded-3xl p-8 shadow-xl shadow-blue-900/5 border border-slate-100 flex flex-col"
      >
        <h2
          class="text-lg md:text-2xl font-bold text-slate-900 mb-8 flex items-center gap-3"
        >
          <div class="p-2 bg-blue-100 text-blue-600 rounded-lg">
            <svg
              class="w-6 h-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"
              /></svg
            >
          </div>
          Keterangan Menu Dashboard
        </h2>

        <div class="space-y-6">
          <div class="flex gap-4 items-start group">
            <div
              class="p-3 bg-slate-50 rounded-xl group-hover:bg-blue-50 transition-colors border border-slate-100"
            >
              <svg
                class="w-6 h-6 text-blue-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"
                /></svg
              >
            </div>
            <div>
              <h4 class="font-bold text-lg text-slate-900 mb-1">Todolist</h4>
              <p class="text-slate-600 text-sm leading-relaxed">
                Digunakan untuk melihat apa saja yang sudah dihafal oleh siswa.
              </p>
            </div>
          </div>

          <div class="flex gap-4 items-start group">
            <div
              class="p-3 bg-slate-50 rounded-xl group-hover:bg-emerald-50 transition-colors border border-slate-100"
            >
              <svg
                class="w-6 h-6 text-emerald-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"
                /></svg
              >
            </div>
            <div>
              <h4 class="font-bold text-lg text-slate-900 mb-1">
                Nilai Harian
              </h4>
              <p class="text-slate-600 text-sm leading-relaxed">
                Digunakan untuk melihat hasil nilai ujian dari soal-soal
                latihan.
              </p>
            </div>
          </div>

          <div class="flex gap-4 items-start group">
            <div
              class="p-3 bg-slate-50 rounded-xl group-hover:bg-purple-50 transition-colors border border-slate-100"
            >
              <svg
                class="w-6 h-6 text-purple-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                /></svg
              >
            </div>
            <div>
              <h4 class="font-bold text-lg text-slate-900 mb-1">
                Kuis dan Nilai
              </h4>
              <p class="text-slate-600 text-sm leading-relaxed">
                Digunakan untuk melihat skor hasil quiz interaktif.
              </p>
            </div>
          </div>

          <div class="flex gap-4 items-start group">
            <div
              class="p-3 bg-slate-50 rounded-xl group-hover:bg-pink-50 transition-colors border border-slate-100"
            >
              <svg
                class="w-6 h-6 text-pink-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                /></svg
              >
            </div>
            <div>
              <h4 class="font-bold text-lg text-slate-900 mb-1">Absensi</h4>
              <p class="text-slate-600 text-sm leading-relaxed">
                Digunakan untuk melihat riwayat kehadiran dan absensi harian.
              </p>
            </div>
          </div>

          <div class="flex gap-4 items-start group">
            <div
              class="p-3 bg-slate-50 rounded-xl group-hover:bg-orange-50 transition-colors border border-slate-100"
            >
              <svg
                class="w-6 h-6 text-orange-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
                /></svg
              >
            </div>
            <div>
              <h4 class="font-bold text-lg text-slate-900 mb-1">
                Perkembangan Menulis
              </h4>
              <p class="text-slate-600 text-sm leading-relaxed">
                Digunakan untuk memantau perkembangan dan kemampuan menulis
                siswa.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
    <!-- Keterangan Menu Utama Section -->
    <section
      class="animate-in fade-in slide-in-from-bottom-8 duration-700 delay-700 mt-8"
    >
      <div
        class="bg-white rounded-3xl p-8 md:p-12 shadow-xl shadow-blue-900/5 border border-slate-100"
      >
        <h2
          class="text-xl md:text-2xl font-bold text-slate-900 mb-8 flex items-center gap-3"
        >
          <div class="p-2 bg-indigo-100 text-indigo-600 rounded-lg">
            <svg
              class="w-6 h-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zm14 0a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zm-9 4a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"
              />
            </svg>
          </div>
          Keterangan Menu Utama
        </h2>

        <div class="grid sm:grid-cols-2 gap-8">
          <div class="flex gap-4 items-start group">
            <div
              class="p-3 bg-slate-50 rounded-xl group-hover:bg-indigo-50 transition-colors border border-slate-100"
            >
              <svg
                class="w-6 h-6 text-indigo-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                /></svg
              >
            </div>
            <div>
              <h4 class="font-bold text-lg text-slate-900 mb-1">
                Mengetik 10 Jari
              </h4>
              <p class="text-slate-600 text-sm leading-relaxed">
                Melatih kecepatan dan ketepatan mengetik menggunakan 10 jari
                agar lebih terampil dan produktif.
              </p>
            </div>
          </div>

          <div class="flex gap-4 items-start group">
            <div
              class="p-3 bg-slate-50 rounded-xl group-hover:bg-cyan-50 transition-colors border border-slate-100"
            >
              <svg
                class="w-6 h-6 text-cyan-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"
                /></svg
              >
            </div>
            <div>
              <h4 class="font-bold text-lg text-slate-900 mb-1">Cetak Code</h4>
              <p class="text-slate-600 text-sm leading-relaxed">
                Fitur untuk mencetak script coding.
              </p>
            </div>
          </div>

          <div class="flex gap-4 items-start group">
            <div
              class="p-3 bg-slate-50 rounded-xl group-hover:bg-rose-50 transition-colors border border-slate-100"
            >
              <svg
                class="w-6 h-6 text-rose-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                /></svg
              >
            </div>
            <div>
              <h4 class="font-bold text-lg text-slate-900 mb-1">Kuis</h4>
              <p class="text-slate-600 text-sm leading-relaxed">
                Menyediakan soal-soal latihan interaktif untuk mengevaluasi dan
                menguji seberapa jauh pemahaman siswa.
              </p>
            </div>
          </div>

          <div class="flex gap-4 items-start group">
            <div
              class="p-3 bg-slate-50 rounded-xl group-hover:bg-amber-50 transition-colors border border-slate-100"
            >
              <svg
                class="w-6 h-6 text-amber-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"
                /></svg
              >
            </div>
            <div>
              <h4 class="font-bold text-lg text-slate-900 mb-1">Berhitung</h4>
              <p class="text-slate-600 text-sm leading-relaxed">
                Latihan berhitung anak PAUD dan TK.
              </p>
            </div>
          </div>

          <div class="flex gap-4 items-start group">
            <div
              class="p-3 bg-slate-50 rounded-xl group-hover:bg-green-50 transition-colors border border-slate-100"
            >
              <svg
                class="w-6 h-6 text-green-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                /></svg
              >
            </div>
            <div>
              <h4 class="font-bold text-lg text-slate-900 mb-1">
                Compress Image
              </h4>
              <p class="text-slate-600 text-sm leading-relaxed">
                Alat untuk mengecilkan ukuran gambar secara aman & privat
                langsung dari browser.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Action Section -->
    <div
      class="flex justify-center pb-8 animate-in fade-in slide-in-from-bottom-8 duration-700 delay-200"
    >
      {#if authChecked}
        {#if !isAuthenticated}
          <a
            href="/login"
            class="w-64 group relative inline-flex items-center justify-center px-8 py-3 text-xs tracking-[0.2em] font-bold uppercase text-white border border-blue-600 hover:border-blue-700 transition-all duration-700 bg-blue-600 hover:bg-blue-700 overflow-hidden cursor-pointer backdrop-blur-md no-underline shadow-md hover:shadow-lg rounded-xl"
          >
            <span class="relative z-10 flex items-center gap-3">
              Masuk Portal
              <svg
                class="w-3.5 h-3.5 transition-transform duration-500 group-hover:translate-x-1 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"
                ></path>
              </svg>
            </span>
            <div
              class="absolute inset-0 -translate-x-full bg-blue-700/50 group-hover:translate-x-0 transition-transform duration-700 ease-out z-0"
            ></div>
          </a>
        {:else}
          <a
            href="/dashboard"
            class="w-64 group relative inline-flex items-center justify-center px-8 py-3 text-xs tracking-[0.2em] font-bold uppercase text-white border border-blue-600 hover:border-blue-700 transition-all duration-700 bg-blue-600 hover:bg-blue-700 overflow-hidden cursor-pointer backdrop-blur-md no-underline shadow-md hover:shadow-lg rounded-xl"
          >
            <span class="relative z-10 flex items-center gap-3">
              Dashboard
              <svg
                class="w-3.5 h-3.5 transition-transform duration-500 group-hover:translate-x-1 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M17 8l4 4m0 0l-4 4m4-4H3"
                ></path>
              </svg>
            </span>
            <div
              class="absolute inset-0 -translate-x-full bg-blue-700/50 group-hover:translate-x-0 transition-transform duration-700 ease-out z-0"
            ></div>
          </a>
        {/if}
      {/if}
    </div>
  </div>
</div>
