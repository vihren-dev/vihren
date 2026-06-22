package embeddedtemporal

import (
	"fmt"
	"os"

	temporalconfig "go.temporal.io/server/common/config"
	commonlog "go.temporal.io/server/common/log"
	sqliteplugin "go.temporal.io/server/common/persistence/sql/sqlplugin/sqlite"
	sqliteschema "go.temporal.io/server/schema/sqlite"
	schematool "go.temporal.io/server/tools/common/schema"
	sqltool "go.temporal.io/server/tools/sql"
)

// temporalSchemaName is the embedded schema path for the Temporal history store,
// the source-of-truth store that ships versioned migrations under
// schema/sqlite/v3/temporal/versioned. The visibility store has only a static
// schema (no versioned migrations) and is rebuildable, so it is not migrated.
const temporalSchemaName = "sqlite/v3/temporal"

// ensureSchema makes the on-disk SQLite database current before the server opens
// it. A fresh database gets the full consolidated schema plus a stamped schema
// version; an existing database is migrated forward to the version shipped by the
// linked go.temporal.io/server. This is what lets a persisted desktop database
// survive server-dependency upgrades.
//
// The work happens here (not inside litekit) so the database file already exists
// when litekit opens it, which makes litekit skip its own unversioned setup.
func ensureSchema(databaseFile string) error {
	fresh := !fileExists(databaseFile)
	sqlConfig := sqliteConfig(databaseFile)

	if fresh {
		// SetupSchema creates the latest consolidated tables for both the history
		// and visibility stores; the file is created by the rwc connection.
		if err := sqliteschema.SetupSchema(sqlConfig); err != nil {
			return fmt.Errorf("set up embedded Temporal schema: %w", err)
		}
		// Stamp the version-bookkeeping tables at the current release so future
		// UpdateTask runs only apply migrations newer than this baseline.
		if err := stampSchemaVersion(sqlConfig); err != nil {
			return err
		}
	}
	return migrateSchema(sqlConfig)
}

// stampSchemaVersion creates the schema-version bookkeeping tables and records
// the current release, without re-applying any schema (no SchemaName is set).
func stampSchemaVersion(sqlConfig *temporalconfig.SQL) error {
	connection, err := sqltool.NewConnection(sqlConfig, schemaLogger())
	if err != nil {
		return fmt.Errorf("open embedded Temporal database for versioning: %w", err)
	}
	defer connection.Close()

	task := schematool.NewSetupSchemaTask(connection, &schematool.SetupConfig{
		InitialVersion: sqliteschema.Version,
	}, schemaLogger())
	if err := task.Run(); err != nil {
		return fmt.Errorf("stamp embedded Temporal schema version: %w", err)
	}
	return nil
}

// migrateSchema applies any history-store migrations newer than the stamped
// version. It is a no-op when the database is already current, and does the real
// upgrade work after the go.temporal.io/server dependency is bumped.
func migrateSchema(sqlConfig *temporalconfig.SQL) error {
	connection, err := sqltool.NewConnection(sqlConfig, schemaLogger())
	if err != nil {
		return fmt.Errorf("open embedded Temporal database for migration: %w", err)
	}
	defer connection.Close()

	task := schematool.NewUpdateSchemaTask(connection, &schematool.UpdateConfig{
		SchemaName: temporalSchemaName,
	}, schemaLogger())
	if err := task.Run(); err != nil {
		return fmt.Errorf("migrate embedded Temporal schema: %w", err)
	}
	return nil
}

// sqliteConfig builds the file-backed SQLite persistence config (read-write-create
// mode), matching what litekit uses for non-ephemeral servers.
func sqliteConfig(databaseFile string) *temporalconfig.SQL {
	return &temporalconfig.SQL{
		PluginName:        sqliteplugin.PluginName,
		DatabaseName:      databaseFile,
		ConnectAttributes: map[string]string{"mode": "rwc"},
	}
}

// schemaLogger returns a quiet logger for schema tooling; setup/migration detail
// is not interesting to a desktop app on the happy path.
func schemaLogger() commonlog.Logger {
	return commonlog.NewNoopLogger()
}

// fileExists reports whether path exists as a regular file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
