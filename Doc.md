MySQL, version : 8.0.27
Go, version : 1.21.1 darwin/amd64

For .env we use : joho/godotenv v1.5.1



We admit that the user is already in MySQL, so we create the user before creating the go program and we don t put this step in the process, its here just for information.

CREATE USER 'new_user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON * . * TO 'new_user'@'localhost';
FLUSH PRIVILEGES;