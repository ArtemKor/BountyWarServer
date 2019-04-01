# README #

This README would normally document whatever steps are necessary to get your application up and running.

### What is this repository for? ###

* Quick summary
* Version
* [Learn Markdown](https://bitbucket.org/tutorials/markdowndemo)

### How do I get set up? ###

* Summary of set up
* Configuration
* Dependencies
* Database configuration
* How to run tests
* Deployment instructions

### Contribution guidelines ###

* Writing tests
* Code review
* Other guidelines


* Repo owner or admin
* Other community or team contact


### Helper ###

###### Frontend ######
* Build frontend: ```cd frontendng; npm install; npm run build; cd ../;```

###### Server ######
* Added in environment: ```EXCONFIG=$GOPATH/src/ex/app.conf```
* Need ```packr``` package: ```go get -u github.com/gobuffalo/packr/packr```
* If server does not work on mac, then fix: ```ulimit -n 4096```
* Build backend for mac x64: ```GOOS=darwin GOARCH=amd64 packr build```
* Build backend for window x64: ```GOOS=windows GOARCH=amd64 packr build```
* Build backend for linux x64: ```GOOS=linux GOARCH=amd64 packr build```
* Build frontend & backend for linux : ```./build_project_linux.sh```

###### Docker: for prod ######
* Added in environment: ```EXPGDATA=/your/path/where/the/database/data/will/be/stored```
* Use: ```docker-compose up --build``` or ```./run_docker_comp.sh```

###### Docker: only DB ######
* Added in environment: ```EXPGDATA=/your/path/where/the/database/data/will/be/stored```
* Use: ```./run_only_db.sh```
