# LRUProject


Brief description of the project.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)


## Installation

### Server

1. Open the project directory in your terminal.
2. Navigate to the server directory:
    ```
    cd server
    ```
3. Run the following command to start the server:
    ```
    go run main.go
    ```

### Frontend

1. Navigate back to the project root directory:
    ```
    cd ..
    ```
2. Navigate to the frontend directory:
    ```
    cd frontend
    ```
3. Install project dependencies using npm:
    ```
    npm install
    ```
4. Start the frontend application:
    ```
    npm start
    ```

### Explanation

#### Frontend Interface:
    Users can interact with the cache management system through a user-friendly React frontend.
    The interface enables users to set key-value pairs as cache entries and retrieve cached values within the specified time limit.

#### Backend Implementation:
    The backend is built using Golang and the Gin framework, providing robust and efficient handling of cache operations.
    Custom implementations, including a doubly linked list, are utilized to store and manage cache entries effectively.

#### Cache Management Policies:
    The system employs a size limit for cache storage, with a maximum capacity of 1024 entries.
    When the storage limit is reached, the least recently used cache entry is automatically removed from the system.

#### Automatic Cleanup Mechanism:
    A parallel cleanup function runs continuously, checking each cache entry every second.
    If the expiration time limit of a cache entry is exceeded, it is promptly removed from the system to ensure efficient memory utilization.