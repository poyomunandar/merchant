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
	MerchantID     string
	MemberID       string
	MerchantCode   string
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
func TestCreateMerchantNormal(t *testing.T) {
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
	MerchantID = response.Id
	MerchantCode = response.MerchantCode
}

// TestCreateMerchantDuplicateMerchantCode is a sample to run an endpoint merchant
func TestCreateMerchantDuplicateMerchantCode(t *testing.T) {
	requestBody, _ := json.Marshal(common.MerchantRequest{
		Address:      "myaddress",
		Name:         "myname",
		MerchantCode: "CODE-1234",
	})
	r, _ := http.NewRequest("POST", URLMerchant, bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A merchant", func() {
			So(response.Message, ShouldStartWith, "Error 1062")
		})
	})
}

// TestCreateMerchantWrongAuth is a sample to run an endpoint merchant
func TestCreateMerchantWrongAuth(t *testing.T) {
	requestBody, _ := json.Marshal(common.MerchantRequest{
		Address:      "myaddress",
		Name:         "myname",
		MerchantCode: uuid.New().String(),
	})
	r, _ := http.NewRequest("POST", URLMerchant, bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication+"XX")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A merchant", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestCreateMemberNormal is a sample to run an endpoint member
func TestCreateMemberNormal(t *testing.T) {
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

// TestCreateMemberEmailWrongFormat is a sample to run an endpoint member
func TestCreateMemberEmailWrongFormat(t *testing.T) {
	requestBody, _ := json.Marshal(common.MemberRequest{
		Address:      "myaddress",
		Name:         "myname",
		Password:     "password",
		EmailAddress: "email",
		Role:         "user",
	})
	r, _ := http.NewRequest("POST", URLMember, bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.Message, ShouldEqual, common.ErrorMessageMap[common.ErrorEmailFormatInvalid])
		})
	})
}

// TestCreateMemberEmailMissing is a sample to run an endpoint member
func TestCreateMemberEmailMissing(t *testing.T) {
	requestBody, _ := json.Marshal(common.MemberRequest{
		Address:  "myaddress",
		Name:     "myname",
		Password: "password",
		Role:     "user",
	})
	r, _ := http.NewRequest("POST", URLMember, bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeRequiredParamEmpty)
		})
	})
}

// TestCreateMemberPasswordMissing is a sample to run an endpoint member
func TestCreateMemberPasswordMissing(t *testing.T) {
	requestBody, _ := json.Marshal(common.MemberRequest{
		Address:      "myaddress",
		Name:         "myname",
		EmailAddress: fmt.Sprintf("%d@merchant.com", time.Now().Unix()),
		Role:         "user",
	})
	r, _ := http.NewRequest("POST", URLMember, bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeRequiredParamEmpty)
		})
	})
}

// TestCreateMemberDuplicateEmail is a sample to run an endpoint member
func TestCreateMemberDuplicateEmail(t *testing.T) {
	requestBody, _ := json.Marshal(common.MemberRequest{
		Address:      "myaddress",
		Name:         "myname",
		Password:     "password",
		EmailAddress: "superadmin@merchant.com",
		Role:         "user",
	})
	r, _ := http.NewRequest("POST", URLMember, bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeEmailExist)
		})
	})
}

// TestCreateMemberWrongAuth is a sample to run an endpoint member
func TestCreateMemberWrongAuth(t *testing.T) {
	requestBody, _ := json.Marshal(common.MemberRequest{
		Address:      "myaddress",
		Name:         "myname",
		Password:     "password",
		EmailAddress: fmt.Sprintf("%d@merchant.com", time.Now().Unix()),
		Role:         "user",
	})
	r, _ := http.NewRequest("POST", URLMember, bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication+"XX")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestGetOneMemberNormal is a sample to run an endpoint member
func TestGetOneMemberNormal(t *testing.T) {
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

// TestGetOneMemberOtherMerchant is a sample to run an endpoint member
func TestGetOneMemberOtherMerchant(t *testing.T) {
	r, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", URLMember, "f55b6262-9ee8-4f07-8855-493b6b5cacb1"), nil)
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestGetOneMemberWrongId is a sample to run an endpoint member
func TestGetOneMemberWrongId(t *testing.T) {
	r, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", URLMember, "xxx-9ee8-4f07-8855-493b6b5cacb1"), nil)
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, 400)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUndefined)
		})
	})
}

// TestGetOneMemberWrongAuth is a sample to run an endpoint member
func TestGetOneMemberWrongAuth(t *testing.T) {
	r, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", URLMember, "xxx-9ee8-4f07-8855-493b6b5cacb1"), nil)
	r.Header.Add("Authorization", Authentication+"XX")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestGetAllMemberNormal is a sample to run an endpoint member
func TestGetAllMemberNormal(t *testing.T) {
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

// TestGetAllMemberWithPagination is a sample to run an endpoint member
func TestGetAllMemberWithPagination(t *testing.T) {
	r, _ := http.NewRequest("GET", fmt.Sprintf("%s?offset=0&limit=1", URLMember), nil)
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
			So(len(response), ShouldEqual, 1)
		})
	})
}

