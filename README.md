# dwimmer

To get: (hopefully) go get github.com/paulfchristiano/dwimmer/repl

To start REPL: with mongod running locally, run repl

### Commands:

view * ("view x"): View the variable with name *.


ask * ("ask what is 3 + 4?"): Ask the question *.


reply * ("reply 7"): Send the message * to your parent process.


tell N * ("tell 2 how reliable is that answer?"): Send the message * to the process that answered question N. Return its reply.


delete * ("delete x"): Delete the variable named x. If there are no remaining references to that object, the GC can free its memory.


close N ("close 7"): Close the process that answered question 7, allowing the GC to free its memory and deleting all of its variable references.


replace N * ("replace 6 the answer to the question is [A]"): 
Replaces output #N with *, removing all of the following inputs and outputs.


correct N ("correct 7"): Remove the rule that produced input N, and provide a new action to be used instead.

### Input format:

Integer literals are entered as: 7

String literals are entered as: "asdf"

Variables are entered as: x

A compound expression is a sequence of words, interspersed with subexpressions.

A variable name, quoted string, or integer will automatically be interpreted as a subexpression. A compound expression can be enclosed in [] or () to be made into a subexpression.


E.g:

What is the sum of squares of the elements of [the list with first element 7 and second element [the length of "hello"]]?

What is the list formed by repeating each element of x 4 times?


### Notes

Essentially nothing will work out of the box; all information about how to answer queries is stored in the mongo instance.

If you press [up] while entering input, you may receive suggested replacements for the most recent output.

