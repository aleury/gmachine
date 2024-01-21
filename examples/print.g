JUMP run
VARB msg "hello world"

.run
SETX msg 	    ; set X register to the address of msg
JUMP print

.print
MOVE *X -> A    ; move the value at address in X to A
OUTA		    ; print A
INCX		    ; increment address stored in X
JANZ print	    ; jump to label 'done' if A = 0
HALT

