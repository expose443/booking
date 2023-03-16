package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/with-insomnia/Hotel/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func NewDB(cfg config.Database) (*sql.DB, error) {
	postgresDbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBhost, cfg.DBport, cfg.DBuser, cfg.DBpass, cfg.DBname)

	db, err := sql.Open("postgres", postgresDbInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = createTable(db)
	if err != nil {
		return nil, err
	}
	err = createRooms(db)
	if err != nil {
		return nil, err
	}
	err = createAdmin(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTable(db *sql.DB) error {
	query := []string{}

	users := `CREATE TABLE IF NOT EXISTS users(
		user_id SERIAL,
		email VARCHAR NOT NULL,
		first_name VARCHAR NOT NULL,
		last_name VARCHAR NOT NULL,
		password VARCHAR NOT NULL,
		role VARCHAR NOT NULL,
		number VARCHAR NOT NULL
	);
	`
	posts := `
	CREATE TABLE IF NOT EXISTS posts(
		post_id SERIAL,
		user_id SERIAL,
		message TEXT NOT NULL
	);
	`
	waitposts := `
	CREATE TABLE IF NOT EXISTS waitposts(
		post_id SERIAL,
		user_id SERIAL,
		message TEXT NOT NULL
	);
	`
	session := `
	CREATE TABLE IF NOT EXISTS sessions(
		user_id SERIAL,
		token TEXT NOT NULL,
		expiry TIMESTAMP NOT NULL
	)
	`
	rooms := `CREATE TABLE IF NOT EXISTS rooms(
		room_id SERIAL NOT NULL,
		room_name VARCHAR NOT NULL,
		room_occupied SERIAL
	);
	`

	waitlist := `CREATE TABLE IF NOT EXISTS waitlist(
		user_id SERIAL,
		room_id SERIAL,
		check_in DATE NOT NULL,
		check_out DATE NOT NULL
	)`

	query = append(query, users, posts, session, rooms, waitlist, waitposts)
	for _, v := range query {
		_, err := db.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func createRooms(db *sql.DB) error {
	createRoom := `INSERT INTO rooms(room_name, room_occupied) values($1, $2)`
	stmt, err := db.Prepare(createRoom)
	defer stmt.Close()
	_, err = stmt.Exec("generalsolo", 0)
	if err != nil {
		return err
	}
	_, err = stmt.Exec("generalduo", 0)
	if err != nil {
		return err
	}
	_, err = stmt.Exec("luxduo", 0)
	if err != nil {
		return err
	}
	_, err = stmt.Exec("luxtrio", 0)
	if err != nil {
		return err
	}
	return nil
}

func createAdmin(db *sql.DB) error {
	query := `INSERT INTO users(email, password, role, first_name, last_name, number) VALUES($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(query, "tanatar-abdurahman@mail.ru", Hasher("Abdu1234"), "admin", "abdurahman", "tanatar", "+12345678")
	if err != nil {
		return err
	}
	return nil
}

func Hasher(password string) string {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	fmt.Println(passwordHash)
	return string(passwordHash)
}
