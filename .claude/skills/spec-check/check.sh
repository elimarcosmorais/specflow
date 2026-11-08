#!/bin/sh
# spec-check: o portão de testes do SpecFlow. Decide o que rodar, roda fora da sessão do modelo,
# guarda o que passou pelo hash do código e grava o resultado. Documentação: SKILL.md, ao lado.
#
# Uso:
#   sh check.sh passo '<comando de teste>' ['<comando>' ...]   testes do passo, build e lint
#   sh check.sh mudanca [--arquivos <arquivo> ...]             suíte, build, lint; integração e tela por impacto
#   sh check.sh versao                                         tudo inteiro, sem cache
#   sh check.sh esperar                                        continua esperando a execução em andamento
#   sh check.sh auditar                                        confere o spec-check.json com o repo
#
# Saída: 0 verde · 1 vermelho · 2 verde com teste pulado · 3 ainda rodando · 4 erro de uso ou config
set -u

ESPERA=${SPEC_CHECK_ESPERA:-540}
CONFIG=.claude/spec-check.json
TAB=$(printf '\t')

case $0 in
  /* | ?:*) self=$0 ;;
  *) self=$PWD/$0 ;;
esac

raiz=$(git rev-parse --show-toplevel 2>/dev/null) || {
  echo "spec-check: fora de um repositório git" >&2
  exit 4
}
cd "$raiz" || exit 4
base=$(git rev-parse --absolute-git-dir)/spec-check
mkdir -p "$base/cache" "$base/execucoes"

uso() {
  sed -n '5,10p' "$self" | sed 's/^# *//' >&2
  exit 4
}

# Lê o JSON plano da config (só texto e lista de textos) e escreve "chave<TAB>valor" por linha.
ler_config() {
  awk '
    function erro(m) { print "#ERRO\t" m; exit 1 }
    function ch() { return substr(texto, i, 1) }
    function pula(  c) {
      while (i <= n) {
        c = ch()
        if (c == " " || c == "\n" || c == "\t" || c == "\r") i++
        else break
      }
    }
    function str(  s, c) {
      i++
      s = ""
      while (i <= n) {
        c = ch()
        i++
        if (c == "\"") return s
        if (c == "\\") {
          c = ch()
          i++
          if (c == "n" || c == "t") c = " "
        }
        s = s c
      }
      erro("texto sem fechar")
    }
    { texto = texto $0 "\n" }
    END {
      n = length(texto)
      i = index(texto, "{")
      if (!i) erro("sem objeto JSON")
      i++
      while (1) {
        pula()
        c = ch()
        if (c == "}") break
        if (c == ",") { i++; continue }
        if (c != "\"") erro("esperava o nome de uma chave")
        k = str()
        pula()
        if (ch() != ":") erro("faltou : depois de " k)
        i++
        pula()
        c = ch()
        if (c == "\"") {
          print k "\t" str()
        } else if (c == "[") {
          i++
          while (1) {
            pula()
            c = ch()
            if (c == "]") { i++; break }
            if (c == ",") { i++; continue }
            if (c != "\"") erro(k " aceita só uma lista de textos")
            print k "\t" str()
          }
        } else {
          erro(k " aceita só texto ou lista de textos")
        }
      }
    }
  ' "$CONFIG"
}

cfg() {
  awk -F'\t' -v k="$1" '$1 == k { print substr($0, length(k) + 2) }' "$dir/config"
}

validar_config() {
  if grep -q '^#ERRO' "$dir/config"; then
    echo "spec-check: $CONFIG: $(grep '^#ERRO' "$dir/config" | cut -f2)" >&2
    exit 4
  fi
  while IFS= read -r linha; do
    chave=${linha%%"$(printf '\t')"*}
    valor=${linha#*"$(printf '\t')"}
    case $chave in
      suite | build | lint | integracao | tela) ;;
      integracao_alvos | tela_alvos | integracao_inteira | tela_inteira) ;;
      integracao_testes | tela_testes) ;;
      integracao_impacto | tela_impacto)
        case $valor in
          *'{alvos}'*) ;;
          *)
            echo "spec-check: $CONFIG: $chave precisa de {alvos}" >&2
            exit 4
            ;;
        esac
        ;;
      integracao_mapa | tela_mapa)
        case $valor in
          *' => '*) ;;
          *)
            echo "spec-check: $CONFIG: $chave: use 'padrão => alvo', recebi: $valor" >&2
            exit 4
            ;;
        esac
        ;;
      *)
        echo "spec-check: $CONFIG: chave desconhecida: $chave" >&2
        exit 4
        ;;
    esac
  done < "$dir/config"
}

