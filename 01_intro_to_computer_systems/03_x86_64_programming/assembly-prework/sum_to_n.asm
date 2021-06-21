section .text
global sum_to_n
sum_to_n:
	; rdi - argument
	; rax - return value

	
;; For loop version:
; 	mov rax, 0
; loop: 
; 	cmp rdi, 0
; 	jle done
; 	add rax, rdi
; 	dec rdi
; 	jmp loop
; done:
; 	ret


;; while loop:
; 	mov rax, 0
; loop: 
; 	add rax, rdi
; 	dec rdi
; 	cmp rdi, 0
; 	jge loop
; 	ret

;; o(1)
	mov rax, rdi
	inc rdi 
	imul rax, rdi
	mov rcx, 2
	mov rdx, 0
	idiv rcx
	ret

