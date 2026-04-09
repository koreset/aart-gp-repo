#!/usr/bin/env bash

echo "Building binary for linux deploy"
env GOOS=linux GOARCH=amd64 go build -o aart_api

# Create migrations directory if it doesn't exist
echo "Ensuring migrations directory exists"
mkdir -p migrations

#Staging
#ssh -l root app.aart-enterprise.com "systemctl stop app1.service; rm /home/aart/api1/api"
#scp aart_api root@app.aart-enterprise.com:/home/aart/api1/api
#ssh -l root app.aart-enterprise.com "systemctl start app1.service"

#ssh -l root app.aart-enterprise.com "systemctl stop app2.service; rm /home/aart/api2/api"
#scp aart_api root@app.aart-enterprise.com:/home/aart/api2/api
#ssh -l root app.aart-enterprise.com "systemctl start app2.service"

#Performance Testing
#ssh -l root 64.226.92.196 "systemctl stop app1.service; rm /home/aart/api1/api"
#scp aart_api root@64.226.92.196:/home/aart/api1/api
#ssh -l root 64.226.92.196 "systemctl start app1.service"
#
#ssh -l root 64.226.92.196 "systemctl stop app2.service; rm /home/aart/api2/api"
#scp aart_api root@64.226.92.196:/home/aart/api2/api
#ssh -l root 64.226.92.196 "systemctl start app2.service"


# Copy migrations directory
#echo "Copying migrations to server"
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo mkdir -p /home/aart/api1/migrations"
#scp -i "aart_key_access.pem" -r migrations/* ubuntu@13.247.15.250:/home/aart/api1/migrations/
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl start app1.service"


#Performance Testing AWS EC2
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl stop app1.service; sudo rm /home/aart/api1/api_bak;  sudo mv /home/aart/api1/api /home/aart/api1/api_bak"
#scp -i "aart_key_access.pem" aart_api ubuntu@13.247.15.250:/home/aart/api1/api
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo chmod +x /home/aart/api1/api"



#
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl stop app2.service; sudo rm /home/aart/api2/api_bak; sudo mv /home/aart/api2/api /home/aart/api2/api_bak; sudo cp /home/aart/api1/api /home/aart/api2/api"
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo chmod +x /home/aart/api2/api"

#
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl stop app3.service; sudo rm /home/aart/api3/api_bak; sudo mv /home/aart/api3/api /home/aart/api3/api_bak; sudo cp /home/aart/api1/api /home/aart/api3/api"
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo chmod +x /home/aart/api3/api"


#echo "Copying migrations to server"
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo mkdir -p /home/aart/api3/migrations"
#scp -i "aart_key_access.pem" -r migrations/* ubuntu@13.247.15.250:/home/aart/api3/migrations/
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl start app3.service"
#
##echo "Copying migrations to server"
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo mkdir -p /home/aart/api2/migrations"
#scp -i "aart_key_access.pem" -r migrations/* ubuntu@13.247.15.250:/home/aart/api2/migrations/
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl start app2.service"
#

