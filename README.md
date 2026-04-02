# 🛒 Store Manager API

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgresql-4169e1?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Clean Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-4ea94b?style=for-the-badge)
![Testes](https://img.shields.io/badge/Testes-TDD-blueviolet?style=for-the-badge)

**O projeto ainda se encontra em construção**

O **Store Manager** é uma API RESTful robusta desenvolvida em Go para o gerenciamento de ponta a ponta de um ecossistema comercial. O sistema orquestra produtos, controle transacional de estoque, processamento de pedidos (checkout) e fluxo logístico de entregas.

## 🎯 Foco Arquitetural e Técnico

Este projeto foi construído com foco extremo em engenharia de software de alta qualidade:

* **Clean Architecture (Ports & Adapters):** Separação estrita entre a regra de negócio (Core) e a infraestrutura (Banco de Dados, HTTP, Cache).
* **Domain-Driven Design (DDD):** Modelagem tática focada nas entidades comerciais (Product, Stock, Order, Delivery).
* **Testes (TDD):** Cobertura de testes unitários e de integração para garantir a confiabilidade das regras de negócio.
* **Controle de Concorrência & Transações (ACID):** Tratamento rigoroso de concorrência no estoque para evitar o problema de "duplo gasto".
* **CI/CD:** Automação de testes e deploy utilizando GitHub Workflows.

## 🚀 Funcionalidades Principais (Use Cases)

* **Catálogo:** Gerenciamento completo de Produtos (CRUD).
* **Inventário Segregado:** Controle de saldo disponível vs. saldo reservado.
* **Motor de Pedidos (Checkout):** Validação atômica de disponibilidade de estoque e cálculo de totais.
* **Cache e Performance:** Utilização do Redis para otimizar consultas frequentes.
* **Logística:** Rastreamento de mudança de status da entrega.

## 🛠️ Tecnologias Utilizadas

* **Linguagem:** Go (1.25.0)
* **Banco de Dados Relacional:** PostgreSQL
* **Cache & Filas:** Redis
* **Infraestrutura:** Docker e Docker Compose
* **CI/CD:** GitHub Actions (Workflows)

## 🏁 Como executar o projeto localmente

**Pré-requisitos:** Go, Docker e Docker Compose instalados.

1. Clone o repositório:
`git clone https://github.com/SeuUsuario/store-manager.git`
`cd store-manager/backend`

2. Suba a infraestrutura via Docker:
`docker-compose up -d`

3. Baixe as dependências e execute os testes:
`go mod tidy`
`go test ./...`

4. Execute a aplicação principal:
`go run cmd/api/main.go`