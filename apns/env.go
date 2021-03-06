// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package apns

import (
	"log"

	"github.com/go-ini/ini"
	"github.com/codegangsta/cli"

	"zooinit/cluster"
	"zooinit/config"
	"gopush/lib"
)

// This basic discovery service bootstrap env info
type EnvInfo struct {
	cluster.BaseInfo

	CertPath          string
	CertPassword      string
	CertENV           string
	CertTopic         string

	PoolConfig        *lib.PoolConfig

	QueueSourceConfig *lib.QueueSourceConfig

	WorkerPool        *lib.WorkerPool
}

func NewEnvInfo(iniobj *ini.File, c *cli.Context) *EnvInfo {
	env := new(EnvInfo)

	sec := iniobj.Section(CONFIG_SECTION)
	env.Service = sec.Key("service").String()
	if env.Service == "" {
		log.Fatalln("Config of service section is empty.")
	}

	// parse base info
	env.ParseConfigFile(sec, c)

	keyNow := "cert.env"
	env.CertENV = config.GetValueString(keyNow, sec, c)
	if env.CertENV == "" {
		log.Fatalln("Config of " + keyNow + " is empty.")
	}

	keyNow = "cert.path"
	env.CertPath = config.GetValueString(keyNow, sec, c)
	if env.CertPath == "" {
		log.Fatalln("Config of " + keyNow + " is empty.")
	}

	keyNow = "cert.password"
	env.CertPassword = config.GetValueString(keyNow, sec, c)
	if env.CertPassword == "" {
		log.Fatalln("Config of " + keyNow + " is empty.")
	}

	keyNow = "cert.topic"
	env.CertTopic = config.GetValueString(keyNow, sec, c)
	if env.CertTopic == "" {
		log.Fatalln("Config of " + keyNow + " is empty.")
	}

	qsConfig:=&lib.QueueSourceConfig{}
	keyNow = "queue.method"
	tmpStr := config.GetValueString(keyNow, sec, c)
	if tmpStr == "" {
		log.Fatalln("Config of " + keyNow + " is empty.")
	}
	if tmpStr!=lib.QUEUE_SOURCE_METHOD_API && tmpStr!=lib.QUEUE_SOURCE_METHOD_FILE && tmpStr!=lib.QUEUE_SOURCE_METHOD_MYSQL {
		log.Fatalln("Config of " + keyNow + " value is not allowed: "+tmpStr)
	}
	qsConfig.Method=tmpStr

	keyNow = "queue.cache.path"
	tmpStr = config.GetValueString(keyNow, sec, c)
	if tmpStr == "" {
		log.Fatalln("Config of " + keyNow + " is empty.")
	}
	qsConfig.CachePath=tmpStr

	if qsConfig.Method == lib.QUEUE_SOURCE_METHOD_API {
		keyNow = "queue.api.uri"
		tmpStr = config.GetValueString(keyNow, sec, c)
		if tmpStr == "" {
			log.Fatalln("Config of " + keyNow + " is empty.")
		}
		qsConfig.ApiPrefix=tmpStr

		//can be empty
		keyNow = "queue.api.default"
		tmpStr = config.GetValueString(keyNow, sec, c)
		qsConfig.Value=tmpStr
	}else if qsConfig.Method == lib.QUEUE_SOURCE_METHOD_MYSQL {
		keyNow = "queue.mysql.dsn"
		tmpStr = config.GetValueString(keyNow, sec, c)
		if tmpStr == "" {
			log.Fatalln("Config of " + keyNow + " is empty.")
		}
		qsConfig.MysqlDsn=tmpStr

		//can be empty
		keyNow = "queue.mysql.sql"
		tmpStr = config.GetValueString(keyNow, sec, c)
		qsConfig.Value=tmpStr
	}else if qsConfig.Method == lib.QUEUE_SOURCE_METHOD_FILE {
		keyNow = "queue.file.path"
		tmpStr = config.GetValueString(keyNow, sec, c)
		if tmpStr == "" {
			log.Fatalln("Config of " + keyNow + " is empty.")
		}
		qsConfig.FilePath=tmpStr

		//can be empty
		keyNow = "queue.file.default"
		tmpStr = config.GetValueString(keyNow, sec, c)
		qsConfig.Value=tmpStr
	}
	//set qsconfig
	env.QueueSourceConfig=qsConfig
	wp, err := lib.NewWorkerPool(env)
	if err != nil {
		log.Fatalln("Create lib.NewWorkerPool error: " + err.Error())
	} else {
		env.WorkerPool = wp
	}

	//create uuid
	env.CreateUUID()

	env.GuaranteeSingleRun()

	//register signal watcher
	env.RegisterSignalWatch()

	return env
}

func (e *EnvInfo) CreateWorker() (lib.Worker, error) {
	worker, err := NewWorker(e)
	if err != nil {
		return nil, err
	}

	return worker, nil
}

// TODO destroy
func (e *EnvInfo) DestroyWorker(worker lib.Worker) (error) {


	return nil
}

func (e *EnvInfo) GetPoolConfig() (*lib.PoolConfig) {
	return e.PoolConfig
}

func (e *EnvInfo) GetWorkerPool() (*lib.WorkerPool) {
	return e.WorkerPool
}

func (e *EnvInfo) GetQueueSourceConfig() (*lib.QueueSourceConfig) {
	return e.QueueSourceConfig
}