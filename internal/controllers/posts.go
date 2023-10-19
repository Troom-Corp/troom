package controllers

import (
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type PostsControllers struct {
	PostServices services.PostInterface
}

func (p PostsControllers) GetPost(c *fiber.Ctx) error {
	postId, _ := strconv.Atoi(c.Query("post_id"))
	p.PostServices = services.Post{PostId: postId}
	if postId == 0 {
		allPosts, err := p.PostServices.ReadAll()
		if err != nil {
			return fiber.NewError(404, "Ошибка при получении постов")
		}
		return c.JSON(allPosts)
	}
	post, err := p.PostServices.ReadById()
	if err != nil {
		return fiber.NewError(404, "Пост не найден")
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
