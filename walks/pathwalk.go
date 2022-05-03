package walks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	path2 "path"
	"strings"
)

var dir, pan string
var Pre bool

type FileInfo struct {
	Name  string
	Size  int64
	Bytes []byte
}

// PreFile 图片简单压缩
func PreFile(name, kind string) []byte {
	if name == "" {
		return nil
	}
	if kind == "png" || kind == "jpg" || kind == "gif" || kind == "mp4" || kind == "txt" {
		if !Pre {
			if kind == "txt" {
			} else {
				fmt.Print(name, "no prepare  ,")
				return []byte("hello")
			}
		}
	} else {
		return nil
	}
	file, err := os.Open(dir + `/` + name)
	if err != nil {
		fmt.Println(err)
		return []byte("error")
	}
	if kind == "jpg" {
		ima, err := jpeg.Decode(file)
		if err != nil {
			fmt.Println(err)
			return []byte("error")
		}
		var b bytes.Buffer
		err = jpeg.Encode(&b, ima, &jpeg.Options{Quality: 5})
		if err != nil {
			fmt.Println(err)
			return []byte("error")
		}

		return b.Bytes()
	}
	if kind == "gif" {
		gifs, err := gif.DecodeAll(file)
		if err != nil {
			return []byte("error")
		}
		var b bytes.Buffer
		err = jpeg.Encode(&b, gifs.Image[0], &jpeg.Options{Quality: 5})
		if err != nil {
			fmt.Println(err)
			return []byte("error")
		}
		return b.Bytes()
	}
	if kind == "png" {
		ima, err := png.Decode(file)
		if err != nil {
			fmt.Println(err)
			return []byte("error")
		}
		var b bytes.Buffer
		err = jpeg.Encode(&b, ima, &jpeg.Options{Quality: 5})
		if err != nil {
			fmt.Println(err)
			return []byte("error")
		}
		return b.Bytes()
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
	if kind == "txt" {
		return []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 32, 0, 0, 0, 32, 8, 6, 0, 0, 0, 115, 122, 122, 244, 0, 0, 0, 151, 73, 68, 65, 84, 120, 156, 236, 150, 81, 10, 128, 48, 12, 67, 163, 120, 6, 17, 239, 127, 36, 81, 60, 141,
			226, 159, 31, 117, 75, 187, 66, 153, 52, 80, 252, 41, 245, 153, 198, 177, 17, 193, 74, 128, 18, 192, 101, 168, 29, 192, 226, 5, 103, 1, 120, 234, 4, 176, 122, 2, 88, 128, 93, 32, 172, 0, 199, 235, 217, 180, 14, 43, 192, 12, 96, 243, 112,
			194, 10, 0, 47, 136, 150, 12, 72, 127, 135, 168, 161, 50, 176, 214, 35, 245, 171, 222, 53, 145, 195, 25, 125, 129, 22, 193, 194, 79, 66, 173, 3, 108, 38, 216, 181, 245, 231, 0, 253, 101, 172, 186, 115, 32, 51, 144, 25, 200, 12, 252, 47,
			3, 9, 192, 100, 64, 115, 43, 82, 43, 220, 129, 112, 133, 59, 16, 14, 112, 7, 0, 0, 255, 255, 30, 221, 72, 54, 38, 7, 188, 168, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130,
		}
	}
	if kind == "mp4" {
		width := 640
		height := 360
		filename := dir + `/` + name
		cmd := exec.Command("ffmpeg", "-i", filename, "-vframes", "1", "-s", fmt.Sprintf("%dx%d", width, height), "-f", "singlejpeg", "-")
		var buffer bytes.Buffer
		cmd.Stdout = &buffer
		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
			return []byte("error")
		}
		return buffer.Bytes()
	}
	return nil
}

// DirMap 遍历文件返回json格式切片,Pre 为true预览图片
func DirMap(fileInfos []fs.FileInfo) []byte {
	var dirMap = make(map[string][]FileInfo)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			dirMap["dir"] = append(dirMap["dir"], FileInfo{Name: fileInfo.Name()})
			continue
		}
		if path2.Ext(fileInfo.Name()) == "" {
			dirMap["none"] = append(dirMap["none"], FileInfo{Name: fileInfo.Name()})
			continue
		}
		name := strings.TrimLeft(path2.Ext(fileInfo.Name()), ".")
		b := PreFile(fileInfo.Name(), name)
		dirMap[name] = append(dirMap[name], FileInfo{Name: fileInfo.Name(), Bytes: b, Size: fileInfo.Size()})
	}
	dir += `/`
	fmt.Println("path:", dir)
	j, err := json.Marshal(dirMap)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return j
}

//Pan 访问盘,Pre 为true预览图片
func Pan(path string) []byte {
	fileInfos, err := ioutil.ReadDir(path)
	pan = path
	dir = path[:2]
	fmt.Println("pan:", pan)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return DirMap(fileInfos)
}

// Walk 访问文件夹,Pre 为true预览图片
func Walk(path string) []byte {
	if path == ".." {
		if dir == pan {
			fmt.Println("back to pan:", dir)
			return Pan(pan)
		} else {
			dir = path2.Dir(path2.Dir(dir))
			if dir+`/` == pan {
				dir = pan
				return Pan(pan)
			}
		}
	} else {
		dir += path
	}
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return DirMap(fileInfos)
}

// WalkFile 写入图片和视频
func WalkFile(name string, w http.ResponseWriter) {
	if name == "" {
		_, _ = w.Write(nil)
		return
	}
	file, err := os.Open(dir + name)
	fileInfo, err := file.Stat()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Length", fmt.Sprint(fileInfo.Size()))
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write(nil)
		return
	} else {
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(file)
	}
	n, err := io.Copy(w, file)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write(nil)
		return
	}
	fmt.Println("      ", "写入", n, "字节...")
}
