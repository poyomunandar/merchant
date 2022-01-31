package test

import (
	"bytes"
	"github.com/beego/beego/v2/client/orm"
	"merchant/common"
	"merchant/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	_ "merchant/routers"
)

var (
	Authentication string
	MerchantCode   string
)

func init() {
	cwd, _ := os.Getwd()
	roots := strings.Split(cwd, string(os.PathSeparator))
	beego.TestBeegoInit(strings.Join(roots[0:len(roots)-1], "/"))
	conn, _ := beego.AppConfig.String("sqlconntest")
	orm.RegisterDataBase("default", "mysql", conn)
}

// TestAuthentication is a sample to run an endpoint authentication
func TestAuthentication(t *testing.T) {
	requestBody, _ := json.Marshal(common.AuthenticationRequest{
		EmailAddress: "superadmin@merchant.com",
		Password:     "password",
	})
	r, _ := http.NewRequest("POST", "/v1/authentication", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response models.Authentication
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 201", func() {
			So(w.Code, ShouldEqual, 201)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be Token", func() {
			So(response.Token, ShouldNotEqual, "")
		})
	})
	Authentication = response.Token
}

// TestCreateMerchant is a sample to run an endpoint merchant
func TestCreateMerchant(t *testing.T) {
	requestBody, _ := json.Marshal(common.MerchantRequest{
		Address: "myaddress",
		Name:    "myname",
	})
	r, _ := http.NewRequest("POST", "/v1/merchant", bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response models.Merchant
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 201", func() {
			So(w.Code, ShouldEqual, 201)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A merchant", func() {
			So(response.MerchantCode, ShouldNotEqual, "")
		})
	})
	MerchantCode = response.MerchantCode
}

// TestGetMember is a sample to run an endpoint member
func TestGetMember(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/member", nil)
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response []models.Member
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(len(response), ShouldNotEqual, 0)
		})
	})
}
