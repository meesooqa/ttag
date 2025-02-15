# TTAG: Save Tags from Telegram Exported HTML

- Export Telegram channel history.
- Move html-files to `var/data/you_channel/*.html`
- `docker compose up`
- `go run ./app/main.go` (go 1.23.4)
- Check `mongodb://localhost:27017` , database: `db_tags`, collection: `tags`.