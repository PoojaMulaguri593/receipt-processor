Receipt Processor API

Overview:
This project is a simple API built with Go that processes receipts and calculates reward points based on a set of predefined rules.

Requirements:
- Go 1.18 or later
- `github.com/google/uuid` package

Installation:
1. Clone this repository:
   ```bash
   git clone https://github.com/PoojaMulaguri593/receipt-processor.git
   cd receipt-processor
   ```

2. Install dependencies (if any):
   ```bash
   go mod tidy
   ```

3. Run the API server:
   ```bash
   go run main.go
   ```

Run Instructions:
- Ensure Go is installed on your machine.
- Open a terminal and navigate to the project directory.
- Use the command `go run main.go` to start the server.
- The server will run on `http://localhost:8080`.
- Use `cURL` or Postman to send requests.

API Endpoints:
1. **Process a Receipt**
   - **Endpoint:** `POST /receipts/process`
   - **Request Body (JSON):**
     ```json
     {
       "retailer": "Target",
       "purchaseDate": "2022-01-01",
       "purchaseTime": "13:01",
       "items": [
         { "shortDescription": "Mountain Dew 12PK", "price": "6.49" }
       ],
       "total": "35.35"
     }
     ```
   - **Response:**
     ```json
     { "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
     ```

2. **Get Points for a Receipt**
   - **Endpoint:** `GET /receipts/{id}/points`
   - **Response:**
     ```json
     { "points": 32 }
     ```

Testing:
Use cURL or Postman to send requests and check responses.
