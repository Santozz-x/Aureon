// Package migrations embeds the SQL migration files so the Gateway binary
// can apply schema changes on boot without shipping .sql files alongside it.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
