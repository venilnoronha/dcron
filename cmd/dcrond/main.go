package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"dcron/election"
	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/clientv3"
)

// electionKey is the key to be used while campaigning for leadership.
const electionKey = "dcrond"

var (
	// etcdHostFlag is the flag for obtaining etcd host address via the CLI.
	etcdHostFlag = flag.String("H", "http://127.0.0.1:2379", "etcd host address")

	// etcdHostFlag is the flag for obtaining etcd host address via the CLI.
	portFlag = flag.Int("P", 9090, "dcrond service port")

	// client is the etcd client.
	client *etcd.Client

	// leaderElection encapsulates logic for managing leadership.
	leaderElection election.LeaderElection

	// campaignCh is the chan that's notified on acquiring leadership.
	campaignCh chan struct{}

	// isLeader is set to true on acquiring leadership.
	isLeader bool = false

	// campaignCancel is a func to be used to cancel the leadership campaign.
	campaignCancel context.CancelFunc
)

func main() {
	flag.Parse()

	initEtcdClient()
	initLeaderElection()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		destroyLeaderElection()
		destroyEtcdClient()
		destroyHTTPServer()
		os.Exit(0)
	}()

	go waitAndServe()
	initHTTPServer(*portFlag)
}

func waitAndServe() {
	campaignCh, campaignCancel = leaderElection.Campaign()
	<-campaignCh
	isLeader = true

	log.Info("Going to serve")
	// TODO: Add cron logic
	log.Info("Finished serving")
}

func initEtcdClient() {
	log.Info("Connecting to etcd host ", *etcdHostFlag)
	var err error
	client, err = etcd.NewFromURL(*etcdHostFlag)
	if err != nil {
		log.WithField("err", err).Fatal("Couldn't create etcd client, exiting")
		os.Exit(1)
	}
	log.Info("Created etcd client")
}

func destroyEtcdClient() {
	log.Info("Closing etcd client")
	client.Close()
	log.Info("Closed etcd client")
}

func initLeaderElection() {
	var err error
	leaderElection, err = election.NewEtcdLeaderElection(client, electionKey)
	if err != nil {
		log.WithField("err", err).Fatal("Failed to init etcd leader election, exiting")
		os.Exit(1)
	}
}

func destroyLeaderElection() {
	if isLeader {
		leaderElection.Resign()
	} else {
		log.Info("Cancelling campaign")
		campaignCancel()
		log.Info("Campaign cancelled")
	}
	leaderElection.Close()
}
