// Package embeddedtemporal runs a Temporal server inside the current process so
// a single Go binary can host durable workflows with no external infrastructure.
// It supports two modes:
//
//   - Ephemeral (default): an in-memory server for demos, tests, and blog-sized
//     examples. State is lost when the process exits.
//   - Persistent (WithDatabaseFile): a file-backed SQLite server for single-user
//     desktop durable-agent apps. Start, do work, stop, and restart later with
//     workflow state intact. The on-disk schema is migrated forward on startup
//     (see schema.go), so the database survives go.temporal.io/server upgrades.
//
// The server is built on the vendored LiteServer (internal/litekit). The
// convenience helpers never close a door: Client, HostPort, and Namespace expose
// the underlying primitives so callers can build any client or worker the
// Temporal SDK allows.
//
// Neither mode is for production: persistent mode assumes a single server
// process owns the database file (enforced by a lock, see lock.go).
package embeddedtemporal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	commonlog "go.temporal.io/server/common/log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/vihren-dev/vihren/platform/embeddedtemporal/internal/litekit"
)

// DefaultNamespace is the namespace registered when WithNamespace is not set. A
// stable default lets a restarted persistent app address the same namespace.
const DefaultNamespace = "default"

// dialTimeout bounds the in-process client dial during Start.
const dialTimeout = 30 * time.Second

// Server is a running in-process Temporal server. Construct it with Start and
// release it with Close. The zero value is not usable.
type Server struct {
	lite      *litekit.LiteServer
	namespace string
	client    client.Client

	workerOptions worker.Options

	mu      sync.Mutex
	workers []worker.Worker
	unlock  func()
}

// config accumulates Start options.
type config struct {
	databaseFile  string // empty => ephemeral
	namespace     string
	clientOptions client.Options
	workerOptions worker.Options
}

// Option customizes an embedded server.
type Option func(*config)

// WithDatabaseFile switches Start to persistent mode, storing all state in the
// SQLite file at path. The file and its parent directory are created if needed,
// and the schema is migrated forward on startup. Without this option the server
// is ephemeral (in-memory).
func WithDatabaseFile(path string) Option {
	return func(c *config) { c.databaseFile = path }
}

// WithNamespace registers and uses namespace instead of DefaultNamespace.
func WithNamespace(namespace string) Option {
	return func(c *config) { c.namespace = namespace }
}

// WithClientOptions sets base client.Options for the client returned by Client.
// HostPort and Namespace are managed by the server and are overridden.
func WithClientOptions(options client.Options) Option {
	return func(c *config) { c.clientOptions = options }
}

// WithWorkerOptions sets base worker.Options used by StartWorker.
func WithWorkerOptions(options worker.Options) Option {
	return func(c *config) { c.workerOptions = options }
}

// Start launches an in-process Temporal server and returns it ready to use. The
// caller owns the returned Server and must call Close to release resources.
//
//	srv, err := embeddedtemporal.Start(embeddedtemporal.WithDatabaseFile(path))
//	if err != nil {
//	    return err
//	}
//	defer srv.Close()
func Start(opts ...Option) (*Server, error) {
	cfg := config{namespace: DefaultNamespace}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.namespace == "" {
		cfg.namespace = DefaultNamespace
	}

	var unlock func()
	if cfg.databaseFile != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.databaseFile), 0o755); err != nil {
			return nil, fmt.Errorf("create embedded Temporal database directory: %w", err)
		}
		release, err := acquireLock(cfg.databaseFile)
		if err != nil {
			return nil, err
		}
		unlock = release
		if err := ensureSchema(cfg.databaseFile); err != nil {
			unlock()
			return nil, err
		}
	}

	server, err := startLite(cfg, unlock)
	if err != nil {
		if unlock != nil {
			unlock()
		}
		return nil, err
	}
	return server, nil
}

// startLite builds and starts the LiteServer, dials the in-process client, and
// recovers any panic from the server libraries into an error.
func startLite(cfg config, unlock func()) (srv *Server, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			srv = nil
			err = fmt.Errorf("start embedded Temporal server: %v", recovered)
		}
	}()

	lite, err := litekit.NewLiteServer(&litekit.LiteServerConfig{
		Ephemeral:        cfg.databaseFile == "",
		DatabaseFilePath: cfg.databaseFile,
		Namespaces:       []string{cfg.namespace},
		FrontendIP:       "127.0.0.1",
		Logger:           commonlog.NewNoopLogger(),
	})
	if err != nil {
		return nil, fmt.Errorf("create embedded Temporal server: %w", err)
	}
	if err := lite.Start(); err != nil {
		return nil, fmt.Errorf("start embedded Temporal server: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	temporalClient, err := lite.NewClientWithOptions(ctx, cfg.clientOptions)
	if err != nil {
		_ = lite.Stop()
		return nil, fmt.Errorf("connect to embedded Temporal server: %w", err)
	}

	return &Server{
		lite:          lite,
		namespace:     cfg.namespace,
		client:        temporalClient,
		workerOptions: cfg.workerOptions,
		unlock:        unlock,
	}, nil
}

// Client returns a Temporal client connected to the embedded server's namespace.
// The client is closed by Close; do not close it yourself.
func (s *Server) Client() client.Client {
	return s.client
}

// HostPort returns the frontend host:port of the embedded server, for callers
// that need to dial it from another process or configure external tooling.
func (s *Server) HostPort() string {
	return s.lite.FrontendHostPort()
}

// Namespace returns the namespace registered on the embedded server.
func (s *Server) Namespace() string {
	return s.namespace
}

// StartWorker registers a worker on taskQueue and starts polling. The register
// callback receives the native worker.Registry, so generated Register functions
// drop straight in:
//
//	srv.StartWorker(myapp.DefaultTaskQueue, func(r worker.Registry) {
//	    myapp.Register(r, &myapp.Activities{...})
//	})
//
// The returned worker is stopped by Close. For full control over worker
// construction, use worker.New(srv.Client(), ...) directly instead.
func (s *Server) StartWorker(taskQueue string, register func(worker.Registry)) (worker.Worker, error) {
	if register == nil {
		return nil, fmt.Errorf("embedded Temporal worker requires a register function")
	}
	w := worker.New(s.client, taskQueue, s.workerOptions)
	register(w)
	if err := w.Start(); err != nil {
		return nil, fmt.Errorf("start embedded Temporal worker on %q: %w", taskQueue, err)
	}
	s.mu.Lock()
	s.workers = append(s.workers, w)
	s.mu.Unlock()
	return w, nil
}

// Close stops the embedded server, its workers, and its client, and releases the
// persistent-mode database lock. It is safe to call once.
func (s *Server) Close() {
	s.mu.Lock()
	workers := s.workers
	s.workers = nil
	s.mu.Unlock()
	for _, w := range workers {
		w.Stop()
	}
	if s.client != nil {
		s.client.Close()
	}
	if s.lite != nil {
		_ = s.lite.Stop()
	}
	if s.unlock != nil {
		s.unlock()
		s.unlock = nil
	}
}
