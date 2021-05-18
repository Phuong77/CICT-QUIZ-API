package controllers

import (
	"cict-quiz-api/app/database"
	"cict-quiz-api/app/models"
	"encoding/json"
	"fmt"
	"github.com/revel/revel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"net/http"
)

//questions
type QuestionController struct {
	*revel.Controller
}
type Todo struct {
	Text string
	parentID string
}

var todos []Todo

func init() {
}

func GetAllCat() []models.Category {
	var result = []models.Category{}
	ctx := context.Background()
	cur, err := database.CategoryCollection.Find(ctx, bson.D{})

	if err != nil {
		return result
	}

	for cur.Next(ctx){
		var cat models.Category
		fmt.Println(cur)
		if err := cur.Decode(&cat); err !=nil{
			return result
		}
		result = append(result, cat)
	}
	if err := cur.Err(); err != nil {
		return result	}
	cur.Close(ctx)
	return result
}

func findChildren(id string, categories []models.Category) []models.Category {
	var result = []models.Category{}

	for i := 0; i < len(categories) ; i++ {
		var cat = categories[i]
		if cat.ParentID == id{
			result =  append(result, cat)
		}
	}
	return result
}

func(c *QuestionController) Index() revel.Result{


	var categories = GetAllCat()
	data := make(map[string]interface{})

	for _, cat := range categories {
		data[cat.ID.String()] = findChildren(cat.ID.String(),categories)

	}
	fmt.Println(data)
	return c.Render(categories)
}

func (c *QuestionController) GetAll() revel.Result {
	defer c.Request.Destroy()
	var result = []models.Question{}
	ctx := context.Background()
	cur, err := database.QuestionCollection.Find(ctx, bson.D{})
	c.Response.Status = http.StatusInternalServerError
	data:= make(map[string]interface{})
	if err != nil {
		data["status"] = "error"
		data["data"] = "Internal Server Error"
		return c.RenderJSON(data)
	}

	for cur.Next(ctx){
		var cat models.Question
		if err:= cur.Decode(&cat); err !=nil{
			data["status"] = "error"
			data["data"] = "Internal Server Error"
			return c.RenderJSON(data)
		}
		result = append(result, cat)
	}
	if err := cur.Err(); err != nil {
		data["status"] = "error"
		data["data"] = "Internal Server Error"
		return c.RenderJSON(data)
	}
	cur.Close(ctx)
	if len(result) == 0 {
		data["status"] = "error"
		data["data"] = mongo.ErrNoDocuments
		return c.RenderJSON(data)
	}
	c.Response.Status = http.StatusOK
	data["status"] = "success"
	data["data"] = result

	return c.RenderJSON(data)
}

func (c QuestionController) Create() revel.Result {
	defer c.Request.Destroy()
	var cat = &models.Question{}
	err := json.NewDecoder(c.Request.GetBody()).Decode(&cat)
	c.Response.Status = http.StatusBadRequest
	data := make(map[string]interface{})
	ctx := context.Background()
	data["status"] = "Error"

	if err != nil{
		data["data"] = "Status Bad Request. Could not decode Question"
		c.RenderJSON(data)
	}

	_, err = database.QuestionCollection.InsertOne(ctx, cat)
	ctx.Done()
	if err !=nil{
		data["status"] = "error"
		data["data"] = "Could insert user"
		return c.RenderJSON(data)
	}

	c.Response.Status = http.StatusCreated
	data["status"] = "success"
	data["data"] = cat

	return c.RenderJSON(data)

}
func (c App) List() revel.Result {
	response := map[string]interface{}{
		"Status": "succeeded",
		"data": todos,
	}

	return c.RenderJSON(response)
}

func (c App) Add(text string) revel.Result {
	if len(text) == 0 {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]string{ "Status": "failed" })
	}

	todos = append(todos, Todo { Text: text, parentID: text })

	response := map[string]interface{}{
		"Status": "succeeded",
		"data": todos,
	}

	return c.RenderJSON(response)
}

func (c App) Remove(index int) revel.Result {
	todos = append(todos[:index], todos[index+1:]...)

	response := map[string]interface{}{
		"Status": "succeeded",
		"data": todos,
	}

	return c.RenderJSON(response)
}

func (c App) RemoveAll() revel.Result {
	todos = nil

	response := map[string]interface{}{
		"Status": "succeeded",
	}

	return c.RenderJSON(response)
}


