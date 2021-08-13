package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

var STATIC_UPLOAD_PATH string = "uploads" //cannot be changed
var CUR_UPLOAD_PATH string = "uploads"    //may update by user

type BasicFileInfo struct {
	Name string `json:"Name"`
	Size int64  `json:"Size"`
	Time string `json:"Time"`
}

func Go_To_Path(ctx *gin.Context) {
	subfolder := ctx.PostForm("subfolder")
	CUR_UPLOAD_PATH = path.Join(CUR_UPLOAD_PATH, subfolder)
	fmt.Println(CUR_UPLOAD_PATH)
}

func Go_Back(ctx *gin.Context) {
	CUR_UPLOAD_PATH = filepath.Dir(CUR_UPLOAD_PATH)
	fmt.Println(CUR_UPLOAD_PATH)
}

func UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	err = ctx.SaveUploadedFile(file, filepath.Join(CUR_UPLOAD_PATH, file.Filename))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.String(http.StatusOK, "upload successful \n")

}

func ListFile(ctx *gin.Context) {
	files, err := ioutil.ReadDir(CUR_UPLOAD_PATH)
	// fmt.Println("Files are", files)
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	data := []BasicFileInfo{}
	folder := []BasicFileInfo{}
	for _, f := range files {
		tmp := BasicFileInfo{
			Name: f.Name(),
			Size: f.Size(),
			Time: f.ModTime().Local().String()}
		if f.IsDir() {
			folder = append(folder, tmp)
		} else {
			data = append(data, tmp)
		}
	}
	current_path := strings.TrimPrefix(CUR_UPLOAD_PATH, STATIC_UPLOAD_PATH)
	ctx.JSON(http.StatusOK, gin.H{
		"file":         data,
		"folder":       folder,
		"current_path": current_path,
	})
}

func DownloadFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	targetPath := filepath.Join(CUR_UPLOAD_PATH, filename)
	f, err := os.Stat(targetPath)
	if err != nil {
		log.Fatal("Cannot find file error: ", err)
		return
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		// If the given path is a directory, then recursively zip and stream to client side.
		new_zip_name := fmt.Sprintf("%s.zip", filename)
		value := fmt.Sprintf("attachment; filename=%s", new_zip_name)
		ctx.Writer.Header().Set("Content-type", "application/octet-stream")
		ctx.Writer.Header().Set("Content-Disposition", value)
		ar := zip.NewWriter(ctx.Writer)
		walker := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			// Trim the prefix string to re-produce the same directory
			base_path_prefix := fmt.Sprintf("%s\\", CUR_UPLOAD_PATH)
			f, err := ar.Create(strings.TrimPrefix(path, base_path_prefix))
			if err != nil {
				return err
			}

			_, err = io.Copy(f, file)
			if err != nil {
				return err
			}
			return nil
		}
		err = filepath.Walk(targetPath, walker)
		if err != nil {
			panic(err)
		}
		ar.Close()

	case mode.IsRegular():
		ctx.Header("Content-Description", "File Transfer")
		ctx.Header("Content-Transfer-Encoding", "binary")
		ctx.Header("Content-Disposition", "attachment; filename="+filename)
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.File(targetPath)
	}
}

func DeleteFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	targetPath := filepath.Join(CUR_UPLOAD_PATH, filename)
	err := os.RemoveAll(targetPath)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"info": "Success",
	})
}

func CreateFolder(ctx *gin.Context) {
	dir_name := ctx.Param("foldername")
	targetPath := filepath.Join(CUR_UPLOAD_PATH, dir_name)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		err := os.MkdirAll(targetPath, os.ModePerm)
		if err != nil {
			log.Fatal("There is error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
	}
	log.Println("Create success")
}

func MoveTo(ctx *gin.Context) {
	oldLocation := ctx.PostForm("oldLocation")
	newLocation := ctx.PostForm("newLocation")
	f, err := os.Stat(oldLocation)
	if err != nil {
		log.Fatal("err is ", err)
		return
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		return
	case mode.IsRegular():
		return
	}
	fmt.Println(oldLocation, newLocation)
	err = os.Rename(oldLocation, newLocation)
	if err != nil {
		log.Fatal(err)
	}
}

func IndexPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 2 << 32
	router.Static("/views", "./views")
	router.LoadHTMLGlob("views/layouts/*")

	router.GET("/", IndexPage)
	router.POST("/", UploadFile)
	router.POST("ChangePath", Go_To_Path)
	router.GET("/Go_Back", Go_Back)
	router.GET("/ls", ListFile)
	router.GET("downloads/:filename", DownloadFile)
	router.POST("/delete/:filename", DeleteFile)
	router.GET("/create/:foldername", CreateFolder)
	router.POST("/move", MoveTo)
	router.Run(":5000")
}
