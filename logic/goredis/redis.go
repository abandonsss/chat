package goredis

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"time"
)

var RedisClient *redis.Client

const expireTime = time.Hour

func InitRedis() error {
	vip := viper.New()
	vip.SetConfigName("redis")
	vip.SetConfigType("yaml")
	vip.AddConfigPath("./config")
	err := vip.ReadInConfig()
	if err != nil {
		return err
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     vip.GetString("address"),
		Password: vip.GetString("password"),
		DB:       vip.GetInt("db"),
	})
	_, err = RedisClient.Ping().Result()
	if err != nil {
		return errors.New("connect redis err")
	}
	return nil
}

func HMSet(key string, value map[string]interface{}) error {
	err := RedisClient.HMSet(key, value).Err()
	RedisClient.Expire(key, expireTime)
	return err
}

func Set(key string, value string) error {
	err := RedisClient.Set(key, value, expireTime).Err()
	return err
}

func HGetALL(key string) (map[string]string, error) {
	val, err := RedisClient.HGetAll(key).Result()
	RedisClient.Expire(key, expireTime)
	return val, err
}

func Del(key string) error {
	err := RedisClient.Del(key).Err()
	return err
}

func LPush(key string, value interface{}) error {
	err := RedisClient.LPush(key, value).Err()
	return err
}
