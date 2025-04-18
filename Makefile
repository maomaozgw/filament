.PHONY: all backend frontend clean

all: backend frontend

backend:
	@echo "Building backend..."
	@mkdir -p dist
	go build -o dist/filament ./main.go

frontend:
	@echo "Building frontend..."
	@cd static && pnpm install && pnpm run build

clean:
	@echo "Cleaning builds..."
	@rm -rf dist static/dist

gen-api:
	@echo "Generating API..."
	swag init --parseDependency