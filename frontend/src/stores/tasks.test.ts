import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { createTask, deleteTask, getTasks, updateTask, type Task } from '../api'
import { useTasksStore } from './tasks'

vi.mock('../api', () => ({
  createTask: vi.fn(),
  deleteTask: vi.fn(),
  getTasks: vi.fn(),
  updateTask: vi.fn(),
}))

describe('tasks store', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    setActivePinia(createPinia())
  })

  it('starts with the expected initial state', () => {
    const store = useTasksStore()

    expect(store.tasks).toEqual([])
    expect(store.input).toBe('')
    expect(store.isLoading).toBe(true)
    expect(store.isCreating).toBe(false)
    expect(store.updatingTaskId).toBeNull()
    expect(store.deletingTaskId).toBeNull()
    expect(store.editingTaskId).toBeNull()
    expect(store.editingTitle).toBe('')
    expect(store.errText).toBe('')
    expect(store.formErrText).toBe('')
    expect(store.actionErrText).toBe('')
  })

  it('loads tasks successfully', async () => {
    const store = useTasksStore()
    const fetchedTasks: Task[] = [
      { id: 1, title: 'Prepare README', completed: true },
      { id: 2, title: 'Check tasks API', completed: false },
    ]

    vi.mocked(getTasks).mockResolvedValue(fetchedTasks)

    expect(store.tasks).toEqual([])

    await store.loadTasks()

    expect(getTasks).toHaveBeenCalledTimes(1)
    expect(store.tasks).toEqual(fetchedTasks)
    expect(store.isLoading).toBe(false)
    expect(store.errText).toBe('')
  })

  it('keeps tasks empty and stores error text when loading fails', async () => {
    const store = useTasksStore()

    vi.mocked(getTasks).mockRejectedValue(new Error('network error'))

    await store.loadTasks()

    expect(getTasks).toHaveBeenCalledTimes(1)
    expect(store.tasks).toEqual([])
    expect(store.isLoading).toBe(false)
    expect(store.errText).not.toBe('')
  })

  it('creates a task successfully', async () => {
    const store = useTasksStore()
    const createdTask: Task = { id: 1, title: 'Write frontend tests', completed: false }

    store.input = '  Write frontend tests  '
    vi.mocked(createTask).mockResolvedValue(createdTask)

    await store.createTask()

    expect(createTask).toHaveBeenCalledTimes(1)
    expect(createTask).toHaveBeenCalledWith({ title: 'Write frontend tests' })
    expect(store.tasks).toEqual([createdTask])
    expect(store.input).toBe('')
    expect(store.isCreating).toBe(false)
    expect(store.formErrText).toBe('')
  })

  it('keeps input and stores form error when task creation fails', async () => {
    const store = useTasksStore()

    store.input = 'Write frontend tests'
    vi.mocked(createTask).mockRejectedValue(new Error('network error'))

    await store.createTask()

    expect(createTask).toHaveBeenCalledTimes(1)
    expect(store.tasks).toEqual([])
    expect(store.input).toBe('Write frontend tests')
    expect(store.isCreating).toBe(false)
    expect(store.formErrText).not.toBe('')
  })

  it('toggles task completed successfully', async () => {
    const store = useTasksStore()
    const task: Task = { id: 1, title: 'Check tasks API', completed: false }
    const updatedTask: Task = { ...task, completed: true }

    store.tasks = [task]
    vi.mocked(updateTask).mockResolvedValue(updatedTask)

    await store.toggleTaskCompleted(task.id)

    expect(updateTask).toHaveBeenCalledTimes(1)
    expect(updateTask).toHaveBeenCalledWith(task.id, { completed: true })
    expect(store.tasks).toEqual([updatedTask])
    expect(store.updatingTaskId).toBeNull()
    expect(store.actionErrText).toBe('')
  })

  it('updates task title successfully', async () => {
    const store = useTasksStore()
    const task: Task = { id: 1, title: 'Prepare README', completed: false }
    const updatedTask: Task = { ...task, title: 'Update README' }

    store.tasks = [task]
    vi.mocked(updateTask).mockResolvedValue(updatedTask)

    store.startEditingTask(task.id)

    expect(store.editingTaskId).toBe(task.id)
    expect(store.editingTitle).toBe(task.title)

    store.editingTitle = '  Update README  '

    await store.updateTaskTitle()

    expect(updateTask).toHaveBeenCalledTimes(1)
    expect(updateTask).toHaveBeenCalledWith(task.id, { title: 'Update README' })
    expect(store.tasks).toEqual([updatedTask])
    expect(store.editingTaskId).toBeNull()
    expect(store.editingTitle).toBe('')
    expect(store.updatingTaskId).toBeNull()
    expect(store.actionErrText).toBe('')
  })

  it('keeps task draft and stores action error when title update fails', async () => {
    const store = useTasksStore()
    const task: Task = { id: 1, title: 'Prepare README', completed: false }

    store.tasks = [task]
    vi.mocked(updateTask).mockRejectedValue(new Error('network error'))

    store.startEditingTask(task.id)
    store.editingTitle = 'Update README'

    await store.updateTaskTitle()

    expect(updateTask).toHaveBeenCalledTimes(1)
    expect(updateTask).toHaveBeenCalledWith(task.id, { title: 'Update README' })
    expect(store.tasks).toEqual([task])
    expect(store.editingTaskId).toBe(task.id)
    expect(store.editingTitle).toBe('Update README')
    expect(store.actionErrText).not.toBe('')
    expect(store.updatingTaskId).toBeNull()
  })

  it('does not call API and keeps editing draft when title is blank', async () => {
    const store = useTasksStore()
    const task: Task = { id: 1, title: 'Prepare README', completed: false }

    store.tasks = [task]

    store.startEditingTask(task.id)
    store.editingTitle = '   '

    await store.updateTaskTitle()

    expect(updateTask).not.toHaveBeenCalled()
    expect(store.tasks).toEqual([task])
    expect(store.editingTaskId).toBe(task.id)
    expect(store.editingTitle).toBe('   ')
    expect(store.actionErrText).not.toBe('')
    expect(store.updatingTaskId).toBeNull()
  })

  it('keeps task unchanged and stores action error when toggle fails', async () => {
    const store = useTasksStore()
    const task: Task = { id: 1, title: 'Check tasks API', completed: false }

    store.tasks = [task]
    vi.mocked(updateTask).mockRejectedValue(new Error('network error'))

    await store.toggleTaskCompleted(task.id)

    expect(updateTask).toHaveBeenCalledTimes(1)
    expect(updateTask).toHaveBeenCalledWith(task.id, { completed: true })
    expect(store.tasks).toEqual([task])
    expect(store.updatingTaskId).toBeNull()
    expect(store.actionErrText).not.toBe('')
  })

  it('deletes a task successfully', async () => {
    const store = useTasksStore()
    const firstTask: Task = { id: 1, title: 'Prepare README', completed: true }
    const secondTask: Task = { id: 2, title: 'Check tasks API', completed: false }

    store.tasks = [firstTask, secondTask]
    vi.mocked(deleteTask).mockResolvedValue()

    await store.deleteTaskById(firstTask.id)

    expect(deleteTask).toHaveBeenCalledTimes(1)
    expect(deleteTask).toHaveBeenCalledWith(firstTask.id)
    expect(store.tasks).toEqual([secondTask])
    expect(store.deletingTaskId).toBeNull()
    expect(store.actionErrText).toBe('')
  })

  it('keeps task unchanged and stores action error when delete fails', async () => {
    const store = useTasksStore()
    const task: Task = { id: 1, title: 'Prepare README', completed: false }

    store.tasks = [task]
    vi.mocked(deleteTask).mockRejectedValue(new Error('network error'))

    await store.deleteTaskById(task.id)

    expect(deleteTask).toHaveBeenCalledTimes(1)
    expect(deleteTask).toHaveBeenCalledWith(task.id)
    expect(store.tasks).toEqual([task])
    expect(store.deletingTaskId).toBeNull()
    expect(store.actionErrText).not.toBe('')
  })

  it('starts editing a task with a title draft', () => {
    const store = useTasksStore()
    const task: Task = { id: 1, title: 'Prepare README', completed: false }

    store.tasks = [task]

    store.startEditingTask(task.id)

    expect(store.editingTaskId).toBe(task.id)
    expect(store.editingTitle).toBe(task.title)
  })

  it('cancels task editing', () => {
    const store = useTasksStore()

    store.editingTaskId = 1
    store.editingTitle = 'Prepare README'
    store.actionErrText = 'Название не может быть пустым'

    store.cancelEditingTask()

    expect(store.editingTaskId).toBeNull()
    expect(store.editingTitle).toBe('')
    expect(store.actionErrText).toBe('')
  })
})
