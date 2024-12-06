package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"os"
	"log"
)

func checkForUpdates(){
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	//Send GET request with user headers(Request won't work without headers set)
	req, err := http.NewRequest("GET", "https://www.minecraft.net/en-us/download/server/bedrock", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	//Open, read then close the response
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//regex to find latest version in the body
	re := regexp.MustCompile(`\d*\.*\d*\.\d*\.\d*\.zip`)
	latestVersion := strings.Trim(re.FindString(string(body)), ".zip")
	
	installedVers, err := os.ReadFile(homeDir + "/mc-bedrock-autoupdater/vers.txt")
	if(err != nil){
		log.Fatal(err)
	}
	
	if latestVersion == string(installedVers) {
		fmt.Println("same version")
	} 
}



func main() {
	checkForUpdates()
}
