/*
 * @Description:
 * @Author: gphper
 * @Date: 2021-09-21 10:08:32
 */
package global

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"goredismanager/common"
	"goredismanager/model"
	"net"
	"path"
	"strings"
	"sync"
)

var RedisServiceStorage map[string]RedisService
var HostName string
var Port string
var ConfigViper *viper.Viper
var Accounts map[string]string
var Limit int64

var UseClient GlobalClient

type GlobalClient struct {
	ConnectName string
	Db          int
	Client      *redis.Client
}

type RedisService struct {
	RedisService string
	Config       *redis.Options
	UseSsh       int
	SSHConfig    model.SSHConfig
	Client       *redis.Client
}

var instance *MysqlClient
var once sync.Once
var Db *sql.DB

//数据库客户端

type MysqlClient struct{}

//单例模式
func GetInstance() *MysqlClient {
	once.Do(func() {
		instance = &MysqlClient{}
	})
	return instance
}

func init() {
	RedisServiceStorage = make(map[string]RedisService)
	Accounts = make(map[string]string)

	//获取配置文件
	var configPath string
	flag.StringVar(&configPath, "c", "./config.yaml", "配置文件路径")
	flag.Parse()

	basePath := path.Base(configPath)
	fileInfo := strings.Split(basePath, ".")
	ConfigViper = viper.New()

	ConfigViper.AddConfigPath(path.Dir(configPath))
	ConfigViper.SetConfigName(fileInfo[0])
	ConfigViper.SetConfigType(fileInfo[1])
	err := ConfigViper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	HostName = ConfigViper.GetString("hostname")
	if HostName == "" {
		HostName = "127.0.0.1"
	}

	Port = ConfigViper.GetString("port")
	if Port == "" {
		Port = "8088"
	}

	Limit = ConfigViper.GetInt64("limit")
	if Limit == 0 {
		Limit = 100
	}
	//redis
	connections := ConfigViper.Get("connections")

	if connections != nil {

		slice_conns := connections.([]interface{})

		for _, v := range slice_conns {
			vv := v.(map[interface{}]interface{})

			optionConfig := &redis.Options{
				Addr:     vv["host"].(string) + ":" + vv["port"].(string),
				Password: vv["password"].(string),
				DB:       0,
			}

			if vv["usessh"] == 1 {
				sshConfig := vv["sshconfig"].(map[interface{}]interface{})
				cli, err := common.GetSSHClient(sshConfig["sshusername"].(string), sshConfig["sshpassword"].(string), sshConfig["sshhost"].(string)+":"+sshConfig["sshport"].(string))
				if nil != err {
					panic(err)
				}
				optionConfig.Dialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
					return cli.Dial(network, addr)
				}
			}

			client := redis.NewClient(optionConfig)

			client.AddHook(common.RedisLog{
				Logger: common.NewLogger(vv["servicename"].(string)),
			})

			_, err := client.Ping(context.Background()).Result()
			if err != nil {
				panic(vv["servicename"].(string) + "连接失败:" + err.Error())
			}

			RsSlice := RedisService{
				RedisService: vv["servicename"].(string),
				Config:       optionConfig,
				Client:       client,
			}

			RedisServiceStorage[vv["servicename"].(string)] = RsSlice

		}

		for name, conn := range RedisServiceStorage {
			//设置全局参数
			UseClient.ConnectName = name
			UseClient.Db = 0
			UseClient.Client = conn.Client
		}

	}

	//mysql
	mysql_con := ConfigViper.Get("mysql_config")
	if mysql_con != nil {

		mysql_slice_conns := mysql_con.([]interface{})
		for _, mv := range mysql_slice_conns {
			mvv := mv.(map[interface{}]interface{})
			//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
			dsn := strings.Join([]string{mvv["username"].(string), ":", mvv["password"].(string), "@tcp(", mvv["host"].(string), ":", mvv["port"].(string), ")/", mvv["database"].(string), "?charset=utf8"}, "")
			//fmt.Print(path)
			Db, err = sql.Open("mysql", dsn)
			if err != nil {
				panic(errors.New("mysql连接失败"))
			}
			Db.SetConnMaxLifetime(100)
			Db.SetMaxIdleConns(10)

			if err := Db.Ping(); err != nil {
				fmt.Print(err.Error())
				panic(mvv["servicename"].(string) + "连接失败:" + err.Error())
			}
		}
	}

	accounts := ConfigViper.Get("accounts")
	if accounts != nil {
		slice_account := accounts.([]interface{})
		for _, account := range slice_account {
			accountMap := account.(map[interface{}]interface{})
			Accounts[accountMap["account"].(string)] = accountMap["password"].(string)
		}
	}

}
