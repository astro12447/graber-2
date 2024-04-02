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
func readDataFromFile(fileName string) (string, error) {
	//read file lines, one by one line in memory
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
			fmt.Println("The sources Files Is Empty!")
		}
	}
	return FileContent, nil
}

// Domain format(https://www.google.com/)
type Domain struct {
	schema   string
	host     string
	hostname string
	path     string
}

// Get Real url link from Experimetals Data Files

func isUrl(input string) bool {
	flag := false
	a, err := url.Parse(input)
	if err != nil {
		fmt.Println(input, "Is Empty!..")
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

// Create Random Files Names, Get names url Files
func GethostnameFromURL(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatal("URL given not correcly!", URL)
	}
	hostname := u.Hostname()
	return hostname, nil
}

// Hostname to by equal Reponse Hostname file
func addFileFormatFromHostName(hostname string) string {
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

// Get Request from Client via FileData(urls)
func RequestReponse(l *list.List, dir string) error {
	for e := l.Front(); e != nil; e = e.Next() {
		host, err := GethostnameFromURL(e.Value.(string))
		hostname := addFileFormatFromHostName(host)
		resp, err := http.Get(e.Value.(string))
		if err != nil {
			fmt.Print(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		f, err := createFileFromDirectory(dir, &hostname)
		f.WriteString(string(body))
		if err != nil {
			fmt.Print(err)
		}
		fmt.Print(hostname, " Was been writted succesfully..")
		fmt.Print("\n")
	}
	return nil
}

// Creata a File given By Newfile Name
func createFileFromDirectory(dir string, filename *string) (*os.File, error) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fmt.Println("Directory does not exist!")
	}
	f, err := os.Create(strings.Join([]string{dir, *filename}, "/"))
	if err != nil {
		fmt.Print("Anable to create File", filename)
	}
	fmt.Println(*filename, " Was been Created")
	return f, nil
}

// testcommit

func main() {
	srcflag := "src"
	dstflag := "dst"
	s1, s2 := GetPathFromCommandLine(srcflag, dstflag)
	if s1 == "None" || s2 == "None" || s1 == "" || s2 == "" {
		fmt.Println("->Introduce correct Command line:(--src=./file.txt  --dst=./)")
	} else {
		l, err := getUrlFromfile(s1)
		if err != nil {
			panic(err)
		}
		start := time.Now()
		RequestReponse(l, s2)
		end := time.Now()
		elapse := end.Sub(start)
		fmt.Println("Duration time elapse:", elapse)
	}
}
