<template>
  <main class="app-shell">
    <section class="task-panel" aria-labelledby="app-title">
      <header class="app-header">
        <p class="app-kicker">Менеджер задач</p>
        <h1 id="app-title">Мои задачи</h1>
      </header>

      <form class="task-form" @submit.prevent>
        <label class="sr-only" for="task-title">Новая задача</label>
        <input id="task-title" type="text" placeholder="Например: проверить API" />
        <button type="submit">Добавить</button>
      </form>

      <ul class="task-list" aria-label="Список задач">
        <li class="task-item task-item--done">
          <span class="task-status" aria-hidden="true"></span>
          <span class="task-title">Подключить PostgreSQL</span>
        </li>
        <li class="task-item">
          <span class="task-status" aria-hidden="true"></span>
          <span class="task-title">Собрать frontend shell</span>
        </li>
        <li class="task-item">
          <span class="task-status" aria-hidden="true"></span>
          <span class="task-title">Подготовить загрузку задач из API</span>
        </li>
      </ul>
    </section>
  </main>
</template>
