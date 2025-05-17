# Nook Script Spec

Nook script is a lisp style functional language.

### basic syntax

Nook uses lisp style s-expressions for most of it's syntax.
Expressions start with parens and an operator.

```
(+ 1 2 3 4)
(print "hello", "world")
(let a 10)
```

Notice that the `print` expression includes a `,` character.
In nook `,` is considered a whitespace character and is ignored.
It's not typical to use `,` in a function call like this but there are other more complex cases where `,` is useful.

### basic values and types

Nook includes the following built-in types
* any - technically the empty trait. All values fit the any type.
* int - a 64-bit integer value
* float - a 64-bit floating point value
* bool - a boolean value, either `true` or `false`
* str - an immutable string value
* byte - an 8 bit value, the underlying type of str
* path - a file path. Must start with `.` or `/`
* flag - a command flag. Must start with `-`
* atom - like atoms from elixer or other functional language. Must start with `'`
* cmd - a command line function. Must start with `$`
* none - like nil, it's the empty / nothing value

New types can be constructed with curly braces, even basic types.

```
{float 1} # creates a floating point number even though 1 is an integer literal
{any true} # create an any value that wraps a bool
{flag "--test"} # create a flag
{none} # create an untyped none value
```

### tuples

Nook allows users to create tuple data as well as tuple types.
Tuples are created using `{}`.

```
(let point {5 10})
```

Tuple types are created using `<>`.

```
(type Point <int int>)
```

Specific tuple values can be accessed by index using `[]`

```
[point 0] # evaluates to 5
[point 1] # evaluates to 10
```

### dictionaries

Nook can also create dictionaries and dictionary types as well.
Dicts are created using `{}`.

```
(let tv_show {.title "buffy: the vampire slayer", .year 1996})
```

Ditionary types are created using `<>`

```
(let tv_show <.title str, .year int>)
```

Specific dict properties can be accessed by index using `[]`

```
[tv_show .title] # evaluates to "buffy: the vampire slayer"
```

Dicts and tuples are differatiated based on the leading element of the expression.
If the leading element is a value (e.g., 0, true, 3.14), a tuple is constructed.
If the leading value is a property (e.g., .name, .title), a dict is constructed.

### slices 

Nook also supports slice literals.

```
(let ints {[int] 5 10 15 20})
```

slice type are named using square braces around the underlying type.

```
[int] # this is an int slice
[bool] # this is a bool slice
[tv_show] # this is a slice of tv_shows
```

`[_]` is a slice type where the underlying type of the slice is infered based on the underlying values.

```
# type of this slice is [str]
{[_] "sleepy" "dopy" "doc"}

# type of this slice is [any]
{[_] 3.14 true 'ok}
```

Slice values can be accessed by index using `[]`

```
[ints 0] # evaluates to 5
```

### functions 

Given that Nook is a lisp variant it embraces it's functional roots.
The `fn` operator can be used to create function literals.
Function paramaters are defined using `[]` and take a single expression for their body.

```
(fn [a b] (+ a b))
```

Named functions are possible using the `let` operator.

```
(let add (fn [a b] (+ a b)))
```

By default functions do not need to include types but they can be typed if needed.
*see the section on Type Inference for more details on how function types are handled*

```
(let add1 (fn [i int] (+ i 1)))
```

s-expressions in Nook consist of an `Operator` and optional `Operands`.
`Operators` in Nook are technically expressions which allows functions to be easily invoked.

```
# evaluaes to {int 200}
((fn [a b] (* a b)) 10 20)
```

### function overloads

Nook supports function overloads that can be used to create multiple implementations of a function.
Overloaded functions are declared with the `impl` keyword.

```
(impl add (fn [a b] (+ a b)))

(type Vec2 <.x .y int>)
(impl add (fn [a b] 
    {Vec2
        .x (+ [a .x] [b .x])
        .y (+ [a .y] [b .y])
    }
))

(type Vec3 <.x .y .z int>)
(impl add (fn [a b] 
    {Vec3
        .x (+ [a .x] [b .x])
        .y (+ [a .y] [b .y])
        .z (+ [a .z] [b .z])
    }
))
```

if the add function were defined using `let` add would be shadowed on each declaration.
With `impl` all the functions are bound as overloaded versions of `add`.
The correct implementation will be chosen at compile time based on the type checker.

# Type Inference

# Controll Flow

Nook supports all the expected control flow types.

```
(if (> a b) (print "a is bigger") (print "a is not bigger"))

(match lang
    ('en (print "English"))
    ('de (print "German"))
    ('vi (print "Vietnamese"))
    (else (print "Unknown"))
)
```

