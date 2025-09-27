package users

const (
	CreateUserQuery = `
		INSERT INTO users (id, name, email, password, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;`

	GetUserByEmailQuery = `SELECT id, name, email, password, created_at FROM users WHERE email = $1;`

	GetUserByIDQuery = `SELECT id, name, email, created_at FROM users WHERE id = $1;`

	CreateJWTQuery = `
		INSERT INTO jwt_tokens (id, user_id, token, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5);`

	GetJWTByTokenQuery = `SELECT user_id, expires_at FROM jwt_tokens WHERE token = $1;`

	DeleteJWTByTokenQuery = `DELETE FROM jwt_tokens WHERE token = $1;`
)
