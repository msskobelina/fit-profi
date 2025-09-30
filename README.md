## API doc with swagger

- To generate doc follow next steps:

`go get github.com/go-swagger/go-swagger/cmd/swagger`

`go install github.com/go-swagger/go-swagger/cmd/swagger@latest`

`mkdir output`

`~/go/bin/swagger generate spec --scan-models -o ./swagger-doc/swagger.yaml`

- To open doc:

`~/go/bin/swagger serve ./swagger-doc/swagger.yaml`