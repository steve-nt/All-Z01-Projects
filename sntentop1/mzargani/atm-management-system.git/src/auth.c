#include <termios.h>
#include "header.h"

char *USERS = "./data/users.txt";

void loginMenu(char a[50], char pass[50])
{
    struct termios oflags, nflags;

    system("clear");
    printf("\n\n\n\t\t\t\t   Bank Management System\n\t\t\t\t\t User Login:");
    scanf("%s", a);

    // disabling echo
    tcgetattr(fileno(stdin), &oflags);
    nflags = oflags;
    nflags.c_lflag &= ~ECHO;
    nflags.c_lflag |= ECHONL;

    if (tcsetattr(fileno(stdin), TCSANOW, &nflags) != 0)
    {
        perror("tcsetattr");
        return exit(1);
    }
    printf("\n\n\n\n\n\t\t\t\tEnter the password to login:");
    scanf("%s", pass);

    // restore terminal
    if (tcsetattr(fileno(stdin), TCSANOW, &oflags) != 0)
    {
        perror("tcsetattr");
        return exit(1);
    }
};

const char *getPassword(struct User u)
{
    FILE *fp;
    struct User userChecker;
    static char password[100];  // Use a static buffer to store the password

    if ((fp = fopen("./data/users.txt", "r")) == NULL)
    {
        printf("Error! opening file");
        exit(1);
    }

    while (fscanf(fp, "%d %s %s", &userChecker.id, userChecker.name, userChecker.password) != EOF)
    {
        if (strcmp(userChecker.name, u.name) == 0)
        {
            fclose(fp);
            strcpy(password, userChecker.password);  // Copy the password to a static buffer
            return password;  // Return the pointer to the static buffer
        }
    }

    fclose(fp);
    return "no user found";  // Return a default message if the user is not found
}

void registerMenu(char a[50], char pass[50]) {
    FILE *fp;
    struct User existingUser;
    int maxId = -1;

    // Open the file in read-write mode
    if ((fp = fopen(USERS, "a+")) == NULL) {
        printf("Error! Unable to open users file.\n");
        exit(1);
    }

    printf("\n\n\n\t\t\t\t   Bank Management System\n\t\t\t\t\t User Registration:");
    printf("\nEnter username: ");
    scanf("%s", a);

    // Check for uniqueness of username
    while (fscanf(fp, "%d %s %s", &existingUser.id, existingUser.name, existingUser.password) != EOF) {
        if (strcmp(existingUser.name, a) == 0) {
            printf("\nError: Username already exists. Please try another one.\n");
            fclose(fp);
            return;
        }
        // Track the highest user id
        if (existingUser.id > maxId) {
            maxId = existingUser.id;
        }
    }

    // Assign the new user id as the next highest id + 1
    int userId = maxId + 1;

    // Ask for password and save the new user
    printf("Enter password: ");
    scanf("%s", pass);

    fprintf(fp, "%d %s %s\n", userId, a, pass);

    printf("\nRegistration successful! Your user ID is %d.\n", userId);

    fclose(fp);
}
