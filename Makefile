# Makefile para Mystery Gifter API

.PHONY: help docs generate-docs serve-docs clean test build run

# Variáveis
SWAGGER_SPEC=docs/specs/swagger.yaml
SWAGGER_CMD=go run github.com/go-swagger/go-swagger/cmd/swagger

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

docs: generate-docs serve-docs ## Gera e serve a documentação Swagger

generate-docs: ## Gera o arquivo swagger.yaml a partir dos comentários no código
	@echo "📝 Gerando documentação Swagger..."
	$(SWAGGER_CMD) generate spec -o $(SWAGGER_SPEC) --scan-models
	@echo "✅ Documentação gerada em $(SWAGGER_SPEC)"

serve-docs: ## Serve a documentação Swagger em http://localhost:8081
	@echo "🚀 Iniciando servidor de documentação..."
	@echo "📖 Acesse: http://localhost:8081"
	@echo "🛑 Para parar: Ctrl+C"
	$(SWAGGER_CMD) serve -F=swagger -p=8081 $(SWAGGER_SPEC)

clean: ## Remove arquivos gerados
	@echo "🧹 Limpando arquivos gerados..."
	rm -f $(SWAGGER_SPEC)
	@echo "✅ Limpeza concluída"

test: ## Executa os testes
	@echo "🧪 Executando testes..."
	go test ./...

build: ## Compila a aplicação
	@echo "🔨 Compilando aplicação..."
	go build -o bin/mystery-gifter-api cmd/api/main.go
	@echo "✅ Aplicação compilada em bin/mystery-gifter-api"

run: ## Executa a aplicação
	@echo "🚀 Iniciando aplicação..."
	go run cmd/api/main.go

install-tools: ## Instala as ferramentas necessárias
	@echo "📦 Instalando ferramentas..."
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	go install go.uber.org/mock/mockgen@latest
	@echo "✅ Ferramentas instaladas"
