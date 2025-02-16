# TTAG: Save Tags from Telegram Exported HTML

1. Export Telegram channel history.
2. Move html-files to `var/data/you_channel/*.html`
3. `docker compose up`
4. `go run ./app/main.go` (go 1.23.4)
5. Check `mongodb://localhost:27017` , database: `db_tags`, collection: `tags`.
6. `docker compose down`