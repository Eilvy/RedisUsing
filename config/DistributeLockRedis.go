package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// DistributeLockRedis 分布式锁结构体
type DistributeLockRedis struct {
	redis      *redis.Client
	ctx        context.Context    //上下文
	cancelFunc context.CancelFunc //context.CancelFunc函数，用于取消与给定上下文相关联的操作
	key        string
	expireTime time.Duration //锁的过期时间
	status     bool          //一个flag用于判断锁的状态
	waitTime   time.Duration //加锁时的等待时间
}

// TryLock 尝试加锁操作
func (d *DistributeLockRedis) TryLock() error {
	if err := d.lock(); err != nil { //加锁失败
		return err
	}

	d.status = true //加锁成功将锁状态设为true
	go d.watchDog()
	return nil
}

// lock 加锁操作（不可导出）用于TryLock调用
func (d *DistributeLockRedis) lock() error { //作为DistributeLockRedis类型的一个方法

	now := time.Now()                   //给出初始时间，开始加锁的时间
	for time.Since(now) <= d.waitTime { //需要在waitTime等待时间内不停的进行加锁操作
		isLock, err := d.redis.SetNX(d.ctx, d.key, "", d.expireTime).Result() //设置redis锁并返回bool类型的isLock是否设置成功和一个err

		if err != nil {
			return err
		}

		if !isLock {
			time.Sleep(1000 * time.Millisecond) //未能上锁成功，停止1000ms
		} else {
			return nil
		}
	}

	return errors.New("try lock time out") //加锁失败，返回超时
}

// watchDog 设置看门狗//需设置新的携程启动看门狗，不随着函数结束而结束//不可导出
// 看门狗，用于判断事务是否进行中来给锁延长存活时长//基于context 的cancelFunc实现
func (d *DistributeLockRedis) watchDog() {
	for {
		select {
		case <-d.ctx.Done(): //上下文的done里如果存在则看门狗取消
			return
		default:
			if d.status { //如果事务尚未完成，看门狗对锁进行续期
				err := d.redis.Set(d.ctx, d.key, "", d.expireTime).Err()
				if err != nil {
					fmt.Println("看门狗续期失败:", err.Error())
					return
				}

				time.Sleep(d.waitTime / 2) //每隔一段时间看门狗检测一次
			}
		}
	}
}

// Unlock 解锁，通常和看门狗联合使用
func (d *DistributeLockRedis) Unlock() error {
	d.cancelFunc() //释放context,给予信号(done)

	if d.status { //若锁的状态是加锁状态则解锁
		err := d.redis.Del(context.Background(), d.key).Err()
		if err != nil {
			return err
		}

		d.status = false //锁删除成功将锁的状态改为false
		return nil
	}

	return errors.New("unlock failed")
}

// NewDistributeLockRedis 创建锁主要函数
func NewDistributeLockRedis(key string, expireTime time.Duration, waitTime ...time.Duration) *DistributeLockRedis {
	wait := time.Second * 3 //设定默认初始等待时间为三秒
	if len(waitTime) > 0 {  //若waitTime切片内有数据则将等待时间赋值为该切片中的第一项
		wait = waitTime[0]
	}

	ctx, cancelFunc := context.WithCancel(context.Background()) //``通过背景上下文创建一个新的上下文，并获取到用于取消该上下文的取消函数

	return &DistributeLockRedis{
		redis:      DB,
		key:        key,
		expireTime: expireTime,
		waitTime:   wait,
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
}
