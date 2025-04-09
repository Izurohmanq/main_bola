package repositories

import (
	orderRepositories "order-service/repositories/order"
	orderFieldRepositories "order-service/repositories/orderfield"
	orderHistoryRepositories "order-service/repositories/orderhistory"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetOrder() orderRepositories.IOrderRepository
	GetOrderHistory() orderHistoryRepositories.IOrderHistoryRepository
	GetOrderField() orderFieldRepositories.IOrderFieldRepository
	GetTx() *gorm.DB
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetOrder() orderRepositories.IOrderRepository {
	return orderRepositories.NewOrderRepository(r.db)
}

func (r *Registry) GetOrderHistory() orderHistoryRepositories.IOrderHistoryRepository {
	return orderHistoryRepositories.NewOrderHistoryRepository(r.db)
}
func (r *Registry) GetOrderField() orderFieldRepositories.IOrderFieldRepository {
	return orderFieldRepositories.NewOrderFieldRepository(r.db)
}

func (r *Registry) GetTx() *gorm.DB {
	return r.db
}
