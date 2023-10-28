package services

import (
	"fmt"
	"strings"

	"github.com/Troom-Corp/troom/internal/storage"
)

type PostInterface interface {
	Create() error
	ReadAll() ([]Post, error)
	ReadById() (Post, error)
	SearchByQuery(searchQuery string) ([]Post, error)
	Update() error
	Delete() error
}

type Post struct {
	PostId   int
	UserId   int
	Time     string
	Title    string
	Blocks   string
	Likes    string
	Dislikes string
}

// Create Создать пост по входным данным
func (p Post) Create() error {
	conn, err := storage.Sql.Open()

	createQuery := fmt.Sprintf("INSERT INTO public.posts (userid, time, blocks, likes, dislikes) VALUES ('%d', '%s', '%s', '%s', '%s');", p.UserId, p.Time, p.Blocks, p.Likes, p.Dislikes)
	_, err = conn.Query(createQuery)

	if err != nil {
		storage.Sql.Close()
		return err
	}

	conn.Close()
	return nil
}

// ReadAll Прочитать все посты из базы данных
func (p Post) ReadAll() ([]Post, error) {
	var posts []Post
	conn, err := storage.Sql.Open()

	conn.Select(&posts, "SELECT * FROM public.posts")

	if err != nil {
		conn.Close()
		return []Post{}, nil
	}

	conn.Close()
	return posts, nil
}

// ReadById Прочитать один пост по ID из базы данных
func (p Post) ReadById() (Post, error) {
	var post Post
	conn, err := storage.Sql.Open()

	readByIdQuery := fmt.Sprintf("SELECT * FROM public.posts WHERE postid=%d", p.PostId)
	err = conn.Get(&post, readByIdQuery)

	if err != nil {
		conn.Close()
		return Post{}, err
	}

	conn.Close()
	return post, nil
}

// SearchByQuery Поиск постов по названию
func (p Post) SearchByQuery(searchQuery string) ([]Post, error) {
	var posts []Post
	conn, err := storage.Sql.Open()

	searchFormat := "%" + strings.ToLower(searchQuery) + "%"
	searchByQuery := fmt.Sprintf("SELECT * FROM public.posts WHERE LOWER(title) LIKE '%s'", searchFormat)
	err = conn.Select(&posts, searchByQuery)

	if err != nil {
		conn.Close()
		return []Post{}, nil
	}

	conn.Close()
	return posts, nil
}

// Update Обновить данные поста по ID
func (p Post) Update() error {
	conn, err := storage.Sql.Open()

	updateByIdQuery := fmt.Sprintf("UPDATE public.posts SET title = '%s' blocks = '%s', WHERE postid = %d", p.Title, p.Blocks, p.PostId)
	_, err = conn.Query(updateByIdQuery)

	if err != nil {
		conn.Close()
		return err
	}

	conn.Close()
	return nil
}

// Delete Удалить все данные поста по ID
func (p Post) Delete() error {
	conn, err := storage.Sql.Open()

	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.posts WHERE postid = %d", p.PostId)
	_, err = conn.Query(deleteByIdQuery)

	if err != nil {
		conn.Close()
		return err
	}
	conn.Close()
	return nil
}
