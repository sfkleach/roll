# About this Application

This is a fun project to build an application in Go and Fyne, which will 
allow you to roll virtual dice, save dice sets, and load them for later use.

## Project Standards

### Temporary Files

- VSCode gets confused by temporary files too easily. And when you try 
  deleted them it often instantly recreates them. So you must always check
  that they are gone after a few seconds have elapsed (i.e. add sleep to
  the command).

- When you need to create temporary files, avoid creating them in the repo 
  folder. It is OK to create them in `/tmp`.

### Coding Guidelines

- Comments should be proper sentences, with correct grammar and punctuation,
  including the use of capitalization and periods.

- Where defensive checks are added, include a comment explaining why they are
  appropriate (not necessary, since defensive checks are not necessary).

## Version 1

You specify the number and rank of dice to roll in a compact expression on a 
single line. When you press the roll button, the application will simulate the
rolling of the dice and display the results with each dice showing its individual
roll value and also the total roll value.

For example, to roll three six-sided dice, you would enter "3d6" in the input
field and press the roll button. The application would then display the 
individual roll values for each die, as well as the total roll value. The
input field would remain the same, allowing you to quickly roll again with
the same parameters.

It should accept the following expressions:

- "d20" for rolling a single twenty-sided dice
- "2d10 d6" for rolling two ten-sided dice and one six-sided die
- "1d20,7d4" for rolling one twenty-sided die and seven four-sided dice

The input field is initially set to greyed out text "e.g. 2d6" which disappears
as soon as the user clicks into the field.

The results are displayed in a 2-column grid. Each individual dice rolls are
shown on their own row, in the order typed by the user. The left column is the
number of faces of the dice, prefixed by the dice character "d", so a six-sided
dice is shown as "d6". The right column is the individual roll value for that
dice. After all the individual dice rolls, a final row shows the total roll value.

I want to be able to run it on my Linux desktop.

## Version 1.1

Some "fancy"" dice are given custom unicode characters:

- "f2". a two sided dice that prints "heads" or "tails"
- "f4", a four sided dice that shows the unicode suit characters
- "f6", a six sided dice that shows the unicode characters for the six faces
- "f7", a seven sided dice that shows the days of the week (Mon, Tue, Wed, Thu, Fri, Sat, Sun)
- "f12", a twelve sided dice that shows the unicode characters for the zodiac signs
- "f13", a thirteen sided dice that shows one of A, 2, 3, 4, 5, 6, 7, 8, 9 10, J, Q, K
- "f52", a 52 sided dice that shows unicode playing card symbols

## Version 1.2

Now for _exclusive_ dice. When you write a group like 3D6 or 5D20 or even 13F52, 
it means that the group will not roll the same number twice! In other words, 
there will be no repeats within that group. This only applies to a single
group.

## Version 1.3

You can run roll from the command line, passing the dice expressions as arguments.
In addition you can pass a --ascending or --descending flag to sort the individual 
dice rolls in ascending or descending order by value.

## Version 1.4

In this version we lean into the idea of it being both a command-line utility
and a GUI utility. Firstly make the cheatsheet pop up in its own separate 
window. At the moment it is being overwritten by the OK button. The cheatsheet
should be printed to the terminal when the --help option is used.

Secondly Add a --version flag to display the current version and 
arrange that this is baked in when we do a release build. The version is 
displayed in the cheatsheet.

In addition add the version to the cheatsheet which pops up when the info (i)
button is pressed. 

## Version 1.5

In this version we both streamline and improve fancy dice. The key idea is
that fancy dice are represented by a mapping from an index to a unicode string
and a value. 

Firstly we streamline the existing implementation of fancy dice to get rid of
the scaling. But we should continue to use Monotype.

Secondly we add a --fancy=GLOB option, which specifies a set of files with 
the following format:

```txt
# coins.fancy
tails, 0
heads, 1
```

```txt
# card_count.fancy, useful for contract bridge.
2, 0
3, 0
4, 0
5, 0
6, 0
7, 0
8, 0
9, 0
10, 0
J: 1
Q: 2
K: 3
A: 4
```

Lines starting with `#` are removed. Whitespace-only lines are removed. There
must be at least one non-blank line. The lines are comma-separated with one or
two fields (name, value). The name will be used as the fancy text and the value
will be used as its scoring value. If the value is omitted then the value 
becomes its position in the list (starting from 1).

So `coins.txt` could equally have been written as:

```txt
tails
heads
```


## Version 2

You can save a dice set for later use by clicking the save button. This will 
store the current dice configuration, allowing you to easily recall and roll 
it again in the future. You simply open the saved dice sets panel and run it 
from there with a single click or load it into the input field for further 
modification.

## Version 3

I want to be able to run it on my 

- MacOS laptop/desktops
- Windows 11 desktops 
- Android phone
- iOS

# Version 4

I want to be