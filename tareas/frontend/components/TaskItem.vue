<!-- components/TaskItem.vue -->
<template>
  <div class="card mb-4" :class="{ 'has-background-success': task.completed }">
    <div class="card-content">
      <div class="level">
        <div class="level-left">
          <p class="is-size-5" :class="{ 'has-text-white-ter line-through': task.completed, 'has-text-black': !task.completed }">
            <strong>{{ task.title }}</strong>
          </p>
        </div>
        <div class="level-right">
          <div class="buttons are-small">
            <button 
              v-if="!task.completed"
              class="button is-success is-light"
              @click="showConfirmModal = true"
            >
              <span class="icon is-small">
                <i class="fas fa-check"></i>
              </span>
              <span>Completar</span>
            </button>
            <span v-else class="tag is-success is-light">
              <span class="icon is-small">
                <i class="fas fa-check-circle"></i>
              </span>
              <span>Completada</span>
            </span>
            <button class="delete" @click="showDeleteConfirmModal = true"></button>
          </div>
        </div>
      </div>

      <!-- Modal de confirmación para completar -->
      <div class="modal" :class="{ 'is-active': showConfirmModal }">
        <div class="modal-background" @click="showConfirmModal = false"></div>
        <div class="modal-card">
          <header class="modal-card-head">
            <p class="modal-card-title">Confirmar completar tarea</p>
            <button 
              class="delete" 
              @click="showConfirmModal = false"
            ></button>
          </header>
          <section class="modal-card-body">
            <p>¿Deseas marcar como completada la tarea <strong>"{{ task.title }}"</strong>?</p>
          </section>
          <footer class="modal-card-foot">
            <button 
              class="button"
              @click="showConfirmModal = false"
            >
              Cancelar
            </button>
            <button 
              class="button is-success"
              @click="confirmComplete"
            >
              Sí, completar
            </button>
          </footer>
        </div>
      </div>

      <!-- Modal de confirmación para eliminar -->
      <div class="modal" :class="{ 'is-active': showDeleteConfirmModal }">
        <div class="modal-background" @click="showDeleteConfirmModal = false"></div>
        <div class="modal-card">
          <header class="modal-card-head">
            <p class="modal-card-title">Confirmar eliminar tarea</p>
            <button 
              class="delete" 
              @click="showDeleteConfirmModal = false"
            ></button>
          </header>
          <section class="modal-card-body">
            <p>¿Deseas eliminar la tarea <strong>"{{ task.title }}"</strong>? Esta acción no se puede deshacer.</p>
          </section>
          <footer class="modal-card-foot">
            <button 
              class="button"
              @click="showDeleteConfirmModal = false"
            >
              Cancelar
            </button>
            <button 
              class="button is-danger"
              @click="confirmDelete"
            >
              Sí, eliminar
            </button>
          </footer>
        </div>
      </div>

      <div v-if="task.description" class="box is-size-6 has-text-white has-background">
        {{ task.description }}
      </div>

      <div class="tags">
        <span
          v-for="tag in task.tags"
          :key="tag"
          class="tag is-info is-light is-size-6 has-text-black"
        >
          #{{ tag }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTaskStore } from '~/stores/task';
import type { Task } from '~/types/task'

const props = defineProps<{ task: Task }>()
const emit = defineEmits(['deleted'])

const store = useTaskStore()
const showConfirmModal = ref(false)
const showDeleteConfirmModal = ref(false)

const markComplete = async () => {
  if (!props.task.id || props.task.completed) return
  try {
    await store.markAsCompleted(props.task.id)
  } catch {
    // El estado de error se guarda en el store para mostrarlo donde corresponda.
  }
}

const confirmComplete = async () => {
  showConfirmModal.value = false
  await markComplete()
}

const confirmDelete = async () => {
  showDeleteConfirmModal.value = false
  if (!props.task.id) return
  try {
    await store.deleteTask(props.task.id)
    emit('deleted')
  } catch {
    // Evita que Vue reporte un error no capturado por evento nativo.
  }
}
</script>