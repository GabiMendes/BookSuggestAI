# BookSuggest AI

Sistema de Recomendação Inteligente de Livros Baseado em Histórico de Leitura com IA.

## Descrição

O BookSuggest-AI conecta-se à sua planilha pessoal do Google Sheets e utiliza IA (Google Gemini) para analisar padrões de preferência baseados no seu histórico real de leitura e avaliações, gerando recomendações personalizadas e precisas.

## Representação visual

https://github.com/user-attachments/assets/0aa9f9a2-3e4c-4680-8d98-e2c6af096eb8

## Funcionalidades

- ✅ Autenticação via Google OAuth 2.0
- ✅ Integração com Google Sheets API
- ✅ Análise inteligente de padrões de leitura
- ✅ Recomendações personalizadas via IA (Google Gemini)
- ✅ Interface web responsiva
- ✅ Histórico de recomendações
- ✅ Sistema de feedback

## Tecnologias Utilizadas

### Backend
- **Go (Golang)** - Linguagem principal
- **Gin** - Framework web
- **GORM** - ORM para banco de dados
- **SQLite** - Banco de dados
- **Google APIs** - Sheets e Drive
- **Google Gemini AI** - Inteligência artificial

### Frontend
- **HTML5 + CSS3 + JavaScript**
- **Bootstrap 5** - Framework CSS
- **Font Awesome** - Ícones

## Estrutura da Planilha

O sistema espera uma planilha com as seguintes colunas:

| Coluna | Título | Descrição | Obrigatório |
|--------|--------|-----------|-------------|
| A | Título | Nome do livro | ✅ |
| B | Autor | Autor principal | ✅ |
| C | Gênero | Categoria/gênero do livro | ✅ |
| D | Data Leitura | Quando foi lido (DD/MM/YYYY) | ✅ |
| E | Avaliação | Nota de 1-5 estrelas | ✅ |
| F | Observações | Comentários opcionais | ❌ |

## Configuração

### Pré-requisitos

1. **Go 1.21+** instalado
2. **Credenciais do Google Cloud Platform**:
   - Google OAuth 2.0 (Client ID e Client Secret)
   - Google Gemini AI API Key

### Instalação

1. Clone o repositório:
```bash
git clone <repo-url>
cd book-suggest-ai
```

2. Copie o arquivo de configuração:
```bash
cp .env.example .env
```

3. Configure suas credenciais no arquivo `.env`:
```env
GOOGLE_CLIENT_ID=your_google_client_id_here
GOOGLE_CLIENT_SECRET=your_google_client_secret_here
GEMINI_API_KEY=your_gemini_api_key_here
SESSION_SECRET=your_random_session_secret_here
```

4. Instale as dependências:
```bash
go mod tidy
```

5. Execute a aplicação:
```bash
go run cmd/main.go
```

6. Acesse `http://localhost:8080`

## Como Usar

1. **Login**: Faça login com sua conta Google
2. **Selecione a Planilha**: Escolha a planilha que contém seu histórico de leitura
3. **Análise**: O sistema analisará seus padrões de leitura
4. **Recomendações**: Receba 3 recomendações personalizadas com justificativas
5. **Feedback**: Avalie as recomendações para melhorar futuras sugestões

## Estrutura do Projeto

```
book-suggest-ai/
├── cmd/
│   └── main.go                 # Ponto de entrada da aplicação
├── internal/
│   ├── config/                 # Configurações
│   ├── handlers/               # Controladores HTTP
│   ├── middleware/             # Middlewares
│   ├── models/                 # Modelos de dados
│   └── services/               # Lógica de negócio
├── pkg/
│   └── database/               # Configuração do banco de dados
├── web/
│   ├── static/                 # Arquivos estáticos (CSS, JS)
│   └── templates/              # Templates HTML
├── configs/                    # Arquivos de configuração
├── .env.example               # Exemplo de configuração
├── go.mod                     # Dependências Go
└── README.md                  # Este arquivo
```

## API Endpoints

### Autenticação
- `GET /auth/login` - Página de login
- `GET /auth/google` - Iniciar OAuth Google
- `GET /auth/google/callback` - Callback OAuth
- `POST /auth/logout` - Fazer logout

### Web Interface
- `GET /` - Página inicial
- `GET /dashboard` - Dashboard (requer autenticação)

### API REST
- `GET /api/v1/spreadsheets` - Listar planilhas do usuário
- `GET /api/v1/spreadsheets/:id/analyze` - Analisar planilha
- `POST /api/v1/recommendations` - Gerar recomendações
- `GET /api/v1/recommendations/history` - Histórico de recomendações

## Configuração do Google Cloud

### 1. Google OAuth 2.0

1. Acesse o [Google Cloud Console](https://console.cloud.google.com/)
2. Crie um novo projeto ou selecione um existente
3. Ative as APIs:
   - Google Sheets API
   - Google Drive API
4. Configure a tela de consentimento OAuth
5. Crie credenciais OAuth 2.0:
   - Tipo: Aplicação Web
   - URIs de redirecionamento: `http://localhost:8080/auth/google/callback`

### 2. Google Gemini AI

1. Acesse [Google AI Studio](https://makersuite.google.com/)
2. Gere uma API Key para o Gemini
3. Adicione a chave no arquivo `.env`

## Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.

## Suporte

Para questões e suporte, abra uma issue no repositório do GitHub.
