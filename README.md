# README #
This README would normally document whatever steps are necessary to get your application up and running.

* Install Go 1.16 above
* Clone the project and set the name as pace_merchant
* Install MySQL
* Create database "pace"
* Run the DDL.sql in the database "pace"
* Update the dbhost, username and password in conf/app.conf for datasource:
   [username]:[password]@tcp([dbhost]:3306)/wallet?charset=utf8
* Install beego:
   * go get github.com/beego/beego/v2@v2.0.0
* Build and run:
   * go mod vendor
   * go build main.go
  ./main
* For testing including documentation, after running it in localhost, just go to http://localhost:8080/swagger/