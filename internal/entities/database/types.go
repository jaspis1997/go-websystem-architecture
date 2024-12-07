package database

type Model interface {
	SetEntity(entity any) error
}
