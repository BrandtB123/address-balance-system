# ABS
An address balance system in Golang, ran on AWS with a Postgres RDS backend that generates a new report every day at midnight.
Domain: http://ec2-3-141-35-217.us-east-2.compute.amazonaws.com/

# APIs

## Endpoints
### Generate Report
URL: /api/generateReport<br>
Example: http://ec2-3-141-35-217.us-east-2.compute.amazonaws.com/api/generateReport
Method: GET<br>
Description: Generates a balance report and returns the latest.<br>
Response:<br>
Status Code 200: Returns the balance of all addresses in a JSON.<br>
Status Code 400: Indicates a missing or invalid query parameter.<br>
### Get Balance of Asset
URL: /api/getReport<br>
Example: http://ec2-3-141-35-217.us-east-2.compute.amazonaws.com/api/getReport?chain=near<br>
Method: GET<br>
Description: Retieves the latest balance report for a certain asset<br>
Query Parameters: chain (required): The asset (chain was a mistake)<br>
Response:<br>
Status Code 200: Returns the latest balance report for that asset.<br>
Status Code 400: Indicates a missing or invalid query parameter.<br>
### Get Balance of Address
URL: /api/getBalance<br>
Example: http://ec2-3-141-35-217.us-east-2.compute.amazonaws.com/api/getBalance?address=demid1.near<br>
Method: GET<br>
Description: Retrieves the balance of specifc address.<br>
Query Parameters: address (required): The address to retrieve.<br>
Response:<br>
Status Code 200: Returns the address' balance.<br>
Status Code 400: Indicates a missing or invalid query parameter.<br>
