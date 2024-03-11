package main

import (
	"database/sql"
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
