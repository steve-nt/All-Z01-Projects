#include <termios.h>
#include "header.h"

extern const char *USERS;

void loginMenu(char a[50], char pass[50])
{
    struct termios oflags, nflags;

    system("clear");
    printf("\n\n\n\t\t\t\t   Bank Management System\n\t\t\t\t\t User Login:");
    scanf("%s", a);
    clearInputBuffer(); // Clear the buffer after reading username

    // disabling echo
    // save current terminal settings to oflags
    tcgetattr(fileno(stdin), &oflags);
    // set new terminal settings to nflags
    nflags = oflags;
    nflags.c_lflag &= ~ECHO; // disable echo
    nflags.c_lflag |= ECHONL; // enable newline echo

    // apply new settings
    // TCSANOW means change settings immediately
    if (tcsetattr(fileno(stdin), TCSANOW, &nflags) != 0)
    {
        perror("tcsetattr");
        exit(1);
    }
    printf("\n\n\n\n\n\t\t\t\tEnter the password to login:");
    scanf("%s", pass);
    clearInputBuffer(); // Clear the buffer after reading password

    // restore terminal
    if (tcsetattr(fileno(stdin), TCSANOW, &oflags) != 0)
    {
        perror("tcsetattr");
        exit(1);
    }
}

const char *getPassword(struct User u)
{
    FILE *fp;
    static struct User userChecker;

    if ((fp = fopen(USERS, "r")) == NULL)
    {
        printf("Error! opening file");
        exit(1);
    }

    while (fscanf(fp, "%d %s %s", &userChecker.id, userChecker.name, userChecker.password) != EOF)
    {
        if (strcmp(userChecker.name, u.name) == 0)
        {
            fclose(fp);
            return userChecker.password;
        }
    }

    fclose(fp);
    return "no user found";
}

void registerMenu(struct User *u)
{
    FILE *fp;
    struct User userChecker;
    int userCount = 0;
    
    system("clear");
    printf("\n\n\n\t\t\t\t   Bank Management System\n\t\t\t\t\t User Registration\n");
    
    // Get username
    printf("\n\t\tEnter Username: ");
    scanf("%s", u->name);
    clearInputBuffer(); // Clear the buffer after reading username
    
    // Check if username already exists
    if ((fp = fopen(USERS, "r")) == NULL)
    {
        printf("Error! opening file");
        exit(1);
    }
    
    // Count users and check for existing username
    while (fscanf(fp, "%d %s %s", &userChecker.id, userChecker.name, userChecker.password) != EOF)
    {
        userCount++;
        if (strcmp(userChecker.name, u->name) == 0)
        {
            printf("\n\t\tUsername already exists! Please choose another name.\n");
            fclose(fp);
            sleep(3);
            initMenu(u);
            return;
        }
    }
    
    fclose(fp);
    
    // Set user ID
    u->id = userCount;
    
    struct termios oflags, nflags;
    
    // disabling echo for password
    tcgetattr(fileno(stdin), &oflags);
    nflags = oflags;
    nflags.c_lflag &= ~ECHO;
    nflags.c_lflag |= ECHONL;
    
    if (tcsetattr(fileno(stdin), TCSANOW, &nflags) != 0)
    {
        perror("tcsetattr");
        exit(1);
    }
    
    printf("\n\t\tEnter Password: ");
    scanf("%s", u->password);
    clearInputBuffer(); // Clear the buffer after reading password
    
    // restore terminal
    if (tcsetattr(fileno(stdin), TCSANOW, &oflags) != 0)
    {
        perror("tcsetattr");
        exit(1);
    }
    
    // Save user to file
    // open in append mode to add new user
    if ((fp = fopen(USERS, "a")) == NULL)
    {
        printf("Error! opening file");
        exit(1);
    }
    
    fprintf(fp, "%d %s %s\n", u->id, u->name, u->password);
    fclose(fp);
    
    printf("\n\t\tRegistration successful!\n");
    printf("\n\t\tReturning to login menu...\n");
    sleep(2);
    
    // Clear the user struct to prevent auto-login
    // set all bytes to to 0 to ensure the user must explicitly login again
    // after registration
    u->id = 0;
    memset(u->name, 0, sizeof(u->name));
    memset(u->password, 0, sizeof(u->password));
}