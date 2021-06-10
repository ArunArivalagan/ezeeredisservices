package redisclient

import (
	"fmt"

	"github.com/ezeeredisservices/logger"
	"github.com/go-redis/redis"
)

var redisclient *redis.Client

func init() {
	if redisclient == nil {
		GetClient()
	}
}

func GetClient() (*redis.Client, error) {
	var err error
	if redisclient == nil {
		redisclient = redis.NewClient(&redis.Options{
			/** The network type, either tcp or unix.*/
			/** Default is tcp.*/
			Network: "tcp",
			/** host:port address.*/
			Addr: "localhost:6379",
			/** Optional password. Must match the password specified in the*/
			/** requirepass server configuration option.*/
			Password: "",
			/** Database to be selected after connecting to the server.*/
			//	DB: 0,

			/** Maximum number of retries before giving up.*/
			/** Default is to not retry failed commands.*/
			// MaxRetries: 0,
			/** Minimum backoff between each retry.*/
			/** Default is 8 milliseconds; -1 disables backoff.*/
			// MinRetryBackoff: -1,
			/** Maximum backoff between each retry.*/
			/** Default is 512 milliseconds; -1 disables backoff.*/
			// MaxRetryBackoff: -1,

			/** Dial timeout for establishing new connections.*/
			/** Default is 5 seconds.*/
			// DialTimeout: 180,
			/** Timeout for socket reads. If reached, commands will fail*/
			/** with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.*/
			/** Default is 3 seconds.*/
			// ReadTimeout: 180,
			/** Timeout for socket writes. If reached, commands will fail*/
			/** with a timeout instead of blocking.*/
			/** Default is ReadTimeout.*/
			// WriteTimeout: 180,

			/** Maximum number of socket connections.*/
			/** Default is 10 connections per every CPU as reported by runtime.NumCPU.*/
			// PoolSize: 10,
			/** Minimum number of idle connections which is useful when establishing*/
			/** new connection is slow.*/
			//	MinIdleConns: 3,
			/** Connection age at which client retires (closes) the connection.*/
			/** Default is to not close aged connections.*/
			//MaxConnAge: 180,
			/** Amount of time client waits for connection if all connections*/
			/** are busy before returning an error.*/
			/** Default is ReadTimeout + 1 second.*/
			//	PoolTimeout: 120,
			/** Amount of time after which client closes idle connections.*/
			/** Should be less than server's timeout.*/
			/** Default is 5 minutes. -1 disables idle timeout check.*/
			//IdleTimeout: 180,
			/** Frequency of idle checks made by idle connections reaper.*/
			/** Default is 1 minute. -1 disables idle connections reaper,*/
			/** but idle connections are still discarded by the client*/
			/** if IdleTimeout is set.*/
			//	IdleCheckFrequency: 180,

			/** Enables read only queries on slave nodes.*/
			/** TLS Config to use. When set TLS will be negotiated.*/
		})

		_, err = redisclient.Ping().Result()
		if err != nil {
			fmt.Println("*************** Failed to Connect Redis ********" + redisclient.Options().Addr)
		} else {
			fmt.Println("*************** Redis host name ********" + redisclient.Options().Addr)
		}
	}
	_, err = redisclient.Ping().Result()
	if err == redis.Nil {
		logger.ErrorLogger.Println(err)
	}

	return redisclient, err
}

func PutLocation(key string, data []byte) error {
	var err error
	client, err := GetClient()
	if err == nil {
		err = client.Set(key, data, 0).Err()
	}
	return err
}

func GetLocations(key string) (string, error) {
	var err error
	var cacheData string
	client, err := GetClient()
	if err == nil {
		cacheData, err = client.Get(key).Result()
	}
	return cacheData, err
}

func GetAllKeys() (*redis.ScanIterator, error) {
	var err error
	client, err := GetClient()
	var iter *redis.ScanIterator
	if err == nil {
		var cursor uint64
		iter = client.Scan(cursor, "LOC_*", 0).Iterator()
	}
	return iter, err
}
