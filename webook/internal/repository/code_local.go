package repository

//
//import (
//	"context"
//	"errors"
//	"github.com/hashicorp/golang-lru"
//	"sync"
//	"time"
//)
//
//type LocalCodeCache struct {
//	cache      *lru.Cache
//	lock       sync.Mutex
//	expiration time.Duration
//}
//
//func NewLocalCodeCache(cache *lru.Cache, expiration time.Duration) LocalCodeCache {
//	return LocalCodeCache{
//		cache:      cache,
//		lock:       sync.Mutex{},
//		expiration: expiration,
//	}
//}
//
//func (l *LocalCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
//	l.lock.Lock()
//	defer l.lock.Unlock()
//
//	key := l.key(biz, phone)
//	now := time.Now()
//	val, ok := l.cache.Get(key)
//
//	// 说明没有验证码发送记录
//	if !ok {
//		l.cache.Add(key, codeItem{
//			code:       code,
//			cnt:        3,
//			expiration: now.Add(l.expiration),
//		})
//
//		return nil
//	}
//	itm, ok = val.(codeItem)
//	if !ok {
//		return errors.New("系统错误")
//	}
//	if itm.expiration.Sub(now) > time.Minute*9 {
//		return ErrCodeSendTooMany
//	}
//
//	l.cache.Add(key, codeItem{
//		code:       code,
//		cnt:        3,
//		expiration: now.Add(l.expiration),
//	})
//	return nil
//}
