# Go "telnet" Honeypot
* A lightweight TCP honeypot that mimics a BusyBox telnet login and logs connection attempts.
* I say "telnet" because currently it only handles raw tcp and theres no IAC negotiation.
* Make sure you only run this in a container with the Dockerfile I provided.
## HOW-TO
1. Build the image:

    ​```docker build -t honeypot .​```
1. Run it:

    ​```docker run --rm -p 23:2323 honeypot​```
1. Connect (choose "all interfaces" when prompted):

    ```nc [local_ip] 23```
Note: This is to reach it from your home network to really deploy this you need to port forward the host machine to the internet (would not do with the current state of this). 
## TO-DO
* IAC Negotiation, right now if you try to connect via telnet it'll just immediately throw an incorrect login on any input.
* Username field.
* Rate limiting IP's.
* Password try limit.
* Json logs.
* Database and log analysis containers.
