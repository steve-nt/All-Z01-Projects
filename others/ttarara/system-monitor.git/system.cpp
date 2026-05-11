#include <fstream>      // Για άνοιγμα και ανάγνωση αρχείων κειμένου
#include <string>       // Για χρήση std::string
#include <cstdlib>      // Για τη χρήση της getenv (ανάγνωση μεταβλητών περιβάλλοντος)
#include <unistd.h>     // Για την gethostname (λήψη hostname από το σύστημα)
#include "header.h"  
#include "cpu.h"
   // (Προαιρετικό – αν έχεις δηλώσεις για τις συναρτήσεις)


// 🔎 Συνάρτηση για λήψη του ονόματος του λειτουργικού συστήματος
std::string getOSName() {
    std::ifstream file("/etc/os-release");  // Άνοιγμα αρχείου συστήματος για ανάγνωση
    std::string line;

    // Διαβάζουμε το αρχείο γραμμή-γραμμή
    while (std::getline(file, line)) {
        // Αν η γραμμή ξεκινάει με "PRETTY_NAME="
        if (line.rfind("PRETTY_NAME=", 0) == 0) {
            // Αφαιρούμε το "PRETTY_NAME=" και τα εισαγωγικά από την τιμή
            std::string value = line.substr(13, line.length() - 14);
            return value; // Επιστρέφουμε το όνομα του OS
        }
    }

    return "Unknown OS";  // Σε περίπτωση που δεν βρεθεί το PRETTY_NAME
}


// 👤 Συνάρτηση για λήψη του ονόματος του συνδεδεμένου χρήστη
std::string getUserName() {
    const char* user = std::getenv("USER");  // Ανάγνωση της μεταβλητής περιβάλλοντος USER
    return user ? user : "Unknown User";     // Αν υπάρχει, επιστρέφεται· αλλιώς "Unknown User"
}


// 💻 Συνάρτηση για λήψη του hostname του υπολογιστή
std::string getHostName() {
    char hostname[1024];                     // Πίνακας χαρακτήρων για να χωρέσει το hostname
    gethostname(hostname, 1024);             // Συστημική κλήση για λήψη του hostname
    return std::string(hostname);            // Μετατροπή σε std::string και επιστροφή
}

// 🧠 Συνάρτηση για λήψη του μοντέλου του CPU από το /proc/cpuinfo
//std::string getCPUModel() {
//    std::ifstream file("/proc/cpuinfo");
//    std::string line;
//
//    while (std::getline(file, line)) {
//        // Αν η γραμμή ξεκινάει με "model name"
//        if (line.rfind("model name", 0) == 0) {
//            size_t colon = line.find(':');
//            if (colon != std::string::npos) {
//                std::string model = line.substr(colon + 1);
//                // Αφαίρεση κενών στην αρχή
//                while (!model.empty() && model.front() == ' ')
//                    model.erase(0, 1);
//                return model;
//            }
//        }
//    }
//
//    return "Unknown CPU";
//}

