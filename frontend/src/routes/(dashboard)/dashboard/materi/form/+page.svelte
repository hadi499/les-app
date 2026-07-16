<script lang="ts">
  import { onMount, tick } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import RichEditor from "$lib/components/RichEditor.svelte";
  import { toast } from "$lib/stores/toast.svelte";

  interface RichEditorRef {
    getHTML(): string;
    setHTML(html: string): void;
  }

  type Subject = {
    id: number;
    name: string;
  };

  type User = {
    id: number;
    username: string;
    role: string;
  };

  let subjects = $state<Subject[]>([]);
  let users = $state<User[]>([]);
  
  let editingId = $state<number | null>(null);
  let formTitle = $state("");
  let formSubjectId = $state<number | "">("");
  let formUserIds = $state<number[]>([]);
  let editorRef = $state<RichEditorRef>();

  let availableStudents = $derived(users.filter(u => u.role !== 'teacher' && !formUserIds.includes(u.id)));
  let availableTeachers = $derived(users.filter(u => u.role === 'teacher' && !formUserIds.includes(u.id)));

  onMount(async () => {
    try {
      const resMe = await fetch("/me", { credentials: "include" });
      if (!resMe.ok) {
        goto("/dashboard");
        return;
      }
      const dataMe = await resMe.json();
      if (dataMe.user?.role !== "teacher") {
        goto("/dashboard/materi");
        return;
      }

      const resSubj = await fetch("/api/subjects", { credentials: "include" });
      if (resSubj.ok) subjects = await resSubj.json();

      const resUsers = await fetch("/api/users", { credentials: "include" });
      if (resUsers.ok) {
        const data = await resUsers.json();
        users = data.users || [];
      }

      const editIdParam = page.url.searchParams.get("id");
      if (editIdParam) {
        editingId = Number(editIdParam);
        const resMateri = await fetch(`/api/materis/${editingId}`, { credentials: "include" });
        if (resMateri.ok) {
          const materiData = await resMateri.json();
          formTitle = materiData.title;
          formSubjectId = materiData.subject_id;
          formUserIds = materiData.users?.map((u: any) => u.id) || [];
          
          await tick();
          if (editorRef) {
            editorRef.setHTML(materiData.content);
          }
        }
      }
    } catch (e) {
      console.error(e);
      toast.error("Gagal memuat data formulir");
    }
  });

  async function saveMateri(e: Event) {
    e.preventDefault();
    const html = editorRef?.getHTML() ?? "";
    if (!formTitle.trim() || !formSubjectId || formUserIds.length === 0 || !(html.trim() || html.includes("<img"))) {
      toast.error("Harap isi semua bidang");
      return;
    }

    const payload = {
      title: formTitle,
      subject_id: Number(formSubjectId),
      user_ids: formUserIds,
      content: html,
    };

    try {
      if (editingId) {
        const res = await fetch(`/api/materis/${editingId}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify(payload),
        });
        if (res.ok) {
          toast.success("Materi berhasil diupdate");
          goto("/dashboard/materi");
        } else {
          toast.error("Gagal mengupdate materi");
        }
      } else {
        const res = await fetch("/api/materis", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify(payload),
        });
        if (res.ok) {
          toast.success("Materi berhasil dibuat");
          goto("/dashboard/materi");
        } else {
          toast.error("Gagal membuat materi");
        }
      }
    } catch (e) {
      console.error(e);
      toast.error("Gagal menyimpan materi");
    }
  }
</script>

<svelte:head>
  <title>{editingId ? 'Edit Materi' : 'Materi Baru'} - Les Balongarut</title>
</svelte:head>

<div class="max-w-5xl mx-auto space-y-6 animate-in fade-in duration-300">
  <!-- Header -->
  <div class="flex items-center gap-4">
    <button onclick={() => goto("/dashboard/materi")} class="p-2 text-slate-500 hover:text-slate-900 bg-white hover:bg-slate-50 rounded-xl border border-slate-200 transition-colors cursor-pointer">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
    </button>
    <div>
      <h1 class="text-2xl font-bold text-slate-900 tracking-tight">{editingId ? 'Edit Materi' : 'Materi Baru'}</h1>
      <p class="text-sm text-slate-500 mt-1">Isi formulir materi pembelajaran.</p>
    </div>
  </div>

  <form onsubmit={saveMateri} class="bg-white rounded-2xl border border-slate-200 p-6 sm:p-8 shadow-sm space-y-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="space-y-2 md:col-span-2">
        <label for="title" class="text-sm font-semibold text-slate-700">Judul Materi</label>
        <input id="title" type="text" bind:value={formTitle} class="w-full px-4 py-3 rounded-xl border border-slate-300 focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 transition-all bg-white text-slate-800 font-medium" placeholder="Contoh: Pengenalan Aljabar" required />
      </div>

      <div class="space-y-2">
        <label for="subject" class="text-sm font-semibold text-slate-700">Mata Pelajaran</label>
        <div class="relative">
          <select id="subject" bind:value={formSubjectId} class="w-full pl-4 pr-10 py-3 rounded-xl border border-slate-300 focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 transition-all bg-white text-slate-800 appearance-none cursor-pointer" required>
            <option value="">Pilih Mata Pelajaran</option>
            {#each subjects as subject}
              <option value={subject.id}>{subject.name}</option>
            {/each}
          </select>
          <div class="absolute inset-y-0 right-0 flex items-center pr-4 pointer-events-none">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
          </div>
        </div>
      </div>

      <div class="space-y-2">
        <label for="user" class="text-sm font-semibold text-slate-700">Tugaskan ke</label>
        
        {#if formUserIds.length > 0}
          <div class="flex flex-wrap gap-2 mb-3">
            {#each formUserIds as id}
              {@const user = users.find(u => u.id === id)}
              {#if user}
                <div class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl {user.role === 'teacher' ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-blue-50 text-blue-700 border-blue-200'} text-sm font-semibold border shadow-sm">
                  {user.username} {user.role === 'teacher' ? '(Guru)' : ''}
                  <button type="button" onclick={() => formUserIds = formUserIds.filter(uid => uid !== id)} class="{user.role === 'teacher' ? 'hover:text-emerald-900 hover:bg-emerald-200' : 'hover:text-blue-900 hover:bg-blue-200'} cursor-pointer p-0.5 rounded-full transition-colors" title="Hapus">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
                  </button>
                </div>
              {/if}
            {/each}
          </div>
        {/if}

        <div class="relative">
          <select id="user" onchange={(e) => {
            const val = e.currentTarget.value;
            if (val === "all_students") {
              const allStudents = users.filter(u => u.role !== 'teacher');
              formUserIds = [...new Set([...formUserIds, ...allStudents.map(u => u.id)])];
            } else if (val === "all_teachers") {
              const allTeachers = users.filter(u => u.role === 'teacher');
              formUserIds = [...new Set([...formUserIds, ...allTeachers.map(u => u.id)])];
            } else {
              const numVal = Number(val);
              if (numVal && !formUserIds.includes(numVal)) {
                formUserIds = [...formUserIds, numVal];
              }
            }
            e.currentTarget.value = "";
          }} class="w-full pl-4 pr-10 py-3 rounded-xl border border-slate-300 focus:border-blue-500 focus:ring-4 focus:ring-blue-500/10 transition-all bg-white text-slate-800 appearance-none cursor-pointer font-medium" required={formUserIds.length === 0}>
            <option value="" disabled selected>+ Tambah Penerima</option>
            
            <optgroup label="Aksi Cepat">
              {#if availableStudents.length > 0}
                <option value="all_students" class="font-bold text-blue-600 bg-blue-50">-- Pilih Semua Siswa --</option>
              {/if}
              {#if availableTeachers.length > 0}
                <option value="all_teachers" class="font-bold text-emerald-600 bg-emerald-50">-- Pilih Semua Guru --</option>
              {/if}
            </optgroup>

            {#if availableTeachers.length > 0}
              <optgroup label="Guru">
                {#each availableTeachers as u}
                  <option value={u.id}>{u.username}</option>
                {/each}
              </optgroup>
            {/if}

            {#if availableStudents.length > 0}
              <optgroup label="Siswa">
                {#each availableStudents as u}
                  <option value={u.id}>{u.username}</option>
                {/each}
              </optgroup>
            {/if}
          </select>
          <div class="absolute inset-y-0 right-0 flex items-center pr-4 pointer-events-none">
            <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
          </div>
        </div>
      </div>
    </div>

    <div class="space-y-2 relative z-20">
      <label for="content" class="text-sm font-semibold text-slate-700">Isi Konten Materi</label>
      <RichEditor 
        bind:this={editorRef}
        minHeight="min-h-[400px] sm:min-h-[550px]"
        containerMinHeight="min-h-[450px] sm:min-h-[600px]"
        textSize="prose-base"
      />
    </div>

    <div class="flex items-center justify-end gap-3 pt-6 border-t border-slate-100">
      <button type="button" onclick={() => goto("/dashboard/materi")} class="px-6 py-3 text-sm font-bold rounded-xl border border-slate-300 text-slate-700 hover:bg-slate-50 transition-colors cursor-pointer">Batal</button>
      <button type="submit" class="px-6 py-3 text-sm font-bold rounded-xl bg-blue-600 text-white hover:bg-blue-700 shadow-md shadow-blue-500/20 transition-all cursor-pointer">Simpan Materi</button>
    </div>
  </form>
</div>
