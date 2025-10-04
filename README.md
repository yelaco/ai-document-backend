# AI Document Backend

A powerful NestJS-based backend application that enables users to upload PDF documents and interact with them using AI-powered question-answering capabilities. The system leverages Retrieval-Augmented Generation (RAG) to provide accurate, context-aware responses based on document content.

## 🚀 Features

- **PDF Document Processing**: Upload and parse PDF documents with automatic text extraction
- **AI-Powered Q&A**: Ask questions about uploaded documents and receive intelligent answers
- **Vector Search**: Advanced semantic search using ChromaDB for document retrieval
- **Multiple AI Providers**: Support for both Google Gemini and Anthropic Claude models
- **User Authentication**: Secure JWT-based authentication system
- **Real-time Streaming**: Server-sent events for streaming AI responses
- **Docker Support**: Complete containerized deployment with Docker Compose
- **Observability**: Built-in OpenTelemetry integration with Zipkin tracing

## 🏗️ Architecture

The application implements a modern RAG (Retrieval-Augmented Generation) architecture:

### Indexing Pipeline

1. **Document Upload** - PDF files are uploaded via REST API
2. **Text Extraction** - Content is extracted using pdf-parse
3. **Text Chunking** - Documents are split into semantic chunks using LangChain
4. **Embedding Generation** - Text chunks are converted to vectors using AI models
5. **Vector Storage** - Embeddings are stored in ChromaDB for fast retrieval

### Query Pipeline

1. **Question Processing** - User questions are embedded using the same AI model
2. **Semantic Search** - ChromaDB finds the most relevant document chunks
3. **Context Construction** - Retrieved chunks are combined with the question
4. **AI Generation** - AI model generates answers based on the provided context
5. **Streaming Response** - Answers are streamed back to the client in real-time

## 🛠️ Tech Stack

- **Framework**: NestJS (Node.js)
- **Runtime**: Bun
- **Database**: PostgreSQL with TypeORM
- **Vector Database**: ChromaDB
- **AI Models**: Google Gemini, Anthropic Claude
- **Authentication**: JWT
- **Text Processing**: LangChain, pdf-parse
- **Containerization**: Docker & Docker Compose
- **Observability**: OpenTelemetry, Zipkin
- **Language**: TypeScript

## 📋 Prerequisites

- [Bun](https://bun.sh/) (latest version)
- [Docker](https://www.docker.com/) and Docker Compose
- AI API Key (Google Gemini or Anthropic Claude)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd ai-document-backend
```

### 2. Environment Setup

Copy the example environment file and configure your settings:

```bash
cp .env.example .env.development
```

Edit `.env.development` with your configuration:

```env
# App
APP_ENV=development

# PostgreSQL
DB_HOST=db
DB_USER=aidocument
DB_PASS=password
DB_NAME=aidocument

# AI Service Configuration
AI_SERVICE_TYPE=gemini  # or 'anthropic'
GEMINI_API_KEY=your_gemini_api_key_here
# ANTHROPIC_API_KEY=your_anthropic_api_key_here
```

### 3. Install Dependencies

```bash
bun install
```

### 4. Start the Application

**Development Mode** (with hot reload):

```bash
bun run compose:dev
```

**Production Mode**:

```bash
bun run compose:prod
```

The application will be available at:

- **API**: http://localhost:3000
- **ChromaDB**: http://localhost:8000
- **Zipkin Tracing**: http://localhost:9411

### 5. Stop the Application

```bash
# Development
bun run compose:down:dev

# Production
bun run compose:down:prod
```

## 🧪 Development

### Running Tests

```bash
# Unit tests
bun test

# End-to-end tests
bun run test:e2e

# Test coverage
bun run test:cov
```

### Code Quality

```bash
# Lint code
bun run lint

# Format code
bun run format
```

### Local Development (without Docker)

1. Start PostgreSQL and ChromaDB separately
2. Update environment variables to point to local services
3. Run the application:

```bash
bun run start:dev
```

## 🏗️ Project Structure

```
src/
├── ai/                 # AI service implementations
│   ├── ai.service.ts   # AI service interface
│   ├── gemini.service.ts
│   └── anthropic.service.ts
├── auth/               # Authentication module
├── documents/          # Document processing and Q&A
├── users/              # User management
├── shared/             # Shared utilities and interceptors
├── app.module.ts       # Main application module
└── main.ts            # Application entry point

config/                 # Configuration files
doc/                   # Documentation
test/                  # End-to-end tests
```

## 🔧 Configuration

### Environment Variables

| Variable            | Description                           | Default       |
| ------------------- | ------------------------------------- | ------------- |
| `APP_ENV`           | Application environment               | `development` |
| `APP_PORT`          | Server port                           | `3000`        |
| `DB_HOST`           | PostgreSQL host                       | `db`          |
| `DB_PORT`           | PostgreSQL port                       | `5432`        |
| `DB_USER`           | Database username                     | -             |
| `DB_PASS`           | Database password                     | -             |
| `DB_NAME`           | Database name                         | -             |
| `AI_SERVICE_TYPE`   | AI provider (`gemini` or `anthropic`) | `gemini`      |
| `GEMINI_API_KEY`    | Google Gemini API key                 | -             |
| `ANTHROPIC_API_KEY` | Anthropic Claude API key              | -             |

### AI Service Configuration

The application supports multiple AI providers. Configure your preferred provider in the environment:

- **Google Gemini**: Set `AI_SERVICE_TYPE=gemini` and provide `GEMINI_API_KEY`
- **Anthropic Claude**: Set `AI_SERVICE_TYPE=anthropic` and provide `ANTHROPIC_API_KEY`

## 📊 Monitoring and Observability

The application includes comprehensive observability features:

- **Logging**: Structured logging with request interceptors
- **Tracing**: OpenTelemetry integration with Zipkin
- **Health Checks**: Database and ChromaDB connectivity monitoring

Access Zipkin UI at http://localhost:9411 to view distributed traces.

## 🚀 Deployment

### Production Deployment

1. Configure production environment variables in `.env.production`
2. Deploy using Docker Compose:

```bash
bun run compose:prod
```

### Container Configuration

The application runs in a multi-container setup:

- **app**: NestJS application
- **db**: PostgreSQL database
- **vector_db**: ChromaDB vector database
- **otel-collector**: OpenTelemetry collector
- **zipkin**: Distributed tracing UI

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Author

Created by [yelaco](https://github.com/yelaco)

## 🔗 Related Documentation

- [Workflow Documentation](doc/workflow.md) - Detailed explanation of the RAG pipeline
- [NestJS Documentation](https://docs.nestjs.com/)
- [ChromaDB Documentation](https://docs.trychroma.com/)
- [Google Gemini API](https://ai.google.dev/)
- [Anthropic Claude API](https://docs.anthropic.com/)
