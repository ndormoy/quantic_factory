# THE QUANTIC FACTORY TEST

## Version used

**MySQL**, version : **8.0.27**
**Go**, version : **1.21.1 darwin/amd64**

**joho/godotenv v1.5.1** : for getting .env variables

TO create fake DATA : **https://www.mockaroo.com**

ProgressBar :
<s> https://github.com/schollz/progressbar </s> -> Not used anymore because of overlapping

## Before launching program

TO have access to the localfiles :
mysql -u root -h 127.0.0.1 --protocol=tcp -p
SHOW GLOBAL VARIABLES LIKE 'local_infile';
SET GLOBAL local_infile = 'ON';

## To build and run the program

To compile and run .go
go run main.go init_db.go utils.go treat.go quantile.go export.go

## Some infos about the program

- We see in the main that we have a commented code, its for testing  
the fact that if we add some data in our database that the export table will be modied.  
So first start the program without these lines, check the table, and then restart the program to see the changes.

- To see in a easy way if the export create a new table if the date change :  
We have to uncomment in export.go the ligne 16
and comment the line 13
It simulate a new day (2023 09 15), change the day if we are on the same day.

- I use bulk insert to insert data in the database, its faster than insert one by one. that respect the subject.

- The Email is not added in the struct of the export table because i think i misunderstood how channeltype has to be implemented.  
For me we don t had to fill ChannelType, so its just a number between 1 and 5 representing how the client was contacted.  
Same thing for EventType, number between 1 and 6 representing the type of event (visit, purchase...)