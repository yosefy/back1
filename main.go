package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strconv"
	"time"
	"html/template"
	"log"
	"sort"
	"os/exec"
	"strings"
)

var Netmaski = "192.168."
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

var freeIpsCollector []int
var takenIpsCollector []int

func pingOne(a int, count int) () {
	string_to_ping := Netmaski + strconv.Itoa(a) + "." + strconv.Itoa(count)
	out, _ := exec.Command("ping", string_to_ping, "-c 2", "-W 3").Output()
	if strings.Contains(string(out), "2 packets transmitted, 2 received") {
		takenIpsCollector = append(takenIpsCollector, count)
	} else {
		freeIpsCollector = append(freeIpsCollector, count)
	}

}

func checkAllIps(a string) {
	freeIpsCollector = nil
	takenIpsCollector = nil
	subsub, _ := strconv.Atoi(a)
	for i := 1; i < 255; i++ {
		go pingOne(subsub, i)
	}
	time.Sleep(time.Second * 6)
}

type Lists2 struct {
	Taken []int
	Free  []int
}

type Data struct {
	Maski string
	Sub   string
	Lists2
}

func foo(w http.ResponseWriter, req *http.Request) {
	var data Data
	data.Maski = Netmaski
	s := req.FormValue("subnet")
	data.Sub = s
	if req.Method == http.MethodPost {

		checkAllIps(s)
		fmt.Print(takenIpsCollector)
		sort.Ints(takenIpsCollector)
		sort.Ints(freeIpsCollector)
		data.Taken = takenIpsCollector
		data.Free = freeIpsCollector

		//              time.Sleep(time.Second * 20)
	}

	err := tpl.ExecuteTemplate(w, "index.gohtml", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}

}

type person struct {
	Fname string
	Lname string
	Items []string
}

func encd(w http.ResponseWriter, req *http.Request) {

	var data Data
	data.Maski = Netmaski
	s := "10"
	data.Sub = s
	if req.Method == http.MethodPost {

		checkAllIps(s)
		fmt.Print(takenIpsCollector)
		sort.Ints(takenIpsCollector)
		sort.Ints(freeIpsCollector)
		data.Taken = takenIpsCollector
		data.Free = freeIpsCollector
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data.Free)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/free_enc", encd)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":7070", nil)
}
