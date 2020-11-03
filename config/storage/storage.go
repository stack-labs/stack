package storage

type Storage interface {
	Write(namespace string, config []byte) (err error)
	Load(namespace string) (config []byte, err error)
}
