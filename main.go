package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
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

	// for range 로 업로드한 파일 순회
	for _, file := range files {
		log.Println(file.Filename)
		c.SaveUploadedFile(file, filepath.Join("./uploaded", file.Filename))
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}
