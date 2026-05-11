import sys # This module provides access to some variables used or maintained by the interpreter and to functions that interact strongly with the interpreter.
import math # This module provides access to the mathematical functions defined by the C standard.

def read_data(file_path):  # Define a function that reads data from a file
    with open(file_path, 'r') as file: # Open the file in read mode
        return [int(line.strip()) for line in file.readlines()] #Read numbers from the file, convert them to integers, and return as a list

def calculate_average(data): # Define a function that calculates the average of a list of numbers
    return sum(data) / len(data) # Return the sum of the numbers divided by the count of numbers

def calculate_median(data): # Define a function that calculates the median of a list of numbers
    sorted_data = sorted(data) # Sort the data in ascending order
    n = len(sorted_data) # Get the count of elements in the sorted data
    middle = n // 2 # Calculate the index of the middle element
    if n % 2 == 0: # Check if the count of elements is even
        
        return (sorted_data[middle - 1] + sorted_data[middle]) / 2 # If even number of elements, return the average of the two middle ones
    else: 
        
        return sorted_data[middle] # If odd, return the middle element

def calculate_variance(data, average): # Define a function that calculates the variance of a list of numbers
    return sum((x - average) ** 2 for x in data) / len(data) # Return the sum of squared differences from the average divided by the count of numbers
 
def calculate_std_deviation(variance): # Define a function that calculates the standard deviation from the variance
    return math.sqrt(variance) # Return the square root of the variance

def custom_round(value): # Define a function that rounds a number to the nearest integer
    return math.ceil(value) if value % 1 >= 0.5 else math.floor(value) # Return the ceiling of the number if the decimal part is greater than or equal to 0.5, otherwise return the floor

def main(): # Define the main function
    if len(sys.argv) != 2: # Check if the number of command-line arguments is not equal to 2
        print("Usage: python your_program.py data.txt") # Print a usage message
        return # Exit the program

    file_path = sys.argv[1] # Get the file path from the command-line arguments
    data = read_data(file_path) # Read data from the file

    average = calculate_average(data) # Calculate the average of the data using the sum of the data and the count of numbers
    median = calculate_median(data) # Calculate the median of the data using the sorted data
    variance = calculate_variance(data, average) # Calculate the variance of the data using the average
    std_deviation = calculate_std_deviation(variance) # Calculate the standard deviation from the variance 

    print(f"Average: {custom_round(average)}") # Print the average rounded to the nearest integer 
    print(f"Median: {custom_round(median)}") # Print the median rounded to the nearest integer 
    print(f"Variance: {custom_round(variance)}") # Print the variance rounded to the nearest integer 
    print(f"Standard Deviation: {custom_round(std_deviation)}") # Print the standard deviation rounded to the nearest integer 

if __name__ == "__main__": # Check if the script is being run directly 
    main()  # Call the main function 
