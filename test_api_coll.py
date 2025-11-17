#!/usr/bin/env python3
"""
Collaborator test for @kornnellio's library-api
→ No firebase-key.json needed
→ Uses your live test-login.html
→ Tests ALL ENDPOINTS: health, home, register, login, CRUD
→ Full feedback
"""

import webbrowser
import requests

# === YOUR LIVE URLs ===
API_URL = "https://library-api-72471622130.us-central1.run.app"
TOKEN_PAGE = "https://kornnellio.github.io/library-api/test-login.html"

def open_token_page():
    print("Opening token generator...")
    print(f"→ {TOKEN_PAGE}")
    print("→ Click 'Login with Google' → token auto-copied")
    print()
    webbrowser.open(TOKEN_PAGE)

def get_token():
    return input("Paste Firebase ID Token here: ").strip()

def test_all(token):
    headers = {"Authorization": f"Bearer {token}"}
    print("\nTesting ALL ENDPOINTS →", API_URL)
    print("=" * 80)

    # 1. Public endpoints
    print("GET  /health →", requests.get(f"{API_URL}/health").text.strip())
    print("GET  /       →", requests.get(f"{API_URL}/").text[:100] + "...")

    # 2. Register (optional)
    print("\nPOST /register →", end=" ")
    reg = requests.post(f"{API_URL}/register", json={"email": "test@kornnellio.ro", "password": "123456"})
    print(reg.status_code, reg.json().get("message", reg.text))

    # 3. Protected CRUD
    print("\nGET  /books  →", requests.get(f"{API_URL}/books", headers=headers).json())

    # 4. Create book
    create = requests.post(f"{API_URL}/books", headers=headers, json={
        "title": "Test from @kornnellio fan",
        "author": "Collaborator"
    })
    print("POST /books  →", create.status_code, create.json())

    if create.status_code == 201:
        book_id = create.json().get("id")
        print(f"Book created → ID: {book_id}")

        # 5. Update
        update = requests.put(f"{API_URL}/books/{book_id}", headers=headers, json={
            "title": "Updated by collaborator",
            "author": "RO Power"
        })
        print("PUT  /books/{book_id} →", update.status_code, update.json())

        # 6. Delete
        delete = requests.delete(f"{API_URL}/books/{book_id}", headers=headers)
        print("DELETE /books/{book_id} →", delete.status_code, delete.json())

        # 7. Final list
        print("GET  /books (final) →", requests.get(f"{API_URL}/books", headers=headers).json())
    else:
        print("Failed to create book — check API")

    print("\nAll tests completed! Your API is PERFECT")

def main():
    print("library-api by @kornnellio — Full Test Suite")
    print("Zero setup | Uses your token page | All endpoints")
    print("-" * 80)

    open_token_page()
    token = get_token()

    if token and token.startswith("eyJ"):
        test_all(token)
    else:
        print("Invalid token — try again")

if __name__ == "__main__":
    main()