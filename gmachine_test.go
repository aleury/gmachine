package gmachine_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"gmachine"

	"github.com/google/go-cmp/cmp"
)

// TODO(adam): Add CLI tests using testscript

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
	err := g.AssembleAndRun("HALT")
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
	err := g.AssembleAndRun("NOOP")
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
	err := g.AssembleAndRun("INCA")
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
	err := g.AssembleAndRun("NOOP")
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
	err := g.AssembleAndRun("DECA")
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
	err := g.AssembleAndRun("SETA 5")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	if wantA != g.A {
		t.Errorf("want initial A value %d, got %d", wantA, g.A)
	}
}

func TestAssemble(t *testing.T) {
	t.Parallel()
	want := []gmachine.Word{gmachine.OpINCA, gmachine.OpHALT}
	program, err := gmachine.Assemble("INCA\nHALT")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	if !cmp.Equal(want, program) {
		t.Errorf(cmp.Diff(want, program))
	}
}

func TestAssemble_ReturnsErrorForUnknownInstruction(t *testing.T) {
	t.Parallel()
	_, err := gmachine.Assemble("ILLEGAL")
	wantErr := gmachine.ErrUnknownIdentifier
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}

func TestAssembleAndRun(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := g.AssembleAndRun("INCA\nHALT")
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
	err := g.AssembleAndRun("SETA 2a")
	wantErr := gmachine.ErrUnknownIdentifier
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
	err := g.AssembleAndRun("SETA 'h'")
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	wantA := 'h'
	gotA := rune(g.A)
	if wantA != gotA {
		t.Errorf("want A %d, got %d", wantA, gotA)
	}
}

func TestOUTA(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	out := io.Writer(&buf)
	g := gmachine.New(out)
	err := g.AssembleAndRun(`
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
	want := "hello world!"
	got := buf.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestJUMP(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantA gmachine.Word = 42
	err := g.AssembleAndRun(`
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
	err := g.AssembleAndRun("JUMP 2a")
	wantErr := gmachine.ErrUnknownIdentifier
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
	got, err := gmachine.Assemble("; this is a comment")
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
	err := g.AssembleAndRun("SETA 42\nPSHA")
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
	err := g.AssembleAndRun(`
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

func TestMOVA(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := g.AssembleAndRun("SETA 42\nMOVA X\n")
	if err != nil {
		t.Fatal("didn't expect an error:", err)
	}
	var wantX gmachine.Word = 42
	if wantX != g.X {
		t.Errorf("want X %d, got %d", wantX, g.X)
	}
}

func TestMOVA_FailsForInvalidRegister(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := g.AssembleAndRun("MOVA Z")
	wantErr := gmachine.ErrInvalidRegister
	if err == nil {
		t.Fatal("expected an error to be returned for invalid argument to MOVA")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}

func TestADDA(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	var wantA gmachine.Word = 10
	err := g.AssembleAndRun(`
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

func TestADDA_FailsForInvalidRegister(t *testing.T) {
	t.Parallel()
	g := gmachine.New(nil)
	err := g.AssembleAndRun("ADDA Z")
	wantErr := gmachine.ErrInvalidRegister
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
	err := g.AssembleAndRun(`
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

func TestSubroutineLabel(t *testing.T) {
	t.Parallel()
	want := []gmachine.Word{gmachine.OpSETA, gmachine.Word(42), gmachine.OpOUTA}
	got, err := gmachine.Assemble(`
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
	_, err := gmachine.Assemble("JUMP foo")
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
	_, err := gmachine.Assemble("JUMP 'a'")
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
	got, err := gmachine.Assemble(`
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
