package main

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type user struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// getUsers returns an empty slice if any user doesn't exist
func getUsers(db *sql.DB) ([]*user, error) {
	users := []*user{}
	rows, err := db.Query("select * from users")
	if err != nil {
		return users, err
	}

	for rows.Next() {
		u := new(user)
		err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func createUser(db *sql.DB, u *user) (*user, error) {
	var id string
	err := db.QueryRow("insert into users (name, email) values ($1, $2) returning id", u.Name, u.Email).Scan(&id)
	if err != nil {
		return nil, err
	}
	return getUser(db, id)
}

func getUser(db *sql.DB, id string) (*user, error) {
	row := db.QueryRow("select * from users where id = $1", id)
	u := new(user)
	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func updateUser(db *sql.DB, id string, u *user) (*user, error) {
	err := db.QueryRow(
		"update users set name = $1, email = $2, updated_at = $3 where id = $4 returning id",
		u.Name, u.Email, time.Now().Format(time.RFC3339Nano), id,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return getUser(db, id)
}

func deleteUser(db *sql.DB, id string) error {
	err := db.QueryRow("delete from users where id = $1 returning id", id).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
