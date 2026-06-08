import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { createTask, getTasks, updateTask, type Task } from '../api'
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
})
