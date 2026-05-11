
# 🎨 ASCII Art Justify

**ASCII Art Justify** is a CLI-based tool written in Go for generating beautifully formatted ASCII art banners. It offers flexible alignment options (left, center, right, justify) and supports multiple font styles, making it ideal for terminal-based presentations, decorations, or fun text rendering.

---

## ✨ Features

- 🔄 **Custom Alignments**: Supports:
  - `left` 🠔
  - `center` 🎯
  - `right` 🠖
  - `justify` 📏
- 🖋️ **Banner Font Styles**: Includes `standard`, `shadow`, and `thinkertoy`.
- 📐 **Terminal Width Adaptation**: Dynamically adjusts to your terminal width for a perfect fit.
- 🛡️ **ASCII Validation**: Ensures only valid ASCII characters are processed.

---

## 🚀 Getting Started

### 🔧 Prerequisites

- Install Go on your machine. [Download Go](https://golang.org/dl/)

### 📥 Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/ascii-art-justify.git
   ```
2. Navigate to the project directory:
   ```bash
   cd ascii-art-justify
   ```
3. Ensure the banner font files (`standard`, `shadow`, `thinkertoy`) are in the `banner/` directory.

---

## 💻 Usage

Run the program using the following syntax:

```bash
go run . [OPTIONS] [STRING] [BANNER]
```

### 📋 Options

| Option              | Description                                  | Example                                |
|---------------------|----------------------------------------------|----------------------------------------|
| `--align=<align>`   | Text alignment (`left`, `center`, `right`, `justify`) | `--align=right`                        |
| `--type=<font>`     | Font type (`standard`, `shadow`, `thinkertoy`)         | `--type=shadow`                        |

### 🖼️ Examples

1. **Right Align Example**:
   ```bash
   go run . --align=right "Hello, World!" standard
   ```

2. **Justify Example**:
   ```bash
   go run . --align=justify "ASCII Art Justify!" thinkertoy
   ```

---

## 📂 Project Structure

- **`main.go`**: Entry point of the application.
- **Alignment Functions**: Handles text alignment logic (`left`, `center`, `right`, `justify`).
- **Font Management**: Reads and parses banner font files (`standard`, `shadow`, `thinkertoy`).
- **Terminal Adaptation**: Detects terminal width and adapts output dynamically.

---

## 🛠️ How It Works

1. **Input Processing**:
   - Reads the user input string and ensures all characters are valid ASCII.
2. **Font Mapping**:
   - Loads the ASCII art representation of characters from the specified banner font file.
3. **Alignment Logic**:
   - Aligns the ASCII art text according to the user's selection (`left`, `center`, `right`, or `justify`).
4. **Terminal Output**:
   - Dynamically adjusts the formatted output to match the terminal's width.

---

## 📘 Notes

- Ensure the terminal supports font files located in the `banner/` directory.
- Add custom fonts by following the same structure as the provided files.

---

## 🖇️ Links

- 📘 [Go Documentation](https://golang.org/doc/)
- 🎨 [ASCII Art Reference](https://en.wikipedia.org/wiki/ASCII_art)
