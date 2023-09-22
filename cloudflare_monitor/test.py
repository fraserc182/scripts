import json

json_data = '''
[
  {
    "id": "jh8785777qww",
    "name": "Cloudflare Pages — Functions live logs issues",
    "status": "investigating",
    "created_at": "2023-07-04T14:35:20.010Z",
    "updated_at": "2023-07-04T14:36:26.924Z",
    "monitoring_at": null,
    "resolved_at": null,
    "impact": "minor",
    "shortlink": "https://stspg.io/ym4z3nth12r2",
    "started_at": "2023-07-04T14:35:20.000Z",
    "page_id": "yh6f0r4529hb",
    "incident_updates": [
      {
        "id": "ws508wv4gzsg",
        "status": "investigating",
        "body": "Cloudflare is aware of and investigating an issue with streaming Functions logs on Cloudflare Pages deployments. Traffic to those deployments is not affected.",
        "incident_id": "jh8785777qww",
        "created_at": "2023-07-04T14:35:20.111Z",
        "updated_at": "2023-07-04T14:36:26.861Z",
        "display_at": "2023-07-04T14:35:20.000Z",
        "affected_components": [
          {
            "code": "vgxj684rcw7t",
            "name": "Cloudflare Sites and Services - Pages",
            "old_status": "operational",
            "new_status": "degraded_performance"
          }
        ],
        "deliver_notifications": true,
        "custom_tweet": null,
        "tweet_id": null
      }
    ],
    "components": [
      {
        "id": "vgxj684rcw7t",
        "name": "Pages",
        "status": "degraded_performance",
        "created_at": "2021-04-16T21:04:08.018Z",
        "updated_at": "2023-07-04T14:35:20.051Z",
        "position": 44,
        "description": null,
        "showcase": false,
        "start_date": "2021-04-16",
        "group_id": "1km35smx8p41",
        "page_id": "yh6f0r4529hb",
        "group": false,
        "only_show_if_degraded": false
      }
    ]
  },
  {
    "id": "jh8785777qww",
    "name": "test",
    "status": "test",
    "created_at": "1234",
    "updated_at": "12345",
    "monitoring_at": null,
    "resolved_at": null,
    "impact": "big boy",
    "shortlink": "https://stspg.io/sdfasdf",
    "started_at": "2023-07-04T14:35:20.000Z",
    "page_id": "yh6f0r4529hb",
    "incident_updates": [
      {
        "id": "ws508wv4gzsg",
        "status": "investigating",
        "body": "Cloudflare is aware of and investigating an issue with streaming Functions logs on Cloudflare Pages deployments. Traffic to those deployments is not affected.",
        "incident_id": "jh8785777qww",
        "created_at": "2023-07-04T14:35:20.111Z",
        "updated_at": "2023-07-04T14:36:26.861Z",
        "display_at": "2023-07-04T14:35:20.000Z",
        "affected_components": [
          {
            "code": "vgxj684rcw7t",
            "name": "Cloudflare Sites and Services - Pages",
            "old_status": "operational",
            "new_status": "degraded_performance"
          }
        ],
        "deliver_notifications": true,
        "custom_tweet": null,
        "tweet_id": null
      }
    ],
    "components": [
      {
        "id": "vgxj684rcw7t",
        "name": "Pages",
        "status": "degraded_performance",
        "created_at": "2021-04-16T21:04:08.018Z",
        "updated_at": "2023-07-04T14:35:20.051Z",
        "position": 44,
        "description": null,
        "showcase": false,
        "start_date": "2021-04-16",
        "group_id": "1km35smx8p41",
        "page_id": "yh6f0r4529hb",
        "group": false,
        "only_show_if_degraded": false
      }
    ]
  }  
]

'''

# Parse the JSON data
data = json.loads(json_data)

# Extract the fields from each dictionary in the list
extracted_data = []
for item in data:
    extracted_item = {
        'name': item['name'],
        'status': item['status'],
        'created_at': item['created_at'],
        'updated_at': item['updated_at'],
        'shortlink': item['shortlink']
    }
    extracted_data.append(extracted_item)

# Print the extracted data
for item in extracted_data:
    print(item)
