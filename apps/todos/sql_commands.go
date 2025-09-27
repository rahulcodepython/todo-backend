package todos

const (
	CreateTodoQuery = `
		INSERT INTO todos (id, user_id, todo, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;`

	CountTodosByUserIDQuery = `SELECT COUNT(*) FROM todos WHERE user_id = $1;`

	GetTodosByUserIDQuery = `
		SELECT id, todo, completed, created_at, updated_at
		FROM todos WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;`

	GetTodosByUserIDFilteredComplitionQuery = `
		SELECT id, todo, completed, created_at, updated_at
		FROM todos WHERE user_id = $1 AND Completed = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4;`

	GetTodoByIDAndUserIDQuery = `
		SELECT id, todo, completed, created_at, updated_at
		FROM todos WHERE id = $1 AND user_id = $2;`

	UpdateTodoQuery = `
		UPDATE todos SET todo = $1, updated_at = $2
		WHERE id = $3 AND user_id = $4
		RETURNING id, todo, completed, created_at, updated_at;`

	UpdateTodoStatusQuery = `
		UPDATE todos SET completed = NOT completed, updated_at = $1
		WHERE id = $2 AND user_id = $3;`

	DeleteTodoQuery = `DELETE FROM todos WHERE id = $1 AND user_id = $2;`
)
