package etcd

import (
	"context"
	"github.com/sun-fight/zinx-websocket/global"
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

var etcdKvClient *clientv3.Client
var mu sync.Mutex

func GetInstance() *clientv3.Client {
	if etcdKvClient == nil {
		if client, err := clientv3.New(clientv3.Config{
			Endpoints:   global.Object.EtcdConfig.Endpoints,
			DialTimeout: 5 * time.Second,
		}); err != nil {
			global.Glog.Error("etcd", zap.Error(err))
			return nil
		} else {
			//创建时才加锁
			mu.Lock()
			defer mu.Unlock()
			etcdKvClient = client
			return etcdKvClient
		}

	}
	return etcdKvClient
}

func Put(key, value string) error {
	_, err := GetInstance().Put(context.Background(), key, value)
	return err
}

func Get(key string) (resp *clientv3.GetResponse, err error) {
	resp, err = GetInstance().Get(context.Background(), key)
	return resp, err
}
