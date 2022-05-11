# Nexus-Analytics-BE

## Description
This is the backend of Nexus Analytics. 

## Built With
Language Used: Go

Dependencies Used:
* [Gin](https://github.com/gin-gonic/gin) (go get github.com/gin-gonic/gin)
* [GoDotEnv](https://github.com/joho/godotenv) (go get github.com/joho/godotenv)
* [pq](https://github.com/lib/pq) (go get github.com/lib/pq)

## Api Endpoints 

GET("/api/orders/multistops/") - get total number of DOs in past *6 months*
Required Params: None
Sample Response
    { 
    "1-2022": 936, 
    "11-2021": 1029, 
    "12-2021": 785, 
    "2-2022": 1391, 
    "3-2022": 2424, 
    "4-2022": 10475 
    } 
GET("/api/orders/multistops/average/") - get average number of stops per DO in past *6 months*
Required Params: None
Sample Response
    { 
    "1-2022": 3.474359, 
    "11-2021": 2.1788144, 
    "12-2021": 2.2917197, 
    "2-2022": 7.311287, 
    "3-2022": 16.955858, 
    "4-2022": 4.1030073 
     } 



