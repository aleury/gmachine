; Tasks:
; 1. Implement MULA
; 2. Implement SETX
; 3. Implement DECX (and DECA, DECY)
; 4. Implement JXNZ
; 5. Implement CALL and RTRN

.factorial
MULA X
DECX
JXNZ factorial ; jump to .factorial if register X is not zero
RTRN ; exit subroutine

.start
SETA 1
SETX 6
CALL factorial ; push PC on stack and call subroutine.
OUTA
HALT

