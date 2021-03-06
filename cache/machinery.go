package cache

import (
	"errors"
	"sync"

	"github.com/RichardKnop/machinery/v1"
	"github.com/go-redis/redis"
)

var (
	machineryRedisClient   *redis.Client
	machineryServer        *machinery.Server
	machineryServerMutex   sync.RWMutex
	machineryActiveWorkers []*machinery.Worker
)

func SetMachineryServer(s *machinery.Server) {
	machineryServerMutex.Lock()
	machineryServer = s
	machineryServerMutex.Unlock()
}

func GetMachineryServer() *machinery.Server {
	machineryServerMutex.RLock()
	defer machineryServerMutex.RUnlock()

	if machineryServer == nil {
		panic(errors.New("Tried to get machinery server before cache#SetMachineryServer() was called"))
	}

	return machineryServer
}

func SetMachineryRedisClient(s *redis.Client) {
	machineryServerMutex.Lock()
	defer machineryServerMutex.Unlock()

	machineryRedisClient = s
}

func GetMachineryRedisClient() *redis.Client {
	machineryServerMutex.RLock()
	defer machineryServerMutex.RUnlock()

	if machineryRedisClient == nil {
		panic(errors.New("Tried to get machinery redis client before cache#SetMachineryRedisClient() was called"))
	}

	return machineryRedisClient
}

func HasMachineryRedisClient() bool {
	machineryServerMutex.RLock()
	defer machineryServerMutex.RUnlock()

	if machineryRedisClient == nil {
		return false
	}
	return true
}

func AddMachineryActiveWorker(worker *machinery.Worker) {
	machineryActiveWorkers = append(machineryActiveWorkers, worker)
}

func RemoveMachineryActiveWorker(worker *machinery.Worker) {
	newWorkers := make([]*machinery.Worker, 0)
	for _, activeWorker := range machineryActiveWorkers {
		if activeWorker.ConsumerTag == worker.ConsumerTag {
			continue
		}

		newWorkers = append(newWorkers, activeWorker)
	}
	machineryActiveWorkers = newWorkers
}

func GetMachineryActiveWorkers() (workers []*machinery.Worker) {
	return machineryActiveWorkers
}
