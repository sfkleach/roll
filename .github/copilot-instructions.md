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

Add an unobtrusive way to switch between dice order and sort-by-value. This 
will be a simple 'hotkey' toggle, allowing users to quickly switch modes without
having to click buttons.

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