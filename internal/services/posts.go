package services

import (
	"context"
	"fmt"
	"github.com/Troom-Corp/troom/internal/storage"
)

type PostInterface interface {
	Create() error
	ReadAll() ([]Post, error)
	ReadById() (Post, error)
	Update() error
	Delete() error
}

type Post struct {
	PostId   int
	UserId   int
	Time     string
	Blocks   string
	Likes    string
	Dislikes string
}

// Create Создать пост по входным данным
func (p Post) Create() error {

	conn := storage.SqlInterface.New()

	createQuery := fmt.Sprintf("INSERT INTO public.posts (userid, time, blocks, likes, dislikes) VALUES ('%d', '%s', '%s', '%s', '%s');", p.UserId, p.Time, p.Blocks, p.Likes, p.Dislikes)
	_, err := conn.Query(context.Background(), createQuery)
	if err != nil {
		return err
	}

	storage.SqlInterface.Close(conn)
	return nil
}

// ReadAll Прочитать все посты из базы данных
func (p Post) ReadAll() ([]Post, error) {
	var posts []Post

	conn := storage.SqlInterface.New()

	rows, _ := conn.Query(context.Background(), "SELECT * FROM public.posts")

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostId, &post.UserId, &post.Time, &post.Blocks, &post.Likes, &post.Dislikes)
		if err != nil {
			return []Post{}, err
		}
		posts = append(posts, post)
	}
	storage.SqlInterface.Close(conn)
	return posts, nil
}

// ReadById Прочитать один пост по ID из базы данных
func (p Post) ReadById() (Post, error) {
	var post Post

	conn := storage.SqlInterface.New()

	readByIdQuery := fmt.Sprintf("SELECT * FROM public.posts WHERE postid=%d", p.PostId)
	err := conn.QueryRow(context.Background(), readByIdQuery).Scan(&post.PostId, &post.UserId, &post.Time, &post.Blocks, &post.Likes, &post.Dislikes)
	if err != nil {
		return Post{}, err
	}
	storage.SqlInterface.Close(conn)
	return post, nil
}

// Update Обновить данные поста по ID
func (p Post) Update() error {

	conn := storage.SqlInterface.New()

	updateByIdQuery := fmt.Sprintf("UPDATE public.posts SET blocks = '%s' WHERE postid = %d", p.Blocks, p.PostId)
	_, err := conn.Query(context.Background(), updateByIdQuery)
	if err != nil {
		return err
	}
	storage.SqlInterface.Close(conn)
	return nil
}

// Delete Удалить все данные поста по ID
func (p Post) Delete() error {

	conn := storage.SqlInterface.New()

	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.posts WHERE postid = %d", p.PostId)
	_, err := conn.Query(context.Background(), deleteByIdQuery)
	if err != nil {
		return err
	}
	storage.SqlInterface.Close(conn)
	return nil
}
