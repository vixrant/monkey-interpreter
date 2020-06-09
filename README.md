# 🐒 Monkey Language interpreter
##### A programming language from the book "Write an interpreter in Go"

Still WIP.

## Build

`go build`

## Use

`./mkc` will start the REPL.

`./mkc <filename>` will interpret a file.

## Features
- let statements
- expression evaluation
- first class functions
- conditions construct: if
- operators: + - / * **
- dynamic type system
- functions and closures (first class functions)
    ```rust
    let factorial = fn(x) {
        if(x <= 1) {
            1;
        } else {
            factorial(x-1) * factorial(x-2);
        };
    };
    ```
- variable scoping

## TODO other than book
- [ ] if-else-if ladder
- [ ] loop constructs
- [ ] llvm code generation
- [ ] fix the return statement

## Reference

[1] Ball, Thorsten. _Writing An Interpreter In Go_. 2016. Web.
