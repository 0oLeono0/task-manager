import { defineStore } from 'pinia'
import {
  createTask as createTaskRequest,
  deleteTask,
  getTasks,
  updateTask,
  type Task,
} from '../api'

interface TasksState {
  tasks: Task[]
  input: string
  isLoading: boolean
  isCreating: boolean
  updatingTaskId: number | null
  deletingTaskId: number | null
  errText: string
  formErrText: string
  actionErrText: string
}

export const useTasksStore = defineStore('tasks', {
  state: (): TasksState => ({
    tasks: [],
    input: '',
    isLoading: true,
    isCreating: false,
    updatingTaskId: null,
    deletingTaskId: null,
    errText: '',
    formErrText: '',
    actionErrText: '',
  }),

  actions: {
    async loadTasks() {
      try {
        const fetchedTasks = await getTasks()
        this.tasks = fetchedTasks
      } catch {
        this.errText = 'Не удалось загрузить задачи. Попробуйте обновить страницу.'
      } finally {
        this.isLoading = false
      }
    },

    async createTask() {
      if (this.isCreating) return

      const trimmedTitle = this.input.trim()
      if (trimmedTitle === '') return

      this.formErrText = ''
      this.isCreating = true

      try {
        const createdTask = await createTaskRequest({ title: trimmedTitle })
        this.tasks.push(createdTask)
        this.input = ''
      } catch {
        this.formErrText = 'Не удалось создать задачу. Попробуйте еще раз.'
      } finally {
        this.isCreating = false
      }
    },

    async toggleTaskCompleted(id: number) {
      if (this.updatingTaskId !== null) return

      const task = this.tasks.find((task) => task.id === id)
      if (!task) return

      this.actionErrText = ''
      this.updatingTaskId = id

      try {
        const updatedTask = await updateTask(id, { completed: !task.completed })
        const taskIndex = this.tasks.indexOf(task)
        if (taskIndex === -1) return

        this.tasks[taskIndex] = updatedTask
      } catch {
        this.actionErrText = 'Не удалось обновить задачу. Попробуйте еще раз.'
      } finally {
        this.updatingTaskId = null
      }
    },

    async deleteTaskById(id: number) {
      if (this.deletingTaskId !== null) return

      this.actionErrText = ''
      this.deletingTaskId = id

      try {
        await deleteTask(id)
        this.tasks = this.tasks.filter((task) => task.id !== id)
      } catch {
        this.actionErrText = 'Не удалось удалить задачу. Попробуйте еще раз.'
      } finally {
        this.deletingTaskId = null
      }
    },
  },
})
