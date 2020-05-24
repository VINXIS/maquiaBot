package gencommands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	config "../../config"
	"github.com/bwmarrin/discordgo"
)

// WeatherData holds information for the weather data given
type WeatherData struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Latitude       float64 `json:"lat"`
		Longtitude     float64 `json:"lon"`
		TimeZoneName   string  `json:"tz_id"`
		LocalTimeEpoch int64   `json:"localtime_epoch"`
	} `json:"location"`
	Weather struct {
		LastUpdatedEpoch int64   `json:"last_updated_epoch"`
		Celsius          float64 `json:"temp_c"`
		Fahrenheit       float64 `json:"temp_f"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
		} `json:"condition"`
		WindMPH         float64 `json:"wind_mph"`
		WindKPH         float64 `json:"wind_kph"`
		WindDeg         float64 `json:"wind_deg"`
		WindDir         string  `json:"wind_dir"`
		PressureMB      float64 `json:"pressure_mb"`
		PressurePSI     float64 `json:"pressure_in"`
		PrecipitationMM float64 `json:"precip_mm"`
		PrecipitationIN float64 `json:"precip_in"`
		Humidity        int     `json:"humidity"`
		Cloud           int     `json:"cloud"`
		FeelsLikeC      float64 `json:"feelslike_c"`
		FeelsLikeF      float64 `json:"feelslike_f"`
		VisibilityKM    float64 `json:"vis_km"`
		VisibilityMI    float64 `json:"vis_miles"`
		UV              float64 `json:"uv"`
		GustMPH         float64 `json:"gust_mph"`
		GustKPH         float64 `json:"gust_kph"`
	} `json:"current"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// Weather details the weather of the location given
func Weather(s *discordgo.Session, m *discordgo.MessageCreate) {
	weatherRegex, _ := regexp.Compile(`w(eather)?\s+(.+)`)
	paramRegex, _ := regexp.Compile(`&(.+)=(\S+)`)
	languageRegex, _ := regexp.Compile(`-l\s+(.+)`)

	// Get location
	if !weatherRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "No location given!")
		return
	}
	location := weatherRegex.FindStringSubmatch(m.Content)[2]
	location = paramRegex.ReplaceAllString(location, "")

	// Get language if given
	var lang string
	if languageRegex.MatchString(location) {
		lang = languageRegex.FindStringSubmatch(location)[1]
		location = languageRegex.ReplaceAllString(location, "")
	}

	location = strings.Replace(location, " ", "_", -1)

	// API request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://api.weatherapi.com/v1/current.json?key="+config.Conf.Weather+"&q="+strings.TrimSpace(location)+"&lang="+lang, nil)
	res, err := client.Do(req)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in obtaining weather information! Error 1")
		return
	}
	defer res.Body.Close()
	byteArray, err := ioutil.ReadAll(res.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in obtaining weather information! Error 2")
		return
	}

	// Parse Response
	var weatherData WeatherData
	err = json.Unmarshal(byteArray, &weatherData)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Bad weather request")
		return
	}

	// Check error
	if weatherData.Error.Code != 0 {
		s.ChannelMessageSend(m.ChannelID, weatherData.Error.Message)
		return
	}

	// Create embed
	embed := &discordgo.MessageEmbed{
		Description: weatherData.Weather.Condition.Text + " | **" + strconv.Itoa(weatherData.Weather.Cloud) + "%** Clouds\nUV Index: **" + strconv.FormatFloat(weatherData.Weather.UV, 'f', 1, 64) + "** | **" + strconv.Itoa(weatherData.Weather.Humidity) + "%** Humidity",
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: "https://raw.githubusercontent.com/gosquared/flags/master/flags/flags/flat/64/" + strings.Replace(weatherData.Location.Country, " ", "-", -1) + ".png",
			Name:    weatherData.Location.Name + ", " + weatherData.Location.Region + ", " + weatherData.Location.Country,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "http:" + weatherData.Weather.Condition.Icon,
		},
		Fields: []*discordgo.MessageEmbedField{
			{ // Temperature
				Name: "Temperature",
				Value: "**" + strconv.FormatFloat(weatherData.Weather.Celsius, 'f', 1, 64) + "C (" + strconv.FormatFloat(weatherData.Weather.Fahrenheit, 'f', 1, 64) + "F)**\n" +
					"Feels like: **" + strconv.FormatFloat(weatherData.Weather.FeelsLikeC, 'f', 1, 64) + "C (" + strconv.FormatFloat(weatherData.Weather.FeelsLikeF, 'f', 1, 64) + "F)**",
			},
			{ // Wind
				Name: "Wind",
				Value: "**" + strconv.FormatFloat(weatherData.Weather.WindKPH, 'f', 1, 64) + "KPH (" + strconv.FormatFloat(weatherData.Weather.WindMPH, 'f', 1, 64) + "MPH)**\n" +
					"Direction: **" + weatherData.Weather.WindDir + " (" + strconv.FormatFloat(weatherData.Weather.WindDeg, 'f', 0, 64) + ")**\n" +
					"Gust speed: **" + strconv.FormatFloat(weatherData.Weather.GustKPH, 'f', 1, 64) + "KPH (" + strconv.FormatFloat(weatherData.Weather.GustMPH, 'f', 1, 64) + "MPH)**",
			},
			{ // Other
				Name: "Other",
				Value: "Pressure **" + strconv.FormatFloat(weatherData.Weather.PressureMB/1000, 'f', 2, 64) + "Bar (" + strconv.FormatFloat(weatherData.Weather.PressurePSI, 'f', 2, 64) + "psi)**\n" +
					"Precipitation: **" + strconv.FormatFloat(weatherData.Weather.PrecipitationMM/10, 'f', 2, 64) + "cm (" + strconv.FormatFloat(weatherData.Weather.PrecipitationIN, 'f', 2, 64) + "in)**\n" +
					"Visibility **" + strconv.FormatFloat(weatherData.Weather.VisibilityKM, 'f', 2, 64) + "km (" + strconv.FormatFloat(weatherData.Weather.VisibilityMI, 'f', 2, 64) + "mi)**",
			},
		},
	}

	// Exceptions
	if weatherData.Location.Country == "United States of America" {
		embed.Author.IconURL = "https://raw.githubusercontent.com/gosquared/flags/master/flags/flags/flat/64/United-States.png"
	}
	if weatherData.Location.Region == "" {
		embed.Author.Name = weatherData.Location.Name + ", " + weatherData.Location.Country
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
