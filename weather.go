package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/joho/godotenv"
)

type Geocoder struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Bounds struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"bounds"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
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
		panic("Error loading .env file")
	}

	DARK_SKY_API_KEY := os.Getenv("DARK_SKY_API_KEY")
	GOOGLE_MAPS_API_KEY := os.Getenv("GOOGLE_MAPS_API_KEY")

	// zipCode := flag.String("zip", randomdata.PostalCode("SE"), "Zip code to obtain weather forecast")
	location := flag.String("location", randomdata.PostalCode("SE"), "Enter a location (e.g. address, zip) to obtain weather forecast")
	flag.Parse()

	// if len(*zipCode) != 5 {
	// 	panic("Please enter a 5 digit zip code")
	// }
	// Look at status from maps request

	// https://maps.googleapis.com/maps/api/geocode/json?address=[ZIP_OR_ADDRESS]&key=[YOUR_API_KEY]
	googleMapsURL := url.URL{Scheme: "https", Host: "maps.googleapis.com", Path: "maps/api/geocode/json"}
	q := googleMapsURL.Query()
	locationEncoded := url.QueryEscape(*location)
	// fmt.Println(locationEncoded)
	q.Set("address", locationEncoded)
	q.Set("key", GOOGLE_MAPS_API_KEY)
	googleMapsURL.RawQuery = q.Encode()
	googleMapsData := Geocoder{}
	makeGetRequest(googleMapsURL.String(), &googleMapsData)

	if googleMapsData.Status != "OK" {
		log.Fatalln("Invalid location. Please try again.")
	}

	lat := googleMapsData.Results[0].Geometry.Location.Lat
	lon := googleMapsData.Results[0].Geometry.Location.Lng
	address := googleMapsData.Results[0].FormattedAddress

	// https://api.darksky.net/forecast/[KEY]/[LATITUDE],[LONGITUDE]
	urlPath := fmt.Sprintf("forecast/%v/%v,%v", DARK_SKY_API_KEY, lat, lon)
	darkSkyURL := url.URL{Scheme: "https", Host: "api.darksky.net", Path: urlPath}
	// fmt.Println("Request URL:", darkSkyURL.String())

	darkSkyData := Forecast{}
	makeGetRequest(darkSkyURL.String(), &darkSkyData)
	current := darkSkyData.Currently
	icon := parseIcon(current.Icon)
	// fmt.Println(icon)

	time := parseTime(current.Time, darkSkyData.Timezone)
	// fmt.Println(time)

	fmt.Printf("Current conditions for %v (last updated on %v %v %v %v:%v)\n", address, time.Weekday(), time.Month(), time.Day(), time.Hour(), time.Minute())
	fmt.Println("--------------------------------")
	fmt.Printf("%v°   %v %v\n", current.Temperature, current.Summary, icon)
	fmt.Printf("Wind: %v mph   Humidity: %v%%   Dew Pt: %v°   UV Index: %v   Visibility: %v+ mi   Pressure: %v mb\n", math.Round(current.WindSpeed), math.Round(current.Humidity*100), math.Round(current.DewPoint), current.UvIndex, current.Visibility, math.Round(current.Pressure))
	fmt.Println("--------------------------------")
}

func parseTime(currentTime int, timezone string) time.Time {
	parsedTime, err := strconv.ParseInt(strconv.Itoa(currentTime), 10, 64)
	if err != nil {
		log.Fatalln("Error with time conversion:", err)
	}
	utcTime := time.Unix(parsedTime, 0)

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatalln("Error with timezone conversion:", err)
	}

	return utcTime.In(loc)
}

func makeGetRequest(url string, target interface{}) {
	res, err := httpClient.Get(url)
	if err != nil {
		log.Fatalln("Error in request to", url)
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
