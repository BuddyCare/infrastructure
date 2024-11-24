package main

import (
	"context"
	"fmt"
	config2 "github.com/BuddyCare/infrastructure/ec2/create/config"
	"github.com/BuddyCare/infrastructure/ec2/create/internal"
)

func main() {
	ctx := context.Background()
	config := config2.GetConfig("ec2_config", "json")
	ec2Src, err := internal.NewEc2Svc(ctx, config)
	if err != nil {
		panic(err)
	}
	creator := internal.NewCreator(ctx, config, ec2Src)
	err = creator.CreateEc2(ctx)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	fmt.Println("ec2 instance creation done...")
}
