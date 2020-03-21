package main
// Command line application to POST location information to the WhereAmI Server

import (
	"fmt"
	"log"
	"time"
	"github.com/go-resty/resty/v2"
)

type LocationResponse struct {
    Status string `json:"status"`
    Country string `json:"country"`
    CountryCode string `json:"countryCode"`
		Region string `json:"region"`
		RegionName string `json:"regionName"`
		City string `json:"city"`
		Zip string `json:"zip"`
		Lat float32 `json:"lat"`
		Lon float32 `json:"lon"`
		Timezone string `json:"timezone"`
		Isp string `json:"isp"`
		Org string `json:"org"`
		As string `json:"as"`
		Query string `json:"query"`
}

type LocationResult struct {
	CountryLong string `json:"country_long"`
	Country string `json:"country"`
	City string `json:"city"`
	Region string `json:"region"`
	Lat float32 `json:"latitude"`
	Lon float32 `json:"longitude"`
	Timestamp int32 `json:"timestamp"`
	Ip string `json:"ip"`
}

// Uploads the user's approximate location to whereami.mattgowie.com
//
// 1. Fetches the client's IP Address from icanhazip.com
// 2. Does a lookup of the fetched IP Address from the ip-api.com API
// 3. POSTs the client's location infomration to whereami.mattgowie.com
func main() {
	client := resty.New()

	ip_resp, err := client.R().Get("https://icanhazip.com")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("IP Address: ", ip_resp)


	location_url := fmt.Sprintf("http://ip-api.com/json/%s", ip_resp)
	location_resp, location_err := client.R().SetResult(LocationResponse{}).Get(location_url)

	if location_err != nil {
		log.Fatal(location_err)
	}

	fmt.Println("Location Response: ", location_resp)

	var location = location_resp.Result().(*LocationResponse)

	var result = LocationResult {
		CountryLong: location.Country,
		Country    : location.CountryCode,
		City       : location.City,
		Region     : location.Region,
		Lat        : location.Lat,
		Lon        : location.Lon,
		Timestamp  : int32(time.Now().Unix()),
		Ip         : fmt.Sprintf("%s", ip_resp) }

	fmt.Println("Result: ", result)

	post_resp, post_err := client.R().
      SetBody(result).
      Post("https://docker-micro.mattgowie.com/location")

	if post_err != nil {
		log.Fatal(post_err)
	}

	fmt.Println("Location Update Response: ", post_resp)
	fmt.Println("Completed successfully!")
}
