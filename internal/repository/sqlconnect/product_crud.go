package sqlconnect

import (
	"chickenTrade/API/internal/models"
	"chickenTrade/API/pkg/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"

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

	}

	var totalProduct int
	err = db.QueryRow("SELECT COUNT(*) FROM Product").Scan(&totalProduct)
	if err != nil {
		utils.ErrorHandler(err, "")
		totalProduct = 0
	}
	return productsList, totalProduct, nil
}

func GetOneProductDBHandler(id string) (models.Product, error) {
	db, err := ConnectDB()
	if err != nil {
		utils.ErrorHandler(err, "Couldn't connect to db")
	}

	defer db.Close()
	var product models.Product
	var lastUpdate []uint8 // Scan into byte slice first

	err = db.QueryRow(`SELECT id, breed, age_week, price_per_hen, is_vaccinated, 
                      NULLIF(last_update, '0000-00-00 00:00:00'), collection_id 
                      FROM Product WHERE id = ?`, id).Scan(
		&product.ID, &product.Breed, &product.AgeWeek, &product.PricePerHen,
		&product.IsVaccinated, &lastUpdate, &product.CollectionID)

	if err == sql.ErrNoRows {
		return models.Product{}, utils.ErrorHandler(err, "No rows found")
	} else if err != nil {
		return models.Product{}, utils.ErrorHandler(err, "Database query error")
	}
	return product, nil
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
		product.ID = uuid.New().String()
		values := utils.GetStructValues(product)
		_, err := stmt.Exec(values...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid request")
		}

		// id, err := res.LastInsertId()
		// if err != nil {
		// 	utils.ErrorHandler(err, "Couldnt fetch the ID")
		// }

		// product.ID = uuid.New().String()
		addProducts[i] = product
	}

	return addProducts, nil
}

func UpdateProductDBHandler(id string, updatedProduct models.Product) (models.Product, error) {

	db, err := ConnectDB()
	if err != nil {
		return models.Product{}, utils.ErrorHandler(err, "Couldn't connect to db")
	}

	defer db.Close()

	var existingProduct models.Product
	row := db.QueryRow(`SELECT id, breed, age_week, price_per_hen, is_vaccinated, collection_id FROM Product WHERE id = ?`, id)
	err = row.Scan(&existingProduct.ID, &existingProduct.Breed, &existingProduct.AgeWeek, &existingProduct.PricePerHen, &existingProduct.IsVaccinated, &existingProduct.CollectionID)

	if err == sql.ErrNoRows {
		return models.Product{}, utils.ErrorHandler(err, "No rows found with this id")
	} else if err != nil {
		return models.Product{}, utils.ErrorHandler(err, "Couldnt fetch the rows")
	}

	updatedProduct.ID = existingProduct.ID
	res, err := db.Exec("UPDATE Product SET breed = ?, age_week = ?, price_per_hen = ?, is_vaccinated = ?, collection_id = ? WHERE id = ?",
		updatedProduct.Breed, updatedProduct.AgeWeek, updatedProduct.PricePerHen, updatedProduct.IsVaccinated, updatedProduct.CollectionID, updatedProduct.ID)
	if err != nil {
		return models.Product{}, utils.ErrorHandler(err, "Couldn't update the rows")
	}

	rowsAffected, _ := res.RowsAffected()
	fmt.Println(rowsAffected, "row's are affected")

	return updatedProduct, nil
}

func PatchProductDBHandler(id string, updates map[string]interface{}) (models.Product, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Product{}, utils.ErrorHandler(err, "Couldn't connect to db")
	}

	defer db.Close()

	var existingProduct models.Product
	row := db.QueryRow(`SELECT id, breed, age_week, price_per_hen, is_vaccinated, collection_id FROM Product WHERE id = ?`, id)
	err = row.Scan(&existingProduct.ID, &existingProduct.Breed, &existingProduct.AgeWeek, &existingProduct.PricePerHen, &existingProduct.IsVaccinated, &existingProduct.CollectionID)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Product{}, utils.ErrorHandler(err, "no rows found witht this id")
		}
		return models.Product{}, utils.ErrorHandler(err, "couldn't scan the existing rows")
	}

	productVal := reflect.ValueOf(&existingProduct).Elem()
	productType := productVal.Type()

	for k, v := range updates {
		for i := 0; i < productVal.NumField(); i++ {
			field := productType.Field(i)

			if field.Tag.Get("json") == k || field.Tag.Get("json") == k+",omitempty" {
				if productVal.Field(i).CanSet() {
					fieldVal := productVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(productVal.Field(i).Type()))
				}
			}
		}
	}

	res, err := db.Exec(`UPDATE Product SET breed = ?, age_week = ?, price_per_hen = ?, is_vaccinated = ?, collection_id = ? WHERE id = ?`, &existingProduct.Breed, &existingProduct.AgeWeek, &existingProduct.PricePerHen, &existingProduct.IsVaccinated, &existingProduct.CollectionID, &existingProduct.ID)
	if err != nil {
		return models.Product{}, utils.ErrorHandler(err, "couldn't update the request")
	}
	rowsAffected, _ := res.RowsAffected()
	fmt.Println("Rows updated: ", rowsAffected)
	return existingProduct, nil
}

