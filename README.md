# Bot Builder Engine

Este é o repositório do Bot Builder Engine, uma aplicação escrita em Go para construção de bots.

## Estrutura do Projeto

A estrutura do projeto é organizada da seguinte forma:

```
bot-builder-engine/
├── application
│   ├── api_whats_application.go
│   └── runner_application.go
├── data
│   ├── engine_runner.go
│   └── whats_node.go
├── go.mod
├── infra
│   ├── config
│   │   └── cors.go
│   └── server.go
├── README.md
├── repository
│   └── engine_local_storage.go
├── routes
│   └── api_whats_router.go
├── services
│   └── api_whats_service.go
└── utils
    └── utils.go
```

- **go.mod**: O arquivo de manifesto do módulo Go, que especifica as dependências do projeto.
- **infra/**: Pasta que inclui o arquivo `server.go`, onde a função `main` está localizada.

## Pré-requisitos

Para executar este projeto, você precisará ter o Go instalado em seu sistema. Você pode encontrar as instruções de instalação em [golang.org](https://golang.org/doc/install).

## Executando a Aplicação

Para executar a aplicação, siga estas etapas:

1. Clone este repositório:

2. Navegue até o diretório do projeto:

    ```bash
    cd bot-builder-engine
    ```

3. Execute o comando para iniciar o servidor:

    ```bash
    go run infra/server.go
    ```

Isso iniciará o servidor da aplicação.
