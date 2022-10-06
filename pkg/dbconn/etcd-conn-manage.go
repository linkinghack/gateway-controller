package dbconn

import (
	"log"
	"sync"
	"time"

	"github.com/linkinghack/gateway-controller/config"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConnStore struct {
	ServerEndpoints []string
	client          *clientv3.Client
}

func (e *EtcdConnStore) GetClient() *clientv3.Client {
	return e.client
}

var etcdConnSingleton *EtcdConnStore
var etcdConnSingletonLock *sync.Mutex

func init() {
	etcdConnSingletonLock = &sync.Mutex{}
}

func NewETCDConnStoreWithGlobalConfig() (*EtcdConnStore, error) {
	etcdConf := config.GetGlobalConfig().EtcdConfig
	cliConf := clientv3.Config{
		Endpoints:   etcdConf.Endpoints,
		DialTimeout: 10 * time.Second,
	}
	if etcdConf.AuthType == config.EtcdAuthTypePassowrd {
		cliConf.Username = etcdConf.Username
		cliConf.Password = etcdConf.Password
	}
	cli, err := clientv3.New(cliConf)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create ETCD client")
	}

	store := EtcdConnStore{
		ServerEndpoints: etcdConf.Endpoints,
		client:          cli,
	}
	return &store, nil
}

func initEtcdConnSingleton() {
	etcdConnSingletonLock.Lock()
	defer etcdConnSingletonLock.Unlock()

	if etcdConnSingleton == nil {
		etcdConn, err := NewETCDConnStoreWithGlobalConfig()
		if err != nil {
			log.Fatal(err, "Init ETCD client failed")
			return
		}
		etcdConnSingleton = etcdConn
	}
}

// MustGetEtcdConn returns the singleton instance of EtcdConnStore
func MustGetEtcdConn() *EtcdConnStore {
	if etcdConnSingleton != nil {
		return etcdConnSingleton
	}
	initEtcdConnSingleton()

	return etcdConnSingleton
}