func PatchProductsDBHandler(w http.ResponseWriter, updates []map[string]interface{}) ([]models.Product, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "couldn't connect to db")
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error starting the transaction")
	}

	var updatedProducts []models.Product

	for _, update := range updates {
		id, ok := update["id"]
		if !ok {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Invalid id")
		}

		var existingProduct models.Product
		err = db.QueryRow(`SELECT id, breed, age_week, price_per_hen, is_vaccinated, collection_id FROM Product WHERE id = ?`, id).Scan(
			&existingProduct.ID, &existingProduct.Breed, &existingProduct.AgeWeek, &existingProduct.PricePerHen, &existingProduct.IsVaccinated, &existingProduct.CollectionID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, utils.ErrorHandler(err, "couldn't found a row with this id")
			}
			return nil, utils.ErrorHandler(err, "Error retriving the products")
		}

		productVal := reflect.ValueOf(&existingProduct).Elem()
		productType := productVal.Type()

		for k, v := range update {
			for i := 0; i < productVal.NumField(); i++ {
				field := productType.Field(i)
				if field.Tag.Get("json") == k || field.Tag.Get("json") == k+",omitempty" {
					if productVal.Field(i).CanSet() {
						fieldval := productVal.Field(i)
						fieldval.Set(reflect.ValueOf(v).Convert(fieldval.Type()))
					}
					break
				}
			}
		}
		res, err := tx.Exec(`UPDATE Product SET breed = ?, age_week = ?, price_per_hen = ?, is_vaccinated = ?, collection_id = ? WHERE id = ?`, &existingProduct.Breed, &existingProduct.AgeWeek, &existingProduct.PricePerHen, &existingProduct.IsVaccinated, &existingProduct.CollectionID, &existingProduct.ID)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Couldn't update the rows")
		}

		affectedRows, err := res.RowsAffected()
		if err != nil {
			return nil, utils.ErrorHandler(err, "There is no affectedRows")
		}

		fmt.Println(affectedRows, "Affected Rows")

		updatedProducts = append(updatedProducts, existingProduct)
	}
	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error commiting the transaction")
	}

	return updatedProducts, nil
}

func DeleteProductDBHandler(w http.ResponseWriter, id string) (string, error) {
	db, err := ConnectDB()
	if err != nil {
		return "", utils.ErrorHandler(err, "couldn't connect to db")
	}

	defer db.Close()

	res, err := db.Exec("DELETE FROM Product WHERE id = ?", id)
	if err != nil {
		return "", utils.ErrorHandler(err, "couldn't delete the requested value")
	}

	affectedRow, err := res.RowsAffected()
	if err != nil {
		return "", utils.ErrorHandler(err, "there ia no affected row")
	}
	fmt.Println(affectedRow, "affected rows")

	return id, nil
}

func DeleteProductsDBHandler(w http.ResponseWriter, ids []string) ([]string, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "couldn't connect to db")
	}

	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "Invalid transaction")
	}

	stmt, err := tx.Prepare("DELETE FROM Product WHERE id = ?")
	if err != nil {
		utils.ErrorHandler(err, "Coldn't delete on this transaction")
	}

	defer stmt.Close()

	var deletedIDs []string
	for _, id := range ids {
		res, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			log.Println("Invlid transaction")
			return nil, utils.ErrorHandler(err, "Invalid transaction")
		}

		affectedRow, err := res.RowsAffected()
		if err != nil {
			return nil, utils.ErrorHandler(err, "no affected rows")
		}

		if affectedRow > 0 {
			deletedIDs = append(deletedIDs, id)
		}

		if affectedRow < 1 {
			return nil, utils.ErrorHandler(err, "Invalid transaction")
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println("err")
		return nil, utils.ErrorHandler(err, "Invalid transaction")
	}

	if len(deletedIDs) < 1 {
		return nil, utils.ErrorHandler(err, "IDs do not exists")
	}

	return deletedIDs, nil
}
