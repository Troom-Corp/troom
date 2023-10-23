package controllers

import (
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type CommentControllers struct {
	CommentServices services.CommentInterface
}

func (cc CommentControllers) CreateComment(c *fiber.Ctx) error {
	var comment services.Comment
	c.BodyParser(&comment)

	cc.CommentServices = comment

	err := cc.CommentServices.Create()
	if err != nil {
		return fiber.NewError(500, "Ошибка при создании комментария")
	}

	return fiber.NewError(200, "Комментарий успешно создан")
}

func (cc CommentControllers) GetComment(c *fiber.Ctx) error {
	postId, _ := strconv.Atoi(c.Query("post_id"))
	if postId != 0 {
		cc.CommentServices = services.Comment{PostId: postId}
		commentsByPostId, err := cc.CommentServices.ReadByPostId()
		if err != nil {
			return fiber.NewError(404, "Комментарии отсутсвуют")
		}
		return c.JSON(commentsByPostId)
	}
	return c.JSON([]services.Comment{})
}

func (cc CommentControllers) DeleteComment(c *fiber.Ctx) error {
	commentId, _ := strconv.Atoi(c.Query("comment_id"))
	cc.CommentServices = services.Comment{CommentId: commentId}

	err := cc.CommentServices.Delete()
	if err != nil {
		return fiber.NewError(500, "Ошибка при удалении комментария")
	}

	return fiber.NewError(200, "Комментарий успешно удален")
}

func (cc CommentControllers) PatchComment(c *fiber.Ctx) error {
	var newComment services.Comment

	c.BodyParser(&newComment)

	cc.CommentServices = newComment

	err := cc.CommentServices.Update()
	if err != nil {
		return fiber.NewError(500, "Ошибка при редактировании комментария")
	}

	return fiber.NewError(200, "Комментарий успешно отредоктирован")
}
