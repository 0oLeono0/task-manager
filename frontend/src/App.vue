<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getTasks, type Task } from './api'

const tasks = ref<Task[]>([])
const input = ref('')
const isLoading = ref(true)
const errText = ref('')

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

const onSubmit = () => {
  const trimmedTitle = input.value.trim()
  if (trimmedTitle === '') return

  const task = { id: tasks.value.length + 1, title: trimmedTitle, completed: false }
  tasks.value.push(task)
  input.value = ''
}

const toggleTaskCompleted = (id: number) => {
  const task = tasks.value.find((task) => task.id === id)
  if (!task) return

  task.completed = !task.completed
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

      <p v-if="isLoading" class="task-message">Загружаем задачи...</p>
      <p v-else-if="errText" class="task-message task-message--error">{{ errText }}</p>

      <ul v-else class="task-list" aria-label="Список задач">
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
        </li>
      </ul>
    </section>
  </main>
</template>
