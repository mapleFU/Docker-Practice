// https://stackoverflow.com/questions/3050518/what-http-status-response-code-should-i-use-if-the-request-is-missing-a-required
// 耗时 1.5 天
package Server

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	//"github.com/sirupsen/logrus"

	"github.com/mapleFU/GoSQLServerDemo/persistenceEntity"

	"fmt"
	"net/url"
	//"time"
	//"github.com/gin-gonic/contrib/ginrus"
	"github.com/sirupsen/logrus"
)


func main() {

	//gin.DefaultWriter = log.Logger{}
	r := gin.Default()


	// 获得订单信息, 利用uid来get
	r.GET("/orderForm", func(context *gin.Context) {

		queryUID, ifQuery := context.GetQuery("id")
		fmt.Println(queryUID, ifQuery)
		if ifQuery {
			uid, err := uuid.FromString(queryUID)
			if err != nil {
				// 参数错误
				context.Status(422)
			}
			form := persistenceEntity.GetFormPersistence(uid)
			if form == nil {
				// doesn't exists
				context.Status(404)
			} else {
				// write json

				form := persistenceEntity.GetFormPersistence(uid)
				context.JSON(200, form)
				logrus.WithFields(logrus.Fields{
					"event": "read",
					"id": form.OrderFormId,
				}).Info("read a form")
			}
		} else {
			//GetForm(uuid.UUID(queryUID))
			forms := persistenceEntity.GetFormsPersistence()

			context.JSON(200, gin.H{
				"order_form_num": len(forms),
				"forms":          forms,
			})
			logrus.WithFields(logrus.Fields{
				"event": "read",
			}).Info("read all forms")
		}
	})

	ret := url.URL{
		Path: "/orderForm",
	}

	r.POST("/orderForm", func(context *gin.Context) {
		formGood, ok := context.GetPostForm("good")

		if !ok {
			//context.Header()
			context.JSON(422, gin.H{
				"error": "good not exists",
			})
			return
		}
		curRet := ret.Query()
		form := persistenceEntity.NewFormPersistence(formGood)
		curRet.Set("id", form.OrderFormId.String())
		ret.RawQuery = curRet.Encode()
		context.Header("Location", ret.String())
		logrus.WithFields(logrus.Fields{
			"event": "add",
			"formID": form.OrderFormId,
		}).Info("add a form.")
		context.JSON(201, *form)
	})

	r.DELETE("/orderForm", func(context *gin.Context) {
		formGood, ok := context.GetQuery("id")

		if ok {
			uid, err := uuid.FromString(formGood)
			if err != nil {
				// 参数错误
				context.JSON(422, gin.H{
					"error": "good not exists",
				})
				return
			}
			formExisted := persistenceEntity.DeleteFormPersistence(uid)
			if formExisted == false {
				// doesn't exists --> 404/410
				context.Status(410)
			} else {
				// success --> 204
				logrus.WithFields(logrus.Fields{
					"event": "delete",
					"formID": uid,
				}).Info("add a form.")
				context.Status(204)
			}
		} else {
			context.JSON(422, gin.H{
				"error": "argument id not in",
			})
			return
		}
	})

	// UPDATE
	r.PATCH("/orderForm", func(context *gin.Context) {

	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
