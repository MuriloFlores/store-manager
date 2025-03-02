# Store Manager

Store Manager é uma aplicação web desenvolvida em Go para gerenciamento de lojas e inventários. O projeto demonstra boas práticas de desenvolvimento utilizando Clean Architecture, GORM para acesso a dados (PostgreSQL), Swagger para documentação interativa e Zap para logging.

# Recursos

- Gerenciamento de Produtos: Criação, atualização, listagem, busca por IDs e exclusão.
- Gerenciamento de Matérias-Primas: Operações similares para gerenciar matérias-primas.
- API RESTful Documentada: Endpoints documentados com Swagger para facilitar a integração.
* Arquitetura Limpa: Código organizado em camadas (domínio, aplicação, infraestrutura).

# Instalação

Pré-requisitos:

- Go (1.16+)
- PostgreSQL
- Swag para geração da documentação

#  Passos

Clone o repositório:

```sh
git clone https://github.com/seu-usuario/store-manager.git
cd store-manager
```

# Configure as variáveis de ambiente:

Crie um arquivo config.env conforme necessário (não versionado).

# Instale as dependências:

```
go mod tidy
```

# Gere a documentação Swagger:

```
swag init -g cmd/store-manager/main.go
```
# Execute a aplicação:

```
    go run cmd/store-manager/main.go
```
# Uso

A API estará disponível em: http://localhost:8080
Acesse a documentação interativa em: http://localhost:8080/swagger/index.html

# Contribuição

Contribuições são bem-vindas!
Por favor, abra uma issue ou envie um pull request para sugerir melhorias ou correções.
# Licença

Este software é licenciado sob a MIT License with Commons Clause.
Para uso comercial, favor contatar o autor para obter uma licença apropriada.

------------------------------------------------

# Store Manager

Store Manager is a web application developed in Go for managing stores and inventories.The project demonstrates best practices using Clean Architecture, GORM for data access (PostgreSQL), Swagger for interactive documentation, and Zap for logging.

# Features

- Product Management: Creation, update, listing, searching by IDs, and deletion.
- Raw Material Management: Similar operations for managing raw materials.
- Documented RESTful API: Endpoints documented with Swagger to facilitate integration.
- Clean Architecture: Code organized in layers (domain, application, infrastructure).

# Installation

Prerequisites:

- Go (1.16+)
- PostgreSQL
- Swag for documentation generation

# Steps

Clone the repository:

```sh
git clone https://github.com/seu-usuario/store-manager.git
cd store-manager
```

# Configure environment variables:

Create a config.env file as needed (do not commit it to version control).

Install dependencies:
```
go mod tidy
```
# Generate Swagger documentation:

```
swag init -g cmd/store-manager/main.go
```

# Run the application:
```
    go run cmd/store-manager/main.go
```
# Usage

The API is available at: http://localhost:8080
Access the interactive documentation at: http://localhost:8080/swagger/index.html

# Contribution

Contributions are welcome!
Please open an issue or submit a pull request to suggest improvements or fixes.

# License

This software is licensed under the MIT License with Commons Clause.
For commercial use, please contact the author for an appropriate commercial license.


