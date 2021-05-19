package controllers

import (
	"cict-quiz-api/app/database"
	"cict-quiz-api/app/models"
	"cict-quiz-api/app/routes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/revel/revel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

//login
type AuthController struct {
	* revel.Controller
}

func (c AuthController) Index() revel.Result{
	return c.Render()

}

//Api
func (c AuthController) Login() revel.Result {
	defer c.Request.Destroy()
	login := &models.Login{}
	c.Response.Status = http.StatusBadRequest
	data := make(map[string]interface{})
	data["status"] = "error"
	if err := json.NewDecoder(c.Request.GetBody()).Decode(&login); err != nil{
		data["data"] = "Could parse request"
		return c.RenderJSON(data)
	}

	if govalidator.IsNull(login.Username) || govalidator.IsNull(login.Password){
		data["data"] = "Username or Password is not correct"
		return c.RenderJSON(data)
	}

	result := &models.User{}
	ctx := context.Background()
	filter := bson.D{primitive.E{Key: "username", Value: login.Username}}

	if err := database.UserCollection.FindOne(ctx,filter).Decode(&result); err != nil{
		data["data"] = "Username not exist"
		return c.RenderJSON(data)
	}

	if err := models.CheckHashAndPassword(login.Password,result.Password); err !=nil{
		data["data"] = "Username or Password is not correct"
		return c.RenderJSON(data)
	}

	c.Response.Status = http.StatusOK
	body := make(map[string]interface{})
	body["user"] = result
	body["token"] = ""
	data["data"] = body
	return c.RenderJSON(data)
}

func  (c AuthController) LoginForm() revel.Result {

	login := &models.Login{Username: c.Params.Get("username"), Password: c.Params.Get("password")}

	fmt.Println(login)

	if govalidator.IsNull(login.Username) || govalidator.IsNull(login.Password){
		return c.RenderText("Not validate")
	}

	result := &models.User{}
	ctx := context.Background()
	filter := bson.D{primitive.E{Key: "username", Value: login.Username}}
	filter2 := bson.D{primitive.E{Key: "password", Value: login.Password}}
	if err := database.UserCollection.FindOne(ctx,filter).Decode(&result); err != nil{
		c.Flash.Error("Mật khẩu hoặc tên người dùng chưa đúng!")

		if err := database.UserCollection.FindOne(ctx,filter2).Decode(&result); err != nil{
			c.Flash.Error("Mật khẩu hoặc tên người dùng chưa đúng!")
			return  c.Redirect(AuthController.Index)

		}
		return c.Redirect(routes.QuestionController.Index())
	}


	return c.Redirect(routes.QuestionController.Index())
}



