// speedtest-go project main.go
package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
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
	"time"

	"github.com/oschwald/geoip2-golang"
)

var (
	addr      string
	wwwroot   string
	mime      map[string]string = make(map[string]string)
	isWin                       = os.IsPathSeparator('\\')
	randomdat []byte            = make([]byte, 1024*1024)
	appdir    string
	ipserver  *geoIP = nil
)

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
	appdir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	wwwroot = appdir + string(os.PathSeparator) + "speedtest"
	io.ReadFull(rand.Reader, randomdat)
	prepareMime()
}

type geoIP struct {
	db *geoip2.Reader
}

func (this *geoIP) GetCity(ip string) string {
	tcpaddr, e := net.ResolveTCPAddr("tcp", ip)
	if e == nil {
		city, e := this.db.City(tcpaddr.IP)
		//isp, e2 := this.db.ISP(tcpaddr.IP)
		if e == nil {
			result := fmt.Sprintln(city.City.Names["en"], city.Country.Names["en"])
			result = strings.ReplaceAll(result, "\x0a", "")
			if len(result) >= 2 {
				return result
			}
		}
	}
	return "Unknown ISP"
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
	http.HandleFunc("/backend/empty.php", ulHandler)
	http.HandleFunc("/backend/garbage.php", dlHandler)
	http.HandleFunc("/backend/getIP.php", ipHandler)

	http.HandleFunc("/", fileHandler)
	fmt.Println("Listening on", addr)
	fmt.Println("webroot at", wwwroot)

	db, err := geoip2.Open(path.Join(appdir, "ip.dat"))
	if err == nil {
		ipserver = &geoIP{db}
	}
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
		if ckSize > 5000 {
			ckSize = 5000
		}
	}
	for ; ckSize > 0; ckSize-- {
		w.Write(randomdat)
	}
}

func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

func getIP(r *http.Request) string {
	remoteip := r.Header.Get("X-FORWARDED-FOR")
	iplist := strings.Split(remoteip, ",")
	for len(iplist) > 0 && strings.TrimSpace(iplist[len(iplist)-1]) == "127.0.0.1" {
		iplist = iplist[0 : len(iplist)-1] // remove all the localhosts forwarder
	}

	if len(iplist) > 0 {
		if ip, err := net.ResolveIPAddr("ip", strings.TrimSpace(iplist[len(iplist)-1])); err == nil {
			return ip.IP.String() // use ip from x-forwarded-for
		}
	}

	// fallback to r.RemoteAddr
	if ip, err := net.ResolveIPAddr("ip", r.RemoteAddr); err == nil {
		return ip.IP.String()
	} else {
		return "127.0.0.1" // fallback to localhost if r.RemoteAddr is mocked
	}
}

type Ipinfo struct {
	Ip      string
	City    string
	Region  string
	Country string
	Loc     string
	Org     string
}

var httpsClient = &http.Client{
	Timeout: 15 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func getIpInfoOnline(ip string) string {
	resp, err := httpsClient.Get("https://ipinfo.io/" + ip + "/json")
	if err == nil {
		defer resp.Body.Close()
		var ipnfo Ipinfo
		jd := json.NewDecoder(resp.Body)
		if err = jd.Decode(&ipnfo); err == nil {
			return ipnfo.City + " - " + ipnfo.Org
		}
	}
	return "Error " + err.Error()
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	sip := getIP(r)

	locinfo := ""
	if ipserver != nil {
		locinfo = ipserver.GetCity(sip)
	} else {
		locinfo = getIpInfoOnline(sip)
	}
	w.Write([]byte("{\"processedString\": \"" + sip + " - " + locinfo + "\", \"rawIspInfo\":\"\"}"))
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
