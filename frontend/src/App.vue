<script setup lang="ts">
import { ref } from 'vue'

interface Task {
  id: number
  title: string
  completed: boolean
}

const tasks = ref<Task[]>([
  {
    id: 1,
    title: 'Подключить PostgreSQL',
    completed: true,
  },
  {
    id: 2,
    title: 'Собрать frontend shell',
    completed: false,
  },
  {
    id: 3,
    title: 'Подготовить загрузку задач из API',
    completed: false,
  },
])

const input = ref('')

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
        </li>
      </ul>
    </section>
  </main>
</template>
