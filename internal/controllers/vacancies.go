package controllers

import (
	"strconv"

	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
)

type VacancyControllers struct {
	VacancyServices services.VacancyInterface
}

func (v VacancyControllers) AllVacancies(c *fiber.Ctx) error {
	// получаем все queries
	searchQuery := c.Query("search_query")
	if searchQuery != "" {
		v.VacancyServices = services.Vacancy{}
		resultVacansies, err := v.VacancyServices.SearchByQuery(searchQuery)
		if err != nil {
			return fiber.NewError(500, "Ошибка при поиске вакансий")
		}
		return c.JSON(resultVacansies)
	}

	v.VacancyServices = services.Vacancy{}
	vacancies, err := v.VacancyServices.ReadAll()

	if err != nil {
		return fiber.NewError(404, "Вакансии не найдены")
	}

	return c.JSON(vacancies)
}

func (v VacancyControllers) VacancyId(c *fiber.Ctx) error {
	vacancyId, _ := strconv.Atoi(c.Params("id"))

	v.VacancyServices = services.Vacancy{VacancyId: vacancyId}
	vacancies, err := v.VacancyServices.ReadById()

	if err != nil {
		return fiber.NewError(404, "Вакансия не найдена")
	}

	return c.JSON(vacancies)
}

func (v VacancyControllers) PatchVacancy(c *fiber.Ctx) error {
	var newVacancy services.Vacancy
	c.BodyParser(&newVacancy)

	v.VacancyServices = newVacancy
	err := v.VacancyServices.Update()

	if err != nil {
		return fiber.NewError(404, "Вакансия не найдена")
	}

	return fiber.NewError(200, "Вакансия успешно обновлена")
}

func (v VacancyControllers) DeleteVacancy(c *fiber.Ctx) error {
	vacancyId, _ := strconv.Atoi(c.Query("vacancy_id"))

	v.VacancyServices = services.Vacancy{VacancyId: vacancyId}
	err := v.VacancyServices.Delete()

	if err != nil {
		return fiber.NewError(404, "Вакансия не найдена")
	}

	return fiber.NewError(200, "Вакансия успешно удалена")
}
