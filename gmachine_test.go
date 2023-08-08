package gmachine_test

import (
	"errors"
	"testing"

	"gmachine"

	"github.com/google/go-cmp/cmp"
)

// TODO(adam): Add CLI tests using testscript

func TestNew(t *testing.T) {
	t.Parallel()
	g := gmachine.New()
	var wantP gmachine.Word = 0
	if wantP != g.P {
		t.Errorf("want initial P value %d, got %d", wantP, g.P)
	}
	var wantA gmachine.Word = 0
	if wantA != g.A {
		t.Errorf("want initial A value %d, got %d", wantA, g.A)
	}
	var wantMemValue gmachine.Word = 0
	gotMemValue := g.Memory[gmachine.MemSize-1]
	if wantMemValue != gotMemValue {
		t.Errorf("want last memory location to contain %d, got %d", wantMemValue, gotMemValue)
	}
}

func TestHALT(t *testing.T) {
	t.Parallel()
	g := gmachine.New()
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
	g := gmachine.New()
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
	g := gmachine.New()
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
	g := gmachine.New()
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
	g := gmachine.New()
	g.P = gmachine.MemSize - 1
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
	g := gmachine.New()
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
	g := gmachine.New()
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

func TestAssembleInvalidSourceCode(t *testing.T) {
	t.Parallel()
	_, err := gmachine.Assemble("ILLEGAL")
	wantErr := gmachine.ErrUndefinedInstruction
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}

func TestAssembleAndRun(t *testing.T) {
	t.Parallel()
	g := gmachine.New()
	err := g.AssembleAndRun("INCA\nHALT")
	if err != nil {
		t.Fatal("didn't expect an error", err)
	}
	var wantA gmachine.Word = 1
	if g.A != wantA {
		t.Errorf("want A value %d, got %d", wantA, g.A)
	}
}

func TestSETAWithInvalidNumber(t *testing.T) {
	t.Parallel()
	g := gmachine.New()
	err := g.AssembleAndRun("SETA a")
	wantErr := gmachine.ErrInvalidNumber
	if err == nil {
		t.Fatal("expected an error to be returned for invalid argument to SETA")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("wanted error %v, got %v", wantErr, err)
	}
}
