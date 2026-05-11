// header.h

#ifndef HEADER_H
#define HEADER_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

struct Date {
    int month, day, year;
};

struct Record {
    int id;
    int userId;
    char name[100];
    char country[100];
    int phone;
    char accountType[10];
    int accountNbr;
    double amount;
    struct Date deposit;
    struct Date withdraw;
};

struct User {
    int id;
    char name[50];
    char password[50];
};

// Authentication functions
void loginMenu(char a[50], char pass[50]);
void registerMenu(char a[50], char pass[50]);
const char *getPassword(struct User u);

// System functions
void createNewAcc(struct User u);
void updateAccount(struct User u);
void checkAccountDetails(struct User u);
void checkAllAccounts(struct User u);
void makeTransaction(struct User u);
void removeAccount(struct User u);
void transferOwnership(struct User u);
void mainMenu(struct User u);
void initMenu(struct User *u);

// Function declaration for success
void success(struct User u);  // Add this line

#endif // HEADER_H
