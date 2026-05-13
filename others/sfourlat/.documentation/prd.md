### PRD: ASCII Art Color 
The goal of this project is to extend the current functionality to support terminal coloring using flags. We will focus on modularity and user experience.

# Project Objectives
 - Flag Parsing: Implement the ability to detect --color=<color>.

 - Substring Targeting: Identify and apply color only to specific parts of the string if requested.

 - Color Systems: Choose a system (like ANSI escape codes) to render colors in the terminal.

 - Error Handling: Enforce strict usage messages as defined in the requirements.


# Proposed Logic & Implementation
To get started, we need to think about how the program will distinguish between the color,the substring, the main text and the banner.

In your current main.go, you handle arguments like this:
go run . [STRING] [SUBSTRING] [BANNER]

For the new version, the arguments might look like:
go run . --color=red "hello" standard

#  Step 1 : Argument Prccesing

    - HANDLING THE COLOR FLAG:
        To support the new requirement, we need to check if an argument starts with --color=.

    We need logic to : 
    1. Identify that --color=red is the flag.
    2. Extract the string "red" from that flag.
    3. Check if the value of the color is right.
    

# Step 2 : Terminal color implementation

    In a terminal, we change text color using Ansi Escape Codes. These are special sequences of characters the the terminal doesn't print literally, but it interprets them as commands to change the text's appearance.
    A standard Ansi color code looks like this : \033[<COLOR_CODE>m.
        - Escape Character : \033 tells the terminal a command is starting
        - Color Code: A number represanting the color 
        - reset: Is very important to use a reset code at the end,because the entire color will stayu for ever.

 # Step 3 : Changes to Printing

    Currently, your code uses a RENDER function to print 8-line characters.To add color, we don't change the ascii art itself, but we simply wrap the output string with the ANSI codes.

        we would print :

        \033[31m@@@@@@\033[0m (This makes it Red)

    So we have to make map to save the codes of each color to search on it: 
 ```
    var colorMap = map[string]string{
    "red":     "\033[31m",
    "green":   "\033[32m",
    "yellow":  "\033[33m",
    "blue":    "\033[34m",
    "magenta": "\033[35m",
    "cyan":    "\033[36m",
    "white":   "\033[37m",
    "reset":   "\033[0m",
}
```


## EDGE CASES

    1. Flag Formatting Errors 
        The instructions state the flag must have exactly the same format: --color=<color>.

        Missing Equals: What if they type --color blue instead of --color=blue?

        Typochs: What if they type --colour=blue or -color=blue?

        Empty Color: What happens with --color=?

    2. Substring Ambiguities 
        This is where the logic gets tricky because the substring is optional.

        Substring Not Found: If the user runs go run . --color=red "z" "hello", but "z" isn't in "hello," what should the output look like?

        Multiple Occurrences: If the user wants to color "l" in "hello," should both 'l's be colored or just the first one?

        Overlapping/Identical Strings: What if the substring and the main string are the same?

    3. Terminal & Banner Limits 
        Empty Strings: go run . --color=red "".