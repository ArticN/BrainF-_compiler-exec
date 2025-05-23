// bfc.go
package main

import (
    "fmt"
    "io"
    "os"
    "strings"
    "strconv"
    "unicode"
)

func main() {
    data, err := io.ReadAll(os.Stdin)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    text := strings.TrimSpace(string(data))
    parts := strings.SplitN(text, "=", 2)
    if len(parts) != 2 {
        fmt.Fprintln(os.Stderr, "uso: VAR=EXPR")
        os.Exit(1)
    }
    varName := parts[0]
    exprStr := parts[1]

    // Parse
    p := &Parser{s: exprStr}
    ast := p.parseExpr()

    // Gen BF
    gen := &BFGen{}

    // 1) imprime VAR=
    for _, c := range varName + "=" {
        for _, b := range []byte(string(c)) {
            gen.moveTo(10)
            gen.zero()
            gen.inc(int(b))
            gen.sb.WriteByte('.')
        }
    }

    // 2) computa a expressão em cell 0
    ast.Gen(gen, 0)

    // 3) converte para decimal e imprime
    gen.emitPrintDecimal()

    // 4) saída
    fmt.Print(gen.String())
}

// ——————————————————————————————————————————
// Parser simples (descida recursiva) + AST
// ——————————————————————————————————————————

type Parser struct {
    s   string
    pos int
}

func (p *Parser) peek() rune {
    if p.pos >= len(p.s) {
        return 0
    }
    return rune(p.s[p.pos])
}

func (p *Parser) consume() rune {
    ch := p.peek()
    if ch != 0 {
        p.pos++
    }
    return ch
}

func (p *Parser) parseExpr() Node {
    node := p.parseTerm()
    for {
        switch p.peek() {
        case '+', '-':
            op := byte(p.consume())
            right := p.parseTerm()
            node = &BinOp{op: op, left: node, right: right}
        default:
            return node
        }
    }
}

func (p *Parser) parseTerm() Node {
    node := p.parseFactor()
    for {
        if p.peek() == '*' {
            p.consume()
            right := p.parseFactor()
            node = &BinOp{op: '*', left: node, right: right}
        } else {
            return node
        }
    }
}

func (p *Parser) parseFactor() Node {
    if p.peek() == '(' {
        p.consume()
        node := p.parseExpr()
        if p.peek() == ')' {
            p.consume()
        }
        return node
    }
    start := p.pos
    for unicode.IsDigit(p.peek()) {
        p.consume()
    }
    num, err := strconv.Atoi(p.s[start:p.pos])
    if err != nil {
        fmt.Fprintf(os.Stderr, "número inválido: %s\n", p.s[start:p.pos])
        os.Exit(1)
    }
    return &Number{val: num}
}

type Node interface {
    Gen(g *BFGen, cell int)
}

type Number struct{ val int }

func (n *Number) Gen(g *BFGen, cell int) {
    g.moveTo(cell)
    g.zero()
    g.inc(n.val)
}

type BinOp struct {
    op         byte
    left, right Node
}

func (b *BinOp) Gen(g *BFGen, cell int) {
    switch b.op {
    case '+':
        b.left.Gen(g, cell)
        b.right.Gen(g, cell+1)
        g.emitAdd(cell+1, cell)
    case '-':
        b.left.Gen(g, cell)
        b.right.Gen(g, cell+1)
        g.emitSub(cell+1, cell)
    case '*':
        b.left.Gen(g, cell)
        b.right.Gen(g, cell+1)
        // usamos cell+2 como res e cell+3 como backup
        g.emitMul(cell, cell+1, cell+2, cell+3)
    }
}

// ——————————————————————————————————————————
// Gerador de Brainfuck
// ——————————————————————————————————————————

type BFGen struct {
    sb  strings.Builder
    pos int // posição atual (célula) do ponteiro
}

// move o ponteiro até a célula `c`
func (g *BFGen) moveTo(c int) {
    for g.pos < c {
        g.sb.WriteByte('>')
        g.pos++
    }
    for g.pos > c {
        g.sb.WriteByte('<')
        g.pos--
    }
}

// zera a célula atual
func (g *BFGen) zero() {
    g.sb.WriteString("[-]")
}

// incrementa a célula atual `n` vezes
func (g *BFGen) inc(n int) {
    for i := 0; i < n; i++ {
        g.sb.WriteByte('+')
    }
}

// loop em torno da célula `c`: [ body ]
func (g *BFGen) emitLoop(c int, body func()) {
    g.moveTo(c)
    g.sb.WriteByte('[')
    body()
    g.moveTo(c)
    g.sb.WriteByte(']')
}

// soma SRC → DST: enquanto SRC>0, --SRC, ++DST
func (g *BFGen) emitAdd(src, dst int) {
    g.emitLoop(src, func() {
        g.sb.WriteByte('-')
        g.moveTo(dst)
        g.sb.WriteByte('+')
        g.moveTo(src)
    })
    g.moveTo(dst)
}

// subtração SRC de DST: enquanto SRC>0, --SRC, --DST
func (g *BFGen) emitSub(src, dst int) {
    g.emitLoop(src, func() {
        g.sb.WriteByte('-')
        g.moveTo(dst)
        g.sb.WriteByte('-')
        g.moveTo(src)
    })
    g.moveTo(dst)
}

// multiplicação A * B em RES, usando TMP como backup, e coloca em A
func (g *BFGen) emitMul(a, b, res, tmp int) {
    // zera RES e TMP
    g.moveTo(res); g.zero()
    g.moveTo(tmp); g.zero()

    // enquanto A>0
    g.emitLoop(a, func() {
        // decrementa A
        g.moveTo(a); g.sb.WriteByte('-')

        // copia B → RES e TMP
        g.emitLoop(b, func() {
            g.sb.WriteByte('-')
            g.moveTo(res); g.sb.WriteByte('+')
            g.moveTo(tmp); g.sb.WriteByte('+')
            g.moveTo(b)
        })

        // restaura B de TMP
        g.emitLoop(tmp, func() {
            g.sb.WriteByte('-')
            g.moveTo(b); g.sb.WriteByte('+')
            g.moveTo(tmp)
        })
    })

    // move RES de volta para A
    g.emitLoop(res, func() {
        g.sb.WriteByte('-')
        g.moveTo(a); g.sb.WriteByte('+')
        g.moveTo(res)
    })
    g.moveTo(a)
}

// converte a valor em cell 0 de base10 e imprime dígitos
// usa cell1 como quociente
func (g *BFGen) emitPrintDecimal() {
    // limpa quociente
    g.moveTo(1); g.zero()

    // dividir por 10: enquanto cell0 ≥ 10, subtrai 10 e ++cell1
    g.moveTo(0)
    g.sb.WriteByte('[')
    for i := 0; i < 10; i++ { g.sb.WriteByte('-') }
    g.moveTo(1); g.sb.WriteByte('+')
    g.moveTo(0)
    g.sb.WriteByte(']')

    // se quociente>0, imprime dígito (Q + '0')
    g.moveTo(1)
    g.sb.WriteByte('[')
    for i := 0; i < 48; i++ { g.sb.WriteByte('+') }
    g.sb.WriteByte('.')
    g.zero()
    g.sb.WriteByte(']')

    // imprime resto (cell0 + '0')
    g.moveTo(0)
    for i := 0; i < 48; i++ { g.sb.WriteByte('+') }
    g.sb.WriteByte('.')
    g.zero()
}

// String retorna todo o Brainfuck gerado
func (g *BFGen) String() string {
    return g.sb.String()
}
