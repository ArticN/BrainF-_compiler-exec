## Brainfuck Compiler (`bfc`) & Executor (`bfe`)

Este repositório implementa dois utilitários em Go para avaliar expressões e imprimir resultados usando Brainfuck em tempo de execução:

1. **`bfc`**: compila uma entrada no formato `VAR=EXPR`, gera um programa Brainfuck que realiza **todas** as operações de `EXPR` em tempo de execução.
2. **`bfe`**: interpreta o Brainfuck produzido pelo `bfc` e imprime o resultado final na saída padrão.

---

### Processo

* **Parser**: `bfc` faz parsing recursivo-descendente (LL(1)) de `EXPR` (`+`, `-`, `*`, parênteses, literais numéricos).
* **Codegen**: constrói código Brainfuck que executa soma, subtração, multiplicação (via loops), conversão para decimal e impressão de dígitos.
* **Execução**: `bfe` mapeia saltos (`[`/`]`), gerencia fita de 30.000 células e executa os comandos Brainfuck, produzindo o texto `VAR=valor`.

---

### Como compilar

No diretório raiz, supondo as pastas `bfc/` e `bfe/` contendo os `.go`:

```bash
# Compilar o compilador Brainfuck
go build -o bfc ./bfc

# Compilar o executor Brainfuck
go build -o bfe ./bfe
```

---

### Exemplo de teste

```bash
artic@aart1c:~/cod/Heredia$ go build -o bfc ./bfc
artic@aart1c:~/cod/Heredia$ go build -o bfe ./bfe
artic@aart1c:~/cod/Heredia$ echo 'CRÉDITO=2*5+10' | bfc/bfc
>>>>>>>>>>[-]+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.[-]++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.[-]+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.[-]+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.[-]++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.[-]+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.[-]++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.[-]+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.[-]+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.<<<<<<<<<<[-]++>[-]+++++>[-]>[-]<<<[->[->+>+<<]>>[-<<+>>]<<<]>>[-<<+>>]<<>[-]++++++++++[-<+>]<>[-]<[---------->+<]>[++++++++++++++++++++++++++++++++++++++++++++++++.[-]]<++++++++++++++++++++++++++++++++++++++++++++++++.[-]
artic@aart1c:~/cod/Heredia$ echo 'CRÉDITO=2*5+10' | bfc/bfc | bfe/bfe
CRÉDITO=20
```
