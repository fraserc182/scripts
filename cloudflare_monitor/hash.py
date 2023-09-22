import hashlib
import json

def generate_unique_id(data):
    # Convert the dictionary to a JSON string
    json_string = json.dumps(data, sort_keys=True)

    # Generate the SHA-256 hash of the JSON string
    sha256_hash = hashlib.sha256(json_string.encode()).hexdigest()

    return sha256_hash

# Example dictionary
dictionary = {
    "name": "test",
    "status": "test",
    "created_at": "1234",
    "updated_at": "12345",
    "shortlink": "https://stspg.io/sdfasdf"
}

# Generate a unique ID based on the dictionary
unique_id = generate_unique_id(dictionary)

# append the unique ID to the dictionary
dictionary["unique_id"] = unique_id

print(dictionary)
