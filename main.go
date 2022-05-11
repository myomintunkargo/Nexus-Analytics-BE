package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"net/http"

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

func contains (s []int, finder int) bool {
	for _, v := range s{
		if v == finder{
			return true
		}
	}

	return false
}

func handleJulyOnwards(month, year int, order *deliveryOrder) bool{
	queryYear := year
	queryMonths := []int{month-6, month-5, month-4, month-3, month-2, month-1} //assumed to be 6 for now
	if order.Year == queryYear && contains(queryMonths, order.Month){
		return true
	}
	return false
}

func handleJune(month, year int, order *deliveryOrder) bool{
	queryYear := year
	queryMonths := []int{month-5, month-4, month-3, month-2, month-1} //assumed to be 6 for now
	queryPastYearMonths := []int{12}
	if (order.Year == queryYear && contains(queryMonths, order.Month)) || (order.Year == year-1 && contains(queryPastYearMonths, order.Month)){
		return true
	}
	return false
}

func handleMay(month, year int, order *deliveryOrder) bool{
	queryYear := year
	queryMonths := []int{month-4, month-3, month-2, month-1} //assumed to be 6 for now
	queryPastYearMonths := []int{12, 11}
	

	if (order.Year == queryYear && contains(queryMonths, order.Month)) || (order.Year == year-1 && contains(queryPastYearMonths, order.Month)){
		return true
	}

	return false
}

func handleApril(month, year int, order *deliveryOrder) bool{
	queryYear := year
	queryMonths := []int{month-3, month-2, month-1} //assumed to be 6 for now
	queryPastYearMonths := []int{12, 11, 10}
	if (order.Year == queryYear && contains(queryMonths, order.Month)) || (order.Year == year-1 && contains(queryPastYearMonths, order.Month)){
		return true
	}
	return false
}

func handleMarch(month, year int, order *deliveryOrder) bool{
	queryYear := year
	queryMonths := []int{month-2, month-1} //assumed to be 6 for now
	queryPastYearMonths := []int{12, 11, 10, 9}
	if (order.Year == queryYear && contains(queryMonths, order.Month)) || (order.Year == year-1 && contains(queryPastYearMonths, order.Month)){
		return true
	}
	return false
}

func handleFebruary(month, year int, order *deliveryOrder) bool{
	queryYear := year
	queryMonths := []int{month-1} //assumed to be 6 for now
	queryPastYearMonths := []int{12, 11, 10, 9, 8}
	if (order.Year == queryYear && contains(queryMonths, order.Month)) || (order.Year == year-1 && contains(queryPastYearMonths, order.Month)){
		return true
	}
	return false
}

func handleJanuary(month, year int, order *deliveryOrder) bool{
	queryYear := year
	queryMonths := []int{} //assumed to be 6 for now
	queryPastYearMonths := []int{12, 11, 10, 9, 8, 7}
	// Check for  current year's corresponding month orders and previous year's corresponding orders 
	if (order.Year == queryYear && contains(queryMonths, order.Month)) || (order.Year == year-1 && contains(queryPastYearMonths, order.Month)){
		return true
	}
	return false
}

//function to retrieve all orders from past N months -- for now only assumed past 6 months
func getPastMonthsDeliveryOrder(monthsPlaceholder string)  (requiredOrders []deliveryOrder){

	// months, _ := strconv.Atoi(monthsPlaceholder)
	t := time.Now()
	year := t.Year()   // type int
	mthPlaceholder := t.Month() // type time.Month
	month := int(mthPlaceholder) // convert to type int

	julToDec := []int{7,8,9,10,11,12}

	for i, order := range allRows{
		
		if contains(julToDec, month){ // handle logic from july to december
			if handleJulyOnwards(month, year, &order){
				requiredOrders = append(requiredOrders, allRows[i])
			}
		} else if month == 6{
			if handleJune(month, year, &order){
				requiredOrders = append(requiredOrders, allRows[i])
			}
		} else if month == 5{

			if handleMay(month, year, &order){
				requiredOrders = append(requiredOrders, allRows[i])
			}
		} else if month == 4{
			if handleApril(month, year, &order){
				requiredOrders = append(requiredOrders, allRows[i])
			}
		} else if month == 3{
			if handleMarch(month, year, &order){
				requiredOrders = append(requiredOrders, allRows[i])
			}
		} else if month == 2{
			if handleFebruary(month, year, &order){
				requiredOrders = append(requiredOrders, allRows[i])
			}
		} else if month == 1{
			if handleJanuary(month, year, &order){
				requiredOrders = append(requiredOrders, allRows[i])
			}
		}		
		}
		return
		
	}

//GET request function to retrieve number of orders from past N months -- for now only assumed past 6 months
func pastNMonths(c *gin.Context){
	months := c.Param("months")
	orders:= getPastMonthsDeliveryOrder(months)

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
	months := c.Param("months")
	orders := getPastMonthsDeliveryOrder(months)

	monthlyOrdersMap := make(map[string]int) // map monthYear to total number of orders in that month
	monthlyStopsMap := make(map[string]int) // map monthYear to total number of stops in that month
	averageStopsMap := make(map[string]float32) //average per month

	for _, order := range orders{
		var monthYearRepresentation = fmt.Sprintf("%v-%v", order.Month, order.Year)
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

	router := gin.Default()
	router.GET("/api/orders/multistops/:months", pastNMonths)
	router.GET("/api/orders/multistops/average/:months", averagePastNMonthsNumberOfStops)
	router.Run("localhost:5000")

 }


