<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { createTask, deleteTask, getTasks, updateTask, type Task } from './api'

const tasks = ref<Task[]>([])
const input = ref('')
const isLoading = ref(true)
const errText = ref('')
const formErrText = ref('')
const actionErrText = ref('')

onMounted(async () => {
  try {
    const fetchedTasks = await getTasks()
    tasks.value = fetchedTasks
  } catch {
    errText.value = 'Не удалось загрузить задачи. Попробуйте обновить страницу.'
  } finally {
    isLoading.value = false
  }
})

const onSubmit = async () => {
  const trimmedTitle = input.value.trim()
  if (trimmedTitle === '') return

  formErrText.value = ''

  try {
    const createdTask = await createTask({ title: trimmedTitle })
    tasks.value.push(createdTask)
    input.value = ''
  } catch {
    formErrText.value = 'Не удалось создать задачу. Попробуйте еще раз.'
  }
}

const toggleTaskCompleted = async (id: number) => {
  const task = tasks.value.find((task) => task.id === id)
  if (!task) return

  actionErrText.value = ''

  try {
    const updatedTask = await updateTask(id, { completed: !task.completed })
    const taskIndex = tasks.value.indexOf(task)
    if (taskIndex === -1) return

    tasks.value[taskIndex] = updatedTask
  } catch {
    actionErrText.value = 'Не удалось обновить задачу. Попробуйте еще раз.'
  }
}

const deleteTaskById = async (id: number) => {
  actionErrText.value = ''

  try {
    await deleteTask(id)
    tasks.value = tasks.value.filter((task) => task.id !== id)
  } catch {
    actionErrText.value = 'Не удалось удалить задачу. Попробуйте еще раз.'
  }
}
</script>

<template>
  <main class="app-shell">
    <section class="task-panel" aria-labelledby="app-title">
      <header class="app-header">
        <p class="app-kicker">Менеджер задач</p>
        <h1 id="app-title">Мои задачи</h1>
      </header>

      <form class="task-form" @submit.prevent="onSubmit">
        <label class="sr-only" for="task-title">Новая задача</label>
        <input id="task-title" v-model="input" type="text" placeholder="Например: проверить API" />
        <button type="submit">Добавить</button>
      </form>
      <p v-if="formErrText" class="task-message task-message--error">{{ formErrText }}</p>

      <p v-if="isLoading" class="task-message">Загружаем задачи...</p>
      <p v-else-if="errText" class="task-message task-message--error">{{ errText }}</p>

      <template v-else>
        <p v-if="actionErrText" class="task-message task-message--error">{{ actionErrText }}</p>

        <ul class="task-list" aria-label="Список задач">
          <li
            v-for="task in tasks"
            :key="task.id"
            class="task-item"
            :class="{ 'task-item--done': task.completed }"
          >
            <button
              type="button"
              class="task-status"
              :aria-pressed="task.completed"
              :aria-label="
                task.completed
                  ? `Отметить невыполненной: ${task.title}`
                  : `Отметить выполненной: ${task.title}`
              "
              @click="toggleTaskCompleted(task.id)"
            ></button>
            <span class="task-title">{{ task.title }}</span>
            <button
              type="button"
              class="task-delete"
              :aria-label="`Удалить задачу: ${task.title}`"
              @click="deleteTaskById(task.id)"
            >
              Удалить
            </button>
          </li>
        </ul>
      </template>
    </section>
  </main>
</template>
