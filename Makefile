migrate-up:
	goose -dir migrations postgres "$(DATABASE_DSN)" up

migrate-down:
	goose -dir migrations postgres "$(DATABASE_DSN)" down
