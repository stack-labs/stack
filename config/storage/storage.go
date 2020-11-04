package storage

type Storage interface {
	Write(config []byte) (err error)
	Load() (config []byte, err error)
}
