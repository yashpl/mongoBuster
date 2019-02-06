# mongoBuster
Hunt Open MongoDB instances

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

|Flag| Description 
|---|
|--max-rate= (int)| Defines maximum rate at which packets are generated and sent. Default is 1000.
