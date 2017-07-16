package config

import (
	"context"

	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/clientv3"
)

type EtcdCronConfigService struct {
	CronConfigService
	client *etcd.Client
	key    string
}

func NewEtcdCronConfigService(client *etcd.Client, key string) *EtcdCronConfigService {
	return &EtcdCronConfigService{client: client, key: key}
}

func (c *EtcdCronConfigService) Load() (*CronConfig, error) {
	log.Info("Loading cron config from etcd")
	resp, _ := c.client.Get(context.TODO(), c.key)
	var conf CronConfig
	if len(resp.Kvs) == 0 {
		conf = CronConfig{Config: "", Version: 0}
	} else {
		kv := resp.Kvs[0]
		conf = CronConfig{Config: string(kv.Value), Version: kv.ModRevision}
	}
	log.WithField("conf", conf).Info("Loaded cron config from etcd")
	return &conf, nil
}

func (c *EtcdCronConfigService) Save(conf *CronConfig) error {
	log.Info("Saving cron config to etcd")

	_, err := c.client.
		Txn(context.TODO()).
		If(etcd.ModRevision(string(conf.Version))).
		Then(etcd.OpPut(c.key, conf.Config)).
		Commit()

	if err != nil {
		log.WithField("err", err).Fatal("Failed to save cron config to etcd")
		return err
	}
	log.Info("Saved cron config to etcd")
	return nil
}
