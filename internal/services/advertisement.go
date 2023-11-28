package services

import (
	"fmt"

	"github.com/Troom-Corp/troom/internal/storage"
)

type AdvertisementInterface interface {
	Create() error
	Update() error
	ReadAll() ([]Advertisement, error)
	ReadById() (Advertisement, error)
	Delete() error
}

type Advertisement struct {
	AdvertId  int
	AuthorId  int
	Time      string
	Title     string
	Blocks    string
	StartTime string
	EndTime   string
}

func (a Advertisement) Create() error {
	conn, err := storage.Sql.Open()
	if err != nil {
		conn.Close()
		return err
	}
	createQuery := fmt.Sprintf("INSERT INTO public.advertisement (advertid, authorid, time, title, blocks, starttime, endtime) VALUES (%d, %d, '%s', '%s', '%s', '%s', '%s')",
		a.AdvertId,
		a.AuthorId,
		a.Time,
		a.Title,
		a.Blocks,
		a.StartTime,
		a.EndTime)
	_, err = conn.Query(createQuery)
	if err != nil {
		conn.Close()
		return err
	}
	conn.Close()
	return nil
}

func (a Advertisement) Update() error {
	conn, err := storage.Sql.Open()

	if err != nil {
		conn.Close()
		return err
	}

	updateQuery := fmt.Sprintf("UPDATE public.advertisement SET title = '%s' blocks = '%s' WHERE advertid = %d",
		a.Title,
		a.Blocks,
		a.AdvertId)

	_, err = conn.Query(updateQuery)

	if err != nil {
		conn.Close()
		return err
	}

	conn.Close()
	return nil
}

func (a Advertisement) ReadAll() ([]Advertisement, error) {
	var Advertisments []Advertisement
	conn, err := storage.Sql.Open()

	if err != nil {
		conn.Close()
		return []Advertisement{}, err
	}

	readAllQuery := fmt.Sprintf("SELECT * FROM public.advertisement")
	conn.Select(&Advertisments, readAllQuery)

	conn.Close()
	return Advertisments, nil
}

func (a Advertisement) ReadById() (Advertisement, error) {
	var Advert Advertisement
	conn, err := storage.Sql.Open()

	if err != nil {
		conn.Close()
		return Advertisement{}, nil
	}

	readByIdQuery := fmt.Sprintf("SELECT * FROM public.advertisement WHERE advertid = %d", a.AdvertId)
	err = conn.Get(&Advert, readByIdQuery)

	if err != nil {
		conn.Close()
		return Advertisement{}, nil
	}

	conn.Close()
	return Advert, nil
}

func (a Advertisement) Delete() error {
	conn, err := storage.Sql.Open()

	if err != nil {
		conn.Close()
		return err
	}

	deleteQuery := fmt.Sprintf("DELETE FROM public.advertisement WHERE advertid = %d", a.AdvertId)
	_, err = conn.Query(deleteQuery)

	if err != nil {
		conn.Close()
		return err
	}

	conn.Close()
	return nil
}
