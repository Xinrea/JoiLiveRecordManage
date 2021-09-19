package main

import (
	"joirecord/internal/api"
	"joirecord/internal/logger"

	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/spf13/viper"
)

var log = logger.Log

func main() {
	logger.Init()
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")
	viper.SetDefault("paths", []string{"S1/轴伊Joi_Channel/", "S2/轴伊Joi_Channel/", "S3/轴伊Joi_Channel/"})
	viper.SetDefault("ak", "")
	viper.SetDefault("sk", "")
	viper.SetDefault("endpoint", "https://gz.bcebos.com")
	viper.SetDefault("bucket", "winks")
	err := viper.ReadInConfig()
	if err != nil {
		viper.SafeWriteConfig()
	}
	log.Info("JoiRecord Backend Start")
	AK, SK := viper.GetString("ak"), viper.GetString("sk")

	// 用户指定的Endpoint
	ENDPOINT := viper.GetString("endpoint")

	clientConfig := bos.BosClientConfiguration{
		Ak:               AK,
		Sk:               SK,
		Endpoint:         ENDPOINT,
		RedirectDisabled: false,
	}

	// 初始化一个BosClient
	bosClient, err := bos.NewClientWithConfig(&clientConfig)
	if err != nil {
		log.Fatal(bosClient)
	}
	server := api.New(bosClient)
	server.Run()
}
