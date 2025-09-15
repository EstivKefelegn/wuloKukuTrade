package sqlconnect

import (
	"chickenTrade/API/internal/models"
	"chickenTrade/API/pkg/utils"
	"net/http"

	"github.com/google/uuid"
)

func GetUsersDBHandler(users []models.Users, r *http.Request) ([]models.Users, error) {

	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "couldn't connect to db")
	}

	defer db.Close()

	rows, err := db.Query(`SELECT id, first_name, last_name, email, username, phone, user_created_at, in_active_status, role FROM Users WHERE 1=1`)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Couldn't retrive the values")
	}

	defer rows.Close()

	usersList := make([]models.Users, 0)
	for rows.Next() {
		var user models.Users
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.UserName, &user.Phone, &user.UserCreatedAt, &user.InactiveStatus, &user.Role)

		if err != nil {
			return nil, utils.ErrorHandler(err, "unknown error")
		}

		usersList = append(usersList, user)
	}

	return usersList, nil
}

func RegisterUsersDBHandler(newUser models.Users) (models.Users, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Users{}, utils.ErrorHandler(err, "couldn't connect to db")
	}

	defer db.Close()

	stmt, err := db.Prepare(utils.GenerateInsertQuery("Users", models.Users{}))
	if err != nil {
		return models.Users{}, utils.ErrorHandler(err, "Internal erorr")
	}

	defer stmt.Close()
	newUser.Password, err = utils.HashPassword(newUser.Password)
	if err != nil {
		return models.Users{}, utils.ErrorHandler(err, "couldn't hash the password")
	}

	newUser.ID = uuid.New().String()
	user := utils.GetStructValues(newUser)

	res, err := stmt.Exec(user...)
	if err != nil {
		return models.Users{}, utils.ErrorHandler(err, "couldn't excute the queries")
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return models.Users{}, utils.ErrorHandler(nil, "User not registered")
	}

	return newUser, nil
}
