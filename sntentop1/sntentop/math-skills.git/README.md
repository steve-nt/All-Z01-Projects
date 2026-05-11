# README

## Overview

This project is a simple program designed to calculate basic statistical measures: **Average**, **Median**, **Variance**, and **Standard Deviation**. It reads numerical data from a file and performs the calculations, printing the results in a user-friendly format. The program is implemented in both **Python** 🐍 and **Go** 🐿.

There is a script in both folder that automates the following tasks you can read more about it further down:

- 📥 **Downloads** a zip file containing the `stat-bin` project.
- 📂 **Unzips** the downloaded file.
- 📋 **Copies** additional files (`mathskills.go` and `commands-to-run.sh`) into the `stat-bin` directory.
- ⚙️ **Executes** the `commands-to-run.sh` script.
- 📋 **Prompts** the user to clean up files and directories starting with `stat-bin`.

---

## Features

- Calculates **Average**
- Calculates **Median**
- Calculates **Variance**
- Calculates **Standard Deviation**
- Handles input from a file containing numerical data (one value per line)

---

## Prerequisites - Requirements

To run this program, ensure you have the following installed:

- A programming environment for **Go** 🐿 or **Python** 🐍
- Basic command-line interface (CLI) knowledge
- **VS Code** or another editor
- 🖥️ **Bash**: The script is written for Bash and requires a Linux or macOS environment (or WSL on Windows).
- 🌐 **wget**: For downloading files.
- 📦 **unzip**: For extracting the zip archive.

---

---

## Usage

### 1. Clone or Download the Repository

```bash
git clone https://github.com/StephanosNt/math-skills.git
```

### 2. Prepare Your Data File

Place your data file in the same directory as the program (or provide the correct path to it). Or use the bash script

### 3. Run the Program

For **Go** 🐿:

```bash
go run math_skills.go data.txt
```

For **Python** 🐍:

```bash
python3 math_skills.py data.txt
```

### Example Input File (data.txt):

```text
189
113
121
114
145
110
...
```

### Example Output:

```text
Average: 132
Median: 121
Variance: 627
Standard Deviation: 25
```

## Notes

- Results are rounded to the nearest integer.
- 🌐 Ensure you have an active internet connection for downloading the zip file.
- 📋 The cleanup prompt helps avoid accidental deletions.
- ✏️ Review and customize the script as needed for your specific requirements.

---

## How It Works

### File Input:

The program reads numerical data from the file specified as an argument.

### Calculations:

1. **Average:**

   - Sum of all values divided by the count of values.

2. **Median:**

   - Middle value in the sorted data (or average of the two middle values if count is even).

3. **Variance:**

   - Measure of how spread out the data is, calculated as the average of squared differences from the Mean.

4. **Standard Deviation:**

   - Square root of the Variance.

### Output:

The program prints results for each calculation in a readable format.

---

## Testing

- Run the program with different data files to verify the results.
- Compare the output with manually calculated values or use a third-party calculator to ensure accuracy.

---

## Contributing

If you'd like to contribute to this project:

1. **Fork** this repository.

2. Create a feature branch:

   ```bash
   git checkout -b feature/AmazingFeature
   ```

3. Commit your changes:

   ```bash
   git commit -m 'Add some AmazingFeature'
   ```

4. Push to the branch:

   ```bash
   git push origin feature/AmazingFeature
   ```

5. Open a pull request.

---

## License

This project is licensed under the **MIT License**.
---

## Author

**Stefanos Ntentopoulos and ChatGPT**

Feel free to reach out with any questions or suggestions!

---

Happy coding! 🚀




## Script Usage Instructions

### 1. Prepare the Environment

🛠️ Make sure `wget` and `unzip` are installed:

```bash
sudo apt-get install wget unzip   # For Debian-based systems
sudo yum install wget unzip       # For Red Hat-based systems
```

✅ Ensure `mathskills.go` and `commands-to-run.sh` are in the same directory as the script.

### 2. Run the Script

🏃 Make the script executable:

```bash
chmod +x your_script_name.sh
```

🚀 Execute the script:

```bash
./your_script_name.sh
```

### 3. Script Execution Workflow

The script will:

1. 📥 **Download** the zip file from the specified URL.
2. 📂 **Unzip** the file and list the extracted contents.
3. 📋 **Copy** `mathskills.go` and `commands-to-run.sh` into the `stat-bin` directory.
4. ⚙️ **Execute** the `commands-to-run.sh` script inside the `stat-bin` directory.
5. ❓ **Prompt** you with the following question:
   ```bash
   Do you want to remove all files and folders starting with 'stat-bin'? (yes/no):
   ```

### 4. Respond to Cleanup Prompt

- ✅ Enter `yes` or `y` to delete `stat-bin` files and directories.
- ❌ Enter `no` or `n` to keep the files.

### 5. Verify Output

🔍 Check the logs printed during the script's execution.

- ✅ Confirm that the `commands-to-run.sh` script completed its tasks.
- 📋 If cleanup was performed, ensure that files starting with `stat-bin` are removed.

---
## Troubleshooting

### 🛑 Permission Denied:

Ensure the script has execution permissions:

```bash
chmod +x your_script_name.sh
```

### 🔍 Command Not Found:

Verify that `wget` and `unzip` are installed and accessible in your PATH.

### ❓ File Not Found Errors:

Confirm that `mathskills.go` and `commands-to-run.sh` exist in the same directory as the script.

---

