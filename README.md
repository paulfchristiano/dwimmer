# dwimmer

To get: (hopefully) go get github.com/paulfchristiano/dwimmer/repl

To start REPL: with mongod running locally, run repl

### Commands:

- view X ("view x"): View the variable with name X.
- ask Q ("ask what is 3 + 4?"): Ask the question Q. 
The return values are the result and @C, where C is a channel that can be used for follow-up questions.
- reply A ("reply 7"): Send the message A to your parent process.
- ask@C Q ("ask@y how reliable is that answer?"): Send the message Q to the channel bound to X.
- delete X ("delete x"): Delete the variable named X. If there are no remaining references to that term, the GC can free it.
- replace N X ("replace 6 the answer to the question is [A]"): 
Replaces output #N with X, removing all of the subsequent inputs and outputs.
- correct N ("correct 7"): Remove the rule that produced input N, and provide a new action to be used instead.

### Input format:

- Integer literals are entered as: 7
- String literals are entered as: "asdf"
- Variables are entered as: x (or #x)
- A compound expression is a sequence of words, interspersed with subexpressions.
- A variable name, quoted string, or integer will automatically be interpreted as a subexpression. A compound expression can be enclosed in [] or () to be made into a subexpression.

E.g:

- What is the sum of squares of the elements of [the list with first element 7 and second element [the length of "hello"]]?
- What is the list formed by repeating each element of x 4 times?

### Notes

Essentially nothing will work out of the box; all information about how to answer queries is stored in the mongo instance.

If you press [up] while entering input, you may receive suggested replacements for the most recent output.

The current version is a rough prototype. Contributions are likely to be discarded.
