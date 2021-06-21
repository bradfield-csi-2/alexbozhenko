section .text
global pangram
pangram:
       ; rdi - pointer to char array
       ; rdx - holds char read from memory
       ; rax - will hold '1' in the places of every char ever seen 
       mov rax, 0 ;; rax will hold the result
process_next_char:
       movzx rdx, byte [rdi] ;; read single char(byte) to rdx
       add rdi, 1 ; move pointer to the next char

       cmp rdx, 0 ;; stop if we reached end of the string (NULL char)
       je done 
       or rdx, 0b00100000 ;; set 6th bit to make sure we have a lowercase letter
                                               ; see https://upload.wikimedia.org/wikipedia/commons/c/cf/USASCII_code_chart.png
       sub rdx, [a_ascii_code] ;; substract 'a' ASCII code from the char code to get the 0-based character index

       cmp rdx, 0                      ; make sure we get a char
       jl process_next_char

       cmp rdx, 25                     ; make sure we get a char
       jg process_next_char

	   bts rax, rdx 				;; set bit in rax at position rdx to 1. 
	   								;; 4 operations below do the same ;-)

;;       mov rbp, 1 				; 1 will be shiften (char_number) times
;;       mov rcx, rdx				; shl accepts only cl as the second argument, so move char number to cl
;;       shl rbp, cl              ; rbp now has a bit set corresponding to character nubmer
;;       or rax, rbp              ; set the bit corresponding to the current char in the resulting bit pattern 

	   jmp process_next_char


done:
	   cmp rax, 0b11111111111111111111111111 ;; check if all 26 characters were seen
	   mov rcx, 0							 ;; cmovnee needs a register 		 
	   cmovne rax, rcx					     ;; set return code(rax) to 0(false) if not all chars were seen
       ret									 ;; return 26 ones bit pattern as true, and 0 as false

       section .data
a_ascii_code:  dq 0x61

