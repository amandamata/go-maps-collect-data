package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Country struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type OverpassResponse struct {
	Elements []struct {
		Tags struct {
			Name string `json:"name"`
			ISO  string `json:"ISO3166-2"`
		} `json:"tags"`
	} `json:"elements"`
}

func main() {
	countries := []Country{
		{ID: "BR", Name: "Brasil"},
		{ID: "US", Name: "Estados Unidos"},
		{ID: "AE", Name: "Emirados Árabes Unidos"},
		{ID: "AL", Name: "Albânia"},
		{ID: "SE", Name: "Suécia"},
		{ID: "ZA", Name: "Africa do Sul"},
		{ID: "VE", Name: "Venezuela"},
		{ID: "BA", Name: "Bósnia e Herzegóvina"},
		{ID: "SG", Name: "Singapura"},
		{ID: "KR", Name: "República da Coréia"},
		{ID: "BD", Name: "Bangladesh"},
		{ID: "IR", Name: "Irã"},
	}

	outputFile, err := os.Create("states_and_cities.csv")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"country_code", "state_code", "state_name", "city_name"}
	if err := writer.Write(header); err != nil {
		fmt.Println("Error writing header:", err)
		return
	}

	for _, country := range countries {
		fmt.Printf("Processing country: %s (%s)\n", country.Name, country.ID)
		query := fmt.Sprintf(`[out:json];area["ISO3166-1"="%s"][admin_level=2];(relation["admin_level"="4"](area););out body;`, country.ID)
		encodedQuery := url.QueryEscape(query)
		apiURL := fmt.Sprintf("http://overpass-api.de/api/interpreter?data=%s", encodedQuery)
		resp, err := http.Get(apiURL)
		if err != nil {
			fmt.Println("Error fetching data for country:", country.Name, err)
			continue
		}
		defer resp.Body.Close()

		var overpassResponse OverpassResponse
		if err := json.NewDecoder(resp.Body).Decode(&overpassResponse); err != nil {
			fmt.Println("Error decoding response for country:", country.Name, err)
			continue
		}

		for _, element := range overpassResponse.Elements {
			stateCode := element.Tags.ISO
			stateName := element.Tags.Name
			if stateCode != "" && stateName != "" {
				fmt.Printf("  Found state: %s (%s)\n", stateName, stateCode)
				cityQuery := fmt.Sprintf(`[out:json];area["ISO3166-2"="%s"][admin_level=4];(node["place"~"city|town|village|hamlet|suburb|neighbourhood"](area););out body;`, stateCode)
				encodedCityQuery := url.QueryEscape(cityQuery)
				cityAPIURL := fmt.Sprintf("http://overpass-api.de/api/interpreter?data=%s", encodedCityQuery)
				cityResp, err := http.Get(cityAPIURL)
				if err != nil {
					fmt.Println("Error fetching cities for state:", stateName, err)
					continue
				}
				defer cityResp.Body.Close()

				var cityResponse OverpassResponse
				if err := json.NewDecoder(cityResp.Body).Decode(&cityResponse); err != nil {
					fmt.Println("Error decoding city response for state:", stateName, err)
					continue
				}

				for _, cityElement := range cityResponse.Elements {
					cityName := cityElement.Tags.Name
					if cityName != "" {
						fmt.Printf("    Found city: %s\n", cityName)
						newRecord := []string{
							country.ID,
							stateCode,
							stateName,
							cityName,
						}
						if err := writer.Write(newRecord); err != nil {
							fmt.Println("Error writing record:", err)
							return
						}
					}
				}
			}
		}
	}

	fmt.Println("CSV file processed successfully!")
}
