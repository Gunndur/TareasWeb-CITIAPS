<!-- pages/index.vue -->
<template>
  <div class="columns">
    <!-- Formulario -->
    <div class="column is-4">
      <TaskForm @created="store.fetchTasks" />
    </div>

    <!-- Lista -->
    <div class="column is-8">
      <h1 class="title is-3 mb-5">Mis Tareas</h1>
      <TaskList />
    </div>
  </div>
</template>

<script setup lang="ts">
import TaskForm from '~/components/TaskForm.vue'
import TaskList from '~/components/TaskList.vue'
import { useTaskStore } from '~/stores/task'

const store = useTaskStore()

onMounted(() => {
  store.fetchTasks()
})
</script>