from flask import Flask, request, jsonify
from flask_cors import CORS
from PIL import Image
import threading
import uuid
from datetime import datetime
from instabot import Bot
import os
import shutil
from dotenv import load_dotenv


app = Flask(__name__)
load_dotenv('./.env')
IGusername = os.getenv("IGUSERNAME")
IGpassword = os.getenv("IGPASSWORD")

# Enable CORS for your app
CORS(app)

def upload_to_instagram(unique_filename):
    try:
        bot = Bot()
        bot.login(username=IGusername, password=IGpassword, is_threaded=True)
        image_path = os.path.join('uploads', unique_filename)  # Construct the full image path
        bot.upload_photo(image_path, caption="Tag Him/Her in comments. ! Love from a Senior Automated the Page !")
    except Exception as e:
        print(e)
        pass


def clean_up():
    dir = "config"
    if os.path.exists(dir):
        try:
            shutil.rmtree(dir)
        except OSError as e:
            print(e)
    
    directory = "uploads"
    if os.path.exists(directory):
        for filename in os.listdir(directory):
            if filename.endswith(".REMOVE_ME"):
                file_path = os.path.join(directory, filename)
                try:
                    os.remove(file_path)
                    print(f"Deleted {filename}")
                except OSError as e:
                    print(f"Error deleting {filename}: {e}")

def generate_unique_filename():
    # Generate a unique filename using a combination of timestamp and a random UUID
    timestamp = datetime.now().strftime('%Y%m%d%H%M%S')
    unique_id = str(uuid.uuid4().hex)[:8]  # Use the first 8 characters of the UUID
    return f'screenshot_{timestamp}_{unique_id}.jpg'


@app.route('/upload', methods=['POST'])
def upload_screenshot():
    try:
        screenshot_data = request.files['screenshot']
        if screenshot_data:
            image = Image.open(screenshot_data)

            if image.mode == 'RGBA':
                image = image.convert('RGB')
            aspect_ratio = (1080, 820)
            image.thumbnail(aspect_ratio)
            
            unique_filename = generate_unique_filename()
            image.save(os.path.join('uploads', unique_filename))

            clean_up()
            instagram_thread = threading.Thread(target=upload_to_instagram, args=(unique_filename,))
            instagram_thread.start()

            return jsonify({'message': 'Screenshot uploaded successfully'}), 200
        else:
            return jsonify({'error': 'No screenshot data received'}), 400
    except Exception as e:
        return jsonify({'error': str(e)}), 500

