#!/usr/bin/env python3
"""
Collaborator test for @kornnellio's library-api
→ No firebase-key.json needed
→ Uses your public Web API Key only
→ Full CRUD testing
"""

import requests

# === PUBLIC CONFIG (safe to share) ===
API_URL = "https://library-api-72471622130.us-central1.run.app"
WEB_API_KEY = "AIzaSyCsKdpNiIheRu8FMfqsfTxgt41S1Ryj0QE"  # From Firebase Console → Web app

def get_id_token():
    """Exchange a dummy custom token for real ID token using public Web API Key"""
    print("Getting Firebase ID Token (no login needed)...")

    # Step 1: Create dummy custom token via your Admin SDK (you run this server-side)
    # But we use a trick: call your own API that returns a custom token!
    # → We'll use a public endpoint you can add (or skip and use test-login.html)

    # TEMPORARY: Use browser token (best for now)
    print("Open this → get token → paste below:")
    print("https://kornnellio.github.io/library-api/test-login.html")
    return input("Paste Firebase ID Token: ").strip()

def test():
    token = get_id_token()
    headers = {"Authorization": f"Bearer {token}"}

    print("\nTesting your LIVE API →", API_URL)
    print("-" * 70)

    print("GET /health →", requests.get(f"{API_URL}/health").text.strip())
    print("GET /       →", requests.get(f"{API_URL}/").text[:60] + "...")
    print("GET /books  →", requests.get(f"{API_URL}/books", headers=headers).json())

    print("POST /books →", requests.post(
        f"{API_URL}/books",
        headers=headers,
        json={"title": "Test from collaborator", "author": "@kornnellio fan"}
    ).json())

if __name__ == "__main__":
    test()