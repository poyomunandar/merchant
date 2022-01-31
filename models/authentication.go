package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

type Authentication struct {
	Id         int    `orm:"column(id);auto"`
	Token      string `orm:"column(token);size(255)"`
	MemberId   string `orm:"column(account_id);size(255)"`
	ExpiryTime int64  `orm:"column(expiry_time)"`
}

func (t *Authentication) TableName() string {
	return "authentication"
}

func init() {
	orm.RegisterModel(new(Authentication))
}

// AddAuthentication insert a new Authentication into database and returns
// last inserted Id on success.
func AddAuthentication(m *Authentication) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAuthenticationByToken retrieves Authentication by token. Returns error if
// token doesn't exist
func GetAuthenticationByToken(token string) (v *Authentication, err error) {
	o := orm.NewOrm()
	v = &Authentication{Token: token}
	qs := o.QueryTable(v)
	err = qs.Filter("token", token).One(v)
	return
}

// GetAuthenticationById retrieves Authentication by Id. Returns error if
// Id doesn't exist
func GetAuthenticationById(id int) (v *Authentication, err error) {
	o := orm.NewOrm()
	v = &Authentication{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAuthentication retrieves all Authentication matches certain condition. Returns empty list if
// no records exist
func GetAllAuthentication(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Authentication))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Authentication
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateAuthentication updates Authentication by Id and returns error if
// the record to be updated doesn't exist
func UpdateAuthenticationById(m *Authentication) (err error) {
	o := orm.NewOrm()
	v := Authentication{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}
