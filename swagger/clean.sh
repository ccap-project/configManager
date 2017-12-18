perl -pi.bak -e s!/configManager!!g cmd/config-manager-server/main.go restapi/operations/config_manager_api.go
perl -pi.bak -e s!configManager!..!g restapi/operations/*/*.go
perl -pi.bak -e s!../configManager/!!g client/*.go
perl -pi.bak -e s!/configManager!!g client/*/*.go
