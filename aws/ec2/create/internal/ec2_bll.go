package internal

import (
	"context"
	"github.com/BuddyCare/infrastructure/ec2/create/config"
)

type Ec2Service interface {
	CreateEc2Instance(ctx context.Context) error
}

type Creator struct {
	config config.EC2Config
	ec2Svc Ec2Service
}

func NewCreator(ctx context.Context, config config.EC2Config, ec2Svc Ec2Service) Creator {
	creator := Creator{
		config: config,
		ec2Svc: ec2Svc,
	}

	return creator
}

func (cr Creator) CreateEc2(ctx context.Context) error {
	err := cr.ec2Svc.CreateEc2Instance(ctx)
	if err != nil {
		return err
	}
	return nil
}
