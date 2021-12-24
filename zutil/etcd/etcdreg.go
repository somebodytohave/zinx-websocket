package etcd

import (
	"context"
	"github.com/sun-fight/zinx-websocket/global"
	"go.etcd.io/etcd/client/v3"
	"time"
)

//注册租约服务
type ServiceReg struct {
	client        *clientv3.Client
	lease         clientv3.Lease //租约
	leaseResp     *clientv3.LeaseGrantResponse
	canclefunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

func NewServiceReg(addr []string, timeNum time.Duration) (*ServiceReg, error) {
	var (
		err    error
		client *clientv3.Client
	)

	if client, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}); err != nil {
		return nil, err
	}

	ser := &ServiceReg{
		client: client,
	}

	if err := ser.setLease(timeNum); err != nil {
		return nil, err
	}
	go ser.ListenLeaseRespChan()
	return ser, nil
}

//设置租约
func (this *ServiceReg) setLease(timeNum time.Duration) error {
	lease := clientv3.NewLease(this.client)

	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()
	leaseResp, err := lease.Grant(ctx, int64(timeNum.Seconds()))
	if err != nil {
		cancel()
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.TODO())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}

	this.lease = lease
	this.leaseResp = leaseResp
	this.canclefunc = cancelFunc
	this.keepAliveChan = leaseRespChan
	return nil
}

//监听续租情况
func (this *ServiceReg) ListenLeaseRespChan() {
	for {
		select {
		case leaseKeepResp := <-this.keepAliveChan:
			if leaseKeepResp == nil {
				global.Glog.Error("已经关闭续租功能")
				return
			} else {
				global.Glog.Info("续租成功")
			}
		}
	}
}

//注册租约
func (this *ServiceReg) PutService(key, val string) error {
	kv := clientv3.NewKV(this.client)
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(this.leaseResp.ID))
	return err
}

//撤销租约
func (this *ServiceReg) RevokeLease() error {
	this.canclefunc()
	time.Sleep(2 * time.Second)
	_, err := this.lease.Revoke(context.TODO(), this.leaseResp.ID)
	return err
}
