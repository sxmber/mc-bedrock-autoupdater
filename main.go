package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func updateServer(latestVersZip string, homeDir string) {
	serverURL := "https://www.minecraft.net/bedrockdedicatedserver/bin-linux/" + latestVersZip
	client := &http.Client{}
	latestVersTrim := strings.TrimSuffix(strings.TrimPrefix(latestVersZip, "bedrock-server-"), ".zip")

	//Set the logger to send errors to log.txt
	f, err := os.OpenFile(homeDir+"/mc-be-logs/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Error opening log file", err)
	}
	log.SetOutput(f)

	//Request the latest server file
	fmt.Println("Requesting to download server file...")
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

	//create the new file and then write its content to disk
	fmt.Println("Writing new server to disk")
	serverFile, err := os.Create(homeDir + "/" + latestVersZip)
	if err != nil {
		log.Fatal(err)
	}
	defer serverFile.Close()
	if _, err := io.Copy(serverFile, resp.Body); err != nil {
		log.Fatal(err)
	}

	//Backup the current server
	fmt.Println("Backing up current server")
	t := time.Now()
	timeF := t.Format("01-02-2006")
	backupDir := homeDir + "/bedrock-server-backup-" + timeF //change to wherever if needed
	bedrockServerDir := homeDir + "/bedrock-server"
	latestServerDir := homeDir + "/" + latestVersZip
	cmd := exec.Command("cp", "--recursive", bedrockServerDir, backupDir)
	cmd.Run()

	//Kill the running server
	fmt.Println("Killing server...")
	pid, err := exec.Command("pidof", "bedrock_server").Output()
	if err != nil {
		log.Fatal("Error killing the server, maybe its not running? ", err)
	}
	//converting pid from byte slice to int...
	pidstr := strings.TrimSpace(string(pid))

	pidint, err := strconv.Atoi(pidstr)
	fmt.Println("PID:", pidint)
	if err != nil {
		log.Fatal(err)
	}

	//Find process and kill it
	process, err := os.FindProcess(pidint)
	if err != nil {
		log.Fatal(err)
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		log.Fatal("Error sending signal", err)
	}

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
	fmt.Println("Copying files from backup directory")
	cmd = exec.Command("cp", "-r", backupDir+"/worlds", backupDir+"/server.properties", backupDir+"/permissions.json", backupDir+"/allowlist.json", bedrockServerDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal("Error copying files from backup")
	}

	//Create tar ball of the backup directory then delete the original
	fmt.Println("Creating tar ball of backup dir")
	cmd = exec.Command("sh", "-c", "tar -czvf "+backupDir+".tar.gz "+backupDir+" && rm -rf "+backupDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal("Error creating tar ball")
	}

	//Deleting the downloaded bedrock zip
	err = os.RemoveAll(latestServerDir)
	if err != nil {
		log.Fatal("Error deleting downloaded bedrock directory: ", err)
	}

	//Run the new bedrock server

	fmt.Println("Attemping to run the bedrock server...")
	err = os.Chdir(bedrockServerDir)
	if err != nil {
		log.Fatal("Error changing directory to bedrock-server:", err)
	}

	cmd = exec.Command("screen", "-dmS", "bedrock_server_session", "./bedrock_server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal("Error running the bedrock server:", err)
	}
	fmt.Println("Bedrock server started successfully")

	//Log that the server updated successfully
	if _, err := f.WriteString(t.String()); err != nil {
		log.Fatal(err)
	}
	if _, err := f.WriteString(" - Bedrock server updated successfully to: " + latestVersTrim); err != nil {
		log.Fatal(err)
	}
	if _, err := f.WriteString("\n"); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	//Update the current version in vers.txt
	f, err = os.OpenFile(homeDir+"/mc-be-logs/vers.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Error opening the version file", err)
	}

	if _, err := f.WriteString(latestVersTrim); err != nil {
		log.Fatal("Error updating vers.txt file: ", err)
	}

	if err := f.Close(); err != nil {
		log.Fatal("Error closing the version file")
	}

}

func checkForUpdates() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting the user home directory", err)
	}
	//open the log file and set the logger to it
	f, err := os.OpenFile(homeDir+"/mc-be-logs/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	//Send GET request with user headers(Request won't work without headers set)
	client := &http.Client{}
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
	//regex to find latest version in the response body
	latestVersionZipRegex := regexp.MustCompile(`bedrock-server-\d*\.*\d*\.\d*\.\d*\.zip`)
	latestVersionZip := latestVersionZipRegex.FindString(string(body))
	if latestVersionZip == "" {
		log.Fatal("Could not find the latest version zip in the response")
	}

	latestVersion := strings.TrimSuffix(strings.TrimPrefix(latestVersionZip, "bedrock-server-"), ".zip")
	installedVers, err := os.ReadFile(homeDir + "/mc-be-logs/vers.txt")
	if err != nil {
		log.Fatal(err)
	}
	if string(installedVers) == "" {
		log.Fatal("Set the installed minecraft version in vers.txt", err)
	}

	//Check if the bedrock server needs to be updated and log it to the log.txt file
	if latestVersion == string(installedVers) {
		t := time.Now()
		if _, err := f.WriteString(t.String()); err != nil {
			log.Fatal(err)
		}
		if _, err := f.WriteString(" - Bedrock server is up to date: " + latestVersion); err != nil {
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
