MySQL, version : 8.0.27
Go, version : 1.21.1 darwin/amd64

For .env we use : joho/godotenv v1.5.1

TO create fake DATA : https://www.mockaroo.com


TO have access to the localfiles :
mysql -u root -h 127.0.0.1 --protocol=tcp -p
SHOW GLOBAL VARIABLES LIKE 'local_infile';
SET GLOBAL local_infile = 'ON';

ProgressBar :
https://github.com/schollz/progressbar

To compile and run .go
go run main.go

To compile :
go build main.go