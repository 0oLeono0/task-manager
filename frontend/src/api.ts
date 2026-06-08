import axios from 'axios'

export interface Task {
  id: number
  title: string
  completed: boolean
}

export interface CreateTaskInput {
  title: string
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
