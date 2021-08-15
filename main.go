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

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var STATIC_UPLOAD_PATH string = "uploads" //cannot be changed
var CUR_UPLOAD_PATH string = "uploads"    //may update by user

type BasicFileInfo struct {
	Name string `json:"Name"`
	Size int64  `json:"Size"`
	Time string `json:"Time"`
}

// goto folder by absolute path
func Go_abs_Path(ctx *gin.Context) {
	pathname := ctx.PostForm("pathname")
	CUR_UPLOAD_PATH = path.Join(STATIC_UPLOAD_PATH, pathname)
	fmt.Println("Goto abs path", CUR_UPLOAD_PATH)
}

// goto subfolder
func Go_To_Path(ctx *gin.Context) {
	subfolder := ctx.PostForm("subfolder")
	CUR_UPLOAD_PATH = path.Join(CUR_UPLOAD_PATH, subfolder)
	fmt.Println("Goto subfolder ", CUR_UPLOAD_PATH)
}

func Go_Back(ctx *gin.Context) {
	if CUR_UPLOAD_PATH == STATIC_UPLOAD_PATH {
		return
	}
	CUR_UPLOAD_PATH = filepath.Dir(CUR_UPLOAD_PATH)
	CUR_UPLOAD_PATH = strings.ReplaceAll(CUR_UPLOAD_PATH, "\\", "/")
	fmt.Println("Back to parent folder", CUR_UPLOAD_PATH)
}

func UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
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
	session := sessions.Default(ctx)
	current_path := strings.TrimPrefix(CUR_UPLOAD_PATH, STATIC_UPLOAD_PATH)
	ctx.JSON(http.StatusOK, gin.H{
		"file":         data,
		"folder":       folder,
		"current_path": current_path,
		"permission":   session.Get("permission"),
	})
}

func DownloadFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	targetPath := filepath.Join(CUR_UPLOAD_PATH, filename)
	log.Println("DownloadFile@", targetPath)
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
			f, err := ar.Create(strings.TrimPrefix(strings.ReplaceAll(path, "\\", "/"), CUR_UPLOAD_PATH))
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

func RenamePath(ctx *gin.Context) {
	oldname := ctx.PostForm("oldname")
	newname := ctx.PostForm("newname")
	oldLocation := path.Join(CUR_UPLOAD_PATH, oldname)
	f, err := os.Stat(oldLocation)
	if err != nil {
		log.Fatal("Err on filename ", oldLocation, err, newname)
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		parent := filepath.Dir(oldLocation)
		newLocation := path.Join(parent, newname)
		err := os.Rename(oldLocation, newLocation)
		if err != nil {
			log.Fatal("Error when rename", err)
			ctx.String(http.StatusBadRequest, err.Error())
		}
		return
	case mode.IsRegular():
		extension := filepath.Ext(f.Name())
		parent := filepath.Dir(oldLocation)
		newLocation := path.Join(parent, newname+extension)
		log.Printf("The file name is %s and its extension is %s", parent, extension)
		err := os.Rename(oldLocation, newLocation)
		if err != nil {
			log.Fatal("Error when rename", err)
			ctx.String(http.StatusBadRequest, err.Error())
		}
		return
	}
}

func MoveToFolder(ctx *gin.Context) {
	filename := ctx.PostForm("filename")
	foldername := ctx.PostForm("foldername")
	fileLocation := path.Join(CUR_UPLOAD_PATH, filename)
	if foldername == "$-parent-$" {
		if CUR_UPLOAD_PATH == STATIC_UPLOAD_PATH {
			ctx.String(http.StatusBadRequest, "不可以將檔案移動至根目錄之外")
			return
		} else {
			parent := filepath.Dir(CUR_UPLOAD_PATH)
			newLocation := path.Join(parent, filename)
			os.Rename(fileLocation, newLocation)
			return
		}
	}
	folderLocation := path.Join(CUR_UPLOAD_PATH, foldername)
	f, err := os.Stat(fileLocation)
	if err != nil {
		log.Fatal("Err on filename ", fileLocation, err)
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		{
			newLocation := path.Join(folderLocation, filename)
			os.Rename(fileLocation, newLocation)
			return
		}
	case mode.IsRegular():
		{
			_, err = os.Stat(folderLocation)
			if err != nil {
				log.Fatal("Err on foldername ", folderLocation, err)
				return
			}
			newLocation := path.Join(folderLocation, filename)
			err = os.Rename(fileLocation, newLocation)
			if err != nil {
				log.Fatal("Error when rename", err)
				ctx.String(http.StatusBadRequest, err.Error())
			}
		}
	}

}

func IndexPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func AuthRequired(ctx *gin.Context) {
	session := sessions.Default(ctx)
	permission := session.Get("permission")
	if permission == nil {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}
	ctx.Next()
}

// Check the correctness of guuid and give the permission
func Login(ctx *gin.Context) {
	session := sessions.Default(ctx)
	guuid := ctx.PostForm("guuid")
	if guuid == "" {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}
	log.Println("Get guuid", guuid)
	switch guuid {
	case "alankingdom":
		session.Set("permission", "admin")
	case "visitor":
		session.Set("permission", "visitor")
	default:
		{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": "無法識別的GUUID",
			})
			log.Println("Get guuid", "無法識別的GUUID")
			return
		}
	}
	session.Options(sessions.Options{
		MaxAge: 3600 * 24, // 24hrs session
	})
	if err := session.Save(); err != nil {
		session.Delete("permission")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
	}
	ctx.Redirect(http.StatusFound, "/")
}

// Show Login Page
func Login_Render(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", nil)
}

func engine() *gin.Engine {
	r := gin.New()
	r.MaxMultipartMemory = 2 << 32
	r.Static("/views", "./views")
	r.LoadHTMLGlob("views/layouts/*")

	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/login", Login)
	r.GET("/login", Login_Render)

	r.Use(AuthRequired)
	{
		r.GET("/", IndexPage)
		r.POST("/upload", UploadFile)
		r.POST("/ChangePath", Go_To_Path)
		r.POST("/Go_abs_Path", Go_abs_Path)
		r.GET("/Go_Back", Go_Back)
		r.GET("/ls", ListFile)
		r.GET("downloads/:filename", DownloadFile)
		r.POST("/delete/:filename", DeleteFile)
		r.GET("/create/:foldername", CreateFolder)
		r.POST("/rename", RenamePath)
		r.POST("/movetofolder", MoveToFolder)
	}
	return r
}

func main() {
	router := engine()
	router.Use(gin.Logger())
	if err := engine().Run("0.0.0.0:5000"); err != nil {
		log.Fatal("無法啟動伺服器", err)
	}
}
