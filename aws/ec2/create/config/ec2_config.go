package config

import (
	"github.com/spf13/viper"
	"log"
)

type EC2Config struct {
	Region       string
	User         string
	Name         string
	ImageId      string
	InstanceType string
	MinCount     int32
	MaxCount     int32
	KeyName      string
	KeyLocation  string
	Credentials  Credentials
}

type Credentials struct {
	AccessKeyId string
	SecretKeyId string
}

func GetConfig(configName, configType string) EC2Config {
	conf := EC2Config{}
	viper.SetConfigName(configName)
	viper.AddConfigPath("./aws/ec2/create/config")
	viper.SetConfigType(configType)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	viper.BindEnv("AWS_ACCESS_KEY_ID")
	viper.BindEnv("AWS_SECRET_ACCESS_KEY")
	accessKey := viper.GetString("AWS_ACCESS_KEY_ID")
	secretKey := viper.GetString("AWS_SECRET_ACCESS_KEY")
	conf.Credentials = Credentials{
		AccessKeyId: accessKey,
		SecretKeyId: secretKey,
	}
	return conf
}
