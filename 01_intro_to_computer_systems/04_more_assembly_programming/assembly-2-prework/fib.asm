section .text
global fib
fib:
  cmp rdi, 1 ; compare n to 1, if n<=1, goto done
  mov rax, rdi ; init return value to n
  jle done ; if we reached 1, return

  push rbx ; save rbx to stack
  mov rbx, rdi ;  save function argument to callee-saved register
  dec rdi  ; set n = n - 1
  call fib

  lea rdi, [-2+rbx] ; set n = n - 2
  push rax ; save fib(n-1) to stack
  call fib
  pop r10 ; restore fib(n-1) from stack
  
  add rax, r10
  pop rbx
done: 
  ret