package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type TranslateResponse struct {
	TranslatedText string `json:"translatedText"`
}

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
	translateCities()
}

func collectStates() {
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

	outputFile, err := os.Create("states.csv")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"country_code", "state_code", "state_name"}
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

				newRecord := []string{
					country.ID,
					stateCode,
					stateName,
				}
				if err := writer.Write(newRecord); err != nil {
					fmt.Println("Error writing record:", err)
					return
				}
			}
		}
	}

	fmt.Println("CSV file processed successfully!")
}

func collectCities() {
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
				cityQuery := fmt.Sprintf(`[out:json];area["ISO3166-2"="%s"][admin_level=4];(node["place"~"city|town|village"](area););out body;`, stateCode)
				encodedCityQuery := url.QueryEscape(cityQuery)
				cityURL := fmt.Sprintf("http://overpass-api.de/api/interpreter?data=%s", encodedCityQuery)
				cityResp, err := http.Get(cityURL)
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

func translateStates() {
	inputFile, err := os.Open("states.csv")
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	reader := csv.NewReader(inputFile)

	header, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading header:", err)
		return
	}

	newHeader := append(header, "state_name_br", "state_name_en")

	languageMap := map[string]string{
		"BR": "pt", // Portuguese
		"US": "en", // English
		"AE": "ar", // Arabic
		"AL": "sq", // Albanian
		"SE": "sv", // Swedish
		"ZA": "en", // English
		"VE": "es", // Spanish
		"BA": "bs", // Bosnian
		"SG": "en", // English
		"KR": "ko", // Korean
		"BD": "bn", // Bengali
		"IR": "fa", // Persian
	}

	outputFile, err := os.Create("formula_states.csv")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write(newHeader); err != nil {
		fmt.Println("Error writing header:", err)
		return
	}

	rowNumber := 2

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		countryCode := record[0]
		stateName := record[2]

		sourceLang, exists := languageMap[countryCode]
		if !exists {
			fmt.Println("Language not found for country code:", countryCode)
			continue
		}

		var stateNameBRFormula, stateNameENFormula string

		if sourceLang == "en" {
			stateNameBRFormula = fmt.Sprintf(`=GOOGLETRANSLATE(C%d;"%s";"PT")`, rowNumber, sourceLang)
			stateNameENFormula = stateName
		}

		if sourceLang == "pt" {
			stateNameBRFormula = stateName
			stateNameENFormula = fmt.Sprintf(`=GOOGLETRANSLATE(C%d;"%s";"EN")`, rowNumber, sourceLang)
		} else {
			stateNameBRFormula = fmt.Sprintf(`=GOOGLETRANSLATE(C%d;"%s";"PT")`, rowNumber, sourceLang)
			stateNameENFormula = fmt.Sprintf(`=GOOGLETRANSLATE(C%d;"%s";"EN")`, rowNumber, sourceLang)
		}

		newRecord := append(record, stateNameBRFormula, stateNameENFormula)

		if err := writer.Write(newRecord); err != nil {
			fmt.Println("Error writing record:", err)
			return
		}

		rowNumber++
	}

	fmt.Println("CSV file processed successfully!")
}

func translateCities() {
	stateTranslations := make(map[string][2]string)
	stateFile, err := os.Open("translated_states.csv")
	if err != nil {
		fmt.Println("Error opening translated states file:", err)
		return
	}
	defer stateFile.Close()

	stateReader := csv.NewReader(stateFile)
	stateHeader, err := stateReader.Read()
	if err != nil {
		fmt.Println("Error reading translated states header:", err)
		return
	}

	stateNameIndex := -1
	stateNameBRIndex := -1
	stateNameENIndex := -1
	for i, h := range stateHeader {
		switch h {
		case "state_name":
			stateNameIndex = i
		case "state_name_br":
			stateNameBRIndex = i
		case "state_name_en":
			stateNameENIndex = i
		}
	}

	if stateNameIndex == -1 || stateNameBRIndex == -1 || stateNameENIndex == -1 {
		fmt.Println("Error: Missing required columns in translated states file")
		return
	}

	for {
		record, err := stateReader.Read()
		if err != nil {
			break
		}
		stateName := record[stateNameIndex]
		stateNameBR := record[stateNameBRIndex]
		stateNameEN := record[stateNameENIndex]
		stateTranslations[stateName] = [2]string{stateNameBR, stateNameEN}
	}

	inputFile, err := os.Open("country_states_cities.csv")
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	reader := csv.NewReader(inputFile)

	header, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading header:", err)
		return
	}

	newHeader := append(header, "state_name_br", "state_name_en", "city_name_br", "city_name_en")

	languageMap := map[string]string{
		"BR": "pt", // Portuguese
		"US": "en", // English
		"AE": "ar", // Arabic
		"AL": "sq", // Albanian
		"SE": "sv", // Swedish
		"ZA": "en", // English
		"VE": "es", // Spanish
		"BA": "bs", // Bosnian
		"SG": "en", // English
		"KR": "ko", // Korean
		"BD": "bn", // Bengali
		"IR": "fa", // Persian
	}

	outputFile, err := os.Create("formula_country_states_cities.csv")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write(newHeader); err != nil {
		fmt.Println("Error writing header:", err)
		return
	}

	rowNumber := 2

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		countryCode := record[0]
		stateName := record[2]
		cityName := record[3]

		sourceLang, exists := languageMap[countryCode]
		if !exists {
			fmt.Println("Language not found for country code:", countryCode)
			continue
		}

		var stateNameBR, stateNameEN, cityNameBRFormula, cityNameENFormula string

		if translations, found := stateTranslations[stateName]; found {
			stateNameBR = translations[0]
			stateNameEN = translations[1]
		}

		if sourceLang == "en" {
			cityNameBRFormula = fmt.Sprintf(`=GOOGLETRANSLATE(D%d;"%s";"PT")`, rowNumber, sourceLang)
			cityNameENFormula = cityName
		} else if sourceLang == "pt" {
			cityNameBRFormula = cityName
			cityNameENFormula = fmt.Sprintf(`=GOOGLETRANSLATE(D%d;"%s";"EN")`, rowNumber, sourceLang)
		} else {
			cityNameBRFormula = fmt.Sprintf(`=GOOGLETRANSLATE(D%d;"%s";"PT")`, rowNumber, sourceLang)
			cityNameENFormula = fmt.Sprintf(`=GOOGLETRANSLATE(D%d;"%s";"EN")`, rowNumber, sourceLang)
		}

		newRecord := append(record, stateNameBR, stateNameEN, cityNameBRFormula, cityNameENFormula)

		if err := writer.Write(newRecord); err != nil {
			fmt.Println("Error writing record:", err)
			return
		}

		rowNumber++
	}

	fmt.Println("CSV file processed successfully!")
}
