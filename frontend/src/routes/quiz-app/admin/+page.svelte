<script lang="ts">
  import { onMount, getContext } from "svelte";
  import { fade } from "svelte/transition";
  import { toast } from "$lib/stores/toast.svelte";

  const quizUser: any = getContext("quizUser");

  let stats = $state({
    totalQuizzes: 0,
    totalCategories: 0,
    totalParticipants: 0,
  });

  let categories = $state<any[]>([]);
  let quizzes = $state<any[]>([]);
  let allResults = $state<any[]>([]);
  let activeTab = $state("quizzes");
  let isLoading = $state(true);

  // Pagination for Results
  let resultsPage = $state(1);
  const resultsPerPage = 30;
  let totalResults = $state(0);
  let isFetchingResults = $state(false);

  // Forms states
  let showCategoryForm = $state(false);
  let newCategory = $state({ name: "" });

  let showDeleteQuizModal = $state(false);
  let quizToDeleteId = $state<number | null>(null);
  let showResetPointsModal = $state(false);
  let showDeleteCategoryModal = $state(false);
  let categoryToDeleteId = $state<number | null>(null);

  let showQuizForm = $state(false);
  let editingQuizId = $state<number | null>(null);
  let newQuiz = $state({
    title: "",
    timeLimit: 30,
    category_id: null as number | null,
    is_published: true,
  });

  // Users states
  let allUsers = $state<any[]>([]);
  let uniqueParticipants = $derived(
    allUsers
      .filter((u) => u.points > 0)
      .sort((a: any, b: any) => b.points - a.points),
  );
  let userSearchQuery = $state("");
  let searchTimeout: any;

  async function fetchUsers() {
    try {
      const q = userSearchQuery
        ? `?search=${encodeURIComponent(userSearchQuery)}`
        : "";
      const res = await fetch(`/api/kuisapp/users${q}`, {
        credentials: "include",
      });
      if (res.ok) {
        const d = await res.json();
        allUsers = d.data || [];
      }
    } catch (e) {
      console.error(e);
    }
  }

  function handleSearchInput() {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      fetchUsers();
    }, 300);
  }

  let showUserForm = $state(false);
  let isEditingUser = $state(false);
  let editingUserId = $state<number | null>(null);
  let newUserForm = $state({
    username: "",
    password: "",
    role: "user",
    status: "umum",
    points: 0,
  });
  let showPassword = $state(false);

  let showDeleteUserModal = $state(false);
  let userToDeleteId = $state<number | null>(null);

  let showResetUserModal = $state(false);
  let userToResetId = $state<number | null>(null);
  let resetUserPassword = $state("");
  let showResetPassword = $state(false);

  onMount(async () => {
    const savedTab = localStorage.getItem("kuisappAdminTab");
    if (savedTab) activeTab = savedTab;
    await fetchData();
  });

  $effect(() => {
    localStorage.setItem("kuisappAdminTab", activeTab);
  });

  async function fetchResults(page = 1, append = false) {
    if (isFetchingResults) return;
    isFetchingResults = true;
    try {
      const res = await fetch(
        `/api/kuisapp/all-results?page=${page}&limit=${resultsPerPage}`,
        { credentials: "include" },
      );
      if (res.ok) {
        const d = await res.json();
        const results = d.data || [];
        if (append) {
          allResults = [...allResults, ...results];
        } else {
          allResults = results;
        }
        totalResults = d.total || 0;
        resultsPage = page;
      }
    } catch (e) {
      console.error(e);
    } finally {
      isFetchingResults = false;
    }
  }

  async function fetchData() {
    isLoading = true;
    try {
      const [catRes, quizRes, userRes] = await Promise.all([
        fetch("/api/kuisapp/categories", { credentials: "include" }),
        fetch("/api/kuisapp/quizzes", { credentials: "include" }),
        fetch("/api/kuisapp/users", { credentials: "include" }),
      ]);

      await fetchResults(1, false);

      if (catRes.ok) {
        const d = await catRes.json();
        categories = d.data || [];
        stats.totalCategories = categories.length;
      }
      if (quizRes.ok) {
        const d = await quizRes.json();
        quizzes = d.data || [];
        stats.totalQuizzes = quizzes.length;
      }
      if (userRes.ok) {
        const d = await userRes.json();
        allUsers = d.data || [];
        stats.totalParticipants = allUsers.length;
      }
    } catch (e) {
      console.error(e);
    } finally {
      isLoading = false;
    }
  }

  async function handleAddCategory(e: Event) {
    e.preventDefault();
    const res = await fetch("/api/kuisapp/categories", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify(newCategory),
    });
    if (res.ok) {
      showCategoryForm = false;
      newCategory = { name: "" };
      await fetchData();
    }
  }

  async function handleDeleteCategory() {
    if (!categoryToDeleteId) return;
    try {
      const res = await fetch(`/api/kuisapp/categories/${categoryToDeleteId}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (res.ok) {
        toast.success("Kategori berhasil dihapus");
        await fetchData();
      } else {
        toast.error("Gagal menghapus kategori");
      }
    } catch (err) {
      toast.error("Terjadi kesalahan saat menghapus kategori");
    } finally {
      showDeleteCategoryModal = false;
      categoryToDeleteId = null;
    }
  }

  let isDuplicating: Record<number, boolean> = $state({});

  async function duplicateQuiz(id: number) {
    isDuplicating[id] = true;
    try {
      const res = await fetch(`/api/kuisapp/quizzes/${id}/duplicate`, {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
        await fetchData();
      } else {
        alert("Gagal menduplikasi kuis");
      }
    } catch (e) {
      console.error(e);
      alert("Terjadi kesalahan jaringan");
    } finally {
      isDuplicating[id] = false;
    }
  }

  async function handleDeleteQuiz() {
    if (quizToDeleteId === null) return;
    await fetch(`/api/kuisapp/quizzes/${quizToDeleteId}`, {
      method: "DELETE",
      credentials: "include",
    });
    await fetchData();
    showDeleteQuizModal = false;
    quizToDeleteId = null;
  }

  function formatDate(dateString: string) {
    const d = new Date(dateString);
    const datePart = new Intl.DateTimeFormat("id-ID", {
      day: "numeric",
      month: "short",
      year: "numeric",
    }).format(d);

    const timePart = [
      d.getHours().toString().padStart(2, "0"),
      d.getMinutes().toString().padStart(2, "0"),
      d.getSeconds().toString().padStart(2, "0"),
    ].join(":");

    return `${datePart}, ${timePart} WIB`;
  }

  function formatDuration(seconds: number) {
    if (!seconds || seconds < 0) return "0 Detik";
    return `${seconds} Detik`;
  }

  function getScoreColor(score: number) {
    if (score >= 80) return "text-emerald-600 bg-emerald-50 border-emerald-100";
    if (score >= 60) return "text-yellow-600 bg-yellow-50 border-yellow-100";
    return "text-red-600 bg-red-50 border-red-100";
  }

  // --- User Management Handlers ---

  function openAddUser() {
    isEditingUser = false;
    showPassword = false;
    newUserForm = {
      username: "",
      password: "",
      role: "user",
      status: "umum",
      points: 0,
    };
    showUserForm = true;
  }

  function openEditUser(user: any) {
    isEditingUser = true;
    editingUserId = user.id;
    newUserForm = {
      username: user.username,
      password: "",
      role: user.role,
      status: user.status,
      points: user.points,
    };
    showUserForm = true;
  }

  async function handleSaveUser(e: Event) {
    e.preventDefault();
    const url = isEditingUser
      ? `/api/kuisapp/users/${editingUserId}`
      : "/api/kuisapp/users";
    const method = isEditingUser ? "PUT" : "POST";

    // Validasi saat Add user
    if (!isEditingUser && !newUserForm.password) {
      alert("Password wajib diisi untuk user baru");
      return;
    }

    const payload: any = {
      username: newUserForm.username,
      role: newUserForm.role,
      status: newUserForm.status,
      points: Number(newUserForm.points),
    };
    if (!isEditingUser) payload.password = newUserForm.password;

    try {
      const res = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        showUserForm = false;
        await fetchData();
      } else {
        const d = await res.json();
        alert(d.error || "Gagal menyimpan user");
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan jaringan");
    }
  }

  async function handleDeleteUser() {
    if (!userToDeleteId) return;
    try {
      const res = await fetch(`/api/kuisapp/users/${userToDeleteId}`, {
        method: "DELETE",
        credentials: "include",
      });
      if (res.ok) {
        showDeleteUserModal = false;
        await fetchData();
      } else {
        alert("Gagal menghapus user");
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan");
    }
  }

  async function handleResetPassword(e: Event) {
    e.preventDefault();
    if (!userToResetId || !resetUserPassword) return;
    try {
      const res = await fetch(
        `/api/kuisapp/users/${userToResetId}/reset-password`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({ password: resetUserPassword }),
        },
      );
      if (res.ok) {
        showResetUserModal = false;
        resetUserPassword = "";
        alert("Password berhasil direset!");
      } else {
        alert("Gagal reset password");
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan");
    }
  }

  async function handleToggleSuspend(user: any) {
    if (
      !confirm(
        `Yakin ingin ${user.is_suspended ? "mengaktifkan" : "men-suspend"} akun ${user.username}?`,
      )
    )
      return;
    try {
      const res = await fetch(`/api/kuisapp/users/${user.id}/suspend`, {
        method: "PUT",
        credentials: "include",
      });
      if (res.ok) {
        await fetchData();
      } else {
        alert("Gagal mengubah status suspend user");
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan");
    }
  }

  async function handleResetAllPointsAndHistory() {
    try {
      const res = await fetch(`/api/kuisapp/reset-all-points-and-history`, {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
        await fetchData();
        toast.success(`Seluruh poin berhasil direset!`);
        showResetPointsModal = false;
      } else {
        toast.error("Gagal mereset poin");
        showResetPointsModal = false;
      }
    } catch (err) {
      console.error(err);
      toast.error("Terjadi kesalahan saat mereset");
      showResetPointsModal = false;
    }
  }
</script>

<svelte:head>
  <title>Quiz | Admin</title>
</svelte:head>

<div class="space-y-8 pb-12">
  <!-- Admin Header -->
  <div
    class="bg-white rounded-3xl p-4 sm:p-6 border border-slate-100 shadow-sm flex flex-col md:flex-row md:items-center justify-between gap-6"
  >
    <div>
      <div class="relative z-10 flex items-center justify-between mb-2">
        <div class="flex items-center">
          <span
            class="font-black tracking-wider uppercase text-md sm:text-lg drop-shadow-sm text-blue-600"
            >LB Quiz</span
          >
        </div>
      </div>

      <h1 class="text-2xl sm:text-3xl font-black text-slate-800 mb-3">
        Panel Administrator
      </h1>
      <p
        class="text-slate-500 text-sm font-semibold px-2 rounded-md border border-slate-100 shadow-sm shadow-amber-600 w-fit"
      >
        {quizUser?.current?.username || "Admin"}
      </p>
    </div>
    <div class="flex gap-4">
      <div
        class="bg-indigo-50 px-5 py-3 rounded-2xl border border-indigo-100 text-center min-w-30"
      >
        <div
          class="text-indigo-600 text-sm font-bold uppercase tracking-wider mb-1"
        >
          Total Kuis
        </div>
        <div class="text-3xl font-black text-indigo-900">
          {stats.totalQuizzes}
        </div>
      </div>
      <div
        class="bg-violet-50 px-5 py-3 rounded-2xl border border-violet-100 text-center min-w-30"
      >
        <div
          class="text-violet-600 text-sm font-bold uppercase tracking-wider mb-1"
        >
          Peserta
        </div>
        <div class="text-3xl font-black text-violet-900">
          {stats.totalParticipants}
        </div>
      </div>
    </div>
  </div>

  <!-- Tabs Navigation -->
  <div class="mb-6 sm:mb-8">
    <!-- Mobile Dropdown -->
    <div class="sm:hidden mb-2 relative">
      <label for="tabs" class="sr-only">Pilih Menu Admin</label>
      <select
        id="tabs"
        bind:value={activeTab}
        class="block w-full appearance-none rounded-xl border border-slate-200 py-3 pl-4 pr-12 text-base focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 bg-white shadow-sm font-bold text-slate-800"
      >
        <option value="quizzes">Kuis</option>
        <option value="categories">Kategori</option>
        <option value="results">Score Peserta</option>
        <option value="leaderboard">Poin Peserta</option>
        <option value="users">Manajemen User</option>
      </select>
      <div
        class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-4 text-slate-500"
      >
        <svg
          class="w-5 h-5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 9l-7 7-7-7"
          /></svg
        >
      </div>
    </div>

    <!-- Desktop Tabs -->
    <div class="hidden sm:flex overflow-x-auto border-b border-slate-200">
      <button
        onclick={() => (activeTab = "quizzes")}
        class="px-6 py-4 text-sm font-bold transition-colors border-b-2 whitespace-nowrap shrink-0 {activeTab ===
        'quizzes'
          ? 'border-indigo-600 text-indigo-600'
          : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'}"
      >
        Kuis
      </button>
      <button
        onclick={() => (activeTab = "categories")}
        class="px-6 py-4 text-sm font-bold transition-colors border-b-2 whitespace-nowrap shrink-0 {activeTab ===
        'categories'
          ? 'border-indigo-600 text-indigo-600'
          : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'}"
      >
        Kategori
      </button>
      <button
        onclick={() => (activeTab = "results")}
        class="px-6 py-4 text-sm font-bold transition-colors border-b-2 whitespace-nowrap shrink-0 {activeTab ===
        'results'
          ? 'border-indigo-600 text-indigo-600'
          : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'}"
      >
        Score Peserta
      </button>
      <button
        onclick={() => (activeTab = "leaderboard")}
        class="px-6 py-4 text-sm font-bold transition-colors border-b-2 whitespace-nowrap shrink-0 {activeTab ===
        'leaderboard'
          ? 'border-indigo-600 text-indigo-600'
          : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'}"
      >
        Poin Peserta
      </button>
      <button
        onclick={() => (activeTab = "users")}
        class="px-6 py-4 text-sm font-bold transition-colors border-b-2 whitespace-nowrap shrink-0 {activeTab ===
        'users'
          ? 'border-indigo-600 text-indigo-600'
          : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'}"
      >
        Manajemen User
      </button>
    </div>
  </div>

  <!-- Loading State -->
  {#if isLoading}
    <div class="flex justify-center py-12">
      <div
        class="w-10 h-10 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin"
      ></div>
    </div>
  {:else}
    <!-- Content Section -->
    {#if activeTab === "quizzes"}
      <section class="animate-in fade-in slide-in-from-bottom-2 duration-300">
        <div
          class="flex flex-col sm:flex-row sm:justify-between items-start sm:items-center gap-4 mb-6"
        >
          <h2 class="text-xl font-bold text-slate-800">Daftar Kuis</h2>
          <a
            href="/quiz-app/admin/quizzes/create"
            class="inline-flex items-center gap-2 px-4 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-bold rounded-xl shadow-sm transition-all no-underline"
          >
            <svg
              class="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 4v16m8-8H4"
              /></svg
            >
            Buat Kuis Baru
          </a>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {#each quizzes as quiz}
            <div
              class="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm hover:shadow-md transition-shadow relative flex flex-col h-full"
            >
              <div class="flex justify-between items-start mb-2">
                <div>
                  <div class="flex items-center gap-2 mb-1">
                    {#if quiz.is_published}
                      <span
                        class="px-1.5 py-0.5 rounded bg-green-100 text-green-700 text-[10px] font-black uppercase tracking-wider"
                        >Publish</span
                      >
                    {:else}
                      <span
                        class="px-1.5 py-0.5 rounded bg-slate-200 text-slate-700 text-[10px] font-black uppercase tracking-wider"
                        >Draft</span
                      >
                    {/if}
                  </div>
                  <h3 class="text-lg font-bold text-slate-900 leading-tight">
                    {quiz.title}
                  </h3>
                </div>
                {#if quiz.category_id}
                  <span
                    class="px-2 py-1 bg-slate-100 text-slate-600 rounded-md text-[10px] font-bold border border-slate-200 whitespace-nowrap uppercase tracking-wide"
                  >
                    {categories.find((c) => c.id === quiz.category_id)?.name ||
                      "Tanpa Kategori"}
                  </span>
                {/if}
              </div>

              <div
                class="flex items-center text-xs font-medium text-slate-500 mb-6 bg-slate-50 w-max px-2.5 py-1.5 rounded-md border border-slate-100 mt-2"
              >
                {quiz.timeLimit} Detik / Soal
              </div>

              <div
                class="mt-auto flex flex-wrap gap-2 pt-4 border-t border-slate-100"
              >
                <a
                  href="/quiz-app/admin/quizzes/{quiz.id}/edit"
                  class="flex-2 bg-indigo-50 hover:bg-indigo-100 text-indigo-700 font-semibold py-2 px-3 rounded-xl transition-colors text-center text-sm flex items-center justify-center gap-2 no-underline"
                >
                  Edit Kuis & Soal
                </a>
                <button
                  onclick={() => duplicateQuiz(quiz.id)}
                  disabled={isDuplicating[quiz.id]}
                  class="flex-1 bg-amber-50 hover:bg-amber-100 text-amber-700 p-2 rounded-xl transition-colors flex items-center justify-center {isDuplicating[
                    quiz.id
                  ]
                    ? 'opacity-50 cursor-not-allowed'
                    : ''}"
                  title="Duplikasi Kuis"
                >
                  {#if isDuplicating[quiz.id]}
                    <div
                      class="w-4 h-4 border-2 border-amber-300 border-t-amber-600 rounded-full animate-spin"
                    ></div>
                  {:else}
                    <svg
                      class="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      ><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                      /></svg
                    >
                  {/if}
                </button>
                <button
                  onclick={() => {
                    quizToDeleteId = quiz.id;
                    showDeleteQuizModal = true;
                  }}
                  class="flex-1 bg-red-50 hover:bg-red-100 text-red-600 p-2 rounded-xl transition-colors flex items-center justify-center"
                  title="Hapus Kuis"
                >
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                    /></svg
                  >
                </button>
              </div>
            </div>
          {:else}
            <div
              class="col-span-full text-center p-12 bg-white rounded-2xl border border-slate-200 shadow-sm text-sm text-slate-500"
            >
              Belum ada kuis yang dibuat.
            </div>
          {/each}
        </div>
      </section>
    {:else if activeTab === "categories"}
      <section class="animate-in fade-in slide-in-from-bottom-2 duration-300">
        <div
          class="flex flex-col sm:flex-row sm:justify-between items-start sm:items-center gap-4 mb-6"
        >
          <h2 class="text-xl font-bold text-slate-800">Daftar Kategori</h2>
          <button
            onclick={() => (showCategoryForm = true)}
            class="inline-flex items-center gap-2 px-4 py-2.5 bg-white border border-slate-300 hover:bg-slate-50 text-slate-900 text-sm font-bold rounded-xl shadow-sm transition-all"
          >
            <svg
              class="w-5 h-5 text-slate-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 4v16m8-8H4"
              /></svg
            >
            Tambah Kategori
          </button>
        </div>

        {#if showCategoryForm}
          <div
            class="mb-8 p-6 bg-slate-50 rounded-2xl border border-slate-200 max-w-sm"
          >
            <h3 class="font-bold text-lg mb-4">Tambah Kategori Baru</h3>
            <form onsubmit={handleAddCategory} class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-slate-700"
                  >Nama Kategori</label
                >
                <input
                  type="text"
                  required
                  bind:value={newCategory.name}
                  class="mt-1 block w-full rounded-lg border-slate-300 bg-white text-slate-900 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                />
              </div>
              <div class="flex gap-3 pt-2">
                <button
                  type="submit"
                  class="bg-indigo-600 text-white px-4 py-2 rounded-lg text-sm font-bold hover:bg-indigo-700"
                  >Simpan</button
                >
                <button
                  type="button"
                  onclick={() => (showCategoryForm = false)}
                  class="bg-white border border-slate-300 text-slate-700 px-4 py-2 rounded-lg text-sm font-bold hover:bg-slate-50"
                  >Batal</button
                >
              </div>
            </form>
          </div>
        {/if}

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
          {#if categories.length === 0}
            <div
              class="col-span-full py-8 text-center text-slate-500 bg-white rounded-2xl border border-slate-100"
            >
              Belum ada kategori.
            </div>
          {/if}
          {#each categories as cat}
            <div
              class="bg-white rounded-2xl p-6 border border-slate-100 shadow-sm relative group"
            >
              <div class="absolute top-4 right-4 flex gap-2">
                <button
                  onclick={() => {
                    categoryToDeleteId = cat.id;
                    showDeleteCategoryModal = true;
                  }}
                  class="p-1.5 text-slate-400 hover:text-red-600 bg-slate-50 hover:bg-red-50 rounded-lg transition-colors"
                  title="Hapus"
                  ><svg
                    class="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                    /></svg
                  ></button
                >
              </div>
              <h3 class="text-lg font-bold text-slate-800 mb-4 pr-10">
                {cat.name}
              </h3>
              <div
                class="text-xs font-semibold text-indigo-500 uppercase tracking-wider bg-indigo-50 inline-block px-3 py-1 rounded-full"
              >
                {quizzes.filter((q) => q.category_id === cat.id).length} Kuis Terhubung
              </div>
            </div>
          {/each}
        </div>
      </section>
    {:else if activeTab === "results"}
      <section class="animate-in fade-in slide-in-from-bottom-2 duration-300">
        <div
          class="flex flex-col sm:flex-row sm:justify-between items-start sm:items-center gap-4 mb-6"
        >
          <h2 class="text-xl font-bold text-slate-800">Riwayat Kuis Peserta</h2>
        </div>

        {#if allResults.length === 0}
          <div
            class="text-center p-12 bg-white rounded-2xl border border-slate-200 shadow-sm text-sm text-slate-500"
          >
            Belum ada data riwayat peserta.
          </div>
        {:else}
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
            {#each allResults as result}
              <div
                class="bg-white rounded-2xl p-6 border border-slate-100 shadow-sm hover:shadow-md transition-shadow"
              >
                <div class="flex justify-between items-start mb-4">
                  <div>
                    <div class="font-bold text-slate-900 text-lg">
                      {result.user?.username || "Tidak Diketahui"}
                    </div>
                    <div class="text-sm font-medium text-indigo-600 mt-1">
                      {result.quiz?.title || "Kuis Tidak Diketahui"}
                    </div>
                  </div>
                  <span
                    class="inline-flex items-center justify-center px-3 py-1 rounded-xl text-sm font-bold border shrink-0 ml-4 {getScoreColor(
                      result.score,
                    )}"
                  >
                    {Math.round(result.score)}
                  </span>
                </div>

                <div
                  class="flex items-center justify-between mt-4 pt-4 border-t border-slate-100"
                >
                  <div class="flex flex-col gap-1">
                    <div
                      class="text-xs font-medium text-slate-500 flex items-center gap-1.5"
                    >
                      <svg
                        class="w-4 h-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        ><path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                        /></svg
                      >
                      {formatDate(result.finished_at)}
                    </div>
                    {#if result.duration > 0}
                      <div
                        class="text-xs font-medium text-slate-500 flex items-center gap-1.5"
                      >
                        <svg
                          class="w-4 h-4"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                          ><path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                          ></path></svg
                        >
                        {formatDuration(result.duration)}
                      </div>
                    {/if}
                  </div>
                  <div>
                    {#if result.points_earned > 0}
                      <span
                        class="inline-flex items-center justify-center px-2.5 py-1 rounded-lg text-xs font-bold bg-emerald-100 text-emerald-700 border border-emerald-200"
                      >
                        +{result.points_earned} Poin
                      </span>
                    {:else}
                      <span class="text-slate-400 text-xs font-medium"
                        >0 Poin</span
                      >
                    {/if}
                  </div>
                </div>
              </div>
            {/each}
          </div>

          <!-- Load More Control -->
          {#if allResults.length < totalResults}
            <div class="flex justify-center mt-8">
              <button
                onclick={() => fetchResults(resultsPage + 1, true)}
                disabled={isFetchingResults}
                class="px-6 py-2.5 bg-white border border-slate-200 text-slate-700 rounded-xl font-bold hover:bg-slate-50 hover:border-slate-300 shadow-sm transition-all flex items-center gap-2 disabled:opacity-50"
              >
                {isFetchingResults ? "Memuat..." : "Lihat lebih banyak"}
                {#if !isFetchingResults}
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M19 9l-7 7-7-7"
                    /></svg
                  >
                {/if}
              </button>
            </div>
          {/if}
        {/if}
      </section>
    {:else if activeTab === "leaderboard"}
      <section class="animate-in fade-in slide-in-from-bottom-2 duration-300">
        <div
          class="flex flex-col sm:flex-row sm:justify-between items-start sm:items-center gap-4 mb-6"
        >
          <h2 class="text-xl font-bold text-slate-800">Daftar Poin Peserta</h2>
          <button
            onclick={() => (showResetPointsModal = true)}
            class="py-2 px-4 bg-red-100 hover:bg-red-200 text-red-700 rounded-lg text-sm font-bold flex items-center justify-center gap-2 transition-colors border border-red-200"
            title="Reset Seluruh Poin & Riwayat"
          >
            <svg
              class="w-4 h-4 shrink-0"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              /></svg
            >
            <span>Reset Poin</span>
          </button>
        </div>

        {#if uniqueParticipants.length === 0}
          <div
            class="text-center p-12 bg-white rounded-2xl border border-slate-200 shadow-sm text-sm text-slate-500"
          >
            Belum ada peserta yang mengumpulkan poin.
          </div>
        {:else}
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
            {#each uniqueParticipants as user, index}
              <div
                class="bg-white rounded-2xl p-6 border border-slate-100 shadow-sm hover:shadow-md transition-shadow relative overflow-hidden flex flex-col items-center text-center group"
              >
                <div
                  class="w-16 h-16 rounded-full bg-indigo-50 flex items-center justify-center text-indigo-600 font-bold text-xl mb-4 group-hover:scale-110 transition-transform"
                >
                  <svg class="w-8 h-8" fill="currentColor" viewBox="0 0 24 24"
                    ><path
                      d="M12 12a5 5 0 100-10 5 5 0 000 10zm0 2c-5.33 0-8 2.67-8 8v1h16v-1c0-5.33-2.67-8-8-8z"
                    /></svg
                  >
                </div>
                <h3 class="font-bold text-slate-900 text-lg mb-1">
                  {user.username}
                </h3>
                <span class="text-sm text-slate-500 capitalize mb-4"
                  >{user.status}</span
                >

                <div
                  class="mt-auto w-full bg-slate-50 py-3 rounded-xl border border-slate-100 mb-3"
                >
                  <div
                    class="text-amber-500 font-black text-2xl drop-shadow-sm"
                  >
                    {user.points}
                  </div>
                  <div
                    class="text-xs font-bold text-slate-400 uppercase tracking-wider mt-1"
                  >
                    Total Poin
                  </div>
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </section>
    {/if}

    <!-- USER MANAGEMENT SECTION -->
    {#if activeTab === "users"}
      <section class="animate-in fade-in slide-in-from-bottom-2 duration-300">
        <div
          class="flex flex-col sm:flex-row sm:justify-between items-start sm:items-center gap-4 mb-6"
        >
          <h2 class="text-xl font-bold text-slate-800">Daftar Pengguna</h2>
          <div class="flex items-center gap-3 w-full sm:w-auto">
            <div class="relative flex-1 sm:w-64">
              <svg
                class="w-5 h-5 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                /></svg
              >
              <input
                type="text"
                bind:value={userSearchQuery}
                oninput={handleSearchInput}
                placeholder="Cari user..."
                class="w-full pl-10 pr-10 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all text-sm"
              />
              {#if userSearchQuery}
                <button
                  type="button"
                  onclick={() => {
                    userSearchQuery = "";
                    fetchUsers();
                  }}
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 focus:outline-none"
                  title="Hapus pencarian"
                >
                  <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M6 18L18 6M6 6l12 12"
                    /></svg
                  >
                </button>
              {/if}
            </div>
            <button
              onclick={openAddUser}
              class="shrink-0 bg-indigo-600 hover:bg-indigo-700 text-white px-4 sm:px-5 py-2.5 rounded-xl font-bold shadow-sm hover:shadow-md transition-all flex items-center gap-2"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                /></svg
              >
              <span class="hidden sm:inline">Tambah User</span>
            </button>
          </div>
        </div>

        {#if allUsers.length === 0}
          {#if userSearchQuery}
            <div
              class="bg-white rounded-3xl p-12 text-center border border-slate-100 shadow-sm"
            >
              <div class="text-6xl mb-4">🔍</div>
              <h3 class="text-xl font-bold text-slate-800 mb-2">
                Pencarian Tidak Ditemukan
              </h3>
              <p class="text-slate-500">
                Tidak ada user yang cocok dengan "{userSearchQuery}".
              </p>
            </div>
          {:else}
            <div
              class="bg-white rounded-3xl p-12 text-center border border-slate-100 shadow-sm"
            >
              <div class="text-6xl mb-4">👥</div>
              <h3 class="text-xl font-bold text-slate-800 mb-2">
                Belum ada user
              </h3>
              <p class="text-slate-500">
                Tambahkan user untuk mengizinkan mereka login dan bermain kuis.
              </p>
            </div>
          {/if}
        {:else}
          <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {#each allUsers as user}
              <div
                class="bg-white p-5 rounded-2xl border {user.is_suspended
                  ? 'border-red-200 bg-red-50/30 opacity-75'
                  : 'border-slate-100'} shadow-sm hover:shadow-md transition-all"
              >
                <div class="flex items-center gap-3 mb-4">
                  <div
                    class="w-12 h-12 rounded-full {user.is_suspended
                      ? 'bg-red-100 text-red-600'
                      : 'bg-indigo-50 text-indigo-600'} flex items-center justify-center font-bold text-lg"
                  >
                    {user.username.substring(0, 2).toUpperCase()}
                  </div>
                  <div class="flex-1">
                    <h4
                      class="font-bold {user.is_suspended
                        ? 'text-red-800 line-through'
                        : 'text-slate-800'}"
                    >
                      {user.username}
                    </h4>
                    <div class="flex flex-col gap-1 items-start mt-1">
                      <span
                        class="inline-flex items-center px-2 py-0.5 rounded text-xs font-semibold uppercase {user.is_suspended
                          ? 'bg-red-100 text-red-700'
                          : user.role === 'admin'
                            ? 'bg-amber-100 text-amber-800'
                            : 'bg-slate-100 text-slate-700'}"
                      >
                        {user.is_suspended ? "Disuspend" : user.role} &bull; {user.status}
                      </span>
                      {#if user.created_at}
                        <span
                          class="text-[11px] text-slate-500 font-medium flex items-center gap-1"
                        >
                          <svg
                            class="w-3 h-3"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                            ><path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                            ></path></svg
                          >
                          {formatDate(user.created_at)}
                        </span>
                      {/if}
                    </div>
                  </div>
                </div>
                <div class="flex gap-2">
                  <button
                    onclick={() => openEditUser(user)}
                    class="flex-1 flex justify-center items-center py-2 bg-slate-50 hover:bg-slate-100 text-slate-500 hover:text-slate-700 rounded-lg transition-colors border border-slate-200"
                    title="Edit User"
                  >
                    <svg
                      class="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      ><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                      /></svg
                    >
                  </button>
                  <button
                    onclick={() => {
                      userToResetId = user.id;
                      showResetPassword = false;
                      resetUserPassword = "";
                      showResetUserModal = true;
                    }}
                    class="flex-1 flex justify-center items-center py-2 bg-blue-50 hover:bg-blue-100 text-blue-500 hover:text-blue-700 rounded-lg transition-colors border border-blue-200"
                    title="Reset Password"
                  >
                    <svg
                      class="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      ><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                      /></svg
                    >
                  </button>
                  <button
                    onclick={() => handleToggleSuspend(user)}
                    class="flex-1 flex justify-center items-center py-2 {user.is_suspended
                      ? 'bg-emerald-50 hover:bg-emerald-100 text-emerald-500 hover:text-emerald-700 border-emerald-200'
                      : 'bg-orange-50 hover:bg-orange-100 text-orange-500 hover:text-orange-700 border-orange-200'} rounded-lg transition-colors border"
                    title={user.is_suspended ? "Aktifkan User" : "Suspend User"}
                  >
                    {#if user.is_suspended}
                      <svg
                        class="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        ><path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M5 13l4 4L19 7"
                        /></svg
                      >
                    {:else}
                      <svg
                        class="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        ><path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
                        /></svg
                      >
                    {/if}
                  </button>
                  <button
                    onclick={() => {
                      userToDeleteId = user.id;
                      showDeleteUserModal = true;
                    }}
                    class="flex-1 flex justify-center items-center py-2 bg-red-50 hover:bg-red-100 text-red-500 hover:text-red-700 rounded-lg transition-colors border border-red-200"
                    title="Hapus User"
                  >
                    <svg
                      class="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      ><path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      /></svg
                    >
                  </button>
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </section>
    {/if}
  {/if}
</div>

<!-- Modal Form User -->
{#if showUserForm}
  <div
    class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
  >
    <div
      class="bg-white rounded-3xl w-full max-w-md shadow-2xl overflow-hidden"
      in:fade={{ duration: 200 }}
    >
      <div class="p-6 border-b border-slate-100">
        <h3 class="text-xl font-bold text-slate-800">
          {isEditingUser ? "Edit User" : "Tambah User Baru"}
        </h3>
      </div>
      <form onsubmit={handleSaveUser} class="p-6 space-y-4">
        <div>
          <label
            class="block text-sm font-bold text-slate-700 mb-1"
            for="username">Username</label
          >
          <input
            type="text"
            id="username"
            bind:value={newUserForm.username}
            required
            class="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
          />
        </div>
        {#if !isEditingUser}
          <div>
            <label
              class="block text-sm font-bold text-slate-700 mb-1"
              for="password">Password</label
            >
            <div class="relative">
              <input
                type={showPassword ? "text" : "password"}
                id="password"
                bind:value={newUserForm.password}
                required
                class="w-full pl-4 pr-10 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
              />
              <button
                type="button"
                onclick={() => (showPassword = !showPassword)}
                class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 focus:outline-none"
                title={showPassword
                  ? "Sembunyikan password"
                  : "Tampilkan password"}
              >
                {#if showPassword}
                  <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                    /></svg
                  >
                {:else}
                  <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    ><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    /><path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                    /></svg
                  >
                {/if}
              </button>
            </div>
          </div>
        {/if}
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label
              class="block text-sm font-bold text-slate-700 mb-1"
              for="role">Peran (Role)</label
            >
            <select
              id="role"
              bind:value={newUserForm.role}
              class="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all bg-white"
            >
              <option value="user">User</option>
              <option value="admin">Admin</option>
            </select>
          </div>
          <div>
            <label
              class="block text-sm font-bold text-slate-700 mb-1"
              for="status">Status</label
            >
            <select
              id="status"
              bind:value={newUserForm.status}
              class="w-full px-4 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all bg-white"
            >
              <option value="umum">Umum</option>
              <option value="pelajar">Pelajar</option>
            </select>
          </div>
        </div>
        <div class="flex gap-3 pt-4 border-t border-slate-100">
          <button
            type="button"
            onclick={() => (showUserForm = false)}
            class="flex-1 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold rounded-xl transition-colors"
            >Batal</button
          >
          <button
            type="submit"
            class="flex-1 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl transition-colors"
            >Simpan</button
          >
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- Modal Delete User -->
{#if showDeleteUserModal}
  <div
    class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
  >
    <div
      class="bg-white rounded-3xl w-full max-w-sm shadow-2xl overflow-hidden p-6 text-center"
    >
      <div
        class="w-16 h-16 rounded-full bg-red-100 text-red-600 flex items-center justify-center mx-auto mb-4"
      >
        <svg
          class="w-8 h-8"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
          /></svg
        >
      </div>
      <h3 class="text-xl font-bold text-slate-800 mb-2">Hapus User?</h3>
      <p class="text-slate-500 mb-6">
        User <strong class="text-slate-800"
          >{allUsers.find((u) => u.id === userToDeleteId)?.username ||
            ""}</strong
        > dan semua riwayat pengerjaan kuisnya akan dihapus permanen. Lanjutkan?
      </p>
      <div class="flex gap-3">
        <button
          onclick={() => (showDeleteUserModal = false)}
          class="flex-1 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold rounded-xl transition-colors"
          >Batal</button
        >
        <button
          onclick={handleDeleteUser}
          class="flex-1 px-4 py-2 bg-red-600 hover:bg-red-700 text-white font-bold rounded-xl transition-colors"
          >Ya, Hapus</button
        >
      </div>
    </div>
  </div>
{/if}

<!-- Modal Delete Quiz -->
{#if showDeleteQuizModal}
  <div
    class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
  >
    <div
      class="bg-white rounded-3xl w-full max-w-sm shadow-2xl overflow-hidden p-6 text-center"
    >
      <div
        class="w-16 h-16 rounded-full bg-red-100 text-red-600 flex items-center justify-center mx-auto mb-4"
      >
        <svg
          class="w-8 h-8"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
          /></svg
        >
      </div>
      <h3 class="text-xl font-bold text-slate-800 mb-2">Hapus Kuis?</h3>
      <p class="text-slate-500 mb-6">
        Kuis <strong class="text-slate-800"
          >{quizzes.find((q) => q.id === quizToDeleteId)?.title || ""}</strong
        > beserta seluruh soalnya akan dihapus permanen. Lanjutkan?
      </p>
      <div class="flex gap-3">
        <button
          onclick={() => {
            showDeleteQuizModal = false;
            quizToDeleteId = null;
          }}
          class="flex-1 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold rounded-xl transition-colors"
          >Batal</button
        >
        <button
          onclick={handleDeleteQuiz}
          class="flex-1 px-4 py-2 bg-red-600 hover:bg-red-700 text-white font-bold rounded-xl transition-colors"
          >Ya, Hapus</button
        >
      </div>
    </div>
  </div>
{/if}

<!-- Modal Reset Password -->
{#if showResetUserModal}
  <div
    class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
  >
    <div
      class="bg-white rounded-3xl w-full max-w-sm shadow-2xl overflow-hidden p-6 text-center"
    >
      <h3 class="text-xl font-bold text-slate-800 mb-2">Reset Password</h3>
      <p class="text-slate-500 mb-4 text-sm">
        Masukkan password baru untuk <strong class="text-slate-800"
          >{allUsers.find((u) => u.id === userToResetId)?.username ||
            "user ini"}</strong
        >.
      </p>
      <form onsubmit={handleResetPassword} class="space-y-4">
        <label for="reset-password" class="sr-only">Password Baru</label>
        <div class="relative">
          <input
            type={showResetPassword ? "text" : "password"}
            id="reset-password"
            bind:value={resetUserPassword}
            placeholder="Password baru"
            required
            class="w-full pl-4 pr-10 py-2 rounded-xl border border-slate-200 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
          />
          <button
            type="button"
            onclick={() => (showResetPassword = !showResetPassword)}
            class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 focus:outline-none"
            title={showResetPassword
              ? "Sembunyikan password"
              : "Tampilkan password"}
          >
            {#if showResetPassword}
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                /></svg
              >
            {:else}
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                /><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                /></svg
              >
            {/if}
          </button>
        </div>
        <div class="flex gap-3">
          <button
            type="button"
            onclick={() => {
              showResetUserModal = false;
              resetUserPassword = "";
            }}
            class="flex-1 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold rounded-xl transition-colors"
            >Batal</button
          >
          <button
            type="submit"
            class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-bold rounded-xl transition-colors"
            >Simpan</button
          >
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- Modal Reset All Points and History -->
{#if showResetPointsModal}
  <div
    class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
  >
    <div
      class="bg-white rounded-3xl w-full max-w-sm shadow-2xl overflow-hidden p-6 text-center"
    >
      <div
        class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4 text-red-600"
      >
        <svg
          class="w-8 h-8"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          /></svg
        >
      </div>
      <h3 class="text-xl font-bold text-slate-800 mb-2">
        Reset Semua Poin & Riwayat?
      </h3>
      <p class="text-slate-500 mb-6 text-sm">
        PERINGATAN: Seluruh poin peserta dan riwayat jawaban kuis akan dihapus.
        Tindakan ini tidak dapat dibatalkan!
      </p>
      <div class="flex gap-3">
        <button
          onclick={() => (showResetPointsModal = false)}
          class="flex-1 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold rounded-xl transition-colors"
          >Batal</button
        >
        <button
          onclick={handleResetAllPointsAndHistory}
          class="flex-1 px-4 py-2 bg-red-600 hover:bg-red-700 text-white font-bold rounded-xl transition-colors"
          >Ya, Reset</button
        >
      </div>
    </div>
  </div>
{/if}

<!-- Modal Hapus Kategori -->
{#if showDeleteCategoryModal}
  <div
    class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
  >
    <div
      class="bg-white rounded-3xl w-full max-w-sm shadow-2xl overflow-hidden p-6 text-center"
    >
      <div
        class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4 text-red-600"
      >
        <svg
          class="w-8 h-8"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
          /></svg
        >
      </div>
      <h3 class="text-xl font-bold text-slate-800 mb-2">Hapus Kategori?</h3>
      <p class="text-slate-500 mb-6 text-sm">
        Kategori <strong class="text-slate-800"
          >{categories.find((c) => c.id === categoryToDeleteId)?.name ||
            ""}</strong
        > akan dihapus permanen. Kuis yang ada di dalamnya mungkin akan kehilangan
        label kategori.
      </p>
      <div class="flex gap-3">
        <button
          onclick={() => {
            showDeleteCategoryModal = false;
            categoryToDeleteId = null;
          }}
          class="flex-1 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold rounded-xl transition-colors"
          >Batal</button
        >
        <button
          onclick={handleDeleteCategory}
          class="flex-1 px-4 py-2 bg-red-600 hover:bg-red-700 text-white font-bold rounded-xl transition-colors"
          >Ya, Hapus</button
        >
      </div>
    </div>
  </div>
{/if}
