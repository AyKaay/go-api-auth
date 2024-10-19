**User Authentication API**
API that allows for user registration, login, and transaction management. Built using Go and the Gin framework, it provides a simple and secure way to manage user accounts and transactions.

Install:
1. Navigate to the project root directory
2. Initialize the modules:
`go mod init`
3. Install dependencies:
`go mod tidy`
4. Build the program:
`go mod build`
Run:
1. To run the api:
`go run main.go`

API endpoints:
1. Register/Signup new user
    POST /auth/signup

    ```
        {
        "username": "your_username",
        "password": "your_password"
        }
    ```

2. Log in to existing/registered user
    POST /auth/login

    ```
        {
        "username": "your_username",
        "password": "your_password"
        }
    ```

3. Log out from current session
    POST /auth/logout

4. Create new transaction/payment (must be logged in)
    POST /transaction/create
    
    ```
    {
        "amount": "100.00",
        "type": "credit",
        "user_id": "user_id"
    }
    ```

Database (JSON) setup:
    The API uses a JSON file (db.json) for data storage.

Technologies used:
1. Golang   : Main language used
2. Gin      : Go framework for routing, middlewares, and HTTP request
3. Bcrypt   : Password hashing before storing in db.json
4. JSON     : To store user and transaction data in a readable format

Directory structures:
1. controllers/
    - handling incoming requests and returning responses.
2. initializers/
    - initializing the application and loading resources such as database connections
3. models/
    - defines the data structures used 
4. main.go
    - entry point for the app
5. go.mod
    - package dependecies list
6. db.json
    - mock database
7. .env
    - environment variables
