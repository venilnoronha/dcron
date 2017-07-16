package rest

// REST represents the HTTP service.
type REST interface {
	// Init starts the REST service.
	Init() error

	// Destroy stops the REST service.
	Destroy() error
}
