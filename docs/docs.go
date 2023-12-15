// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Alex Bezverkhniy",
            "email": "alexander.bezverkhniy@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/process/": {
            "post": {
                "description": "Submits/Creates new process",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "process"
                ],
                "summary": "Creates new process",
                "parameters": [
                    {
                        "description": "ProcessRequest",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.ProcessDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.ProcessSubmitResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "Assign/move the process to the status",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "process"
                ],
                "summary": "Assign the process to the status",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Code of Process",
                        "name": "code",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "UUID of Process",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "ProcessStatus",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.ProcessStatusDTO"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/api/v1/process/{code}/list": {
            "get": {
                "description": "Get list of processes",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "process"
                ],
                "summary": "Get list of processes",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Code of Process",
                        "name": "code",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "X-Page",
                        "in": "header"
                    },
                    {
                        "type": "integer",
                        "description": "Page size",
                        "name": "X-Page-Size",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.ProcessDTO"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/process/{code}/{uuid}": {
            "get": {
                "description": "Get process by UUID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "process"
                ],
                "summary": "Get process",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Code of Process",
                        "name": "code",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "UUID of Process",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.ProcessDTO"
                            }
                        }
                    }
                }
            }
        },
        "/v1/Health": {
            "get": {
                "description": "get the status of server.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Show the status of server.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.Payload": {
            "type": "object",
            "additionalProperties": true
        },
        "api.ProcessDTO": {
            "type": "object",
            "properties": {
                "changed_at": {
                    "type": "string",
                    "example": "2023-12-10T12:30:55.442484002-06:00"
                },
                "code": {
                    "type": "string",
                    "example": "requests"
                },
                "created_at": {
                    "type": "string",
                    "example": "2023-12-08T11:33:55.418484002-06:00"
                },
                "current_status": {
                    "$ref": "#/definitions/api.ProcessStatusDTO"
                },
                "payload": {
                    "$ref": "#/definitions/api.Payload"
                },
                "statuses": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.ProcessStatusDTO"
                    }
                },
                "uuid": {
                    "type": "string",
                    "example": "23c968a6-5fc5-4e42-8f59-a7f9c0d4999c"
                }
            }
        },
        "api.ProcessStatusDTO": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2023-12-08T11:33:55.418484002-06:00"
                },
                "name": {
                    "type": "string",
                    "example": "created"
                },
                "payload": {
                    "$ref": "#/definitions/api.Payload"
                }
            }
        },
        "api.ProcessSubmitResponse": {
            "description": "Response with UUID of created process.",
            "type": "object",
            "properties": {
                "uuid": {
                    "type": "string",
                    "example": "23c968a6-5fc5-4e42-8f59-a7f9c0d4999c"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Business Process Engine API",
	Description:      "This is the Business Process Engine API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
