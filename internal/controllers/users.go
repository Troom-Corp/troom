package controllers

import (
	"fmt"
	"github.com/Troom-Corp/troom/internal/models"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strconv"
)

type UserControllers struct {
	UserServices store.InterfaceUser
}

func (u UserControllers) GetUserByLogin(c *fiber.Ctx) error {
	login := c.Params("login")

	user, err := u.UserServices.FindOne("login", login)

	if user.UserId == 0 {
		return c.Status(404).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "404",
				Message: "User not found",
			},
		})
	}

	if err != nil {
		return c.Status(500).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "500",
				Message: "An error while getting an user",
			},
		})
	}

	return c.JSON(models.HttpResponse{
		Error: models.Error{
			Status:  "200",
			Message: "User",
		},
		Data: user,
	})
}

func (u UserControllers) SearchByQuery(c *fiber.Ctx) error {
	queries := c.Queries()
	limit, _ := strconv.Atoi(queries["limit"])
	page, _ := strconv.Atoi(queries["page"])

	queryUser, err := u.UserServices.QuerySearch(queries["q"], limit, page)
	if err != nil {
		return c.Status(500).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "500",
				Message: "An error while getting users",
			},
		})
	}

	return c.JSON(models.HttpResponse{
		Error: models.Error{
			Status:  "200",
			Message: "Users",
		},
		Data: queryUser,
	})
}

func (u UserControllers) SetAvatar(c *fiber.Ctx) error {
	jwt := c.Cookies("token")
	ID, err := pkg.GetIdentity(jwt)

	avatar, err := c.FormFile("avatar")

	if err != nil {
		return c.Status(400).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "400",
				Message: "Bad request",
			},
		})
	}

	avatarInstance := models.Avatar{
		Image: avatar,
	}

	isValidAvatar := avatarInstance.Validate()

	if isValidAvatar.Extension != "" || isValidAvatar.Size != "" {
		return c.Status(409).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "409",
				Message: "Bad Data",
			},
			Data: isValidAvatar,
		})
	}

	imageUUID := pkg.GenerateUUID()
	imageExt := filepath.Ext(avatar.Filename)

	oldAvatar, err := u.UserServices.UploadAvatar(ID, imageUUID+imageExt)
	if err != nil {
		return c.Status(500).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "500",
				Message: "An error while uploading an avatar",
			},
		})
	}

	if oldAvatar != "" {
		os.Remove(fmt.Sprintf("./uploads/%s", oldAvatar))
	}

	err = c.SaveFile(avatar, fmt.Sprintf("./uploads/%s", imageUUID+imageExt))
	if err != nil {
		return c.Status(500).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "500",
				Message: "An error while saving an avatar",
			},
		})
	}

	return c.JSON(models.HttpResponse{
		Error: models.Error{
			Status:  "201",
			Message: "An avatar was uploaded successfully",
		},
	})
}

func (u UserControllers) SetLayout(c *fiber.Ctx) error {
	jwt := c.Cookies("token")
	ID, _ := pkg.GetIdentity(jwt)

	layout, err := c.FormFile("layout")

	if err != nil {
		return c.Status(400).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "400",
				Message: "Bad request",
			},
		})
	}

	layoutInstance := models.Avatar{
		Image: layout,
	}

	isValidAvatar := layoutInstance.Validate()

	if isValidAvatar.Extension != "" || isValidAvatar.Size != "" {
		return c.Status(400).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "400",
				Message: "Bad request",
			},
			Data: isValidAvatar,
		})
	}

	imageUUID := pkg.GenerateUUID()
	imageExt := filepath.Ext(layout.Filename)

	oldLayout, err := u.UserServices.UploadLayout(ID, imageUUID+imageExt)
	if err != nil {
		return c.Status(500).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "500",
				Message: "An error while uploading a layout",
			},
		})
	}
	if oldLayout != "" {
		os.Remove(fmt.Sprintf("./uploads/%s", oldLayout))
	}

	err = c.SaveFile(layout, fmt.Sprintf("./uploads/%s", imageUUID+imageExt))
	if err != nil {
		return c.Status(500).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "500",
				Message: "An error while saving a layout",
			},
		})
	}

	return c.JSON(models.HttpResponse{
		Error: models.Error{
			Status:  "201",
			Message: "Layout was uploaded successfully",
		},
	})
}

func (u UserControllers) DeleteAvatar(c *fiber.Ctx) error {
	jwt := c.Cookies("token")
	ID, _ := pkg.GetIdentity(jwt)

	deletedAvatar, err := u.UserServices.DeleteAvatar(ID)
	if err != nil {
		return c.Status(500).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "500",
				Message: "An error while deleting an avatar",
			},
		})
	}

	if deletedAvatar != "" {
		os.Remove(fmt.Sprintf("./uploads/%s", deletedAvatar))
	}
	return c.JSON(models.HttpResponse{
		Error: models.Error{
			Status:  "200",
			Message: "Avatar was deleted successfully",
		},
	})
}

func (u UserControllers) DeleteLayout(c *fiber.Ctx) error {
	jwt := c.Cookies("token")
	ID, _ := pkg.GetIdentity(jwt)

	deletedLayout, err := u.UserServices.DeleteLayout(ID)
	if err != nil {
		return c.Status(500).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "500",
				Message: "An error while deleting a layout",
			},
		})
	}
	if deletedLayout != "" {
		os.Remove(fmt.Sprintf("./uploads/%s", deletedLayout))
	}

	return c.JSON(models.HttpResponse{
		Error: models.Error{
			Status:  "200",
			Message: "Layout was deleted successfully",
		},
	})
}

func GetUserControllers(store store.InterfaceStore) UserControllers {
	return UserControllers{
		UserServices: store.Users(),
	}
}
