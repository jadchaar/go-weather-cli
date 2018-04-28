package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jasonwinn/geocoder"
	"github.com/joho/godotenv"
)

type MapQuest struct {
	Info struct {
		Statuscode int `json:"statuscode"`
		Copyright  struct {
			Text         string `json:"text"`
			ImageURL     string `json:"imageUrl"`
			ImageAltText string `json:"imageAltText"`
		} `json:"copyright"`
		Messages []interface{} `json:"messages"`
	} `json:"info"`
	Options struct {
		MaxResults        int  `json:"maxResults"`
		ThumbMaps         bool `json:"thumbMaps"`
		IgnoreLatLngInput bool `json:"ignoreLatLngInput"`
	} `json:"options"`
	Results []struct {
		ProvidedLocation struct {
			Location string `json:"location"`
		} `json:"providedLocation"`
		Locations []struct {
			Street             string `json:"street"`
			AdminArea6         string `json:"adminArea6"`
			AdminArea6Type     string `json:"adminArea6Type"`
			AdminArea5         string `json:"adminArea5"`
			AdminArea5Type     string `json:"adminArea5Type"`
			AdminArea4         string `json:"adminArea4"`
			AdminArea4Type     string `json:"adminArea4Type"`
			AdminArea3         string `json:"adminArea3"`
			AdminArea3Type     string `json:"adminArea3Type"`
			AdminArea1         string `json:"adminArea1"`
			AdminArea1Type     string `json:"adminArea1Type"`
			PostalCode         string `json:"postalCode"`
			GeocodeQualityCode string `json:"geocodeQualityCode"`
			GeocodeQuality     string `json:"geocodeQuality"`
			DragPoint          bool   `json:"dragPoint"`
			SideOfStreet       string `json:"sideOfStreet"`
			LinkID             string `json:"linkId"`
			UnknownInput       string `json:"unknownInput"`
			Type               string `json:"type"`
			LatLng             struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"latLng"`
			DisplayLatLng struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"displayLatLng"`
			MapURL string `json:"mapUrl"`
		} `json:"locations"`
	} `json:"results"`
}

