# Migrations

This directory contains `*.sql` migration files required to generate
the `github.com/m1kola/shipsterbot/internal/migrations` package.
This package is a pure `bindata` package and it should not contain
any other go code.

The package can be generated using `make migrations`.
