package gmachine_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gmachine/parser"
	"io"
	"os"
	"strings"
	"testing"

	"gmachine"

	"github.com/google/go-cmp/cmp"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"gc": gmachine.MainCompile,
		"gr": gmachine.MainRun,
	}))
}

func Test(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
	})
}

func TestNew(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantP gmachine.Word = 0
	if wantP != g.P {
		t.Errorf("want initial P value %d, got %d", wantP, g.P)
	}
	var wantS gmachine.Word = 0
	if wantS != g.S {
		t.Errorf("want initial S value %d, got %d", wantS, g.S)
	}
	var wantA gmachine.Word = 0
	if wantA != g.A {
		t.Errorf("want initial A value %d, got %d", wantA, g.A)
	}
	var wantX gmachine.Word = 0
	if wantX != g.X {
		t.Errorf("want initial X value %d, got %d", wantX, g.X)
	}
	var wantY gmachine.Word = 0
	if wantY != g.Y {
		t.Errorf("want initial Y value %d, got %d", wantY, g.Y)
	}
	var wantMemValue gmachine.Word = 0
	gotMemValue := g.Memory[gmachine.MemSize-1]
	if wantMemValue != gotMemValue {
		t.Errorf("want last memory location to contain %d, got %d", wantMemValue, gotMemValue)
	}
	var wantMemOffset gmachine.Word = gmachine.StackSize
	gotMemOffset := g.MemOffset
	if wantMemOffset != gotMemOffset {
		t.Errorf("want memory offset %d, got %d", wantMemValue, gotMemValue)
	}
}

func TestHALT(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantP gmachine.Word = 1
	err := assembleAndRunFromString(g, "HALT")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	if wantP != g.P {
		t.Errorf("want P == %d, got P == %d", wantP, g.P)
	}
}

func TestNOOP(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantP gmachine.Word = 2
	err := assembleAndRunFromString(g, "NOOP")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	if wantP != g.P {
		t.Errorf("want P == %d, got P == %d", wantP, g.P)
	}
}

func TestINCA(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, "INCA")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	var wantA gmachine.Word = 1
	if wantA != g.A {
		t.Errorf("want initial A value %d, got %d", wantA, g.A)
	}
}

func TestIllegalInstruction(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	g.RunProgram([]gmachine.Word{
		0,
	})
	var wantE = gmachine.ExceptionIllegalInstruction
	if wantE != g.E {
		t.Errorf("want error code value %d, got %d", wantE, g.E)
	}
}

func TestOutOfMemoryException(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	g.P = gmachine.MemSize - gmachine.StackSize - 1
	err := assembleAndRunFromString(g, "NOOP")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	var wantE = gmachine.ExceptionOutOfMemory
	if wantE != g.E {
		t.Errorf("want error code value %d, got %d", wantE, g.E)
	}
}

func TestDECA(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	g.A = 1
	var wantA gmachine.Word = 0
	err := assembleAndRunFromString(g, "DECA")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	if wantA != g.A {
		t.Errorf("want initial A value %d, got %d", wantA, g.A)
	}
}

func TestSETA(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantA gmachine.Word = 5
	err := assembleAndRunFromString(g, "SETA 5")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	if wantA != g.A {
		t.Errorf("want initial A value %d, got %d", wantA, g.A)
	}
}

func TestSETX(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantX gmachine.Word = 5
	err := assembleAndRunFromString(g, "SETX 5")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	if wantX != g.X {
		t.Errorf("want X value %d, got %d", wantX, g.X)
	}
}

func TestSETY(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantY gmachine.Word = 5
	err := assembleAndRunFromString(g, "SETY 5")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	if wantY != g.Y {
		t.Errorf("want Y value %d, got %d", wantY, g.Y)
	}
}

func TestAssemble(t *testing.T) {
	t.Parallel()
	want := []gmachine.Word{gmachine.OpINCA, gmachine.OpHALT}
	program, err := assembleFromString("INCA\nHALT")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	if !cmp.Equal(want, program) {
		t.Errorf(cmp.Diff(want, program))
	}
}

func TestAssemble_ReturnsErrorForUnknownInstruction(t *testing.T) {
	t.Skip("TODO: Determine how to handle invalid statements")
	t.Parallel()
	_, err := assembleFromString("ILLEGAL")
	wantErr := gmachine.ErrUnknownIdentifier
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}

func TestAssembleAndRun(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, "INCA\nHALT")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	var wantA gmachine.Word = 1
	if g.A != wantA {
		t.Errorf("want A value %d, got %d", wantA, g.A)
	}
}

func TestSETA_ReturnsErrorForInvalidNumber(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, "SETA 2a")
	wantErr := parser.ErrInvalidIntegerLiteral
	if err == nil {
		t.Fatal("expected an error to be returned for invalid argument to SETA")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}

func TestSETX_ReturnsErrorForInvalidNumber(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, "SETX 2a")
	wantErr := parser.ErrInvalidIntegerLiteral
	if err == nil {
		t.Fatal("expected an error to be returned for invalid argument to SETA")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}

