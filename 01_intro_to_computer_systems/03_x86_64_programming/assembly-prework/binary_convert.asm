	section .text
global binary_convert
binary_convert:
	; rdi - pointer to char array
	; rax - return value
	; rdx - holds char read from memory
	mov rax, 0 ;; rax will hold the result
process_next_char:
	movzx rdx, byte [rdi] ;; read single char(byte) to rdx
	cmp rdx, 0 ;; stop if we reached end of the string (NULL char)
	je done 
	sub rdx, [zero_ascii_code] ;; substract '0' ASCII code from the char to get the number
	mov rbp, rdx ;; save number to rdx
	or rbp, 1 ;; OR number with 1 to check if we get a valid 
	cmp rbp, 1  ;; If not valid char - finish processing the string
	jne done
	shl rax, 1 ;; multiply result that we accumulated so far by 2
	add rax, rdx ;; add number that we just read
	add rdi, 1 ; move pointer to the next char
	jmp process_next_char
done:	
	ret

	section .data
zero_ascii_code: 	dq 0x30
