// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Acho Arnold",
            "email": "arnold@httpmock.dev"
        },
        "license": {
            "name": "AGPL-3.0",
            "url": "https://raw.githubusercontent.com/NdoleStudio/httpmock/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/lemonsqueezy/event": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Publish a lemonsqueezy event to the registered listeners",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Lemonsqueezy"
                ],
                "summary": "Consume a lemonsqueezy event",
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "$ref": "#/definitions/responses.NoContent"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.BadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.Unauthorized"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/responses.UnprocessableEntity"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.InternalServerError"
                        }
                    }
                }
            }
        },
        "/v1/projects": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Fetches the list of all projects available to the currently authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Projects"
                ],
                "summary": "List of projects",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.Ok-array_entities_Project"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.BadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.Unauthorized"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/responses.UnprocessableEntity"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.InternalServerError"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint creates a new project for a user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Projects"
                ],
                "summary": "Create a project",
                "parameters": [
                    {
                        "description": "project create payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.ProjectCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.Ok-entities_Project"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.BadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.Unauthorized"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/responses.UnprocessableEntity"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.InternalServerError"
                        }
                    }
                }
            }
        },
        "/v1/projects/{projectID}": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint updates a project for a user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Projects"
                ],
                "summary": "Update a project",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Project ID",
                        "name": "projectID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "project update payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.ProjectUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.Ok-entities_Project"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.BadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.Unauthorized"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/responses.UnprocessableEntity"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.InternalServerError"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint deletes a project",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Projects"
                ],
                "summary": "Delete a project",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Project ID",
                        "name": "projectID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.NoContent"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.BadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.Unauthorized"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.NotFound"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/responses.UnprocessableEntity"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.InternalServerError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entities.Project": {
            "type": "object",
            "required": [
                "created_at",
                "description",
                "id",
                "name",
                "subdomain",
                "updated_at",
                "user_id"
            ],
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2022-06-05T14:26:02.302718+03:00"
                },
                "description": {
                    "type": "string",
                    "example": "Mock API for an online store for selling shoes"
                },
                "id": {
                    "type": "string",
                    "example": "8f9c71b8-b84e-4417-8408-a62274f65a08"
                },
                "name": {
                    "type": "string",
                    "example": "Mock Stripe API"
                },
                "subdomain": {
                    "type": "string",
                    "example": "api"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2022-06-05T14:26:10.303278+03:00"
                },
                "user_id": {
                    "type": "string",
                    "example": "WB7DRDWrJZRGbYrv2CKGkqbzvqdC"
                }
            }
        },
        "requests.ProjectCreateRequest": {
            "type": "object",
            "required": [
                "description",
                "name",
                "subdomain"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "subdomain": {
                    "type": "string"
                }
            }
        },
        "requests.ProjectUpdateRequest": {
            "type": "object",
            "required": [
                "description",
                "name"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "responses.BadRequest": {
            "type": "object",
            "required": [
                "data",
                "message",
                "status"
            ],
            "properties": {
                "data": {
                    "type": "string",
                    "example": "The request body is not a valid JSON string"
                },
                "message": {
                    "type": "string",
                    "example": "The request isn't properly formed"
                },
                "status": {
                    "type": "string",
                    "example": "error"
                }
            }
        },
        "responses.InternalServerError": {
            "type": "object",
            "required": [
                "message",
                "status"
            ],
            "properties": {
                "message": {
                    "type": "string",
                    "example": "We ran into an internal error while handling the request."
                },
                "status": {
                    "type": "string",
                    "example": "error"
                }
            }
        },
        "responses.NoContent": {
            "type": "object",
            "required": [
                "message",
                "status"
            ],
            "properties": {
                "message": {
                    "type": "string",
                    "example": "item deleted successfully"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "responses.NotFound": {
            "type": "object",
            "required": [
                "message",
                "status"
            ],
            "properties": {
                "message": {
                    "type": "string",
                    "example": "cannot find item with ID [32343a19-da5e-4b1b-a767-3298a73703ca]"
                },
                "status": {
                    "type": "string",
                    "example": "error"
                }
            }
        },
        "responses.Ok-array_entities_Project": {
            "type": "object",
            "required": [
                "data",
                "message",
                "status"
            ],
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Project"
                    }
                },
                "message": {
                    "type": "string",
                    "example": "Request handled successfully"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "responses.Ok-entities_Project": {
            "type": "object",
            "required": [
                "data",
                "message",
                "status"
            ],
            "properties": {
                "data": {
                    "$ref": "#/definitions/entities.Project"
                },
                "message": {
                    "type": "string",
                    "example": "Request handled successfully"
                },
                "status": {
                    "type": "string",
                    "example": "success"
                }
            }
        },
        "responses.Unauthorized": {
            "type": "object",
            "required": [
                "data",
                "message",
                "status"
            ],
            "properties": {
                "data": {
                    "type": "string",
                    "example": "Make sure your Bearer token is set in the [Bearer] header in the request"
                },
                "message": {
                    "type": "string",
                    "example": "You are not authorized to carry out this request."
                },
                "status": {
                    "type": "string",
                    "example": "error"
                }
            }
        },
        "responses.UnprocessableEntity": {
            "type": "object",
            "required": [
                "data",
                "message",
                "status"
            ],
            "properties": {
                "data": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    }
                },
                "message": {
                    "type": "string",
                    "example": "validation errors while sending message"
                },
                "status": {
                    "type": "string",
                    "example": "error"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "",
	Schemes:          []string{"http", "https"},
	Title:            "HTTP Mock API",
	Description:      "Backend HTTP Mock API server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
