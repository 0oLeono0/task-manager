package main

type task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type createTaskInput struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type updateTaskInput struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}
