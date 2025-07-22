package discovery

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type QueueWorker struct {
	key    string
	kqConf kq.KqConf        // 存储在etcd中的kafka配置信息, 连接信息
	client *clientv3.Client // etcd client
}

func NewQueueWorker(key string, endpoints []string, kqConf kq.KqConf) *QueueWorker {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 3,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return &QueueWorker{
		key:    key,
		client: etcdClient,
		kqConf: kqConf,
	}
}

func (q *QueueWorker) HeartBeat() {
	value, err := json.Marshal(q.kqConf)
	if err != nil {
		panic(err)
	}
	q.register(string(value), context.Background())
}
func (q *QueueWorker) register(value string, ctx context.Context) {
	leaseGrantResp, err := q.client.Grant(ctx, 600)
	if err != nil {
		panic(err)
	}
	leaseId := leaseGrantResp.ID
	logx.Infof("leaseID: %x", leaseId)
	kv := clientv3.NewKV(q.client)
	putResp, err := kv.Put(ctx, q.key, value, clientv3.WithLease(leaseId))
	if err != nil {
		panic(err)
	}

	keepRespChan, err := q.client.KeepAlive(context.TODO(), leaseId)
	if err != nil {
		panic(err)
	}

	go func() {
		//nolint:gocritic
		for {
			select {
			case keepResp, ok := <-keepRespChan:
				if !ok {
					logx.Infof("租约已经失效:%x", leaseId)
					q.register(value, ctx)
					return
				} else {
					logx.Infof("收到自动续租应答:%x", keepResp.ID)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	logx.Info("写入成功:", putResp.Header.Revision)
}
