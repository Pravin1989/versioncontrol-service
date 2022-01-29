# versioncontrol-service

This repository contains the utility methods to create branch, make changes in existing files, create commit, push and 
create pull request.
## Clone the project

```
$ git clone git@github.com:Pravin1989/versioncontrol-service.git

```

## Two ways to run it 

```
1. If you have a Docker installed on your machine then run below commands
$ cd versioncontrol-service
$ docker compose build
$ docker compose up
```
```
2. If you have installed only Go on your machine then run below commands but before that you need to set the environment variables which are present in `env.example` file
$ cd versioncontrol-service
$ go build .\src\
$ run the .exe file
```

## API Details
* Endpoint : http://localhost:8080/versioncontrol-service/
* Http Method : GET

Once you hit the above URL you will be redirected to github login page if you are not already logged in then you need to login and you can see in the console the PR link has been created with some dummy changes in `sample-repo`.