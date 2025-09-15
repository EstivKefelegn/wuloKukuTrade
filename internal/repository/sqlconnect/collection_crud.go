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
	}

	var totalCollection int
	err = db.QueryRow("SELECT COUNT(*) FROM Collection").Scan(&totalCollection)
	if err != nil {
		utils.ErrorHandler(err, "Error scanning databse results")
		totalCollection = 0
	}

	return collectionModels, totalCollection, nil
}

func GetCollectionDBHandler(id string) (models.Collection, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Collection{}, utils.ErrorHandler(err, "Couldn't connect to db")
	}

	defer db.Close()

	var collection models.Collection
	err = db.QueryRow(`SELECT id, title, featured_product_id FROM Collection WHERE id = ?`, id).Scan(&collection.ID, &collection.Title, &collection.FeaturedProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Collection{}, utils.ErrorHandler(err, "there is no collection with this id")
		}
		return models.Collection{}, utils.ErrorHandler(err, "could't query")
	}

	return collection, nil
}

func AddCollectionDBHandler(newCollections []models.Collection) ([]models.Collection, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Couldn't connect to db")
	}

	defer db.Close()

	stmt, err := db.Prepare(utils.GenerateInsertQuery("Collection", models.Collection{}))
	fmt.Println("stmt:", stmt)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database preparation failed")
	}

	defer stmt.Close()

	addCollection := make([]models.Collection, len(newCollections))
	for i, collection := range newCollections {
		collection.ID = uuid.New().String()
		values := utils.GetStructValues(collection)

		_, err := stmt.Exec(values...)

		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid request")
		}

		// id, err := res.LastInsertId()
		// if err != nil {
		// 	utils.ErrorHandler(err, "Couldnt fetch the ID")
		// }

		addCollection[i] = collection
	}

	return addCollection, nil
}

func UpdateCollectionDBHandler(id string, update models.Collection) (models.Collection, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Collection{}, utils.ErrorHandler(err, "Couldn't connect db")
	}

	defer db.Close()

	var existingCollection models.Collection
	err = db.QueryRow(`SELECT id, title, featured_product_id FROM Collection WHERE id = ?`, id).Scan(&existingCollection.ID, &existingCollection.Title, &existingCollection.FeaturedProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Collection{}, utils.ErrorHandler(err, "NO rows found with this id")
		}

		return models.Collection{}, utils.ErrorHandler(err, "couldn't retrive the data")
	}

	fmt.Println("Existing ID:", existingCollection.ID)

	update.ID = existingCollection.ID
	res, err := db.Exec("UPDATE Collection SET title = ?, featured_product_id = ? WHERE id = ?", update.Title, update.FeaturedProductID, update.ID)

	if err != nil {
		return models.Collection{}, utils.ErrorHandler(err, "couldn't update this row")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Collection{}, utils.ErrorHandler(err, "no affected row found")
	}
	fmt.Println(rowsAffected, "row's are affected")

	return update, nil
}

func PatchCollectionDBHandler(id string, updates map[string]interface{}) (models.Collection, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Collection{}, utils.ErrorHandler(err, "couldn't connect to db")
	}

	defer db.Close()

	var existingCollection models.Collection
	row := db.QueryRow(`SELECT id, title, featured_product_id FROM Collection WHERE id = ?`, id)
	err = row.Scan(&existingCollection.ID, &existingCollection.Title, &existingCollection.FeaturedProductID)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Collection{}, utils.ErrorHandler(err, "no rows with this id")
		}
		return models.Collection{}, utils.ErrorHandler(err, "Couldn't scan the query")
	}

	collectionVL := reflect.ValueOf(&existingCollection).Elem()
	collectionType := collectionVL.Type()

	for k, v := range updates {
		for i := 0; i < collectionVL.NumField(); i++ {
			field := collectionType.Field(i)

			if field.Tag.Get("json") == k || field.Tag.Get("json") == k+",omitempty" {
				if collectionVL.Field(i).CanSet() {
					fieldVal := collectionVL.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(collectionVL.Field(i).Type()))
				}
			}
		}
	}

	res, err := db.Exec(`UPDATE Collection SET title = ?, featured_product_id = ? WHERE id = ?`, &existingCollection.Title, &existingCollection.FeaturedProductID, &existingCollection.ID)
	if err != nil {
		return models.Collection{}, nil
	}
	affectedRows, _ := res.RowsAffected()
	fmt.Println(affectedRows, "affected rows")

	return existingCollection, nil
}

