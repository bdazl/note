install-prereq:
	# https://github.com/mattn/go-sqlite3?tab=readme-ov-file#installation
	CGO_ENABLED=1 go install github.com/mattn/go-sqlite3
