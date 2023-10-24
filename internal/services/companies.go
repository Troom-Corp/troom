package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/Troom-Corp/troom/internal/storage"
)

type CompanyInterface interface {
	Create() error
	ReadAll() ([]Company, error)
	ReadById() (Company, error)
	SearchByQuery(string) ([]Company, error)
	Update() error
	Delete() error
}

type Company struct {
	CompanyId    int
	Name         string
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
	conn := storage.SqlInterface.New()

	createQuery := fmt.Sprintf("INSERT INTO public.companies (name, companybio, companyphoto, contacts, followers, location, employees, vacansies, reviews) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') RETURNING companyid",
		c.Name,
		c.CompanyBio,
		c.CompanyPhoto,
		c.Contacts,
		c.Followers,
		c.Location,
		c.Employees,
		c.Vacancies,
		c.Reviews)
	rows, err := conn.Query(context.Background(), createQuery)
	rows.Scan(&companyId)

	storage.SqlInterface.Close(conn)
	return companyId, err
}

// ReadAll Прочитать все компании и вернуть их слайсом
func (c Company) ReadAll() ([]Company, error) {
	var companies []Company
	conn := storage.SqlInterface.New()

	rows, _ := conn.Query(context.Background(), "SELECT * FROM public.companies;")
	for rows.Next() {
		var company Company
		err := rows.Scan(
			&company.CompanyId,
			&company.CompanyBio,
			&company.CompanyPhoto,
			&company.Contacts,
			&company.Followers,
			&company.Location,
			&company.Employees,
			&company.Vacancies,
			&company.Reviews)

		if err != nil {
			storage.SqlInterface.Close(conn)
			return []Company{}, err
		}
		companies = append(companies, company)
	}

	storage.SqlInterface.Close(conn)
	return companies, nil
}

// ReadById Прочитать одну компанию по ID
func (c Company) ReadById() (Company, error) {
	var company Company
	conn := storage.SqlInterface.New()

	readByIdQuery := fmt.Sprintf("SELECT * FROM public.companies WHERE companyid=%d", c.CompanyId)
	conn.QueryRow(context.Background(), readByIdQuery).Scan(
		&company.CompanyId,
		&company.CompanyBio,
		&company.CompanyPhoto,
		&company.Contacts,
		&company.Followers,
		&company.Location,
		&company.Employees,
		&company.Vacancies,
		&company.Reviews)

	storage.SqlInterface.Close(conn)
	return company, nil
}

// SearchByQuery Поиск компаний по названию
func (c Company) SearchByQuery(searchQuery string) ([]Company, error) {
	var companies []Company
	conn := storage.SqlInterface.New()

	searchFormat := "%" + strings.ToLower(searchQuery) + "%"
	searchByQuery := fmt.Sprintf("SELECT * FROM public.companies WHERE LOWER(name) LIKE '%s'", searchFormat)
	rows, err := conn.Query(context.Background(), searchByQuery)

	if err != nil {
		return []Company{}, nil
	}

	for rows.Next() {
		var company Company
		err = rows.Scan(
			&company.CompanyId,
			&company.CompanyBio,
			&company.CompanyPhoto,
			&company.Contacts,
			&company.Followers,
			&company.Location,
			&company.Employees,
			&company.Vacancies,
			&company.Reviews)

		if err != nil {
			storage.SqlInterface.Close(conn)
			return []Company{}, err
		}
		companies = append(companies, company)
	}

	storage.SqlInterface.Close(conn)
	return companies, nil
}

// Update Обновить данные компании по ID
func (c Company) Update() error {
	conn := storage.SqlInterface.New()

	updateByIdQuery := fmt.Sprintf("UPDATE public.companies SET name = '%s', companybio = '%s', companyphoto = '%s', contacts = '%s', followers = '%s', location = '%s', employees = '%s', vacansies = '%s', reviews = '%s' WHERE companyid = %d",
		c.Name,
		c.CompanyBio,
		c.CompanyPhoto,
		c.Contacts,
		c.Followers,
		c.Location,
		c.Employees,
		c.Vacancies,
		c.Reviews,
		c.CompanyId)
	_, err := conn.Query(context.Background(), updateByIdQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return err
	}

	storage.SqlInterface.Close(conn)
	return nil
}

// Delete Удалить все данные компании по ID
func (c Company) Delete() error {
	conn := storage.SqlInterface.New()

	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.companies WHERE companyid = %d", c.CompanyId)
	_, err := conn.Query(context.Background(), deleteByIdQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return err
	}

	storage.SqlInterface.Close(conn)
	return nil
}
