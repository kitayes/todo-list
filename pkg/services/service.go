package service

import (
	"context"
	"log"
)

type Service interface {
	Init() error
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}

type ServiceManager struct {
	Services []Service
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{}
}

func (sm *ServiceManager) AddService(services ...Service) {
	for _, service := range services {
		sm.Services = append(sm.Services, service)
	}
}

func (sm *ServiceManager) Run(ctx context.Context) error {
	var err error
	for _, service := range sm.Services {
		err = service.Init()
		if err != nil {
			err = sm.Stop(ctx)
			if err != nil {
				log.Println(err.Error())
			}
			return err
		}
		go func() {
			err = service.Run(ctx)
			if err != nil {
				err = sm.Stop(ctx)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}()
	}
	return nil
}

func (sm *ServiceManager) Stop(ctx context.Context) error {
	var err error
	for _, service := range sm.Services {
		err = service.Stop(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
