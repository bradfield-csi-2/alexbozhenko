section .text
global index
index:
	; rdi: matrix
	; rsi: rows
	; rdx: cols
	; rcx: rindex
	; r8: cindex
	mov rax, rcx  ; rax <- rindex
	imul rax, rdx ; rax <- rindex*cols
	imul rax, 4   ; rax <- 4 bytes(int size)*rindex*cols
	imul r8, 4    ; r8 <- cindex * 4 bytes (int size)
	add rax, r8	  ; rax <- 4*rindex*cols + 4*cindex
	add rax, rdi  ; rax <- matrix + 4*rindex*cols + 4*cindex
	mov rax, [rax] ; move memory from that address

	ret
