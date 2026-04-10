<!-- components/TaskList.vue -->
<template>
  <div>
    <div class="field">
      <label class="label">Filtrar por etiqueta</label>
      <input
        v-model="tagFilter"
        class="input"
        type="text"
        placeholder="Ej: Frontend"
      >
    </div>

    <div v-if="store.loading" class="has-text-centered py-6">
      <span class="loader is-loading"></span>
    </div>

    <div v-else-if="filteredTasks.length === 0" class="has-text-centered py-6">
      <p class="has-text-grey">No hay tareas aún</p>
    </div>

    <TaskItem
      v-for="task in paginatedTasks"
      :key="task.id"
      :task="task"
      @deleted="store.fetchTasks"
    />

    <nav
      v-if="filteredTasks.length > 0 && totalPages > 1"
      class="pagination is-centered mt-5"
      role="navigation"
      aria-label="pagination"
    >
      <button
        class="pagination-previous"
        :disabled="currentPage === 1"
        @click="goToPreviousPage"
      >
        Anterior
      </button>

      <button
        class="pagination-next"
        :disabled="currentPage === totalPages"
        @click="goToNextPage"
      >
        Siguiente
      </button>

      <ul class="pagination-list">
        <li v-for="page in pages" :key="page">
          <button
            class="pagination-link"
            :class="{ 'is-current': page === currentPage }"
            @click="goToPage(page)"
          >
            {{ page }}
          </button>
        </li>
      </ul>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { useTaskStore } from '~/stores/task'
import TaskItem from './TaskItem.vue'

const store = useTaskStore()
const tagFilter = ref('')
const currentPage = ref(1)
const pageSize = 5

const filteredTasks = computed(() => store.filteredTasks(tagFilter.value))
const totalPages = computed(() => Math.ceil(filteredTasks.value.length / pageSize))

const paginatedTasks = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  return filteredTasks.value.slice(start, end)
})

const pages = computed(() =>
  Array.from({ length: totalPages.value }, (_, index) => index + 1)
)

const goToPage = (page: number) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
}

const goToPreviousPage = () => {
  goToPage(currentPage.value - 1)
}

const goToNextPage = () => {
  goToPage(currentPage.value + 1)
}

watch(tagFilter, () => {
  currentPage.value = 1
})

watch(totalPages, (newTotalPages) => {
  if (newTotalPages === 0) {
    currentPage.value = 1
    return
  }
  if (currentPage.value > newTotalPages) {
    currentPage.value = newTotalPages
  }
})
</script>