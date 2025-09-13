package sqlconnect

import (
	"chickenTrade/API/internal/models"
	"chickenTrade/API/pkg/utils"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func GetProductsDBHandler(products []models.Product, r *http.Request) ([]models.Product, int, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "Couldn't connect to db")
	}

	defer db.Close()

	query := `SELECT id, breed, age_week, price_per_hen, is_vaccinated, collection_id FROM Product WHERE 1=1`
	rows, err := db.Query(query)

	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "Database query error")
	}

	defer rows.Close()

	productsList := make([]models.Product, 0)
	for rows.Next() {
		var prodcut models.Product

		err := rows.Scan(&prodcut.ID, &prodcut.Breed, &prodcut.AgeWeek, &prodcut.PricePerHen, &prodcut.IsVaccinated, &prodcut.CollectionID)
		if err != nil {
			return nil, 0, utils.ErrorHandler(err, "Error scanning databse results")
		}
		productsList = append(productsList, prodcut)
		fmt.Println("Products: ", productsList)
	}

	var totalProduct int
	err = db.QueryRow("SELECT COUNT(*) FROM Product").Scan(&totalProduct)
	if err != nil {
		utils.ErrorHandler(err, "")
		totalProduct = 0
	}
	return productsList, totalProduct, nil
}

func AddProductsDBHandler(newProducts []models.Product) ([]models.Product, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Couldn't connect to db")
	}

	defer db.Close()

	stmt, err := db.Prepare(utils.GenerateInsertQuery("Product", models.Product{}))

	if err != nil {
		return nil, utils.ErrorHandler(err, "Database preparation failed")
	}

	defer stmt.Close()

	addProducts := make([]models.Product, len(newProducts))

	for i, product := range newProducts {
		values := utils.GetStructValues(product)
		_, err := stmt.Exec(values...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid request")
		}

		// id, err := res.LastInsertId()
		// if err != nil {
		// 	utils.ErrorHandler(err, "Couldnt fetch the ID")
		// }

		product.ID = uuid.New().String()
		addProducts[i] = product
	}

	return addProducts, nil
}
