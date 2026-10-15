# Regras de comportamento

Valem em qualquer tarefa, em qualquer fase do fluxo. Não são negociáveis.

1. **Incerteza → leia o arquivo ou rode `grep`. Nunca chute.** Vale para função, campo, API,
   comando, caminho, e também para o que dizem as skills, os templates e os documentos da
   feature. Procurou e não achou → diga que não achou, em vez de inventar.
2. **Escopo é o pedido.** Resolva só o que foi pedido, com o menor delta que o resolva bem.
   Melhoria ou problema que você notar fora dele → sinalize em separado; não execute sem
   autorização explícita. Menor delta não é entregar pela metade: dentro do escopo, faça
   completo e robusto.

   **Exceção: correção de impacto pequeno**, ao implementar pelo `/spec-change` ou pelo
   `/spec-build`. Pode corrigir sem perguntar, e informar no relatório como corrigida, o
   problema fora do escopo que cumpre **todos** estes critérios:
   - **óbvio:** só há um jeito certo de corrigir (erro de digitação, import ou link quebrado,
     comentário ou spec que contradiz o código, aviso de lint, código morto evidente);
   - **localizado:** poucas linhas, num ou dois arquivos;
   - **sem efeito para quem usa:** não muda comportamento de que alguém depende nem contrato
     (API, formato de dados, schema, tela);
   - **fora de área sensível:** autenticação, autorização, pagamentos, segredos, dados pessoais;
   - **provado:** os testes, o build e o lint que a tarefa já roda cobrem o que mudou.

   **ADR fica fora desta exceção.** Mesmo o erro óbvio numa ADR só é corrigido sem confirmação
   pela via da Exceção 2 — um `spec-fixer`, na varredura —, porque é lá que a rule `docs.md`
   ("Correção verificável") a autoriza. Notou um → não corrija aqui: anote para a varredura.

   **Exceção 2: item isolado delegado**, só na varredura fora do escopo — `/spec-change`,
   Etapa 3, item 6, e `/spec-build`, Etapa 3, item 2 —, depois de a mudança pedida estar feita e
   provada. Pode resolver, num `spec-fixer` por item, o problema que falha a exceção acima **só
   por tamanho, por precisar de teste próprio ou por ser numa ADR** (regra 4) e cumpre **todos**
   estes critérios:
   - **uma correção certa só:** nenhuma escolha de comportamento ou de produto que seja do
     usuário;
   - **isolado:** lista fechada de arquivos, disjunta da mudança pedida e dos outros itens (dois
     fixers nunca escrevem no mesmo arquivo);
   - **sem contrato novo:** não muda API, formato de dados, schema, migração nem tela;
   - **fora de área sensível:** autenticação, autorização, pagamentos, segredos, dados pessoais;
   - **provável:** os testes, o build e o lint que o projeto já tem provam que entrou verde.

   **Item só de documentação** — hoje, a correção verificável numa ADR, que a rule `docs.md`
   autoriza — não tem teste que o prove: a prova é a conferência mecânica que o `Pronto quando`
   descreve (abrir o link e ver que resolve, conferir que o arquivo citado existe). Suíte e
   build seguem valendo como gate, porque ele não pode deixá-los vermelhos.

   No máximo 3 por rodada, **um nível só** (o que um `spec-fixer` notar nunca vira outro item na
   mesma rodada) e gate no fim: suíte ou build vermelho → desfaz-se só o que os fixers tocaram.

   Falhou um critério, ou ficou na dúvida → só sinalize, com a recomendação e o comando pronto.
3. **"Funciona" só com evidência.** Ao afirmar que algo funciona, mostre o output do comando ou
   do teste que prova. Não rodou → diga "não rodei".

   **Teste que pulou não é evidência.** `SKIP`, `skipped`, `pending`, `todo`, `it.skip`,
   `t.Skip`, `@pytest.mark.skip`, suíte que não achou o pacote, o servidor ou a porta que
   precisava: a suíte fica verde e o comportamento **não foi exercido**. Conte o pulado junto
   com o não rodado, sempre com o motivo e o nome do teste, e nunca o apresente como prova. O
   item que ele deveria provar continua **sem prova**.
4. **Ação destrutiva ou irreversível exige confirmação antes.** Ex.: apagar arquivos ou pastas,
   editar plano `aprovado` (exceto os checkboxes e o status que o build marca), ADR `aceito` ou
   `substituído`, `git reset`/`git restore`/`push --force`, commit e push, migração de schema,
   delete de dados em massa. Não pedem confirmação, porque a rule `docs.md` já as autoriza:
   mexer nas specs e no `docs/architecture.md`, a correção verificável numa ADR `proposto` ou
   `aceito` na varredura, e o build apagar o plano que executou. Também não pedem
   confirmação o build descartar a rodada de sugestões que a revisão reprovou e a varredura
   desfazer o item que deixou o gate vermelho: cada uma desfaz só o que ela mesma acabou de
   fazer, nos arquivos daquele item, e o que já estava aprovado não é tocado.
5. **Precisou de admin → eleve só aquele comando.** Não passe a tarefa ao usuário nem rode a
   sessão inteira como admin. Antes, diga no chat o comando exato e por quê: o aviso do UAC não
   mostra o comando. O processo elevado não devolve a saída, então grave-a num log e leia-o
   depois (é a evidência da regra 3). Usuário recusou o UAC → diga o comando e peça que ele o
   rode.

   ```powershell
   $p = Start-Process pwsh -Verb RunAs -Wait -PassThru -WindowStyle Hidden -ArgumentList '-NoProfile','-Command',"& { <comando> } *> '<log>'; exit `$LASTEXITCODE"
   $p.ExitCode; Get-Content '<log>'
   ```

Como evoluir esta lista: errou em algo que uma regra já cobria → reforce a redação dela; errou
em algo novo → acrescente uma regra.
