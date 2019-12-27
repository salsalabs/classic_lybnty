package main

import (
	"fmt"
	"log"
	"os"
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
		app       = kingpin.New("LYBNTY Report for Classic", "A command-line app to create an LYBNTY for a Salsa Classic instance")
		login     = app.Flag("login", "YAML file with API token").String()
		org       = app.Flag("org", "Organization name (for output file)").String()
		lastStart = app.Flag("last-year-start", p1).Default(v1).String()
		lastEnd   = app.Flag("last-year-end", p2).Default(v2).String()
		thisStart = app.Flag("this-year-start", p3).Default(v3).String()
		thisEnd   = app.Flag("this-year-end", p4).Default(v4).String()
	)
	app.Parse(os.Args[1:])

	fmt.Printf("login: %s\n", *login)
	a, err := godig.YAMLAuth(*login)
	if err != nil {
		log.Fatalf("Authentication error: %+v\n", err)
	}
	rt, err := lybnty.NewRuntime(a, org, lastStart, lastEnd, thisStart, thisEnd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Runtime: %+v\n", rt)

	c1 := make(chan map[string]string)
	c2 := make(chan lybnty.Data)
	var wg sync.WaitGroup

	log.Printf("main: Done")

	//Start the writer.
	go (function(rt RunTime, c2 chan lybnty.Data, wg &sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := rt.Write(c2)
		if err != nil {
			log.Fatal(err)
		}
	})(rt, c2, &wg)
	log.Printf("main: writer started")

	//Start the filter.
	go (function(rt RunTime, c1 chan map[string]string, c2 chan lybnty.Data, wg &sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := rt.Filter(c1, c2)
		if err != nil {
			log.Fatal(err)
		}
	})(rt, c1, c2, &wg)
	log.Printf("main: filter started")
	
	//Start the push.
	go (function(rt RunTime, c1 chan map[string]string, wg &sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err := rt.Push(c1)
		if err != nil {
			log.Fatal(err)
		}
	})(rt, c1, &wg)
	log.Printf("main: push started")

	//Wait for things to terminate.
	wg.Wait()
	log.Printf("main: Done")
}
