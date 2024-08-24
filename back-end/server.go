package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

type CompanyDetails struct {
	CompanyName  string `json:"companyName"`
	Designation  string `json:"designation"`
	JoinDate     string `json:"joinDate"`
	RelievedDate string `json:"relievedDate"`
	Duration     string `json:"duration"`
	id           int    `json:"id`
}

type Employee struct {
	ID                     int            `json:"id"`
	EmpId                  int            `json:"empId"`
	FirstName              string         `json:"firstName"`
	LastName               string         `json:"lastName"`
	Email                  string         `json:"email"`
	PhoneNo                int            `json:"phoneNo"`
	FatherName             string         `json:"fatherName"`
	EmergencyContact       int            `json:"emergencyContact"`
	DateOfBirth            string         `json:"dateOfBirth"`
	Address                string         `json:"address"`
	Qualification          string         `json:"qualification"`
	Experience             bool           `json:"experience"`
	CreatedTime            string         `json:"created_time"`
	CompanyName            string         `json:"companyName"`
	Designation            string         `json:"designation"`
	JoinDate               string         `json:"joinDate"`
	RelievedDate           string         `json:"relievedDate"`
	TotalDuration          string         `json:"totalDuration"`
	SecondCompanyFormValue CompanyDetails `json:"secondCompanyFormValue"`
}

type EmployeeWithCompany struct {
	Employee Employee         `json:"employee"`
	Company  []CompanyDetails `json:"company"`
}

//

//get

func dataHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(` SELECT empId,firstName,lastName,email,phoneNo,fatherName,emergencyContact,dateOfBirth,address,experience,qualification,created_time FROM emply`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.EmpId, &emp.FirstName, &emp.LastName, &emp.Email, &emp.PhoneNo, &emp.FatherName, &emp.EmergencyContact, &emp.DateOfBirth, &emp.Address, &emp.Experience, &emp.Qualification, &emp.CreatedTime); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved data: %v", employees)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

//
//

// func prevcompanybyid(w http.ResponseWriter, r *http.Request) {
// 	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
// 	if err != nil {
// 		log.Printf("Error opening database: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	params := mux.Vars(r)
// 	idStr := params["id"]

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		log.Printf("Invalid employee ID: %v", err)
// 		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
// 		return
// 	}
// 	row, err := db.Query(` SELECT id,companyName,position,startDate,endDate,duration FROM prevcompany where empId=?`, id)
// 	if err != nil {
// 		log.Printf("Error executing query: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	defer row.Close()

// 	var company []CompanyDetails
// 	for row.Next() {
// 		var emp CompanyDetails
// 		if err := row.Scan(&emp.id, &emp.CompanyName, &emp.Designation, &emp.JoinDate, &emp.RelievedDate, &emp.Duration); err != nil {
// 			log.Printf("Error scanning row: %v", err)
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		company = append(company, emp)
// 	}

