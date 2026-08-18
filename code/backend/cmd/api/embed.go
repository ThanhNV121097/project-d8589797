package main

import "embed"

// Files holds every *.sql migration under migrations/. Up files are applied in
// filename order (timestamp prefix) on boot; down files exist for rollback and
// are never applied automatically.
//
//go:embed migrations/*.sql
var Files embed.FS
