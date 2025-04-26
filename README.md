# Nook

Nook is a small terminal program that is driven by nook script.
Nook script is a simple lisp style programming langauge that can be used to interact with the underlying OS.

![screenshot](/assets/screenshot.png)

# TODO:
* [x] add git status information to the header
* [x] add a footer
* [x] figure out a good color scheme to make picking colors easier
* [ ] spend some time thinking through debugging so the system is better. It's really ad-hoc right now.
* [ ] take another pass through the layout code. It probably could use a little more love.
* [ ] give colors better names. The current ones aren't really useful.
* [ ] do I create a new github for this?
* [ ] add indentation level to the footer
* [ ] should nook-script be in it's own repo?
* [ ] do a round of cleanup on the VM code. It's all over the place right now.
* [ ] add support for more things in nook script vm.
    * [x] `cd` to change directories.
    * [x] `ex` to execute a command with argument.
    * [x] `ls` to get files in a datstructure that nook script can work with.
    * [x] support for `-f` and `--flag` as values in the language.
    * [x] support for `./path` and `/root/path` as value in the language.
    * [ ] `cp` to copy files.
    * [ ] `mv` to move files.
    * [ ] `mkdir` to make a directory
    * [ ] `touch` to make a file
    * [ ] `fn` to support functions
    * [ ] `do` for multi-expression expressions
    * [ ] `{"value" :ok}` for tuples.
    * [ ] `{name="lex" status=:ok}` for named tuples.
    * [ ] support type checking and type inference systems at "compile" time.
    * [ ] import support
    * [ ] configuration for VM startup (writtein in nook-script so it's like IEx)
    * [ ] how dumb is it if the operator for the s-expr can be an expression so that functions can be called just using their identifier?
        (It's probably fine if we do type checking ahead of time)
* [ ] show script errors directly inline with the code.
* [ ] just a way better error system in general would be great.
* [ ] simple syntax hilighting would be great.