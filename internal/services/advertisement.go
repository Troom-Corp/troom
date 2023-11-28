package services

import (
	"github.com/Troom-Corp/troom/internal/storage"
)

type AdvertisementInterface interface {
	Create() error
	Update() error
	ReadAll() ([]Advertisement, error)
	ReadById(AdvertId int) (Advertisement, error)
	Delete() error
}

type Advertisement struct {
	AdvertId  int
	UserId    int
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
	createQuery := "INSERT INTO public.advertisment (advertid, userid, time, title, blocks, starttime, endtime) VALUES (%d, %d, '%s', '%s', '%s', '%s', '%s')"
}
