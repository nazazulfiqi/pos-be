swag:
	@swag init -g cmd/server/main.go --parseInternal -d .

echo:
	@echo "Use 'make swag' to generate Swagger docs (requires swag installed)"
