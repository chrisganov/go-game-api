package data

import (
	"database/sql"
	"errors"
	"fmt"

	validator "go_game_api.com/internal/validators"
)

type userRole string

const (
	UserRole      userRole = "USER"
	SuperUserRole userRole = "SUPERUSER"
	AdminRole     userRole = "ADMIN"
)

type User struct {
	Id        int      `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	CreatedAt string   `json:"createdAt"`
	Passhash  string   `json:"-"`
	Role      userRole `json:"-"`
	UpdatedAt string   `json:"-"`
}

type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserModel struct {
	DB *sql.DB
}

func ValidateUserInput(v *validator.Validator, userInput *UserInput) {
	v.Check(userInput.Username != "", "username", "Username cannot be empty")
	v.Check(validator.UsernameRegex.MatchString(userInput.Username), "username", "Username does not match pattern")

	v.Check(userInput.Email != "", "email", "Email cannot be empty")
	v.Check(validator.EmailRegex.MatchString(userInput.Email), "email", "Email does not match pattern")

	v.Check(len(userInput.Password) >= 8, "password", "Password cannot be less than 8 characters")
	v.Check(validator.LowerCaseLettersRegex.MatchString(userInput.Password), "password", "Password must include lower case")
	v.Check(validator.UpperCaseLettersRegex.MatchString(userInput.Password), "password", "Password must include upper case")
	v.Check(validator.DigitsRegex.MatchString(userInput.Password), "password", "Password must include a number")
	v.Check(validator.SpecialCharactersRegex.MatchString(userInput.Password), "password", "Password must include a special character")

}

// func ValidateUser(v *validator.Validator, user *User) {
// 	// TODO add more validatons
// v.Check(user.Username != "", "username", "Username cannot be empty")
// v.Check(validator.UsernameRegex.MatchString(user.Username), "username", "Username does not match pattern")

// v.Check(user.Email != "", "email", "Email cannot be empty")
// v.Check(validator.EmailRegex.MatchString(user.Email), "email", "Email does not match pattern")

// v.Check(len(user.Passhash) < )

// }

const userColumns = "id, username, email, passhash, role, created_at, updated_at"

func (m UserModel) Insert(user *User) error {
	args := []interface{}{user.Username, user.Email, user.Passhash, user.Role}

	query := fmt.Sprintf(`
		INSERT INTO users (username, email, passhash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING %s
	`, userColumns)

	err := m.DB.QueryRow(query, args...).Scan(&user.Id, &user.Username, &user.Email, &user.Passhash, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (m UserModel) GetById(id int) (*User, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM users
		WHERE id = $1
`, userColumns)

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	var user User

	err := m.DB.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Email, &user.Passhash, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) GetAll() ([]User, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM users
	`, userColumns)

	rows, err := m.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Passhash, &user.Role, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			// TODO
			continue
		}

		users = append(users, user)
	}

	return users, nil
}
