package daemonwaiter

import (
	"context"
	"sync"

	"github.com/linkinghack/gateway-controller/pkg/log"
)

// DaemonWaiter manages multiple blocking daemon goroutines
type DaemonWaiter struct {
	stop         chan struct{}
	ctx          context.Context
	cancleCtx    context.CancelFunc
	daemonsQueue chan DaemonToWait
	daemons      []DaemonToWait
}

var defaultDaemonWaiter *DaemonWaiter
var deafultDwLock *sync.Mutex

func init() {
	deafultDwLock = &sync.Mutex{}
}

type DaemonToWait interface {
	Start()
	Stop()
}

func NewDaemonWaiter() *DaemonWaiter {
	ctx, cancle := context.WithCancel(context.Background())

	return &DaemonWaiter{
		stop:         make(chan struct{}),
		daemonsQueue: make(chan DaemonToWait),
		ctx:          ctx,
		cancleCtx:    cancle,
	}
}

func GetDefaultDaemonWaiter() *DaemonWaiter {
	if defaultDaemonWaiter == nil {
		initDefaultDw()
	}
	return defaultDaemonWaiter
}

func initDefaultDw() {
	deafultDwLock.Lock()
	defer deafultDwLock.Unlock()

	if defaultDaemonWaiter == nil {
		defaultDaemonWaiter = NewDaemonWaiter()
	}
}

func (d *DaemonWaiter) GetContext() context.Context {
	return d.ctx
}

func (d *DaemonWaiter) AddAndStart(dw DaemonToWait) {
	if dw == nil {
		return
	}
	d.daemons = append(d.daemons, dw)
	d.daemonsQueue <- dw
}

func (d *DaemonWaiter) Run() {
	logger := log.GetSpecificLogger("pkg/kit/daemonwaiter/DaemonWaiter.Run")
	for {
		select {
		case di := <-d.daemonsQueue:
			go di.Start()
			logger.Debug("A daemon has started")
		case <-d.stop:
			logger.Info("Stopping DaemonWaiter")
			for _, item := range d.daemons {
				item.Stop()
			}
			return
		}
	}
}

func (d *DaemonWaiter) Stop() {
	d.cancleCtx()
	close(d.daemonsQueue)
	d.stop <- struct{}{}
}
