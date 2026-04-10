import { defineStore } from 'pinia'
import axios from 'axios'
import type { Task } from '~/types/task'

const normalizeTask = (raw: any): Task => ({
  id: raw?.id ?? raw?._id ?? '',
  title: raw?.title ?? '',
  description: raw?.description ?? '',
  completed: Boolean(raw?.completed),
  tags: Array.isArray(raw?.tags) ? raw.tags : [],
  createdAt: raw?.createdAt
})

const normalizeTaskList = (payload: any): Task[] => {
  const source = Array.isArray(payload) ? payload : payload?.items
  if (!Array.isArray(source)) return []
  return source.map(normalizeTask).filter(task => Boolean(task.id))
}

export const useTaskStore = defineStore('task', {
  state: () => ({
    tasks: [] as Task[],
    loading: false,
    error: null as string | null
  }),

  getters: {
    filteredTasks: (state) => (tagFilter: string) => {
      const sortedTasks = [...state.tasks].sort((a, b) =>
        a.title.localeCompare(b.title, 'es', { sensitivity: 'base' })
      )

      if (!tagFilter) return sortedTasks
      return sortedTasks.filter(task =>
        task.tags.some(tag => tag.toLowerCase().includes(tagFilter.toLowerCase()))
      )
    }
  },

  actions: {
    async fetchTasks() {
      this.loading = true
      this.error = null
      try {
        const { public: { apiBase } } = useRuntimeConfig()
        const { data } = await axios.get(`${apiBase}/tasks`)
        this.tasks = normalizeTaskList(data)
      } catch (err: any) {
        this.error = err.response?.data?.message || 'Error al cargar tareas'
      } finally {
        this.loading = false
      }
    },

    // Crea una nueva tarea
    async createTask(taskData: Omit<Task, 'id' | 'completed'>) {
      this.error = null
      try {
        const { public: { apiBase } } = useRuntimeConfig()
        const { data } = await axios.post(`${apiBase}/tasks`, {
          ...taskData,
          completed: false
        })
        const createdTask = normalizeTask(data)
        if (!createdTask.id) {
          throw new Error('Respuesta invalida del servidor al crear la tarea')
        }
        this.tasks.unshift(createdTask)
        return createdTask
      } catch (err: any) {
        this.error = err.response?.data?.message || 'No se pudo crear la tarea. Revisa la conexion con el backend.'
        throw err
      }
    },

    // Marca una tarea como completada
    async markAsCompleted(id: string) {
      if (!id) {
        this.error = 'No se puede completar una tarea sin ID.'
        return
      }
      const { public: { apiBase } } = useRuntimeConfig()
      await axios.put(`${apiBase}/tasks/${id}/complete`)
      const task = this.tasks.find(t => t.id === id)
      if (task) task.completed = true
    },

    // Elimina una tarea
    async deleteTask(id: string) {
      if (!id) {
        this.error = 'No se puede eliminar una tarea sin ID.'
        return
      }
      const { public: { apiBase } } = useRuntimeConfig()
      await axios.delete(`${apiBase}/tasks/${id}`)
      this.tasks = this.tasks.filter(t => t.id !== id)
    }
  }
})