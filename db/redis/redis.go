package redis

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	rediscache "github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/go-redsync/redsync"
	redis2 "github.com/gomodule/redigo/redis"
	"github.com/vmihailenco/msgpack"
)

type Options struct {
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string
	// host:port address.
	Host string
	Port string
	// Dialer creates new network connection and has priority over
	// Network and Addr options.
	Dialer func() (net.Conn, error)

	// Hook that is called when new connection is established.
	OnConnect func(*redis.Conn) error

	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string
	// Database to be selected after connecting to the server.
	DB int

	// Maximum number of retries before giving up.
	// Default is to not retry failed commands.
	MaxRetries int
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	// but idle connections are still discarded by the client
	// if IdleTimeout is set.
	IdleCheckFrequency time.Duration

	// Enables read only queries on slave nodes.
	readOnly bool

	// TLS Config to use. When set TLS will be negotiated.
	TLSConfig *tls.Config
}

var defaultOptions = Options{
	Host:     "127.0.0.1",
	Port:     "6379",
	Password: "",
	DB:       0,
}

func NewRedisOption() *Options {
	return &Options{}
}

func (o *Options) SetHost(add string) *Options {
	o.Host = add
	return o
}

func (o *Options) SetPort(port string) *Options {
	o.Port = port
	return o
}

func (o *Options) SetPwd(pwd string) *Options {
	o.Password = pwd
	return o
}

func (o *Options) SetDB(db int) *Options {
	o.DB = db
	return o
}

func (o *Options) SetPoolSize(size int) *Options {
	o.PoolSize = size
	return o
}

func getOptions(options ...*Options) *redis.Options {
	option := &redis.Options{}
	option.Password = defaultOptions.Password
	option.DB = defaultOptions.DB
	option.Addr = fmt.Sprintf("%s:%s", defaultOptions.Host, defaultOptions.Port)

	for _, op := range options {
		if op.Host != "" && op.Port != "" {
			option.Addr = fmt.Sprintf("%s:%s", op.Host, op.Port)
		}
		if op.DB != 0 {
			option.DB = op.DB
		}

		if op.Password != "" {
			option.Password = op.Password
		}
		if op.PoolSize != 0 {
			option.PoolSize = op.PoolSize
		}
	}

	return option
}

func Client(options ...*Options) (*redis.Client, error) {
	rdb := redis.NewClient(getOptions(options...))
	status := rdb.Ping()
	if err := status.Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}

// CacheClient 	is a redis client
func CacheClient(options ...*Options) (*rediscache.Codec, error) {
	client, err := Client(options...)
	if err != nil {
		return nil, err
	}

	codec := &rediscache.Codec{
		Redis: client,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
	return codec, nil
}

func Sync(addr, password string) *redsync.Redsync {
	pools := []redsync.Pool{&redis2.Pool{
		MaxIdle:     2,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis2.Conn, error) {
			c, err := redis2.Dial("tcp", addr)
			if err != nil {
				panic(fmt.Errorf("初始化Redis分布式锁失败:%s", err.Error()))
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					panic(fmt.Errorf("初始化Redis分布式锁Auth失败:%s", err.Error()))
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis2.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}}
	return redsync.New(pools)
}
