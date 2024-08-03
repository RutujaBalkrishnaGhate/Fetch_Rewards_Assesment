# Fetch_Rewards_Assessment
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
</head>
<body>

<h1>Receipt Processor</h1>

<h2>Introduction</h2>
<p>This project is a receipt processor web service built with Go. The service fulfills the documented API for processing receipts and calculating points based on specific rules. The API is defined in the <code>api.yml</code> file and described in the documentation below.</p>

<h2>Requirements</h2>
<ul>
    <li>Go (version 1.20 or later)</li>
    <li>Docker</li>
</ul>

<h2>Getting Started</h2>

<h3>Running the Application</h3>

<h4>Using Go</h4>
<ol>
    <li>Clone the repository:
        <pre><code>git clone git@github.com:RutujaBalkrishnaGhate/Fetch_Rewards_Assesment.git
</code></pre>
    </li>
    <li>Build and run the application:
        <pre><code>go run .</code></pre>
    </li>
    <li>The application will start on port 8080 by default.</li>
</ol>

<h4>Using Docker</h4>
<ol>
    <li>Clone the repository:
        <pre><code>git clone https://github.com/your-username/receipt-processor.git
cd receipt-processor</code></pre>
    </li>
    <li>Build the Docker image:
        <pre><code>docker build -t receipt-processor .</code></pre>
    </li>
    <li>Run the Docker container:
        <pre><code>docker run -p 8080:8080 receipt-processor</code></pre>
    </li>
    <li>The application will be accessible on port 8080.</li>
</ol>

<h2>API Endpoints</h2>

<h3>Process Receipts</h3>
<ul>
    <li><strong>Endpoint:</strong> <code>/receipts/process</code></li>
    <li><strong>Method:</strong> <code>POST</code></li>
    <li><strong>Payload:</strong> JSON object representing a receipt</li>
    <li><strong>Response:</strong> JSON object containing an <code>id</code> for the receipt</li>
</ul>

<h4>Example Request</h4>
<pre><code>{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },
    {
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    }
  ],
  "total": "18.74"
}</code></pre>

<h4>Example Response</h4>
<pre><code>{
  "id": "7fb1377b-b223-49d9-a31a-5a02701dd310"
}</code></pre>

<h3>Get Points</h3>
<ul>
    <li><strong>Endpoint:</strong> <code>/receipts/{id}/points</code></li>
    <li><strong>Method:</strong> <code>GET</code></li>
    <li><strong>Response:</strong> JSON object containing the number of points awarded</li>
</ul>

<h4>Example Response</h4>
<pre><code>{
  "points": 28
}</code></pre>

<h2>Points Calculation Rules</h2>
<ol>
    <li>One point for every alphanumeric character in the retailer name.</li>
    <li>50 points if the total is a round dollar amount with no cents.</li>
    <li>25 points if the total is a multiple of 0.25.</li>
    <li>5 points for every two items on the receipt.</li>
    <li>If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.</li>
    <li>6 points if the day in the purchase date is odd.</li>
    <li>10 points if the time of purchase is after 2:00pm and before 4:00pm.</li>
</ol>

<h2>Examples</h2>
<h4>Example 1</h4>
<pre><code>{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}</code></pre>
<p><strong>Total Points:</strong> 28</p>
<p><strong>Breakdown:</strong></p>
<ul>
    <li>6 points - retailer name has 6 characters</li>
    <li>10 points - 4 items (2 pairs @ 5 points each)</li>
    <li>3 points - "Emils Cheese Pizza" is 18 characters (a multiple of 3)
        <ul>
            <li>item price of 12.25 * 0.2 = 2.45, rounded up is 3 points</li>
        </ul>
    </li>
    <li>3 points - "Klarbrunn 12-PK 12 FL OZ" is 24 characters (a multiple of 3)
        <ul>
            <li>item price of 12.00 * 0.2 = 2.4, rounded up is 3 points</li>
        </ul>
    </li>
    <li>6 points - purchase day is odd</li>
</ul>

<h4>Example 2</h4>
<pre><code>{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}</code></pre>
<p><strong>Total Points:</strong> 109</p>
<p><strong>Breakdown:</strong></p>
<ul>
    <li>50 points - total is a round dollar amount</li>
    <li>25 points - total is a multiple of 0.25</li>
    <li>14 points - retailer name (M&M Corner Market) has 14 alphanumeric characters
        <ul>
            <li>note: '&' is not alphanumeric</li>
        </ul>
    </li>
    <li>10 points - 2:33pm is between 2:00pm and 4:00pm</li>
    <li>10 points - 4 items (2 pairs @ 5 points each)</li>
</ul>

<h2>Running Tests</h2>
<p>To run the tests for the application, use the following command:</p>
<pre><code>go test ./...</code></pre>

<h2>Directory Structure</h2>
<pre><code>.
├── main.go          # Entry point of the application
├── router.go        # Router setup and route handlers
├── handlers.go      # Receipt processing logic
├── main_test.go     # Tests for the application
├── Dockerfile       # Dockerfile for containerization
└── README.md        # This readme file
</code></pre>

<h2>License</h2>
<p>This project is licensed under the MIT License. See the <a href="LICENSE">LICENSE</a> file for details.</p>

</body>
</html>

