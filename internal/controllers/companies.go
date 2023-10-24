package controllers

import (
	"fmt"
	"strconv"

	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
)

type CompanyControllers struct {
	CompanyServices services.CompanyInterface
}

// PatchCompany Метод для обновления данных компании
func (comp CompanyControllers) PatchCompany(c *fiber.Ctx) error {
	var company services.Company
	c.BodyParser(&company)

	comp.CompanyServices = company
	err := comp.CompanyServices.Update()

	if err != nil {
		return fiber.NewError(500, "Ошибка при обновлении данных")
	}

	return fiber.NewError(200, "Данные успешно обновлены")
}

// CompanyId Метод для получения компании по id
func (comp CompanyControllers) CompanyId(c *fiber.Ctx) error {
	companyId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return fiber.NewError(500, "Неизвестная ошибка")
	}

	comp.CompanyServices = services.Company{CompanyId: companyId}
	company, err := comp.CompanyServices.ReadById()

	if err != nil {
		return fiber.NewError(500, "Неизвестная ошибка")
	}

	return c.JSON(company)
}

// AllCompanies Методя для получения компаний
func (comp CompanyControllers) AllCompanies(c *fiber.Ctx) error {
	searchQuery := c.Query("search_query")
	if searchQuery != "" {
		comp.CompanyServices = services.Company{}
		resultCompanies, err := comp.CompanyServices.SearchByQuery(searchQuery)

		if err != nil {
			return fiber.NewError(500, "Ошибка при поиске компаний")
		}

		return c.JSON(resultCompanies)
	}

	comp.CompanyServices = services.Company{}
	allCompanies, err := comp.CompanyServices.ReadAll()
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(404, "Ошибка при поиске компаний")
	}

	return c.JSON(allCompanies)
}

// DeleteCompany Метод для удаления компаний
func (comp CompanyControllers) DeleteCompany(c *fiber.Ctx) error {
	companyId, _ := strconv.Atoi(c.Query("company_id"))
	comp.CompanyServices = services.Company{CompanyId: companyId}

	err := comp.CompanyServices.Delete()
	if err != nil {
		return fiber.NewError(500, "Ошибка при удалении компании")
	}

	return fiber.NewError(200, "Компания успешно удалена")
}
