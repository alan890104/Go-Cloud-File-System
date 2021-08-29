package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fileserver/packages/dbfunc"
	"fileserver/packages/utils"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type BasicFileInfo struct {
	Name string `json:"Name"`
	Size int64  `json:"Size"`
	Time string `json:"Time"`
}

/* ####################Configuration##########################*/

var Config = utils.InitConfig()
var TRASH_PATH = "trash"           //Path of trash
var ROOT_UPLOAD_PATH = "uploads"   //cannot be changed
var WhiteList = Config.IPWhiteList // A set of whitelist
var SERVER_IP = fmt.Sprint(utils.GetOutboundIP())

/* ####################Configuration##########################*/

// This is not a route. Return "" if no CUR_UPLOAD_PATH in session
func Load_CUR_UPLOAD_PATH(ctx *gin.Context) string {
	session := sessions.Default(ctx)
	v := session.Get("CUR_UPLOAD_PATH")
	if v == nil {
		return ""
	}
	return v.(string)
}

// This is not a route
func Set_CUR_UPLOAD_PATH(ctx *gin.Context, newpath string) {
	session := sessions.Default(ctx)
	session.Set("CUR_UPLOAD_PATH", newpath)
	session.Options(sessions.Options{
		MaxAge: 0,
	})
	session.Save()
}

// goto folder by absolute path
func Go_abs_Path(ctx *gin.Context) {
	pathname := ctx.PostForm("pathname")
	CUR_UPLOAD_PATH := path.Join(ROOT_UPLOAD_PATH, pathname)
	Set_CUR_UPLOAD_PATH(ctx, CUR_UPLOAD_PATH)
	fmt.Println("Goto abs path", CUR_UPLOAD_PATH)
}

// goto subfolder
func Go_To_Path(ctx *gin.Context) {
	subfolder := ctx.PostForm("subfolder")
	CUR_UPLOAD_PATH := Load_CUR_UPLOAD_PATH(ctx)
	CUR_UPLOAD_PATH = path.Join(CUR_UPLOAD_PATH, subfolder)
	Set_CUR_UPLOAD_PATH(ctx, CUR_UPLOAD_PATH)
	fmt.Println("Goto subfolder ", CUR_UPLOAD_PATH)
}

func Go_Back(ctx *gin.Context) {
	CUR_UPLOAD_PATH := Load_CUR_UPLOAD_PATH(ctx)
	if CUR_UPLOAD_PATH == ROOT_UPLOAD_PATH {
		return
	}
	CUR_UPLOAD_PATH = filepath.Dir(CUR_UPLOAD_PATH)
	CUR_UPLOAD_PATH = strings.ReplaceAll(CUR_UPLOAD_PATH, "\\", "/")
	Set_CUR_UPLOAD_PATH(ctx, CUR_UPLOAD_PATH)
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
	CUR_UPLOAD_PATH := Load_CUR_UPLOAD_PATH(ctx)
	err = ctx.SaveUploadedFile(file, filepath.Join(CUR_UPLOAD_PATH, file.Filename))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.String(http.StatusOK, "upload successful \n")

}

func ListFile(ctx *gin.Context) {
	CUR_UPLOAD_PATH := Load_CUR_UPLOAD_PATH(ctx)
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
	current_path := strings.TrimPrefix(CUR_UPLOAD_PATH, ROOT_UPLOAD_PATH)
	ctx.JSON(http.StatusOK, gin.H{
		"file":         data,
		"folder":       folder,
		"current_path": current_path,
		"permission":   session.Get("permission"),
	})
}

