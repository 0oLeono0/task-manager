import axios from 'axios'

export interface Task {
  id: number
  title: string
  completed: boolean
}

export interface CreateTaskInput {
  title: string
}

export interface UpdateTaskInput {
  title?: string
  completed?: boolean
}

const apiClient = axios.create({
  baseURL: '/api',
})

export const getTasks = async (): Promise<Task[]> => {
  const { data } = await apiClient.get<Task[]>('/tasks')
  return data
}

export const createTask = async (input: CreateTaskInput): Promise<Task> => {
  const { data } = await apiClient.post<Task>('/tasks', input)
  return data
}

export const updateTask = async (id: number, input: UpdateTaskInput): Promise<Task> => {
  const { data } = await apiClient.patch<Task>(`/tasks/${id}`, input)
  return data
}
