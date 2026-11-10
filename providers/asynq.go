package providers

import (
	"github.com/cabbagen/wgenerator/v2/conf"
	"github.com/hibiken/asynq"
)

var (
	AsynqClient *asynq.Client = nil
	AsynqServer *asynq.Server = nil
)

func GetAsynqClientSingleInstance() *asynq.Client {
	if AsynqClient != nil {
		return AsynqClient
	}

	config, error := conf.ScanfBuildinYamlConfig()

	if error != nil {
		panic(error)
	}

	AsynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     config["asynq"]["address"].(string),
		Password: config["asynq"]["password"].(string),
		DB:       config["asynq"]["db"].(int),
	})

	return AsynqClient
}
