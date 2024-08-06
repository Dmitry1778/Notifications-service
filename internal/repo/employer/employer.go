package employer

import (
	"context"
	"notify/internal/domain"
)

//go:generate mockgen -source=employer.go -destination=mocks/mock.go

func (r Repo) List(ctx context.Context, id int) ([]int, error) {
	return r.userStorage.GetPublisherList(ctx, id)
}

func (r Repo) Create(ctx context.Context, employee domain.Employee) error {
	// to something
	return r.userStorage.Register(ctx, employee)
}

func (r Repo) Unsub(ctx context.Context, sub, pub int) error {
	return r.userStorage.Unsubscribe(ctx, sub, pub)
}

func (r Repo) Sub(ctx context.Context, sub, pub int) error {
	return r.userStorage.Subscribe(ctx, sub, pub)
}

func (r Repo) Get(ctx context.Context) (*[]domain.ResponseEmployee, error) {
	return r.userStorage.GetAllEmployee(ctx)
}

func (r Repo) GetID(ctx context.Context, id int) (*domain.ResponseEmployee, error) {
	return r.userStorage.GetEmployeeID(ctx, id)
}

func New(storage userStorage) *Repo {
	return &Repo{
		storage,
	}
}

type Repo struct {
	userStorage userStorage
}

type userStorage interface {
	GetPublisherList(ctx context.Context, id int) ([]int, error)
	Unsubscribe(ctx context.Context, sub, pub int) error
	Subscribe(ctx context.Context, sub, pub int) error
	Register(ctx context.Context, emp domain.Employee) error
	GetAllEmployee(ctx context.Context) (*[]domain.ResponseEmployee, error)
	GetEmployeeID(ctx context.Context, id int) (*domain.ResponseEmployee, error)
}
