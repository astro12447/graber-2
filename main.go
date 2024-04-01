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
	"time"
)

// Get path location from command line(cmd).

func GetPathFromCommandLine(src string, dst string) (string, string) {
	var sources *string
	var destination *string
	sources = flag.String(src, "None", "")
	destination = flag.String(dst, "None", "")
	flag.Parse()

	return *sources, *destination
}

// Get format  URS->(https://HostDomainName) links from  file data
func ReadDataFromFile(fileName string) string {
	//read file lines one by one line in memory
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("anable to Open File: %v", err)
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
		}
	}
	return FileContent
}

// Domain format(https://www.google.com/)
type Domain struct {
	schema   string
	host     string
	hostname string
	path     string
}

// Get Real url link from Experimetals Data Files
func IsUrl(input string) bool {
	flag := false
	a, err := url.Parse(input)
	if err != nil {
		fmt.Print(err)
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

// Get a valid Url from Random Data
func GetUrlFromFile(FileName string) *list.List {
	data := ReadDataFromFile(FileName)
	lines := strings.Split(data, "\n")
	l := list.New()
	for _, item := range lines {
		if IsUrl(item) == true {
			l.PushBack(item)
		}
	}
	return l
}

// Create Random Files Names, Get names url Files
func GethostnameFromURL(URL string) string {
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatal("URL given not correcly!", URL)
	}
	hostname := u.Hostname()
	return hostname
}

// Hostname to by equal Reponse Hostname file
func ADDFileFormatFromHostName(hostname string) string {
	lines := strings.Split(hostname, ".")
	return lines[0] + ".txt"
}

// Create Files from giving path from command-lines(cmd)
func CheckIsFileExist(filepath string) bool {
	_, err := os.Stat(filepath)
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
func CreateFromCurrrentDirectory(Dir string) {

}
func IsEmpty(l *list.List) bool {
	return l.Len() == 0
}

// Get Request from Client via FileData(urls)
func RequestReponse(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		host := GethostnameFromURL(e.Value.(string))
		hostname := ADDFileFormatFromHostName(host)
		resp, err := http.Get(e.Value.(string))
		if err != nil {
			fmt.Print(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		f := CreateFile(&hostname)
		defer f.Close()
		f.WriteString(string(body))
		if err != nil {
			fmt.Print(err)
		}
		defer f.Close()
		fmt.Print(f.Name(), " Was been writted succesfully.. ")
		fmt.Print("\n")
	}
}

// Creata a File given By Newfile Name
func CreateFile(Newfile *string) *os.File {
	file, err := os.Create(*Newfile)
	if err != nil {
		log.Fatalf("%v", file)
	}
	fmt.Print(file.Name(), "Was been created succesfully...")
	fmt.Println()
	return file
}

// testcommit

func main() {
	srcflag := "src"
	dstflag := "dst"
	s1, s2 := GetPathFromCommandLine(srcflag, dstflag)
	if s1 == "None" || s2 == "None" || s1 == "" || s2 == "" {
		fmt.Println("->Introduce correct Command line:(--src=./file.txt  --dst=./)")
	} else {
		l := GetUrlFromFile(s1)
		start := time.Now()
		RequestReponse(l)
		end := time.Now()
		elapse := end.Sub(start)
		fmt.Println("Duration time elapse:", elapse)
	}
}
