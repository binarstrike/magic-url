package repository

const createUser = `INSERT INTO "user"
  (user_id,username,email,hashed_password) VALUES
  ($1,$2,$3,$4) RETURNING *`

const updateUser = `UPDATE "user" SET
  username = COALESCE(NULLIF($1, ''), username),
  updated_at = now()
  WHERE user_id = $2 RETURNING *`

const getUserById = `SELECT * FROM "user" WHERE user_id = $1 LIMIT 1`

const getUserByEmail = `SELECT * FROM "user" WHERE email = $1 LIMIT 1`

const deleteUser = `DELETE FROM "user" WHERE user_id = $1`
