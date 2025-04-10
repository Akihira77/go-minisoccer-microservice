package controllers

import (
	"user-service/controllers/user"
	"user-service/services"
)

type Registry struct {
	service services.IServiceRegistry
}

type IControllerRegistry interface {
	GetUserController() user.IUserController
}

func NewRegistryController(service services.IServiceRegistry) IControllerRegistry {
	return &Registry{
		service: service,
	}
}

func (r *Registry) GetUserController() user.IUserController {
	return user.NewUserController(r.service)
}
