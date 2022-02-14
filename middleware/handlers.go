package middleware

import (
	"CRUD-app/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct {
    ID      int64  `json:"id,omitempty"`
    Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

    if err != nil {
        panic(err)
    }

    err = db.Ping()

    if err != nil {
        panic(err)
    }

    fmt.Println("Successfully connected!")
    return db
}

func CreateHarga(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    var daftar models.Harga

    err := json.NewDecoder(r.Body).Decode(&daftar)

    if err != nil {
        log.Fatalf("Unable to decode the request body.  %v", err)
    }

    insertID := insertHarga(daftar)

    res := response{
        ID:      insertID,
        Message: "Harga created successfully",
    }

    json.NewEncoder(w).Encode(res)
}

func GetHarga(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
   
    params := mux.Vars(r)

    
    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    harga, err := getHarga(int64(id))

    if err != nil {
        log.Fatalf("Unable to get harga. %v", err)
    }

    json.NewEncoder(w).Encode(harga)
}

func GetAllHarga(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
   
    harga, err := getAllHarga()

    if err != nil {
        log.Fatalf("Unable to get all harga. %v", err)
    }

    json.NewEncoder(w).Encode(harga)
}

func UpdateHarga(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    params := mux.Vars(r)

    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    var harga models.Harga

    err = json.NewDecoder(r.Body).Decode(&harga)

    if err != nil {
        log.Fatalf("Unable to decode the request body.  %v", err)
    }

    updatedRows := updateHarga(int64(id), harga)

    msg := fmt.Sprintf("Harga updated successfully. Total rows/record affected %v", updatedRows)

    res := response{
        ID:      int64(id),
        Message: msg,
    }

    json.NewEncoder(w).Encode(res)
}

func DeleteHarga(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    params := mux.Vars(r)

    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    deletedRows := deleteHarga(int64(id))

    msg := fmt.Sprintf("Harga deleted successfully. Total rows/record affected %v", deletedRows)

    res := response{
        ID:      int64(id),
        Message: msg,
    }

    json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------

func insertHarga(harga models.Harga) int64 {

    db := createConnection()

    defer db.Close()

    sqlStatement := `INSERT INTO bahan_bakar (liter, premium, pertalite) VALUES ($1, $2, $3) RETURNING idharga`

    var id int64

    err := db.QueryRow(sqlStatement, harga.Liter, harga.Premium, harga.Pertalite).Scan(&id)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    fmt.Printf("Inserted a single record %v", id)

    return id
}

func getHarga(id int64) (models.Harga, error) {
    db := createConnection()

    defer db.Close()

    var harga models.Harga

    sqlStatement := `SELECT * FROM bahan_bakar WHERE idharga=$1`

    row := db.QueryRow(sqlStatement, id)

    err := row.Scan(&harga.ID, &harga.Liter, &harga.Premium, &harga.Pertalite)

    switch err {
    case sql.ErrNoRows:
        fmt.Println("No rows were returned!")
        return harga, nil
    case nil:
        return harga, nil
    default:
        log.Fatalf("Unable to scan the row. %v", err)
    }

    return harga, err
}

func getAllHarga() ([]models.Harga, error) {
    db := createConnection()

    defer db.Close()

    var harga2 []models.Harga

    sqlStatement := `SELECT * FROM bahan_bakar`

    rows, err := db.Query(sqlStatement)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    defer rows.Close()

    for rows.Next() {
        var harga models.Harga

        err = rows.Scan(&harga.ID, &harga.Liter, &harga.Premium, &harga.Pertalite)

        if err != nil {
            log.Fatalf("Unable to scan the row. %v", err)
        }

        harga2 = append(harga2, harga)

    }

    return harga2, err
}

func updateHarga(id int64, harga models.Harga) int64 {
    db := createConnection()

    defer db.Close()

    sqlStatement := `UPDATE baan_bakar SET liter=$2, premium=$3, pertalite=$4 WHERE idharga=$1`

    res, err := db.Exec(sqlStatement, id, harga.Liter, harga.Premium, harga.Pertalite)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    rowsAffected, err := res.RowsAffected()

    if err != nil {
        log.Fatalf("Error while checking the affected rows. %v", err)
    }

    fmt.Printf("Total rows/record affected %v", rowsAffected)

    return rowsAffected
}

func deleteHarga(id int64) int64 {

    db := createConnection()

    defer db.Close()

    sqlStatement := `DELETE FROM bahan_bakar WHERE idharga=$1`

    res, err := db.Exec(sqlStatement, id)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    rowsAffected, err := res.RowsAffected()

    if err != nil {
        log.Fatalf("Error while checking the affected rows. %v", err)
    }

    fmt.Printf("Total rows/record affected %v", rowsAffected)

    return rowsAffected
}