func TestSETY_ReturnsErrorForInvalidNumber(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, "SETY 2a")
	wantErr := parser.ErrInvalidIntegerLiteral
	if err == nil {
		t.Fatal("expected an error to be returned for invalid argument to SETA")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}

func TestSETA_AcceptsCharacterLiteral(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, "SETA 'h'")
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	wantA := 'h'
	gotA := rune(g.A)
	if wantA != gotA {
		t.Errorf("want A %d, got %d", wantA, gotA)
	}
}

func TestSETX_AcceptsCharacterLiteral(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, "SETX 'h'")
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	wantX := 'h'
	gotX := rune(g.X)
	if wantX != gotX {
		t.Errorf("want X %d, got %d", wantX, gotX)
	}
}

func TestSETY_AcceptsCharacterLiteral(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, "SETY 'h'")
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	wantY := 'h'
	gotY := rune(g.Y)
	if wantY != gotY {
		t.Errorf("want Y %d, got %d", wantY, gotY)
	}
}

func TestOUTA_SerializesValueAsBytesInBigEndianOrder(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	out := io.Writer(&buf)
	g := gmachine.New(out)
	err := assembleAndRunFromString(g, "SETA 1\nOUTA")
	if err != nil {
		t.Fatalf("didn't expect an error: %v", err)
	}

	want := []byte{0, 0, 0, 0, 0, 0, 0, 1}
	got := buf.Bytes()
	if !cmp.Equal(want, got) {
		t.Errorf(cmp.Diff(want, got))
	}
}

func TestOUTA(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	out := io.Writer(&buf)
	g := gmachine.New(out)
	err := assembleAndRunFromString(g, `
SETA 'h'
OUTA
SETA 'e'
OUTA
SETA 'l'
OUTA
SETA 'l'
OUTA
SETA 'o'
OUTA
SETA ' '
OUTA
SETA 'w'
OUTA
SETA 'o'
OUTA
SETA 'r'
OUTA
SETA 'l'
OUTA
SETA 'd'
OUTA
SETA '!'
OUTA`)
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}

	output := []uint64{}
	for _, r := range "hello world!" {
		output = append(output, uint64(r))
	}

	wantBuf := bytes.Buffer{}
	binary.Write(&wantBuf, binary.BigEndian, output)

	want := wantBuf.Bytes()
	got := buf.Bytes()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestJUMP(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantA gmachine.Word = 42
	err := assembleAndRunFromString(g, `
JUMP 3
HALT
SETA 41
INCA
`)
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	if wantA != g.A {
		t.Errorf("want A %d, got %d", wantA, g.A)
	}
}

func TestJUMPWithInvalidNumber(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, "JUMP 2a")
	wantErr := parser.ErrInvalidIntegerLiteral
	if err == nil {
		t.Fatal("expected an error to be returned for invalid argument to JUMP")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}

func TestAssemble_SkipsComments(t *testing.T) {
	t.Parallel()
	want := []gmachine.Word{}
	got, err := assembleFromString("; this is a comment")
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestPSHA(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantA gmachine.Word = 42
	var wantS gmachine.Word = 1
	var want gmachine.Word = 42
	err := assembleAndRunFromString(g, "SETA 42\nPSHA")
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	if wantA != g.A {
		t.Errorf("wanted A %d, got %d", wantA, g.A)
	}
	if wantS != g.S {
		t.Errorf("wanted S %d, got %d", wantS, g.S)
	}
	if want != g.Memory[wantS-1] {
		t.Errorf("wanted stack value %d, got %d", want, g.Memory[wantS])
	}
}

func TestPOPA(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantA gmachine.Word = 42
	var wantS gmachine.Word = 0
	err := assembleAndRunFromString(g, `
SETA 42
PSHA
SETA 3
POPA
`)
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	if wantA != g.A {
		t.Errorf("wanted A %d, got %d", wantA, g.A)
	}
	if wantS != g.S {
		t.Errorf("wanted S %d, got %d", wantS, g.S)
	}
}

func TestMOVAX(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantX gmachine.Word = 42
	err := assembleAndRunFromString(g, "SETA 42\nMOVA X\n")
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	if wantX != g.X {
		t.Errorf("want %d, got %d", wantX, g.X)
	}
}

func TestMOVAY(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantY gmachine.Word = 42
	err := assembleAndRunFromString(g, "SETA 42\nMOVA Y\n")
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	if wantY != g.Y {
		t.Errorf("want %d, got %d", wantY, g.Y)
	}
}

func TestMOVA_FailsForInvalidRegister(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, "MOVA Z")
	wantErr := gmachine.ErrInvalidOperand
	if err == nil {
		t.Fatal("expected an error to be returned for invalid argument to MOVA")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}

func TestADDAX(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantA gmachine.Word = 10
	err := assembleAndRunFromString(g, `
SETA 6
MOVA X
SETA 4
ADDA X
`)
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	if wantA != g.A {
		t.Errorf("want A %d, got %d", wantA, g.A)
	}
}

func TestADDAY(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantA gmachine.Word = 10
	err := assembleAndRunFromString(g, `
SETA 6
MOVA Y
SETA 4
ADDA Y
`)
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	if wantA != g.A {
		t.Errorf("want A %d, got %d", wantA, g.A)
	}
}

func TestADDA_FailsForInvalidRegister(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, "ADDA Z")
	wantErr := gmachine.ErrInvalidOperand
	if err == nil {
		t.Fatal("expected an error to be returned for invalid argument to ADDA")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}

func TestAddTwoNumbers(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, `
; x = 4, y = 6
SETA 4
PSHA
SETA 6
PSHA
; add x y
POPA
MOVA X
POPA
ADDA X
`)
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	var wantA gmachine.Word = 10
	if wantA != g.A {
		t.Errorf("want A %d, got %d", wantA, g.A)
	}
}

