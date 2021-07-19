package tcsrdrserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"tcsrdr/internal/tcsrdrconfig"
	"time"
)

var config tcsrdrconfig.TCSRDRCfgFile
var defaultConfigPath = "config/config.yaml"

const (
	linkPortfolio  = "portfolio"
	linkOrderbook  = "market/orderbook"
	layoutDateTime = "02.01.2006 15:04 MST"
	layoutDateOnly = "02.01.2006"
)

func Init(isReload bool) error {
	err := tcsrdrconfig.GetConfig(&config, &defaultConfigPath)
	if err != nil {
		return err
	}

	if isReload {
		log.Println("Config has been reloaded")
	} else {
		log.Println("Config init done")
	}

	return nil
}

func getLastPrice(figi *string) (float64, error) {

	var orderbook_resp TCSPriceResponse

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	url := fmt.Sprintf("%s/%s?figi=%s&depth=%d", config.Url, linkOrderbook, *figi, config.Depth)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	req.Header.Add("Authorization", "Bearer "+config.Token)
	req.Header.Add("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	defer resp.Body.Close()

	if err != nil {
		log.Printf("Something wrong during update %s: %v", *figi, err)
		return 0, err
	}

	err = json.NewDecoder(resp.Body).Decode(&orderbook_resp)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return orderbook_resp.Payload.LastPrice, nil
}

func getPortfolioData(response *TCSPortfolioResponse) error {

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	url := fmt.Sprintf("%s/%s", config.Url, linkPortfolio)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	req.Header.Add("Authorization", "Bearer "+config.Token)
	req.Header.Add("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Incorrect response: %v\n", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func httpUpdateHandler(response *TCSPortfolioResponse, currentDateTime *string, currentDate *string, cachedResponse *map[string]float64) error {

	var waitGr sync.WaitGroup

	// Step 1. Get the portfolio's current state.

	waitGr.Add(1)
	go func() error {
		defer waitGr.Done()
		err := getPortfolioData(response)

		if err != nil {
			return err
		}

		return nil
	}()
	waitGr.Wait()

	// Step 2. As TinkoffAPI can't provide prices in portfolio, we have to request them separately.
	//         (You may get some 0.00 in average instead of real prices, if you get the stocks due to corporate actions).

	// Step 2.1. Clear cache

	*cachedResponse = make(map[string]float64)

	// Step 2.2. Request prices

	for _, item := range response.Positions {

		if item.InstrumentType != "Currency" {
			waitGr.Add(1)
			go func(ticker string, figi string) error {
				defer waitGr.Done()

				req, err := getLastPrice(&figi)

				if err != nil {
					return err
				}

				(*cachedResponse)[ticker] = req
				return nil

			}(item.Ticker, item.Figi)
		}
	}

	waitGr.Wait()

	*currentDateTime = time.Now().Format("02.01.2006 15:04 MST")

	timeParser, err := time.Parse(layoutDateTime, *currentDateTime)
	if err != nil {
		return nil
	}
	*currentDate = timeParser.Format(layoutDateOnly)

	return nil
}

func HttpTickerServer(wg *sync.WaitGroup) {

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", config.Ip, config.Port),
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	}

	var response TCSPortfolioResponse
	var currentDateTime string
	var currentDate string
	var cachedResponse map[string]float64 = map[string]float64{}

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {

		err := httpUpdateHandler(&response, &currentDateTime, &currentDate, &cachedResponse)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Something went wrong during info request, please check logs")
			log.Println(err)
		} else {
			fmt.Fprintf(w, "Done at %v\n", currentDateTime)
			log.Println("Info has been updated")
		}
	})

	http.HandleFunc("/getPortfolio", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%+v\n", response)
	})

	http.HandleFunc("/getTicker", func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["name"]

		if !ok {
			log.Fatalf("Something went wrong during param parse: %v\n", keys)
		}

		val, ok := cachedResponse[keys[0]]

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Ticker %s is not found!\n", keys[0])
		} else {
			fmt.Fprintf(w, "date: %s, ticker: %s, price: %f\n", currentDate, keys[0], val)
		}
	})

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		err := Init(true)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Something went wrong, check the logs\n")
		}

		fmt.Fprintf(w, "Config has been reloaded\n")
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Server shutdown initiated")
		s.Shutdown(context.Background())
	})

	http.HandleFunc("/getCache", func(w http.ResponseWriter, r *http.Request) {
		if len(cachedResponse) == 0 {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Cache is empty!\n")
		} else {
			fmt.Fprintf(w, "Updated at %v\n", currentDateTime)
			for key, value := range cachedResponse {
				fmt.Fprintf(w, "%s: %f\n", key, value)
			}
		}
	})

	http.HandleFunc("/getConfig", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%+v\n", config)))
	})

	go func() {
		defer wg.Done()

		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error during websrv execution: %v", err)
		}
	}()
}
