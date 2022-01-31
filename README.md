# README #
This README would normally document whatever steps are necessary to get your application up and running.

* Install Go 1.16 above
* Clone the project and set the name as merchant
* Install MySQL
* Create database "pace" and "pace_test"
* Run the DDL.sql in the database "pace" and "pace_test"
* Update the dbhost, username and password in conf/app.conf for datasource:
    * [username]:[password]@tcp([dbhost]:3306)/pace?charset=utf8
    * [username]:[password]@tcp([dbhost]:3306)/pace_test?charset=utf8
* Install beego:
    * go get github.com/beego/beego/v2@v2.0.0
* Build and run:
    * go mod vendor
    * go build main.go
    * ./main
* For testing including documentation, after running it in localhost, just go to http://localhost:8080/swagger/
* For running unit test, run:
    * cd tests
    * go test -v
* Notes: use email: superadmin@merchant.com and password: password for creating new merchant.
    * There are three roles available: user, administrator, and superadmin
    * Superadmin can do anything, and only superadmin can create new merchant
    * Once merchant is created, it will automatically create the administrator member for that merchant as well
    * administrator member can create, view, update, delete all members in the same merchant, also can view, update and delete their merchant.
    * user member can only view and update their self only, and can only view their merchant as well