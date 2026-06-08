import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { getTasks, type Task } from '../api'
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
})
