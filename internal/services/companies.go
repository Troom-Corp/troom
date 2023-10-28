package services

import (
	"fmt"

	"github.com/Troom-Corp/troom/internal/storage"
)

type CompanyInterface interface {
	Create() (int, error)
	ReadAll() ([]Company, error)
	ReadById() (Company, error)
	SearchByQuery(string) ([]Company, error)
	Update() error
	Delete() error
}

type Company struct {
	CompanyId    int
	CompanyName  string
	CompanyBio   string
	CompanyPhoto string
	Contacts     string
	Followers    string
	Location     string
	Employees    string
	Vacancies    string
	Reviews      string
}

// Create Создать компанию по входным данным и получить ID этой компании
func (c Company) Create() (int, error) {
	var companyId int
	conn, err := storage.Sql.Open()

	createQuery := fmt.Sprintf("INSERT INTO public.companies (name, companybio, companyphoto, contacts, followers, location, employees, vacansies, reviews) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') RETURNING companyid",
		c.CompanyName,
		c.CompanyBio,
		c.CompanyPhoto,
		c.Contacts,
		c.Followers,
		c.Location,
		c.Employees,
		c.Vacancies,
		c.Reviews)
	err = conn.Get(&companyId, createQuery)

	conn.Close()
	return companyId, err
}

// ReadAll Прочитать все компании и вернуть их слайсом
func (c Company) ReadAll() ([]Company, error) {
	var companies []Company
	conn, err := storage.Sql.Open()

	if err != nil {
		return []Company{}, err
	}

	err = conn.Select(&companies, "SELECT * FROM public.companies;")

	conn.Close()
	return companies, nil
}

// ReadById Прочитать одну компанию по ID
func (c Company) ReadById() (Company, error) {
	var company Company
	conn, err := storage.Sql.Open()
	if err != nil {
		return Company{}, err
	}

	readByIdQuery := fmt.Sprintf("SELECT * FROM public.companies WHERE companyid=%d", c.CompanyId)
	err = conn.Get(&company, readByIdQuery)

	conn.Close()
	return company, nil
}

// SearchByQuery Поиск компаний по названию
func (c Company) SearchByQuery(searchQuery string) ([]Company, error) {
	var companies []Company
	conn, err := storage.Sql.Open()

	searchFormat := "%" + searchQuery + "%"
	searchByQuery := fmt.Sprintf("SELECT * FROM public.companies WHERE LOWER(companyname) LIKE LOWER('%s')", searchFormat)
	err = conn.Select(&companies, searchByQuery)

	if err != nil {
		conn.Close()
		return []Company{}, nil
	}

	conn.Close()
	return companies, nil
}

// Update Обновить данные компании по ID
func (c Company) Update() error {
	conn, err := storage.Sql.Open()

	updateByIdQuery := fmt.Sprintf("UPDATE public.companies SET companyname = '%s', companybio = '%s', companyphoto = '%s', contacts = '%s', followers = '%s', location = '%s', employees = '%s', vacancies = '%s', reviews = '%s' WHERE companyid = %d",
		c.CompanyName,
		c.CompanyBio,
		c.CompanyPhoto,
		c.Contacts,
		c.Followers,
		c.Location,
		c.Employees,
		c.Vacancies,
		c.Reviews,
		c.CompanyId)
	_, err = conn.Query(updateByIdQuery)

	if err != nil {
		conn.Close()
		return err
	}

	conn.Close()
	return nil
}

// Delete Удалить все данные компании по ID
func (c Company) Delete() error {
	conn, err := storage.Sql.Open()

	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.companies WHERE companyid = %d", c.CompanyId)
	_, err = conn.Query(deleteByIdQuery)

	if err != nil {
		conn.Close()
		return err
	}

	conn.Close()
	return nil
}
