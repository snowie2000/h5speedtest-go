// speedtest-go project main.go
package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var addr string
var wwwroot string
var mime map[string]string = make(map[string]string)
var isWin = os.IsPathSeparator('\\')
var randomdat []byte = make([]byte, 1024*1024)

func prepareMime() {
	mime[".aac"] = "audio/aac"
	mime[".abw"] = "application/x-abiword"
	mime[".arc"] = "application/octet-stream"
	mime[".avi"] = "video/x-msvideo"
	mime[".azw"] = "application/vnd.amazon.ebook"
	mime[".bin"] = "application/octet-stream"
	mime[".bz"] = "application/x-bzip"
	mime[".bz2"] = "application/x-bzip2"
	mime[".csh"] = "application/x-csh"
	mime[".css"] = "text/css"
	mime[".csv"] = "text/csv"
	mime[".doc"] = "application/msword"
	mime[".epub"] = "application/epub+zip"
	mime[".gif"] = "image/gif"
	mime[".htm"] = "text/html"
	mime[".html"] = "text/html"
	mime[".ico"] = "image/x-icon"
	mime[".ics"] = "text/calendar"
	mime[".jar"] = "application/java-archive"
	mime[".jpeg"] = "image/jpeg"
	mime[".jpg"] = "image/jpeg"
	mime[".js"] = "application/javascript"
	mime[".json"] = "application/json"
	mime[".json"] = "application/javascript"
	mime[".mid"] = "audio/midi"
	mime[".midi"] = "audio/midi"
	mime[".mpeg"] = "video/mpeg"
	mime[".mpkg"] = "application/vnd.apple.installer+xml"
	mime[".odp"] = "application/vnd.oasis.opendocument.presentation"
	mime[".ods"] = "application/vnd.oasis.opendocument.spreadsheet"
	mime[".odt"] = "application/vnd.oasis.opendocument.text"
	mime[".oga"] = "audio/ogg"
	mime[".ogv"] = "video/ogg"
	mime[".ogx"] = "application/ogg"
	mime[".pdf"] = "application/pdf"
	mime[".ppt"] = "application/vnd.ms-powerpoint"
	mime[".rar"] = "application/x-rar-compressed"
	mime[".rtf"] = "application/rtf"
	mime[".sh"] = "application/x-sh"
	mime[".svg"] = "image/svg+xml"
	mime[".swf"] = "application/x-shockwave-flash"
	mime[".tar"] = "application/x-tar"
	mime[".tif"] = "image/tiff"
	mime[".tiff"] = "image/tiff"
	mime[".ttf"] = "application/x-font-ttf"
	mime[".vsd"] = "application/vnd.visio"
	mime[".wav"] = "audio/x-wav"
	mime[".weba"] = "audio/webm"
	mime[".webm"] = "video/webm"
	mime[".webp"] = "image/webp"
	mime[".woff"] = "application/x-font-woff"
	mime[".xhtml"] = "application/xhtml+xml"
	mime[".xls"] = "application/vnd.ms-excel"
	mime[".xlsx"] = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	mime[".xml"] = "application/xml"
	mime[".xul"] = "application/vnd.mozilla.xul+xml"
	mime[".zip"] = "application/zip"
	mime[".3gp"] = "video/3gpp"
	mime[".3g2"] = "video/3gpp2"
	mime[".7z"] = "application/x-7z-compressed"
}

func initDefault() {
	p, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	wwwroot = p + string(os.PathSeparator) + "speedtest"
	io.ReadFull(rand.Reader, randomdat)
	prepareMime()
}

func main() {
	initDefault()
	fmt.Println("Speedtest by golang!")
	flag.StringVar(&addr, "l", ":80", "binding address")
	flag.StringVar(&wwwroot, "r", wwwroot, "web root folder")
	flag.Parse()
	http.HandleFunc("/empty.php", ulHandler)
	http.HandleFunc("/garbage.php", dlHandler)
	http.HandleFunc("/getIP.php", ipHandler)
	http.HandleFunc("/", fileHandler)
	fmt.Println("Listening on", addr)
	fmt.Println("webroot at", wwwroot)
	//comment the lines below if you're running Windows
	if !isWin {
		syscall.Chroot(wwwroot)
		wwwroot = ""
	}

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal(e)
	}
	http.Serve(l, nil)
}

func ulHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Add("Cache-Control", "post-check=0, pre-check=0")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Connection", "keep-alive")
	io.Copy(ioutil.Discard, r.Body)
	w.WriteHeader(200)
}

func dlHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Add("Cache-Control", "post-check=0, pre-check=0")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Content-Description", "File Transfer")
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", "attachment; filename=random.dat")
	w.Header().Add("Content-Transfer-Encoding", "binary")
	w.WriteHeader(200)

	var e error
	szCkSize := r.URL.Query().Get("ckSize")
	ckSize := 4
	if szCkSize != "" {
		if ckSize, e = strconv.Atoi(szCkSize); e != nil {
			ckSize = 4
		}
		if ckSize > 100 {
			ckSize = 100
		}
	}
	for ; ckSize > 0; ckSize-- {
		w.Write(randomdat)
	}
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	sip := r.Header.Get("X-FORWARDED-FOR")
	if sip == "" {
		sip = r.RemoteAddr
	}
	if !strings.Contains(sip, ":") {
		sip += ":0"
	}
	if ip, e := net.ResolveTCPAddr("tcp", sip); e == nil {
		w.Write([]byte("{\"processedString\": \"" + ip.IP.String() + "\", \"rawIspInfo\":\"\"}"))
	} else {
		w.Write([]byte(e.Error()))
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var fileurl string
	fileurl = r.URL.Path
	if fileurl == "" || fileurl == "/" || fileurl == "\\" {
		fileurl = "/index.html"
	}
	fileurl = wwwroot + fileurl
	f, e := os.Open(fileurl)
	if e == nil {
		defer f.Close()

		if m, ok := mime[strings.ToLower(path.Ext(fileurl))]; ok {
			w.Header().Set("Content-Type", m)
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}
		w.WriteHeader(200)
		io.Copy(w, f)
	} else {
		w.WriteHeader(404)
		w.Write([]byte(e.Error()))
	}
}
