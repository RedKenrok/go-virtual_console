# Language

A simple functional language.

TODO:

- Don't import the same file twice.
- What if a value from an imported file gets overwritten in the global environment. Shouldn't each file have its own environment which includes the items specified in the import statement as well as the ones created in the file and nothing more.
- Add distinction between functions and procedures. Procedures have side effects and are called in place. Whereas functions should only be evaluated when their value is used.
- Add an optimise step for after the parser that removes dead code.
- Add a step that can be ran after the parser that checks if the program is valid.
- Add structs that have a predefined list of keys. It's syntax should be `[define fileData [struct [integer size-in-bytes] [string name] [string extension]]]`
- Create values using their types. For example `1.0` should be `[float 1.0]`.