type Forecast struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Currently struct {
		Time                 int     `json:"time"`
		Summary              string  `json:"summary"`
		Icon                 string  `json:"icon"`
		NearestStormDistance float64 `json:"nearestStormDistance"`
		NearestStormBearing  float64 `json:"nearestStormBearing"`
		PrecipIntensity      float64 `json:"precipIntensity"`
		PrecipProbability    float64 `json:"precipProbability"`
		Temperature          float64 `json:"temperature"`
		ApparentTemperature  float64 `json:"apparentTemperature"`
		DewPoint             float64 `json:"dewPoint"`
		Humidity             float64 `json:"humidity"`
		Pressure             float64 `json:"pressure"`
		WindSpeed            float64 `json:"windSpeed"`
		WindGust             float64 `json:"windGust"`
		WindBearing          float64 `json:"windBearing"`
		CloudCover           float64 `json:"cloudCover"`
		UvIndex              float64 `json:"uvIndex"`
		Visibility           float64 `json:"visibility"`
		Ozone                float64 `json:"ozone"`
	} `json:"currently"`
	Minutely struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
		Data    []struct {
			Time              int     `json:"time"`
			PrecipIntensity   float64 `json:"precipIntensity"`
			PrecipProbability float64 `json:"precipProbability"`
		} `json:"data"`
	} `json:"minutely"`
	Hourly struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
		Data    []struct {
			Time                int     `json:"time"`
			Summary             string  `json:"summary"`
			Icon                string  `json:"icon"`
			PrecipIntensity     float64 `json:"precipIntensity"`
			PrecipProbability   float64 `json:"precipProbability"`
			Temperature         float64 `json:"temperature"`
			ApparentTemperature float64 `json:"apparentTemperature"`
			DewPoint            float64 `json:"dewPoint"`
			Humidity            float64 `json:"humidity"`
			Pressure            float64 `json:"pressure"`
			WindSpeed           float64 `json:"windSpeed"`
			WindGust            float64 `json:"windGust"`
			WindBearing         float64 `json:"windBearing"`
			CloudCover          float64 `json:"cloudCover"`
			UvIndex             float64 `json:"uvIndex"`
			Visibility          float64 `json:"visibility"`
			Ozone               float64 `json:"ozone"`
			PrecipType          string  `json:"precipType,omitempty"`
		} `json:"data"`
	} `json:"hourly"`
	Daily struct {
		Summary string `json:"summary"`
		Icon    string `json:"icon"`
		Data    []struct {
			Time                        int     `json:"time"`
			Summary                     string  `json:"summary"`
			Icon                        string  `json:"icon"`
			SunriseTime                 float64 `json:"sunriseTime"`
			SunsetTime                  float64 `json:"sunsetTime"`
			MoonPhase                   float64 `json:"moonPhase"`
			PrecipIntensity             float64 `json:"precipIntensity"`
			PrecipIntensityMax          float64 `json:"precipIntensityMax"`
			PrecipIntensityMaxTime      float64 `json:"precipIntensityMaxTime"`
			PrecipProbability           float64 `json:"precipProbability"`
			PrecipType                  string  `json:"precipType,omitempty"`
			TemperatureHigh             float64 `json:"temperatureHigh"`
			TemperatureHighTime         float64 `json:"temperatureHighTime"`
			TemperatureLow              float64 `json:"temperatureLow"`
			TemperatureLowTime          float64 `json:"temperatureLowTime"`
			ApparentTemperatureHigh     float64 `json:"apparentTemperatureHigh"`
			ApparentTemperatureHighTime float64 `json:"apparentTemperatureHighTime"`
			ApparentTemperatureLow      float64 `json:"apparentTemperatureLow"`
			ApparentTemperatureLowTime  float64 `json:"apparentTemperatureLowTime"`
			DewPoint                    float64 `json:"dewPoint"`
			Humidity                    float64 `json:"humidity"`
			Pressure                    float64 `json:"pressure"`
			WindSpeed                   float64 `json:"windSpeed"`
			WindGust                    float64 `json:"windGust"`
			WindGustTime                float64 `json:"windGustTime"`
			WindBearing                 float64 `json:"windBearing"`
			CloudCover                  float64 `json:"cloudCover"`
			UvIndex                     float64 `json:"uvIndex"`
			UvIndexTime                 float64 `json:"uvIndexTime"`
			Visibility                  float64 `json:"visibility,omitempty"`
			Ozone                       float64 `json:"ozone"`
			TemperatureMin              float64 `json:"temperatureMin"`
			TemperatureMinTime          float64 `json:"temperatureMinTime"`
			TemperatureMax              float64 `json:"temperatureMax"`
			TemperatureMaxTime          float64 `json:"temperatureMaxTime"`
			ApparentTemperatureMin      float64 `json:"apparentTemperatureMin"`
			ApparentTemperatureMinTime  float64 `json:"apparentTemperatureMinTime"`
			ApparentTemperatureMax      float64 `json:"apparentTemperatureMax"`
			ApparentTemperatureMaxTime  float64 `json:"apparentTemperatureMaxTime"`
		} `json:"data"`
	} `json:"daily"`
	Flags struct {
		Sources     []string `json:"sources"`
		IsdStations []string `json:"isd-stations"`
		Units       string   `json:"units"`
	} `json:"flags"`
	Offset float64 `json:"offset"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	darkSkyApiKey := os.Getenv("DARK_SKY_API_KEY")
	geocoder.SetAPIKey(os.Getenv("MAPQUEST_API_KEY"))

	zipCode := flag.String("zip", "48108", "Zip code to obtain weather forecast")
	flag.Parse()

	if len(*zipCode) != 5 {
		log.Fatalln("Please enter a 5 digit zip code")
	}

	fmt.Println("Obtaining weather for", *zipCode)

	lat, lon, err := geocoder.Geocode(*zipCode)
	if err != nil {
		log.Fatalln("Error with geocoder")
	}

	// https://api.darksky.net/forecast/[key]/[latitude],[longitude]
	urlPath := fmt.Sprintf("forecast/%v/%v,%v", darkSkyApiKey, lat, lon)
	darkSkyURL := (url.URL{Scheme: "https", Host: "api.darksky.net", Path: urlPath})
	// fmt.Println("Request URL:", darkSkyURL.String())

	data := Forecast{}
	weatherRequest(darkSkyURL.String(), &data)
	current := data.Currently
	icon := parseIcon(current.Icon)
	fmt.Println(icon)
}

func weatherRequest(darkSkyURL string, target interface{}) {
	res, err := httpClient.Get(darkSkyURL)
	if err != nil {
		log.Fatalln("Error request to", darkSkyURL)
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(target)

	if err != nil {
		log.Fatalln("Error in JSON decoding:", err)
	}
}

func parseIcon(icon string) string {
	switch icon {
	case "clear-day":
		return "☀️"
	case "clear-night":
		return "🌚"
	case "rain":
		return "☔️"
	case "snow":
		return "❄️"
	case "sleet":
		return "🌨"
	case "wind":
		return "🌬"
	case "fog":
		return "🌫"
	case "cloudy":
		return "☁️"
	case "partly-cloudy-day":
		return "⛅️"
	case "partly-cloudy-night":
		return "🌚"
	default:
		return "⭐️"
	}
}
