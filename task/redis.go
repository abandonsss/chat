package task

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"log"
	"time"
)

var redisClient *redis.Client

func (t *Task) InitRedis() error {
	vip := viper.New()
	vip.SetConfigType("yaml")
	vip.SetConfigName("redis")
	vip.AddConfigPath("./config")
	if err := vip.ReadInConfig(); err != nil {
		return err
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     vip.GetString("address"),
		Password: vip.GetString("password"),
		DB:       vip.GetInt("db"),
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		return errors.New("connect redis err")
	}
	go func() {
		for {
			val, err := BRPop("chat")
			if err != nil {
				log.Println(err)
			}
			if len(val) >= 2 {
				t.Push(val[1])
			}
		}
	}()
	return nil
}

func BRPop(key string) ([]string, error) {
	val, err := redisClient.BRPop(60*time.Second, key).Result()
	fmt.Println("task's msg is ", val, err)
	return val, err
}
