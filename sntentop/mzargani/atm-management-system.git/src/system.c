#include "header.h"

const char *RECORDS = "./data/records.txt";

// Function to read account data from file
int getAccountFromFile(FILE *ptr, char name[50], struct Record *r) {
    return fscanf(ptr, "%d %d %s %d %d/%d/%d %s %d %lf %s",
                  &r->id,
                  &r->userId,
                  name,
                  &r->accountNbr,
                  &r->deposit.month,
                  &r->deposit.day,
                  &r->deposit.year,
                  r->country,
                  &r->phone,
                  &r->amount,
                  r->accountType) != EOF;
}

// Function to save account data to file
void saveAccountToFile(FILE *ptr, struct User u, struct Record r) {
    fprintf(ptr, "%d %d %s %d %d/%d/%d %s %d %.2lf %s\n",
            r.id,
            u.id,
            u.name,
            r.accountNbr,
            r.deposit.month,
            r.deposit.day,
            r.deposit.year,
            r.country,
            r.phone,
            r.amount,
            r.accountType);
}

// Function to create a new account
void createNewAcc(struct User u) {
    struct Record r;
    struct Record cr;
    char userName[50];
    FILE *pf = fopen(RECORDS, "a+");

    if (pf == NULL) {
        printf("Error! Unable to open records file.\n");
        exit(1);
    }

noAccount:
    system("clear");
    printf("\t\t\t===== New record =====\n");

    printf("\nEnter today's date(mm/dd/yyyy):");
    scanf("%d/%d/%d", &r.deposit.month, &r.deposit.day, &r.deposit.year);
    printf("\nEnter the account number:");
    scanf("%d", &r.accountNbr);

    // Check if the account number already exists for this user
    while (getAccountFromFile(pf, userName, &cr)) {
        if (strcmp(userName, u.name) == 0 && cr.accountNbr == r.accountNbr) {
            printf("✖ This Account already exists for this user\n\n");
            goto noAccount;  // Restart if account exists
        }
    }

    // After verifying, save the new record
    printf("\nEnter the country:");
    scanf("%s", r.country);
    printf("\nEnter the phone number:");
    scanf("%d", &r.phone);
    printf("\nEnter amount to deposit: $");
    scanf("%lf", &r.amount);
    printf("\nChoose the type of account:\n\t-> saving\n\t-> current\n\t-> fixed01(for 1 year)\n\t-> fixed02(for 2 years)\n\t-> fixed03(for 3 years)\n\n\tEnter your choice:");
    scanf("%s", r.accountType);

    // Save the new record to file
    saveAccountToFile(pf, u, r);
    fclose(pf);
    success(u);
}


// Function to update account information
void updateAccount(struct User u) {
    struct Record r;
    char userName[50];
    int accountId, choice, found = 0;
    FILE *pf = fopen(RECORDS, "r");
    FILE *temp = fopen("./data/temp.txt", "w");

    if (pf == NULL || temp == NULL) {
        printf("Error! Unable to open records file.\n");
        exit(1);
    }

    system("clear");
    printf("\t\t====== Update Account Information ======\n\n");
    printf("Enter the account ID you want to update: ");
    scanf("%d", &accountId);

    while (getAccountFromFile(pf, userName, &r)) {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountId) {
            found = 1;
            printf("\nWhat would you like to update?\n");
            printf("1. Country\n");
            printf("2. Phone Number\n");
            printf("Enter your choice: ");
            scanf("%d", &choice);

            if (choice == 1) {
                printf("\nEnter the new country: ");
                scanf("%s", r.country);
            } else if (choice == 2) {
                printf("\nEnter the new phone number: ");
                scanf("%d", &r.phone);
            } else {
                printf("\nInvalid choice!\n");
                fclose(pf);
                fclose(temp);
                remove("./data/temp.txt");
                return;
            }
            printf("\n✔ Account updated successfully!\n");
        }
        saveAccountToFile(temp, u, r);
    }

    fclose(pf);
    fclose(temp);

    if (!found) {
        printf("\n✖ Account not found!\n");
        remove("./data/temp.txt");
        return;
    }

    remove(RECORDS);
    rename("./data/temp.txt", RECORDS);
}