func PatchCollectionsDBHandler(w http.ResponseWriter, updates []map[string]interface{}) ([]models.Collection, error) {

	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "couldn't connect to db")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Transaction can't begin")
	}

	var updatedCollection []models.Collection

	for _, update := range updates {
		id, ok := update["id"]
		if !ok {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Invalid id")
		}

		var existingCollection models.Collection
		row := db.QueryRow(`SELECT id, title, featured_product_id FROM Collection WHERE id = ?`, id)
		err = row.Scan(&existingCollection.ID, &existingCollection.Title, &existingCollection.FeaturedProductID)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, utils.ErrorHandler(err, "Couldn't find any row with this id")
			}
			return nil, utils.ErrorHandler(err, "Couldnt scan the values")
		}

		collectionVal := reflect.ValueOf(&existingCollection).Elem()
		colletionType := collectionVal.Type()

		for k, v := range update {
			for i := 0; i < collectionVal.NumField(); i++ {
				field := colletionType.Field(i)
				if field.Tag.Get("json") == k || field.Tag.Get("json") == k+",omiempty" {
					if collectionVal.Field(i).CanSet() {
						fieldVal := collectionVal.Field(i)
						fieldVal.Set(reflect.ValueOf(v).Convert(fieldVal.Type()))
					}
					break
				}
			}
		}

		res, err := tx.Exec("UPDATE Collection SET title = ?, featured_product_id = ? WHERE id = ?", &existingCollection.Title, &existingCollection.FeaturedProductID, &existingCollection.ID)
		if err != nil {
			return nil, utils.ErrorHandler(err, "couldn't update the collections")
		}

		affectedRow, _ := res.RowsAffected()
		fmt.Println(affectedRow, "Affected rows")

		updatedCollection = append(updatedCollection, existingCollection)
	}

	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error commiting the transaction")
	}

	defer db.Close()
	return updatedCollection, nil
}

func DeleteCollectionDBHandler(w http.ResponseWriter, id string) (string, error) {
	db, err := ConnectDB()
	if err != nil {
		return "", utils.ErrorHandler(err, "couldn't connect to db")
	}

	defer db.Close()

	res, er := db.Exec("DELETE FROM Collection WHERE id = ?", id)
	if er != nil {
		return "", utils.ErrorHandler(err, "couldn't delete the collection")
	}

	affectedRow, err := res.RowsAffected()
	if err != nil {
		return "", utils.ErrorHandler(err, "there ia no affected row")
	}
	fmt.Println(affectedRow, "affected rows")

	return id, nil
}

func DeleteCollectionsDBHandler(w http.ResponseWriter, ids []string) ([]string, error) {

	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Couldn't connect to db")
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "transaction canceled")
	}

	stmt, err := tx.Prepare("DELETE FROM Collection WHERE id = ?")
	if err != nil {
		return nil, utils.ErrorHandler(err, "couldnt prepare the statement")
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

		affectedRow, _ := res.RowsAffected()

		if affectedRow > 0 {
			deletedIDs = append(ids, id)
		}

		if affectedRow < 1 {
			return nil, utils.ErrorHandler(err, "invalid transaction there is no affected row")
		}
		fmt.Println(affectedRow, "rows are affected")
	}

	err = tx.Commit()
	if err != nil {
		log.Println("err")
		return nil, utils.ErrorHandler(err, "Invalid transaction")
	}

	return deletedIDs, nil
}
