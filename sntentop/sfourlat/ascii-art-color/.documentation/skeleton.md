### 1. main.go

This file handles the initial input and delegates tasks to the other components.
Go

```
package main

import (
    "os"
    "strings"
)

func main() {
    // 1. Get arguments from os.Args[1:]
    
    // 2. Identify the --color= flag
    // Check if the first argument starts with the prefix "--color=".
    // If it does:
    //    - Extract the color name (e.g., "red").
    //    - Validate if the color exists in our map.
    //    - If the format is wrong (e.g., --color red), print a usage message and exit.

    // 3. Determine Substring, Main String, and Banner
    // Logic challenge: How do we know if the user provided a substring to color, 
    // or if they want to color the whole string?
    
    // 4. Load the Banner data (Call GetBannerData from banner.go)

    // 5. Call the Render function (from render.go) with:
    // (mainString, subString, colorCode, bannerData)
}
```

### 2. banner.go

This file focuses on file I/O and preparing the ASCII characters in memory.
Go

```
package main

// GetBannerData reads the requested .txt file (standard, shadow, etc.)
// and returns a map where each character (rune) points to its 8-line ASCII art.
func GetBannerData(bannerName string) (map[rune][]string, error) {
    // 1. Determine the file path based on bannerName (default to "standard.txt").
    
    // 2. Read the file and split it by newline. 
    // Remember: Each character is 8 lines of art plus 1 empty line separator.
    
    // 3. Map ASCII characters (runes 32-126) to the corresponding 8 lines.
    
    return nil, nil
}
```

### 3. render.go

This is where the terminal coloring and the actual printing happen.
Go

```

package main

import "fmt"

// Define the ANSI color codes
var colorMap = map[string]string{
    "red":     "\033[31m",
    "green":   "\033[32m",
    "yellow":  "\033[33m",
    "reset":   "\033[0m",   // Crucial for stopping the color
}

// Render processes the text and prints it line by line.
func Render(mainStr, subStr, color string, banner map[rune][]string) {
    // 1. Split mainStr by "\n" to handle multi-line inputs.
    
    // 2. For each "line" of text, we must print 8 lines of ASCII art.
    
    // 3. COLOR LOGIC:
    // As you loop through the characters of the line:
    // - Check if the current sequence of characters matches the subStr.
    // - If it matches: Print the color ANSI code -> ASCII art -> Reset ANSI code.
    // - If it doesn't: Print error message.
}
```