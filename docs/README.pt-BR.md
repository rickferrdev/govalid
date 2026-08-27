# govalid

[Página inicial](../README.md) · [English](README.en.md) · [Referência completa da API](API.md)

`govalid` oferece regras combináveis para validar campos selecionados de
structs Go sem utilizar tags. O escopo da v0.1.0 inclui booleanos, strings,
inteiros, floats, maps, slices e arrays.

> Este é um projeto pre-v1. Nomes públicos e comportamentos ainda podem mudar
> antes da estabilização da API.

## Instalação

```bash
go get github.com/rickferrdev/govalid@v0.1.0
```

Utilize o caminho do módulo sem versão para instalar a revisão atual de
desenvolvimento.

## Uso básico

```go
validator := govalid.New()

err := validator.Validate(
	user,
	govalid.Field("Name", govalid.StringRequired(), govalid.StringMinLength(3)),
	govalid.Field("Age", govalid.IntBetween(18, 130)),
	govalid.Field("Active", govalid.BoolTrue()),
	govalid.Field("Scores", govalid.CollectionEach(govalid.IntBetween(0, 100))),
)
```

Campos de structs aninhadas podem ser acessados por caminhos como
`"Profile.Email"`. Por padrão, todas as falhas são coletadas. Utilize
`New(WithStopOnFirstError())` para encerrar no primeiro problema.

```go
var validationErr *govalid.FieldIssueError

if errors.As(err, &validationErr) {
	for _, issue := range validationErr.RulesIssues {
		fmt.Printf("%s: %s\n", issue.Path, issue.Message)
	}
}
```

## Resumo das regras

### Boolean

`BoolTrue`, `BoolFalse` e `BoolEqual` validam estado e igualdade de booleanos.

### String

As regras de string cobrem:

- valor obrigatório e comprimento baseado em runas Unicode;
- conteúdo alfabético, numérico, alfanumérico e ASCII;
- texto em minúsculas, maiúsculas ou sem espaços externos;
- prefixos, sufixos, conteúdo e conjuntos permitidos;
- regex, email, URL absoluta e UUID.

Regras somente de formato aceitam string vazia. Combine-as com
`StringRequired` quando o campo for obrigatório.

### Integer

As regras de inteiros aceitam todos os tipos assinados e não assinados,
`uintptr` e tipos definidos pelo usuário. Elas cobrem:

- comparação, igualdade, intervalos inclusivos e conjuntos;
- valores positivos, negativos e zero;
- paridade, divisibilidade, primos, compostos, potências de dois e quadrados
  perfeitos;
- portas, porcentagens e status HTTP.

Comparações mistas entre signed e unsigned preservam todo o intervalo de
`uint64`.

### Float

As regras de float aceitam `float32`, `float64` e tipos derivados. Elas cobrem:

- comparação, igualdade, intervalos e tolerância;
- sinal e estados de zero;
- representação de 32 e 64 bits;
- `NaN`, infinito positivo/negativo e finitude.

### Map

As regras de map cobrem:

- nil, vazio e comprimento;
- chaves obrigatórias, proibidas, permitidas ou exatas;
- presença de valores e valores zero/nil;
- regras aplicadas a cada chave, cada valor ou uma chave selecionada;
- igualdade, subset e superset.

```go
govalid.Field(
	"Labels",
	govalid.MapNotEmpty(),
	govalid.MapKeys(govalid.StringLowercase()),
	govalid.MapValues(govalid.StringRequired()),
)
```

### Collection

As regras de collection aceitam slices e arrays. Elas cobrem:

- nil, vazio e comprimento;
- presença, unicidade, valores zero e nil;
- igualdade e diferença;
- regras aplicadas a todos, qualquer, nenhum ou um índice específico.

Arrays nunca são nil; somente slices podem satisfazer `CollectionNil`.

## Limitações da v0.1.0

- Ainda não existem regras dedicadas `Bytes*`; `[]byte` pode usar apenas regras
  de collection.
- Ainda não existem regras recursivas `Struct*`; caminhos aninhados funcionam.
- Ainda não existem regras específicas para tempo, duração, ponteiros ou
  presença universal.
- Maps não possuem ordem de iteração garantida; a primeira falha aninhada pode
  variar.

Para consultar todas as assinaturas e aliases legados, veja a
[referência completa da API](API.md), escrita em inglês.
