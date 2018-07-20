/**
不采用数据库的web服务形式
*/
package entity

import (
	"github.com/satori/go.uuid"

	"time"
)

type OrderForm struct {
	OrderFormId uuid.UUID `json:"order_form_id"` // 对应的uuid
	//Good Goods	// 商品
	Good string    `json:"good"` // 购买的商品
	T    time.Time `json:"created_time"`
}

// var uidMap map[string]*OrderForm{} 是不可以的
var uidMap = map[uuid.UUID]*OrderForm{}

func NewForm(goodName string) *OrderForm {
	// raise a panic
	uid := uuid.Must(uuid.NewV4())
	newForm := OrderForm{
		OrderFormId: uid,
		Good:        goodName,
		T:           time.Now(),
	}
	uidMap[uid] = &newForm

	return &newForm
}

func GetForm(uuidF uuid.UUID) *OrderForm {
	formPtr, ok := uidMap[uuidF]
	if !ok {
		return nil
	}
	return formPtr
}

func GetForms() []*OrderForm {
	var forms []*OrderForm
	for _, value := range uidMap {
		forms = append(forms, value)
	}
	return forms
}

func DeleteForm(uuid2 uuid.UUID) *OrderForm {
	formPtr, ok := uidMap[uuid2]
	if !ok {
		return nil
	}
	delete(uidMap, uuid2)
	return formPtr
}
