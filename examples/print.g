@data msg "hello world"

.run
SETX msg 	; set X register to the address of msg
JUMP print

.done
HALT

.print
MOVX A		; move the value of address in X to A
JAEZ done	; jump to label 'done' if A = 0
OUTA		; print A
INCX		; increment address stored in X
JUMP print	; jump back to the start of .print to print the next character

