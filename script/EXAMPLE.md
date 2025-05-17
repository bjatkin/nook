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
(type Person {dict name str, age int})
(let a {dict name "alexis", age 32})

# slice types and array types are specificed using []
(let names {[str] "Alexis" "John" "Jill"}) # string slice
(let ages {[int 3] 32 38 53}) # int array with length 3
```

### scope

```
(let greet (fn [name str, lang str]
    (match lang
        (case :en {print "hello" name})
        (case :de {print "guten tag" name})
        (case :vi {print "xin chao" name})
    )
))
```

### ideas again...

```
(let name "hello")
(let greet "xin chao")
(let greet_fn (fn [name: str, lang: atom]
    (match lang
        (case :en (print "hello" name))
        (case :de (print "guten tag" name))
        (case :vi (print "xin chao" name))
    )
))

(type Person [name: str, age: int])
(let bob {Person name: "Bob", age: 43})

# postfix atoms are also allowed...
(let new_person (fn [name: str, age: int] {Person name: name, age: age}))
(let people {[Person] bob {name: "Jil", age: 27}})

# this creates a tuple
(let tup {0 :ok})

# the rules are
{0 1 2}                       # if the first operand is a value, it creates a tuple
{name: "hello", lang: 'de}    # if the first operand is a property then it creates a dict typed by it's fields
{name: str, lang: atom}       # dicts can contain types as value, however those types must be resolved at compile time
{Person name: "Bob", age: 27} # if the first operand is a type then it constructs a new value of that type
{[int] 0 1 2 3 4 5 6}         # this can be used to create slice literals

[bob name] # get the name property from bob
[people 3] # get the third person

# the rules are if the first operand is an expression, it accesses fields
[repo status]
# if the first operand is a property it creates a list of function arguments
[n: int]

# in here is the incomming arguments as a tuple
(let fib (fn [n: int] int) (match in
    (0 0)
    (1 1)
    (true (do
        (let a (- n 1))
        (let b (- n 2))
        (+ a b)
    ))
))

# unlike in other lisps, let does not create it's own scope.
# instead, each paren creates it's own scope, but let operates on it's parent scope
# this means that you can use `do` for example to build up a set of variables gradually

# the type of a can be omitted here and it will be infered as a `trait`
# that matches any type that has a associated `+` function
(let add_one (fn [a:] (+ a 1)))

(trait String
    # trait is any type where there is some function string whos first argument
    # is the type in question and returns a str
    string [s:] str,
)

# types get evaluated at compile time meaning the types in the
# associated dict will be handled at complie time so there is no syntax error
(type Person {name: str, age: int})
```

### name spaces / modules

I think that modules should be pretty small in nook in the same way they are in Elixir.
Then you have have name spaces that are larger and contain the smaller modules.
Either that or you can have modules that contain modules.

```
{mod vec
    (let new2 (fn [x y] {.x x, .y y}))
    (let new3 (fn [x y z] {.x x, .y y, .z z}))
    (let new4 (fn [x y z w] {.x x, .y y, .z z, .w w}))

    (let add (fn [a b] (match in
        ([{0 0} {.x .y}] b)
        ([{.x .y} {0 0}] a)
        ([{.x .y} {.x .y}] (new2 (+ [a .x] [b .x]) (+ [a .y] [b .y])))
        ([{.x .y .z} {.x .y .z}] (new3 (+ [a .x] [b .x]) (+ [a .y] [b .y]) (+ [a .z] [b .z])))
        ([{.x .y .z .w} {.x .y .z .w}] (new3 (+ [a .x] [b .x]) (+ [a .y] [b .y]) (+ [a .z] [b .z]) (+ [a .w] [b .w])))
    )))
}

# the type of this value is {.x int .y int}
(let a (vec.new2 0 1))

# the type of this value is {.x int .y float}
(let b (vec.new2 0 1.0))

# this fails at complie time becuase .y of int and .y of float can not be directly added
# and the traits that constrain the types of a and b are
# (trait [a .x]
#     + <fn [T <typeof [b .x]>] any>
# )
# (trait [b .x]
#     + <fn [<typeof [a .x]> T] any>
# )
# (trait [a .y]
#     + <fn [T <typeof [b .y]>] any>
# )
# (trait [b .y]
#     + <fn [<typeof [a .y]> T] any>
# )
# and there is no function with a signature that is compatible with <fn [int float] any> named '+'
(vec.add a b)

(let c (vec.new2 3 4))

# this is successful because there is a function named '+' that has a compatable type singature
# namely <fn [int int] int>, it is technically more narrow than <fn [int int] any> but that is
# allowed. This also means we know the return type of this call is {.x int .y int} rather than
# {.x any .y any} which would be the more generic case
(vec.add a c)


(trait String
    # T here becomes a generic type referencing the implementation of the trait
    string <fn [T] str>
)

(trait Add
    # T can be reused for other types
    add <fn [T T] T>
)
```

```
(type Vec2 <.x, .y int>)
(ns vec2
    (let add (fn [a b]
        ({vec2 
            .x (+ [a .x] [b .x])
            .y (+ [a .y] [b .y])
        })
    ))
)

(type Vec3 <.x, .y, .z int>)
(ns vec3
    (let add (fn [a b]
        ({vec3 
            .x (+ [a .x] [b .x])
            .y (+ [a .y] [b .y])
            .z (+ [a .z] [b .z])
        })
    ))
)

(trait add [T]
    <fn add [T T] T>
)
```

```
(impl add (fn [a b] (+ a b)))

(type Vec2 <.x, .y int>)
(impl add (fn [a b] 
    {Vec2
        .x (+ [a .x] [b .x])
        .y (+ [a .y] [b .y])
    }
))

(type Vec3 <x, y, z int>)
(impl add (fn [a b] 
    {Vec3
        .x (+ [a .x] [b .x])
        .y (+ [a .y] [b .y])
        .z (+ [a .z] [b .z])
    }
))
```