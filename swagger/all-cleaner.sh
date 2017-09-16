perl -pi.bak -e s!/client!!g client/cmd/*/main.go client/restapi/operations/config_manager_api.go
perl -pi.bak -e s!client!..!g client/restapi/operations/*/*.go