# Hash do código como está agora, com o que ainda não está em stage. Documentação não conta: marcar
# um checkbox no plano não muda o que o teste provou. Usa um índice temporário: o stage fica intacto.
arvore_atual() {
  idx="$base/index.$$"
  real=$(git rev-parse --git-path index)
  [ -f "$real" ] && cp "$real" "$idx"
  if GIT_INDEX_FILE=$idx git add -A . 2>/dev/null; then
    GIT_INDEX_FILE=$idx git -c core.quotepath=off ls-files -s \
      | grep -Ev '	docs/|\.md$' \
      | git hash-object --stdin
  fi
  rm -f "$idx" "$idx.lock"
}

registrar() {
  printf '%s\t%s\t%s\t%s\t%s\t%s\n' "$@" >> "$dir/etapas"
}

etapa() { # <nome> <comando> <alcance>
  nome=$1
  comando=$2
  alcance=$3
  log="$dir/$nome.log"
  chave=$(printf '%s\n%s\n' "$arvore" "$comando" | git hash-object --stdin)
  if [ "$usa_cache" = 1 ] && [ -f "$base/cache/$chave" ]; then
    registrar "$nome" cache 0 "$alcance" "" "$comando"
    return 0
  fi
  inicio=$(date +%s)
  sh -c "$comando" > "$log" 2>&1 < /dev/null
  rc=$?
  dur=$(($(date +%s) - inicio))
  grep -En '(^|[^0-9])[1-9][0-9]* ([Ss]kipped|[Pp]ending|[Tt]odo)([^a-z]|$)|--- SKIP|[Ss]kipped:[[:space:]]*[1-9]|(^|[^A-Za-z])SKIPPED([^A-Za-z]|$)' "$log" > "$dir/$nome.pulos"
  if [ "$rc" -ne 0 ]; then
    geral=vermelho
    registrar "$nome" falhou "$dur" "$alcance" "$log" "$comando"
    return 1
  fi
  if [ -s "$dir/$nome.pulos" ]; then
    [ "$geral" = vermelho ] || geral=pulou
    registrar "$nome" pulou "$dur" "$alcance" "$log" "$comando"
    return 1
  fi
  if [ "$usa_cache" = 1 ] && [ -n "$arvore" ]; then
    : > "$base/cache/$chave"
  fi
  registrar "$nome" ok "$dur" "$alcance" "$log" "$comando"
  return 0
}

regex_de() {
  printf '%s\n' "$1" | sed -E \
    -e 's/[.+?()|^$]/\\&/g' \
    -e 's/\*\*/@@DUAS@@/g' \
    -e 's/\*/[^\/]*/g' \
    -e 's/@@DUAS@@/.*/g' \
    -e 's/\{x\}/([^\/]+)/g' \
    -e 's/\/$/\/.*/'
}

alvos_do_mapa() { # <tipo>
  cfg "$1_mapa" | while IFS= read -r m; do
    re=$(regex_de "${m%% => *}")
    destino=$(printf '%s' "${m#* => }" | sed -e 's/[&#\\]/\\&/g' -e 's/{x}/\\1/g')
    sed -nE "s#^$re\$#$destino#p" "$dir/arquivos"
  done | sort -u | while IFS= read -r g; do
    for f in $g; do
      [ -e "$f" ] && printf '%s\n' "$f"
    done
  done
}

