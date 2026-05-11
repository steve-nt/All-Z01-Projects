#include "header.h"

void createNewAcc(struct User u)
{
    // store new account information
    struct Record r;
    // used for checking existing records
    struct Record cr;
    char userName[50];
    // open the file in append and read mode
    FILE *pf = fopen(RECORDS, "a+");

    if (pf == NULL)
    {
        printf("\n✖ Cannot open file!\n");
        exit(1);
    }

    // label for retrying account creation
    // the first time it is disregarded
    noAccount:
    system("clear");
    printf("\t\t\t===== New record =====\n");

    // Date input with validation
    int month, day, year;
    char dateInput[20];
    int validDate = 0;
    
    while (!validDate) {
        printf("\nEnter today's date (mm/dd/yyyy): ");
        scanf("%s", dateInput);
        
        // Try to parse the date
        if (sscanf(dateInput, "%d/%d/%d", &month, &day, &year) == 3) {
            // Basic date validation
            if (month >= 1 && month <= 12 && day >= 1 && day <= 31 && year >= 1900) {
                validDate = 1;
                r.deposit.month = month;
                r.deposit.day = day;
                r.deposit.year = year;
                clearInputBuffer();
            } else {
                printf("\n✖ Invalid date! Please use a valid date.\n");
                clearInputBuffer();
            }
        } else {
            printf("\n✖ Invalid format! Please use mm/dd/yyyy format.\n");
            clearInputBuffer();
        }
    }
    
    // Account number input with validation
    int validAccount = 0;
    while (!validAccount) {
        printf("\nEnter the account number: ");
        if (scanf("%d", &r.accountNbr) == 1) {
            validAccount = 1;
            clearInputBuffer();
        } else {
            printf("\n✖ Invalid input! Please enter a number.\n");
            clearInputBuffer();
        }
    }

    // Go to beginning of file to check for existing accounts
    rewind(pf);
    while (getAccountFromFile(pf, userName, &cr))
    {
        if (strcmp(userName, u.name) == 0 && cr.accountNbr == r.accountNbr)
        {
            printf("✖ This Account already exists for this user\n\n");
            sleep(2);
            goto noAccount;
        }
    }
    
    // Set user ID correctly from the User struct
    r.userId = u.id;
    
    // Get the highest record ID and increment
    rewind(pf);
    int maxId = -1;
    while (getAccountFromFile(pf, userName, &cr))
    {
        if (cr.id > maxId)
            maxId = cr.id;
    }
    r.id = maxId + 1;
    
    printf("\nEnter the country: ");
    scanf("%s", r.country);
    
    // Phone number input with validation
    int validPhone = 0;
    while (!validPhone) {
        printf("\nEnter the phone number: ");
        char phoneStr[20];
        scanf("%s", phoneStr);
        
        // Check if all characters are digits
        int allDigits = 1;
        for (int i = 0; phoneStr[i] != '\0'; i++) {
            if (phoneStr[i] < '0' || phoneStr[i] > '9') {
                allDigits = 0;
                break;
            }
        }
        
        if (allDigits) {
            // convert string to long long integer
            r.phone = strtoll(phoneStr, NULL, 10);
            validPhone = 1;
        } else {
            printf("\n✖ Invalid phone number! Please enter digits only.\n");
        }
    }
    clearInputBuffer();
    
    // Amount input with validation
    int validAmount = 0;
    char amountStr[50];
    while (!validAmount) {
        printf("\nEnter amount to deposit: $");
        scanf("%s", amountStr);
        
        // Check if the input is a valid number
        char *endptr;
        double value = strtod(amountStr, &endptr);
        
        if (endptr != amountStr && *endptr == '\0' && value >= 0) {
            r.amount = value;
            validAmount = 1;
        } else {
            printf("\n✖ Invalid amount! Please enter a positive number.\n");
        }
    }
    clearInputBuffer();
    
    // Account type input with validation
    int validType = 0;
    while (!validType) {
        printf("\nChoose the type of account:\n");
        printf("\t-> saving\n");
        printf("\t-> current\n");
        printf("\t-> fixed01 (for 1 year)\n");
        printf("\t-> fixed02 (for 2 years)\n");
        printf("\t-> fixed03 (for 3 years)\n");
        printf("\n\tEnter your choice: ");
        
        scanf("%s", r.accountType);
        
        if (strcmp(r.accountType, "saving") == 0 || 
            strcmp(r.accountType, "current") == 0 || 
            strcmp(r.accountType, "fixed01") == 0 || 
            strcmp(r.accountType, "fixed02") == 0 || 
            strcmp(r.accountType, "fixed03") == 0) {
            validType = 1;
        } else {
            printf("\n✖ Invalid account type! Please choose from the options listed.\n");
        }
    }
    clearInputBuffer();

    // Move to end of file for appending
    fseek(pf, 0, SEEK_END);
    
    // Write directly to the file with the correct format
    fprintf(pf, "%d %d %s %d %d/%d/%d %s %lld %.2lf %s\n",
            r.id, r.userId, u.name, r.accountNbr,
            r.deposit.month, r.deposit.day, r.deposit.year,
            r.country, r.phone, r.amount, r.accountType);

    fclose(pf);
    success(u);
}

