# Language

A simple language of square brackets where you never have to press the shift keys. It has the following features:
- Option instead of null.
- Explicit types.
- Functions are pure and only evaluated when used.
- Procedures cause side-effects.

TODO:

- What if a value from an imported file gets overwritten in the global environment. Shouldn't each file have its own environment which includes the items specified in the import statement as well as the ones created in the file and nothing more. Except that we don't want to import the same file twice so these references should be re-used in multiple environments.
- After adding type defs to function and procedure parameters also add function overloading based on parameters. So if multiple functions with the same definition are called it should automatically pick the right one to use. In addition like Elixir it should support the "when" clause to prevent early returns with guard clauses.
- Add an optimise step for after the parser that removes dead code.
- Add a step that can be ran after the parser that checks if the program is valid.
- Add structs that have a predefined list of keys. It's syntax should be `[define fileData [struct [integer size-in-bytes] [string name] [string extension]]]`.
- Create values using their types. For example `1.0` should be `[float 1.0]`.
- How are errors handled? They should be values.
- Allow for multiple return types.

```
[define calc-fib
  [function int / returns an integer.
    [int n
      [greater n 0] / if n is 0 or lower then this function should not match and an overload ought to be specified.
    ]
    [match n
      [0 0]
      [1 1]
      [_
        [add
          [calc-fib
            [subtract n 1]
          ]
          [calc-fib
            [subtract n 2]
          ]
        ]
      ]
    ]
  ]
]
```

```
[define unwrap
  [function
    [some[int] value] / Is this how types of some are specified?
    [match value / How do I check for some?
      [some [some value]] / How do I unwrap a value?
      [none 1]
    ]
  ]
]
```