lista_repo() {
  {
    git -c core.quotepath=off ls-files
    git -c core.quotepath=off ls-files -o --exclude-standard
  } | tr -d '\r' | sort -u
}

# Lê linhas do mapa na entrada e escreve "alvo<TAB>origem" para cada arquivo do <arquivo> que casa.
pares_do_mapa() { # <arquivo com a lista>
  while IFS= read -r m; do
    re=$(regex_de "${m%% => *}")
    destino=$(printf '%s' "${m#* => }" | sed -e 's/[&#\\]/\\&/g' -e 's/{x}/\\2/g')
    sed -nE "s#^($re)\$#$destino$TAB\\1#p" "$1"
  done
}

# Testes (arquivo ou pasta) que algum arquivo do repo alcança pelo mapa. Origem que é o próprio
# alvo, ou está dentro dele, não conta: com a linha e2e/{x}.spec.ts => e2e/{x}.spec.ts, ou
# tests/{x}.Tests/ => tests/{x}.Tests/, o teste só roda quando é editado.
alcancados() { # <tipo>
  cfg "$1_mapa" | pares_do_mapa "$dir/repo" | sort -u > "$dir/pares-$1"
  cut -f1 "$dir/pares-$1" | sort -u | while IFS= read -r alvo; do
    for f in $alvo; do
      [ -e "$f" ] || continue
      f=${f%/}
      awk -F'\t' -v a="$alvo" -v f="$f" '
        $1 == a && $2 != f && index($2, f "/") != 1 { achou = 1; exit }
        END { exit !achou }
      ' "$dir/pares-$1" && printf '%s\n' "$f"
    done
  done | sort -u
}

# Testes da lista que nenhuma linha do mapa alcança. Com *_alvos, quem acha os testes é a
# ferramenta do ecossistema (go list, --findRelatedTests), e não há órfão a apontar.
orfaos() { # <tipo> <arquivo com a lista>
  [ -n "$(cfg "$1_testes")" ] || return 0
  [ -z "$(cfg "$1_alvos")" ] || return 0
  cfg "$1_testes" | while IFS= read -r p; do
    grep -E "^$(regex_de "$p")\$" "$2"
  done | sort -u | while IFS= read -r f; do
    [ -e "$f" ] && printf '%s\n' "$f"
  done > "$dir/orfaos-$1.candidatos"
  [ -s "$dir/orfaos-$1.candidatos" ] || return 0
  alcancados "$1" > "$dir/orfaos-$1.alcancados"
  # Alcançado é o próprio arquivo ou uma pasta acima dele (o projeto de teste no .NET).
  awk '
    NR == FNR { a[$0] = 1; next }
    {
      ok = ($0 in a)
      p = $0
      while (!ok && sub(/\/[^\/]*$/, "", p)) if (p in a) ok = 1
      if (!ok) print
    }
  ' "$dir/orfaos-$1.alcancados" "$dir/orfaos-$1.candidatos"
}

auditar() {
  lista_repo > "$dir/repo"
  problemas="$dir/problemas"
  : > "$problemas"
  for chave in suite build lint integracao tela integracao_impacto tela_impacto integracao_alvos tela_alvos; do
    c=$(cfg "$chave")
    [ -n "$c" ] || continue
    prog=$(printf '%s\n' "$c" | awk '{ for (i = 1; i <= NF; i++) if ($i !~ /^[A-Za-z_][A-Za-z0-9_]*=/) { print $i; exit } }')
    [ -n "$prog" ] || continue
    command -v "$prog" > /dev/null 2>&1 \
      || echo "$chave: o comando '$prog' não existe nesta máquina" >> "$problemas"
  done
  for tipo in integracao tela; do
    cfg "${tipo}_mapa" | while IFS= read -r m; do
      if ! grep -Eq "^$(regex_de "${m%% => *}")\$" "$dir/repo"; then
        echo "${tipo}_mapa: '$m' não casa com nenhum arquivo do repo"
        continue
      fi
      existe=$(printf '%s\n' "$m" | pares_do_mapa "$dir/repo" | cut -f1 | sort -u | while IFS= read -r g; do
        for f in $g; do
          if [ -e "$f" ]; then
            echo sim
            break 2
          fi
        done
      done)
      [ -n "$existe" ] || echo "${tipo}_mapa: '$m' não leva a nenhum teste que existe"
    done >> "$problemas"
    cfg "${tipo}_testes" | while IFS= read -r p; do
      grep -Eq "^$(regex_de "$p")\$" "$dir/repo" \
        || echo "${tipo}_testes: '$p' não casa com nenhum arquivo do repo"
    done >> "$problemas"
    orfaos "$tipo" "$dir/repo" | sed "s/^/${tipo}: nenhuma linha do mapa alcança /" >> "$problemas"
  done
  if [ -s "$problemas" ]; then
    echo "spec-check auditar: $(wc -l < "$problemas" | tr -d ' ') problema(s) em $CONFIG"
    sed 's/^/  - /' "$problemas"
    return 1
  fi
  echo "spec-check auditar: ok ($CONFIG confere com o repo)"
  return 0
}

