# Generating code from swagger
swagger generate server -P models.Customer --skip-validation -f swagger/swagger.yml &&  bash swagger/clean.sh 