func DownloadFile(ctx *gin.Context) {
	filename := ctx.Param("filename")
	CUR_UPLOAD_PATH := Load_CUR_UPLOAD_PATH(ctx)
	targetPath := filepath.Join(CUR_UPLOAD_PATH, filename)
	log.Println("DownloadFile@", targetPath)
	f, err := os.Stat(targetPath)
	if err != nil {
		log.Println("Cannot find file error: ", err)
		ctx.String(http.StatusInternalServerError, err.Error())
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

// Not acutally delete file, move the file to trash instead
func DeleteFile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if permission := session.Get("permission"); permission != "admin" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "No permission",
		})
		return
	}
	filename := ctx.Param("filename")
	mode := ctx.PostForm("mode")
	switch mode {
	case "delete":
		{
			targetPath := filepath.Join(TRASH_PATH, filename)
			log.Printf("Delete file mode: delete, path:%s", targetPath)
			err := os.RemoveAll(targetPath)
			if err != nil {
				log.Println("Error when delete file (delete mode)", err)
				return
			}
			err = dbfunc.DeleteHistory(targetPath)
			if err != nil {
				log.Println("Error when delete history (delete mode)", err)
				return
			}
		}
	default:
		{
			CUR_UPLOAD_PATH := Load_CUR_UPLOAD_PATH(ctx)
			targetPath := filepath.Join(CUR_UPLOAD_PATH, filename)
			log.Printf("Delete file mode: trash, path:%s, parent: %s", targetPath, filepath.Dir(targetPath))
			newPath := path.Join(TRASH_PATH, filename)
			err := os.Rename(targetPath, newPath)
			if err != nil {
				log.Println("Error when delete file (trash mode)", err)
				return
			}
			dbfunc.InsertHistory(targetPath, newPath)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"info": "Success",
	})
}

func TrashFilesRender(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("permission") != "admin" {
		log.Println("No permission to access ListTrash")
		ctx.Redirect(http.StatusPermanentRedirect, "/")
		// ctx.AbortWithStatus(http.StatusPermanentRedirect)
		return
	}
	ctx.HTML(http.StatusOK, "trash.html", nil)
}

func ListTrash(ctx *gin.Context) {
	session := sessions.Default(ctx)
	data := []BasicFileInfo{}
	folder := []BasicFileInfo{}
	if session.Get("permission") != "admin" {
		log.Println("No permission to access ListTrash")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "No permission to access ListTrash",
		})
		return
	}
	files, err := ioutil.ReadDir(TRASH_PATH)
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
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
	ctx.JSON(http.StatusOK, gin.H{
		"file":   data,
		"folder": folder,
	})
}

func Recover(ctx *gin.Context) {
	filename := ctx.Param("filename")
	targetPath := path.Join(TRASH_PATH, filename)
	recoverPath, err := dbfunc.FindHistory(targetPath)
	if err != nil {
		log.Println("Error when Find History", err)
		ctx.String(http.StatusBadRequest, "檔案歷史紀錄遺失，無法復原檔案")
		return
	}
	fmt.Printf("Recovery %s to %s\n", targetPath, recoverPath)
	err = os.Rename(targetPath, recoverPath)
	if err != nil {
		log.Println("Error when Recover Rename", err)
		parent_folder_name := filepath.Base(filepath.Dir(recoverPath))
		err_string := fmt.Sprintf("原始檔案的母目錄已經被刪除，無法復原。可以試著在回收桶中尋找名稱為%s的資料夾並將其還原。", parent_folder_name)
		ctx.String(http.StatusBadRequest, err_string)
		return
	}
	err = dbfunc.DeleteHistory(targetPath)
	if err != nil {
		log.Println("Error when delete history", err)
	}
}

