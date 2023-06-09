package dispersed_lock

import (
	"context"
	"fmt"
	"github.com/lfxnxf/zdy_tools/logging"
	"github.com/lfxnxf/zdy_tools/resource/redis"
	"github.com/lfxnxf/zdy_tools/utils"
	"go.uber.org/zap"
	"reflect"
	"sync"
	"time"
)

//github.com/lfxnxf/while

const (
	// 解锁lua
	unLockScript = "if redis.call('get', KEYS[1]) == ARGV[1] " +
		"then redis.call('del', KEYS[1]) return 1 " +
		"else " +
		"return 0 " +
		"end"

	// 看门狗lua
	watchLogScript = "if redis.call('get', KEYS[1]) == ARGV[1] " +
		"then return redis.call('expire', KEYS[1], ARGV[2]) " +
		"else " +
		"return 0 " +
		"end"

	lockMaxLoopNum = 1000 //加锁最大循环数量
)

var scriptMap sync.Map

type option func() (bool, error)

type DispersedLock struct {
	key            string        // 锁key
	value          string        // 锁的值，随机数
	expire         int           // 锁过期时间,单位秒
	lockClient     *redis.Redis  // 锁客户端，暂时只有redis
	unLockScript   string        // lua脚本
	watchLogScript string        // 看门狗lua
	options        []option      // 事件
	unlockCh       chan struct{} // 解锁通知通道
}

func New(ctx context.Context, client *redis.Redis, key string, expire int) *DispersedLock {
	d := &DispersedLock{
		key:    key,
		expire: expire,
		value:  fmt.Sprintf("%d", utils.Random(100000000, 999999999)), // 随机值作为锁的值
	}

	//初始化连接
	d.lockClient = client

	//初始化lua script
	lockScript, _ := scriptMap.LoadOrStore("dispersed_lock", d.getScript(ctx, unLockScript))
	watchLogScript, _ := scriptMap.LoadOrStore("watch_log", d.getScript(ctx, watchLogScript))

	d.unLockScript = lockScript.(string)
	d.watchLogScript = watchLogScript.(string)

	d.unlockCh = make(chan struct{}, 0)

	return d
}

func (d *DispersedLock) getScript(ctx context.Context, script string) string {
	scriptString, _ := d.lockClient.ScriptLoad(ctx, script)
	return scriptString
}

// 注册事件
func (d *DispersedLock) RegisterOptions(f ...option) {
	d.options = append(d.options, f...)
}

// 加锁
func (d *DispersedLock) Lock(ctx context.Context) bool {
	ok, _ := d.lockClient.SetnxEx(ctx, d.key, d.value, d.expire)
	if ok {
		go d.watchDog(ctx)
	}
	return ok
}

// 循环加锁
func (d *DispersedLock) LoopLock(ctx context.Context, sleepTime int) bool {
	t := time.NewTicker(time.Duration(sleepTime) * time.Millisecond)
	w := utils.NewWhile(lockMaxLoopNum)
	w.For(func() {
		if d.Lock(ctx) {
			t.Stop()
			w.Break()
		} else {
			<-t.C
		}
	})
	if !w.IsNormal() {
		return false
	}
	return true
}

// 循环获取锁并且绑定事件
// eg:单个线程获取缓存、其它线程等待
func (d *DispersedLock) LoopLockWithOption(ctx context.Context, sleepTime int) (bool, error) {
	t := time.NewTicker(time.Duration(sleepTime) * time.Millisecond)
	w := utils.NewWhile(lockMaxLoopNum)
	var err error
	w.For(func() {
		locked := d.Lock(ctx)
		if locked { // 获取到锁，跳出循环
			t.Stop()
			w.Break()
		}

		var flag bool
		for _, option := range d.options {
			flag, err = option()
			if err != nil { //事件代码出现异常，跳出循环
				t.Stop()
				w.Break()
			}
			if !flag {
				break
			}
		}

		//所有事件全部为true，不用等到获取锁，直接跳出
		if flag {
			t.Stop()
			w.Break()
		}

		<-t.C
	})
	return true, err
}

// 解锁
func (d *DispersedLock) Unlock(ctx context.Context) bool {
	args := []interface{}{
		d.value, // 脚本中的argv
	}
	flag, err := d.lockClient.EvalSha(ctx, d.unLockScript, []string{d.key}, args...)
	// 关闭看门狗
	close(d.unlockCh)
	return lockRes(flag, err)
}

// 看门狗
func (d *DispersedLock) watchDog(ctx context.Context) {
	// 创建一个定时器NewTicker, 每过期时间的3分之2触发一次
	loopTime := time.Duration(d.expire*1e3/3) * time.Millisecond
	expTicker := time.NewTicker(loopTime)
	//确认锁与锁续期打包原子化
	for {
		select {
		case <-expTicker.C:
			args := []interface{}{
				d.value,
				d.expire,
			}
			res, err := d.lockClient.EvalSha(ctx, d.watchLogScript, []string{d.key}, args...)
			if err != nil {
				logging.Error("watchDog error", zap.Error(err))
				return
			}
			r, ok := res.(int64)
			if !ok {
				return
			}
			if r == 0 {
				return
			}
		case <-d.unlockCh: //任务完成后用户解锁通知看门狗退出
			return
		}
	}
}

func lockRes(flag interface{}, err error) bool {
	if err != nil {
		logging.Errorw("unlock error", zap.Error(err))
		return false
	}
	r, ok := flag.(int64)
	if !ok {
		logging.Errorw("flag types error",
			zap.String("type", reflect.TypeOf(flag).String()),
			zap.Any("flag", flag),
		)
		return false
	}
	if r > 0 {
		return true
	} else {
		return false
	}
}