# Copy migrations directory to second instance
#echo "Copying migrations to second instance"
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo mkdir -p /home/aart/api2/migrations; sudo cp -r /home/aart/api1/migrations/* /home/aart/api2/migrations/"
#ssh -i "aart_key_access.pem" -l ubuntu 13.247.15.250 "sudo systemctl start app2.service"

#AWS EC2 Zambia
#ssh -i "aart_key_access.pem" -l ubuntu alliance-api.aart-enterprise.com "sudo systemctl stop app1.service; sudo rm /home/aart/api1/api"
#scp -i "aart_key_access.pem" aart_api ubuntu@alliance-api.aart-enterprise.com:/home/aart/api1/api
# Copy migrations directory to Zambia server
#echo "Copying migrations to Zambia server"
#ssh -i "aart_key_access.pem" -l ubuntu alliance-api.aart-enterprise.com "sudo mkdir -p /home/aart/api1/migrations"
#scp -i "aart_key_access.pem" -r migrations/* ubuntu@alliance-api.aart-enterprise.com:/home/aart/api1/migrations/
#ssh -i "aart_key_access.pem" -l ubuntu alliance-api.aart-enterprise.com "sudo systemctl start app1.service"

#ssh -i "aart_key_access.pem" -l ubuntu alliance-api.aart-enterprise.com "sudo systemctl stop app2.service; sudo rm /home/aart/api2/api; sudo cp /home/aart/api1/api /home/aart/api2/api"
# Copy migrations directory to second instance on Zambia server
#echo "Copying migrations to second instance on Zambia server"
#ssh -i "aart_key_access.pem" -l ubuntu alliance-api.aart-enterprise.com "sudo mkdir -p /home/aart/api2/migrations; sudo cp -r /home/aart/api1/migrations/* /home/aart/api2/migrations/"
#ssh -i "aart_key_access.pem" -l ubuntu alliance-api.aart-enterprise.com "sudo systemctl start app2.service"

#ssh -i "aart_key_access.pem" -l ubuntu alliance-api.aart-enterprise.com "sudo systemctl stop app3.service; sudo rm /home/aart/api3/api; sudo cp /home/aart/api1/api /home/aart/api3/api"
# Copy migrations directory to third instance on Zambia server
#echo "Copying migrations to third instance on Zambia server"
#ssh -i "aart_key_access.pem" -l ubuntu alliance-api.aart-enterprise.com "sudo mkdir -p /home/aart/api3/migrations; sudo cp -r /home/aart/api1/migrations/* /home/aart/api3/migrations/"
#ssh -i "aart_key_access.pem" -l ubuntu alliance-api.aart-enterprise.com "sudo systemctl start app3.service"



#Afrihost Copy migration files
scp -r migrations/* aartadmin@aartserver.dedicated.co.za:/opt/aart/api1/migrations/
#scp -r migrations/* aartadmin@aartserver.dedicated.co.za:/opt/aart/api2/migrations/
#scp -r migrations/* aartadmin@aartserver.dedicated.co.za:/opt/aart/api3/migrations/

#Afrihost
ssh -l aartadmin aartserver.dedicated.co.za "sudo systemctl stop app1.service; sudo systemctl stop app2.service; sudo systemctl stop app3.service; sudo rm /opt/aart/api1/aart-api; rm /opt/aart/api2/aart-api; rm /opt/aart/api3/aart-api"
scp aart_api aartadmin@aartserver.dedicated.co.za:/opt/aart/api1/aart-api
ssh -l aartadmin aartserver.dedicated.co.za "cp /opt/aart/api1/aart-api /opt/aart/api2/aart-api; cp /opt/aart/api1/aart-api /opt/aart/api3/aart-api; "
ssh -l aartadmin aartserver.dedicated.co.za "sudo chmod +x /opt/aart/api1/aart-api; sudo chmod +x /opt/aart/api2/aart-api; sudo chmod +x /opt/aart/api3/aart-api"
ssh -l aartadmin aartserver.dedicated.co.za "sudo systemctl start app1.service"
ssh -l aartadmin aartserver.dedicated.co.za "sudo systemctl start app1.service; sudo systemctl start app2.service; sudo systemctl start app3.service"
#ssh -l aartadmin aartserver.dedicated.co.za "sudo systemctl start app1.service"


#TMU Copy migration files
#scp -r migrations/* azureadmin@4.221.211.46:/opt/aart/apps/api1/migrations/
#scp -r migrations/* azureadmin@4.221.211.46:/opt/aart/apps/api2/migrations/
#scp -r migrations/* azureadmin@4.221.211.46:/opt/aart/apps/api1/migrations/

#TMU
#ssh -l azureadmin 4.221.211.46 "sudo systemctl stop app1.service; sudo systemctl stop app2.service; sudo systemctl stop app3.service; sudo rm /opt/aart/apps/api1/aart-api; rm /opt/aart/apps/api2/aart-api; rm /opt/aart/apps/api3/aart-api"
#scp aart_api azureadmin@4.221.211.46:/opt/aart/apps/api1/aart-api
#ssh -l azureadmin 4.221.211.46 "sudo cp /opt/aart/apps/api1/aart-api /opt/aart/apps/api2/aart-api; sudo cp /opt/aart/apps/api1/aart-api /opt/aart/apps/api3/aart-api; "
#ssh -l azureadmin 4.221.211.46 "sudo chmod +x /opt/aart/apps/api1/aart-api; sudo chmod +x /opt/aart/apps/api1/aart-api; sudo chmod +x /opt/aart/apps/api1/aart-api"
#ssh -l azureadmin 4.221.211.46 "sudo systemctl start app1.service; sudo systemctl start app2.service; sudo systemctl start app3.service"




#Afrihost 2
#ssh -l root 197.242.147.13 "systemctl stop aart_app.service; rm /home/aart/api1/api"
#scp aart_api root@197.242.147.13:/home/aart/api1/api
#ssh -l root 197.242.147.13 "cp /home/aart/api1/api /home/aart/api2/api; systemctl start aart_app.service; systemctl start aart2_app.service"
#




rm aart_api

#[Unit]
#Description=AART APP1
#After=network.target mysql.service
#
#[Service]
#Type=simple
#Restart=always
#RestartSec=5s
#WorkingDirectory=/home/aart/api1
#ExecStart=/home/aart/api1/api
#Environment=PORT=9090
#
#[Install]
#WantedBy=multi-user.target


#frontend haproxy-main
#    bind *:80
#    option forwardfor
#    default_backend go_servers
#
#backend go_servers
#    balance roundrobin
#    server app1		localhost:9090 check
#    server app2         localhost:9091 check
