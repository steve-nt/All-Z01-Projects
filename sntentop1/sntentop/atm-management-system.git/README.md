# ATM Management System

## Objective

This project demonstrates programming logic and the ability to adapt to new programming languages. The application is written in C and implements an ATM management system with user authentication and account management features.

The project involves both optimizing existing code and implementing new features for a complete ATM management system.

## Project Overview

The ATM Management System allows users to:

- Login and register for accounts
- Create new bank accounts
- Check account details and calculate interest
- Update account information
- Remove accounts
- View all owned accounts
- Make transactions (deposits and withdrawals)
- Transfer account ownership to other users

## File System Structure

```
atm-system
│
├── data/
│   ├── records.txt          (Account information storage)
│   └── users.txt            (User credentials storage)
├── Makefile
└── src/
    ├── auth.c               (Authentication logic)
    ├── header.h             (Header definitions)
    ├── main.c               (Main program entry)
    └── system.c             (Core system functions)
```

## Data Format

### users.txt
Format: `id name password`

Example:
```
0 Alice 1234password
1 Michel password1234
```

### records.txt
Format: `id user_id user_name account_id date_of_creation country phone_number balance account_type`

Example:
```
0 0 Alice 0 10/02/2020 german 986134231 11090830.00 current
1 1 Michel 2 10/10/2021 portugal 914134431 1920.42 savings
2 0 Alice 1 10/10/2000 finland 986134231 1234.21 savings
```

## Implemented Features

### 1. User Registration
- Register new users with unique names
- Prevent duplicate usernames
- Save user credentials to users.txt

### 2. Update Account Information
- Allow users to update their account details
- Users select account ID and choose which field to modify
- Permitted fields for update: country and phone number
- Changes must be saved to records.txt

### 3. Check Account Details
- Allow users to view individual account information
- Users input the account ID to view
- Display calculated interest based on account type:
  - Savings: 7% interest rate
  - Fixed01 (1 year): 4% interest rate
  - Fixed02 (2 year): 5% interest rate
  - Fixed03 (3 year): 8% interest rate
  - Current: No interest accrual
- Interest is calculated and displayed as monthly gain from the account creation date

### 4. Make Transactions
- Allow users to deposit or withdraw money
- Block transactions on fixed-term accounts (fixed01, fixed02, fixed03)
- Display error when attempting transactions on locked accounts
- Prevent withdrawals exceeding available balance
- Update records.txt with all transaction changes

### 5. Remove Account
- Allow users to delete their own accounts
- Remove account from records.txt
- Display confirmation of deletion

### 6. Transfer Account Ownership
- Allow users to transfer account ownership to another user
- Select account and target user
- Update user_id in records.txt
- Save changes to data storage

## Building and Running

### Compile the Project
```bash
make
```

### Run the Application
```bash
./atm
```

### Clean Build Artifacts
```bash
make clean
```

## Usage

1. Start the application
2. Choose to login or register a new account
3. After authentication, access the main menu with available operations
4. Follow prompts for each operation

## Bonus Features

The following features are optional enhancements:

### Process Communication and Instant Notifications
Implement instant account transfer notifications between multiple terminals using pipes and child processes. When one user transfers an account to another, the receiving user is immediately notified even if logged in separately.

### Enhanced User Interface
Improve the terminal user interface (TUI) with better formatting and user experience.

### Password Encryption
Encrypt passwords before storing them in users.txt for enhanced security.

### Custom Makefile
Create and optimize the Makefile for building the project.

### Database Integration
Integrate a relational database (SQLite recommended) as an alternative to text file storage.

### Additional Features
Implement extra features beyond the core requirements or optimize existing code for better performance and maintainability.

## Functional Testing Checklist

- User registration with unique name validation
- Prevention of duplicate user registration
- Successful user login
- Account creation with multiple accounts per user
- Account update with error handling for non-existent accounts
- Update phone number and country fields
- Interest calculation for different account types
- Transaction restrictions on fixed-term accounts
- Deposit and withdrawal operations with balance validation
- Account removal and deletion verification
- Account ownership transfer between users
- Data persistence across application restarts

## Technical Notes

- The application uses C language for implementation
- All data is stored in text files within the data/ directory
- The system uses file I/O for persistent storage
- User sessions are managed during runtime
- Interest calculations consider account creation date

## Requirements

- GCC compiler or compatible C compiler
- Make utility
- Standard C library

## Notes

- All changes to accounts must be saved to the respective data files
- User names must be unique in the system
- Account IDs are unique identifiers for accounts within the system
- Fixed-term accounts cannot be used for transactions
