#include "header.h"
#include <unistd.h>

// Add the external declaration
extern const char *USERS;

void mainMenu(struct User u)
{
    int option;
    system("clear");
    printf("\n\n\t\t======= ATM =======\n\n");
    printf("\n\t\t-->> Feel free to choose one of the options below <<--\n");
    printf("\n\t\t[1]- Create a new account\n");
    printf("\n\t\t[2]- Update account information\n");
    printf("\n\t\t[3]- Check accounts\n");
    printf("\n\t\t[4]- Check list of owned account\n");
    printf("\n\t\t[5]- Make Transaction\n");
    printf("\n\t\t[6]- Remove existing account\n");
    printf("\n\t\t[7]- Transfer ownership\n");
    printf("\n\t\t[8]- Logout\n");
    printf("\n\t\t[9]- Exit\n");
    printf("\n\t\tEnter your choice: ");
    scanf("%d", &option);
    clearInputBuffer(); // Clear the buffer after reading option

    switch (option)
    {
    case 1:
        createNewAcc(u);
        break;
    case 2:
        updateAccount(u);
        break;
    case 3:
        checkAccount(u);
        break;
    case 4:
        checkAllAccounts(u);
        break;
    case 5:
        makeTransaction(u);
        break;
    case 6:
        removeAccount(u);
        break;
    case 7:
        transferOwnership(u);
        break;
    case 8:
        printf("\n\t\t✅ Logged out successfully!\n");
        sleep(2);
        // Create a new User struct to clear current user data
        struct User newUser = {0};
        initMenu(&newUser);
        break;
    case 9:
        exit(1);
        break;
    default:
        printf("Invalid operation!\n");
        sleep(2);
        mainMenu(u);
    }
}

void initMenu(struct User *u)
{
    int option;
    system("clear");
    printf("\n\n\t\t======= ATM =======\n");
    printf("\n\t\t-->> Feel free to login / register :\n");
    printf("\n\t\t[1]- login\n");
    printf("\n\t\t[2]- register\n");
    printf("\n\t\t[3]- exit\n");
    printf("\n\t\tEnter your choice: ");
    
    scanf("%d", &option);
    clearInputBuffer(); // Clear the buffer after reading option
    
    switch (option)
    {
    case 1:
        // "pass by reference" to the username and password
        loginMenu(u->name, u->password);
        const char* pwd = getPassword(*u);
        if (strcmp(u->password, pwd) == 0)
        {
            // Retrieve user ID from file
            FILE *userFile = fopen(USERS, "r");
            if (userFile != NULL) {
                int id;
                char name[50], password[50];
                while (fscanf(userFile, "%d %s %s", &id, name, password) != EOF) {
                    if (strcmp(name, u->name) == 0) {
                        u->id = id;
                        break;
                    }
                }
                fclose(userFile);
            }
            
            system("clear");
            printf("\n\n\n\t\t\t\t   Bank Management System\n\n");
            printf("\n\n\t\t✅ Login Successful! Welcome, %s !\n\n", u->name);
            printf("\n\t\tPress Enter to continue...");
            getchar(); // Wait for user to press Enter
            mainMenu(*u);
        }
        else
        {
            printf("\nWrong password or Username\n");
            sleep(2);
            initMenu(u);
            return;
        }
        break;
    case 2:
        registerMenu(u);
        // restart the init menu to allow login
        initMenu(u);
        return;
        break;
    case 3:
        exit(1);
        break;
    default:
        printf("Insert a valid operation!\n");
        sleep(2);
        initMenu(u);
        return;
    }
}

int main()
{
    struct User u;
    
    initMenu(&u);
    return 0;
}