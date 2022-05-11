# Nexus-Analytics-BE

## Description
This is the backend of Nexus Analytics. 

## Built With
Language Used: Go

Dependencies Used:
* [Gin](https://github.com/gin-gonic/gin) (go get github.com/gin-gonic/gin)
* [GoDotEnv](https://github.com/joho/godotenv) (go get github.com/joho/godotenv)
* [pq](https://github.com/lib/pq) (go get github.com/lib/pq)

## Api Calls

GET("/api/orders/multistops/:months") - get total number of DOs in past *months*
GET("/api/orders/multistops/average/:months") - get average number of stops per DO in past *months*




