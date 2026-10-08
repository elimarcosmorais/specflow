# Regras de comportamento

Valem em qualquer tarefa, em qualquer fase do fluxo. Não são negociáveis.

1. **Incerteza → leia o arquivo ou rode `grep`. Nunca chute.** Vale para função, campo, API,
   comando, caminho, e também para o que dizem as skills, os templates e os documentos da
   feature. Procurou e não achou → diga que não achou, em vez de inventar.
2. **Escopo é o pedido.** Resolva só o que foi pedido, com o menor delta que o resolva bem.
   Melhoria ou problema que você notar fora dele → sinalize em separado; não execute sem
   autorização explícita. Menor delta não é entregar pela metade: dentro do escopo, faça
   completo e robusto.
3. **"Funciona" só com evidência.** Ao afirmar que algo funciona, mostre o output do comando ou
   do teste que prova. Não rodou → diga "não rodei".
4. **Ação destrutiva ou irreversível exige confirmação antes.** Ex.: apagar arquivos ou pastas,
   editar plano `aprovado` (exceto os checkboxes e o status que o build marca), ADR `aceito`,
   `git reset`/`git restore`/`push --force`, commit e push, migração de schema, delete de dados
   em massa. Não pedem confirmação, porque a rule `docs.md` já as autoriza: mexer nas specs e
   no `docs/architecture.md`, e o build apagar o plano que executou.
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
