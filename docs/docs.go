// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/feeds": {
            "get": {
                "description": "Retrieve all feeds from the store",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "feeds"
                ],
                "summary": "Get all feeds",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/feed.Presenter"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new feed with the data provided in the request",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "feeds"
                ],
                "summary": "Create a new feed",
                "parameters": [
                    {
                        "description": "CreateRequest data",
                        "name": "feed",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/feed.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/feed.Presenter"
                        }
                    }
                }
            }
        },
        "/feeds/{id}": {
            "delete": {
                "description": "Delete a feed with the given feed ID",
                "tags": [
                    "feeds"
                ],
                "summary": "Delete a feed",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Feed deleted successfully"
                    }
                }
            }
        },
        "/feeds/{id}/sync": {
            "put": {
                "description": "Sync a feed with the given feed ID",
                "tags": [
                    "feeds"
                ],
                "summary": "Sync a feed",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Feed ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Feed synced successfully"
                    }
                }
            }
        }
    },
    "definitions": {
        "feed.CreateRequest": {
            "type": "object",
            "required": [
                "title",
                "url"
            ],
            "properties": {
                "title": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "feed.Presenter": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "syncedAt": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "PCast REST-API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
