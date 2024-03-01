# kidgo
The chatting app for very young kids (kindergarten or below), prioritizing safety, simplicity, and ease of use.

## Code Structure
The `messages` folder contains the `messages` package as well as the code that actually runs on the Lambda. The `messages` package primarily allows you to access all messages that have been sent to a user and also allows you to access messages they've sent. The `lambda` subfolder contains the code that runs on the lambda, like a handler and picking different API functions depending on the request.

The `frontend` folder contains the code that runs on the frontend. There's my Python web server, HTML, CSS, and JavaScript, as well as the navigator helpers' audio files.

## Running
First, upload the `bootstrap.zip` file to Lambda, then connect the Lambda functions to API Gateway. To run the client, you'll need to make a stage deployment of the KidGo API. Get the stage deployment's invoke url and set the environment variable `API_ENDPOINT` to that. Add your AWS credentials to the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables. You'll also need to make the DynamoDB databases and S3 buckets needed by the API. Once setup, run `main.py` to start the webserver. 