resolver() { # <tipo> → define cmd_final e alcance
  cmd_final=
  inteiro=$(cfg "$1")
  impacto=$(cfg "$1_impacto")
  if [ -z "$inteiro" ]; then
    alcance=n/a
    return
  fi
  if [ "$modo" = versao ]; then
    cmd_final=$inteiro
    alcance=inteira
    return
  fi
  if [ -z "$impacto" ]; then
    cmd_final=$inteiro
    alcance="inteira: sem por impacto"
    return
  fi
  pendente="$base/pendente-$1"
  if [ "$dados" = 1 ] && [ -f "$pendente" ] && grep -qx INTEIRA "$pendente"; then
    cmd_final=$inteiro
    alcance="inteira: falhou na tentativa anterior"
    return
  fi
  if grep -qxF "$CONFIG" "$dir/arquivos"; then
    cmd_final=$inteiro
    alcance="inteira: $CONFIG mudou"
    return
  fi
  motivo=$(cfg "$1_inteira" | while IFS= read -r p; do
    achou=$(grep -Em1 "^$(regex_de "$p")\$" "$dir/arquivos") && {
      printf '%s\n' "$achou"
      break
    }
  done)
  if [ -n "$motivo" ]; then
    cmd_final=$inteiro
    alcance="inteira: $motivo"
    return
  fi
  alvos=$({
    alvos_do_mapa "$1"
    comando_alvos=$(cfg "$1_alvos")
    if [ -n "$comando_alvos" ] && [ -s "$dir/arquivos" ]; then
      sh -c "$comando_alvos" < "$dir/arquivos"
    fi
    if [ "$dados" = 1 ] && [ -f "$pendente" ]; then
      cat "$pendente"
    fi
  } | tr -d '\r' | sort -u | tr '\n' ' ' | sed 's/ *$//')
  if [ -z "$alvos" ]; then
    alcance="nenhum alvo"
    return
  fi
  cmd_final=${impacto%%\{alvos\}*}$alvos${impacto#*\{alvos\}}
  alcance="por impacto: $alvos"
}

por_impacto() { # <tipo>
  tipo=$1
  resolver "$tipo"
  pendente="$base/pendente-$tipo"
  if [ -z "$cmd_final" ]; then
    case $alcance in
      n/a) resultado=n/a ;;
      *)
        resultado=nenhum_alvo
        rm -f "$pendente"
        ;;
    esac
    registrar "$tipo" "$resultado" 0 "$alcance" "" ""
    return
  fi
  if etapa "$tipo" "$cmd_final" "$alcance"; then
    rm -f "$pendente"
    return
  fi
  [ "$modo" = versao ] && return
  case $alcance in
    inteira*) echo INTEIRA > "$pendente" ;;
    *) printf '%s\n' "$alvos" | tr ' ' '\n' > "$pendente" ;;
  esac
}

