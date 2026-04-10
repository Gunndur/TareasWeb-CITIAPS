<!-- components/TaskForm.vue -->
<template>
  <div class="box">
    <h2 class="title is-4">Nueva Tarea</h2>
    <form @submit.prevent="handleSubmit">
      <div class="field">
        <label class="label">Título</label>
        <div class="control">
          <input v-model="form.title" class="input" type="text" placeholder="¿Qué vas a hacer?" required>
        </div>
      </div>

      <div class="field">
        <label class="label">Descripción (opcional)</label>
        <div class="control">
          <textarea v-model="form.description" class="textarea" rows="2"></textarea>
        </div>
      </div>

      <div class="field">
        <label class="label">Etiquetas</label>
        <div class="field has-addons">
          <div class="control is-expanded">
            <input
              v-model="newTag"
              class="input"
              type="text"
              placeholder="Ej: urgente, trabajo"
              @keyup.enter="addTag"
            >
          </div>
          <div class="control">
            <button type="button" class="button is-info" @click="addTag">
              Agregar
            </button>
          </div>
        </div>
        <div class="tags mt-2">
          <span
            v-for="(tag, i) in form.tags"
            :key="i"
            class="tag is-primary is-medium"
          >
            {{ tag }}
            <button class="delete is-small" type="button" @click="removeTag(i)"></button>
          </span>
        </div>
      </div>

      <button
        :disabled="loading"
        class="button is-success is-fullwidth"
        type="submit"
      >
        <span v-if="loading" class="loader is-loading"></span>
        Crear Tarea
      </button>
    </form>

    <p v-if="submitError" class="help is-danger mt-3">
      {{ submitError }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { useTaskStore } from '~/stores/task'

const emit = defineEmits(['created'])

const form = ref({
  title: '',
  description: '',
  tags: [] as string[]
})

const newTag = ref('')

const addTag = () => {
  if (newTag.value.trim()) {
    form.value.tags.push(newTag.value.trim())
    newTag.value = ''
  }
}

const removeTag = (index: number) => {
  form.value.tags.splice(index, 1)
}

const loading = ref(false)
const submitError = ref('')

const handleSubmit = async () => {
  if (!form.value.title) return

  loading.value = true
  submitError.value = ''
  const store = useTaskStore()
  try {
    await store.createTask(form.value)
    emit('created')

    // Reset form
    form.value = { title: '', description: '', tags: [] }
  } catch {
    submitError.value = store.error || 'No se pudo crear la tarea.'
  } finally {
    loading.value = false
  }
}
</script>