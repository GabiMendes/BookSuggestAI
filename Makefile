# Makefile para BookSuggest AI

# Variáveis
BINARY_NAME=book-suggest-ai
MAIN_PATH=./cmd/main.go

# Comandos padrão
.PHONY: build run test clean deps help

# Compilar a aplicação
build:
	@echo "🔨 Compilando a aplicação..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Executar a aplicação
run:
	@echo "🚀 Iniciando a aplicação..."
	go run $(MAIN_PATH)

# Executar testes
test:
	@echo "🧪 Executando testes..."
	go test ./...

# Limpar arquivos gerados
clean:
	@echo "🧹 Limpando arquivos..."
	go clean
	rm -rf bin/
	rm -rf data/

# Instalar dependências
deps:
	@echo "📦 Instalando dependências..."
	go mod download
	go mod tidy

# Verificar código
lint:
	@echo "🔍 Verificando código..."
	go vet ./...
	go fmt ./...

# Executar em modo de desenvolvimento
dev:
	@echo "💻 Modo desenvolvimento..."
	go run -race $(MAIN_PATH)

# Gerar documentação
docs:
	@echo "📚 Gerando documentação..."
	godoc -http=:6060

# Setup inicial do projeto
setup:
	@echo "⚙️ Configuração inicial..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ Arquivo .env criado. Configure suas credenciais!"; \
	else \
		echo "⚠️ Arquivo .env já existe."; \
	fi
	make deps
	@echo "✅ Setup concluído!"

# Compilar para produção
build-prod:
	@echo "🏭 Compilando para produção..."
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Mostrar ajuda
help:
	@echo "📖 Comandos disponíveis:"
	@echo "  build      - Compilar a aplicação"
	@echo "  run        - Executar a aplicação"
	@echo "  dev        - Executar em modo desenvolvimento"
	@echo "  test       - Executar testes"
	@echo "  lint       - Verificar código"
	@echo "  clean      - Limpar arquivos gerados"
	@echo "  deps       - Instalar dependências"
	@echo "  setup      - Configuração inicial do projeto"
	@echo "  build-prod - Compilar para produção"
	@echo "  docs       - Gerar documentação"
	@echo "  help       - Mostrar esta ajuda"