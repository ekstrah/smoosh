package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	// Mongodb Setup Finished
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.LoadHTMLGlob("templates/*")
	router.GET("/", serveUploadPage)
	router.POST("/upload", uploadHandler)
	router.Run(":8080")
}

// 다중 파일 업로드 폼 HTML
const uploadPage string = `<html>
<head>
	<title>파일 업로드</title>
	<meta charset="utf-8">
</head>
<body>
	<h2>파일 업로드</h2>
	<form action="/upload" method="POST" enctype="multipart/form-data">
		upload file: <input type="file" name="upload[]" multiple>
		<input type="submit" value="Submit">
	</form>
</body>
</html>`

// 업로드 폼 제공 핸들러
func serveUploadPage(c *gin.Context) {
	c.Data(http.StatusOK, "text/html", []byte(uploadPage))
}

// 다중 파일 업로드 핸들러
func uploadHandler(c *gin.Context) {
	// Multipart form
	form, _ := c.MultipartForm()

	// Uploaded files
	files := form.File["upload[]"]

	// Setting up MongoDB Client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary()) // Primary DB에 대한 연결 체크

	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database("filer").Collection("userTest")
	// for range 로 업로드한 파일 순회
	for _, file := range files {
		log.Println(file.Filename)
		c.SaveUploadedFile(file, filepath.Join("./uploaded", file.Filename))
		doc := bson.D{{"fileName", file.Filename}, {"filePath", filepath.Join("./uploaded", file.Filename)}}
		result, _ := coll.InsertOne(context.TODO(), doc)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)

	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	fmt.Println(files)
}
