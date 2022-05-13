package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

//Defining the deliveryOrder struct
type deliveryOrder struct{
	Job_ksuid string `json:"job_ksuid"`
	NumStops string  `json:"numStops"`
	Month int `json:"month"`
	Year int `json:"year"`
}
	
func dbConnection() []deliveryOrder{
	// Loading environment variables from local.env file
	err1 := godotenv.Load("local.env")
	if err1 != nil {
		log.Fatalf("Some error occured. Err: %s", err1)
	}
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// Database connection string
	dbURI := fmt.Sprintf("port=%s host=%s user=%s "+ "password=%s dbname=%s sslmode=disable", dbPort, host, user, password, dbName)
	

	// Create database object
	db, err := sql.Open(dialect,
		dbURI)

	if err != nil {
		log.Fatal(err)
	}
	
	

	// Declare variables for query responses 
	var (
		job_ksuid string
		numStops string
		month string
		year string
	)

	var allRows = []deliveryOrder{}

	// SQL query to get the 1) job_ksuid 2) Count = number of rows 3) Month of the DO 4) Year of the DO
	query := `
		SELECT 
			job_ksuid, 
			count(*), 
			EXTRACT(MONTH FROM(inserted_at)) as mth,
			EXTRACT(YEAR FROM(inserted_at)) as yr
		FROM 
			load.shipment_points
		GROUP BY
			job_ksuid, inserted_at
		HAVING
			COUNT(job_ksuid)>1
		;
	`

	//Get rows using the query
	rows, err := db.Query(query)
	if err != nil { //Log if error
		log.Fatal(err)
	}
	defer rows.Close()

	// Add each row into the "allRows" slice
	for rows.Next() {

		err := rows.Scan(&job_ksuid, &numStops, &month, &year)		
		if err != nil {
			log.Fatal(err)
		}

		monthValue, _ := strconv.Atoi(month)
		yearValue, _ := strconv.Atoi(year)

		//Create new deliverOrder struct with the received data
		row := deliveryOrder{
			Job_ksuid: job_ksuid,
			NumStops: numStops,
			Month: monthValue ,
			Year: yearValue,
		}
		allRows = append(allRows, row)
	}
		//Log if error
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}


	defer db.Close()

	return allRows
}

var allRows = dbConnection()

// Get the number of months between firstDate and secondDate
func monthDifference(firstDate, secondDate time.Time) (month int) {
	if firstDate.Location() != secondDate.Location() {
		secondDate = secondDate.In(firstDate.Location())
	}

	y1, M1, _ := firstDate.Date()
	y2, M2, _ := secondDate.Date()

	year := int(y2 - y1)
	month = int(M2 - M1)


	// Normalize negative values
	for year > 0{
		month += 12
		year--
	}

	return
}

//function to retrieve all orders from past N months -- for now only assumed past 6 months
func getPastMonthsDeliveryOrder()  (requiredOrders []deliveryOrder){

	numMonths, _ := strconv.Atoi(monthsPlaceholder)
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), 1, 1, 1, 1, 1, time.Local)

	for _, order := range allRows{

		orderDate := time.Date(order.Year, time.Month(order.Month), 1, 1, 1, 1, 1, time.Local) //convert the order date into Date struct
		monthsApart := monthDifference(orderDate , today)
		if monthsApart <= numMonths{
			requiredOrders = append(requiredOrders, order)
		}

	}		
	return
}

	
//GET request function to retrieve number of orders from past N months -- for now only assumed past 6 months
func pastNMonths(c *gin.Context){
	// months := c.Param("months")
	orders:= getPastMonthsDeliveryOrder()

	monthlyOrdersMap := make(map[string]int) // map monthYear to total number of orders in that month

	for _, order := range orders{
		var monthYearRepresentation = fmt.Sprintf("%v-%v", order.Month, order.Year)
		var monthYearOrders = fmt.Sprintf("%v", monthYearRepresentation)
		monthlyOrdersMap[monthYearOrders] += 1
	}

	c.IndentedJSON(http.StatusOK, monthlyOrdersMap)
}


//GET request function to retrieve average stops of all orders from past N months -- for now only assumed past 6 months
func averagePastNMonthsNumberOfStops(c *gin.Context){ 
	// months := c.Param("months")
	orders := getPastMonthsDeliveryOrder()

	monthlyOrdersMap := make(map[string]int) // map monthYear to total number of orders in that month
	monthlyStopsMap := make(map[string]int) // map monthYear to total number of stops in that month
	averageStopsMap := make(map[string]float32) //average per month

	for _, order := range orders{
		var monthYearRepresentation = fmt.Sprintf("%v-%v", order.Month, order.Year) // represent dates as month-year
		var monthYearOrders = fmt.Sprintf("%v", monthYearRepresentation) 
		var monthYearStops = fmt.Sprintf("%v", monthYearRepresentation)
		monthlyOrdersMap[monthYearOrders] += 1
		numberOfStops, _ := strconv.Atoi(order.NumStops)
		monthlyStopsMap[monthYearStops] += numberOfStops
	}

	for monthYear := range monthlyOrdersMap{

		averageStopsMap[monthYear] = float32(monthlyStopsMap[monthYear]) /float32(monthlyOrdersMap[monthYear])
	}

	c.IndentedJSON(http.StatusOK, averageStopsMap)
}


func main() {

	// getPastMonthsDeliveryOrder("5")
	router := gin.Default()

    router.Use(cors.New(cors.Config{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{"GET"},
        AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
    }))

	router.GET("/api/orders/multistops/:months", pastNMonths)
	router.GET("/api/orders/multistops/average/:months", averagePastNMonthsNumberOfStops)
	router.Run("localhost:5000")

 }


