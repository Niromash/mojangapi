package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	// Send 300 req to the url https://verceltest-sooty.vercel.app/api/profile?uuid=65a854778e6e42728f4f766c3beb1734 in go routines

	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// req to https://verceltest-sooty.vercel.app/api/profile?uuid=65a854778e6e42728f4f766c3beb1734
			resp, err := http.Get("https://verceltest-sooty.vercel.app/api/v1/profile?uuid=65a854778e6e42728f4f766c3beb1734")
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			fmt.Println(resp.StatusCode)
		}()
	}
	wg.Wait()
}
