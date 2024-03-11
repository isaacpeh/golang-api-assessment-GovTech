package main

import (
	"database/sql"
	"strconv"
	"strings"
)

func registerTeacherStudent(teacherEmail string, studentEmails []string) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Rollback if any errors
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		// Commit transaction if no errors
		err = tx.Commit()
	}()

	// Register teacher if not exists
	var teacherID int
	err = tx.QueryRow("SELECT id FROM teachers WHERE email = $1", teacherEmail).Scan(&teacherID)
	if err == sql.ErrNoRows {
		err = tx.QueryRow("INSERT INTO teachers (email) VALUES ($1) RETURNING id", teacherEmail).Scan(&teacherID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Register each student if not exists
	for _, studentEmail := range studentEmails {
		var studentID int
		err = tx.QueryRow("SELECT id FROM students WHERE email = $1", studentEmail).Scan(&studentID)
		if err == sql.ErrNoRows {
			err = tx.QueryRow("INSERT INTO students (email) VALUES ($1) RETURNING id", studentEmail).Scan(&studentID)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		// Register the student to the teacher
		_, err = tx.Exec("INSERT INTO teacher_student (teacher_id, student_id) VALUES ($1, $2)", teacherID, studentID)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCommonStudents(teachers []string) ([]string, error) {
	// Enclose each email in single quotes
	quotedEmails := make([]string, len(teachers))
	for i, teacher := range teachers {
		quotedEmails[i] = "'" + teacher + "'"
	}

	query := `
		SELECT s.email
		FROM students s
		INNER JOIN teacher_student ts ON s.id = ts.student_id
		INNER JOIN teachers t ON ts.teacher_id = t.id
		WHERE t.email IN (` + strings.Join(quotedEmails, ", ") + `)
		GROUP BY s.email
		HAVING COUNT(DISTINCT t.email) = ` + strconv.Itoa(len(teachers))

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect the common students into a slice
	var commonStudents []string

	// Loop each row returned by db
	for rows.Next() {
		var studentEmail string
		if err := rows.Scan(&studentEmail); err != nil {
			return nil, err
		}
		commonStudents = append(commonStudents, studentEmail)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commonStudents, nil
}
