package main

import (
	
	"io"
	"net/http"
	"regexp"
	"strings"
	"os"
	"log"
	"time"
	"fmt"
	
)

func updateServer(latestVersZip string, homeDir string) {
	serverURL := "https://www.minecraft.net/bedrockdedicatedserver/bin-linux/bedrock-server-" + latestVersZip
	client := &http.Client{}
	//Request the latest server file
	req, err := http.NewRequest("GET", serverURL, nil)
		if err != nil {
			log.Fatal(err)
		}
	

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
        return
    } 

	fmt.Println("status", resp.Status)
    //create the file and then write its content to disk
	serverFile, err := os.Create(homeDir + "/bedrock-server-" + latestVersZip)
	if err != nil {
		log.Fatal(err)
	}
	defer serverFile.Close()
	if _, err := io.Copy(serverFile, resp.Body); err != nil {log.Fatal(err)}

	//TODO
	//Backup the running "bedrock-server"
	//stop the minecraft server?
	//Add logic to install new server:
	//	Delete running server, unzip the latest version into "bedrock-server"
	//	copy server.properties, worlds, permission.json and allowlist.json, recursively!!! from the backup and into the new bedrock-server
	//	delete the latest .zip
	//	Run the bedrock server with screen.
	//	Log that the server was updated sucessfully
}


func checkForUpdates(){
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	//Send GET request with user headers(Request won't work without headers set)
	req, err := http.NewRequest("GET", "https://www.minecraft.net/en-us/download/server/bedrock", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	//Open, read then close the response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
        return
    } 
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//regex to find latest version in the body
	latestVersionZipRegex := regexp.MustCompile(`\d*\.*\d*\.\d*\.\d*\.zip`)
	latestVersionZip := latestVersionZipRegex.FindString(string(body))
	if latestVersionZip == "" {
		log.Fatal("Could not find the latest version zip in the response")
	}

	latestVersion := strings.TrimSuffix(latestVersionZip, ".zip")
	
	installedVers, err := os.ReadFile(homeDir + "/mc-bedrock-autoupdater/vers.txt")
	if(err != nil){
		log.Fatal(err)
	}
	
	if latestVersion == string(installedVers) {
		//open the log file, append the status then close it
		f, err := os.OpenFile(homeDir + "/mc-bedrock-autoupdater/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		t := time.Now()
		if _, err := f.WriteString(t.String()); err != nil {log.Fatal(err)}
		if _, err := f.WriteString(" - Latest version installed: " + latestVersion); err != nil {log.Fatal(err)}
		if _, err := f.WriteString("\n"); err != nil {log.Fatal(err)}
		if err := f.Close(); err != nil {log.Fatal(err)}
		
		
	} else {
		updateServer(latestVersionZip, homeDir)
	}
}



func main() {
	checkForUpdates()
}
