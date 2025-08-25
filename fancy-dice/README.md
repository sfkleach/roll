# Custom Fancy Dice Files

This directory contains examples of custom fancy dice files that can be loaded into the Roll dice application using the `--fancy` command-line option.

## Fancy Dice File Format

Custom fancy dice files use a simple text format:

### Basic Format
- One value per line
- Lines starting with `#` are comments and are ignored
- Empty lines are ignored
- Each line can be either:
  - `name` - Display name only (value defaults to line position starting from 1)
  - `name, value` - Display name with explicit scoring value

### File Naming and Dice Type
The dice type is determined by the **number of valid lines** in the file (excluding comments and empty lines):
- A file with 6 values creates an `f6` dice
- A file with 8 values creates an `f8` dice
- A file with 12 values creates an `f12` dice
- etc.

Custom fancy dice **override** built-in fancy dice of the same type.

## Example Files

### colors.dice (6-sided die with explicit values)
```
# Custom color dice
Red, 3
Blue, 2
Green, 5
Yellow, 1
Purple, 4
Orange, 6
```
Creates an `f6` dice where each color has a specific scoring value.

### animals.dice (6-sided die with position-based values)
```
# Animals without explicit values (will default to position-based)
Cat
Dog
Elephant
Lion
Tiger
Bear
```
Creates an `f6` dice where Cat=1, Dog=2, Elephant=3, Lion=4, Tiger=5, Bear=6.

### fruits.dice (5-sided die with mixed format)
```
# Mixed format - some with values, some without
Apple, 10
Banana
Cherry, 15
Date
Elderberry, 5
```
Creates an `f5` dice where:
- Apple = 10 (explicit value)
- Banana = 2 (position-based, since it's the 2nd line)
- Cherry = 15 (explicit value)
- Date = 4 (position-based, since it's the 4th line)
- Elderberry = 5 (explicit value)

## Usage Examples

### Loading a Single File
```bash
# Load colors.dice and roll the custom f6 dice
./roll --fancy="colors.dice" f6

# Result might be: f6: Purple, Total: 4
```

### Loading Multiple Files with Glob Patterns
```bash
# Load all .dice files in current directory
./roll --fancy="*.dice" f6 f5

# Load all .dice files in fancy-dice directory
./roll --fancy="fancy-dice/*.dice" 2f6 f5

# Load specific pattern of files
./roll --fancy="*color*.dice" f6
```

### Complex Examples
```bash
# Mix custom fancy dice with regular dice
./roll --fancy="*.dice" 2d6 f6 f5

# Use with exclusive dice
./roll --fancy="colors.dice" 3F6

# Use with sorting
./roll --fancy="*.dice" --ascending f6 f5 2d6

# Load from multiple directories
./roll --fancy="fancy-dice/*.dice" --fancy="custom/*.dice" f6 f8
```

## Precedence Rules

When multiple fancy dice sources are available:
1. **Custom fancy dice** (from files) - highest precedence
2. **Built-in fancy dice** (f2, f4, f6, f7, f12, f13, f52) - medium precedence  
3. **Default scoring values** - lowest precedence

If you load a custom `f6` dice, it will override the built-in unicode dice faces f6.

## Tips

- Use meaningful names for your dice values
- Values can be negative numbers if needed
- Position-based values start from 1, not 0
- Custom dice work with all Roll features: exclusive dice, sorting, mixing with regular dice
- Test your files with a single roll first: `./roll --fancy="myfile.dice" f6`