package embeddedtemporal

import (
	"fmt"
	"os"
)

// acquireLock takes a single-instance lock for a persistent database file. A
// SQLite-backed Temporal server assumes exactly one server process owns the
// file; a second concurrent server would corrupt it. The lock is an exclusively
// created sibling file holding the owner PID; it is removed on Close.
//
// The lock is advisory and does not detect a stale lock left by a crashed
// process: the returned error explains how to clear it. (PID-liveness recovery
// is deferred to avoid platform-specific code.)
func acquireLock(databaseFile string) (func(), error) {
	lockPath := databaseFile + ".lock"
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf(
				"embedded Temporal database %q is locked by %q; another instance may be running. If no instance is running, the lock is stale and can be removed",
				databaseFile, lockPath,
			)
		}
		return nil, fmt.Errorf("lock embedded Temporal database: %w", err)
	}
	fmt.Fprintf(file, "%d\n", os.Getpid())
	_ = file.Close()

	return func() { _ = os.Remove(lockPath) }, nil
}
