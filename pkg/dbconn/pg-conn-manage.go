package dbconn

import (
	"database/sql"
	"fmt"
	"reflect"
	"sync"

	"github.com/linkinghack/gateway-controller/config"
	"github.com/linkinghack/gateway-controller/pkg/log"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DBConnStore struct {
	connStr    string // [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	gormEngine *gorm.DB
	pool       *sql.DB
}

var dbConnStoreSingleton *DBConnStore
var initLock *sync.Mutex

func init() {
	initLock = &sync.Mutex{}
}

// InitDBConn initiate the DB connection singleton object with global configs
//
//	It is safe to call this method multiple times incidentally.
func InitDBConn() {
	initLock.Lock()
	defer initLock.Unlock()
	logger := log.GetSpecificLogger("pkg/database/relational/InitDB")

	if dbConnStoreSingleton == nil {
		dbConnStore := DBConnStore{}

		dbconf := config.GetGlobalConfig().DBConfig
		dbConnStore.connStr = generateDBUrl(dbconf.EngineType, dbconf.User, dbconf.Password, dbconf.Database, dbconf.Host, dbconf.Port, marshalDBConnParam(dbconf.ConnectionParameters))

		logger.Infof("Initiating db connections. addr=%s:%d", dbconf.Host, dbconf.Port)
		logger.Debugf("Initiating db connections. dburl=%s", dbConnStore.connStr)

		// Init gorm helper
		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			DSN: dbConnStore.connStr,
		}), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   dbconf.TablesPrefix,
				SingularTable: true,
			},
		})
		if err != nil {
			logger.WithError(err).Error("Database connection cannot establish")
			return
		}

		// config DB pool size
		pool, err := gormDB.DB()
		if err != nil {
			logger.WithError(err).Error("Get original DB connection failed")
			return
		}
		pool.SetMaxIdleConns(dbconf.MaxIdleSize)
		pool.SetMaxOpenConns(dbconf.MaxPoolSize)

		dbConnStore.pool = pool
		dbConnStore.gormEngine = gormDB

		dbConnStoreSingleton = &dbConnStore
	}
}

// AdaptTableName returns the table name based on current global db config
// and orm engine config. Add prefixs and split words.
func (db *DBConnStore) AdaptTableName(tableObj interface{}) string {
	return db.GetGOrmEngine().NamingStrategy.TableName(reflect.TypeOf(tableObj).Name())
}

func (db *DBConnStore) AssureTables(dataModels ...interface{}) error {
	err := db.gormEngine.AutoMigrate(dataModels...)

	if err != nil {
		return errors.Wrap(err, "Create tables failed")
	}
	return nil
}

func (db *DBConnStore) GetGOrmEngine() *gorm.DB {
	if db.gormEngine == nil {
		InitDBConn()
	}
	return db.gormEngine
}

func (db *DBConnStore) GetOriginalGoSqlDB() *sql.DB {
	return db.pool
}

// GetDBConn Get the singleton DB connection store object
func GetDBConn() *DBConnStore {
	if dbConnStoreSingleton == nil {
		InitDBConn()
	}
	return dbConnStoreSingleton
}

func marshalDBConnParam(params map[string]string) string {
	result := ""
	for k, v := range params {
		if len(result) > 0 {
			result = result + "&"
		}
		result = result + k + "=" + v
	}
	return result
}

func generateDBUrl(engineType, dbUser, password, dbName, host string, port int, otherParams string) string {
	url := fmt.Sprintf("%s://%s:%s@%s:%d/%s", engineType, dbUser, password, host, port, dbName)
	if len(otherParams) > 0 {
		url = fmt.Sprintf("%s?%s", url, otherParams)
	}
	// if engineType == "postgres" {
	// 	url = fmt.Sprintf("%s://%s", engineType, url)
	// }
	return url
}
