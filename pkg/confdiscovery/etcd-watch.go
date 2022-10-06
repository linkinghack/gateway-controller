package confdiscovery

import (
	"context"

	"github.com/linkinghack/gateway-controller/pkg/log"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConfWatcher struct {
	ctx       context.Context
	etcdCli   *clientv3.Client
	eventChan clientv3.WatchChan
	stopChan  chan bool
	Key       string
	isPrefix  bool

	handlers map[mvccpb.Event_EventType][]func(event *clientv3.Event)
}

func NewEtcdConfWatcher(ctx context.Context, cli *clientv3.Client, key string, isPrefix bool) *EtcdConfWatcher {
	return &EtcdConfWatcher{
		ctx:      ctx,
		etcdCli:  cli,
		Key:      key,
		isPrefix: isPrefix,
		handlers: make(map[mvccpb.Event_EventType][]func(event *clientv3.Event)),
	}
}

func (e *EtcdConfWatcher) RegisterHandler(eventType mvccpb.Event_EventType, handler func(evt *clientv3.Event)) {
	e.handlers[eventType] = append(e.handlers[eventType], handler)
}

func (e *EtcdConfWatcher) Start() {
	log := log.GetSpecificLogger("pkg/confdiscovery/EtcdConfWatcher.Start()")
	opts := []clientv3.OpOption{}
	if e.isPrefix {
		opts = append(opts, clientv3.WithPrefix())
	}
	watchChan := e.etcdCli.Watch(e.ctx, e.Key, opts...)
	e.eventChan = watchChan
	log.WithField("key", e.Key).Info("start to watch")

	for {
		select {
		case resp := <-watchChan:
			if err := resp.Err(); err != nil {
				log.WithError(err).Error("error occurred while watching ETCD key")
			} else {
				for _, evt := range resp.Events {
					// call event handlers
					for _, handle := range e.handlers[evt.Type] {
						handle(evt)
					}
				}
			}
		case <-e.stopChan:
			return
		}
	}
}

func (e *EtcdConfWatcher) Stop() {
	e.stopChan <- true
}
