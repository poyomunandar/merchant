package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/google/uuid"
)

type Member struct {
	Id           string    `orm:"column(id);pk"`
	Merchant     *Merchant `orm:"column(merchant_id);rel(fk)"`
	EmailAddress string    `orm:"column(email_address);size(255)"`
	Name         string    `orm:"column(name);size(100);null"`
	Address      string    `orm:"column(address);size(255);null"`
	Role         string    `orm:"column(role);size(20);null"`
	IsDeleted    int8      `orm:"column(is_deleted)"`
	CreatedTime  int64     `orm:"column(created_time)"`
	UpdatedTime  int64     `orm:"column(updated_time)"`
	Password     string    `orm:"column(password)"`
}

func (t *Member) TableName() string {
	return "member"
}

func init() {
	orm.RegisterModel(new(Member))
}

// AddMember insert a new Member into database and returns
// last inserted Id on success.
func AddMember(m *Member) (id int64, err error) {
	o := orm.NewOrm()
	m.Id = uuid.New().String()
	m.CreatedTime = time.Now().Unix()
	id, err = o.Insert(m)
	return
}

// GetMemberByMerchantId retrieves Member by MerchantId. Returns error if
// MerchantId doesn't exist
func GetMemberByMerchantId(merchantId string) (v *[]Member, err error) {
	v = &[]Member{}
	o := orm.NewOrm()
	member := Member{}
	qs := o.QueryTable(member)
	_, err = qs.Filter("Merchant__Id", merchantId).Filter("is_deleted", 0).All(v)
	return
}

// GetMemberById retrieves Member by Id. Returns error if
// Id doesn't exist
func GetMemberById(id string) (v *Member, err error) {
	o := orm.NewOrm()
	v = &Member{Id: id}
	qs := o.QueryTable(v)
	err = qs.Filter("id", id).Filter("is_deleted", 0).One(v)
	return
}

// GetMemberByEmail retrieves Member by Email and returns error if
// Email doesn't exist
func GetMemberByEmail(emailAddress string) (v *Member, err error) {
	o := orm.NewOrm()
	v = &Member{EmailAddress: emailAddress}
	qs := o.QueryTable(v)
	err = qs.Filter("email_address", emailAddress).Filter("is_deleted", 0).One(v)
	return
}

// GetAllMember retrieves all Member matches certain condition. Returns empty list if
// no records exist
func GetAllMember(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Member))
	qs = qs.Filter("is_deleted", 0)
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

	var l []Member
	qs = qs.OrderBy(sortFields...)
	_, err = qs.Limit(limit, offset).All(&l, fields...)
	if err == nil {
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

// UpdateMember updates Member by Id and returns error if
// the record to be updated doesn't exist
func UpdateMemberById(m *Member) (err error) {
	o := orm.NewOrm()
	m.UpdatedTime = time.Now().Unix()
	v := Member{Id: m.Id}
	// assert id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		m.CreatedTime = v.CreatedTime
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteMember deletes Member by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMember(id string) (err error) {
	o := orm.NewOrm()
	v := Member{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		v.IsDeleted = 1
		v.UpdatedTime = time.Now().Unix()
		_, err = o.Update(&v, "is_deleted", "updated_time")
	}
	return
}
