package main

import (
	"books/m/v2/models"
	"books/m/v2/storage"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct{

	Author string       `json:"author" `
	Title  string        `json:"title"`
	Publisher string      `json:"publisher"`

}


type Repository struct {

	DB *gorm.DB
}



func (r *Repository)CreateBook(context *fiber.Ctx) {
	book := Book{}

err :=	context.BodyParser(&book)

if err != nil{
	context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message":"Request Failed"})
	//return err
}
err = r.DB.Create(&book).Error
if err != nil {

	context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message":"could not create book "})

	//return err
}
context.Status(http.StatusOK).JSON(&fiber.Map{"message":"Created Book sucessfully"})

//return nil


}

func (r *Repository) GetBooks(context *fiber.Ctx) {
bookModels := &[]models.Books{}

err := r.DB.Find(bookModels).Error

if err != nil {
	context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message":"could not get books"})

	//return err
}

context.Status(http.StatusOK).JSON(&fiber.Map{"message":"books fetched sucessfully ","data":bookModels})

//return nil
 
}


func (r *Repository) DeleteBook(context *fiber.Ctx)  {

	bookModel := models.Books{

	} 

	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":"Id cannot be empty"})
		//return nil
	}

	err := r.DB.Delete(bookModel,id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"messsage" : "Could not Delete Book" })
	  //return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"messsage" : "Book Deleted Sucessfully"}) 

		//return nil
}

func  (r *Repository) GetBookByID(context *fiber.Ctx)   {

	bookModel := models.Books{}
    id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":"Id cannot be empty"})
		//return nil
	}

	err := r.DB.Where("id = ?",id).First(bookModel).Error

	if err.Error != nil{
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not get Book"})
		//return err
	}

	

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"messsage" : "Book retrived Sucessfully",
		"data": bookModel,}) 

		//return nil

	
}


func (r * Repository)SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	api.Post("/create_book",r.CreateBook)
	api.Delete("delete_book:id", r.DeleteBook)
	api.Get("/get_books/:id" , r.GetBookByID)
	api.Get("/books", r.GetBooks)


}

func main(){

	err := godotenv.Load(".env")

	if err != nil{
		log.Fatal(err)
	}
      config := &storage.Config{
		Host:  os.Getenv("DB_HOST"),
		Port :   os.Getenv("DB_PORT"),
		Password:  os.Getenv("DB_PASSWORD"),
		User :     os.Getenv("DB_USER"),
		DBName:    os.Getenv("DB_NAME"),
		SSLMode:   os.Getenv("DB_SSL"),

	  }
	db,err := storage.NewConnection(config)

	if err != nil{
		log.Fatal("Could not load the data base")
	}

	err = models.MigrateBooks(db)

	if err != nil {
		log.Fatal("Could not migrate DB ")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
    r.SetupRoutes(app)

	app.Listen(":8080")


}