func CreateFolder(ctx *gin.Context) {
	dir_name := ctx.Param("foldername")
	CUR_UPLOAD_PATH := Load_CUR_UPLOAD_PATH(ctx)
	targetPath := filepath.Join(CUR_UPLOAD_PATH, dir_name)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		err := os.MkdirAll(targetPath, os.ModePerm)
		if err != nil {
			log.Println("There is error", err)
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
	CUR_UPLOAD_PATH := Load_CUR_UPLOAD_PATH(ctx)
	oldLocation := path.Join(CUR_UPLOAD_PATH, oldname)
	f, err := os.Stat(oldLocation)
	if err != nil {
		log.Println("Err on filename ", oldLocation, err, newname)
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		parent := filepath.Dir(oldLocation)
		newLocation := path.Join(parent, newname)
		err := os.Rename(oldLocation, newLocation)
		if err != nil {
			log.Println("Error when rename", err)
			ctx.String(http.StatusBadRequest, err.Error())
		}
		return
	case mode.IsRegular():
		extension := filepath.Ext(f.Name())
		parent := filepath.Dir(oldLocation)
		newLocation := path.Join(parent, newname+extension)
		log.Printf("The file name is %s and its extension is %s\n", parent, extension)
		err := os.Rename(oldLocation, newLocation)
		if err != nil {
			log.Println("Error when rename", err)
			ctx.String(http.StatusBadRequest, err.Error())
		}
		return
	}
}

func MoveToFolder(ctx *gin.Context) {
	filename := ctx.PostForm("filename")
	foldername := ctx.PostForm("foldername")
	CUR_UPLOAD_PATH := Load_CUR_UPLOAD_PATH(ctx)
	fileLocation := path.Join(CUR_UPLOAD_PATH, filename)
	if foldername == "$-parent-$" {
		if CUR_UPLOAD_PATH == ROOT_UPLOAD_PATH {
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
		log.Println("Err on filename ", fileLocation, err)
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
				log.Println("Err on foldername ", folderLocation, err)
				return
			}
			newLocation := path.Join(folderLocation, filename)
			err = os.Rename(fileLocation, newLocation)
			if err != nil {
				log.Println("Error when rename", err)
				ctx.String(http.StatusBadRequest, err.Error())
			}
		}
	}

}

func FreeSpace(ctx *gin.Context) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println("Error on os.Getwd", err)
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	ds, err := utils.Disk_Space(dir)
	if err != nil {
		log.Println("Error on Disk_Space", err)
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	usedpercent := fmt.Sprintf("%.3f%%", utils.SpaceUsedPercent(ds)*100)
	ctx.JSON(http.StatusOK, gin.H{
		"percent": usedpercent,
		"use":     utils.ByteFormat(utils.SpaceUsed(ds)),
		"free":    utils.ByteFormat(ds.FreeByte),
	})
	// log.Println("Return disk usage percent :", usedpercent)
}

func IndexPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func SettingPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "setting.html", gin.H{
		"ADMIN_GUUID":     Config.ADMIN_GUUID,
		"StopDelaySecond": Config.StopDelaySecond,
		"WhiteList":       WhiteList,
	})
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

func IPRestriction(ctx *gin.Context) {
	if _, exist := WhiteList[ctx.ClientIP()]; exist {
		ctx.Next()
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "IP is not allowed",
		})
		ctx.Abort()
	}
}

// Check user permission
func CheckPermission(ctx *gin.Context) {
	session := sessions.Default(ctx)
	permission := session.Get("permission")
	switch permission {
	case "admin":
		ctx.Next()
	default:
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

// A PUT request to add ip whitelist
func AddIPRestriction(ctx *gin.Context) {
	allow_ip := ctx.Param("allow_ip")
	if allow_ip != "127.0.0.1" && allow_ip != SERVER_IP {
		WhiteList[allow_ip] = false
		ctx.Status(http.StatusOK)
		return
	}
	ctx.Status(http.StatusNotAcceptable)
}

// A DELETE request to delete ip whitelist
func RemoveIPRestriction(ctx *gin.Context) {
	remove_ip := ctx.Param("remove_ip")
	if remove_ip != "127.0.0.1" && remove_ip != SERVER_IP {
		delete(WhiteList, remove_ip)
		ctx.Status(http.StatusOK)
		return
	}
	ctx.Status(http.StatusNotAcceptable)
}

// Check the correctness of guuid and give the permission
func Login(ctx *gin.Context) {
	session := sessions.Default(ctx)
	guuid := ctx.PostForm("guuid")
	if guuid == "" {
		log.Println("Found guuid empty, redirect to /login")
		ctx.Redirect(http.StatusFound, "/login")
		return
	}
	log.Println("Get guuid", guuid)
	switch guuid {
	case Config.ADMIN_GUUID:
		session.Set("permission", "admin")
	case "visitor":
		session.Set("permission", "visitor")
	default:
		{
			ctx.String(http.StatusUnauthorized, "無法識別的GUUID")
			log.Println("Get guuid", "無法識別的GUUID")
			return
		}
	}
	session.Set("CUR_UPLOAD_PATH", ROOT_UPLOAD_PATH) // personal path
	session.Options(sessions.Options{
		MaxAge: 0,
	})
	session.Save()
	if err := session.Save(); err != nil {
		session.Clear()
		session.Options(sessions.Options{MaxAge: -1})
		session.Save()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func GetSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"permission": session.Get("permission"),
	})
}

// Show Login Page
func Login_Render(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", nil)
}

func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	log.Println("User logout.")
	ctx.Redirect(http.StatusMovedPermanently, "/login")
}

