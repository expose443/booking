package repository

import (
	"database/sql"
	"errors"

	"github.com/with-insomnia/Hotel/internal/model"
)

type PostQuery interface {
	CreatePost(post model.Post) error
	CreateAdminPost(post model.Post) error
	DeleteAdminPost(postId int64) error
	GetAllPost() ([]model.Post, error)
	GetAllAdminPost() ([]model.Post, error)
	GetPostById(postId int64) (model.Post, error)
	GetAdminPostById(postId int64) (model.Post, error)
}

type postQuery struct {
	db *sql.DB
}

func (p *postQuery) CreatePost(post model.Post) error {
	_, err := p.db.Exec("INSERT INTO posts ( message, user_id) VALUES($1, $2)", post.Content, post.Author.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *postQuery) CreateAdminPost(post model.Post) error {
	_, err := p.db.Exec("INSERT INTO waitposts ( message, user_id) VALUES($1, $2)", post.Content, post.Author.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *postQuery) GetAllPost() ([]model.Post, error) {
	rows, err := p.db.Query("SELECT * FROM posts")
	if err != nil {
		return []model.Post{}, err
	}
	defer rows.Close()
	var all []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Author.ID, &post.Content); err != nil {
			return []model.Post{}, err
		}
		all = append(all, post)
	}
	return all, nil
}

func (p *postQuery) GetAllAdminPost() ([]model.Post, error) {
	rows, err := p.db.Query("SELECT * FROM waitposts")
	if err != nil {
		return []model.Post{}, err
	}
	defer rows.Close()
	var all []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Author.ID, &post.Content); err != nil {
			return []model.Post{}, err
		}
		all = append(all, post)
	}
	return all, nil
}

func (p *postQuery) GetPostById(postId int64) (model.Post, error) {
	row := p.db.QueryRow("SELECT post_id, message, user_id FROM posts WHERE post_id = $1 ", postId)
	var post model.Post
	if err := row.Scan(&post.ID, &post.Content, &post.Author.ID); err != nil {
		return model.Post{}, err
	}
	return post, nil
}

func (p *postQuery) GetAdminPostById(postId int64) (model.Post, error) {
	row := p.db.QueryRow("SELECT post_id, message, user_id FROM waitposts WHERE post_id = $1 ", postId)
	var post model.Post
	if err := row.Scan(&post.ID, &post.Content, &post.Author.ID); err != nil {
		return model.Post{}, err
	}
	return post, nil
}

func (p *postQuery) DeleteAdminPost(postId int64) error {
	res, err := p.db.Exec("DELETE FROM waitposts WHERE post_id = $1", postId)
	if err != nil {
		return err
	}
	rowsaffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsaffected == 0 {
		return errors.New("post delete was fail")
	}
	return nil
}
