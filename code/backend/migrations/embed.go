// Package migrations embeds the SQL migration files so the compiled server is
// self-contained: the runtime image ships only the binary, not this directory.
package migrations

import "embed"

// Files holds every *.sql migration in this directory. Up files are applied in
// filename order (timestamp prefix) on boot; down files exist for rollback and
// are never applied automatically.
//
// The //go:embed path is resolved relative to THIS file, so the directive lives
// beside the files it embeds rather than in cmd/api where it would look in
// cmd/api/migrations/.
//
//go:embed *.sql
var Files embed.FS
