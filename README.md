# Golang API Assessment - OneCV

Develop a backend application that will be part of a system which teachers can use to perform administrative functions for their students. Teachers and students are identified by their email addresses.

## Built with
- Golang (_server_)
- PostgreSQL (_database_)
- Go Fiber (_github.com/gofiber/fiber/v2_)

## Getting started

#### Prerequisites 

Ensure that the following packages are installed using the `go get` command.

- github.com/gofiber/fiber/v2
- github.com/joho/godotenv
- github.com/lib/pq

#### Environmental Variables

Create a `.env` file in the roof of the project, and ensure that the following variables are created. An example `.env` file is displayed below.
```env
DB_USER=username
DB_PASSWORD=password
DB_NAME=databasename
DB_HOST=localhost
DB_PORT=5432
```

#### Database

A script is included in the root of the project and can be found [here](dbscript). Ensure to run this script to initialize the database before running the server.

dbscript:
```sql
-- Drop tables if they exist
DROP TABLE IF EXISTS teacher_student;
DROP TABLE IF EXISTS teachers;
DROP TABLE IF EXISTS students;

-- Create teachers table
CREATE TABLE teachers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE
);

-- Create students table
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    isSuspended BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create table to represent the many-to-many relationship between teachers and students
CREATE TABLE teacher_student (
    teacher_id INT REFERENCES teachers(id),
    student_id INT REFERENCES students(id),
    PRIMARY KEY (teacher_id, student_id)
);
```

#### Running the project

If you are doing a fresh installation of the project, please initialize the project by running the `go mod init` command first.

- ` go mod init example/golang-api-assessment-govtech`

And to start the server, navigate to the root of the project and run the following command:

- `go run main.go server.go database.go queries.go`
