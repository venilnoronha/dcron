package election

import "context"

// LeaderElection represents an abstraction of a leader election implementation.
type LeaderElection interface {
	// Campaign campaigns for leadership and notifies it via the returned chan.
	// The campaign can be cancelled using the returned CancelFunc.
	Campaign() (chan struct{}, context.CancelFunc)

	// Resigns resigns from a previously acquired leadership.
	Resign() error

	// Close releases all resources associated with the leadership.
	Close()
}
