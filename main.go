package main

import (
	"container/list"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// Получение  пути из командной строки (cmd).
func GetPathFromCommandLine(src string, dst string) (string, string, error) {
	if src == "None" || dst == "None" || src == "" || dst == "" {
		fmt.Println("->Введите правильную командную строку:(--src=./file.txt  --dst=./)")
	}
	var sources *string
	var destination *string
	sources = flag.String(src, "None", "")
	destination = flag.String(dst, "None", "")
	flag.Parse()
	return *sources, *destination, nil
}

// Получение формат url ->(https://HostDomainName) ссылки из данных файла
func readDataFromFile(fileName string) (string, error) {
	//читать строки файла, одну за другой в памяти
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Невозможно открыть файл: %v", err)
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
		fmt.Println("URL дано не правильно!", URL)
	}
	hostname := u.Hostname()
	return hostname, nil
}

// Имя хоста для равного файла имени хоста ответа
func addFileFormatFromHostName(hostname string) string {
	arrayline := strings.Split(hostname, ".")
	return arrayline[0] + ".txt"
}

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
		fmt.Print("Невозможно создать файл", *filename)
	}
	return f, nil
}

func requestFromServer(s string) (*http.Response, error) {
	resp, err := http.Get(s)
	if err != nil {
		panic(err)
	}
	return resp, nil
}
func readRepose(r io.Reader) ([]byte, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		fmt.Println(err)
	}
	return body, nil
}
func writeString(f *os.File, s string) (int, error) {
	fi, err := f.WriteString(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f.Name(), "Был создан и записан успешно!")
	return fi, nil
}
func closefile(file *os.File) {
	defer file.Close()
}
func respClose(close io.Closer) {
	defer close.Close()
}
func readAll(r io.Reader) ([]byte, error) {
	scan, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return scan, nil
}

func main() {
	srcflag := "src"
	dstflag := "dst"
	src, dst, err := GetPathFromCommandLine(srcflag, dstflag)
	fmt.Println(src, dst)
	if err != nil {
		panic(err)
	}
	start := time.Now()
	f, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
	}
	scan, err := readAll(f)
	if err != nil {
		fmt.Println(err)
	}
	var wg sync.WaitGroup
	urlarray := strings.Split(string(scan), "\n")
	for domain := range urlarray {
		if !isUrl(urlarray[domain]) || urlarray[domain] == "" {
			domain++
			continue
		}
		wg.Add(1)
		var url = urlarray[domain]
		go func(url string) {
			defer wg.Done()
			host, err := GethostnameFromURL(url)
			hostname := addFileFormatFromHostName(host)
			resp, err := requestFromServer(url)
			if err != nil {
				fmt.Print(err)
			}
			body, err := readRepose(resp.Body)
			respClose(resp.Body)
			if err != nil {
				panic(err)
			}
			f, err := createFileFromDirectory(dst, &hostname)
			if err != nil {
				panic(err)
			}
			fi, err := writeString(f, string(body))
			closefile(f)
			if err != nil {
				fmt.Print(err, fi)
			}
			domain++
		}(url)
		wg.Wait()
	}
	end := time.Now()
	elapse := end.Sub(start)
	fmt.Println("время выполнения программы:", elapse)
}
