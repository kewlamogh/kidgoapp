from flask import Flask, render_template, request
import requests
import string, random, boto3, os
from urllib.parse import quote

app = Flask(__name__)

# AWS S3 configuration
AWS_ACCESS_KEY_ID = os.environ.get('AWS_ACCESS_KEY_ID')
AWS_SECRET_ACCESS_KEY = os.environ.get('AWS_SECRET_ACCESS_KEY')
S3_BUCKET_NAME = 'kidgomsgs'
API_ENDPOINT = os.environ.get('API_ENDPOINT')   

@app.route('/')
def home():
    return render_template('index.html')

@app.route('/inbox', methods=['GET'])
def inbox():    
    response = requests.get(
     API_ENDPOINT + '/inbox?username=' + request.args.get('user').lower() + "&classroom=" + request.args.get("class").lower())
    
    print(response.status_code)

    # If the request was successful, response.status_code will be 200
    if response.status_code == 200:
        print(response.json())
        return response.json()  # Print the json data
    else:
        return {}
    


def generate_random_string(length):
    # 'abcdefghijklmnopqrstuvwxyz0123456789'
    characters = string.ascii_lowercase + string.digits
    return ''.join(random.choice(characters) for _ in range(length))

@app.route("/sendmessage", methods=["POST"]) 
def sendmessage():
    from_param = request.form['from'].lower()
    to_param = request.form['to'].lower()
    audio_blob = request.files['audioBlob']
    classid = request.form['class']
    print(classid)

    # Save the audio Blob to a temporary file
    temp_file_path = 'temp_audio.wav'
    audio_blob.save(temp_file_path)

    # Upload the temporary file to Amazon S3
    s3 = boto3.client('s3', aws_access_key_id=AWS_ACCESS_KEY_ID,
                      aws_secret_access_key=AWS_SECRET_ACCESS_KEY)

    s3_url = f"{from_param}_{to_param}_{generate_random_string(10)}.wav'"
    s3.upload_file(temp_file_path, S3_BUCKET_NAME, s3_url)

    s3_url = "https://kidgomsgs.s3.ap-south-1.amazonaws.com/" + s3_url

    # Clean up the temporary file
    os.remove(temp_file_path)
    response = requests.get(
         API_ENDPOINT + '/send?from=' + from_param + "&to=" + to_param + "&s3url=" + quote(s3_url) + "&classroom=" + classid.lower())

    return 'Audio Blob uploaded to S3 successfully'

@app.route('/outbox', methods=['GET'])
def outbox():    
    response = requests.get(
         API_ENDPOINT + '/outbox?username=' + request.args.get('user').lower() + "&classroom=" + request.args.get("class").lower())

    print('https://xujhofn50c.execute-api.ap-south-1.amazonaws.com/production/outbox?username=' +
          request.args.get('user').lower())
    
    print(response.status_code)

    # If the request was successful, response.status_code will be 200
    if response.status_code == 200:
        print(response.json())
        return response.json()  # Print the json data
    else:
        return {}

if __name__ == '__main__':
    app.run(debug=True)