// Function to check account details
void checkAccountDetails(struct User u) {
    struct Record r;
    char userName[50];
    int accountId, found = 0;

    FILE *pf = fopen(RECORDS, "r");

    system("clear");
    printf("\t\t====== Check Account Details ======\n\n");
    printf("Enter the account ID you want to view: ");
    scanf("%d", &accountId);

    while (getAccountFromFile(pf, userName, &r)) {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountId) {
            found = 1;
            printf("\nAccount Details:\n");
            printf("Account Number: %d\n", r.accountNbr);
            printf("Deposit Date: %d/%d/%d\n", r.deposit.day, r.deposit.month, r.deposit.year);
            printf("Country: %s\n", r.country);
            printf("Phone Number: %d\n", r.phone);
            printf("Balance: $%.2f\n", r.amount);
            printf("Account Type: %s\n", r.accountType);

            double interestRate = 0.0;
            if (strcmp(r.accountType, "savings") == 0) interestRate = 0.07;
            else if (strcmp(r.accountType, "fixed01") == 0) interestRate = 0.04;
            else if (strcmp(r.accountType, "fixed02") == 0) interestRate = 0.05;
            else if (strcmp(r.accountType, "fixed03") == 0) interestRate = 0.08;

            if (interestRate > 0) {
                double interest = r.amount * interestRate / 12; // Monthly interest
                printf("You will get $%.2f as interest on day %d of every month.\n", interest, r.deposit.day);
            } else if (strcmp(r.accountType, "current") == 0) {
                printf("You will not get interests because the account is of type 'current'.\n");
            }
            break;
        }
    }

    fclose(pf);
    if (!found) {
        printf("\n✖ Account not found!\n");
    }
}

// Function to show success and return to main menu or exit
void success(struct User u) {
    int option;
    printf("\n✔ Success!\n\n");
invalid:
    printf("Enter 1 to go to the main menu and 0 to exit!\n");
    scanf("%d", &option);
    system("clear");
    if (option == 1) {
        mainMenu(u);
    } else if (option == 0) {
        exit(1);
    } else {
        printf("Insert a valid operation!\n");
        goto invalid;
    }
}

// Function to remove an account
void removeAccount(struct User u) {
    struct Record r;
    char userName[50];
    int accountId, found = 0;

    FILE *pf = fopen(RECORDS, "r");
    FILE *temp = fopen("./data/temp.txt", "w");

    if (pf == NULL || temp == NULL) {
        printf("Error! Unable to open records file.\n");
        exit(1);
    }

    system("clear");
    printf("\t\t====== Remove Account ======\n\n");
    printf("Enter the account ID you want to remove: ");
    scanf("%d", &accountId);

    while (getAccountFromFile(pf, userName, &r)) {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountId) {
            found = 1;
            printf("\nAre you sure you want to delete the account with ID %d? (y/n): ", r.accountNbr);
            char confirm;
            scanf(" %c", &confirm);  // The space before %c is to consume the newline left by previous input

            if (confirm == 'y' || confirm == 'Y') {
                printf("\n✔ Account deleted successfully!\n");
                continue; // Skip writing this account to the temp file
            } else {
                printf("\n✖ Account deletion canceled!\n");
                break;
            }
        }
        // Write record to temp file if it's not the deleted account
        saveAccountToFile(temp, u, r);
    }

    fclose(pf);
    fclose(temp);

    if (!found) {
        printf("\n✖ Account not found!\n");
        remove("./data/temp.txt");
        return;
    }

    remove(RECORDS);
    rename("./data/temp.txt", RECORDS);
}

