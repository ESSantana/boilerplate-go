# 🛠️ Go Backend Boilerplate (WIP)

Este repositório está em desenvolvimento e tem como objetivo se tornar um **boilerplate moderno para backends em Go**, com funcionalidades comuns para projetos web como autenticação, gerenciamento de clientes e integração com gateways de pagamento.

🔄 **Código em processo de refatoração.**  
Baseado em um projeto anterior descontinuado:  
🔗 [ghcr.io/ESSantana/api](https://github.com/application-ellas)

---

## 🚧 Status: Work in Progress

Este projeto está sendo modificado e estruturado para servir como base reutilizável. Abaixo, os próximos passos planejados:

### ✅ Objetivos (TODO)

- [ ] Integração com gateway de pagamento (**Stripe** ou **Mercado Pago**)
- [ ] Envio de notificações por e-mail (**AWS SES**, **Resend** ou **SendGrid**)
- [ ] Configuração de produção para **Kubernetes**
- [ ] Ambiente de desenvolvimento com **Docker Compose** ou **Kubernetes minimalista**

---

## 🧰 Tecnologias e Ferramentas (previstas)

- **Golang** (Go)
- **Docker / Docker Compose**
- **Kubernetes**
- **Stripe / Mercado Pago**
- **AWS SES / Resend / SendGrid**
- **MySQL (ou PostgreSQL) / Redis(ou Valkey)** (dependendo da necessidade futura)

---

## 📦 Objetivo Final

Criar um **template limpo, escalável e pronto para produção**, com:

- Boas práticas de arquitetura
- Configuração CI/CD amigável
- Modularização de recursos
- Suporte a ambientes locais e em nuvem

---

## 📂 Estrutura do Projeto (em breve)

```bash
.
├── cmd/               # Entrada da aplicação
├── internal/          # Lógica interna (domínios, serviços, handlers)
├── packages/          # Pacotes reutilizáveis
├── deployments/       # Kubernetes YAMLs
├── docker/            # Configurações Docker
├── .env.example       # Arquivo de exemplo de variáveis
└── README.md
