package service

import (
	"errors"
	"strings"

	"github.com/with-insomnia/Hotel/internal/model"
	"github.com/with-insomnia/Hotel/internal/repository"
)

type PostService interface {
	GetAllPosts() ([]model.Post, error)
	GetAllAdminPost() ([]model.Post, error)
	DeleteAdminPost(postId int64) error
	CreatePost(post *model.Post) (int, error)
	CreateAdminPost(post *model.Post) (int, error)
	GetAdminPostById(postId int64) (model.Post, error)
}

type postService struct {
	repository repository.PostQuery
}

func NewPostService(dao repository.DAO) PostService {
	return &postService{
		dao.NewPostQuery(),
	}
}

// 400 - http status Bad request
// 500 - http status Internal server error
// 200 - http status Ok

func (p *postService) GetAllPosts() ([]model.Post, error) {
	posts, err := p.repository.GetAllPost()
	if err != nil {
		return nil, err
	}
	result := []model.Post{}
	for i := len(posts) - 1; i >= 0; i-- {
		result = append(result, posts[i])
	}
	return result, nil
}

func (p *postService) GetAllAdminPost() ([]model.Post, error) {
	posts, err := p.repository.GetAllAdminPost()
	if err != nil {
		return nil, err
	}
	result := []model.Post{}
	for i := len(posts) - 1; i >= 0; i-- {
		result = append(result, posts[i])
	}
	return result, nil
}

func (p *postService) CreatePost(post *model.Post) (int, error) {
	if ok := validDataString(post.Content); !ok {
		return 400, errors.New("content is invalid")
	}
	err := p.repository.CreatePost(*post)
	if err != nil {
		return 500, errors.New("create post was failed")
	}

	return 200, nil
}

func (p *postService) CreateAdminPost(post *model.Post) (int, error) {
	if ok := validDataString(post.Content); !ok {
		return 400, errors.New("content is invalid")
	}
	err := p.repository.CreateAdminPost(*post)
	if err != nil {
		return 500, errors.New("create post was failed")
	}

	return 200, nil
}

func (p *postService) DeleteAdminPost(post int64) error {
	return p.repository.DeleteAdminPost(post)
}

func (p *postService) GetAdminPostById(postId int64) (model.Post, error) {
	return p.repository.GetAdminPostById(postId)
}

func validDataString(s string) bool {
	str := strings.TrimSpace(s)
	if len(str) == 0 {
		return false
	}
	for _, v := range str {
		if v < rune(32) {
			return false
		}
	}
	return true
}
