# ATM Management System

## Overview
This ATM Management System is a command-line application that simulates the functionality of a banking system. It allows users to create accounts, manage their finances, and perform basic banking operations. The application uses file-based storage to maintain user data and account records.

## Features
- **User Authentication:** Secure login and registration
- **Account Management:** Create, update, and remove accounts
- **Transaction Processing:** Deposit and withdraw funds
- **Account Inquiries:** Check account details and list all owned accounts
- **Interest Calculation:** View potential interest earnings based on account type
- **Ownership Transfer:** Transfer account ownership between users

## Account Types
The system supports multiple account types, each with its own interest rate:
- **Saving:** 7% annual interest rate (calculated monthly)
- **Current:** No interest
- **Fixed01:** 4% interest for 1-year term
- **Fixed02:** 5% interest per year for 2-year term (10% total)
- **Fixed03:** 8% interest per year for 3-year term (24% total)

## Technology Stack
- **Language:** C
- **Storage:** Text file-based database
- **Build System:** Make

## Project Structure
```
├── Makefile              # Build configuration
├── src/
│   ├── main.c            # Main program entry point
│   ├── header.h          # Common declarations and structures
│   ├── auth.c            # Authentication functions
│   ├── system.c          # Core utility functions
│   ├── account_utils.c   # Account management functions
│   ├── transactions.c    # Transaction processing
└── data/
    ├── users.txt         # User credentials storage
    └── records.txt       # Account records storage
```

## File Format
The system uses two text files to store data:

### users.txt
Stores user credentials in the format:
```
[id] [name] [password]
```
Example:
```
0 Alice 1234password
1 Michel password1234
```

### records.txt
Stores account records in the format:
```
[id] [user_id] [user_name] [account_id] [mm/dd/yyyy] [country] [phone] [balance] [account_type]
```
Example:
```
0 0 Alice 0 10/02/2020 german 986134231 11090830.00 current
1 1 Michel 2 10/10/2021 portugal 914134431 1920.42 savings
```

## Installation and Setup

### Prerequisites
- GCC compiler
- Make

### Building the Project
1. Clone the repository:
   ```
   git clone https://platform.zone01.gr/git/ychaniot/atm-management-system.git
   cd atm-management-system
   ```

2. Build the project:
   ```
   make
   ```

3. Run the application:
   ```
   ./atm
   ```

### First-time Setup
On first run, ensure the data directory exists:
```
mkdir -p data
touch data/users.txt data/records.txt
```

## Usage Guide

### Registration and Login
1. Start the application
2. Select option 2 to register a new user
3. Enter a unique username and password
4. Return to the main menu and select option 1 to login
5. Enter your credentials

### Creating a New Account
1. Log in to the system
2. Select option 1 from the main menu
3. Enter the requested information:
   - Creation date
   - Account number
   - Country
   - Phone number
   - Initial deposit amount
   - Account type

### Checking Account Details
1. Select option 3 from the main menu
2. Enter your account number
3. The system will display account details and potential interest earnings

### Making Transactions
1. Select option 5 from the main menu
2. Enter your account number
3. Choose the transaction type (deposit or withdraw)
4. Enter the amount
5. The system will update your balance accordingly

Note: Fixed deposit accounts do not allow transactions.

### Transferring Account Ownership
1. Select option 7 from the main menu
2. Enter the account number to transfer
3. Enter the username of the recipient
4. Confirm the transfer

## Security Features
- Passwords are hidden during entry
- Input validation for all user inputs
- Account ownership verification for all operations

## Future Enhancements
- Password encryption
- Transaction history logging
- Improved terminal UI
- Database integration

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is licensed under the MIT License - see the LICENSE file for details.

---

© 2025 ATM Management System