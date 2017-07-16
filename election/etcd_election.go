package election

import (
	"context"
	"fmt"
	"math/rand"

	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

// EtcdLeaderElection is an implementation of LeaderElection based on etcd.
type EtcdLeaderElection struct {
	LeaderElection
	session  *concurrency.Session
	election *concurrency.Election
}

// NewEtcdLeaderElection returns a new instance of EtcdLeaderElection.
func NewEtcdLeaderElection(client *etcd.Client, key string) (*EtcdLeaderElection, error) {
	log.Info("Creating etcd session")
	session, err := concurrency.NewSession(client)
	if err != nil {
		log.WithField("err", err).Error("Couldn't create etcd session")
		return nil, err
	}
	log.Info("Created etcd session")

	election := concurrency.NewElection(session, key)
	return &EtcdLeaderElection{session: session, election: election}, nil
}

// Campaign initiates leader election which can be observed using the returned chan.
// The campaign can be canceled using the returned CancelFunc.
func (e *EtcdLeaderElection) Campaign() (chan struct{}, context.CancelFunc) {
	ch := make(chan struct{})
	ctx, cancelFunc := context.WithCancel(context.TODO())
	go func() {
		for {
			value := fmt.Sprintf("value-%d", rand.Int())
			log.Info("Campaining for leadership")
			if err := e.election.Campaign(ctx, value); err != nil {
				if err.Error() == "context canceled" {
					// Return normally when the context is canceled using the CancelFunc.
					return
				}
				log.WithField("err", err).Error("Failed to campaign for leadership, retrying.. ")
			} else {
				log.Info("Obtained leadership")
				ch <- struct{}{}
				return
			}
		}
	}()
	return ch, cancelFunc
}

// Resign resigns from a previously obtained leadership.
func (e *EtcdLeaderElection) Resign() error {
	log.Info("Resigning from leadership")
	if err := e.election.Resign(context.TODO()); err != nil {
		log.WithField("err", err).Error("Failed to resign")
		return err
	}
	log.Info("Resigned from leadership")
	return nil
}

// Close closes the etcd session.
func (e *EtcdLeaderElection) Close() {
	log.Info("Closing etcd session")
	e.session.Close()
	log.Info("Closed etcd session")
}
