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
