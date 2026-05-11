#ifndef HEADER_H
#define HEADER_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

/* Constants */
extern const char *RECORDS;
extern const char *USERS;

/* Structures */
struct Date
{
    int month, day, year;
};

// all fields for each record of an account
struct Record
{
    int id;
    int userId;
    char name[100];
    char country[100];
    long long phone;  // Changed from int to long long to handle larger numbers
    char accountType[10];
    int accountNbr;
    double amount;
    struct Date deposit;
};

struct User
{
    int id;
    char name[50];
    char password[50];
};

/* Helper functions */
void clearInputBuffer(); // Declared once here

/* Forward declarations */
void mainMenu(struct User u);
void initMenu(struct User *u);
void success(struct User u);

/* Authentication functions (auth.c) */
void loginMenu(char a[50], char pass[50]);
void registerMenu(struct User *u);
const char *getPassword(struct User u);

/* Account functions (account_utils.c) */
void createNewAcc(struct User u);
void checkAllAccounts(struct User u);
void checkAccount(struct User u);
void updateAccount(struct User u);
void removeAccount(struct User u);
void transferOwnership(struct User u);

/* Transaction functions (transactions.c) */
void makeTransaction(struct User u);

/* Common utility functions (system.c) */
int getAccountFromFile(FILE *ptr, char name[50], struct Record *r);
void saveAccountToFile(FILE *ptr, struct User u, struct Record r);


#endif /* HEADER_H */