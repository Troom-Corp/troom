package services

import (
	"fmt"

	"github.com/Troom-Corp/troom/internal/storage"
)

type CommentInterface interface {
	Create() error
	ReadByPostId() ([]Comment, error)
	Update() error
	Delete() error
}

type Comment struct {
	CommentId int
	PostId    int
	UserId    int
	Text      string
	Likes     string
	Replies   string
}

// Create Создать комментарий по входным данным
func (c Comment) Create() error {
	conn, err := storage.Sql.Open()

	createQuery := fmt.Sprintf("INSERT INTO public.comments (postid, userid, text, likes, replies) VALUES (%d, %d, '%s', '%s', '%s')", c.PostId, c.UserId, c.Text, c.Likes, c.Replies)
	_, err = conn.Query(createQuery)

	if err != nil {
		conn.Close()
		return err
	}

	conn.Close()
	return nil
}

// ReadByPostId Прочитать все коментарии
func (c Comment) ReadByPostId() ([]Comment, error) {
	var comments []Comment
	conn, err := storage.Sql.Open()

	if err != nil {
		return []Comment{}, err
	}

	byPostIdQuery := fmt.Sprintf("SELECT * FROM public.comments WHERE postid=%d", c.PostId)
	err = conn.Select(&comments, byPostIdQuery)

	conn.Close()
	return comments, nil
}

// Update Обновить данные коментария по ID
func (c Comment) Update() error {
	conn, err := storage.Sql.Open()

	updateByIdQuery := fmt.Sprintf("UPDATE public.comments SET text = '%s', likes = '%s', replies = '%s' WHERE commentid=%d", c.Text, c.Likes, c.Replies, c.CommentId)
	_, err = conn.Query(updateByIdQuery)

	if err != nil {
		conn.Close()
		return err
	}

	conn.Close()
	return nil
}

// Delete Удалить комментарий по ID
func (c Comment) Delete() error {
	conn, err := storage.Sql.Open()

	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.comments WHERE commentid = %d", c.CommentId)
	_, err = conn.Query(deleteByIdQuery)

	if err != nil {
		conn.Close()
		return err
	}

	conn.Close()
	return nil
}