// 	if err := row.Err(); err != nil {
// 		log.Printf("Rows error: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	log.Printf("Retrieved company data: %v", company)

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(company); err != nil {
// 		log.Printf("Error encoding response: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func prevcompanybyid(tx *sql.Tx, empId string) ([]CompanyDetails, error) {
	rows, err := tx.Query(`SELECT id, companyName, position, startDate, endDate, duration FROM prevcompany WHERE empId=?`, empId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []CompanyDetails
	for rows.Next() {
		var company CompanyDetails
		if err := rows.Scan(&company.id, &company.CompanyName, &company.Designation, &company.JoinDate, &company.RelievedDate, &company.Duration); err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}

func emplybyid(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}
	row, err := db.Query(` SELECT id,empId,firstName,lastName,email,phoneNo,fatherName,emergencyContact,dateOfBirth,address,experience,qualification FROM emply where empId=?`, id)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer row.Close()

	var emply []Employee
	for row.Next() {
		var emp Employee
		if err := row.Scan(&emp.ID, &emp.EmpId, &emp.FirstName, &emp.LastName, &emp.Email, &emp.PhoneNo, &emp.FatherName, &emp.EmergencyContact, &emp.DateOfBirth, &emp.Address, &emp.Experience, &emp.Qualification); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		emply = append(emply, emp)
	}

	if err := row.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved company data: %v", emply)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(emply); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//
//

//employee and company by id

func employeeWithCompanyById(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var result EmployeeWithCompany

	// Query for employee details
	err = db.QueryRow("SELECT empId,firstName,lastName,email,phoneNo,fatherName,emergencyContact,dateOfBirth,address,experience,qualification,created_time FROM emply WHERE empId = ?", id).Scan(
		&result.Employee.EmpId,
		&result.Employee.FirstName,
		&result.Employee.LastName,
		&result.Employee.Email,
		&result.Employee.PhoneNo,
		&result.Employee.FatherName,
		&result.Employee.EmergencyContact,
		&result.Employee.DateOfBirth,
		&result.Employee.Address,
		&result.Employee.Experience,
		&result.Employee.Qualification,
		&result.Employee.CreatedTime,
	)
	if err != nil {
		log.Printf("Error executing query for employee: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query for previous company details
	rows, err := db.Query("SELECT companyName, position, startDate, endDate, duration FROM prevcompany WHERE empId = ?", id)
	if err != nil {
		log.Printf("Error executing query for previous company: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var company CompanyDetails
		if err := rows.Scan(&company.CompanyName, &company.Designation, &company.JoinDate, &company.RelievedDate, &company.Duration); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result.Company = append(result.Company, company)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved employee with company data: %v", result)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//

//add

func addEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Optionally, validate or process emp data here

	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Rollback if any error occurs before commit

	stmt, err := tx.Prepare("INSERT INTO emply (empId, firstName, lastName, email, phoneNo, fatherName, emergencyContact, dateOfBirth, address, qualification,experience) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute insert into emply table
	_, err = stmt.Exec(emp.EmpId, emp.FirstName, emp.LastName, emp.Email, emp.PhoneNo, emp.FatherName, emp.EmergencyContact, emp.DateOfBirth, emp.Address, emp.Qualification, emp.Experience)
	if err != nil {
		log.Printf("Error executing insert statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if emp.Experience {
		// Insert into prevcompany table
		companyStmt, err := tx.Prepare("INSERT INTO prevcompany (companyName, position, startDate, endDate, duration, empId) VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Printf("Error preparing company statement: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer companyStmt.Close()

		_, err = companyStmt.Exec(emp.CompanyName, emp.Designation, emp.JoinDate, emp.RelievedDate, emp.TotalDuration, emp.EmpId)
		if err != nil {
			log.Printf("Error executing company insert statement: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// log.Printf(_)
		// Insert into second company table if needed
		if emp.SecondCompanyFormValue.CompanyName != "" {
			_, err = companyStmt.Exec(emp.SecondCompanyFormValue.CompanyName, emp.SecondCompanyFormValue.Designation, emp.SecondCompanyFormValue.JoinDate, emp.SecondCompanyFormValue.RelievedDate, emp.SecondCompanyFormValue.Duration, emp.EmpId)
			if err != nil {
				log.Printf("Error executing second company insert statement: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Employee added successfully: %+v", emp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

//delete

func deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Delete from prevcompany table first
	_, err = db.Exec("DELETE FROM prevcompany WHERE empId = ?", id)
	if err != nil {
		log.Printf("Error deleting from prevcompany: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM employeeRole WHERE employee_id = ?", id)
	if err != nil {
		log.Printf("Error deleting from prevcompany: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Prepare("DELETE FROM employee_logs WHERE empId = ?")
	if err != nil {
		log.Printf("Error preparing delete statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("DELETE FROM emply WHERE empId = ?")
	if err != nil {
		log.Printf("Error preparing delete statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = db.Prepare("DELETE FROM timesheet WHERE empId = ?")
	if err != nil {
		log.Printf("Error preparing delete statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		log.Printf("Error executing delete statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Printf("No employee found with ID: %d", id)
		http.Error(w, "No employee found with the given ID", http.StatusNotFound)
		return
	}

	log.Printf("Employee deleted successfully: ID %d", id)
	w.WriteHeader(http.StatusNoContent)
}

//update

func updateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var emp Employee
	// var company CompanyDetails

	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE emply SET firstName=?, lastName=?, email=?, phoneNo=?, fatherName=?, emergencyContact=?, dateOfBirth=?, address=?, experience=? WHERE empId=?", emp.FirstName, emp.LastName, emp.Email, emp.PhoneNo, emp.FatherName, emp.EmergencyContact, emp.DateOfBirth, emp.Address, emp.Experience, id)
	if err != nil {
		log.Printf("Error updating employee: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the employee has experience, update their previous company records
	if emp.Experience {
		// Fetch company records for the employee
		companies, err := prevcompanybyid(tx, id)
		if err != nil {
			log.Printf("Error retrieving company records: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update the first company record
		if len(companies) > 0 {
			_, err = tx.Exec(
				"UPDATE prevcompany SET companyName=?, position=?, startDate=?, endDate=?, duration=? WHERE id=?",
				emp.CompanyName,
				emp.Designation,
				emp.JoinDate,
				emp.RelievedDate,
				emp.TotalDuration,
				companies[0].id, // Use the first company's ID
			)
			if err != nil {
				log.Printf("Error updating first company: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Update the second company record if necessary
		if len(companies) > 1 && emp.SecondCompanyFormValue.CompanyName != "" {
			_, err = tx.Exec(
				"UPDATE prevcompany SET companyName=?, position=?, startDate=?, endDate=?, duration=? WHERE id=?",
				emp.SecondCompanyFormValue.CompanyName,
				emp.SecondCompanyFormValue.Designation,
				emp.SecondCompanyFormValue.JoinDate,
				emp.SecondCompanyFormValue.RelievedDate,
				emp.SecondCompanyFormValue.Duration,
				companies[1].id, // Use the second company's ID
			)
			if err != nil {
				log.Printf("Error updating second company: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {

		_, err = db.Exec("DELETE FROM prevcompany WHERE empId = ?", id)
		if err != nil {
			log.Printf("Error deleting from prevcompany: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	log.Printf("Employee updated successfully: %+v", emp)

	if emp.Experience {
		_, err = tx.Exec(
			"INSERT INTO prevcompany (companyName, position, startDate, endDate, duration, empId) VALUES (?, ?, ?, ?, ?, ?)",
			emp.CompanyName, emp.Designation, emp.JoinDate, emp.RelievedDate, emp.TotalDuration, id)
		if err != nil {
			log.Printf("Error inserting company record: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if emp.SecondCompanyFormValue.CompanyName != "" {
			_, err = tx.Exec(
				"INSERT INTO prevcompany (companyName, position, startDate, endDate, duration, empId) VALUES (?, ?, ?, ?, ?, ?)",
				emp.SecondCompanyFormValue.CompanyName, emp.SecondCompanyFormValue.Designation, emp.SecondCompanyFormValue.JoinDate, emp.SecondCompanyFormValue.RelievedDate, emp.SecondCompanyFormValue.Duration, id)
			if err != nil {
				log.Printf("Error inserting second company record: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	log.Printf("Employee updated successfully: %+v", emp)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Employee added successfully: %+v", emp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emp)
}

type Role struct {
	RoleID   int    `json:"role_id"`
	RoleName string `json:"roleName"`
}

func getRoleById(w http.ResponseWriter, r *http.Request) {
	// Extract role ID from URL parameters
	vars := mux.Vars(r)
	roleid, ok := vars["id"]
	if !ok {
		http.Error(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	// Open database connection
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Prepare the query
	query := `SELECT role_id,roleName FROM role WHERE role_id = ?`

	// Execute the query with parameter
	row := db.QueryRow(query, roleid)

	// Scan the result into a Role struct
	var role Role
	if err := row.Scan(&role.RoleID, &role.RoleName); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No role found with the given ID", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Send the result as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getRolesByDepartment handles requests to fetch roles by department ID
func getRolesByDepartment(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	vars := mux.Vars(r)
	deptID, ok := vars["id"]
	if !ok {
		http.Error(w, "Department ID is required", http.StatusBadRequest)
		return
	}

	// Query to get roles based on department ID
	query := `SELECT role_id, roleName FROM role WHERE dept_id = ?`
	rows, err := db.Query(query, deptID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.RoleID, &role.RoleName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		roles = append(roles, role)
	}

	// Convert the result to JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(roles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Department represents the department structure
type Department struct {
	DeptID         int    `json:"dept_id"`
	DepartmentName string `json:"department"`
}

// getDepartments handles requests to fetch all departments
func getDepartments(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	query := `SELECT dept_id, department FROM department`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var departments []Department
	for rows.Next() {
		var dept Department
		if err := rows.Scan(&dept.DeptID, &dept.DepartmentName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		departments = append(departments, dept)
	}
	// Convert the result to JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(departments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getDepartmentsById(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(r)
	deptID, ok := vars["id"]
	if !ok {
		http.Error(w, "Department ID is required", http.StatusBadRequest)
		return
	}

	query := `SELECT dept_id,department FROM department where dept_id = ?`
	rows, err := db.Query(query, deptID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var departments []Department
	for rows.Next() {
		var dept Department
		if err := rows.Scan(&dept.DeptID, &dept.DepartmentName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		departments = append(departments, dept)
	}
	// 	// Convert the result to JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(departments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Manager represents the manager structure
type Manager struct {
	ManagerID   int    `json:"manager_id"`
	ManagerName string `json:"managerName"`
}

// getManagers handles requests to fetch all managers
func getManagers(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `SELECT manager_id, managerName FROM manager`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var managers []Manager
	for rows.Next() {
		var mgr Manager
		if err := rows.Scan(&mgr.ManagerID, &mgr.ManagerName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		managers = append(managers, mgr)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(managers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getManagerByDepartment handles requests to fetch the manager by department ID
func getManagerByDepartment(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(r)
	deptID, ok := vars["id"]
	if !ok {
		http.Error(w, "Department ID is required", http.StatusBadRequest)
		return
	}
	// deptID := r.URL.Query().Get("dept_id")
	query := `SELECT m.manager_id, m.managerName 
              FROM manager m
              JOIN department d ON m.manager_id = d.manager_id
              WHERE d.dept_id = ?`
	row := db.QueryRow(query, deptID)

	var mgr Manager
	if err := row.Scan(&mgr.ManagerID, &mgr.ManagerName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mgr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getEmployeeAsManager(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	vars := mux.Vars(r)
	deptID, ok := vars["id"]
	if !ok {
		http.Error(w, "Department ID is required", http.StatusBadRequest)
		return
	}
	query := `select firstName,lastName from employeeRole where tech_lead is NUll and dept_id = ?`
	row := db.QueryRow(query, deptID)

	var mgr Employee
	if err := row.Scan(&mgr.FirstName, &mgr.LastName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mgr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

type AssignEmployee struct {
	EmpID    int `json:"employee_id"`
	RoleID   int `json:"role_id"`
	DeptID   int `json:"dept_id"`
	TechLead int `json:"tech_lead"`
	// TechLead struct {
	// 	Type string `json:"type"`
	// 	ID   int    `json:"id"`
	// } `json:"tech_lead"`
}

func assignEmployee(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var ass AssignEmployee
	log.Printf("%v", ass)
	// Decode JSON request body into AssignEmployee struct
	err = json.NewDecoder(r.Body).Decode(&ass)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Assigning Employee: %+v", ass)

	// Check if the employee is already assigned to the same role and department
	var exists bool
	query := `SELECT EXISTS(SELECT 6 FROM employeeRole WHERE employee_id = ?)`
	err = db.QueryRow(query, ass.EmpID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking existing assignment: %v", err)
		http.Error(w, "Error checking existing assignment", http.StatusInternalServerError)
		return
	}

	if exists {
		log.Printf("Employee is already assigned to this role and department")
		http.Error(w, "Employee is already assigned to this role and department", http.StatusConflict)
		return
	}
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO employeeRole (employee_id, role_id, dept_id,tech_lead) VALUES (? , ? , ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, "Error preparing statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(ass.EmpID, ass.RoleID, ass.DeptID, ass.TechLead)
	if err != nil {
		log.Printf("Error executing insert statement: %v", err)
		http.Error(w, "Error executing insert statement", http.StatusInternalServerError)
		return
	}

	// Prepare the SQL statement for UPDATE
	stmt, err = db.Prepare("UPDATE emply SET role_id = ?, dept_id = ?, tech_lead = ? WHERE empId = ?")
	if err != nil {
		log.Printf("Error preparing update statement: %v", err)
		http.Error(w, "Error preparing update statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement for UPDATE
	_, err = stmt.Exec(ass.RoleID, ass.DeptID, ass.TechLead, ass.EmpID)
	if err != nil {
		log.Printf("Error executing update statement: %v", err)
		http.Error(w, "Error executing update statement", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ass); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type TechLead struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
}

type HierarchyData struct {
	EmpID       int    `json:"empId"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	RoleName    string `json:"roleName"`
	RoleID      int    `json:"role_id"`
	Department  string `json:"department"`
	ManagerName string `json:"managerName"`
	TechLead    int    `json:"techLead"`
}

// Fetch hierarchy data from the database
func getHierarchyData(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// SQL query to fetch employee hierarchy data
	query := `
		SELECT 
			e.empId,
			e.firstName,
			e.lastName,
			r.roleName,
			r.role_id,
			d.department,
			m.managerName,
			er.tech_lead
		FROM 
			employeeRole er
		JOIN 
			emply e ON er.employee_id = e.empId
		JOIN 
			role r ON er.role_id = r.role_id
		JOIN 
			department d ON er.dept_id = d.dept_id
		JOIN 
			manager m ON d.manager_id = m.manager_id
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var hierarchyData []HierarchyData
	for rows.Next() {
		var techLeadJSON string
		var data HierarchyData

		if err := rows.Scan(&data.EmpID, &data.FirstName, &data.LastName, &data.RoleName, &data.RoleID, &data.Department, &data.ManagerName, &techLeadJSON); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		// Log the tech_lead JSON string
		log.Printf("TechLead JSON string: %s", techLeadJSON)
		// Deserialize tech_lead JSON string into TechLead struct
		if err := json.Unmarshal([]byte(techLeadJSON), &data.TechLead); err != nil {
			log.Printf("Error unmarshalling tech_lead JSON: %v", err)
			http.Error(w, "Error processing tech_lead data", http.StatusInternalServerError)
			return
		}

		hierarchyData = append(hierarchyData, data)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		http.Error(w, "Error reading data", http.StatusInternalServerError)
		return
	}

	// Send a success response with the hierarchy data
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(hierarchyData); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func getHierarchyDataById(w http.ResponseWriter, r *http.Request) {
	// Get the employee ID from the URL parameters
	vars := mux.Vars(r)
	empIdStr := vars["id"]

	// Convert the employee ID to an integer
	empId, err := strconv.Atoi(empIdStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// SQL query to fetch employee hierarchy data by ID
	query := `
		SELECT 
			e.empId,
			e.firstName,
			e.lastName,
			r.roleName,
            r.role_id,
			d.department,
			m.managerName
		FROM 
			employeeRole er
		JOIN 
			emply e ON er.employee_id = e.empId
		JOIN 
			role r ON er.role_id = r.role_id
		JOIN 
			department d ON er.dept_id = d.dept_id
		JOIN 
			manager m ON d.manager_id = m.manager_id
		WHERE 
			er.role_id = ?
	`

	// Execute the query with the employee ID
	row := db.QueryRow(query, empId)

	var data HierarchyData
	if err := row.Scan(&data.EmpID, &data.FirstName, &data.LastName, &data.RoleName, &data.RoleID, &data.Department, &data.ManagerName); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Employee not found", http.StatusNotFound)
		} else {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
		}
		return
	}

	// Send a success response with the hierarchy data
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

type Documents struct {
	ID       int    `json:"id"`
	EmpID    int    `json:"empId"`
	Filename string `json:"filename"`
	FileData []byte `json:"filedata"`
}

// /////
func uploadDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values
	empId := r.FormValue("empId")

	// Extract files
	aadharFile, aadharHeader, err := r.FormFile("aadhar")
	if err != nil {
		http.Error(w, "Unable to get aadhar file", http.StatusBadRequest)
		return
	}
	defer aadharFile.Close()

	profilePhoto, _, err := r.FormFile("profilephoto")
	if err != nil {
		http.Error(w, "Unable to get profile photo", http.StatusBadRequest)
		return
	}
	defer profilePhoto.Close()

	// Handle aadhar file: check if it's a DOC/DOCX file and convert to PDF if necessary
	aadharFileBytes, err := handleFileConversion(aadharFile, aadharHeader.Filename)
	if err != nil {
		http.Error(w, "Unable to process aadhar file", http.StatusInternalServerError)
		return
	}

	// Read profile photo into memory
	profilePhotoBytes, err := io.ReadAll(profilePhoto)
	if err != nil {
		http.Error(w, "Unable to read profile photo", http.StatusInternalServerError)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert file data into database
	_, err = db.Exec("INSERT INTO documents (empId, aadharFile, profilePhoto) VALUES (?, ?, ?)",
		empId, aadharFileBytes, profilePhotoBytes)
	if err != nil {
		http.Error(w, "Error saving files to database", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Files uploaded successfully"))
}

// handleFileConversion checks the file type and converts DOC/DOCX files to PDF
func handleFileConversion(file multipart.File, filename string) ([]byte, error) {
	// Check file extension
	ext := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(filepath.Ext(filename)), "."))

	if ext == "doc" || ext == "docx" {
		// Convert DOC/DOCX to PDF
		return convertDocToPDF(file)
	}

	// For other file types, simply read the file and return the bytes
	return io.ReadAll(file)
}

// convertDocToPDF converts a DOC/DOCX file to PDF and returns the PDF as []byte
func convertDocToPDF(file multipart.File) ([]byte, error) {
	// Create a temporary file to save the uploaded DOC/DOCX
	tempFile, err := os.CreateTemp("", "upload-*.doc")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name()) // Clean up the file after conversion

	// Copy the uploaded file to the temp file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	// Convert DOC/DOCX to PDF using unoconv or LibreOffice
	outputPDF := tempFile.Name() + ".pdf"
	cmd := exec.Command("unoconv", "-f", "pdf", "-o", outputPDF, tempFile.Name())
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error converting DOC to PDF: %w", err)
	}
	defer os.Remove(outputPDF) // Clean up the output PDF file after reading it

	// Read the converted PDF file into a byte slice
	pdfFile, err := os.Open(outputPDF)
	if err != nil {
		return nil, err
	}
	defer pdfFile.Close()

	return io.ReadAll(pdfFile)
}

type PersonalDetails struct {
	EmpID        int    `json:"empId"`
	Gender       string `json:"gender"`
	Relationship string `json:"relationship"`
	BloodGroup   string `json:"bloodgroup"`
}

func personaldetails(w http.ResponseWriter, r *http.Request) {
	var emp PersonalDetails

	// Decode the JSON request body into the PersonalDetails struct
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Printf("Error decoding JSON request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Rollback if any error occurs before commit

	// Prepare the SQL insert statement
	stmt, err := tx.Prepare("INSERT INTO personaldetails (empId, gender, relationship) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the insert statement
	_, err = stmt.Exec(emp.EmpID, emp.Gender, emp.Relationship)
	if err != nil {
		log.Printf("Error executing insert statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Personal details added successfully: %+v", emp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

type doc struct {
	Aadharfilename       string `json:"aadhar_filename"`
	ProfilePhotofilename string `json:"profilephoto_filename"`
}

// ////////////
func handlePersonalDetailsAndDocuments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	var aadharfilename string
	var profilephotofilename string
	// Extract form values
	empId := r.FormValue("empId")
	aadharfilename = r.FormValue("aadharfilename")
	profilephotofilename = r.FormValue("profilephotofilename")

	// Extract JSON data from form data
	jsonData := r.FormValue("personalDetails")
	if jsonData == "" {
		http.Error(w, "Missing personal details", http.StatusBadRequest)
		return
	}

	// Parse JSON data into PersonalDetails struct
	var emp PersonalDetails
	err = json.Unmarshal([]byte(jsonData), &emp)
	if err != nil {
		http.Error(w, "Error decoding JSON request body", http.StatusBadRequest)
		return
	}

	// Extract files
	aadharFile, _, err := r.FormFile("aadhar")
	if err != nil {
		http.Error(w, "Unable to get aadhar file", http.StatusBadRequest)
		return
	}
	defer aadharFile.Close()

	profilePhoto, _, err := r.FormFile("profilephoto")
	if err != nil {
		http.Error(w, "Unable to get profile photo", http.StatusBadRequest)
		return
	}
	defer profilePhoto.Close()

	// Read files into memory
	aadharFileBytes, err := io.ReadAll(aadharFile)
	if err != nil {
		http.Error(w, "Unable to read aadhar file", http.StatusInternalServerError)
		return
	}

	profilePhotoBytes, err := io.ReadAll(profilePhoto)
	if err != nil {
		http.Error(w, "Unable to read profile photo", http.StatusInternalServerError)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Rollback if any error occurs before commit

	// Prepare and execute the SQL insert statement for personal details
	stmt, err := tx.Prepare("INSERT INTO personaldetails (empId, gender, relationship, bloodgroup) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(emp.EmpID, emp.Gender, emp.Relationship, emp.BloodGroup)
	if err != nil {
		log.Printf("Error executing personal details insert statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare and execute the SQL insert statement for file data
	_, err = tx.Exec("INSERT INTO documents (empId, aadharFile, profilePhoto,aadhar_filename,profilephoto_filename) VALUES (?, ? ,?, ?, ?)",
		empId, aadharFileBytes, profilePhotoBytes, aadharfilename, profilephotofilename)
	if err != nil {
		log.Printf("Error executing file data insert statement: %v", err)
		http.Error(w, "Error saving files to database", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Personal details and files added successfully: %+v", emp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

func handleUpdatePersonalDetailsAndDocuments(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values
	empId := r.FormValue("empId")
	if empId == "" {
		http.Error(w, "Employee ID is missing", http.StatusBadRequest)
		return
	}

	var aadharfilename string
	var profilephotofilename string
	// Extract form values
	aadharfilename = r.FormValue("aadharfilename")
	profilephotofilename = r.FormValue("profilephotofilename")

	// Extract JSON data from form data
	jsonData := r.FormValue("personalDetails")
	if jsonData == "" {
		http.Error(w, "Missing personal details", http.StatusBadRequest)
		return
	}

	// Parse JSON data into PersonalDetails struct
	var emp PersonalDetails
	err = json.Unmarshal([]byte(jsonData), &emp)
	if err != nil {
		http.Error(w, "Error decoding JSON request body", http.StatusBadRequest)
		return
	}

	// Extract files
	aadharFile, _, err := r.FormFile("aadhar")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Unable to get aadhar file", http.StatusBadRequest)
		return
	}
	defer aadharFile.Close()

	profilePhoto, _, err := r.FormFile("profilephoto")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Unable to get profile photo", http.StatusBadRequest)
		return
	}
	defer profilePhoto.Close()

	var aadharFileBytes []byte
	if aadharFile != nil {
		aadharFileBytes, err = io.ReadAll(aadharFile)
		if err != nil {
			http.Error(w, "Unable to read aadhar file", http.StatusInternalServerError)
			return
		}
	}

	var profilePhotoBytes []byte
	if profilePhoto != nil {
		profilePhotoBytes, err = io.ReadAll(profilePhoto)
		if err != nil {
			http.Error(w, "Unable to read profile photo", http.StatusInternalServerError)
			return
		}
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Rollback if any error occurs before commit

	// Prepare and execute the SQL update statement for personal details
	stmt, err := tx.Prepare(`
        UPDATE personaldetails
        SET gender = ?, relationship = ?, bloodgroup = ?
        WHERE empId = ?`)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(emp.Gender, emp.Relationship, emp.BloodGroup, empId)
	if err != nil {
		log.Printf("Error executing personal details update statement: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare and execute the SQL update statement for file data
	_, err = tx.Exec(`
        UPDATE documents
        SET aadharFile = ?, profilePhoto = ?,aadhar_filename = ?,profilephoto_filename = ?
        WHERE empId = ?`,
		aadharFileBytes, profilePhotoBytes, aadharfilename, profilephotofilename, empId)
	if err != nil {
		log.Printf("Error executing file data update statement: %v", err)
		http.Error(w, "Error saving files to database", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Personal details and files updated successfully: %+v", emp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emp)
}

type DocumentResponse struct {
	AadharImage          string `json:"aadharImage"`
	ProfilePhoto         string `json:"profilePhoto"`
	Gender               string `json:"gender"`
	Relationship         string `json:"relationship"`
	BloodGroup           string `json:"bloodgroup"`
	AadharFileName       string `json:"aadhar_filename"`
	ProfilePhotoFileName string `json:"profilephoto_filename"`
}

func getDocuments(w http.ResponseWriter, r *http.Request) {
	// Get the employee ID from the URL parameters
	vars := mux.Vars(r)
	empIdStr := vars["id"]

	// Convert the employee ID to an integer
	empId, err := strconv.Atoi(empIdStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// type doc struct{
	// 	Aadharfilename sql.NullString `json:"aadhar_filename"`
	// 	Profilephotoname sql.NullString `json:"profilephoto_filename"`
	// 	Gender
	// }
	// Fetch the document and additional data from the database
	var aadharFileBytes, profilePhotoBytes []byte
	var gender, relationship, bloodGroup string
	var aadharfilename, profilephotoname string
	err = db.QueryRow(`
		SELECT d.aadharFile, d.profilePhoto,d.aadhar_filename,d.profilephoto_filename ,a.gender,a.relationship,a.bloodgroup
		FROM documents d
		JOIN personaldetails a ON d.empId = a.empId
		WHERE d.empId = ?
	`, empId).Scan(&aadharFileBytes, &profilePhotoBytes, &aadharfilename, &profilephotoname, &gender, &relationship, &bloodGroup)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No rows found for empId %d: %v", empId, err)
			http.Error(w, "Document or additional data not found", http.StatusNotFound)
		} else {
			log.Printf("Failed to fetch data for empId %d: %v", empId, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Encode files to Base64
	aadharBase64 := base64.StdEncoding.EncodeToString(aadharFileBytes)
	profilePhotoBase64 := base64.StdEncoding.EncodeToString(profilePhotoBytes)
	// log.Printf("aadhar base 64 %s", aadharBase64)

	// Create a JSON response
	response := DocumentResponse{
		AadharImage:          aadharBase64,
		ProfilePhoto:         profilePhotoBase64,
		Gender:               gender,
		Relationship:         relationship,
		BloodGroup:           bloodGroup,
		AadharFileName:       aadharfilename,
		ProfilePhotoFileName: profilephotoname,
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// /
// setManager handles fetching managers from the database based on department and tech lead
func setManager(w http.ResponseWriter, r *http.Request) {
	// Get the employee ID from the URL parameters
	vars := mux.Vars(r)
	empIdStr := vars["id"]

	// Convert the employee ID to an integer
	id, err := strconv.Atoi(empIdStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare the query with parameters
	query := `SELECT empId FROM emply WHERE dept_id=? AND tech_lead IS NULL`

	// Query the database with the provided department ID
	rows, err := db.Query(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var managers []Employee
	for rows.Next() {
		var mgr Employee
		if err := rows.Scan(&mgr.EmpId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		managers = append(managers, mgr)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(managers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getAssignData(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(` SELECT * FROM employeeRole`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var employees []AssignEmployee
	for rows.Next() {
		var emp AssignEmployee
		if err := rows.Scan(&emp.EmpID, &emp.RoleID, &emp.DeptID, &emp.TechLead); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved data: %v", employees)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

type User struct {
	EmpId          int            `json:"empId"`
	FirstName      string         `json:"firstName"`
	LastName       string         `json:"lastName"`
	Email          string         `json:"email"`
	LoginTime      string         `json:"login_time"`
	LogoutTime     sql.NullString `json:"logout_time"`
	WorkingHours   sql.NullString `json:"working_hours"`
	Password       string         `json:"password"`
	FaceEmbeddings []float32      `json:"faceEmbeddings"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON payload
	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the empId already exists
	var existingUserId int
	checkQuery := "SELECT empId FROM timesheet WHERE empId = ?"
	err = db.QueryRow(checkQuery, user.EmpId).Scan(&existingUserId)
	if err == nil {
		// If no error is returned, it means the empId exists
		http.Error(w, "User with this empId already exists", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		// If there's an error other than no rows found, log it
		log.Printf("Error checking existing user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert face embeddings to JSON
	faceEmbeddingsJSON, err := json.Marshal(user.FaceEmbeddings)
	if err != nil {
		log.Printf("Error marshalling face embeddings: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Insert into database
	query := "INSERT INTO timesheet (empId, firstName, lastName, email, password, face_embeddings) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, user.EmpId, user.FirstName, user.LastName, user.Email, string(hashedPassword), string(faceEmbeddingsJSON))
	if err != nil {
		log.Printf("Error saving user: %v", err)
		http.Error(w, "Failed to register: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Registration successful")
}

func getTimesheet(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(` SELECT 
    e.empId,
    e.firstName,
    e.lastName,
	e.email,
    el.login_time,
    el.logout_time,
    el.working_hours
FROM 
    employee_logs el
JOIN 
    emply e ON e.empId = el.empId;`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var employees []User
	for rows.Next() {
		var emp User
		if err := rows.Scan(&emp.EmpId, &emp.FirstName, &emp.LastName, &emp.Email, &emp.LoginTime, &emp.LogoutTime, &emp.WorkingHours); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved data: %v", employees)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func getTimesheetById(w http.ResponseWriter, r *http.Request) {
	// Get the employee ID from the URL parameters
	vars := mux.Vars(r)
	empIdStr := vars["id"]

	// Convert the employee ID to an integer
	empId, err := strconv.Atoi(empIdStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT 
		e.empId,
		e.firstName,
		e.lastName,
		e.email,
		el.login_time,
		el.logout_time,
		el.working_hours
	FROM 
		employee_logs el
	JOIN 
		emply e ON e.empId = el.empId
	WHERE 
		e.empId = ?`, empId)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []User
	for rows.Next() {
		var emp User
		if err := rows.Scan(&emp.EmpId, &emp.FirstName, &emp.LastName, &emp.Email, &emp.LoginTime, &emp.LogoutTime, &emp.WorkingHours); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved data: %v", employees)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getregisterdEmployee(w http.ResponseWriter, r *http.Request) {
	// Get the employee ID from the URL parameters
	vars := mux.Vars(r)
	empIdStr := vars["id"]

	// Convert the employee ID to an integer
	empId, err := strconv.Atoi(empIdStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT 
		empId,
		firstName,
		lastName,
		email
	FROM 
		timesheet 
		WHERE
		empId = ?`, empId)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []User
	for rows.Next() {
		var emp User
		if err := rows.Scan(&emp.EmpId, &emp.FirstName, &emp.LastName, &emp.Email); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved data: %v", employees)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// var store = sessions.NewCookieStore([]byte("secret-key"))

func loginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var login_time sql.NullString
	var logout_time sql.NullString

	err = db.QueryRow("SELECT login_time, logout_time FROM employee_logs WHERE empId = ? ORDER BY login_time DESC LIMIT 1", user.EmpId).Scan(&login_time, &logout_time)
	if err != nil {
		if err == sql.ErrNoRows {
			// No previous logins found, allow login
			// Proceed with the login process
		} else {
			log.Printf("Error fetching user data: %v", err)
			http.Error(w, "Error retrieving user data", http.StatusInternalServerError)
			return
		}
	}
	log.Printf("%v,%v", login_time, logout_time)

	// If there is a login_time and logout_time is null, prevent the new login
	if login_time.Valid && !logout_time.Valid {
		http.Error(w, "Cannot log in again without logging out", http.StatusForbidden)
		return
	}

	// Proceed with the login process
	// Your login logic here

	// Proceed with the login process
	// Your login logic here

	var hashedPassword string
	var empId int

	err = db.QueryRow("SELECT empId, password FROM timesheet WHERE email = ?", user.Email).Scan(&empId, &hashedPassword)
	if err != nil {
		log.Printf("Error fetching user data: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Log the system time and UTC time
	localTime := time.Now()
	utcTime := localTime.UTC()

	log.Printf("System Time: %v", localTime)
	log.Printf("UTC Time: %v", utcTime)

	// loginTime := time.Now() // Local time
	// utcTime := loginTime.UTC()

	_, err = db.Exec("INSERT INTO employee_logs (empId) VALUES (?)", empId)
	if err != nil {
		log.Printf("Error inserting into employee_logs: %v", err)
		http.Error(w, "Error tracking login", http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Login successful")
}

func loginEmployeeWithFace(w http.ResponseWriter, r *http.Request) {
	var inputData struct {
		FaceEmbeddings []float32 `json:"faceEmbeddings"` // Face embeddings coming from the frontend
	}

	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Fetch all face embeddings from the database
	rows, err := db.Query("SELECT empId, face_embeddings FROM timesheet")
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var matchedEmpId int
	threshold := 0.6 // Adjust based on required accuracy
	foundMatch := false

	for rows.Next() {
		var empId int
		var faceData []byte
		err := rows.Scan(&empId, &faceData)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		var storedEmbeddings []float32
		err = json.Unmarshal(faceData, &storedEmbeddings)
		if err != nil {
			log.Printf("Error unmarshalling face embeddings: %v", err)
			continue
		}

		if compareFaces(storedEmbeddings, inputData.FaceEmbeddings, threshold) {
			matchedEmpId = empId
			foundMatch = true
			break
		}
	}

	if foundMatch {
		// Check if the last login has been logged out
		var lastLogoutStr string
		err := db.QueryRow("SELECT logout_time FROM employee_logs WHERE empId = ? ORDER BY login_time DESC LIMIT 1", matchedEmpId).Scan(&lastLogoutStr)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error checking last logout: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var lastLogout time.Time
		if lastLogoutStr != "" {
			lastLogout, err = time.Parse("2006-01-02 15:04:05", lastLogoutStr)
			if err != nil {
				log.Printf("Error parsing logout time: %v", err)
				http.Error(w, "Error parsing logout time", http.StatusInternalServerError)
				return
			}

			if lastLogout.IsZero() {
				// Last session is still active, deny new login
				http.Error(w, "Previous session is still active", http.StatusConflict)
				return
			}
		}

		// Insert the login record
		_, err = db.Exec("INSERT INTO employee_logs (empId) VALUES (?)", matchedEmpId)
		if err != nil {
			log.Printf("Error inserting into employee_logs: %v", err)
			http.Error(w, "Error tracking login", http.StatusInternalServerError)
			return
		}

		// Successful login
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"empId": matchedEmpId})
	} else {
		// Face does not match
		http.Error(w, "Face not recognized", http.StatusUnauthorized)
	}
}

func compareFaces(stored, input []float32, threshold float64) bool {
	if len(stored) != len(input) {
		return false
	}
	sum := 0.0
	for i := range stored {
		diff := stored[i] - input[i]
		sum += float64(diff * diff)
	}
	return math.Sqrt(sum) <= threshold
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the employee ID from the URL parameters
	vars := mux.Vars(r)
	empIdStr := vars["id"]

	// Convert the employee ID to an integer
	empId, err := strconv.Atoi(empIdStr)
	if err != nil {
		log.Printf("Invalid employee ID: %v", err)
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Fetch user login time
	var loginTimeStr string
	err = db.QueryRow("SELECT login_time FROM employee_logs WHERE empId = ? AND logout_time IS NULL", empId).Scan(&loginTimeStr)
	if err != nil {
		log.Printf("Error retrieving login time: %v", err)
		http.Error(w, "Error retrieving login time", http.StatusInternalServerError)
		return
	}

	// Convert loginTimeStr to time.Time
	layout := "2006-01-02 15:04:05" // Adjust format based on your database datetime format
	loginTime, err := time.Parse(layout, loginTimeStr)
	if err != nil {
		log.Printf("Error parsing login time: %v", err)
		http.Error(w, "Error parsing login time", http.StatusInternalServerError)
		return
	}

	// Update logout time
	_, err = db.Exec(`UPDATE employee_logs SET logout_time = NOW() WHERE empId = ? AND logout_time IS NULL`, empId)
	if err != nil {
		log.Printf("Error updating employee_logs: %v", err)
		http.Error(w, "Error updating logout time", http.StatusInternalServerError)
		return
	}

	// Fetch logout time
	var logoutTimeStr string
	err = db.QueryRow("SELECT logout_time FROM employee_logs WHERE empId = ? AND logout_time IS NOT NULL ORDER BY logout_time DESC LIMIT 1", empId).Scan(&logoutTimeStr)
	if err != nil {
		log.Printf("Error retrieving logout time: %v", err)
		http.Error(w, "Error retrieving logout time", http.StatusInternalServerError)
		return
	}

	// Convert logoutTimeStr to time.Time
	logoutTime, err := time.Parse(layout, logoutTimeStr)
	if err != nil {
		log.Printf("Error parsing logout time: %v", err)
		http.Error(w, "Error parsing logout time", http.StatusInternalServerError)
		return
	}

	// Calculate working hours
	duration := logoutTime.Sub(loginTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	workingHours := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	// Update working hours in the database
	_, err = db.Exec(`UPDATE employee_logs SET working_hours = ? WHERE empId = ? AND logout_time = ?`, workingHours, empId, logoutTimeStr)
	if err != nil {
		log.Printf("Error updating working_hours: %v", err)
		http.Error(w, "Error updating working hours", http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Logout successful")
}

func logoutWithFaceHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request
	var inputData struct {
		FaceEmbeddings []float32 `json:"faceEmbeddings"` // Face embeddings coming from the frontend
	}

	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Fetch all face embeddings from the database
	rows, err := db.Query("SELECT empId, face_embeddings FROM timesheet")
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var matchedEmpId int
	threshold := 0.6 // Adjust based on required accuracy
	foundMatch := false

	for rows.Next() {
		var empId int
		var faceData []byte
		err := rows.Scan(&empId, &faceData)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		log.Printf("Raw face embeddings data for empId %d: %s", empId, string(faceData))

		// Skip employees without face data
		if len(faceData) == 0 {
			log.Printf("No face data for empId %d, skipping", empId)
			continue
		}

		var storedEmbeddings []float32
		err = json.Unmarshal(faceData, &storedEmbeddings)
		log.Printf("%v", storedEmbeddings)
		if err != nil {
			log.Printf("Error unmarshalling face embeddings: %v", err)
			continue
		}

		if compareFaces(storedEmbeddings, inputData.FaceEmbeddings, threshold) {
			matchedEmpId = empId
			log.Printf("%v", matchedEmpId)
			foundMatch = true
			break
		}
	}

	if !foundMatch {
		http.Error(w, "Face not recognized", http.StatusUnauthorized)
		return
	}

	// Fetch the login time for this employee
	var loginTimeStr string
	err = db.QueryRow("SELECT login_time FROM employee_logs WHERE empId = ? AND logout_time IS NULL", matchedEmpId).Scan(&loginTimeStr)
	if err != nil {
		log.Printf("Error retrieving login time: %v", err)
		http.Error(w, "Error retrieving login time", http.StatusInternalServerError)
		return
	}

	// Convert loginTimeStr to time.Time
	layout := "2006-01-02 15:04:05" // Adjust format based on your database datetime format
	loginTime, err := time.Parse(layout, loginTimeStr)
	if err != nil {
		log.Printf("Error parsing login time: %v", err)
		http.Error(w, "Error parsing login time", http.StatusInternalServerError)
		return
	}

	// Update logout time
	_, err = db.Exec(`UPDATE employee_logs SET logout_time = NOW() WHERE empId = ? AND logout_time IS NULL`, matchedEmpId)
	if err != nil {
		log.Printf("Error updating employee_logs: %v", err)
		http.Error(w, "Error updating logout time", http.StatusInternalServerError)
		return
	}

	// Fetch the updated logout time
	var logoutTimeStr string
	err = db.QueryRow("SELECT logout_time FROM employee_logs WHERE empId = ? AND logout_time IS NOT NULL ORDER BY logout_time DESC LIMIT 1", matchedEmpId).Scan(&logoutTimeStr)
	if err != nil {
		log.Printf("Error retrieving logout time: %v", err)
		http.Error(w, "Error retrieving logout time", http.StatusInternalServerError)
		return
	}

	// Convert logoutTimeStr to time.Time
	logoutTime, err := time.Parse(layout, logoutTimeStr)
	if err != nil {
		log.Printf("Error parsing logout time: %v", err)
		http.Error(w, "Error parsing logout time", http.StatusInternalServerError)
		return
	}

	// Calculate working hours
	duration := logoutTime.Sub(loginTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	workingHours := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	// Update working hours in the database
	_, err = db.Exec(`UPDATE employee_logs SET working_hours = ? WHERE empId = ? AND logout_time = ?`, workingHours, matchedEmpId, logoutTimeStr)
	if err != nil {
		log.Printf("Error updating working_hours: %v", err)
		http.Error(w, "Error updating working hours", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "Logout successful for employee ID %d", matchedEmpId)
	json.NewEncoder(w).Encode(map[string]int{"empId": matchedEmpId})

}

func handleLoginLogout(w http.ResponseWriter, r *http.Request) {
	var inputData struct {
		FaceEmbeddings []float32 `json:"faceEmbeddings"` // Face embeddings coming from the frontend
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Open a connection to the database
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Fetch face embeddings from the database
	rows, err := db.Query("SELECT empId, face_embeddings FROM timesheet")
	if err != nil {
		log.Printf("Failed to query database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var matchedEmpId int
	threshold := 0.6 // Adjust based on required accuracy
	foundMatch := false

	// Compare incoming face embeddings with stored ones
	for rows.Next() {
		var empId int
		var faceData []byte
		err := rows.Scan(&empId, &faceData)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		if len(faceData) == 0 {
			log.Printf("No face data for empId %d, skipping", empId)
			continue
		}

		var storedEmbeddings []float32
		err = json.Unmarshal(faceData, &storedEmbeddings)
		if err != nil {
			log.Printf("Error unmarshalling face embeddings: %v", err)
			continue
		}

		if compareFaces(storedEmbeddings, inputData.FaceEmbeddings, threshold) {
			matchedEmpId = empId
			log.Printf("Matched empID is %d", matchedEmpId)
			foundMatch = true
			break
		}
	}

	if !foundMatch {
		http.Error(w, "Face not recognized", http.StatusUnauthorized)
		return
	}

	// Check the most recent session status
	var loginTimeStr sql.NullString
	var logoutTimeStr sql.NullString
	err = db.QueryRow("SELECT login_time, logout_time FROM employee_logs WHERE empId = ? ORDER BY login_time DESC LIMIT 1", matchedEmpId).Scan(&loginTimeStr, &logoutTimeStr)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking session status: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err == sql.ErrNoRows || (loginTimeStr.Valid && logoutTimeStr.Valid) {
		// If no previous records or last record has both login and logout times, proceed with login
		if err == sql.ErrNoRows || (loginTimeStr.Valid && logoutTimeStr.Valid) {
			_, err = db.Exec("INSERT INTO employee_logs (empId) VALUES (?)", matchedEmpId)
			if err != nil {
				log.Printf("Error inserting into employee_logs for login: %v", err)
				http.Error(w, "Error tracking login", http.StatusInternalServerError)
				return
			}

			// Respond with a success message
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "empId": fmt.Sprintf("%d", matchedEmpId), "Login": time.Now().Format("2006-01-02 15:04:05")})
			return
		}
	}

	if loginTimeStr.Valid && !logoutTimeStr.Valid {
		// If user is currently logged in, log them out
		_, err = db.Exec(`UPDATE employee_logs SET logout_time = NOW() WHERE empId = ? AND logout_time IS NULL`, matchedEmpId)
		if err != nil {
			log.Printf("Error updating employee_logs for logout: %v", err)
			http.Error(w, "Error updating logout time", http.StatusInternalServerError)
			return
		}

		// Calculate working hours
		layout := "2006-01-02 15:04:05" // Adjust format based on your database datetime format
		loginTime, err := time.Parse(layout, loginTimeStr.String)
		if err != nil {
			log.Printf("Error parsing login time: %v", err)
			http.Error(w, "Error parsing login time", http.StatusInternalServerError)
			return
		}

		var logoutTimeStrr string
		err = db.QueryRow("SELECT logout_time FROM employee_logs WHERE empId = ? AND logout_time IS NOT NULL ORDER BY logout_time DESC LIMIT 1", matchedEmpId).Scan(&logoutTimeStrr)
		if err != nil {
			log.Printf("Error retrieving logout time: %v", err)
			http.Error(w, "Error retrieving logout time", http.StatusInternalServerError)
			return
		}

		logoutTime, err := time.Parse(layout, logoutTimeStrr)
		if err != nil {
			log.Printf("Error parsing logout time: %v", err)
			http.Error(w, "Error parsing logout time", http.StatusInternalServerError)
			return
		}

		duration := logoutTime.Sub(loginTime)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		seconds := int(duration.Seconds()) % 60
		workingHours := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

		// Update working hours in the database
		_, err = db.Exec(`UPDATE employee_logs SET working_hours = ? WHERE empId = ? AND logout_time = ?`, workingHours, matchedEmpId, logoutTime.Format(layout))
		if err != nil {
			log.Printf("Error updating working_hours: %v", err)
			http.Error(w, "Error updating working hours", http.StatusInternalServerError)
			return
		}

		// Respond with a success message
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful", "empId": fmt.Sprintf("%d", matchedEmpId), "Login": time.Now().Format("2006-01-02 15:04:05")})
		return
	}

	// If none of the above conditions are met, respond with an error
	http.Error(w, "Invalid session state", http.StatusBadRequest)
}

func checkLoginLogoutStatus(w http.ResponseWriter, r *http.Request) {
	var inputData struct {
		FaceEmbeddings []float32 `json:"faceEmbeddings"`
	}

	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT empId, face_embeddings FROM timesheet")
	if err != nil {
		log.Printf("Failed to query database: %v", err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var matchedEmpId int
	threshold := 0.6
	foundMatch := false

	for rows.Next() {
		var empId int
		var faceData []byte
		err := rows.Scan(&empId, &faceData)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		if len(faceData) == 0 {
			log.Printf("No face data for empId %d, skipping", empId)
			continue
		}

		var storedEmbeddings []float32
		err = json.Unmarshal(faceData, &storedEmbeddings)
		if err != nil {
			log.Printf("Error unmarshalling face embeddings: %v", err)
			continue
		}

		if compareFaces(storedEmbeddings, inputData.FaceEmbeddings, threshold) {
			matchedEmpId = empId
			log.Printf("matched empID IS %v", matchedEmpId)
			foundMatch = true
			break
		}
	}

	if !foundMatch {
		http.Error(w, "Face not recognized", http.StatusUnauthorized)
		return
	}

	var loginTimeStr sql.NullString
	var logoutTimeStr sql.NullString
	err = db.QueryRow("SELECT login_time, logout_time FROM employee_logs WHERE empId = ? ORDER BY login_time DESC LIMIT 1", matchedEmpId).Scan(&loginTimeStr, &logoutTimeStr)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking session status: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Case 1: No records or both login and logout times are present in the last record
	if err == sql.ErrNoRows || (loginTimeStr.Valid && logoutTimeStr.Valid) {
		json.NewEncoder(w).Encode(map[string]string{"message": "Login needed"})
		return
	}

	// Case 2: Logged in but not logged out
	if loginTimeStr.Valid && !logoutTimeStr.Valid {
		json.NewEncoder(w).Encode(map[string]string{"message": "Logout needed"})
		return
	}

	// Case 3: Both login and logout times are present (fully logged out)
	if loginTimeStr.Valid && logoutTimeStr.Valid {
		json.NewEncoder(w).Encode(map[string]string{"message": "Login needed"})
		return
	}
}

// func compareFaces(stored, input []float32, threshold float64) bool {
// 	if len(stored) != len(input) {
// 		return false
// 	}
// 	sum := 0.0
// 	for i := range stored {
// 		diff := stored[i] - input[i]
// 		sum += float64(diff * diff)
// 	}
// 	return math.Sqrt(sum) <= threshold
// }
//main

// Admin represents the admin user
type Admin struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims represents the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// var db *sql.DB
var jwtKey = []byte("your_secret_key")

func createAdmin(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON payload
	var user Admin
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the empId already exists
	var existingUserId int
	checkQuery := "SELECT id FROM admin WHERE id = ?"
	err = db.QueryRow(checkQuery, user.ID).Scan(&existingUserId)
	if err == nil {
		// If no error is returned, it means the empId exists
		http.Error(w, "User with this Admin already exists", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		// If there's an error other than no rows found, log it
		log.Printf("Error checking existing user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// // Convert face embeddings to JSON
	// faceEmbeddingsJSON, err := json.Marshal(user.FaceEmbeddings)
	// if err != nil {
	// 	log.Printf("Error marshalling face embeddings: %v", err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// Insert into database
	query := "INSERT INTO admin (username, password_hash) VALUES (?, ?)"
	_, err = db.Exec(query, user.Username, string(hashedPassword))
	if err != nil {
		log.Printf("Error saving user: %v", err)
		http.Error(w, "Failed to register: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func adminloginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:welcome123@tcp(127.0.0.1:3306)/employee")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var admin Admin
	err = db.QueryRow("SELECT id, username, password_hash FROM admin WHERE username = ?", req.Username).Scan(&admin.ID, &admin.Username, &admin.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Password provided: %s", req.Password)
	log.Printf("Password hash stored: %s", admin.PasswordHash)

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: req.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return JWT token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// func logoutHandler(w http.ResponseWriter, r *http.Request) {
// 	// For JWT-based authentication, logout is typically handled on the client side
// 	// by removing the token from storage. However, you could invalidate tokens on
// 	// the server side by using a token blacklist or short-lived tokens.

// 	w.Write([]byte("Logout successful"))
// }

// JWT Middleware function
func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Store claims in context
		ctx := context.WithValue(r.Context(), "claims", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/employee", dataHandler).Methods("GET")
	// r.HandleFunc("/addemployee", addEmployeeHandler).Methods("POST")
	r.Handle("/addemployee", jwtMiddleware(http.HandlerFunc(addEmployeeHandler))).Methods("POST")

	// r.HandleFunc("/getemployee/employee/{id}", getEmployeeByID).Methods("GET")
	//
	// r.HandleFunc("/getbyid/prevcompany/{id}", prevcompanybyid).Methods("GET")
	r.HandleFunc("/getbyid/emply/{id}", emplybyid).Methods("GET")
	r.HandleFunc("/getbyid/emplywithprevcompany/{id}", employeeWithCompanyById).Methods("GET")
	//
	r.Handle("/delete/employee/{id}", jwtMiddleware(http.HandlerFunc(deleteEmployeeHandler))).Methods("DELETE")
	r.Handle("/update/employee/{id}", jwtMiddleware(http.HandlerFunc(updateEmployeeHandler))).Methods("PUT")
	r.HandleFunc("/getrolesbydepartment/{id}", getRolesByDepartment).Methods("GET")
	r.HandleFunc("/departments", getDepartments).Methods("GET")
	r.HandleFunc("/managers", getManagers).Methods("GET")
	r.HandleFunc("/manager/{id}", getManagerByDepartment).Methods("GET")
	r.Handle("/assignemployee", jwtMiddleware(http.HandlerFunc(assignEmployee))).Methods("POST")
	r.Handle("/uploaddocuments", jwtMiddleware(http.HandlerFunc(uploadDocumentsHandler))).Methods("POST")
	r.Handle("/personaldetails", jwtMiddleware(http.HandlerFunc(personaldetails))).Methods("POST")
	r.Handle("/handlePersonalDetailsAndDocuments", jwtMiddleware(http.HandlerFunc(handlePersonalDetailsAndDocuments))).Methods("POST")

	r.HandleFunc("/hierarchy", getHierarchyData).Methods("GET")
	// Register the route with the router
	r.HandleFunc("/hierarchy/{id}", getHierarchyDataById).Methods("GET")
	r.HandleFunc("/departmentbyid/{id}", getDepartmentsById).Methods("GET")
	r.HandleFunc("/getDocuments/{id}", getDocuments).Methods("GET")
	r.Handle("/handleUpdatePersonalDetailsAndDocuments", jwtMiddleware(http.HandlerFunc(handleUpdatePersonalDetailsAndDocuments))).Methods("PUT")
	r.HandleFunc("/setManager/{id}", setManager).Methods("GET")
	r.HandleFunc("/getEmployeeAsManager/{id}", getEmployeeAsManager).Methods("GET")
	r.HandleFunc("/getRoleById/{id}", getRoleById).Methods("GET")
	r.HandleFunc("/getAssignData", getAssignData).Methods("GET")
	r.Handle("/register", jwtMiddleware(http.HandlerFunc(registerHandler))).Methods("POST")
	r.Handle("/login", jwtMiddleware(http.HandlerFunc(loginHandler))).Methods("POST")
	r.Handle("/loginwithface", jwtMiddleware(http.HandlerFunc(loginEmployeeWithFace))).Methods("POST")

	r.Handle("/logout/{id}", jwtMiddleware(http.HandlerFunc(logoutHandler))).Methods("POST")
	r.Handle("/logoutwithface", jwtMiddleware(http.HandlerFunc(logoutWithFaceHandler))).Methods("POST")

	r.HandleFunc("/gettimesheet", getTimesheet).Methods("GET")
	r.HandleFunc("/gettimesheetbyid/{id}", getTimesheetById).Methods("GET")
	r.HandleFunc("/getregisterdemployee/{id}", getregisterdEmployee).Methods("GET")
	r.Handle("/admin-register", jwtMiddleware(http.HandlerFunc(createAdmin))).Methods("POST")
	r.HandleFunc("/admin-login", adminloginHandler).Methods("POST")
	r.Handle("/facelogin", jwtMiddleware(http.HandlerFunc(handleLoginLogout))).Methods("POST")
	r.Handle("/checkloginlogoutstatus", jwtMiddleware(http.HandlerFunc(checkLoginLogoutStatus))).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
