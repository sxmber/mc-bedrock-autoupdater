# mc-bedrock-autoupdater
Program to automate the process of updating Minecraft Bedrock Edition Server

I'm working on making this more accessible but for now this script assumes your /bedrock-server and this repo is in your home directory and you have installed Golang.

# Requirements
Golang

Linux

An existing Minecraft Bedrock Server in your home directory

# Installation

  1. **Install golang**

  >https://go.dev/doc/install

  2. **git clone this repo into your home directory**

```console
git clone https://github.com/sxmber/mc-bedrock-autoupdater
```
  3. **Run the initial.sh script**

```console
cd mc-bedrock-autoupdater
chmod +x initial.sh
./initial.sh
```

  4. **Install the golang binary**

```console 
go install
```

  5. **Write your Minecraft Bedrock Server into ~/mc-be-logs/vers.txt**

6. **Run the binary manually OR configure a cronjob to do it automatically**

```console
~/go/bin/mc-bedrock-autoupdater


