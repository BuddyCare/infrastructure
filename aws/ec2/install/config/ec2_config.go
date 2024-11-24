package config

import (
	"github.com/spf13/viper"
	"log"
)

type EC2Config struct {
	CicdInstance     InstanceInfo
	ServicesInstance InstanceInfo
	KeyInfo          KeyInfo
}

type InstanceInfo struct {
	User     string
	Name     string
	PublicIp string
}

type KeyInfo struct {
	Name     string
	Location string
}

func GetConfig(configName, configType string) EC2Config {
	conf := EC2Config{}
	viper.SetConfigName(configName)
	viper.AddConfigPath("./aws/ec2/install/config")
	viper.SetConfigType(configType)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return conf
}
