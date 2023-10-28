package services

import (
	"fmt"
	"github.com/Troom-Corp/troom/internal/storage"
)

type VacancyInterface interface {
	Create() error
	ReadAll() ([]Vacancy, error)
	ReadById() (Vacancy, error)
	Update() error
	Delete() error
}

type Vacancy struct {
	VacancyId int
	CompanyId int
	Title     string
	Content   string
	FeedBack  string // это просто counter откликов
	Tags      string
}

func (v Vacancy) Create() error {
	conn, err := storage.Sql.Open()

	createQuery := fmt.Sprintf("INSERT INTO public.vacancies (title, content, feedback, tags) VALUE ('%s', '%s', '%s', '%s'))", v.Title, v.Content, v.FeedBack, v.Tags)

	_, err = conn.Query(createQuery)
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}

func (v Vacancy) ReadById() (Vacancy, error) {
	var vacancy Vacancy
	conn, err := storage.Sql.Open()

	readByIdQuery := fmt.Sprintf("SELECT * FROM public.vacancies WHERE vacancyid=%d", v.VacancyId)
	err = conn.Get(&vacancy, readByIdQuery)

	if err != nil {
		conn.Close()
		return Vacancy{}, err
	}

	conn.Close()
	return vacancy, nil
}

func (v Vacancy) ReadAll() ([]Vacancy, error) {
	var vacancies []Vacancy
	conn, err := storage.Sql.Open()

	err = conn.Select(&vacancies, "SELECT * FROM public.vacancies")

	if err != nil {
		conn.Close()
		return []Vacancy{}, err
	}

	conn.Close()
	return vacancies, nil
}

func (v Vacancy) Update() error {
	conn, err := storage.Sql.Open()

	updateByIdQuery := fmt.Sprintf("UPDATE public.vacancies SET title = '%s', content = '%s', feedback = '%s', tags = '%s'", v.Title, v.Content, v.FeedBack, v.Tags)
	_, err = conn.Query(updateByIdQuery)

	if err != nil {
		conn.Close()
		return err
	}

	conn.Close()
	return nil
}

func (v Vacancy) Delete() error {
	conn, err := storage.Sql.Open()

	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.vacancies WHERE vacancyid = %d", v.VacancyId)
	_, err = conn.Query(deleteByIdQuery)

	if err != nil {
		conn.Close()
		return err
	}

	conn.Close()
	return nil
}
