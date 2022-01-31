package test

import (
	"bytes"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/google/uuid"
	"merchant/common"
	"merchant/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	_ "merchant/routers"
)

var (
	Authentication string
	MerchantCode   string
	MemberID       string
	EmailAddress   = "superadmin@merchant.com"
	Password       = "password"
	URLMember      = "/v1/member"
	URLMerchant    = "/v1/merchant"
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
		EmailAddress: EmailAddress,
		Password:     Password,
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
		Address:      "myaddress",
		Name:         "myname",
		MerchantCode: uuid.New().String(),
	})
	r, _ := http.NewRequest("POST", URLMerchant, bytes.NewReader(requestBody))
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

// TestCreateMember is a sample to run an endpoint member
func TestCreateMember(t *testing.T) {
	EmailAddress = fmt.Sprintf("administrator@%s.com", MerchantCode)
	Password = "Merchant!234"
	TestAuthentication(t)
	requestBody, _ := json.Marshal(common.MemberRequest{
		Address:      "myaddress",
		Name:         "myname",
		Password:     "password",
		EmailAddress: fmt.Sprintf("%d@merchant.com", time.Now().Unix()),
		Role:         "user",
	})
	r, _ := http.NewRequest("POST", URLMember, bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response models.Member
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 201", func() {
			So(w.Code, ShouldEqual, 201)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.Id, ShouldNotEqual, "")
		})
	})
	MemberID = response.Id
}

// TestGetOneMember is a sample to run an endpoint member
func TestGetOneMember(t *testing.T) {
	r, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", URLMember, MemberID), nil)
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response models.Member
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.Id, ShouldNotEqual, "")
		})
	})
}

// TestGetAllMember is a sample to run an endpoint member
func TestGetAllMember(t *testing.T) {
	r, _ := http.NewRequest("GET", URLMember, nil)
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
			So(len(response), ShouldEqual, 2)
		})
	})
}

// TestUpdateMember is a sample to run an endpoint member
func TestUpdateMember(t *testing.T) {
	requestBody, _ := json.Marshal(common.MemberRequest{
		Address:      "youraddress",
		Name:         "yourname",
		Password:     "password",
		EmailAddress: fmt.Sprintf("%d@merchant.com", time.Now().Unix()),
		Role:         "user",
	})
	r, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s", URLMember, MemberID), bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.Message, ShouldEqual, "OK")
		})
	})
}

// TestDeleteMember is a sample to run an endpoint member
func TestDeleteMember(t *testing.T) {
	r, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", URLMember, MemberID), nil)
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.Message, ShouldEqual, "OK")
		})
	})
}

// TestGetMerchant is a sample to run an endpoint member
func TestGetMerchant(t *testing.T) {
	r, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", URLMerchant, MerchantCode), nil)
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response models.Merchant
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.MerchantCode, ShouldNotEqual, "")
		})
	})
}

// TestUpdateMerchant is a sample to run an endpoint member
func TestUpdateMerchant(t *testing.T) {
	requestBody, _ := json.Marshal(common.MerchantRequest{
		Address: "theiraddress",
		Name:    "theirname",
	})
	r, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s", URLMerchant, MerchantCode), bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.Message, ShouldEqual, "OK")
		})
	})
}

// TestGetMerchant is a sample to run an endpoint member
func TestDeleteMerchant(t *testing.T) {
	r, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", URLMerchant, MerchantCode), nil)
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.Message, ShouldEqual, "OK")
		})
	})
}
