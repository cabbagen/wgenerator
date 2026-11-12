package providers

import (
	"context"
	"sync"

	"github.com/cabbagen/wgenerator/v2/conf"
	"github.com/hibiken/asynq"
)

// 互斥锁
var (
	asynqClientMu sync.Mutex
	asynqServerMu sync.Mutex
)

// 单例实例
var (
	asynqClient *asynq.Client
	asynqServer *asynq.Server
)

type AsyncqTask struct {
	TaskId      string
	TaskPayload []byte
}

func getAsynqConfig() (asynq.RedisClientOpt, asynq.Config) {
	config, err := conf.ScanfBuildinYamlConfig()

	if err != nil {
		panic(err)
	}

	return asynq.RedisClientOpt{DB: config.Asynq.Redis.DB, Addr: config.Asynq.Redis.Addr, Password: config.Asynq.Redis.Password}, asynq.Config{Concurrency: config.Asynq.Concurrency}
}

func GetAsynqClientSingleInstance() *asynq.Client {
	asynqClientMu.Lock()

	defer asynqClientMu.Unlock()

	if asynqClient == nil {
		redisConnOpt, _ := getAsynqConfig()
		asynqClient = asynq.NewClient(redisConnOpt)
	}

	return asynqClient
}

func GetAsynqServerSingleInstance() *asynq.Server {
	asynqServerMu.Lock()

	defer asynqServerMu.Unlock()

	if asynqServer == nil {
		redisConnOpt, asynqConfig := getAsynqConfig()
		asynqServer = asynq.NewServer(redisConnOpt, asynqConfig)
	}
	return asynqServer
}

func CloseAsynqClientSingleInstance() {
	asynqClientMu.Lock()

	defer asynqClientMu.Unlock()

	if asynqClient == nil {
		return
	}
	_ = asynqClient.Close()
	asynqClient = nil
}

func CloseAsynqServerSingleInstance() {
	asynqServerMu.Lock()

	defer asynqServerMu.Unlock()

	if asynqServer == nil {
		return
	}
	asynqServer.Shutdown()
	asynqServer = nil
}

func RegisterAsynqTask(client *asynq.Client, task AsyncqTask, option ...asynq.Option) (*asynq.TaskInfo, error) {
	return client.Enqueue(asynq.NewTask(task.TaskId, task.TaskPayload, option...))
}

func StartAsynqTaskMonitoring(handlerMap map[string]func(context.Context, *asynq.Task) error) error {
	server := GetAsynqServerSingleInstance()

	mux := asynq.NewServeMux()

	for taskName, handler := range handlerMap {
		mux.HandleFunc(taskName, handler)
	}

	return server.Start(mux)
}

func ShutdownAsynq() {
	CloseAsynqServerSingleInstance()
	CloseAsynqClientSingleInstance()
}
