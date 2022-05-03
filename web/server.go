package web

import (
	"browse/walks"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var (
	DefaultURL = "http://localhost:8080/handler"
	Addr       = "localhost:8080"
)

//Post 设置传入模板的数据
type Post struct {
	URL string
}

type Pan int

func (p Pan) String() string {
	if p > 4 || p < 1 {
		return ""
	}
	return [...]string{`C:/`, `D:/`, `E:/`, `F:/`}[p-1]
}

type File struct {
	Jpg string
	Png string
	Gif string
	Txt string
}

type BigFile struct {
	Mp4 string
}
type Form struct {
	Test string
	Pan  Pan
	Back int
	Dir  string
	File
	BigFile
}

func (f File) RegFile() (name string) {
	switch {
	case f.Jpg != "":
		return f.Jpg
	case f.Gif != "":
		return f.Gif
	case f.Png != "":
		return f.Png
	case f.Txt != "":
		return f.Txt
	default:
		return ""
	}

}
func (f BigFile) RegBigFile() (name string) {
	switch {
	case f.Mp4 != "":
		{
			return f.Mp4
		}
	default:
		{
			return ""
		}
	}
}
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.RequestURI() == "/favicon.ico" {
		return
	}
	fmt.Println(r.URL)
	if r.Method != "GET" {
		fmt.Println("home----->\n      已拒绝" + r.RemoteAddr + "的" + r.Method + "访问")
		return
	}
	fmt.Println("home----->\n      " + r.RemoteAddr + ":已连接...")
	//加载模板
	temp, err := template.ParseFiles("../templates/start/home.html")
	if err != nil {
		fmt.Println("      传入数据失败!\n      " + err.Error())
		return
	}
	//加载数据
	err = temp.Execute(w, Post{URL: DefaultURL})
	if err != nil {
		fmt.Println("      传入数据失败!\n      " + err.Error())
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		_, _ = w.Write(nil)
		return
	}
	fmt.Println("handler-->\n      " + r.RemoteAddr + ":已连接...")
	var form = Form{}
	//解析请求的数据
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		fmt.Println("      " + err.Error())
	}
	fmt.Printf("      %#v\n", form)
	//请求处理
	fileHandler(form, w)
}

func fileHandler(form Form, w http.ResponseWriter) {
	//test处理

	if form.Test == "connectTest" {
		_, _ = w.Write([]byte("测试成功!"))
		return
	}
	//pan处理

	if form.Pan > 0 && form.Pan < 5 {
		w.Header().Set("content", "no-cache")
		n, err := w.Write(walks.Pan(form.Pan.String()))
		fmt.Println("      ", "写入", n, "字节...")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	//dir处理

	if form.Dir != "" {
		if strings.Contains(form.Dir, "..") {
			_, _ = w.Write(nil)
			return
		}
		w.Header().Set("content", "no-cache")
		n, err := w.Write(walks.Walk(form.Dir))
		fmt.Println("      ", "写入", n, "字节...")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	//back处理

	if form.Back == 1 {

		w.Header().Set("content", "no-cache")
		n, err := w.Write(walks.Walk(".."))
		fmt.Println("      ", "写入", n, "字节...")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	//图片处理

	if name := form.File.RegFile(); name != "" {
		walks.WalkFile(name, w)

	}
	//文件处理

	if name := form.BigFile.RegBigFile(); name != "" {
		w.Header().Set("Content-Length", "")
		walks.WalkFile(name, w)
		return
	}
}

//Listen addr为监听端口，为空默认localhost:8080
//Pre 为true加载缩略图
func Listen(addr string, pre bool) {
	if addr != "" {
		Addr = addr
		DefaultURL = `http://` + Addr + `/handler`
	}
	walks.Pre = pre
	http.HandleFunc("/", home)
	http.HandleFunc("/handler", handler)
	err := http.ListenAndServe(Addr, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
