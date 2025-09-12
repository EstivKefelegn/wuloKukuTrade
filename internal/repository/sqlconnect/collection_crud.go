package sqlconnect

import (
	"chickenTrade/API/internal/models"
	"chickenTrade/API/pkg/utils"
	"fmt"
	"log"
	"net/http"
)

func GetCollectionsDBHandler(collections []models.Collection, r *http.Request) ([]models.Collection, int, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, 0, utils.ErrorHandler(err, "Couldn't connect to the db")
	}

	defer db.Close()

	rows, err := db.Query(`SELECT id, title, featured_product_id FROM Collection WHERE 1=1`)
	if err != nil {
		log.Println("Query error:", err)
		return nil, 0, utils.ErrorHandler(err, "Database query error")
	}
	defer rows.Close()

	collectionModels := make([]models.Collection, 0)
	for rows.Next() {
		var collection models.Collection
		err := rows.Scan(&collection.ID, &collection.Title, &collection.FeaturedProductID)
		if err != nil {
			return nil, 0, utils.ErrorHandler(err, "Error scanning databse results")
		}
		collectionModels = append(collectionModels, collection)
		fmt.Println("Collections:", collectionModels)
	}

	var totalCollection int
	err = db.QueryRow("SELECT COUNT(*) FROM Collection").Scan(&totalCollection)
	if err != nil {
		utils.ErrorHandler(err, "Error scanning databse results")
		totalCollection = 0
	}

	return collectionModels, totalCollection, nil
}

func AddCollectionDBHandler(newCollections []models.Collection) ([]models.Collection, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Couldn't connet to db")
	}

	defer db.Close()

	stmt, err := db.Prepare(utils.GenerateInsertQuery("Collection", models.Collection{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database preparation failed")
	}

	defer stmt.Close()

	addCollection := make([]models.Collection, len(newCollections))
	for i, collection := range newCollections {
		values := utils.GetStructValues(collection)
		fmt.Println("Vals: ", values)
		res, err := stmt.Exec(values...)

		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid request")
		}

		id, err := res.LastInsertId()
		if err != nil {
			utils.ErrorHandler(err, "Couldnt fetch the ID")
		}

		collection.ID = int(id)
		addCollection[i] = collection
	}

	return addCollection, nil
}
