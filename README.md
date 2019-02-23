# mongoBuster
Hunt Open MongoDB instances!

### Features

* Worlds fastest and most efficient scanner ( Uses Masscan ).
* Scans entire internet by default, So fire the tool and chill.
* Hyper efficient - Uses Go-routines which are even lighter than threads.

### Pre-Requisites - 

* Go language ( sudo apt install golang )
* Masscan ( sudo apt install masscan )
* Tested on Ubuntu & Kali linux

### How to install and run - 

```
git clone https://github.com/yashpl/mongoBuster.git

cd mongoBuster

go build mongobuster.go

sudo ./mongobuster
```

Note: Run it with sudo as Masscan requires sudo access.

### Flags - 

|Flag| Description |
|---|---|
|--max-rate= (int)| Defines maximum rate at which packets are generated and sent. Default is 100.|
|--out-file= (string)| Name of file to which vulnerable IPs will be exported.|
|-v| Display error msgs from non-vulnerable servers| 

### NOTE - 

Using ridiculous values for ```max-rate``` flag like 10000+ will *most likely* bring down your own network infrastructure.

Recommended value is to start with ```--max-rate 500``` for consumer Gigabit routers.


#### Happy Hunting ;)

Final Note :- If you find bunch of insecure insances, ( which you will! ) You might wanna explore them with GUI tools like - [Robo 3t](https://robomongo.org/)


Please report these insecure instances to their respective owners, Lets make a safer internet together <3.