package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Uri         string   //完整uri
	User        string   //用户名
	Password    string   //密码
	Hosts       []string //host地址
	Port        int      //端口
	Database    string   //数据库名
	Params      string   //其他相关认证
	MaxPoolSize uint64   //连接池最大数
}

func (t *Config) GetUri() string {
	if len(t.Uri) > 0 {
		return t.Uri
	}
	uri := "mongodb://"
	if len(t.User) > 0 && len(t.Password) > 0 {
		uri = fmt.Sprintf("%v%v:%v@", uri, t.User, t.Password)
	}

	host := ""
	for index, v := range t.Hosts {
		host += v + fmt.Sprintf(":%d", t.Port)
		if index < len(t.Hosts)-1 {
			host += ","
		}
	}

	uri = fmt.Sprintf("%v%v/%v", uri, host, t.Database)
	if len(t.Params) > 0 {
		uri = fmt.Sprintf("%v?%v", uri, t.Params)
	}
	t.Uri = uri
	return uri
}

func (t *Config) GetApplyURI() []*options.ClientOptions {

	var ops []*options.ClientOptions
	var uri string
	uri = "mongodb://"
	if len(t.User) > 0 && len(t.Password) > 0 {
		uri = fmt.Sprintf("%v%v:%v@", uri, t.User, t.Password)
	}
	for index, v := range t.Hosts {
		var host string
		if t.Port != 0 {
			host += v + fmt.Sprintf(":%d", t.Port)
		} else {
			host = v
		}
		if index < len(t.Hosts)-1 {
			host += ","
		}
		uri += host
	}
	uri += fmt.Sprintf("/%s", t.Database)
	if len(t.Params) > 0 {
		uri = fmt.Sprintf("%v?%v", uri, t.Params)
	}
	ops = append(ops, options.Client().ApplyURI(uri))
	if t.MaxPoolSize > 0 {
		ops = append(ops, options.Client().SetMaxPoolSize(t.MaxPoolSize))
	}
	return ops
}

func MongoDatabase(config Config) (*mongo.Database, error) {
	client, err := mongo.Connect(context.TODO(), config.GetApplyURI()...)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client.Database(config.Database), nil
}

func MongoClient(config Config) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), config.GetApplyURI()...)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
