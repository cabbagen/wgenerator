package providers

import (
	"context"
	"log"

	"github.com/cabbagen/wgenerator/v2/conf"
	"github.com/hibiken/asynq"
)

var (
	asynqClient *asynq.Client = nil
	asynqServer *asynq.Server = nil
)

func getAsynqConfig() (asynq.RedisClientOpt, asynq.Config) {
	config, error := conf.ScanfBuildinYamlConfig()

	if error != nil {
		panic(error)
	}

	return asynq.RedisClientOpt{DB: config.Asynq.Redis.DB, Addr: config.Asynq.Redis.Addr, Password: config.Asynq.Redis.Password}, asynq.Config{Concurrency: config.Asynq.Concurrency}
}

func GetAsynqClientSingleInstance() *asynq.Client {
	redisConnOpt, _ := getAsynqConfig()

	if asynqClient == nil {
		asynqClient = asynq.NewClient(redisConnOpt)
	}
	return asynqClient
}

func GetAsynqServerSingleInstance() *asynq.Server {
	redisConnOpt, asynqConfig := getAsynqConfig()

	if asynqServer == nil {
		asynqServer = asynq.NewServer(redisConnOpt, asynqConfig)
	}
	return asynqServer
}

func CloseAsynqClientSingleInstance() {
	if asynqClient != nil {
		asynqClient.Close()
	}
	asynqClient = nil
}

func CloseAsynqServerSingleInstance() {
	if asynqServer != nil {
		asynqServer.Shutdown()
	}
	asynqServer = nil
}

func RegisterAsynqTask(client *asynq.Client, taskName string, taskPayload []byte, option ...asynq.Option) (*asynq.TaskInfo, error) {
	return client.Enqueue(asynq.NewTask(taskName, taskPayload, option...))
}

func ServerAsynqTaskMonitering(server *asynq.Server, handlerMap map[string]func(context.Context, *asynq.Task) error) {
	mux := asynq.NewServeMux()

	for taskName, handler := range handlerMap {
		mux.HandleFunc(taskName, handler)
	}

	if err := server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