func engine() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.MaxMultipartMemory = 2 << 32
	r.Static("/views", "./views")
	r.LoadHTMLGlob("views/layouts/*")

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	// r.Use(gin.Logger()) //Show detail of requests

	low_level := r.Group("/").Use(IPRestriction)
	mid_level := r.Group("/").Use(IPRestriction, AuthRequired)
	high_level := r.Group("/").Use(IPRestriction, CheckPermission)

	low_level.POST("/login", Login)
	low_level.GET("/login", Login_Render)
	low_level.POST("/logout", Logout)

	mid_level.GET("/", IndexPage)
	mid_level.POST("/upload", UploadFile)
	mid_level.POST("/ChangePath", Go_To_Path)
	mid_level.POST("/Go_abs_Path", Go_abs_Path)
	mid_level.GET("/Go_Back", Go_Back)
	mid_level.GET("/ls", ListFile)
	mid_level.GET("downloads/:filename", DownloadFile)
	mid_level.POST("/delete/:filename", DeleteFile)
	mid_level.GET("/trash", TrashFilesRender)
	mid_level.POST("/trash/list", ListTrash)
	mid_level.POST("recover/:filename", Recover)
	mid_level.GET("/create/:foldername", CreateFolder)
	mid_level.POST("/rename", RenamePath)
	mid_level.POST("/movetofolder", MoveToFolder)
	mid_level.GET("/freespace", FreeSpace)
	mid_level.GET("/session", GetSession)

	// Only allow_ip and admin can access the route
	high_level.GET("/setting", SettingPage)
	high_level.PUT("/AddIPRestriction/:allow_ip", AddIPRestriction)
	high_level.DELETE("/RemoveIPRestriction/:remove_ip", RemoveIPRestriction)
	return r
}

// SaveConfig when exit program
func SaveConfig() error {
	delete(WhiteList, SERVER_IP)
	Config.IPWhiteList = WhiteList
	file, err := json.MarshalIndent(Config, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.json", file, 0644)
	if err != nil {
		return errors.New("在儲存config.json時發生錯誤")
	}
	return nil
}

func main() {
	dbfunc.InitializeDB()
	router := engine()
	router.Use(gin.Logger())
	WhiteList[SERVER_IP] = true
	log.Println("您的GUUID為", Config.ADMIN_GUUID)
	log.Printf("將於 http://%s:5000 開啟伺服器", SERVER_IP)
	if dir, err := os.Getwd(); err != nil {
		log.Fatal("工作目錄錯誤: ", err)
		os.Exit(1)
	} else {
		log.Println("檔案根目錄位於 ", dir+"\\"+ROOT_UPLOAD_PATH)
	}
	/* Use graceful shutdown */
	srv := &http.Server{
		Addr:    ":5000",
		Handler: router,
	}
	// Create channel
	ch := make(chan os.Signal, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("無法啟動伺服器", err)
		}
	}()
	// Notify: Send system signal to channel
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	// block channel
	<-ch
	log.Printf("關閉伺服器-緩衝時間%d秒\n", Config.StopDelaySecond)
	ctx := context.Background()
	c, cancel := context.WithTimeout(ctx, time.Second*time.Duration(Config.StopDelaySecond))
	defer cancel()
	if err := srv.Shutdown(c); err != nil {
		log.Println("強制關閉伺服器時發生錯誤:", err)
	}
	// close channel
	<-c.Done()
	close(ch)
	if err := SaveConfig(); err != nil {
		log.Println(err)
	}
	log.Println("已經完全關閉伺服器")
}
