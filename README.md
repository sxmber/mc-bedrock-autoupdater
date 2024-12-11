# mc-bedrock-autoupdater
Program to automate the process of updating Minecraft Bedrock Edition Server

I'm working on making this more accessible but for now this script assumes your /bedrock-server and this repo is in your home directory and you have installed Golang.

# Requirements
Golang

Linux

# Installation

  1. **Install golang**

  >https://go.dev/doc/install

  2. **git clone this repo into your home directory**

```console
git clone https://github.com/sxmber/mc-bedrock-autoupdater
```

  3. **cd into the repo and build the golang binary**

```console 
cd mc-bedrock-autoupdater && go build
```

  4. **Write the latest bedrock server version installed into vers.txt**

5. **Run the binary manually OR configure a cronjob to do it automatically**

```console
./mc-bedrock-autoupdater


