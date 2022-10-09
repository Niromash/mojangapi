package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Send 300 req to the url https://verceltest-sooty.vercel.app/api/profile?uuid=65a854778e6e42728f4f766c3beb1734 in go routines
	start := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// req to https://verceltest-sooty.vercel.app/api/profile?uuid=65a854778e6e42728f4f766c3beb1734
			resp, err := http.Get("https://verceltest-sooty.vercel.app/api/v1/from_username/profile?username=Niromash_")
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				log.Fatalln(resp.StatusCode)
			}
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(start))
}
