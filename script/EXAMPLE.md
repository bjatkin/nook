# Nook Script Reference

Bellow are examples of nook script.
They should be used to guide the developement of the compiler and VM.

### basic syntax

nook takes inspiration

```
# comments start with a hash like in python
# expressions are lisp based
(+ 1 2 3 4)

# let can be used to set variables in the parent scope
(let a "hello world")

# variables can be shadowed but types are not dynamic.
# nook script is fully type checked and relies on type inference.
# this means the developer need no annotate types where they are obvious
(let a 10)

# commands can be executed using the ex operator.
(ex "git" "status")

# command names and arguments can be represented by string literals, but
# expressions, atoms, flags, and paths are also supported.
# see the buitin types section for more information on the basic data types in nook script.
(ex 'ls -la .)

# nook also supports pipelining which is popular functional languages.
# this is different from bash pipelining which link stdin and stdout together between multiple commands.
"hello world" | (print)

# `ex` also has a variant `pex` where the first argument is an io.fdset, which is an io file descriptor set.
# io.fdset is a named tuple (see the datstructures section for more information on tupes), and contains
# the following fields (stdin:io.fd, stdout:io.fd, stderr:io.fd). io.fd is the io file descriptor type
# which you can learn more about by visiting the type system section. This can be used to change where
# the command writes to and reads from. It runs the command it a co-routine and returns it's io.fdset
# with it's stdout file descriptor set to the stdin file descriptor.
(pex {stdin=io.stdin, stderr=io.stderr, stdout=io.dev_null} 'echo "hello?")

# you can use `pex` to replicate bash pipelines
# notice io.stdfd starts off the pipeline so the first command uses stdin/stdout/stderr
(pex io.stdfd 'ls -la) | (pex 'grep "*.go")

# this is equlivilant to thi following nested expression but the pipelined version is
# much easier to read/ reason about
(pex (pex io.stdfd 'ls -la) 'grep "*.go")

# you can also define functions in nook like in most languages
# this expression creates a closure that returns the sum of two integers
(fn {a:int b:int} int (+ a b))

# you can also omitt types and the compiler will infer types for you
# here all the types that support the + operator can use this function
# notice the colon is still required to indicate type expressions, even if
# no type is explicitly provided
((fn {a b} auto) (+ a b))

# functions can be bound to identifiers using `let` since they are just normal
# values in nook, but they can also be bound using the `bind` operator.
# this works like the `let` keyword but the variable can not be shadowed in the current scope. 
(bind div (fn {a b} auto) (/ a b))

# all values in nook follow the format value:type. However, because types and values are almost always
# unambigious the inferred portion of the value can be omitted. For example, these are equivilant
0:int
0

# match args can be used to pattern match against the function arguments
# args is bound to the function inputs.
(bind fib (fn {n:int} int)
    (match args
        {0} 0
        {1} 1
        {n} (+ 
            (fib (- n 2)
            (fib (- n 1))
        ))
    )
)
```

### builtin types

nook has all the basic datatypes you would expect, plus a few special ones that make working with the operating system smoother.

```
# nook has both integer and floating point types
# these are 64 bits by default
(fn {a:int b:int} int (+ a b))
(fn {a:float b:float} float (+ a b))

# nook suppots strings and bools as well
(let a:bool true)
(let b:str "hey ya!")

# nook supports atoms as basic literals
(let a:atom 'ok)

# nook also supports paths as literals to make working with cli tools easier
# for a path to be valid it must start with `/`, `./`, `..` or be the exact value `.`
# this means that path values like `my/dir` are not valid, but allowing such paths
# would cause ambiguity with identifiers in some cases, so they are not permitted.
(let home /User/me/)
(let project ./src/project)
(let cwd .)
(let prev ..)

# nook also has support for `flag` values.
# flags must start with `-` and the second character must be a `-` or an alpha character.
# this is to prevent confusion between negative numbers like `-1` and flags. In the case where
# a flag requires a leading number (e.g. a tool uses a `-0` flag or something similar) a string
# value can be used instead and then cast to a flag type.
(let version --version)
(let recursive -r)

# because -0 is not a valid flag (it would be parsed as a negative number) it can be represented by
# a string literal. :flag here is a type annotation which tells the complier to convert the string
# literal into a flag literal which can be done since flag types can be constructed from string literals
(let level0 "-0":flag)
```

### type system

Nook has strong type inference capabilities meaning the developer need not explicitly identify types when the type is obvious.
There are times where indicating the type is desired though and doing so can make code more clear.
Additionally, nook alows users to define their own types.

```
# nook has type expressions which must be resolved at compile time.
# these expressions are ALWAYS started with a :
# here int is a simple type expression that resolves the type of a to float.
# this tells the compiler to convert the constant 0 into a floating point number.
(let a:float 0)

# this is also allowed but here a will be an in because `0` has the implied type int.
# this is true even though it can represent other data types.
(let a 0)

# nook supports type aliases
(type Age int)

# nook also supports creating struct types. These can also be aliased
(type v2 {x:int y:int})

# you can define union types which are useful for things like result types
(type MaybeInt (union int none))

# nook also support enums
(type RGB {int int int})
(type Color (enum
    Red:RGB
    Blue:RGB
    Green:RGB
    Other
))
```

### trait system

### data structures

Nook supports many datastructures that can be used to group data together.

```
# {} creates structured data in nook. In this example two values are provided
# this creates a "tuple" where values can be retrived by index
(let result {"Jane" 'ok})

# this results in "Jane"
result.0

# {} can also be used to create named tuples. Here the values can be accessed by
# their identifier
(let pet {name="Daisy" type='dog age=3})

# this results in 3
pet.age

# {} can be used to create a struct type as well
# here the `type` keyword is used to create a type alias
# notice that type names are upper case. This is the nook naming convention
(type Result {int atom})

# this type checks
(let a:Result {100 'ok})

# this does not because the 0th element is not an int
(let b:Result {false 'ok})

# juse like {} can be used to create named and unamed tuples it can be used to
# create named and unnamed types
(type Order {name:str item:str price:float})

# this type checks
(let coffee:Order {name="Alexis" item="Coffee" price=5.49})

# so does this
(let muffin:Order {"Alexis" "Blueberry Muffin" 3.19})

# but this does not because items is not a valid key
(let burger:Order {name="Alexis" items=["burger" "fries"] price=12.99})
```

### () vs [] vs {}

() is for executing code
[] is for accessing elements and properties
{} is for creating data structures and types

```
# running a function
(print "hello world")

# run a builtin
(+ 1 2 3 4)

# run a command
($git 'status)
```

```
# create a slice of integers
(let ids {1 2 3 4 5}:[int])

# get the first element of the tuple
[ids 0]

# create a data structure
(let person {name "bob", age 75})

# get the persons name
[person name]
```

```
# create a slice where the first value is a number and the second is an atom
(let return {tup 1 'ok})

# create a struct
(let person {dict name "alexis", age 32})

# create a struct type
(type Person {dict name:str, age:int})
```

### types?

```
# types can be omitted
(let a 10)

# values of a specific type can be created with {} 
(let a {int 10})

# types can also be cast from one type to another
(let a (cast {int 10} float))

# use the dict keyword to create a dict, typed by it's fields
# this looks a little weird but should be fine once syntax hilighting
# kicks in
(let a {dict name str, age int})

# slice types and array types are specificed using []
(let names {[str] "Alexis" "John" "Jill"}) # string slice
(let ages {[int 3] 32 38 53}) # int array with length 3
```