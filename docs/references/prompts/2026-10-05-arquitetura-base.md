# Prompt Pre Refinement

Quero iniciar a construção da arquitetura base de um sistema chamado "Dev Tools", será um programa executável para plataforma Windows com o objetivo de gerenciar diversas ferramentas de desenvolvimento e servir também como servidor web em localhost. Será inspirado em outros programas.

Nesta primeira estapa, o objetivo é definir a estrutura inicial do projeto e desing de código e arquitetura que será utilizada para todo o desenvolvimento do sistema.

## Programas para referências

Os programas abaixo servem como referência para ciração do Dev Tools.

- FlyEnv: https://flyenv.com (Fortemente inspirado no FlyEnv)
- Laragon: https://laragon.org
- Laravel Herd: https://herd.laravel.com/windows

## Stack tecnológica inicial (sugerida)

**Go 1.25+** (motor e CLI) · **Wails v2** (app desktop, WebView2) · front em **React 19 + TypeScript + Vite** · persistência em **SQLite** (`modernc org/sqlite`, sem CGO) · testes com `go test` contra um servidor local que injeta falhas.

## Esquema de diretório do Dev Tools após a instalação (sugerida)

```
C:\dev-tools\
│
├── projects/
│   ├── web/
│   ├── android/
│   └── windows/
│
└── programs/
    ├── README.md
    ├── dev-tools.exe
    ├── backup/
    ├── cache/
    ├── data/
    ├── lang/
    ├── logs/
    ├── temp/
    ├── sessions/
    ├── ...
    └── bin/
        ├── apache/
        ├── go/
        ├── mariabd/
        ├── mysql/
        ├── nodejs/
        |   └── nodejs-vX.X.X
        ├── php/
        |   └── php-vX.X.X
        ├── redis/
        └── .../
```

## Descrição do sistema de recursos e funcionalidade do sistema

A arquitetura inicial deve considerar que o sistema terá:

- Conexão com banco de dados SQLite
- Serviço para identificar a saúde da aplicação
- Estrutura de logs eficiente para o sistema e para o desenvolvimento
- Estrutura de testes unitários e de integração
- Estrutura de documntação de API
- Definição de um framework ou Http Server para a API
- Divisão de responsabilidade entre camadas de negócio, infraestrutura e testes
- O sistema deve ter suporte a mais de 1 idioma: pt-br e en
- Deverá ter opção para selecionar tema Escuro, Claro ou Sistema.
- Ao executar o instalador do programa vai pedir para escolher um local/diretório para instalação, semelhante o instalador do Laragon faz, sugerindo um local de instalação no diretório "C:\dev-tools\" do windows.

A ideia é sempre a criação de uma arquitetura limpa e organizada para o sistema, deve ser criado a arquitetura base para o desenvolvimento e um handler completo de exemplo para o resto do sistema seguir.
