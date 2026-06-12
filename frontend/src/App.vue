<script setup lang="ts">
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useTasksStore } from './stores/tasks'

const tasksStore = useTasksStore()
const {
  tasks,
  input,
  isLoading,
  isCreating,
  updatingTaskId,
  deletingTaskId,
  editingTaskId,
  editingTitle,
  errText,
  formErrText,
  actionErrText,
} = storeToRefs(tasksStore)

const {
  loadTasks,
  createTask: onSubmit,
  toggleTaskCompleted,
  deleteTaskById,
  startEditingTask,
  cancelEditingTask,
  updateTaskTitle,
} = tasksStore

onMounted(loadTasks)
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
        <button type="submit" :disabled="isCreating">Добавить</button>
      </form>
      <p v-if="formErrText" class="task-message task-message--error">{{ formErrText }}</p>

      <p v-if="isLoading" class="task-message">Загружаем задачи...</p>
      <p v-else-if="errText" class="task-message task-message--error">{{ errText }}</p>

      <template v-else>
        <p v-if="actionErrText" class="task-message task-message--error">{{ actionErrText }}</p>

        <p v-if="tasks.length === 0" class="task-message">Задач пока нет</p>

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
              :disabled="updatingTaskId !== null"
              :aria-label="
                task.completed
                  ? `Отметить невыполненной: ${task.title}`
                  : `Отметить выполненной: ${task.title}`
              "
              @click="toggleTaskCompleted(task.id)"
            ></button>
            <input
              v-if="editingTaskId === task.id"
              v-model="editingTitle"
              class="task-title-input"
              type="text"
              :aria-label="`Редактировать задачу: ${task.title}`"
            />
            <span v-else class="task-title">{{ task.title }}</span>
            <button
              v-if="editingTaskId === task.id"
              type="button"
              class="task-save"
              :disabled="updatingTaskId !== null"
              @click="updateTaskTitle"
            >
              Сохранить
            </button>
            <button
              v-if="editingTaskId === task.id"
              type="button"
              class="task-edit"
              :disabled="updatingTaskId !== null"
              @click="cancelEditingTask"
            >
              Отмена
            </button>
            <button
              v-else
              type="button"
              class="task-edit"
              :disabled="updatingTaskId !== null"
              :aria-label="`Редактировать задачу: ${task.title}`"
              @click="startEditingTask(task.id)"
            >
              Редактировать
            </button>
            <button
              type="button"
              class="task-delete"
              :disabled="deletingTaskId !== null || updatingTaskId !== null"
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
