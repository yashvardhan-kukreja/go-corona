package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"strconv"
	//"time"
)

type MailServer struct {
	host string
	port string
}

type CovidStatsObject struct {
	Province   string `json:"province"`
	Country    string `json:"country"`
	LastUpdate string `json:"lastUpdate"`
	Confirmed  int    `json:"confirmed"`
	Deaths     int    `json:"deaths"`
	Recovered  int    `json:"recovered"`
}
type CovidDataObject struct {
	LastChecked  string             `json:"lastChecked"`
	Covid19Stats []CovidStatsObject `json:"covid19Stats"`
}
type CovidAPIResponse struct {
	Err        bool            `json:"error"`
	StatusCode int             `json:"statusCode"`
	Message    string          `json:"message"`
	Data       CovidDataObject `json:"data"`
}

func (mailServer MailServer) getMailServerAddress() string {
	return mailServer.host + ":" + mailServer.port
}

func createAndReturnMailServer(host string, port string) MailServer {
	return MailServer{host: host, port: port}
}

func getCovidData(country string) (data CovidStatsObject, err error) {

	url := "https://covid-19-coronavirus-statistics.p.rapidapi.com/v1/stats?country=" + country

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-host", "covid-19-coronavirus-statistics.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "916bd741ccmshb29b580f9e9981ap11a6f7jsn73d132b41cdf")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	bodyByteValue, _ := ioutil.ReadAll(res.Body)

	var body CovidAPIResponse
	json.Unmarshal(bodyByteValue, &body)

	if body.Message != "OK" {
		fmt.Printf("Country not found")
		return CovidStatsObject{}, errors.New("Country not found")
	}
	if body.Err || body.StatusCode != 200 {
		fmt.Printf("An error occurred!")
		return CovidStatsObject{}, errors.New("An error occurred")
	}

	return body.Data.Covid19Stats[0], nil
}

func sendMail(email string, password string, body string) error {

	gmailServer := createAndReturnMailServer("smtp.gmail.com", "587")
	gmailServerAddress := gmailServer.getMailServerAddress()

	sender := "go_corona@gmail.com"
	receivers := []string{email}

	finalBody := fmt.Sprintf("From: %s\nTo: %s\nSubject: COVID-19 alert from go-corona! \n\n%s", sender, email, body)

	msg := []byte(finalBody)

	auth := smtp.PlainAuth("", email, password, gmailServer.host)

	err := smtp.SendMail(gmailServerAddress, auth, sender, receivers, msg)

	if err != nil {
		return err
	}
	return nil
}

func sendAlert(emailFlag string, passwordFlag string, countryFlag string, timeInSecondsFlag int) error {
	if emailFlag == "" {
		return errors.New("Please enter the E-mail ID where you want to receive the alerts")
	}
	if passwordFlag == "" {
		return errors.New("Please enter the password of the email account")
	}
	if countryFlag == "" {
		return errors.New("Please enter the country name of the country for which you want to find COVID-19 stats")
	}
	if timeInSecondsFlag <= 5 {
		return errors.New("Please make sure that the time rate is greater than 5 seconds")
	}

	data, err := getCovidData(countryFlag)

	if err != nil {
		return err
	}

	finalString := fmt.Sprintf("Country - %s \nTotal Confirmed Cases - %d \nTotal Deaths - %d \nTotal Recoveries - %d", data.Country, data.Confirmed, data.Deaths, data.Recovered)

	err = sendMail(emailFlag, passwordFlag, finalString)

	if err != nil {
		return err
	}
	return nil
}

func killCurrentProcess(pid int) error {
	_, err := exec.Command("kill", strconv.Itoa(pid)).CombinedOutput()
	if err!=nil {
		return err
	}
	return nil
}

func getEnvVariable() (string, error) {
	outputString := os.Getenv("GOCORONA")
	if outputString == "" {
		return "", errors.New("Environment variable not found!")
	}
	return outputString, nil
}

// func runFuncAtRate(rateInSeconds time.Duration, f func(email string, password string, body string)) {
// 	rate := rateInSeconds*time.Millisecond

// 	for i := range 
// }

func main() {
	emailPtr := flag.String("email", "", "Enter the email ID where you want to receive the alerts :)")
	passwordPtr := flag.String("password", "", "Enter the password corresponding to that email account")
	countryPtr := flag.String("country", "India", "Enter the country for which you want the covid stats. (Default: India)")
	timeInSecondsPtr := flag.Int("timeInSeconds", 300, "Enter the rate at which you want to receive alerts (Default: 300 seconds)")
	killCoronaPtr := flag.String("kill", "-", "To kill the existing go-corona alert process!")
	flag.Parse()

	pid := os.Getpid()



	// If the kill flag is passed
	if *killCoronaPtr == "yes" {
		existingPid, _ := getEnvVariable()
		existingPidInt, _ := strconv.Atoi(existingPid)
		err := killCurrentProcess(existingPidInt)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
			return
		}
		fmt.Printf("Killed the existing go-corona process")
		return
	}


	err := os.Setenv("GOCORONA", strconv.Itoa(pid))
	if err!=nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}

	//If the kill flag in not passed in the command
	err = sendAlert(*emailPtr, *passwordPtr, *countryPtr, *timeInSecondsPtr)
	if err!=nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}
}
