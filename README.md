# Rest_Api_Example
Simple rest api example written in GO Lang - No database required for testing

 This has several kinds of API Actions: 
      list all people (GET /people), 
      Display s person via ID (GET /people/{id})),
      Delete a person via ID (DELETE /people/{id})), 
      Create a person record via ID (POST /people/{id}))

 Dependencies: Gorilla/MUX
  go get github.com/gorilla/mux

