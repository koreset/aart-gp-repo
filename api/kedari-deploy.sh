#!/usr/bin/env bash

echo "Building binary for linux deploy"
env GOOS=linux GOARCH=amd64 go build -o aart_api

# Create migrations directory if it doesn't exist
echo "Ensuring migrations directory exists"
mkdir -p migrations


# Copy migrations directory
#echo "Copying migrations to server"
ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo mkdir -p /home/aart/api1/migrations"
scp -i "aart_key_access.pem" -r migrations/* ubuntu@13.247.15.250:/home/aart/api1/migrations/
ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl start app1.service"


#Performance Testing AWS EC2
ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl stop app1.service; sudo rm /home/aart/api1/api_bak;  sudo mv /home/aart/api1/api /home/aart/api1/api_bak"
scp -i "aart_key_access.pem" aart_api ubuntu@13.247.15.250:/home/aart/api1/api
ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo chmod +x /home/aart/api1/api"
ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl start app1.service"



#
ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl stop app2.service; sudo rm /home/aart/api2/api_bak; sudo mv /home/aart/api2/api /home/aart/api2/api_bak; sudo cp /home/aart/api1/api /home/aart/api2/api"
ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo chmod +x /home/aart/api2/api"
ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl start app2.service"

#
ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl stop app3.service; sudo rm /home/aart/api3/api_bak; sudo mv /home/aart/api3/api /home/aart/api3/api_bak; sudo cp /home/aart/api1/api /home/aart/api3/api"
ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo chmod +x /home/aart/api3/api"
ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl start app3.service"




rm aart_api
