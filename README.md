# Monkey
![example workflow](https://github.com/Revolyssup/monkey/actions/workflows/test.yml/badge.svg)

## Implementing monkey interpreter while learning from [Writing an Interpreter in GO](https://1lib.in/book/2927612/f54145) `
## KeyNotes-
1. Pratt top-down parsing is used by parser to populate the AST, as this is a simple and scalable strategy.
2. Walk-the-tree strategy is used at the time of evaluation. A tree-
walking interpreter that recursively evaluates an AST is probably the slowest of all approaches,
but easy to build, extend, reason about and as portable as the language it’s implemented in. We’re going to take the AST our parser
builds for us and interpret it “on the fly”, without any preprocessing or compilation step.

3. No garbage collector of its own does monkey posses. GO's garbage collector does the job good enough for us.;)
