package es_service

const (
	ArticleMapping = `{
  "mappings": {
    "properties": {
      "id": {
        "type": "integer"
      },
      "created_at": {
        "type": "date",
        "null_value": "null"
      },
      "updated_at": {
        "type": "date",
        "null_value": "null"
      },
      "title": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "abstract": {
        "type": "text"
      },
      "content": {
        "type": "text"
      },
      "category_id": {
        "type": "integer"
      },
      "tag_list": {
        "type": "keyword"
      },
      "cover": {
        "type": "keyword"
      },
      "user_id": {
        "type": "integer"
      },
      "look_count": {
        "type": "integer"
      },
      "digg_count": {
        "type": "integer"
      },
      "comment_count": {
        "type": "integer"
      },
      "collect_count": {
        "type": "integer"
      },
      "status": {
        "type": "integer"
      },
      "open_comment": {
        "type": "boolean"
      },
      "comments": {
        "type": "nested",
        "properties": {
          "id": {
            "type": "integer"
          },
          "created_at": {
            "type": "date",
            "null_value": "null"
          },
          "updated_at": {
            "type": "date",
            "null_value": "null"
          },
          "content": {
            "type": "text"
          },
          "user_id": {
            "type": "integer"
          },
          "article_id": {
            "type": "integer"
          },
          "parent_id": {
            "type": "integer"
          },
          "root_parent_id": {
            "type": "integer"
          },
          "digg_count": {
            "type": "integer"
          }
        }
      }
    }
  }
}`

	TextMapping = `{
  "mappings": {
    "properties": {
      "id": {
        "type": "integer"
      },
      "created_at": {
        "type": "date",
        "null_value": "null"
      },
      "updated_at": {
        "type": "date",
        "null_value": "null"
      },
      "head": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "body": {
        "type": "text"
      },
      "article_id": {
        "type": "integer"
      }
    }
  }
}`
)
