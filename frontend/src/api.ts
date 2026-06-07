import axios from 'axios'

export interface Task {
  id: number
  title: string
  completed: boolean
}

const apiClient = axios.create({
  baseURL: '/api',
})

export const getTasks = async (): Promise<Task[]> => {
  const { data } = await apiClient.get<Task[]>('/tasks')
  return data
}
