#include "header.h"

// Function to make transactions
void makeTransaction(struct User u)
{
    char userName[100];
    struct Record r;
    int accountNum, choice, found = 0;
    int transactionSuccess = 0; // Flag to track if a transaction was successfully processed
    double amount;
    FILE *pf = fopen(RECORDS, "r");
    FILE *tempFile = fopen("./data/temp.txt", "w");
    
    if (pf == NULL || tempFile == NULL)
    {
        printf("\n✖ Cannot open file!\n");
        exit(1);
    }

    system("clear");
    printf("\t\t====== Make Transaction =====\n\n");
    
    // Validate account number input
    int validInput = 0;
    while (!validInput) {
        printf("Enter the account number: ");
        if (scanf("%d", &accountNum) == 1) {
            validInput = 1;
        } else {
            printf("\n✖ Invalid input! Please enter a number.\n");
            // Clear input buffer
            clearInputBuffer();
        }
    }
    clearInputBuffer(); // Clear input buffer after successful input
    
    // Process all records
    while (getAccountFromFile(pf, userName, &r))
    {
        // Create a temporary user struct with the original account owner's name
        struct User recordOwner;
        recordOwner.id = r.userId;
        strcpy(recordOwner.name, userName);
        
        // Check if this is the target account for transaction
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountNum)
        {
            found = 1;
            
            // Check if account type allows transactions
            if (strcmp(r.accountType, "fixed01") == 0 || 
                strcmp(r.accountType, "fixed02") == 0 || 
                strcmp(r.accountType, "fixed03") == 0)
            {
                printf("\n✖ Error: Transactions are not allowed for fixed deposit accounts!\n");
                printf("\nAccount type: %s does not support withdrawals or deposits.\n", r.accountType);
                
                // Write the unchanged record with the original owner's name
                saveAccountToFile(tempFile, recordOwner, r);
                continue;
            }
            
            printf("\nCurrent Balance: $%.2f\n", r.amount);
            
            // Transaction type selection with validation
            // repeatedly ask for a valid choice until one is given
            do {
                printf("\nChoose transaction type:\n");
                printf("1. Deposit\n");
                printf("2. Withdraw\n");
                printf("Enter your choice (1-2): ");
                
                if (scanf("%d", &choice) != 1 || (choice != 1 && choice != 2)) {
                    printf("\n✖ Invalid choice! Please enter 1 or 2.\n");
                    // Clear input buffer
                    clearInputBuffer();
                    choice = 0;  // Invalid choice
                }
            } while (choice != 1 && choice != 2);
            clearInputBuffer(); // Clear input buffer after successful input
            
            if (choice == 1) // Deposit
            {
                // Amount validation
                int validAmount = 0;
                while (!validAmount) {
                    printf("\nEnter amount to deposit: $");
                    char amountStr[50];
                    scanf("%s", amountStr);
                    clearInputBuffer(); // Clear input buffer
                    
                    // Check if the input is a valid number
                    char *endptr;
                    // Convert string to double
                    // endptr will point to the first invalid character
                    amount = strtod(amountStr, &endptr);
                    
                    if (endptr != amountStr && *endptr == '\0' && amount > 0) {
                        validAmount = 1;
                    } else {
                        printf("\n✖ Invalid amount! Please enter a positive number.\n");
                    }
                }
                
                r.amount += amount;
                printf("\n✅ Deposit successful! New Balance: $%.2f\n", r.amount);
                transactionSuccess = 1; // Mark transaction as successful
            }
            else if (choice == 2) // Withdraw
            {
                // Amount validation
                int validAmount = 0;
                while (!validAmount) {
                    printf("\nEnter amount to withdraw: $");
                    char amountStr[50];
                    scanf("%s", amountStr);
                    clearInputBuffer(); // Clear input buffer
                    
                    // Check if the input is a valid number
                    char *endptr;
                    amount = strtod(amountStr, &endptr);
                    
                    if (endptr != amountStr && *endptr == '\0' && amount > 0) {
                        // Additional check for sufficient funds
                        if (amount <= r.amount) {
                            validAmount = 1;
                        } else {
                            printf("\n✖ Insufficient funds! Maximum withdrawal amount: $%.2f\n", r.amount);
                        }
                    } else {
                        printf("\n✖ Invalid amount! Please enter a positive number.\n");
                    }
                }
                
                r.amount -= amount;
                printf("\n✅ Withdrawal successful! New Balance: $%.2f\n", r.amount);
                transactionSuccess = 1; // Mark transaction as successful
            }
        }
        
        // Save the record (either modified or as is) with the original owner's name
        saveAccountToFile(tempFile, recordOwner, r);
    }
    
    fclose(pf);
    fclose(tempFile);
    
    // Replace the original file with the updated one
    remove(RECORDS);
    rename("./data/temp.txt", RECORDS);
    
    if (!found)
    {
        printf("\n✖ Account not found or you don't own this account!\n");
        
        int option;
        printf("\nEnter 1 to go to the main menu and 0 to exit: ");
        scanf("%d", &option);
        clearInputBuffer();
        
        if (option == 1) {
            system("clear");
            mainMenu(u);
        } else {
            system("clear");
            exit(1);
        }
    }
    else if (!transactionSuccess)
    {
        // Fixed deposit or other error case
        printf("\n✖ Transaction was not completed due to restrictions or errors.\n");
        
        int option;
        printf("\nEnter 1 to go to the main menu and 0 to exit: ");
        scanf("%d", &option);
        clearInputBuffer();
        
        if (option == 1) {
            system("clear");
            mainMenu(u);
        } else {
            system("clear");
            exit(1);
        }
    }
    else
    {
        // Only call success when transaction was actually successful
        success(u);
    }
}