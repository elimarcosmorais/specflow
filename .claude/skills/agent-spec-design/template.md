---
status: rascunho   # rascunho | aprovado
---

# Design — <nome da feature>

- **Feature:** `NNN-<nome-da-feature>` · **Versão:** vN · **Data:** AAAA-MM-DD
- **Origem:** [prd.md](prd.md)
- **Design system:** [design-system.md](../../../design-system.md) | inexistente
- **Plataforma(s):** web | mobile | desktop
- **Referências:** <links de Figma, caminhos de mockups> | nenhuma

---

## 1. Direção visual

<!-- LLM-ONLY: 2 a 4 linhas: como esta feature se apresenta dentro da identidade do produto.
  Ex.: "tabela densa como elemento central; ações destrutivas sempre com confirmação". Exceção ao
  tema do design system (ex.: tela sempre escura) entra aqui, com o motivo. -->

## 2. Componentes

### 2.1 Reusados

<!-- LLM-ONLY: do design-system.md §7 ou da biblioteca base, com a mesma referência de import. -->

| Componente | Referência (import) | Onde é usado nesta feature |
|------------|---------------------|----------------------------|
|            |                     |                            |

### 2.2 Novos

<!-- LLM-ONLY: só o que comprovadamente NÃO existe no design system nem na biblioteca. Token novo
  também entra aqui. "Promover" = sim quando é reutilizável entre features (confirmado com o
  usuário). Nenhum → "Nenhum.". -->

| Item | Tipo | Por que o existente não atende | Promover ao design system? |
|------|------|--------------------------------|----------------------------|
|      | componente \| token |                 | sim \| não                 |

## 3. Telas

<!-- LLM-ONLY: uma subseção por tela, janela, modal ou painel da feature, descrita de forma que o
  executor implemente sem adivinhar: regiões, hierarquia, o que domina a tela, ordem de leitura,
  componentes usados. Wireframe ASCII quando ajudar. A rota técnica é da Tech Spec. -->

### 3.1 <Nome da tela>

- **Propósito:** <1 frase>
- **Histórias:** US-01, US-02
- **Como se chega:** <menu, botão, link de outra tela>
- **Layout:** <regiões e hierarquia: cabeçalho, conteúdo, ações; o que domina a tela>
- **Componentes:** <da seção 2>

```
<wireframe ASCII opcional>
```

## 4. Estados visuais

<!-- LLM-ONLY: o coração do documento. Cada célula é CONCRETA: "skeleton de 3 linhas no lugar da
  tabela", "banner no topo com ação Tentar novamente", nunca "mostrar loading" ou "mostrar erro".
  Parte do padrão do design system §8; desvio dele diz o motivo. Use o vocabulário da biblioteca
  do projeto (toast, snackbar, banner, dialog). Mobile ou app que funciona offline: acrescente a
  coluna Offline. Estado que não existe na tela → "—". -->

| Tela | Carregando | Sucesso | Erro | Vazio |
|------|------------|---------|------|-------|
|      |            |         |      |       |

### 4.1 Mensagens

<!-- LLM-ONLY: texto literal que o usuário lê. Mensagens de erro vindas de regra do PRD citam a RN. -->

| Situação | Mensagem (literal) | Ação oferecida |
|----------|--------------------|----------------|
|          |                    |                |

## 5. Tamanhos de tela

<!-- LLM-ONLY: use as faixas do design system §5 pelos nomes de lá. Só o que muda NESTA feature.
  Nada muda além do padrão → "Segue o design system §5." -->

| Faixa | O que muda nas telas desta feature |
|-------|------------------------------------|
|       |                                    |

## 6. Interações e motion

<!-- LLM-ONLY: só o que importa para a feature: transições, feedback de ação, gestos (mobile),
  atalhos de teclado (desktop/web). Sem linguagem de animação no projeto, mantenha o mínimo.
  Nada além do padrão → "N/A — usa as transições padrão do design system." -->

| Interação | Comportamento visual |
|-----------|----------------------|
|           |                      |

## 7. Acessibilidade visual

<!-- LLM-ONLY: só o que é específico desta feature (par de cores crítico, foco num componente
  denso, área de toque numa lista). O resto: "Segue o design system §9." -->

## 8. Assets

<!-- LLM-ONLY: ícones, ilustrações, imagens e fontes NOVOS, com a origem. Nenhum → "Nenhum.". -->

| Asset | Tipo | Origem (biblioteca do projeto, a produzir, link) |
|-------|------|--------------------------------------------------|
|       |      |                                                  |

## 9. Pontos em aberto

- **Dependências de produto:** <decisão visual bloqueada por algo que o PRD não decidiu> | nenhuma
- **Para a Tech Spec:** <decisão técnica percebida aqui, fora deste documento> | nenhuma

## 10. Próximo passo

- **Próxima fase:** Tech Alignment → `/agent-spec-tech-alignment docs/specs/features/NNN-<feature>/vN/prd.md`
- **ADR a registrar antes:**
  - `/agent-spec-adr "<título da decisão>"` — <o que ela decide, em poucas palavras>

<!-- LLM-ONLY: comandos prontos para copiar, com o caminho real. Na mudança pequena, a próxima fase
  é a Tech Spec, com o caminho do pre-refinement.md. ADR aqui = decisão de design transversal com
  trade-off real; sem nenhuma → "**ADR a registrar antes:** nenhuma". -->
