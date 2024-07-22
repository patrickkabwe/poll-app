install:
	@go mod download
	@yarn
	@go mod tidy
	@npx tailwindcss -i ./static/input.css -o ./static/styles.css
	@templ generate

dev: 
	@lsof -t -i tcp:5001 | xargs kill -9
	@air -c .air.toml

css:
	@npx tailwindcss -i ./static/input.css -o ./static/styles.css --watch

gen-temp:
	@templ fmt .
	@templ generate --watch --proxy="http://localhost:5001"