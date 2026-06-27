<script lang="ts">
  import { fade, slide, scale } from "svelte/transition";
  import { flip } from "svelte/animate";
  import { onMount } from "svelte";

  type Student = {
    id: string;
    username: string; // From the User model
    role: string;
  };

  type TargetHafalan = {
    id: number;
    student_id: number;
    text: string;
    completed: boolean;
    created_at: string;
  };

  let students = $state<Student[]>([]);
  let targets = $state<TargetHafalan[]>([]);

  let activeStudentId = $state<string>("");
  let newTask = $state("");
  let filter = $state<"all" | "active" | "completed">("all");
  let isLoading = $state(true);
  let errorMsg = $state("");
  let showClearModal = $state(false);

  onMount(async () => {
    try {
      // Ambil daftar murid
      const usersRes = await fetch("/api/users", { credentials: "include" });
      if (!usersRes.ok) throw new Error("Gagal mengambil data murid");
      const usersData = await usersRes.json();

      // Filter hanya murid
      if (usersData.users) {
        students = usersData.users.filter((u: any) => u.role === "student");
      }

      if (students.length > 0) {
        const savedId = localStorage.getItem("lastActiveStudentId");
        if (savedId && students.find((s) => s.id.toString() === savedId)) {
          activeStudentId = savedId;
        } else {
          activeStudentId = students[0].id.toString();
        }
      }

      // Ambil semua target
      await loadTargets();
    } catch (err: any) {
      errorMsg = err.message;
    } finally {
      isLoading = false;
    }
  });

  async function loadTargets() {
    try {
      const res = await fetch("/api/targets", { credentials: "include" });
      if (res.ok) {
        targets = await res.json();
      }
    } catch (err) {
      console.error("Gagal memuat target:", err);
    }
  }

  $effect(() => {
    if (activeStudentId) {
      localStorage.setItem("lastActiveStudentId", activeStudentId);
    }
  });

  let activeStudent = $derived(
    students.find((s) => s.id.toString() === activeStudentId),
  );

  let studentTargets = $derived(
    targets.filter((t) => t.student_id.toString() === activeStudentId),
  );

  let filteredTargets = $derived(
    studentTargets.filter((target) => {
      if (filter === "active") return !target.completed;
      if (filter === "completed") return target.completed;
      return true;
    }),
  );

  let stats = $derived({
    total: studentTargets.length,
    active: studentTargets.filter((t) => !t.completed).length,
    completed: studentTargets.filter((t) => t.completed).length,
  });

  async function addTarget(e?: Event) {
    if (e) e.preventDefault();
    if (!newTask.trim() || !activeStudentId) return;

    try {
      const res = await fetch("/api/targets", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          student_id: parseInt(activeStudentId),
          text: newTask.trim(),
        }),
      });

      if (res.ok) {
        const newTarget = await res.json();
        targets = [...targets, newTarget];
        newTask = "";
      } else {
        const errData = await res.json();
        alert(errData.error || "Gagal menambahkan target");
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan sistem");
    }
  }

  async function toggleTarget(id: number, currentStatus: boolean) {
    // Optimistic update
    const targetIndex = targets.findIndex((t) => t.id === id);
    if (targetIndex > -1) {
      targets[targetIndex].completed = !currentStatus;
    }

    try {
      const res = await fetch(`/api/targets/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ completed: !currentStatus }),
      });

      if (!res.ok) {
        // Revert on fail
        if (targetIndex > -1) {
          targets[targetIndex].completed = currentStatus;
        }
      }
    } catch (err) {
      console.error(err);
      // Revert on fail
      if (targetIndex > -1) {
        targets[targetIndex].completed = currentStatus;
      }
    }
  }

  async function deleteTarget(id: number) {
    // Optimistic delete
    const originalTargets = [...targets];
    targets = targets.filter((t) => t.id !== id);

    try {
      const res = await fetch(`/api/targets/${id}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (!res.ok) {
        targets = originalTargets;
        alert("Gagal menghapus target");
      }
    } catch (err) {
      console.error(err);
      targets = originalTargets;
    }
  }

  function promptClearCompleted() {
    const completedIds = studentTargets
      .filter((t) => t.completed)
      .map((t) => t.id);
    if (completedIds.length === 0) return;
    showClearModal = true;
  }

  async function confirmClearCompleted() {
    showClearModal = false;
    const completedIds = studentTargets
      .filter((t) => t.completed)
      .map((t) => t.id);
    if (completedIds.length === 0) return;

    // We can delete them one by one or create a bulk delete API.
    // Since we only have DELETE /:id, we will delete one by one for now.
    const originalTargets = [...targets];
    targets = targets.filter(
      (t) => !(t.student_id.toString() === activeStudentId && t.completed),
    );

    for (const id of completedIds) {
      try {
        await fetch(`/api/targets/${id}`, {
          method: "DELETE",
          credentials: "include",
        });
      } catch (err) {
        console.error("Gagal menghapus", id);
      }
    }
  }
</script>

<svelte:head>
  <title>Target Hafalan Murid - Portal</title>
</svelte:head>

<div class="max-w-6xl mx-auto p-4 md:p-6 lg:p-8">
  <div class="mb-8 flex flex-col md:flex-row md:items-start justify-between gap-6">
    <div>
      <h1
        class="text-xl md:text-2xl font-bold text-transparent bg-clip-text bg-linear-to-r from-blue-600 to-indigo-600 tracking-tight mb-2"
      >
        Target Hafalan Murid
      </h1>
      <p class="text-slate-500 font-medium">
        Kelola dan pantau target hafalan kartu memori untuk setiap murid secara
        individu.
      </p>
    </div>


  </div>

  {#if isLoading}
    <div class="flex items-center justify-center p-12">
      <div
        class="w-10 h-10 border-4 border-slate-200 border-t-blue-500 rounded-full animate-spin"
      ></div>
    </div>
  {:else if errorMsg}
    <div class="p-6 bg-red-50 text-red-600 rounded-xl border border-red-200">
      <p class="font-semibold">{errorMsg}</p>
    </div>
  {:else}
    <div class="flex flex-col lg:flex-row gap-6">
      <!-- Sidebar Students -->
      <div class="w-full lg:w-1/3 xl:w-1/4 space-y-4">
        <div class="bg-white rounded-2xl border border-slate-200 shadow-sm p-4">
          <div class="flex items-center justify-between mb-4">
            <h2
              class="text-sm font-bold text-slate-800 uppercase tracking-wider"
            >
              Daftar Murid
            </h2>
          </div>

          <div class="space-y-1 max-h-[400px] overflow-y-auto">
            {#if students.length === 0}
              <div class="p-4 text-center text-sm text-slate-500">
                Belum ada murid yang terdaftar.
              </div>
            {/if}
            {#each students as student (student.id)}
              <div
                role="button"
                tabindex="0"
                class="group flex items-center justify-between px-3 py-2.5 rounded-xl transition-all cursor-pointer border outline-none {activeStudentId ===
                student.id.toString()
                  ? 'bg-blue-50 border-blue-200 text-blue-700'
                  : 'bg-transparent border-transparent hover:bg-slate-50 text-slate-600'}"
                onclick={() => (activeStudentId = student.id.toString())}
                onkeydown={(e) => {
                  if (e.key === "Enter" || e.key === " ")
                    activeStudentId = student.id.toString();
                }}
              >
                <span class="font-semibold text-sm truncate capitalize"
                  >{student.username}</span
                >
              </div>
            {/each}
          </div>
        </div>


      </div>

      <!-- Main Target Area -->
      <div class="w-full lg:w-2/3 xl:w-3/4 space-y-6">
        <!-- Active Student Header -->
        {#if activeStudent}
          <div class="flex items-center gap-4 border-b border-slate-200 pb-4">
            <div>
              <h2 class="text-xl font-semibold text-slate-800 capitalize">
                {activeStudent.username}
              </h2>
              <p class="text-slate-500 text-sm font-medium mt-1">
                {stats.completed} dari {stats.total} target selesai 
                {#if stats.total > 0}
                  <span class="text-blue-600 font-bold ml-1">({Math.round((stats.completed / stats.total) * 100)}%)</span>
                {/if}
              </p>
            </div>
          </div>

          <!-- Add Target Form -->
          <div
            class="bg-white/80 backdrop-blur-xl border border-slate-200/60 shadow-lg shadow-slate-200/40 rounded-xl md:rounded-2xl p-1.5 md:p-2 transition-all focus-within:shadow-blue-500/10 focus-within:border-blue-300"
          >
            <form onsubmit={addTarget} class="flex items-center gap-2 md:gap-3">
              <div class="flex-1 flex items-center gap-2 md:gap-3 pl-2 md:pl-3">

                <input
                  type="text"
                  bind:value={newTask}
                  placeholder="Target hafalan ..."
                  class="w-full bg-transparent border-none outline-none text-slate-700 placeholder-slate-400 py-2 md:py-3 text-sm md:text-lg font-normal"
                />
              </div>
              <button
                type="submit"
                disabled={!newTask.trim()}
                class="bg-blue-600 hover:bg-blue-700 disabled:bg-slate-300 disabled:cursor-not-allowed text-white font-semibold text-sm md:text-base p-2 md:p-2.5 rounded-lg md:rounded-xl transition-all duration-300 shadow-md hover:shadow-blue-500/30 shrink-0 cursor-pointer"
              >
                <svg
                  class="w-5 h-5 md:w-6 md:h-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2.5"
                    d="M12 4v16m8-8H4"
                  ></path>
                </svg>
              </button>
            </form>
          </div>

          <!-- Target List -->
          <div
            class="mt-6"
          >
            <!-- List Header / Filters -->
            <div
              class="flex flex-col sm:flex-row sm:items-center justify-between py-2 gap-4"
            >
              <div class="flex gap-2 bg-slate-200/50 p-1 rounded-xl w-max">
                <button
                  onclick={() => (filter = "all")}
                  class="px-4 py-1.5 rounded-lg text-sm font-semibold transition-all cursor-pointer {filter ===
                  'all'
                    ? 'bg-white shadow-sm text-blue-600'
                    : 'text-slate-500 hover:text-slate-700'}"
                >
                  Semua
                </button>
                <button
                  onclick={() => (filter = "active")}
                  class="px-4 py-1.5 rounded-lg text-sm font-semibold transition-all cursor-pointer {filter ===
                  'active'
                    ? 'bg-white shadow-sm text-blue-600'
                    : 'text-slate-500 hover:text-slate-700'}"
                >
                  Belum
                </button>
                <button
                  onclick={() => (filter = "completed")}
                  class="px-4 py-1.5 rounded-lg text-sm font-semibold transition-all cursor-pointer {filter ===
                  'completed'
                    ? 'bg-white shadow-sm text-blue-600'
                    : 'text-slate-500 hover:text-slate-700'}"
                >
                  Sudah Hafal
                </button>
              </div>
              {#if stats.completed > 0}
                <button
                  onclick={promptClearCompleted}
                  class="text-sm font-medium text-red-500 hover:text-red-700 transition-colors flex items-center gap-1.5 cursor-pointer bg-red-50 hover:bg-red-100 px-3 py-1.5 rounded-lg"
                  transition:fade
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
                    ></path></svg
                  >
                  Bersihkan Hafalan
                </button>
              {/if}
            </div>

            <!-- List Items -->
            <div class="min-h-[300px] flex flex-col gap-1">
              {#if filteredTargets.length === 0}
                <div
                  class="p-12 text-center text-slate-500 flex flex-col items-center justify-center h-full"
                  in:fade
                >
                  <div
                    class="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mb-4 text-slate-400"
                  >
                    <svg
                      class="w-8 h-8"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="1.5"
                        d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"
                      ></path>
                    </svg>
                  </div>

                  <p class="font-medium text-lg">Belum ada target</p>
                  <p class="text-sm mt-1">
                    Tambahkan target hafalan baru untuk <span
                      class="font-semibold text-slate-700 capitalize"
                      >{activeStudent.username}</span
                    >!
                  </p>
                </div>
              {:else}
                {#each filteredTargets as target (target.id)}
                  <div
                    animate:flip={{ duration: 300 }}
                    in:slide|local
                    out:slide|local
                    class="group flex items-center justify-between p-3 hover:bg-white rounded-xl transition-colors"
                  >
                    <div
                      role="button"
                      tabindex="0"
                      class="flex items-center gap-4 flex-1 cursor-pointer outline-none"
                      onclick={() => toggleTarget(target.id, target.completed)}
                      onkeydown={(e) => {
                        if (e.key === "Enter" || e.key === " ")
                          toggleTarget(target.id, target.completed);
                      }}
                    >
                      <div
                        class="relative flex items-center justify-center w-6 h-6 shrink-0"
                      >
                        <input
                          type="checkbox"
                          checked={target.completed}
                          tabindex="-1"
                          class="peer appearance-none w-5 h-5 border-2 border-slate-300 rounded-md checked:bg-blue-500 checked:border-blue-500 transition-all cursor-pointer pointer-events-none"
                        />
                        <svg
                          class="absolute w-4 h-4 text-white pointer-events-none opacity-0 peer-checked:opacity-100 transition-opacity"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="3"
                            d="M5 13l4 4L19 7"
                          ></path>
                        </svg>
                      </div>
                      <span
                        class="text-lg font-semibold transition-all duration-300 {target.completed
                          ? 'text-slate-400 italic line-through decoration-slate-500 decoration-2'
                          : 'text-slate-700'}"
                      >
                        {target.text}
                      </span>
                    </div>

                    <button
                      onclick={() => deleteTarget(target.id)}
                      class="opacity-100 md:opacity-0 md:group-hover:opacity-100 p-2 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-all cursor-pointer"
                      title="Hapus target"
                    >
                      <svg
                        class="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        ></path>
                      </svg>
                    </button>
                  </div>
                {/each}
              {/if}
            </div>
          </div>
        {:else}
          <div
            class="bg-white rounded-2xl border border-slate-200 shadow-sm p-12 text-center text-slate-500"
          >
            Pilih murid dari daftar di samping untuk melihat atau mengelola
            target hafalan.
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>

{#if showClearModal}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    transition:fade={{ duration: 150 }}
    onclick={() => (showClearModal = false)}
  >
    <div
      class="bg-white rounded-2xl shadow-xl max-w-sm w-full p-6 border border-slate-100"
      transition:scale={{ duration: 150, start: 0.95 }}
      onclick={(e) => e.stopPropagation()}
    >
      <div
        class="w-12 h-12 rounded-full bg-red-100 text-red-500 flex items-center justify-center mb-4 mx-auto"
      >
        <svg
          class="w-6 h-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          ></path>
        </svg>
      </div>
      <h3 class="text-lg font-bold text-slate-800 text-center mb-2">
        Bersihkan Hafalan?
      </h3>
      <p class="text-slate-600 text-sm text-center mb-6">
        Anda yakin ingin menghapus secara permanen semua target yang sudah
        dihafal oleh <span class="font-semibold text-slate-700 capitalize"
          >{activeStudent?.username}</span
        >? Tindakan ini tidak dapat dibatalkan.
      </p>
      <div class="flex gap-3">
        <button
          onclick={() => (showClearModal = false)}
          class="flex-1 px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold rounded-xl transition-colors cursor-pointer"
        >
          Batal
        </button>
        <button
          onclick={confirmClearCompleted}
          class="flex-1 px-4 py-2 bg-red-500 hover:bg-red-600 text-white font-semibold rounded-xl transition-colors shadow-sm shadow-red-500/30 cursor-pointer"
        >
          Ya, Hapus
        </button>
      </div>
    </div>
  </div>
{/if}
