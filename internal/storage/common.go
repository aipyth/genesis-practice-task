package storage

type Storage interface {
	Connect() error
	Disconnect() error
	Save() error
	Set(string, interface{}) error
	Get(string) (interface{}, error)
	Delete(string) error
}