fechar() {
  case $geral in
    verde) codigo=0 ;;
    vermelho) codigo=1 ;;
    *) codigo=2 ;;
  esac
  awk -F'\t' -v id="$id" -v modo="$modo" -v geral="$geral" -v arvore="$arvore" -v dir="$dir" '
    function j(s,  r, p, c) {
      r = ""
      for (p = 1; p <= length(s); p++) {
        c = substr(s, p, 1)
        if (c == "\\" || c == "\"") r = r "\\"
        if (c != "\r" && c != "\t") r = r c
      }
      return "\"" r "\""
    }
    BEGIN {
      printf "{\n  \"id\": %s,\n  \"modo\": %s,\n  \"status\": %s,\n  \"arvore\": %s,\n  \"etapas\": [", j(id), j(modo), j(geral), j(arvore)
    }
    {
      arq = dir "/" $1 ".pulos"
      pulos = ""
      while ((getline l < arq) > 0) pulos = pulos (pulos == "" ? "" : ", ") j(l)
      close(arq)
      printf "%s\n    {\"nome\": %s, \"resultado\": %s, \"duracao_s\": %d, \"alcance\": %s, \"log\": %s, \"comando\": %s, \"pulos\": [%s]}", (NR > 1 ? "," : ""), j($1), j($2), $3, j($4), j($5), j($6), pulos
    }
    END { print "\n  ]\n}" }
  ' "$dir/etapas" > "$dir/resultado.json"
  cp "$dir/resultado.json" "$base/ultimo.json"
  {
    echo "spec-check $modo: $geral"
    awk -F'\t' '{ printf "  %-12s %-12s %5ss  %s\n", $1, $2, $3, $4 }' "$dir/etapas"
    awk -F'\t' '$2 == "falhou" || $2 == "pulou" { print $1 "\t" $2 "\t" $5 }' "$dir/etapas" \
      | while IFS="$(printf '\t')" read -r nome resultado log; do
        echo
        if [ "$resultado" = pulou ]; then
          echo "--- $nome pulou testes ($log):"
          head -n 40 "$dir/$nome.pulos"
        else
          echo "--- $nome falhou, fim da saída ($log):"
          tail -n 40 "$log"
        fi
      done
    echo
    echo "Resultado: $base/ultimo.json"
  } > "$dir/resumo"
  printf '%s\n' "$codigo" > "$dir/fim.tmp"
  mv "$dir/fim.tmp" "$dir/fim"
}

executar() { # <id>
  id=$1
  dir="$base/execucoes/$id"
  modo=$(cat "$dir/modo")
  dados=$(cat "$dir/dados")
  geral=verde
  usa_cache=1
  arvore=$(arvore_atual)
  [ -n "$arvore" ] || usa_cache=0
  [ "$modo" = versao ] && usa_cache=0
  : > "$dir/etapas"

  rodar_fixa() { # <nome>
    c=$(cfg "$1")
    if [ -n "$c" ]; then
      etapa "$1" "$c" inteira
    else
      registrar "$1" n/a 0 n/a "" ""
    fi
  }

  if [ "$modo" = passo ]; then
    n=0
    while IFS= read -r c; do
      n=$((n + 1))
      etapa "teste-$n" "$c" "do passo"
    done < "$dir/comandos"
    rodar_fixa build
    rodar_fixa lint
  else
    if [ "$modo" = mudanca ]; then
      if [ -n "$(cfg tela_testes)$(cfg integracao_testes)" ]; then
        lista_repo > "$dir/repo"
        {
          orfaos integracao "$dir/arquivos" | sed 's/^/integracao: nenhuma linha do mapa alcança /'
          orfaos tela "$dir/arquivos" | sed 's/^/tela: nenhuma linha do mapa alcança /'
        } > "$dir/mapa.log"
        if [ -s "$dir/mapa.log" ]; then
          geral=vermelho
          registrar mapa falhou 0 "teste do diff fora do mapa" "$dir/mapa.log" ""
        else
          registrar mapa ok 0 "testes do diff no mapa" "" ""
        fi
      else
        registrar mapa n/a 0 "sem *_testes" "" ""
      fi
    fi
    rodar_fixa lint
    rodar_fixa build
    rodar_fixa suite
    if [ "$geral" = verde ]; then
      por_impacto integracao
      por_impacto tela
    else
      registrar integracao nao_rodou 0 "o básico não ficou verde" "" ""
      registrar tela nao_rodou 0 "o básico não ficou verde" "" ""
    fi
  fi
  fechar
}

