package services

import (
	"context"
	"fmt"

	"github.com/Troom-Corp/troom/internal"
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
	createQuery := fmt.Sprintf("INSERT INTO public.posts (userid, time, blocks, likes, dislikes) VALUES ('%d', '%s', '%s', '%s', '%s');", p.UserId, p.Time, p.Blocks, p.Likes, p.Dislikes)
	_, err := internal.Store.Query(context.Background(), createQuery)
	if err != nil {
		return err
	}

	return nil
}

// ReadAll Прочитать все посты из базы данных
func (p Post) ReadAll() ([]Post, error) {
	var posts []Post
	rows, _ := internal.Store.Query(context.Background(), "SELECT * FROM public.posts")

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostId, &post.UserId, &post.Time, &post.Blocks, &post.Likes, &post.Dislikes)
		if err != nil {
			return []Post{}, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// ReadById Прочитать один пост по ID из базы данных
func (p Post) ReadById() (Post, error) {
	var post Post
	readByIdQuery := fmt.Sprintf("SELECT * FROM public.posts WHERE postid=%d", p.PostId)
	err := internal.Store.QueryRow(context.Background(), readByIdQuery).Scan(&post.PostId, &post.UserId, &post.Time, &post.Blocks, &post.Likes, &post.Dislikes)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

// Update Обновить данные поста по ID
func (p Post) Update() error {
	updateByIdQuery := fmt.Sprintf("UPDATE public.posts SET blocks = '%s' WHERE postid = %d", p.Blocks, p.PostId)
	_, err := internal.Store.Query(context.Background(), updateByIdQuery)
	if err != nil {
		return err
	}
	return nil
}

// Delete Удалить все данные поста по ID
func (p Post) Delete() error {
	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.posts WHERE postid = %d", p.PostId)
	_, err := internal.Store.Query(context.Background(), deleteByIdQuery)
	if err != nil {
		return err
	}
	return nil
}
