package svc

import (
	"context"
	"encoding/json"
	"sync"
	"time"
	"zeroim/imrpc/internal/config"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServiceContext struct {
	Config    config.Config
	QueueList *QueueList
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}

type QueueList struct {
	mu  sync.Mutex
	kqs map[string]*kq.Pusher
}

func NewQueueList() *QueueList {
	return &QueueList{
		kqs: make(map[string]*kq.Pusher),
	}
}
func (q *QueueList) Update(key string, data kq.KqConf) {
	edgeQueue := kq.NewPusher(data.Brokers, data.Topic)
	q.mu.Lock()
	q.kqs[key] = edgeQueue
	q.mu.Unlock()
}

func (q *QueueList) Delete(key string) {
	q.mu.Lock()
	delete(q.kqs, key)
	q.mu.Unlock()
}

func (q *QueueList) Load(key string) (*kq.Pusher, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	edgeQueue, ok := q.kqs[key]
	return edgeQueue, ok
}
func GetQueueList(conf discov.EtcdConf) *QueueList {
	ql := NewQueueList()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.Hosts,
		DialTimeout: time.Second * 3,
	})
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := cli.Get(ctx, conf.Key, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, kv := range res.Kvs {
		var data kq.KqConf
		if err := json.Unmarshal(kv.Value, &data); err != nil {
			logx.Errorf("invalid data key is: %s value is: %s", string(kv.Key), string(kv.Value))
			continue
		}
		if len(data.Brokers) == 0 || len(data.Topic) == 0 {
			continue
		}
		edgeQueue := kq.NewPusher(data.Brokers, data.Topic)

		ql.mu.Lock()
		ql.kqs[string(kv.Key)] = edgeQueue
		ql.mu.Unlock()
	}
	return ql
}
