package etcd

import (
	"context"
	"github.com/sun-fight/zinx-websocket/global"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"time"
)

type ClientDis struct {
	client *clientv3.Client
}

func NewClientDis(addr []string) (*ClientDis, error) {
	conf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}

	if client, err := clientv3.New(conf); err == nil {
		return &ClientDis{
			client: client,
		}, nil
	} else {
		return nil, err
	}
}

func (this *ClientDis) GetService(prefix string) ([]string, error) {
	resp, err := this.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	addrs := this.extractAddr(resp)

	go this.watcher(prefix)
	return addrs, nil
}

func (this *ClientDis) watcher(prefix string) {
	rch := this.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				this.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				this.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

func (this *ClientDis) extractAddr(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}
	for _, v := range resp.Kvs {
		this.SetServiceList(string(v.Key), string(v.Value))
		addrs = append(addrs, string(v.Value))
	}
	return addrs
}

func (this *ClientDis) SetServiceList(key, val string) {
	global.Object.EtcdServerConfig.ServerListLock.Lock()
	defer global.Object.EtcdServerConfig.ServerListLock.Unlock()
	global.Object.EtcdServerConfig.ServerList[key] = val
	global.Glog.Info("发现服务：", zap.String("key", key), zap.String(" 地址:", val))
}

func (this *ClientDis) DelServiceList(key string) {
	global.Object.EtcdServerConfig.ServerListLock.Lock()
	defer global.Object.EtcdServerConfig.ServerListLock.Unlock()
	delete(global.Object.EtcdServerConfig.ServerList, key)
	global.Glog.Info("服务下线:", zap.String("key", key))
}
