; Tasks:
; 1. Implement MULA opcode
; 2. Implement SETX opcode
; 3. Implement DECX opcode (and DECA, DECY)
; 4. Implement CMPX opcode and C status register to store result of CMPX.
; 5. Implement BINZ opcode
; 6. Implement CALL and RTRN opcodes

.factorial
MULA X
DECX
CMPX 0
BINZ factorial ; branch to .factorial if result of CMPX is Zero (0)
CLRC ; clear the compare result
RTRN ; exit subroutine

.start
SETA 1
SETX 6
CALL factorial ; push PC on stack and call subroutine.
OUTA
HALT

