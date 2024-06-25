package main

import (
	"bufio"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/zh-blogs/blog-backup/library/config"
	"github.com/zh-blogs/blog-backup/library/zlog"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
	"strconv"
)

const ApiUrl = "https://zhblogs.ohyee.cc/api/blogs"
const DataDir = "./database"
const DataFile = "blogs.dat"
const PerPage = 500

func main() {
	// initialize
	config.Init()
	zlog.Init()

	zlog.L.Info("start backup database")

	// 检查目录是否存在
	if _, err := os.Stat(DataDir); os.IsNotExist(err) {
		if err := os.Mkdir(DataDir, os.ModePerm); err != nil {
			zlog.L.Fatal("create directory failed", zap.Error(err))
		}
	}

	// 创建临时文件
	tmpFile := path.Join(DataDir, DataFile+".tmp")
	file, err := os.Create(tmpFile)
	if err != nil {
		zlog.L.Fatal("create file failed", zap.Error(err))
	}

	defer func() {
		_ = file.Close()
	}()

	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) ZHBlog-Database-Backup/1.0",
	}

	// 分批读取，避免超时以及过高的内存占用
	var offset = 0
	var Bls Blogs
	for {
		Query := map[string]string{
			"search": "",
			"tags":   "",
			"offset": strconv.Itoa(offset),
			"size":   strconv.Itoa(PerPage),
		}
		_, err := resty.New().R().SetHeaders(header).
			SetResult(&Bls).
			SetQueryParams(Query).
			Get(ApiUrl)
		if err != nil {
			zlog.L.Fatal("get data failed", zap.Error(err))
		}

		if Bls.Success == false {
			zlog.L.Fatal("get data failed", zap.Error(err))
		}

		// 写入文件
		for _, blog := range Bls.Data.Blogs {
			jsonData, _ := json.Marshal(blog)
			_, _ = file.WriteString(string(jsonData) + "\n")
		}

		zlog.L.Info("get data", zap.Int("num", len(Bls.Data.Blogs)), zap.Int("now", offset+len(Bls.Data.Blogs)), zap.Int("total", Bls.Data.Total))

		if Bls.Data.Blogs == nil || len(Bls.Data.Blogs) < PerPage {
			break
		}

		offset += PerPage
		Bls = Blogs{}
	}

	_ = file.Close()

	if success := verify(tmpFile); !success {
		zlog.L.Fatal("verify failed")
	}

	// 校验文件
	file, err = os.Open(tmpFile)
	if err != nil {
		zlog.L.Fatal("open file failed", zap.Error(err))
	}

	if err := os.Rename(tmpFile, path.Join(DataDir, DataFile)); err != nil {
		zlog.L.Fatal("rename file failed", zap.Error(err))
	}

	zlog.L.Info("backup database success")
}

// verify 验证文件是否正确
func verify(filename string) bool {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return false
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	reader := bufio.NewReader(f)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		var blog BlogData
		if err := json.Unmarshal(line, &blog); err != nil {
			return false
		}
	}
	return true
}

type Blogs struct {
	Success bool `json:"success"`
	Data    struct {
		Total int `json:"total"`
		Blogs []BlogData
	} `json:"data"`
}

type BlogData struct {
	Id         string   `json:"id"`
	Idx        int      `json:"idx"`
	Name       string   `json:"name"`
	Url        string   `json:"url"`
	Tags       []string `json:"tags"`
	Sign       string   `json:"sign"`
	Feed       string   `json:"feed"`
	Status     string   `json:"status"`
	Repeat     bool     `json:"repeat"`
	Enabled    bool     `json:"enabled"`
	Sitemap    string   `json:"sitemap"`
	Arch       string   `json:"arch"`
	JoinTime   int64    `json:"join_time"`
	UpdateTime int64    `json:"update_time"`
	SaveWebId  string   `json:"saveweb_id"`
	Recommend  bool     `json:"recommend"`
}
