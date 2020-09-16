## Data base information

MongoDB information such as db url, data base name and collection name are mentioned in the .env file

## REST APIs

### GET /enrollees

It returns all the enrollee records from db.

### POST /enrollee

It saves an enrollee record to db along with its dependent information.

Sample request:

{
	"name":"John",
	"is_active": true,
	"phone_number": "1234567890",
	"dependents":[
					{
						"name" : "Jack",
						"birth_date" : "1990-01-01"
					},
					{
						"name" : "Jim",
						"birth_date" : "1992-01-05"
					}
				]
}

### GET /enrollee/:enrolleeId

It returns an enrollee record along with its dependents from db for the provided enrolleeId.

