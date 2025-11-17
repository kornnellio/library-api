# test_api.py — Full test suite for library-api
# Run: python test_api.py

import os
import requests
import firebase_admin
from firebase_admin import auth, credentials

# === CONFIG ===
API_URL = "https://library-api-72471622130.us-central1.run.app"
FIREBASE_KEY_PATH = "firebase-key.json"  # Your Admin SDK key (local only!)
TEST_UID = "test-collab-123"             # Any string — no real user needed

# === Initialize Firebase Admin ===
if not firebase_admin._apps:
    cred = credentials.Certificate(FIREBASE_KEY_PATH)
    firebase_admin.initialize_app(cred)

def get_token():
    """Generate Firebase Custom Token (no Google login!)"""
    token = auth.create_custom_token(TEST_UID)
    return token.decode('utf-8')

def api_request(method, endpoint, token=None, data=None):
    url = f"{API_URL}{endpoint}"
    headers = {}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    if data is not None:
        headers["Content-Type"] = "application/json"

    response = requests.request(method, url, headers=headers, json=data)
    print(f"{method} {endpoint} → {response.status_code}")
    try:
        print(response.json())
    except:
        print(response.text)
    print("-" * 50)
    return response

def main():
    print("library-api Full Test Suite")
    print(f"Target: {API_URL}\n")

    # Generate token
    token = get_token()
    print(f"Token generated for UID: {TEST_UID}\n")

    # Public endpoints
    api_request("GET", "/health")
    api_request("GET", "/")

    # Protected CRUD
    api_request("GET", "/books", token)

    # Create book
    create_resp = api_request("POST", "/books", token, {
        "title": "Test from @kornnellio",
        "author": "Python Test Suite"
    })

    if create_resp.status_code == 201:
        book_id = create_resp.json().get("id")
        print(f"Book created with ID: {book_id}")

        # Read again
        api_request("GET", "/books", token)

        # Update
        api_request("PUT", f"/books/{book_id}", token, {
            "title": "Updated by Python",
            "author": "kornnellio"
        })

        # Delete
        api_request("DELETE", f"/books/{book_id}", token)

        # Final list
        api_request("GET", "/books", token)
    else:
        print("Failed to create book — check API")

if __name__ == "__main__":
    if not os.path.exists(FIREBASE_KEY_PATH):
        print(f"ERROR: {FIREBASE_KEY_PATH} not found!")
        print("Download from Firebase Console → Project Settings → Service Accounts")
    else:
        main()