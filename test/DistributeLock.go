package test

import (
	"RedisUsing/config"
	"fmt"
	"sync"
	"time"
)

func Task1(key string, expireTime time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	lock := config.NewDistributeLockRedis(key, expireTime) //初始化锁//task1未设置请求时间，默认为3s
	err := lock.TryLock()                                  //上锁
	if err != nil {
		fmt.Println("lock failed task1", err.Error())
	}
	fmt.Println("task1 locked success")

	defer func() {
		err := lock.Unlock()
		if err != nil {
			fmt.Println("unlock failed task1", err.Error())
		}
		fmt.Println("task1 unlock success")
	}()

	//写task1需要执行的操作，此处用time.Sleep代替
	time.Sleep(3 * time.Second)
}

// 再开一个task2和task1竞争锁
func Task2(key string, expireTime time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	lock := config.NewDistributeLockRedis(key, expireTime, 7*time.Second) //初始化锁//task2设置请求时间
	err := lock.TryLock()                                                 //上锁
	if err != nil {
		fmt.Println("lock failed task2", err.Error())
		return
	}

	fmt.Println("task2 locked success")
	defer func() {
		err := lock.Unlock()
		if err != nil {
			fmt.Println("unlock failed task2", err.Error())
		}
		fmt.Println("task2 unlock success")
	}()

	//写task2需要执行的操作，此处用time.Sleep代替
	time.Sleep(5 * time.Second)
}
