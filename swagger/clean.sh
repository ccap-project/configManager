perl -pi.bak -e s!/prodx!!g cmd/prodx-server/main.go restapi/operations/prodx_api.go
perl -pi.bak -e s!prodx!..!g restapi/operations/*/*.go
