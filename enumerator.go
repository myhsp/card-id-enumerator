package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"

	"github.com/asmcos/requests"
	_ "github.com/mattn/go-sqlite3"
)

var web_path string = ""
var machine_id string = ""

func binarySearch(left, right int, results chan<- WebServiceResult[PersonCenter], progress chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	if left > right {
		return
	}
	mid := left + (right-left)/2
	progress <- mid

	ws, err := get_student_info[PersonCenter](fmt.Sprintf("%X", mid))
	if err != nil {
		if !ws.HasError() {
			ws.Result.CardID = fmt.Sprintf("%X", mid)
			results <- ws

			binarySearch(left, mid-1, results, progress, wg)
			binarySearch(mid+1, right, results, progress, wg)
		} else {
			binarySearch(left, mid-1, results, progress, wg)
			binarySearch(mid+1, right, results, progress, wg)
		}
	}
}

func parse_remote_url(web_path string, api_name string, zipped bool) string {
	if !zipped {
		return strings.Join(
			[]string{
				web_path,
				"/Services/SmartBoard/",
				api_name,
				"/json",
			}, "")
	} else {
		return strings.Join(
			[]string{
				web_path,
				"/Services/SmartBoard/",
				api_name,
				"/json",
				".gzip",
			}, "")
	}
}

func http_post(api_name string, data requests.Datas, zipped bool) (requests.Response, error) {
	if web_path != "" {
		url := parse_remote_url(web_path, api_name, zipped)
		resp, err := requests.Post(url, data)
		return *resp, err
	} else {
		return requests.Response{}, errors.New("empty web path")
	}
}

func get_student_info[T any](cardId string) (WebServiceResult[T], error) {
	resp, err := http_post("SmartBoardPersonCenterNew", requests.Datas{
		"serial": cardId,
		"tCode":  machine_id,
		"userid": "",
	}, false)

	var ws WebServiceResult[T]
	_ = json.Unmarshal(resp.Content(), &ws)
	return ws, err
}

func main() {
	fmt.Print("Input WebPath:")
	_, err := fmt.Scanln(&web_path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Input MachineID:")
	_, err = fmt.Scanln(&machine_id)
	if err != nil {
		log.Fatal(err)
	}

	if !strings.HasPrefix(web_path, "http://") {
		web_path = "http://" + web_path
	}

	var goroutine_count int
	_, _ = fmt.Scanf("Input Goroutine number:%d", &goroutine_count)

	var start int
	_, _ = fmt.Scanf("Input start card no.(hex):%X", &start)

	var end int
	_, _ = fmt.Scanf("Input end card no.(hex):%X", &end)

	db, err := sql.Open("sqlite3", "./results.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS search_results (id TEXT PRIMARY KEY, uname TEXT, uid TEXT, photo TEXT, classroom TEXT, class TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	result := make(chan WebServiceResult[PersonCenter], 1)
	progress := make(chan int, 1)

	var wg sync.WaitGroup

	slice_len := int(math.Ceil(float64(end-start) / float64(goroutine_count)))

	for i := 0; i < goroutine_count; i++ {
		s := start + int(slice_len)*int(i)
		e := s + slice_len

		if e > end {
			e = end
		}

		wg.Add(1)
		go binarySearch(s, e-1, result, progress, &wg)
	}

	go func() {
		for p := range progress {
			if p%128 == 0 {
				fmt.Println("Current progress:", p)
			}
		}
	}()

	go func() {
		for r := range result {
			_, err := db.Exec("INSERT INTO search_results (id, uname, uid, photo, classroom, class) VALUES (?,?,?,?,?,?)",
				r.Result.CardID, r.Result.Username, r.Result.UserID, r.Result.UserPhoto, r.Result.Classroom, r.Result.ClassName)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Result found at:", r.Result.CardID)
		}
	}()

	wg.Wait()
	close(result)
	close(progress)
}
