package internal

import (
	"context"
	"github.com/BuddyCare/infrastructure/ec2/install/config"
)

type installService interface {
	InstallDocker(ctx context.Context) error
	InstallJenkins(ctx context.Context) error
	InstallAnsible(ctx context.Context) error
	InstallK8s(ctx context.Context) error
}

type Installer struct {
	config     config.EC2Config
	installSrv installService
}

func NewInstaller(ctx context.Context, config config.EC2Config, installSrv installService) Installer {
	creator := Installer{
		config:     config,
		installSrv: installSrv,
	}

	return creator
}

func (cr Installer) PrepareCicdInstance(ctx context.Context) error {
	err := cr.installSrv.InstallDocker(ctx)
	if err != nil {
		return err
	}

	err = cr.installSrv.InstallJenkins(ctx)
	if err != nil {
		return err
	}

	err = cr.installSrv.InstallAnsible(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (cr Installer) PrepareServiceInstance(ctx context.Context) error {
	err := cr.installSrv.InstallDocker(ctx)
	if err != nil {
		return err
	}

	err = cr.installSrv.InstallK8s(ctx)
	if err != nil {
		return err
	}
	return nil
}
