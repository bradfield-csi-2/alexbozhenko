default rel

section .text
global volume
volume:
	; xmm0 - r 
	; xmm1 - h
	mulss xmm0, xmm0
	mulss xmm0, xmm1
	mulss xmm0, [pi]
	mulss xmm0, [one_third]
 	ret
section .rodata
pi:  dd 3.141592654
one_third: dd 0.333333
