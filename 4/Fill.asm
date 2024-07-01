// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/4/Fill.asm

// Runs an infinite loop that listens to the keyboard input. 
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel. When no key is pressed, 
// the screen should be cleared.

(LOOP)
@SCREEN
D=A
@address
M=D

@KBD
D=M

@CLEARING
D; JEQ
@SHADING
D; JGT

(SHADING) // for (i=0;i<8192;++i) shade RAM[M+i]
    @i
    M=0
    (SHADE)
        @8192
        D=A
        @i
        D=D-M
        @LOOP
        D; JEQ

        @i
        D=M
        @address
        A=M+D
        M=-1

        @i
        M=M+1

        @SHADE
        0; JMP


(CLEARING) // for (i=0;i<8192;++i) clear RAM[M+i]
    @i
    M=0
    (CLEAR)
        @8192
        D=A
        @i
        D=D-M
        @LOOP
        D; JEQ

        @i
        D=M
        @address
        A=M+D
        M=0

        @i
        M=M+1
        @CLEAR
        0; JMP

@LOOP
0;JMP