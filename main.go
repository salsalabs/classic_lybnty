package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	lybnty "github.com/salsalabs/classic_lybnty/pkg"
	godig "github.com/salsalabs/godig/pkg"
	"gopkg.in/alecthomas/kingpin.v2"
)

//Application to create a Last Year But Not This Year  (LYBNTY) report for Salsa Classic.
//The user gets to define the start and end of both "Last Year" and "This Year", making this
//more of a general-purpose report than a typicaly LYBUNTY report.

func main() {

	t := time.Now()
	y := t.Year()
	v1 := fmt.Sprintf("%4d-01-01", y-1)
	v2 := fmt.Sprintf("%4d-12-31", y-1)
	v3 := fmt.Sprintf("%4d-01-01", y)
	v4 := fmt.Sprintf("%4d-12-31", y)
	p1 := fmt.Sprintf("Last year start date, default is %s", v1)
	p2 := fmt.Sprintf("Last year end date, default is %s", v2)
	p3 := fmt.Sprintf("This year start date, default is %s", v3)
	p4 := fmt.Sprintf("This year end date, default is %s", v4)

	var (
		app        = kingpin.New("classic_lybnty", "A command-line app to create an LYBNTY for a Salsa Classic instance")
		login      = app.Flag("login", "YAML file with API token").String()
		org        = app.Flag("org", "Organization name (for output file)").String()
		lastStart  = app.Flag("last-year-start", p1).Default(v1).String()
		lastEnd    = app.Flag("last-year-end", p2).Default(v2).String()
		thisStart  = app.Flag("this-year-start", p3).Default(v3).String()
		thisEnd    = app.Flag("this-year-end", p4).Default(v4).String()
		apiVerbose = app.Flag("api-verbose", "Makes the app noisy by showing API calls and results").Bool()
	)
	app.Parse(os.Args[1:])

	fmt.Printf("login: %s\n", *login)
	a, err := godig.YAMLAuth(*login)
	if err != nil {
		log.Fatalf("Authentication error: %+v\n", err)
	}
	if apiVerbose != nil {
		a.Verbose = *apiVerbose
	}
	rt, err := lybnty.NewRuntime(a, org, lastStart, lastEnd, thisStart, thisEnd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Runtime: %+v\n", rt)

	c1 := make(chan map[string]string)
	c2 := make(chan lybnty.Data)
	cached := make(chan bool)
	var wg sync.WaitGroup

	log.Printf("main: Done")

	//Start the writer.
	go (func(rt *lybnty.Runtime, c2 chan lybnty.Data, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := rt.Write(c2)
		if err != nil {
			log.Fatal(err)
		}
	})(rt, c2, &wg)
	log.Printf("main: writer started")

	//Start the filter. Note that Filter waits for cached before processing.
	go (func(rt *lybnty.Runtime, cached chan bool, c2 chan lybnty.Data, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		rt.Filter(cached, c2)
	})(rt, cached, c2, &wg)
	log.Printf("main: filter started")

	//Start the cache.  Note that cache sends to cached when done.
	go (func(rt *lybnty.Runtime, c1 chan map[string]string, cached chan bool, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		rt.Cache(c1, cached)
	})(rt, c1, cached, &wg)
	log.Printf("main: cache started")

	//Start the push.
	go (func(rt *lybnty.Runtime, c1 chan map[string]string, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := rt.Push(c1)
		if err != nil {
			log.Fatal(err)
		}
	})(rt, c1, &wg)
	log.Printf("main: push started")

	//Wait for things to terminate.
	d, err := time.ParseDuration("10s")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("main: napping")
	time.Sleep(d)
	log.Printf("main: nap done")
	log.Printf("main: waiting...")
	wg.Wait()
	log.Printf("main: Done")
}
