package main

import "database/sql"

type postgresTaskStore struct {
	db *sql.DB
}

func (p *postgresTaskStore) getTasks() ([]task, error) {
	rows, err := p.db.Query(`SELECT id, title, completed FROM tasks ORDER BY id`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []task

	for rows.Next() {
		var task task
		err = rows.Scan(&task.ID, &task.Title, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (p *postgresTaskStore) getTaskByID(id int) (task, error) {
	var task task
	row := p.db.QueryRow(`SELECT id, title, completed FROM tasks WHERE id = $1`, id)

	err := row.Scan(&task.ID, &task.Title, &task.Completed)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (p *postgresTaskStore) createTask(title string, completed bool) (task, error) {
	var task task
	row := p.db.QueryRow(`INSERT INTO tasks (title, completed) 
	VALUES ($1, $2)
	RETURNING id, title, completed`, title, completed)

	err := row.Scan(&task.ID, &task.Title, &task.Completed)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (p *postgresTaskStore) updateTask(id int, title *string, completed *bool) (task, error) {
	task, err := p.getTaskByID(id)
	if err != nil {
		return task, err
	}

	if title != nil {
		task.Title = *title
	}
	if completed != nil {
		task.Completed = *completed
	}

	row := p.db.QueryRow(`UPDATE tasks 
	SET title = $1, 
	completed = $2,
	updated_at = now() WHERE id = $3
	RETURNING id, title, completed`, task.Title, task.Completed, id)

	err = row.Scan(&task.ID, &task.Title, &task.Completed)
	if err != nil {
		return task, err
	}
	return task, nil
}
