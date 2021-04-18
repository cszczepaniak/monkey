# monkey

This is my repo for working through [Writing an Interpreter in Go](https://interpreterbook.com/) and subsequently [Writing a Compiler in Go](https://compilerbook.com/), both by Thorsten Ball.

## Why are you learning about interpreters and compilers?

My curiosity began because I used to work professionally quite a bit with [MATLAB](https://www.mathworks.com/products/matlab.html), but sadly MATLAB lacks modern tooling that I became accustomed to - namely, I wanted MATLAB to have a formatter, much like [prettier](https://prettier.io/) for JavaScript/TypeScript or [gofmt](https://golang.org/cmd/gofmt/) for Go. I set out to see what it would take to write one myself, and quickly discovered you need to parse the language to be able to format it... but I didn't know anything about that!

It turns out MATLAB is a major pain to parse because of some of its peculiar syntax, so [my formatter](https://github.com/cszczepaniak/mfmt) never took off, and I since have stopped working with MATLAB. But I was left with my curiosity about lexers, parsers, ASTs, and all the rest of the moving parts that are required to compile a language. I watched most of a [livestreamed series](https://www.youtube.com/watch?v=wgHIkdUQbp0&list=PLRAdsfhKI4OWNOSfS7EUu5GRAVmze1t2y) on YouTube from Immo Landwerth where he builds his own programming language, and he recommended these two books as a great resource for learning - so here we are!
