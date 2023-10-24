package controllers

import (
	"fmt"
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type PostsControllers struct {
	PostServices services.PostInterface
}

func (p PostsControllers) PostId(c *fiber.Ctx) error {
	postId, err := strconv.Atoi(c.Params(":id"))
	if err != nil {
		return fiber.NewError(500, "Неизвестная ошибка")
	}

	p.PostServices = services.Post{PostId: postId}
	post, err := p.PostServices.ReadById()

	if err != nil {
		return fiber.NewError(500, "Неизвестная ошибка")
	}

	return c.JSON(post)
}

func (p PostsControllers) AllPost(c *fiber.Ctx) error {
	p.PostServices = services.Post{}
	post, err := p.PostServices.ReadAll()
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(500, "Неизвестная ошибка")
	}
	return c.JSON(post)
}

func (p PostsControllers) DeletePost(c *fiber.Ctx) error {
	postId, _ := strconv.Atoi(c.Query("post_id"))
	p.PostServices = services.Post{PostId: postId}

	err := p.PostServices.Delete()
	if err != nil {
		return fiber.NewError(500, "Ошибка при удалении поста")
	}
	return fiber.NewError(200, "Пост успешно удален")
}

func (p PostsControllers) PatchPost(c *fiber.Ctx) error {
	var newPost services.Post
	c.BodyParser(&newPost)

	err := newPost.Update()
	if err != nil {
		return fiber.NewError(500, "Ошибка при обновлении поста")
	}
	return fiber.NewError(200, "Пост успешно обновлен")
}
