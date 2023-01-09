# YTSearchAPI
-----------TO RUN IN LOCAL--------------
1) Clone the repo
2) Run "go mod download"
3) Run an instance of postgres DB in Docker or local
4) Populate its credentials in APIs/Worker.go file, Line no.26
5) Go to postman
5) SEARCH API:- 
    curl --location --request POST 'localhost:7777/api/Search' \
        --header 'Content-Type: application/json' \
        --data-raw '{
            "title":"",
            "description":"",
            "exactMatch" : false
        }'  

    Note:- You can get exactMatch if you set the flag as true

6) GET API:-
    curl --location --request GET 'localhost:7777/api/GetVideos?limit=2&start=0' \
    --header 'Content-Type: application/json' \
    '

    Note:- You can give start and limit as per the pagination requirements
7) Additional features:-
    you can pass -query, -apiKey flags as per need

    Note:- The supplied API key has been revoked, please use another key.

    


-----------TO RUN WITH DOCKER[EXPERIMENTAL]------------

Note:- I am new to docker, so this might or  might not work properly
1) sudo docker build -t yt .

2) RUN postgres via docker compose

3) CREATE DATABASE VideoData 

   Inside that:- 
   CREATE TABLE ytvideos (
        vidid varchar(100) PRIMARY KEY,
        title TEXT,
        description TEXT,
        publishdate varchar(100),
        thumbnails varchar(200)
    )
4) Change the host name in ./YTSearchAPI/APIs/Worker.go to local pc IP(ipconfig)
5) sudo docker-compose up
6) You can import data from the provided database db.csv