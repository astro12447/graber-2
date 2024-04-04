package main

import (
	"container/list"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// Получение  пути из командной строки (cmd).
func GetPathFromCommandLine(src string, dst string) (string, string) {
	var sources *string
	var destination *string
	sources = flag.String(src, "None", "")
	destination = flag.String(dst, "None", "")
	flag.Parse()

	return *sources, *destination
}

// Получение формат url ->(https://HostDomainName) ссылки из данных файла
func readDataFromFile(fileName string) (string, error) {
	//читать строки файла, одну за другой в памяти
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Anable to Open File: %v", err)
	}
	defer f.Close()
	buf := make([]byte, 1024)
	var FileContent string
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Print(err)
		}
		if n > 0 {

			FileContent = string(buf[:n])
		} else {
			fmt.Println("Исходной файл пустый!")
		}
	}
	return FileContent, nil
}

// Формат домена(https://www.google.com/)
type Domain struct {
	schema   string
	host     string
	hostname string
	path     string
}

// Получение реальный URL-ссылки из url данных.
func isUrl(input string) bool {
	flag := false
	a, err := url.Parse(input)
	if err != nil {
		fmt.Println(input, "Пусто!..")
	}
	url := Domain{schema: a.Scheme,
		host:     a.Host,
		hostname: a.Hostname(),
		path:     a.Path}
	if url.schema == "" || url.host == "" {
		flag = false
	}
	if url.schema != "" && url.host != "" && url.hostname != "" && url.path != "" || url.path == "" {
		flag = true
	}
	return flag
}

// Получение действительный URL-адрес из файла данных
func getUrlFromfile(FileName string) (*list.List, error) {
	data, err := readDataFromFile(FileName)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(data, "\n")
	l := list.New()
	for _, item := range lines {
		if isUrl(item) == true && item != "" {
			l.PushBack(item)
		}
	}
	return l, nil
}

// Получение названия файла через имени домена
func GethostnameFromURL(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatal("URL given not correcly!", URL)
	}
	hostname := u.Hostname()
	return hostname, nil
}

// Имя хоста для равного файла имени хоста ответа
func addFileFormatFromHostName(hostname string) string {
	arrayline := strings.Split(hostname, ".")
	return arrayline[0] + ".txt"
}

// jfjjfjjjfjjf
// Создание файлы, указав путь из командной строки (cmd)
func CheckIsFileExist(filepath string) bool {
	_, err := os.Stat(filepath)
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// Создание файла по имени Newfile
func createFileFromDirectory(dir string, filename *string) (*os.File, error) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fmt.Println("Каталог не существует!")
	}
	f, err := os.Create(strings.Join([]string{dir, *filename}, "/"))
	if err != nil {
		fmt.Print("Anable to create File", filename)
	}
	fmt.Println(*filename, " Был создан")
	return f, nil
}

// testcommit
func main() {
	srcflag := "src"
	dstflag := "dst"
	src, dst := GetPathFromCommandLine(srcflag, dstflag)
	start := time.Now()
	if src == "None" || dst == "None" || src == "" || dst == "" {
		fmt.Println("->Введите правильную командную строку:(--src=./file.txt  --dst=./)")
	} else {
		f, err := os.Open(src)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		scan, err := io.ReadAll(f)
		if err != nil {
			fmt.Println(err)
		}
		var wg sync.WaitGroup
		urlarray := strings.Split(string(scan), "\n")
		var urlarraycount int = 0
		for urlarraycount < len(urlarray) {
			if !isUrl(urlarray[urlarraycount]) || urlarray[urlarraycount] == "" {
				urlarraycount++
				continue
			}
			wg.Add(1)
			var url = urlarray[urlarraycount]
			go func(url string) {
				defer wg.Done()
				host, err := GethostnameFromURL(url)
				hostname := addFileFormatFromHostName(host)
				resp, err := http.Get(url)
				if err != nil {
					fmt.Print(err)
					return
				}
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if body == nil {
					return
				}
				f, err := createFileFromDirectory(dst, &hostname)
				defer f.Close()
				fi, err := f.WriteString(string(body))
				if err != nil {
					fmt.Print(err, fi)
					return
				}
				urlarraycount++
			}(url)
			wg.Wait()
		}
		end := time.Now()
		elapse := end.Sub(start)
		fmt.Println("время выполнения программы:", elapse)
	}

}
