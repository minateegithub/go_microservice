## Data base information

MongoDB information such as db url, data base name and collection name are mentioned in the .env file

## REST APIs

### GET /enrollees

It returns all the enrollee records from db.

### POST /enrollee

It saves an enrollee record to db along with its dependent information. Unique uuid gets assigned to each enrollee record as 'id' in the code automatically.

Sample request:

{
	"name":"John",
	"is_active": true,
	"birth_date" : "1990-01-01"
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

#### Data validations
* 'name', 'is_active' and 'birth_date' fields are mandatory. 
* 'birth_date' should be in 'YYYY-MM-DD' format.
* 'is_active' is a boolean field and needs a valid boolean value.
* 'phone_number' is not mandatory but it should be a valid one.

### GET /enrollee/:enrolleeId

It returns the enrollee record along with its dependents for the provided enrolleeId from db.

### PUT /enrollee/:enrolleeId

It updates the enrollee record data in db for the provided enrolleeId.

#### Data validations
* 'name', 'is_active' and 'birth_date' fields are mandatory. 
* 'birth_date' should be in 'YYYY-MM-DD' format.
* 'is_active' is a boolean field and needs a valid boolean value.
* 'phone_number' is not mandatory but it should be a valid one.

### DELETE /enrollee/:enrolleeId

It deletes the enrollee record for the provided enrolleeId from db.

### POST /enrollee/:enrolleeId/dependents

It adds dependents to the enrollee record for the provided enrolleeId. Unique uuid gets assigned to each dependent record as 'id' in the code automatically.

Sample request:

[
{
    "name":"Mike",
    "birth_date": "1990-01-01"
},
{
    "name":"Jim",
    "birth_date": "1990-02-02"
}
]

#### Data validations
* 'name' and 'birth_date' fields are mandatory. 
* 'birth_date' should be in 'YYYY-MM-DD' format.

### GET /enrollee/:enrolleeId/dependents

It returns all the dependents of the enrollee record for the provided enrolleeId.

### DELETE /enrollee/:enrolleeId/dependents/:dependentId

It deletes the dependent record from db for the provided enrolleeId and dependentId.

### PUT /enrollee/:enrolleeId/dependents/:dependentId

It updates the dependent data in db for the provided enrolleeId and dependentId.

#### Data validations
* 'name' and 'birth_date' fields are mandatory. 
* 'birth_date' should be in 'YYYY-MM-DD' format.
