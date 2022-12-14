package storage

type ID = string
type URL = string

type Storage interface {
	Add(url string) (id string, err error)
	Get(id string) (string, error)
}
