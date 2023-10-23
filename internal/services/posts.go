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
	conn := storage.SqlInterface.New()

	createQuery := fmt.Sprintf("INSERT INTO public.posts (userid, time, blocks, likes, dislikes) VALUES ('%d', '%s', '%s', '%s', '%s');", p.UserId, p.Time, p.Blocks, p.Likes, p.Dislikes)
	_, err := conn.Query(context.Background(), createQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return err
	}

	storage.SqlInterface.Close(conn)
	return nil
}

// ReadAll Прочитать все посты из базы данных
func (p Post) ReadAll() ([]Post, error) {
	var posts []Post
	conn := storage.SqlInterface.New()

	rows, err := conn.Query(context.Background(), "SELECT * FROM public.posts")

	if err != nil {
		storage.SqlInterface.Close(conn)
		return []Post{}, nil
	}

	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.PostId,
			&post.UserId,
			&post.Time,
			&post.Title,
			&post.Blocks,
			&post.Likes,
			&post.Dislikes)
		if err != nil {
			storage.SqlInterface.Close(conn)
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
		storage.SqlInterface.Close(conn)
		return Post{}, err
	}

	storage.SqlInterface.Close(conn)
	return post, nil
}

// SearchByQuery Поиск постов по названию
func (p Post) SearchByQuery(searchQuery string) ([]Post, error) {
	var posts []Post
	conn := storage.SqlInterface.New()

	searchFormat := "%" + searchQuery + "%"
	searchByQuery := fmt.Sprintf("SELECT * FROM public.posts WHERE LOWER(title) LIKE '%s'", searchFormat)
	rows, err := conn.Query(context.Background(), searchByQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return []Post{}, nil
	}

	for rows.Next() {
		var post Post
		err = rows.Scan(
			&post.PostId,
			&post.UserId,
			&post.Time,
			&post.Title,
			&post.Blocks,
			&post.Likes, &post.Dislikes)
		if err != nil {
			storage.SqlInterface.Close(conn)
			return []Post{}, err
		}
		posts = append(posts, post)
	}

	storage.SqlInterface.Close(conn)
	return posts, nil
}

// Update Обновить данные поста по ID
func (p Post) Update() error {
	conn := storage.SqlInterface.New()

	updateByIdQuery := fmt.Sprintf("UPDATE public.posts SET title = '%s' blocks = '%s', WHERE postid = %d", p.Title, p.Blocks, p.PostId)
	_, err := conn.Query(context.Background(), updateByIdQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
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
		storage.SqlInterface.Close(conn)
		return err
	}
	storage.SqlInterface.Close(conn)
	return nil
}
