## FileVault

### Overview

Imagine you are part of a small team of software engineers at Dropbox tasked with developing a cutting-edge file-sharing system tailored for internal team collaboration. This system will not only streamline file sharing but also integrate robust version control, ensuring team members can track changes, revert to prior versions, and resolve conflicts seamlessly. It’s a high-impact project aimed at elevating collaboration efficiency while maintaining data integrity.

---

### Learning Objectives

By completing this project, you will:

- **Implement Microservices:** Design and develop a modular architecture using gRPC for inter-service communication.
- **Master Version Control:** Build a system to track, store, and manage file versions.
- **Work with SQL Databases:** Design and interact with relational database schemas for efficient data handling.
- **Handle Concurrency in Go:** Address challenges like simultaneous file edits with optimistic concurrency control.
- **Design Scalable Systems:** Develop a scalable and maintainable file-sharing system using modern software engineering practices.
- **Ensure Security and Compliance:** Implement authentication, authorization, and audit logging.
- **Apply Git Best Practices:** Demonstrate proficiency in Git-based collaboration

---

### Role-Play Scenario

Your team at Dropbox is working on a new internal tool called **Dropbox TeamVault**, a file-sharing system designed for small teams to collaborate on sensitive projects. The system must:

- Enable efficient file sharing and version tracking.
- Allow users to resolve file edit conflicts effectively.
- Provide a secure and transparent collaboration environment.

Your job is to design and implement the backend architecture, database, and a minimal frontend to bring this vision to life.

---

### Instructions

1. **Setup the Environment**:
    - Initialize a Go project.
    - Install gRPC and configure the environment for microservices development.
    - Set up an SQL database (e.g., PostgreSQL) with tables for users, teams, files, file versions, and permissions.
2. **Build Core Microservices**:
    - **Authentication Service**:
        - Handles user login, registration, and token-based authentication.
        - Ensures secure access to other services.
    - **File Management Service**:
        - Manages file uploads, storage, retrieval, and deduplication using hashing.
        - Stores file metadata (size, type, owner) in a database table.
    - **Version Control Service**:
        - Tracks file versions and metadata.
        - Provides APIs for viewing history, comparing versions, and reverting changes.
    - **Team Collaboration Service**:
        - Manages team memberships, roles, and file-sharing permissions.
        - Implements fine-grained access control (read-only, edit, admin).
    - **Conflict Resolution Service**:
        - Detects and resolves simultaneous edit conflicts.
        - Suggests or merges changes for supported file types.
    - **Audit and Logging Service**:
        - Logs user actions (uploads, edits, permission changes) immutably for traceability.
3. **Design APIs**:
    - Define gRPC services with protobuf for communication between microservices.
    - Example endpoints:
        - `UploadFile(file: File)`: Creates new file versions.
        - `GetFileHistory(fileId: string)`: Retrieves version history.
        - `RevertToVersion(fileId: string, version: number)`: Restores a file to a previous version.
4. **Implement Frontend**:
    - Build a simple HTML, CSS, and JavaScript interface.
    - Integrate with backend services using REST/gRPC-web.
    - Features include file upload, version history timeline, and permission management.
5. **Testing and Deployment**:
    - Write unit and integration tests for each service.
    - Deploy microservices independently for scalability.
    - Use Docker or Kubernetes for containerization and orchestration.

---

### Recommended Project Repository Structure

```
project-root/
├── auth-service/
├── file-service/
├── version-service/
├── team-service/
├── conflict-service/
├── audit-service/
├── frontend/
├── docs/
└── README.md

```

---

### Architecture Diagram

```split
graph TD
    A[Frontend] -->|gRPC-web/REST| B[API Gateway]
    B --> C[Authentication Service]
    B --> D[File Management Service]
    B --> E[Version Control Service]
    B --> F[Team Collaboration Service]
    B --> G[Conflict Resolution Service]
    B --> H[Audit and Logging Service]
    C --> I[(SQL Database)]
    D --> I
    E --> I
    F --> I
    G --> I
    H --> I
```

<p align="center">
  <a href="https://github.com/tomorrow-school/public/raw/refs/heads/main/image.png">
    <img src="https://github.com/tomorrow-school/public/raw/refs/heads/main/image.png"
         width="500"/>
  </a>
</p>

---

### Git Flow Guidelines

Follow the Git Flow methodology to ensure an organized development process:

1. **Branches**:
    - `main`: Contains the production-ready code.
    - `develop`: Integration branch for features before release.
    - `feature/*`: Individual branches for specific features.
    - `hotfix/*`: Quick fixes to the production code.
    - `release/*`: Stabilization branch before merging into `main`.
2. **Workflow**:
    - Create a `feature` branch from `develop` for each new task.
    - Once complete, merge the `feature` branch back into `develop` via a pull request.
    - For releases, create a `release` branch from `develop`, test, and merge it into `main`.
    - Address urgent production issues using a `hotfix` branch from `main` and merge it back into both `main` and `develop`.

---

### Tips

- **Concurrency Management:** Use optimistic concurrency control to handle simultaneous edits.
- **Database Optimization:** Index frequently queried fields (e.g., file IDs, user IDs) for better performance.
- **Security Best Practices:** Encrypt sensitive data and enforce strict authentication and authorization.
- **Debugging:** Use centralized logging for troubleshooting across services.

---

### Resources

- [gRPC Documentation](https://grpc.io/docs/)
- [Go Official Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [SQL Indexing Best Practices](https://use-the-index-luke.com/)
- [Concurrency in Go](https://go.dev/doc/effective_go#concurrency)
- [Git Branching Strategies](https://www.notion.so/FileVault-1a665d30e34b8028b5dfdd066ecf2a7b?pvs=21)
- [Pull Request Best Practices](https://github.blog/2015-01-21-how-to-write-the-perfect-pull-request/)
- [Git Interactive Learning](https://learngitbranching.js.org/)
- [Resolving Merge Conflicts](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/addressing-merge-conflicts/resolving-a-merge-conflicts)

---

### The task was developed by the engineers at Yandex Kazakhstan.

Good luck! The future of Dropbox TeamVault is in your hands.