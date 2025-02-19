# TTAG: Save Tags from Telegram Exported HTML

## Save data
1. Configure app `etc/config.yml`.
2. Export Telegram channel history.
3. Move html-files to `%system.data_path%/you_channel/*.html`
4. `docker compose up`
5. `go run ./app/save.go` (go 1.23.4)
6. Check `mongodb://localhost:27017`, database: `db_tags`, collection: `tags` (`%mongo.uri%`, `%mongo.database%`, `%mongo.collection_messages%`).
7. `docker compose down`.