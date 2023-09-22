import json

msg = {
        "incidents": [
        {
            "id": "testid12345",
            "name": "Cloudflare Pages — Functions live logs issues",
            "status": "investigating",
            "created_at": "2023-07-04T14: 35: 20.010Z",
            "updated_at": "2023-07-04T14: 36: 26.924Z",
            "shortlink": "https: //stspg.io/ym4z3nth12r2"
        }
        ]
    }

print(msg['incidents'])