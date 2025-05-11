package storage

type Storage interface {
}

type MemoryStorage struct {
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}