vivo() {
  [ -f "$1/pid" ] && kill -0 "$(cat "$1/pid")" 2>/dev/null
}

aguardar() { # <id>
  alvo="$base/execucoes/$1"
  t=0
  while [ ! -f "$alvo/fim" ]; do
    if [ "$t" -ge "$ESPERA" ]; then
      echo "spec-check: ainda rodando ($1). Sem laço de acompanhamento: rode 'sh $self esperar'."
      exit 3
    fi
    if ! vivo "$alvo"; then
      sleep 1
      [ -f "$alvo/fim" ] && break
      echo "spec-check: a execução $1 morreu antes de terminar. Veja $alvo" >&2
      exit 1
    fi
    sleep 2
    t=$((t + 2))
  done
  cat "$alvo/resumo"
  exit "$(cat "$alvo/fim")"
}

comando=${1:-}
[ $# -gt 0 ] && shift

case $comando in
  _worker)
    executar "$1"
    exit 0
    ;;
  esperar)
    [ -f "$base/atual" ] || {
      echo "spec-check: nenhuma execução para esperar" >&2
      exit 4
    }
    aguardar "$(cat "$base/atual")"
    ;;
  passo)
    [ $# -ge 1 ] || uso
    ;;
  mudanca)
    if [ $# -gt 0 ]; then
      [ "$1" = --arquivos ] && [ $# -ge 2 ] || uso
    fi
    ;;
  versao | auditar)
    [ $# -eq 0 ] || uso
    ;;
  *)
    uso
    ;;
esac

if [ "$comando" = auditar ]; then
  dir=$(mktemp -d)
  trap 'rm -rf "$dir"' EXIT
  [ -f "$CONFIG" ] || {
    echo "spec-check: falta $CONFIG. Gere-o como manda a skill spec-check." >&2
    exit 4
  }
  ler_config > "$dir/config"
  validar_config
  auditar
  exit $?
fi

if [ -f "$base/atual" ]; then
  anterior=$(cat "$base/atual")
  if [ ! -f "$base/execucoes/$anterior/fim" ] && vivo "$base/execucoes/$anterior"; then
    echo "spec-check: já existe uma execução em andamento ($anterior). Rode 'sh $self esperar'."
    exit 3
  fi
fi

id=$(date +%Y%m%d-%H%M%S)-$$
dir="$base/execucoes/$id"
mkdir -p "$dir"
printf '%s\n' "$comando" > "$dir/modo"
echo 0 > "$dir/dados"

if [ -f "$CONFIG" ]; then
  ler_config > "$dir/config"
  validar_config
elif [ "$comando" = passo ]; then
  : > "$dir/config"
else
  echo "spec-check: falta $CONFIG. Gere-o como manda a skill spec-check." >&2
  rm -rf "$dir"
  exit 4
fi

: > "$dir/comandos"
if [ "$comando" = passo ]; then
  for c in "$@"; do
    printf '%s\n' "$c" >> "$dir/comandos"
  done
fi

if [ "$comando" = mudanca ] && [ $# -gt 0 ]; then
  shift
  echo 1 > "$dir/dados"
  for a in "$@"; do
    printf '%s\n' "$a"
  done
else
  git -c core.quotepath=off diff --name-only --no-renames HEAD 2>/dev/null
  git -c core.quotepath=off ls-files -o --exclude-standard
fi | tr -d '\r' | grep -Ev '^docs/|\.md$' | sort -u > "$dir/arquivos"

printf '%s\n' "$id" > "$base/atual"
nohup sh "$self" _worker "$id" > "$dir/worker.log" 2>&1 < /dev/null &
echo $! > "$dir/pid"
aguardar "$id"
