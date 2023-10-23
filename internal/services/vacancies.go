package services

import (
	"context"
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
	conn := storage.SqlInterface.New()

	createQuery := fmt.Sprintf("INSERT INTO public.vacancies (title, content, feedback, tags) VALUE ('%s', '%s', '%s', '%s'))", v.Title, v.Content, v.FeedBack, v.Tags)
	_, err := conn.Query(context.Background(), createQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return err
	}

	storage.SqlInterface.Close(conn)
	return nil
}

func (v Vacancy) ReadById() (Vacancy, error) {
	var vacancy Vacancy
	conn := storage.SqlInterface.New()

	readByIdQuery := fmt.Sprintf("SELECT * FROM public.vacancies WHERE vacancyid=%d", v.VacancyId)
	rows, err := conn.Query(context.Background(), readByIdQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return Vacancy{}, err
	}

	err = rows.Scan(&vacancy.VacancyId,
		&vacancy.CompanyId,
		&vacancy.Title,
		&vacancy.Content,
		&vacancy.FeedBack,
		&vacancy.Tags)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return Vacancy{}, err
	}

	storage.SqlInterface.Close(conn)
	return vacancy, nil
}

func (v Vacancy) ReadAll() ([]Vacancy, error) {
	var vacancies []Vacancy
	conn := storage.SqlInterface.New()

	rows, err := conn.Query(context.Background(), "SELECT * FROM public.vacancies")

	if err != nil {
		storage.SqlInterface.Close(conn)
		return []Vacancy{}, err
	}

	for rows.Next() {
		var vacancy Vacancy
		err = rows.Scan(
			&vacancy.VacancyId,
			&vacancy.CompanyId,
			&vacancy.Title,
			&vacancy.Content,
			&vacancy.FeedBack,
			&vacancy.Tags)

		if err != nil {
			storage.SqlInterface.Close(conn)
			return []Vacancy{}, err
		}
		vacancies = append(vacancies, vacancy)
	}

	storage.SqlInterface.Close(conn)
	return vacancies, nil
}

func (v Vacancy) Update() error {
	conn := storage.SqlInterface.New()

	updateByIdQuery := fmt.Sprintf("UPDATE public.vacancies SET title = '%s', content = '%s', feedback = '%s', tags = '%s'", v.Title, v.Content, v.FeedBack, v.Tags)
	_, err := conn.Query(context.Background(), updateByIdQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return err
	}

	storage.SqlInterface.Close(conn)
	return nil
}

func (v Vacancy) Delete() error {
	conn := storage.SqlInterface.New()

	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.vacancies WHERE vacancyid = %d", v.VacancyId)
	_, err := conn.Query(context.Background(), deleteByIdQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return err
	}

	storage.SqlInterface.Close(conn)
	return nil
}
