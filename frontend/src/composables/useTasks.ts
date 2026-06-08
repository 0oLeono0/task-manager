import { onMounted, ref } from 'vue'
import { createTask, deleteTask, getTasks, updateTask, type Task } from '../api'

export const useTasks = () => {
  const tasks = ref<Task[]>([])
  const input = ref('')
  const isLoading = ref(true)
  const isCreating = ref(false)
  const updatingTaskId = ref<number | null>(null)
  const deletingTaskId = ref<number | null>(null)
  const errText = ref('')
  const formErrText = ref('')
  const actionErrText = ref('')

  const loadTasks = async () => {
    try {
      const fetchedTasks = await getTasks()
      tasks.value = fetchedTasks
    } catch {
      errText.value = 'Не удалось загрузить задачи. Попробуйте обновить страницу.'
    } finally {
      isLoading.value = false
    }
  }

  const onSubmit = async () => {
    if (isCreating.value) return

    const trimmedTitle = input.value.trim()
    if (trimmedTitle === '') return

    formErrText.value = ''
    isCreating.value = true

    try {
      const createdTask = await createTask({ title: trimmedTitle })
      tasks.value.push(createdTask)
      input.value = ''
    } catch {
      formErrText.value = 'Не удалось создать задачу. Попробуйте еще раз.'
    } finally {
      isCreating.value = false
    }
  }

  const toggleTaskCompleted = async (id: number) => {
    if (updatingTaskId.value !== null) return

    const task = tasks.value.find((task) => task.id === id)
    if (!task) return

    actionErrText.value = ''
    updatingTaskId.value = id

    try {
      const updatedTask = await updateTask(id, { completed: !task.completed })
      const taskIndex = tasks.value.indexOf(task)
      if (taskIndex === -1) return

      tasks.value[taskIndex] = updatedTask
    } catch {
      actionErrText.value = 'Не удалось обновить задачу. Попробуйте еще раз.'
    } finally {
      updatingTaskId.value = null
    }
  }

  const deleteTaskById = async (id: number) => {
    if (deletingTaskId.value !== null) return

    actionErrText.value = ''
    deletingTaskId.value = id

    try {
      await deleteTask(id)
      tasks.value = tasks.value.filter((task) => task.id !== id)
    } catch {
      actionErrText.value = 'Не удалось удалить задачу. Попробуйте еще раз.'
    } finally {
      deletingTaskId.value = null
    }
  }

  onMounted(loadTasks)

  return {
    tasks,
    input,
    isLoading,
    isCreating,
    updatingTaskId,
    deletingTaskId,
    errText,
    formErrText,
    actionErrText,
    onSubmit,
    toggleTaskCompleted,
    deleteTaskById,
  }
}
