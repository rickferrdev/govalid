# govalid · Visão geral em português

[← Página inicial](../README.md) · [English](README.en.md) · [Referência da API](API.md)

> Validação explícita e combinável para campos selecionados de structs Go.

## Comece aqui

```bash
go get github.com/rickferrdev/govalid@latest
```

```go
validator := govalid.New()

err := validator.Validate(
	user,
	govalid.Field("Name", govalid.Required(), govalid.StringMinLength(3)),
	govalid.Field("Age", govalid.IntBetween(18, 130)),
	govalid.FieldIfPresent("Profile.Email", govalid.StringEmail()),
	govalid.Field("Scores", govalid.CollectionEach(govalid.IntBetween(0, 100))),
)
```

O modelo é propositalmente pequeno:

1. Crie ou reutilize um `Validator`.
2. Associe caminhos e regras com `Field` ou `FieldIfPresent`.
3. Inspecione os erros estruturados retornados.

## O que está incluído

| Família | Cobertura |
| --- | --- |
| Boolean | Estado e igualdade |
| String | Comprimento Unicode, conteúdo, caixa, regex, email, URL e UUID |
| Integer | Comparações signed/unsigned, conjuntos, aritmética e domínios |
| Float | Comparações, tolerância, representação e valores especiais |
| Bytes | Conteúdo binário, encodings, documentos e regras por byte |
| Map | Chaves, valores, comprimento, relações e regras aninhadas |
| Collection | Conteúdo, comprimento, unicidade, comparação e itens aninhados |
| Universal | Presença, nil e zero value |
| Condicional | Condições booleanas, contexto e valores opcionais |

Tipos definidos pelo usuário são aceitos quando possuem um tipo subjacente
suportado. Comparações entre inteiros signed e unsigned preservam todo o
intervalo de `uint64`.

## Padrões de composição

Valide todas as chaves e valores de um map:

```go
govalid.Field(
	"Labels",
	govalid.MapKeys(govalid.StringLowercase()),
	govalid.MapValues(govalid.StringRequired()),
)
```

Ignore regras aninhadas quando o valor estiver ausente:

```go
govalid.Field(
	"Scores",
	govalid.Optional(
		govalid.CollectionEach(govalid.IntBetween(0, 100)),
	),
)
```

Aplique regras com base na struct raiz:

```go
govalid.WhenContext(func(ctx govalid.RuleContext) bool {
	return ctx.RootAny().(Account).Enabled
}, govalid.StringRequired())
```

## Erros estruturados

Por padrão, todas as falhas são retornadas em `*FieldIssueError`:

```go
var validationErr *govalid.FieldIssueError

if errors.As(err, &validationErr) {
	for _, issue := range validationErr.RulesIssues {
		fmt.Printf("%s: %s\n", issue.Path, issue.Message)
	}
}
```

Cada `Issue` expõe `Path`, o `Value` rejeitado, `Rule` e `Message`. Use
`WithStopOnFirstError`, `WithPanicOnFirstError` ou `WithSilenceErrors` para
alterar a execução. `WithIssueHandler` observa falhas independentemente do
comportamento escolhido para o retorno.

## Extensibilidade

- Implemente regras personalizadas com `Rule` e `RuleContext`.
- Use `FieldIfPresent` quando um ponteiro ou interface aninhada puder ser nil.
- Implemente `FieldSource` para descobrir `FieldSpec` em outro formato.
- Reutilize um validator concorrentemente quando callbacks e sources forem
  concurrency-safe.

`FieldSource` é a fronteira de compatibilidade para futuras integrações com
tags, schemas ou código gerado. Nenhuma sintaxe de tags embutida faz parte do
escopo inicial da v1.

## Escopo da v1

A API estável inicial foca em regras explícitas para valores primitivos, maps,
collections, bytes, caminhos aninhados, composição condicional e erros
estruturados. Regras recursivas `Struct*`, famílias de tempo/duração e um parser
de tags embutido podem evoluir em releases posteriores.

Maps não possuem ordem de iteração garantida, portanto a primeira falha
aninhada pode variar. Regras de string somente de formato aceitam valor vazio;
combine-as com `Required` ou `StringRequired` quando a presença for obrigatória.

## Próximos links

- [Referência completa da API](API.md)
- [README principal e exemplo completo](../README.md)
- [Documentação do pacote](https://pkg.go.dev/github.com/rickferrdev/govalid)
