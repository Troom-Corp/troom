package services

import (
	"context"
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
	conn := storage.SqlInterface.New()

	createQuery := fmt.Sprintf("INSERT INTO public.comments (postid, userid, text, likes, replies) VALUES (%d, %d, '%s', '%s', '%s')", c.PostId, c.UserId, c.Text, c.Likes, c.Replies)
	_, err := conn.Query(context.Background(), createQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return err
	}

	storage.SqlInterface.Close(conn)
	return nil
}

// ReadByPostId Прочитать все коментарии
func (c Comment) ReadByPostId() ([]Comment, error) {
	var comments []Comment
	conn := storage.SqlInterface.New()

	byPostIdQuery := fmt.Sprintf("SELECT * FROM public.comments WHERE postid=%d", c.PostId)
	rows, _ := conn.Query(context.Background(), byPostIdQuery)
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.CommentId,
			&comment.PostId,
			&comment.UserId,
			&comment.Text,
			&comment.Likes,
			&comment.Replies)

		if err != nil {
			storage.SqlInterface.Close(conn)
			return []Comment{}, err
		}
		comments = append(comments, comment)
	}

	storage.SqlInterface.Close(conn)
	return comments, nil
}

// Update Обновить данные коментария по ID
func (c Comment) Update() error {
	conn := storage.SqlInterface.New()

	updateByIdQuery := fmt.Sprintf("UPDATE public.comments SET text = '%s', likes = '%s', replies = '%s' WHERE commentid=%d", c.Text, c.Likes, c.Replies, c.CommentId)
	_, err := conn.Query(context.Background(), updateByIdQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return err
	}

	storage.SqlInterface.Close(conn)
	return nil
}

// Delete Удалить комментарий по ID
func (c Comment) Delete() error {
	conn := storage.SqlInterface.New()

	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.comments WHERE commentid = %d", c.CommentId)
	_, err := conn.Query(context.Background(), deleteByIdQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return err
	}

	storage.SqlInterface.Close(conn)
	return nil
}
