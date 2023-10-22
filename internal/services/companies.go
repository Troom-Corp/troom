package services

import (
	"context"
	"fmt"

	"github.com/Troom-Corp/troom/internal"
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
	CompanyId     int
	Name          string
	Company_bio   string
	Company_photo string
	Contacts      string
	Followers     string
	Location      string
	Employees     string
	Vacansies     string // Надо будет поменять на структуру Vacancies
	Reviews       string // Надо будет поменять на структуру Reviews
}

// Create Создать компанию по входным данным и получить ID этой компании
func (c Company) Create() (int, error) {
	var companyId int
	createQuery := fmt.Sprintf("INSERT INTO public.companies (name, companybio, companyphoto, contacts, followers, location, employees, vacansies, reviews) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') RETURNING companyid",
		c.Name,
		c.Company_bio,
		c.Company_photo,
		c.Contacts,
		c.Followers,
		c.Location,
		c.Employees,
		c.Vacansies,
		c.Reviews)
	rows, err := internal.Store().Query(context.Background(), createQuery)
	rows.Scan(&companyId)
	return companyId, err
}

// ReadAll Прочитать все компании и вернуть их слайсом
func (c Company) ReadAll() ([]Company, error) {
	var companies []Company

	rows, _ := internal.Store().Query(context.Background(), "SELECT * FROM public.companies;")
	for rows.Next() {
		var company Company
		err := rows.Scan(
			&company.CompanyId,
			&company.Company_bio,
			&company.Company_photo,
			&company.Contacts,
			&company.Followers,
			&company.Location,
			&company.Employees,
			&company.Vacansies,
			&company.Reviews)
		if err != nil {
			return []Company{}, err
		}
		companies = append(companies, company)
	}

	return companies, nil
}

// ReadById Прочитать одну компанию по ID
func (c Company) ReadById() (Company, error) {
	var company Company
	readByIdQuery := fmt.Sprintf("SELECT * FROM public.companies WHERE companyid=%d", c.CompanyId)
	err := internal.Store().QueryRow(context.Background(), readByIdQuery).Scan(
		&company.CompanyId,
		&company.Company_bio,
		&company.Company_photo,
		&company.Contacts,
		&company.Followers,
		&company.Location,
		&company.Employees,
		&company.Vacansies,
		&company.Reviews)
	if err != nil {
		return Company{}, err
	}
	return company, nil
}

// SearchByQuery Поиск компаний по названию
func (c Company) SearchByQuery(searchQuery string) ([]Company, error) {
	var companies []Company
	searchFormat := "%" + searchQuery + "%"
	searchByQuery := fmt.Sprintf("SELECT * FROM public.companies WHERE LOWER(name) LIKE '%s'", searchFormat)
	rows, err := internal.Store().Query(context.Background(), searchByQuery)
	if err != nil {
		return []Company{}, nil
	}

	for rows.Next() {
		var company Company
		err = rows.Scan(
			&company.CompanyId,
			&company.Company_bio,
			&company.Company_photo,
			&company.Contacts,
			&company.Followers,
			&company.Location,
			&company.Employees,
			&company.Vacansies,
			&company.Reviews)
		if err != nil {
			return []Company{}, err
		}
		companies = append(companies, company)
	}

	return companies, nil
}

// Update Обновить данные компании по ID
func (c Company) Update() error {
	//name, companybio, companyphoto, contacts, followers, location, employees, vacansies, reviews
	updateByIdQuery := fmt.Sprintf("UPDATE public.companies SET name = '%s', companybio = '%s', companyphoto = '%s', contacts = '%s', followers = '%s', location = '%s', employees = '%s', vacansies = '%s', reviews = '%s' WHERE companyid = %d",
		c.Name,
		c.Company_bio,
		c.Company_photo,
		c.Contacts,
		c.Followers,
		c.Location,
		c.Employees,
		c.Vacansies,
		c.Reviews,
		c.CompanyId)
	_, err := internal.Store().Query(context.Background(), updateByIdQuery)
	if err != nil {
		return err
	}
	return nil
}

// Delete Удалить все данные компании по ID
func (c Company) Delete() error {
	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.companies WHERE companyid = %d", c.CompanyId)
	_, err := internal.Store().Query(context.Background(), deleteByIdQuery)
	if err != nil {
		return err
	}
	return nil
}
