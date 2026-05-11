
# 🎨 **ASCII Art Stylize** 🌟

This README provides a comprehensive guide to the ASCII Art Stylize application, a project focused on enhancing the user experience of ASCII Art Web through better design and interactivity.

---

## **Description**

ASCII Art Stylize is an enhanced version of the ASCII Art Web application, designed with the following objectives:

- **More Appealing**: Modern and visually attractive design.
- **Interactive and Intuitive**: Improved navigation and user engagement.
- **User-Friendly**: Simple, accessible interface.
- **Responsive**: Works seamlessly on different devices and screen sizes.
- **Consistent Design**: Uniform appearance and functionality throughout.
- **Color Accessibility**: Ensures text readability regardless of color schemes.

---

## **Features**

1. **Enhanced User Interface**:
   - Dynamic and responsive design using **CSS**.
   - Intuitive navigation and user feedback mechanisms.
   - Improved layout for a seamless user experience.

2. **Interactive Elements**:
   - Buttons, animations, and hover effects enhance interactivity.
   - Instant feedback for user actions.

3. **Accessible and Consistent Design**:
   - Adheres to web accessibility standards.
   - Consistent color schemes and layouts.

4. **Built with Go**:
   - Leverages the simplicity and efficiency of Go.
   - Maintains adherence to good coding practices.

---

## **Usage: How to Run**

### **Prerequisites**

1. **Install Go**: Ensure Go 1.19+ is installed.
2. **Clone the Repository**:

    ```bash
    git clone <repository-url>
    cd ascii-art-stylize
    ```

3. **Initialize `go.mod` File**:

    ```bash
    go mod init ascii-art-stylize
    ```

### **Steps to Run**

1. **Navigate to the Project Root Directory**:

    ```bash
    cd ascii-art-stylize
    ```

2. **Run the Application**:

    ```bash
    go run .
    ```

3. **Open the Application**:

    Navigate to:

    ```
    http://localhost:8080
    ```

4. **Stop the Server**:

    Press **CTRL+C** to stop the application.

---

## **Directory Structure**

```bash
ascii-art-stylize/
├── handlers/
│   ├── handlers.go          # Contains HTTP handler functions and core logic
│   ├── handlers_logger.go   # Contains a logger to keep history (not aplied yet)
│   ├── handlers_static.go   # Contains the handler for static files 
│   ├── handlers_utils.go    # Contains the ASCII logic and validation funcs
├── templates/
│   ├── 400.html            # 400 - Bad Request page
│   ├── 404.html            # 404 - Page Not Found
│   ├── 500.html            # 500 - Internal Server Error
│   ├── index.html          # Main HTML page
│   ├── PgConverter.html    # Converter page
│   ├── PgHome.html         # Home page
│   ├── PgProject.html      # Project page
│   ├── PgTeam.html         # Team page
├── static/
│   ├── css          
│   │   ├── fonts
│   │   │   ├──  Cybrpnuk2.ttf          # custom font
│   │   │   ├──  the_globe.ttf          # custom font
│   │   │   ├──  Vintage_Brother.ttf    # custom fonts
│   │   ├── style.css       # CSS file for styling
│   ├── images/
│   │   ├── 400.png         # 400 error image
│   │   ├── 404.png         # 404 error image
│   │   ├── 500.png         # 500 error image
│   │   ├── sntetop.jpg     # img of team
│   │   ├── ttarara.jpg     # img of team
│   │   ├── xkissas.jpg     # img of team
│   ├── banners/
│   │   ├── standard.txt     # Standard font
│   │   ├── shadow.txt       # Shadow font
│   │   ├── thinkertoy.txt   # Thinkertoy font
├── main.go                 # Entry point of the application
```

---

## **Implementation Details**

### **CSS Integration**

- **Styling**:
  - Defines visual design elements such as colors, fonts, and layouts.
  - Ensures text readability regardless of colors.

- **Responsiveness**:
  - Adapts the layout to different screen sizes using media queries.

- **Consistency**:
  - Applies a uniform style across the application.

### **Go Backend**

- Implements server logic and renders HTML templates.
- Handles HTTP requests and responses efficiently.

---

## **Testing**

### **Unit Testing (Go Code)**

1. **Run Tests**:

    ```bash
    go test ./...
    ```

2. **Coverage**:
   Includes tests for:
   - Template rendering.
   - HTTP request handling.

### **Browser Testing**

1. Open the application in a browser.
2. Verify responsiveness and interactivity.
3. Ensure color accessibility.

---

## **Authors**

- **Theocharoula Tarara** (*ttarara*)
- **Christoforos Kissas** (*xkissas*)
- **Stefanos Ntentopoulos** (*sntentop*)

---

## **License**

This project is licensed under the MIT License. Feel free to use, modify, and distribute it.

---

🎉 Enjoy building your ASCII Art Stylized web application! 🚀
