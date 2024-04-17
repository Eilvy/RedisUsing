package config

import (
	"github.com/redis/go-redis/v9"
	"sync"
)

var DB *redis.Client
var Wg sync.WaitGroup

type User struct {
	Username string `form:"Username"`
	Number   string `form:"Number"`
}
