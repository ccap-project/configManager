# Dependencies sync/download
download govendor (go get -u github.com/kardianos/govendor) and run "govendor sync -v"

# Dynamic files generation
download go-swagger version 0.13.0 binary (https://github.com/go-swagger/go-swagger/releases/tag/0.13.0)
and run "swagger generate server -P models.Customer --skip-validation -f swagger/swagger.yml"
