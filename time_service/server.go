package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Athlete struct {
	ID          uint16
	StartNumber uint16
	FullName    string
}

type RunnerTime struct {
	ID           uint16    `json:"id"`
	CorridorTime time.Time `json:"corridorTime"`
	FinishTime   time.Time `json:"finishTime"`
	TimingPont   string    `json:"timingPoint"`
}

type AthleteTime struct {
	ID          json.Number `json:"id"`
	StartNumber json.Number `json:"startNumber"`
	FullName    string      `json:"fullName"`
	TimingPoint string      `json:"timingPoint"`
	Time        time.Time   `json:"time"`
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan RunnerTime
}

var RunnerTimes []RunnerTime
var Athletes []Athlete
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	generateAthletes()
	pool := newPool()
	go pool.Start()
	handleRequests(pool)
}

func handleRequests(pool *Pool) {
	myRouter := mux.NewRouter().StrictSlash(true)

	// for local testing
	myRouter.HandleFunc("/favicon.ico", doNothing)
	myRouter.HandleFunc("/postRunnerTime",
		func(w http.ResponseWriter, r *http.Request) { postRunner(pool, w, r) }).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/updateRunnerTime",
		func(w http.ResponseWriter, r *http.Request) { updateRunnerTime(pool, w, r) }).Methods("UPDATE", "OPTIONS")
	myRouter.HandleFunc("/getAllRunnerTimes", getAllRunnerTimes)
	myRouter.HandleFunc("/getUnFinishedTimes", getUnFinishedTimes)
	myRouter.HandleFunc("/getFinishedTimes", getFinishedTimes)
	myRouter.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { serveWs(pool, w, r) })
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func getFinishedTimes(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	var athleteTimes []AthleteTime

	if len(RunnerTimes) == 0 {
		return
	}

	for i := range RunnerTimes {
		var runner = RunnerTimes[i]

		if (runner.CorridorTime != time.Time{} && runner.FinishTime != time.Time{}) {
			var athlete Athlete = findAthlete(runner.ID)

			if (athlete != Athlete{}) {
				var athleteTime = AthleteTime{ID: json.Number(strconv.Itoa(int(runner.ID))),
					StartNumber: json.Number(strconv.Itoa(int(athlete.StartNumber))),
					FullName:    athlete.FullName,
					Time:        runner.FinishTime}

				athleteTimes = append(athleteTimes, athleteTime)
			}
		}
	}

	json.NewEncoder(w).Encode(athleteTimes)
}

func getUnFinishedTimes(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	var athleteTimes []AthleteTime

	if len(RunnerTimes) == 0 {
		return
	}

	for i := range RunnerTimes {
		var runner = RunnerTimes[i]

		if (runner.CorridorTime != time.Time{} && runner.FinishTime == time.Time{}) {
			var athlete Athlete = findAthlete(runner.ID)

			if (athlete != Athlete{}) {
				var athleteTime = AthleteTime{ID: json.Number(strconv.Itoa(int(runner.ID))),
					StartNumber: json.Number(strconv.Itoa(int(athlete.StartNumber))),
					FullName:    athlete.FullName,
					Time:        runner.CorridorTime}

				athleteTimes = append(athleteTimes, athleteTime)
			}
		}
	}

	json.NewEncoder(w).Encode(athleteTimes)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Content-Type", "application/json")
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func serveWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.read()
}

func postRunner(pool *Pool, w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	reqBody, _ := ioutil.ReadAll(r.Body)
	if len(reqBody) == 0 {
		return
	}

	var runner RunnerTime
	json.Unmarshal(reqBody, &runner)
	var chipNumber, err = strconv.ParseUint(generateRandomNumber(1, 1000), 16, 16)
	if err != nil {
		fmt.Println(err)
		return
	}

	runner.ID = uint16(chipNumber)
	runner.TimingPont = "Finish Line Corridor"
	// should be client's but not sure how to convert it
	RunnerTimes = append(RunnerTimes, runner)

	json.NewEncoder(w).Encode(runner)
	broadcastRunner(pool, runner)
}

func broadcastRunner(pool *Pool, runner RunnerTime) {
	pool.Broadcast <- runner
}

func generateRandomNumber(min int, max int) string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func getAllRunnerTimes(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var athleteTimes []AthleteTime

	if len(RunnerTimes) == 0 {
		return
	}

	for i := range RunnerTimes {
		var runner = RunnerTimes[i]
		var athlete Athlete = findAthlete(runner.ID)

		if (athlete != Athlete{}) {
			var runnerTimePoint = runner.CorridorTime;
			if (runner.FinishTime != time.Time{}) {
				runnerTimePoint = runner.FinishTime;
			}

			var athleteTime = AthleteTime{ID: json.Number(strconv.Itoa(int(runner.ID))),
				StartNumber: json.Number(strconv.Itoa(int(athlete.StartNumber))),
				FullName:    athlete.FullName,
				TimingPoint: runner.TimingPont,
				Time:        runnerTimePoint}

			athleteTimes = append(athleteTimes, athleteTime)
		}
	}

	json.NewEncoder(w).Encode(athleteTimes)
}

func findAthlete(ID uint16) Athlete {
	athlete := Athlete{}

	for j := range Athletes {
		athlete = Athletes[j]

		if ID == athlete.ID {
			return athlete
		}
	}
	return athlete
}

func updateRunnerTime(pool *Pool, w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	reqBody, _ := ioutil.ReadAll(r.Body)
	if len(reqBody) == 0 {
		return
	}

	var runner RunnerTime
	json.Unmarshal(reqBody, &runner)
	runner.TimingPont = "Finish Line"

	for i := range RunnerTimes {
		if RunnerTimes[i].ID == runner.ID {
			var runnerHelper RunnerTime = RunnerTimes[i]
			runnerHelper.FinishTime = runner.FinishTime
			runnerHelper.TimingPont = runner.TimingPont;
			RunnerTimes[i] = runnerHelper
			break
		}
	}
	pool.Broadcast <- RunnerTime{}
}

func doNothing(w http.ResponseWriter, r *http.Request) {}

func (c *Client) read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()
}

func newPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan RunnerTime),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			break
		case runner := <-pool.Broadcast:
			for client := range pool.Clients {
				// essentially the idea was to send a runner to add to the table, but the render didnt happen after state changed
				// so i just send a meaningless object to trigger a new getallrunners
				if err := client.Conn.WriteJSON(runner); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

func generateAthletes() {
	for i := 1; i <= 1000; i++ {
		var randomNumber, err = strconv.ParseUint(generateRandomNumber(1, 9999), 16, 16)
		if err != nil {
			fmt.Println(err)
			return
		}
		var athlete = Athlete{ID: uint16(i), StartNumber: uint16(randomNumber), FullName: fmt.Sprint("name test", i)}
		Athletes = append(Athletes, athlete)
	}
}
