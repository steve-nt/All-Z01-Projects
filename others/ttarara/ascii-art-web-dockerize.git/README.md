
# 📦 **Dockerized Web Server Project** 🐳

This README provides detailed instructions for a web server application built with Go and containerized using Docker. The project focuses on adhering to best practices for Go and Docker, including the use of metadata, efficient resource management, and robust implementation principles.

---

## **Description**

The Dockerized Web Server Project delivers a high-performance, lightweight web server written in Go and packaged with Docker. The application adheres to the following principles:

- Implements **good coding practices** in Go.
- Ensures the **Dockerfile** follows best practices.
- Includes appropriate **metadata** for Docker objects (images, containers, etc.).
- Accounts for **garbage collection** by managing unused Docker objects.

The application is ready to run with minimal setup and offers a modular, easily extensible design.

---

## **Features**

1. **Web Server in Go:**
   - Built using idiomatic Go practices.
   - Highly efficient and reliable.
   - Logs and monitors requests and responses.

2. **Containerization with Docker:**
   - A single **Dockerfile** for seamless builds.
   - Image and container include descriptive **metadata**.
   - Automatic cleanup of unused objects to prevent resource bloat.

3. **Metadata Management:**
   - All Docker objects are labeled with relevant metadata (e.g., author, description, version).

4. **Garbage Collection:**
   - Includes commands and scripts to manage unused Docker objects.

---

## **Usage: How to Run**

### **Prerequisites**

1. **Install Docker:**
   Ensure Docker is installed and running on your system. [Download Docker](https://www.docker.com/get-started).

2. **Install Go (if modifying the code):**
   Make sure Go 1.19+ is installed.

### **Steps to Run**

1. **Clone the Repository:**

    ```bash
    git clone https://platform.zone01.gr/git/xkissas/ascii-art-web-dockerize.git
    cd ascii-art-web-dockerize
    ```

2. **Build the Docker Image And Verify the Build :**

    ```bash
    docker build -f Dockerfile -t ascii-art-web .
    ```
    ```bash
    docker images
    ```

3. **Run the Container And Verify them:**

    ```bash
    docker run -p 8080:8080 --detach --name ascii-art-web  ascii-art-web 
    ```
    ```bash
    docker ps -a
    ```


4. **Access the Application:**
   Open your browser and navigate to:

   ```
   http://localhost:8080
   ```
   
5. **Inspect the container:**
    ```bash
    docker exec -it ascii-art-web /bin/sh
    ls -l
    exit
    ```

6. **Evaluate Dockerfile metadata:**
   
    ```bash
    docker inspect ascii-art-web 
    ```

7. **Stop and Clean Up:**
   To stop the container:

    ```bash
    docker stop ascii-art-web 
    ```

   To remove unused objects:

    ```bash
    docker system prune -f
    ```

---

## **Implementation Details**

### **Code in Go**

- The server code follows best practices:
  - Proper error handling.
  - Clear and modular code structure.
  - Logging for request handling.
  
- Implements an efficient request-response cycle.

### **Docker Practices**

- **Dockerfile Design:**
  - Minimizes image size.
  - Leverages multi-stage builds for optimization.
  - Ensures proper labeling for metadata.

- **Metadata Example:**
  The Dockerfile includes labels such as:
  ```dockerfile
  LABEL maintainer="Your Name"
  LABEL description="A lightweight web server built with Go"
  LABEL version="1.0.0"
  ```

- **Garbage Collection:**
  - Encourages the use of `docker system prune` to clean up unused objects.

---

## **Authors**

- **Theocharoula Tarara** (*ttarara*)
- **Christoforos Kissas** (*xkissas*)
- **Stefanos Ntentopoulos** (*sntentop*)

---

## **License**

This project is licensed under the MIT License. Feel free to use, modify, and distribute it. 

---

🎉 Happy Coding and Containerizing! 🚀
