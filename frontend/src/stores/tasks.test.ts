import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useTasksStore } from './tasks'

describe('tasks store', () => {
  beforeEach(() => {
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
})
