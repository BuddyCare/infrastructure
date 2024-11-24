package main

import (
	"context"
	"fmt"
	config2 "github.com/BuddyCare/infrastructure/ec2/install/config"
	"github.com/BuddyCare/infrastructure/ec2/install/internal"
)

func main() {
	ctx := context.Background()
	config := config2.GetConfig("ec2_config", "json")

	// cicd instance
	srcForCicdInst, err := internal.NewInstallSvc(ctx, config.CicdInstance, config.KeyInfo)
	if err != nil {
		panic(err)
	}
	installerForCicdInst := internal.NewInstaller(ctx, config, srcForCicdInst)
	err = installerForCicdInst.PrepareCicdInstance(ctx)
	if err != nil {
		panic(err)
	}

	//service instance
	srcForSrvInst, err := internal.NewInstallSvc(ctx, config.ServicesInstance, config.KeyInfo)
	if err != nil {
		panic(err)
	}
	installerForSrvInst := internal.NewInstaller(ctx, config, srcForSrvInst)
	err = installerForSrvInst.PrepareServiceInstance(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("installation into ec2 done...")
}