// Function to transfer ownership of an account
void transferOwnership(struct User u) {
    struct Record r;
    char userName[50];
    int accountId, newUserId;
    char newUserName[50];
    int found = 0;

    FILE *pf = fopen(RECORDS, "r+");  // Open for reading and writing
    FILE *temp = fopen("./data/temp.txt", "w");

    if (pf == NULL || temp == NULL) {
        printf("Error! Unable to open records file.\n");
        exit(1);
    }

    system("clear");
    printf("\t\t====== Transfer Account Ownership ======\n\n");
    printf("Enter the account ID you want to transfer ownership of: ");
    scanf("%d", &accountId);

    while (getAccountFromFile(pf, userName, &r)) {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountId) {
            found = 1;
            printf("\nAccount Details:\n");
            printf("Account Number: %d\n", r.accountNbr);
            printf("Current Owner: %s (ID: %d)\n", u.name, u.id);
            printf("Account Type: %s\n", r.accountType);

            printf("\nEnter the new user's ID: ");
            scanf("%d", &newUserId);
            printf("Enter the new user's name: ");
            scanf("%s", newUserName);

            // Update account owner
            r.userId = newUserId;
            strcpy(r.name, newUserName);

            printf("\n✔ Ownership transferred successfully!\n");

            saveAccountToFile(temp, u, r);
            break;
        }
    }

    fclose(pf);
    fclose(temp);

    if (!found) {
        printf("\n✖ Account not found!\n");
        remove("./data/temp.txt");
        return;
    }

    remove(RECORDS);
    rename("./data/temp.txt", RECORDS);
}

// Function to check all accounts of the user
void checkAllAccounts(struct User u) {
    struct Record r;
    char userName[100];

    FILE *pf = fopen(RECORDS, "r");

    if (pf == NULL) {
        printf("Error! Unable to open records file.\n");
        exit(1);
    }

    system("clear");
    printf("\t\t====== All accounts from user, %s =====\n\n", u.name);
    while (getAccountFromFile(pf, userName, &r)) {
        if (strcmp(userName, u.name) == 0) {
            printf("_____________________\n");
            printf("\nAccount number: %d\nDeposit Date: %d/%d/%d \nCountry: %s \nPhone number: %d \nAmount deposited: $%.2f \nType of Account: %s\n",
                   r.accountNbr,
                   r.deposit.day,
                   r.deposit.month,
                   r.deposit.year,
                   r.country,
                   r.phone,
                   r.amount,
                   r.accountType);
        }
    }

    fclose(pf);
    success(u);
}

// Function to make a transaction on an account
void makeTransaction(struct User u) {
    struct Record r;
    char userName[50];
    int accountId, transactionType;
    double amount;
    int found = 0;

    FILE *pf = fopen(RECORDS, "r+");  // Open for reading and writing
    FILE *temp = fopen("./data/temp.txt", "w");

    if (pf == NULL || temp == NULL) {
        printf("Error! Unable to open records file.\n");
        exit(1);
    }

    system("clear");
    printf("\t\t====== Make Transaction ======\n\n");
    printf("Enter the account ID for the transaction: ");
    scanf("%d", &accountId);

    while (getAccountFromFile(pf, userName, &r)) {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountId) {
            found = 1;
            printf("\nAccount Details:\n");
            printf("Account Number: %d\n", r.accountNbr);
            printf("Account Type: %s\n", r.accountType);
            printf("Current Balance: $%.2f\n", r.amount);

            // Check if transaction is allowed
            if (strcmp(r.accountType, "fixed01") == 0 || 
                strcmp(r.accountType, "fixed02") == 0 || 
                strcmp(r.accountType, "fixed03") == 0) {
                printf("\n✖ Error: Transactions are not allowed for fixed accounts.\n");
                break;
            }

            printf("\nSelect transaction type:\n");
            printf("1. Deposit\n");
            printf("2. Withdraw\n");
            printf("Enter your choice: ");
            scanf("%d", &transactionType);

            if (transactionType == 1) {
                printf("\nEnter amount to deposit: $");
                scanf("%lf", &amount);
                r.amount += amount;
                printf("\n✔ Deposit successful! New balance: $%.2f\n", r.amount);
            } else if (transactionType == 2) {
                printf("\nEnter amount to withdraw: $");
                scanf("%lf", &amount);
                if (amount > r.amount) {
                    printf("\n✖ Insufficient funds!\n");
                } else {
                    r.amount -= amount;
                    printf("\n✔ Withdrawal successful! New balance: $%.2f\n", r.amount);
                }
            } else {
                printf("\nInvalid transaction type!\n");
                break;
            }

            // Write updated account to the temp file
            saveAccountToFile(temp, u, r);
            break;
        }
    }

    fclose(pf);
    fclose(temp);

    if (!found) {
        printf("\n✖ Account not found!\n");
        remove("./data/temp.txt");
        return;
    }

    // Replace old file with updated file
    remove(RECORDS);
    rename("./data/temp.txt", RECORDS);
}
