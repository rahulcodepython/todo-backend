package utils

const (
	// UserTableName defines the constant string for the name of the 'users' table in the database.
	// This constant ensures consistency when referencing the users table across different parts of the application,
	// reducing the risk of typos and making schema changes easier to manage.
	UserTableName = "users"
	// UserTableSchema defines the constant string representing the column names for the 'users' table.
	// This schema is used in SQL INSERT statements to specify the order and names of the columns being populated,
	// ensuring data integrity and correct mapping between application data structures and database columns.
	UserTableSchema = "id, name, email, image, password, jwt, created_at, updated_at"

	// JWTTableName defines the constant string for the name of the 'jwt_tokens' table in the database.
	// This table is typically used to store JSON Web Tokens, often for session management or blacklisting.
	JWTTableName = "jwt_tokens"
	// JWTTableSchema defines the constant string representing the column names for the 'jwt_tokens' table.
	// This schema specifies the structure of the JWT storage, including the token itself and its expiration time,
	// which is crucial for validating and managing user sessions.
	JWTTableSchema = "id, token, expires_at"

	TodoTableName   = "todos"
	TodoTableSchema = "id, title, completed, owner, created_at"
)
