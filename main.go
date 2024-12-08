package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func updateServer(latestVersZip string, homeDir string) {
	serverURL := "https://www.minecraft.net/bedrockdedicatedserver/bin-linux/" + latestVersZip
	client := &http.Client{}

	fmt.Println("Requesting to download server file...")
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
	fmt.Println("Writing new server to disk")
	//create the new file and then write its content to disk

	serverFile, err := os.Create(homeDir + "/" + latestVersZip)
	if err != nil {
		log.Fatal(err)
	}
	defer serverFile.Close()
	if _, err := io.Copy(serverFile, resp.Body); err != nil {
		log.Fatal(err)
	}

	//TODO
	//Backup the running "bedrock-server"
	//stop the minecraft server? *pidof bedrock_server
	//Add logic to install new server:
	//	Delete running server, unzip the latest version into "bedrock-server"
	//	copy server.properties, worlds, permission.json and allowlist.json, recursively!!! from the backup and into the new bedrock-server
	//	delete the latest .zip
	//	Run the bedrock server with screen.
	//	Log that the server was updated sucessfully

	fmt.Println("Backing up current server")
	t := time.Now()
	timeF := t.Format("01-02-2006")
	backupDir := homeDir + "/bedrock-server-backup-" + timeF //change to wherever if needed
	bedrockServerDir := homeDir + "/bedrock-server"
	latestServerDir := homeDir + "/" + latestVersZip
	cmd := exec.Command("cp", "--recursive", bedrockServerDir, backupDir)
	cmd.Run()

	// fmt.Println("Killing server...")
	// pid, err := exec.Command("pidof", "bedrock_server").Output()
	// if err != nil {
	// 	log.Fatal("Error killing the server, maybe its not running? ", err)
	// }
	// //converting pid from byte slice to int...
	// pidstr := strings.TrimSpace(string(pid))

	// pidint, err := strconv.Atoi(pidstr)
	// fmt.Println("PID:", pidint)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// //Find process and kill it
	// process, err := os.FindProcess(pidint)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = process.Signal(syscall.SIGTERM)
	// if err != nil {
	// 	log.Fatal("Error sending signal", err)
	// }

	//Delete the current bedrock dir..
	fmt.Println("Deleting current bedrock server directory")

	err = os.RemoveAll(bedrockServerDir)
	if err != nil {
		log.Fatal("Error deleting bedrock directory: ", err)
	}

	//unzip the new server version into a new bedrock-server directory
	fmt.Println("Unzipping new server into a new bedrock-server directory")

	err = os.Mkdir(bedrockServerDir, 0777)
	if err != nil {
		log.Fatal("Error creating new bedrock server directory: ", err)
	}

	cmd = exec.Command("unzip", latestServerDir, "-d", bedrockServerDir)

	err = cmd.Run()
	if err != nil {
		log.Fatal("Error unzipping, is the zip in the right directory?", err)
	}

	//copy important directories to the new bedrock-server
	fmt.Println(backupDir)
	fmt.Println("Coping files from backup directory")
	cmd = exec.Command("cp", "-r", backupDir+"/worlds", backupDir+"/server.properties", backupDir+"/permission.json", backupDir+"/allowlist.json", bedrockServerDir)
	cmd.Run()

	//Deleting the downloaded bedrock zip
	err = os.RemoveAll(latestServerDir)
	if err != nil {
		log.Fatal("Error deleting downloaded bedrock directory: ", err)
	}

	//Run the new bedrock server

	fmt.Println("Attemping to run the bedrock server...")
	fmt.Println(bedrockServerDir)
	cmd = exec.Command("screen", bedrockServerDir+"/./bedrock_server")
	err = cmd.Run()
	if err != nil {
		log.Fatal("Error running the bedrock server", err)
	}
	fmt.Println("Bedrock server started successfully")

}
func checkForUpdates() {
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
	latestVersionZipRegex := regexp.MustCompile(`bedrock-server-\d*\.*\d*\.\d*\.\d*\.zip`)
	latestVersionZip := latestVersionZipRegex.FindString(string(body))
	if latestVersionZip == "" {
		log.Fatal("Could not find the latest version zip in the response")
	}

	latestVersion := strings.TrimSuffix(strings.TrimPrefix(latestVersionZip, "bedrock-server-"), ".zip")

	installedVers, err := os.ReadFile(homeDir + "/mc-bedrock-autoupdater/vers.txt")
	if err != nil {
		log.Fatal(err)
	}

	if latestVersion == string(installedVers) {
		//open the log file, append the status then close it
		f, err := os.OpenFile(homeDir+"/mc-bedrock-autoupdater/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		t := time.Now()
		if _, err := f.WriteString(t.String()); err != nil {
			log.Fatal(err)
		}
		if _, err := f.WriteString(" - Latest version installed: " + latestVersion); err != nil {
			log.Fatal(err)
		}
		if _, err := f.WriteString("\n"); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}

	} else {
		updateServer(latestVersionZip, homeDir)
	}
}

func main() {
	checkForUpdates()
}
