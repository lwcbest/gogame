go build server.go
cd deploy
pwd
cd local
pwd
pm2 restart 127.0.0.1.json
pm2 logs