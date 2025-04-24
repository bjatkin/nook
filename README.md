# Nook

Nook is a small terminal program that is driven by nook script.
Nook script is a simple lisp style programming langauge that can be used to interact with the underlying OS.

# TODO:
* [ ] add git status information to the header
* [ ] add a footer
* [ ] add indentation level to the footer
* [ ] add the current time to the footer
* [ ] figure out a good color scheme to make picking colors easier
* [ ] should nook-script be in it's own repo?
* [ ] add support for more things in nook script.
    * [ ] `cd` to change directories.
    * [ ] `exe` to execute a command with argument.
    * [ ] `cp` to copy files.
    * [ ] `mv` to move files.
    * [ ] support for `-f` and `--flag` as values in the language.
    * [ ] support for `./path` and `/root/path` as value in the language.
    * [ ] `fn` to support functions
    * [ ] `do` for multi-expression expressions
    * [ ] `{"value" :ok}` for tuples.
    * [ ] `{name="lex" status=:ok}` for named tuples.
    * [ ] support type checking and type inference systems at "compile" time.
    * [ ] how dumb is it if the operator for the s-expr can be an expression so that functions can be called just using their identifier?
* [ ] show script errors directly inline with the code.