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

	var orderbookResp TCSPriceResponse

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

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			if err != nil {
				log.Println(err)
				err = closeErr
			}
		}
	}()

	if err != nil {
		log.Printf("Something wrong during update %s: %v", *figi, err)
		return 0, err
	}

	err = json.NewDecoder(resp.Body).Decode(&orderbookResp)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return orderbookResp.Payload.LastPrice, nil
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

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			if err != nil {
				log.Println(err)
				err = closeErr
			}
		}
	}()

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

	// Step 1. Get the portfolio's current state.

	err := getPortfolioData(response)

	if err != nil {
		return err
	}

	// Step 2. As TinkoffAPI can't provide prices in portfolio, we have to request them separately.
	//         (You may get some 0.00 in average instead of real prices, if you get the stocks due to corporate actions).
	// Step 2.1. Clear cache

	*cachedResponse = make(map[string]float64)

	// Step 2.2. Request prices

	for _, item := range response.Positions {
		err = processStockItem(&item, cachedResponse)
		if err != nil {
			return err
		}
	}

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
	cachedResponse := map[string]float64{}

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {

		err := httpUpdateHandler(&response, &currentDateTime, &currentDate, &cachedResponse)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, printErr := fmt.Fprintf(w, "Something went wrong during info request, please check logs")
			checkPrintErr(&printErr, &err)
		} else {
			_, printErr := fmt.Fprintf(w, "Done at %v\n", currentDateTime)
			if checkPrintErr(&printErr) {
				log.Println("Info has been updated")
			}
		}
	})

	http.HandleFunc("/getPortfolio", func(w http.ResponseWriter, r *http.Request) {
		_, printErr := fmt.Fprintf(w, "%+v\n", response)
		checkPrintErr(&printErr)
	})

	http.HandleFunc("/getTicker", func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["name"]

		if !ok {
			log.Fatalf("Something went wrong during param parse: %v\n", keys)
		}

		val, ok := cachedResponse[keys[0]]

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_, printErr := fmt.Fprintf(w, "Ticker %s is not found!\n", keys[0])
			_ = checkPrintErr(&printErr)
		} else {
			_, printErr := fmt.Fprintf(w, "date: %s, ticker: %s, price: %f\n", currentDate, keys[0], val)
			_ = checkPrintErr(&printErr)
		}
	})

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		err := Init(true)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, printErr := fmt.Fprintf(w, "Something went wrong, check the logs\n")
			checkPrintErr(&printErr, &err)
		}

		_, printErr := fmt.Fprintf(w, "Config has been reloaded\n")
		checkPrintErr(&printErr)
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Server shutdown initiated")
		err := s.Shutdown(context.Background())
		if err != nil {
			log.Fatalln("Error during server shutdown: ", err)
		}

	})

	http.HandleFunc("/getCache", func(w http.ResponseWriter, r *http.Request) {
		if len(cachedResponse) == 0 {
			w.WriteHeader(http.StatusNotFound)
			_, printErr := fmt.Fprintf(w, "Cache is empty!\n")
			checkPrintErr(&printErr)
		} else {
			_, printErr := fmt.Fprintf(w, "Updated at %v\n", currentDateTime)
			checkPrintErr(&printErr)
			for key, value := range cachedResponse {
				_, printErr := fmt.Fprintf(w, "%s: %f\n", key, value)
				checkPrintErr(&printErr)
			}
		}
	})

	http.HandleFunc("/getConfig", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(fmt.Sprintf("%+v\n", config))); err != nil {
			log.Println("Error during config read/parse: ", err)
		}

	})

	go func() {
		defer wg.Done()

		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error during websrv execution: %v", err)
		}
	}()
}

func checkPrintErr(printErr *error, otherErrors ...*error) bool {
	if *printErr != nil {
		log.Printf("Next error has been occured during response writing: %v", printErr)
		if len(otherErrors) != 0 {
			for _, errMsg := range otherErrors {
				log.Println(*errMsg)
			}
		}
		return false
	}
	return true
}

func processStockItem(item *Position, cachedResponse *map[string]float64) error {
	if item.InstrumentType == "Currency" {
		return nil
	}

	req, err := getLastPrice(&item.Figi)
	if err != nil {
		return err
	} else {
		(*cachedResponse)[item.Ticker] = req
	}

	return nil
}