void checkAllAccounts(struct User u)
{
    char userName[100];
    struct Record r;
    int found = 0;

    FILE *pf = fopen(RECORDS, "r");
    
    if (pf == NULL)
    {
        printf("\n✖ Cannot open file!\n");
        exit(1);
    }

    system("clear");
    printf("\t\t====== All accounts from user, %s =====\n\n", u.name);
    while (getAccountFromFile(pf, userName, &r))
    {
        if (strcmp(userName, u.name) == 0)
        {
            found = 1;
            printf("_____________________\n");
            printf("\nAccount number:%d\nDeposit Date:%d/%d/%d \ncountry:%s \nPhone number:%lld \nAmount deposited: $%.2f \nType Of Account:%s\n",
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
    
    if (!found)
    {
        printf("\n✖ No accounts found for this user!\n");
    }
    
    fclose(pf);
    success(u);
}

void checkAccount(struct User u)
{
    char userName[100];
    struct Record r;
    int accountNum, found = 0;
    double interest = 0;

    FILE *pf = fopen(RECORDS, "r");
    
    if (pf == NULL)
    {
        printf("\n✖ Cannot open file!\n");
        exit(1);
    }

    system("clear");
    printf("\t\t====== Check Account Details =====\n\n");
    
    // Validate account number input
    int validInput = 0;
    while (!validInput) {
        printf("Enter the account number: ");
        if (scanf("%d", &accountNum) == 1) {
            validInput = 1;
            clearInputBuffer();
        } else {
            printf("\n✖ Invalid input! Please enter a number.\n");
            clearInputBuffer();
        }
    }
    
    // loop line by line in pf
    while (getAccountFromFile(pf, userName, &r))
    {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountNum)
        {
            found = 1;
            printf("\n_____________________\n");
            printf("\nAccount number: %d\nDeposit Date: %d/%d/%d \nCountry: %s \nPhone number: %lld \nAmount deposited: $%.2f \nType Of Account: %s\n",
                   r.accountNbr,
                   r.deposit.day,
                   r.deposit.month,
                   r.deposit.year,
                   r.country,
                   r.phone,
                   r.amount,
                   r.accountType);
            
            // Calculate interest based on account type
            if (strcmp(r.accountType, "saving") == 0)
            {
                // Savings: interest rate 7%
                interest = r.amount * 0.07 / 12;
                printf("\nYou will get $%.2f as interest on day %d of every month.\n", interest, r.deposit.day);
            }
            else if (strcmp(r.accountType, "fixed01") == 0)
            {
                // Fixed01 (1 year): interest rate 4%
                interest = r.amount * 0.04;
                printf("\nYou will get $%.2f as interest on %d/%d/%d.\n", 
                       interest, r.deposit.day, r.deposit.month, r.deposit.year + 1);
            }
            else if (strcmp(r.accountType, "fixed02") == 0)
            {
                // Fixed02 (2 years): interest rate 5%
                interest = r.amount * 0.05 * 2;
                printf("\nYou will get $%.2f as interest on %d/%d/%d.\n", 
                       interest, r.deposit.day, r.deposit.month, r.deposit.year + 2);
            }
            else if (strcmp(r.accountType, "fixed03") == 0)
            {
                // Fixed03 (3 years): interest rate 8%
                interest = r.amount * 0.08 * 3;
                printf("\nYou will get $%.2f as interest on %d/%d/%d.\n", 
                       interest, r.deposit.day, r.deposit.month, r.deposit.year + 3);
            }
            else if (strcmp(r.accountType, "current") == 0)
            {
                printf("\nYou will not get interests because the account is of type current.\n");
            }
            break;
        }
    }
    
    if (!found)
    {
        printf("\n✖ Account not found or you don't own this account!\n");
    }
    
    fclose(pf);
    success(u);
}

void updateAccount(struct User u)
{
    FILE *pf = fopen(RECORDS, "r");
    // temporary file for updating the records
    // it will be renamed to the original file after the update
    FILE *tempFile = fopen("./data/temp.txt", "w");
    
    if (pf == NULL || tempFile == NULL)
    {
        printf("\n✖ Cannot open file!\n");
        exit(1);
    }

    system("clear");
    printf("\t\t====== Update Account Information =====\n\n");
    
    int accountNum, choice, found = 0;
    printf("Enter the account number: ");
    scanf("%d", &accountNum);
    clearInputBuffer(); 
    
    // Copy each line from the original file to the temp file
    char userName[100];
    struct Record r;
    
    while (getAccountFromFile(pf, userName, &r))
    {
        // Check if this is the record we want to update
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountNum)
        {
            found = 1;
            printf("\nChoose what you want to update:\n");
            printf("1. Phone Number\n");
            printf("2. Country\n");
            printf("Enter your choice: ");
            scanf("%d", &choice);
            clearInputBuffer();
            
            if (choice == 1)
            {
                printf("\nCurrent phone number: %lld", r.phone);
                printf("\nEnter new phone number: ");
                scanf("%lld", &r.phone);
                clearInputBuffer();
                printf("\nPhone number updated successfully!");
            }
            else if (choice == 2)
            {
                printf("\nCurrent country: %s", r.country);
                printf("\nEnter new country: ");
                scanf("%s", r.country);
                clearInputBuffer();
                printf("\nCountry updated successfully!");
            }
            else
            {
                printf("\nInvalid choice!");
                found = -1; // Mark as invalid choice
            }
        }
        
        // Create a temporary user struct with the correct username
        struct User recordOwner;
        recordOwner.id = r.userId;
        strcpy(recordOwner.name, userName);
        
        // Write the record (either updated or as is)
        saveAccountToFile(tempFile, recordOwner, r);
    }
    
    fclose(pf);
    fclose(tempFile);
    
    // Replace the original file with the updated one
    remove(RECORDS);
    rename("./data/temp.txt", RECORDS);
    
    // check if account was found and updated
    if (found <= 0)
    {
        if (found == 0) {
            printf("\n✖ Account not found or you don't own this account!\n");
        }
        
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
        // Only call success when the operation was actually successful
        success(u);
    }
}

void removeAccount(struct User u)
{
    char userName[100];
    struct Record r;
    int accountNum, choice, found = 0;
    FILE *pf = fopen(RECORDS, "r");
    FILE *tempFile = fopen("./data/temp.txt", "w");
    
    if (pf == NULL || tempFile == NULL)
    {
        printf("\n✖ Cannot open file!\n");
        exit(1);
    }

    system("clear");
    printf("\t\t====== Remove Account =====\n\n");
    
    // Validate account number input
    int validInput = 0;
    while (!validInput) {
        printf("Enter the account number you want to remove: ");
        if (scanf("%d", &accountNum) == 1) {
            validInput = 1;
            clearInputBuffer();
        } else {
            printf("\n✖ Invalid input! Please enter a number.\n");
            clearInputBuffer();
        }
    }
    
    // First scan to check if account exists
    while (getAccountFromFile(pf, userName, &r))
    {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountNum)
        {
            found = 1;
            break;
        }
    }
    
    if (!found)
    {
        printf("\n✖ Account not found or you don't own this account!\n");
        fclose(pf);
        fclose(tempFile);
        remove("./data/temp.txt");
        
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
        return;
    }
    
    // Display account details before removal
    printf("\nAccount Details:\n");
    printf("Account number: %d\nDeposit Date: %d/%d/%d\nCountry: %s\nPhone number: %lld\nAmount: $%.2f\nType: %s\n", 
           r.accountNbr, r.deposit.day, r.deposit.month, r.deposit.year, 
           r.country, r.phone, r.amount, r.accountType);
    
    // Confirm removal
    int validChoice = 0;
    while (!validChoice) {
        printf("\nAre you sure you want to remove this account?\n");
        printf("1. Yes\n");
        printf("2. No\n");
        printf("Enter your choice (1-2): ");
        
        if (scanf("%d", &choice) == 1 && (choice == 1 || choice == 2)) {
            validChoice = 1;
            clearInputBuffer();
        } else {
            printf("\n✖ Invalid choice! Please enter 1 or 2.\n");
            clearInputBuffer();
        }
    }
    
    if (choice == 2) {
        printf("\nAccount removal cancelled.\n");
        fclose(pf);
        fclose(tempFile);
        remove("./data/temp.txt");
        
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
        return;
    }
    
    // Reset file position and copy all records except the one to be deleted
    rewind(pf);
    int removedAccount = 0;
    
    while (getAccountFromFile(pf, userName, &r))
    {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountNum)
        {
            // Skip this record (don't write to temp file)
            removedAccount = 1;
            continue;
        }
        else
        {
            // Keep other accounts as they are
            struct User currentOwner;
            currentOwner.id = r.userId;
            strcpy(currentOwner.name, userName);
            saveAccountToFile(tempFile, currentOwner, r);
        }
    }
    
    fclose(pf);
    fclose(tempFile);
    
    // Replace the original file with the updated one
    remove(RECORDS);
    rename("./data/temp.txt", RECORDS);
    
    if (removedAccount)
    {
        printf("\n✅ Account successfully removed!\n");
        // Only call success if an account was actually removed
        success(u);
    }
    else
    {
        printf("\n✖ Failed to remove the account.\n");
        
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
}

