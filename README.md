A 6502 disassembler tuned for NES/Famicom

# Usage

    dasm config.cfg

A configuration file is required.  This disassembler does not do any heuristics
to attempt to automatically determine how to dissamble anything.  Everything
must be explicitly defined.  Blocks of data not defined as data will be
treated as code.

# Config Format

This project uses a custom configuration format that's based around blocks of
options.  The two block types are `Global` and `Bank`.  Only one Global block
is allowed, and Bank blocks are be repeated for as many banks the rom has.

Each block starts with the type name and contains directives wrapped in curly
braces.  Some directives expect a list of further blocks that are wrapped in
square braces.  For example:

    Global {
        Input     "StudyBox.bin"
        MlbOutput "StudyBox.mlb"

        Architecture 6502
        Comments Full

        Labels [
            { Address $30; Name "Argument_A"; Size 2; }
            { Address $32; Name "Argument_B"; Size 2; }
            { Address $34; Name "Argument_C"; Size 2; }
        ]
    }

Comments are start with a double forward slash and continue to the end of the
current line.  There are no multi-line comments.

Hexadecimal numbers are specified with a preceding `0x` or `$`.

## Config Options

All identifiers and non-string options are case-insensitive.

### Global

| Option         | Type   | Description |
|:---------------|:------:|:------------|
| `Architecture` | ident  | Options are `6502`, `Full6502`, and `SBX`.  Only `6502` is currently implemented. |
| `AsmColumn`    | int    | Column to start the instruction printout.  Default is `4`. |
| `AsmCol`       | int    | Alias of `AsmColumn` |
| `AutoVars`     | bool   | Automatically assign names to variables used in `LDA`, `STA`, etc. Destinations of branches, `JMP`, and `JSR` will always auto-generate labels. |
| `CommentColumn`| int    | Column to start the inline comments.  Default is `20`. |
| `CommentCol`   | int    | Alias of `CommentColumn` |
| `Comments`     | ident  | Options are `None`, `Standard`, and `Full`.  `None` omits all comments, `Standard` outputs defined comments, & `Full` outputs defined comments as well as address and raw byte data for each instruction. |
| `Include`      | string | Additional configuration files to include.  Included files can only contain `Bank` blocks. |
| `Input`        | string | Default input file. |
| `Labels`       | list   | List of labels in the global space (eg, zero page & main system RAM. |
| `MlbOutput`    | string | Output file to write Mesen labels. |
| `Output`       | string | File to write the global label definitions. |
| `VerboseColumn`| int    | Column to start the verbose comments (address and raw bytes).  Default is `40`. |
| `VerboseCol`   | int    | Alias of `VerboseColumn` |
| `Windows`      | list   | List of window definitions. |

#### Labels

Same as Bank labels.

#### Windows

This list defines windows.  If the ROM doesn't use a mapper, this list can be
omitted.

Each window listed here defines an range of memory that can be swapped out with
another portion of ROM or RAM.  When more than one window is visible at a time
(ie, not overlapping), references from one bank to another will generate labels
across banks and will use existing labels from the other bank if applicable.

| Option  | Type   | Description |
|:--------|:------:|:------------|
| `Name`  | string | Name of the window that is used in the Bank blocks. |
| `Start` | uint   | Start address for the window. |
| `Size`  | uint   | Size of the window. |
| `Init`  | string | Name of the bank that will be loaded into this window by default. |

### Bank

Bank blocks define what is actually disassembled.  Anything not covered by a
Bank block is ignored.

| Option    | Type   | Description |
|:----------|:------:|:------------|
| `Address` | uint   | CPU address where the disassembly will start. |
| `Input`   | string | Input file if different than the global input. |
| `Labels`  | list   | List of labels defined in this bank. |
| `Offset`  | uint   | Offset into the input file to start processing. |
| `Output`  | string | Output file that the disassembly is written to. |
| `Ranges`  | list   | List of data ranges in this bank. |
| `Size`    | uint   | Size of the bank. |
| `Windows` | list   | List of window swaps. |

#### Labels

Global and Bank blocks use the same format to define labels.

| Option          | Type   | Description |
|:----------------|:------:|:------------|
| `Address`       | uint   | Start address for the label. |
| `Comment`       | string | Alias for `CommentBlock` or `CommentInline`. See below. |
| `CommentBlock`  | string | A Block comment that is output above the dissambly line at the given address.  The comment will start at the begining of the line. |
| `cb`            | string | Alias for `CommentBlock`. |
| `CommentInline` | string | An inline comment that is output after the instruction and before the raw comment on the same line. |
| `ci`            | string | Alias for `CommentInline`. |
| `Name`          | string | Name of the label. Must be alpha-numeric and start with an alpha. Underscores are the only allowed symbol.
| `ParamSize`     | uint   | Size of the inline prameters for the label.  This number of bytes after at `JMP` or `JSR` to this label will be treated as byte data. |
| `Size`          | uint   | Size of the label. |

Block comments start at the beggining of lines whereas inline comments start
immediately after the instruction on the same line.  Both comment types can be
multi-line.

When using `Comment`, if the value contains newlines `CommentBlock` is set.
Otherwise, `CommentInline` is set.

#### Ranges

Any address range not defined as a range is assumed to be code.  The default
type of a defined range is `Bytes`.

A Range of type `Addresses` will resolve all values to labels and create
auto-labels if needed.  Type `words` does not resolve labels.

| Option          | Type   | Description |
|:----------------|:------:|:------------|
| `Address`       | uint   | Start address for range. |
| `Comment`       | string | Create a label with this comment for the whole defined range. |
| `Display`       | ident  | Number display format: `Binary`, `Decimal`, `Hexadecimal`.  |
| `End`           | uint   | End address for range. |
| `Name`          | string | Create a label with this name for the whole defined range. |
| `RtsLabels`     | bool   | Labels in an RTS Trick lookup table |
| `Size`          | uint   | Size of range. |
| `Stride`        | uint   | Number of values to output per line. A stride of `2` will output two bytes with a `Bytes` range, but four bytes with a `Words` range. |
| `Type`          | ident  | Type of range: `Addresses`, `Bytes`, `Words`. Default for defined ranges is `Bytes`. |

#### Windows

Entries here do not define the address ranges for each window.  Instead,
entries define when to change what is assigned to each window.  When an
entry's address is reached during disassembly, swap the given window to the
given bank before continuing disassembly.

| Option    | Type   | Description |
|:----------|:------:|:------------|
| `Address` | uint   | Address to change a window assignment. |
| `Bank`    | string | Name of the bank to assign. |
| `Window`  | string | Name of the window to assign the bank to. |
