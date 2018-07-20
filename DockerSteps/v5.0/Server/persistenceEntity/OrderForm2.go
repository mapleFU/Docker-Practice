/**
采用数据库的web服务形式
*/
package persistenceEntity

import (
	"github.com/satori/go.uuid"
	_ "github.com/go-sql-driver/mysql"

	"github.com/mapleFU/GoSQLServerDemo/entity"

	"database/sql"
	"time"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"log"
)


var db sql.DB
var dbTable string
var insertStmt *sql.Stmt
var deleteStmt *sql.Stmt
var queryStmt *sql.Stmt
var queryAllStmt *sql.Stmt

//const layout  = "2018-06-02 01:16:54.505967419 +0800 CST"
const layout  = "2006-01-02 15:04:05 +0800 CST"

type OrderForm entity.OrderForm

func (form *OrderForm) insert() {

	result , err := insertStmt.Exec(form.Good, form.T.String()[:39], form.OrderFormId.String())
	if err != nil {
		panic(err)
	}
	log.Println(result)
}

func formLoader(good, ctime, uid string) *entity.OrderForm {

	t, err := time.Parse(layout, ctime)
	if err != nil {
		return nil
	}
	nuid ,err := uuid.FromString(uid)
	if err != nil {
		return nil
	}

	var form = entity.OrderForm{
		Good:good,
		T: t,
		OrderFormId:nuid,
	}
	log.Println(good, t, nuid)
	return &form
}

func init()  {

	sqlFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	var sqlData map[string]interface{}
	err = json.Unmarshal(sqlFile, &sqlData)
	if err != nil {
		panic(err)
	}
	var dataMap = sqlData["mysql"].(map[string]interface{})
	dbName := dataMap["database"].(string)
	dbTable = dataMap["table"].(string)

	//dbUri := fmt.Sprintf("%s:%s@/%s", dataMap["username"].(string), dataMap["password"].(string), dbName)
	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dataMap["username"].(string), dataMap["password"].(string),
		dataMap["dbhost"].(string), int(dataMap["dbport"].(float64)),dbName)
	fmt.Println(dbUri)

	db, err := sql.Open("mysql", dbUri)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Println("DB connect error.")
	}

	insertStmt, err = db.Prepare("INSERT INTO "+ dbTable+ " (good, create_time, uuid) VALUES (?, ?, ?)")
	if err != nil {
		panic(err)
	}

	deleteStmt, err = db.Prepare("DELETE FROM " + dbTable +" WHERE (uuid=?)")
	if err != nil {
		panic(err)
	}

	// TODO: make clear what () means in mysql
	queryStmt, err = db.Prepare("SELECT good, create_time, uuid FROM "+ dbTable + " WHERE uuid=?")
	if err != nil {
		panic(err)
	}

	queryAllStmt, err = db.Prepare("SELECT good, create_time, uuid FROM "+ dbTable)
	if err != nil {
		panic(err)
	}
}

func NewFormPersistence(goodName string) *OrderForm {
	// raise a panic
	uid := uuid.Must(uuid.NewV4())
	newForm := OrderForm{
		OrderFormId: uid,
		Good:        goodName,
		T:           time.Now(),
	}

	newForm.insert()

	//uidMap[uid] = &newForm

	return &newForm
}

func loadFormFromRow(row *sql.Row) *OrderForm {
	var good, t, uid string
	err := row.Scan(&good, &t, &uid)

	if err != nil {
		return nil
	}
	form := formLoader(good, t, uid)
	return (*OrderForm)(form)
}

func GetFormPersistence(uuidF uuid.UUID) *OrderForm {
	//formPtr, ok := uidMap[uuidF]
	row := queryStmt.QueryRow(uuidF)
	return loadFormFromRow(row)
}

func GetFormsPersistence() []*OrderForm {
	var forms []*OrderForm

	rows, err := queryAllStmt.Query()
	// why this cannot run?
	//rows, err := db.Query(`SELECT good, create_time, uuid FROM ?`, dbTable)
	//rows, err := db.Query("SELECT * FROM " + dbTable)
	if err != nil {
		return nil
	}

	var good, t, uid string
	for rows.Next() {
		err := rows.Scan(&good, &t, &uid)

		if err != nil {
			return nil
		}
		form := formLoader(good, t, uid)
		// rows and row
		forms = append(forms, (*OrderForm)(form))
	}
	return forms
}

func DeleteFormPersistence(uuid2 uuid.UUID) bool {
	// https://stackoverflow.com/questions/922398/what-does-a-successful-mysql-delete-return-how-to-check-if-delete-was-successfu
	var s string
	err := deleteStmt.QueryRow(uuid2).Scan(&s)
	// what if ... error
	if err == nil {
		return false
	}
	return true
}