void transferOwnership(struct User u)
{
    char userName[100], newOwner[50];
    struct Record r;
    int accountNum, found = 0, targetUserExists = 0;
    FILE *pf = fopen(RECORDS, "r");
    FILE *tempFile = fopen("./data/temp.txt", "w");
    FILE *userFile = fopen(USERS, "r");
    
    if (pf == NULL || tempFile == NULL || userFile == NULL)
    {
        printf("\n✖ Cannot open file!\n");
        exit(1);
    }

    system("clear");
    printf("\t\t====== Transfer Account Ownership =====\n\n");
    
    // Validate account number input
    int validInput = 0;
    while (!validInput) {
        printf("Enter the account number you want to transfer: ");
        if (scanf("%d", &accountNum) == 1) {
            validInput = 1;
            clearInputBuffer();
        } else {
            printf("\n✖ Invalid input! Please enter a number.\n");
            clearInputBuffer();
        }
    }
    
    // First check if the account exists and belongs to the current user
    while (getAccountFromFile(pf, userName, &r))
    {
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountNum)
        {
            found = 1;
            break;
        }
    }
    
    if (!found)
    {
        printf("\n✖ Account not found or you don't own this account!\n");
        fclose(pf);
        fclose(tempFile);
        fclose(userFile);
        remove("./data/temp.txt");
        
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
        return;
    }
    
    // Display account details before transfer
    printf("\nAccount Details:\n");
    printf("Account number: %d\nDeposit Date: %d/%d/%d\nCountry: %s\nPhone number: %lld\nAmount: $%.2f\nType: %s\n", 
           r.accountNbr, r.deposit.day, r.deposit.month, r.deposit.year, 
           r.country, r.phone, r.amount, r.accountType);
    
    printf("\nEnter the username to transfer this account to: ");
    scanf("%s", newOwner);
    clearInputBuffer();
    
    // Check if target user exists and get their user ID
    struct User targetUser;
    rewind(userFile);
    while (fscanf(userFile, "%d %s %s", &targetUser.id, targetUser.name, targetUser.password) != EOF)
    {
        if (strcmp(targetUser.name, newOwner) == 0)
        {
            targetUserExists = 1;
            break;
        }
    }
    
    if (!targetUserExists)
    {
        printf("\n✖ Target user '%s' does not exist!\n", newOwner);
        fclose(pf);
        fclose(tempFile);
        fclose(userFile);
        remove("./data/temp.txt");
        
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
        return;
    }
    
    // Confirm the transfer
    int choice;
    int validChoice = 0;
    while (!validChoice) {
        printf("\nAre you sure you want to transfer this account to %s?\n", newOwner);
        printf("1. Yes\n");
        printf("2. No\n");
        printf("Enter your choice (1-2): ");
        
        if (scanf("%d", &choice) == 1 && (choice == 1 || choice == 2)) {
            validChoice = 1;
            clearInputBuffer();
        } else {
            printf("\n✖ Invalid choice! Please enter 1 or 2.\n");
            clearInputBuffer();
        }
    }
    
    if (choice == 2) {
        printf("\nAccount transfer cancelled.\n");
        fclose(pf);
        fclose(tempFile);
        fclose(userFile);
        remove("./data/temp.txt");
        
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
        return;
    }
    
    // Reset file position and perform the transfer
    rewind(pf);
    int transferCompleted = 0;
    
    while (getAccountFromFile(pf, userName, &r))
    {
        // Create a temporary owner struct based on whether this is the account to transfer
        struct User recordOwner;
        
        if (strcmp(userName, u.name) == 0 && r.accountNbr == accountNum)
        {
            // Update the user ID and name to the new owner
            r.userId = targetUser.id;
            // Use the new owner's info
            recordOwner = targetUser;
            transferCompleted = 1;
        }
        else
        {
            // Keep other accounts as they are
            recordOwner.id = r.userId;
            strcpy(recordOwner.name, userName);
        }
        
        // Save the record (either transferred or as is)
        saveAccountToFile(tempFile, recordOwner, r);
    }
    
    fclose(pf);
    fclose(tempFile);
    fclose(userFile);
    
    // Replace the original file with the updated one
    remove(RECORDS);
    rename("./data/temp.txt", RECORDS);
    
    if (transferCompleted)
    {
        printf("\n✅ Account successfully transferred to %s!\n", newOwner);
        // Only call success if the transfer was completed
        success(u);
    }
    else
    {
        printf("\n✖ Failed to transfer the account.\n");
        
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
}