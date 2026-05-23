# 📊 Math Skills in Go

## 📌 Description

This project is a simple command-line program written in Go that reads a list of numbers from a file and calculates basic statistical values:

* Average (Mean)
* Median
* Variance
* Standard Deviation

Each number in the input file represents one value of a statistical population.

---

## 🎯 Objectives

The goal of this project is to:

* Practice file handling in Go
* Understand basic statistical concepts
* Learn how to structure a program using multiple functions
* Improve problem-solving and algorithmic thinking

---

## 📂 Project Structure

```
.
├── main.go
├── read.go
├── average.go
├── median.go
├── variance.go
├── stddev.go
├── data.txt
└── go.mod
```

---

## 📥 Input Format

The program reads data from a file passed as an argument.

Example (`data.txt`):

```
189
113
121
114
145
110
```

Each line must contain a single integer.

---

## ▶️ How to Run

Make sure you are inside the project folder, then run:

```
go run . data.txt
```

---

## 📤 Output Format

The program prints the results as integers:

```
Average: 132
Median: 117
Variance: 784
Standard Deviation: 28
```

---

## 🧠 Concepts Used

### 1. Average (Mean)

Sum of all numbers divided by the count:

```
average = sum / count
```

---

### 2. Median

* Sort the numbers
* If odd count → middle value
* If even count → average of two middle values

---

### 3. Variance

Measures how far numbers are from the average:

```
variance = Σ(x - mean)² / n
```

---

### 4. Standard Deviation

Square root of the variance:

```
std = √variance
```

---

## ⚙️ Implementation Overview

The program is divided into small functions:

* `readNumbers()` → reads numbers from file
* `average()` → calculates mean
* `median()` → calculates median
* `variance()` → calculates variance
* `stdDev()` → calculates standard deviation

The `main()` function connects everything together.

---

## ⚠️ Notes

* All results are returned as **integers** (rounded down)
* The input file must be provided as a command-line argument
* The program assumes valid integer input

---

## 🧩 Example Run

Input:

```
1
2
2
4
3
1
```

Output:

```
Average: 2
Median: 2
Variance: 1
Standard Deviation: 1
```

---