// TestGetAllMemberWithQuery is a sample to run an endpoint member
func TestGetAllMemberWithQuery(t *testing.T) {
	r, _ := http.NewRequest("GET", fmt.Sprintf("%s?query=name:myname", URLMember), nil)
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
			So(len(response), ShouldEqual, 1)
		})
	})
}

// TestGetAllMemberWrongAuth is a sample to run an endpoint member
func TestGetAllMemberWrongAuth(t *testing.T) {
	r, _ := http.NewRequest("GET", fmt.Sprintf("%s?query=name:myname", URLMember), nil)
	r.Header.Add("Authorization", Authentication+"XX")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestUpdateMemberNormal is a sample to run an endpoint member
func TestUpdateMemberNormal(t *testing.T) {
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

// TestUpdateMemberDifferentMerchant is a sample to run an endpoint member
func TestUpdateMemberDifferentMerchant(t *testing.T) {
	requestBody, _ := json.Marshal(common.MemberRequest{
		Address:      "youraddress",
		Name:         "yourname",
		Password:     "password",
		EmailAddress: fmt.Sprintf("%d@merchant.com", time.Now().Unix()),
		Role:         "user",
	})
	r, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s", URLMember, "f55b6262-9ee8-4f07-8855-493b6b5cacb1"), bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestUpdateMemberDifferentMerchant is a sample to run an endpoint member
func TestUpdateMemberWrongAuth(t *testing.T) {
	requestBody, _ := json.Marshal(common.MemberRequest{
		Address:      "youraddress",
		Name:         "yourname",
		Password:     "password",
		EmailAddress: fmt.Sprintf("%d@merchant.com", time.Now().Unix()),
		Role:         "user",
	})
	r, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s", URLMember, MemberID), bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication+"XX")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestDeleteMemberNormal is a sample to run an endpoint member
func TestDeleteMemberNormal(t *testing.T) {
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

// TestDeleteMemberDifferentMerchant is a sample to run an endpoint member
func TestDeleteMemberDifferentMerchant(t *testing.T) {
	r, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", URLMember, "f55b6262-9ee8-4f07-8855-493b6b5cacb1"), nil)
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestDeleteMemberWrongAuth is a sample to run an endpoint member
func TestDeleteMemberWrongAuth(t *testing.T) {
	r, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", URLMember, MemberID), nil)
	r.Header.Add("Authorization", Authentication+"XX")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestGetMerchantNormal is a sample to run an endpoint member
func TestGetMerchantNormal(t *testing.T) {
	r, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", URLMerchant, MerchantID), nil)
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

// TestGetMerchantNotOnScope is a sample to run an endpoint member
func TestGetMerchantNotOnScope(t *testing.T) {
	r, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", URLMerchant, "c8b49264-93cb-4b19-bad3-81953cf5317e"), nil)
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestGetMerchantWrongAuth is a sample to run an endpoint member
func TestGetMerchantWrongAuth(t *testing.T) {
	r, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", URLMerchant, MerchantID), nil)
	r.Header.Add("Authorization", Authentication+"XX")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestUpdateMerchantNormal is a sample to run an endpoint member
func TestUpdateMerchantNormal(t *testing.T) {
	requestBody, _ := json.Marshal(common.MerchantRequest{
		Address:      "theiraddress",
		Name:         "theirname",
		MerchantCode: uuid.New().String(),
	})
	r, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s", URLMerchant, MerchantID), bytes.NewReader(requestBody))
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

// TestUpdateMerchantNotOnScope is a sample to run an endpoint member
func TestUpdateMerchantNotOnScope(t *testing.T) {
	requestBody, _ := json.Marshal(common.MerchantRequest{
		Address:      "theiraddress",
		Name:         "theirname",
		MerchantCode: uuid.New().String(),
	})
	r, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s", URLMerchant, "c8b49264-93cb-4b19-bad3-81953cf5317e"), bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestUpdateMerchantWrongAuth is a sample to run an endpoint member
func TestUpdateMerchantWrongAuth(t *testing.T) {
	requestBody, _ := json.Marshal(common.MerchantRequest{
		Address:      "theiraddress",
		Name:         "theirname",
		MerchantCode: uuid.New().String(),
	})
	r, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s", URLMerchant, MerchantID), bytes.NewReader(requestBody))
	r.Header.Add("Authorization", Authentication+"XX")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestGetMerchant is a sample to run an endpoint member
func TestDeleteMerchantNormal(t *testing.T) {
	r, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", URLMerchant, MerchantID), nil)
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
func TestDeleteMerchantNotOnScope(t *testing.T) {
	r, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", URLMerchant, "c8b49264-93cb-4b19-bad3-81953cf5317e"), nil)
	r.Header.Add("Authorization", Authentication)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}

// TestGetMerchant is a sample to run an endpoint member
func TestDeleteMerchantWrongAuth(t *testing.T) {
	r, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", URLMerchant, MerchantID), nil)
	r.Header.Add("Authorization", Authentication+"XX")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var response common.ErrorMessage
	json.Unmarshal(w.Body.Bytes(), &response)

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 401", func() {
			So(w.Code, ShouldEqual, 401)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
		Convey("There Should Be A member", func() {
			So(response.ErrorCode, ShouldEqual, common.ErrorCodeUnauthrized)
		})
	})
}