func TestMULAX(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, `
SETA 5
MOVA X
SETA 2
MULA X`)
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	var wantA gmachine.Word = 10
	if wantA != g.A {
		t.Errorf("want A %d, got %d", wantA, g.A)
	}
}
func TestMULAY(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := assembleAndRunFromString(g, `
SETA 5
MOVA Y
SETA 2
MULA Y`)
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	var wantA gmachine.Word = 10
	if wantA != g.A {
		t.Errorf("want A %d, got %d", wantA, g.A)
	}
}

func TestSubroutineLabel(t *testing.T) {
	t.Parallel()
	want := []gmachine.Word{gmachine.OpSETA, gmachine.Word(42), gmachine.OpOUTA}
	got, err := assembleFromString(`
.test
SETA 42
OUTA
`)
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestAssemble_ReturnsErrorWhenGivenAnUndefinedIdentifer(t *testing.T) {
	t.Parallel()
	wantErr := gmachine.ErrUnknownIdentifier
	_, err := assembleFromString("JUMP foo")
	if err == nil {
		t.Fatal("expected an error")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}

func TestAssemble_ReturnsErrorWhenJumpIsPassedInvalidArgument(t *testing.T) {
	t.Parallel()
	wantErr := gmachine.ErrInvalidOperand
	_, err := assembleFromString("JUMP 'a'")
	if err == nil {
		t.Fatal("expected an error")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}

func TestSubRoutineLabelsAreReplacedWithMemoryAddress(t *testing.T) {
	t.Parallel()
	want := []gmachine.Word{
		// Jump to .start
		gmachine.OpJUMP,
		gmachine.Word(11),
		// .test1
		gmachine.OpSETA,
		gmachine.Word(42),
		gmachine.OpOUTA,
		gmachine.OpHALT,
		// .test2
		gmachine.OpSETA,
		gmachine.Word(41),
		gmachine.OpINCA,
		gmachine.OpOUTA,
		gmachine.OpHALT,
		// .start
		gmachine.OpJUMP,
		gmachine.Word(6),
	}
	got, err := assembleFromString(`
JUMP start

.testA
SETA 42
OUTA
HALT

.testB
SETA 41
INCA
OUTA
HALT

.start
JUMP testB
`)
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestSettingAConstant(t *testing.T) {
	t.Parallel()

	input := `
CONS c 1
SETA c
`
	g := gmachine.New(nil)
	assembleAndRunFromString(g, input)
	wantA := gmachine.Word(1)
	gotA := g.A
	if wantA != gotA {
		t.Errorf("wanted %v, got %v", wantA, gotA)
	}
}

func TestConstantReferencesAreReplacedWithValues(t *testing.T) {
	t.Parallel()

	want := []gmachine.Word{
		gmachine.OpSETA,
		gmachine.Word(42),
		gmachine.OpOUTA,
	}
	got, err := assembleFromString(`
CONS c 42
SETA c
OUTA
`)
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestCompile(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	out := io.Writer(&buf)

	input := `
	SETA 42
	OUTA
`
	r := strings.NewReader(input)
	err := gmachine.Compile(r, out)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte{
		0, 0, 0, 0, 0, 0, 0, byte(gmachine.OpSETA),
		0, 0, 0, 0, 0, 0, 0, byte(42),
		0, 0, 0, 0, 0, 0, 0, byte(gmachine.OpOUTA),
	}
	got := buf.Bytes()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestCompile_FailsForInvalidInput(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(`SETA 4a`)
	err := gmachine.Compile(input, io.Discard)
	if err == nil {
		t.Error("expected an error")
	}
}

type errorWriter struct{}

func (w *errorWriter) Write(data []byte) (int, error) {
	return 0, errors.New("failed to write data")
}

func TestCompile_FailsForWriteError(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(`SETA 42`)
	err := gmachine.Compile(input, &errorWriter{})
	if err == nil {
		t.Error("expected an error")
	}
}

func assembleFromString(input string) ([]gmachine.Word, error) {
	return gmachine.Assemble(strings.NewReader(input))
}

func assembleAndRunFromString(g *gmachine.Machine, input string) error {
	return g.AssembleAndRun(strings.NewReader(input))
}
