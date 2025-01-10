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
        "/api/bookmarks/create": {
            "post": {
                "security": [
                    {
                        "CookieAuth": []
                    }
                ],
                "description": "This endpoint allows an authenticated user to create a new bookmark.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bookmark"
                ],
                "summary": "Create a bookmark",
                "parameters": [
                    {
                        "description": "Bookmark details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Bookmark"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Bookmark created successfully",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request - Invalid input",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized - User not authenticated",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    }
                }
            }
        },
        "/api/bookmarks/delete": {
            "delete": {
                "security": [
                    {
                        "CookieAuth": []
                    }
                ],
                "description": "Delete a specific bookmark associated with the user based on the bookmark ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bookmark"
                ],
                "summary": "Delete a bookmark by ID",
                "parameters": [
                    {
                        "description": "Request body with the bookmark ID to delete",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Bookmark"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully deleted the bookmark",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request, invalid input",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "403": {
                        "description": "Forbidden, the user does not have permission to delete this bookmark",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    }
                }
            }
        },
        "/api/bookmarks/get": {
            "get": {
                "security": [
                    {
                        "CookieAuth": []
                    }
                ],
                "description": "Fetch all bookmarks associated with the current user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bookmark"
                ],
                "summary": "Get bookmarks by user ID",
                "responses": {
                    "200": {
                        "description": "List of bookmarks",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.Bookmark"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    }
                }
            }
        },
        "/api/bookmarks/update": {
            "post": {
                "security": [
                    {
                        "CookieAuth": []
                    }
                ],
                "description": "Update a specific bookmark associated with the user based on the bookmark ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bookmark"
                ],
                "summary": "Update a bookmark by ID",
                "parameters": [
                    {
                        "description": "Request body with the bookmark ID to update",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Bookmark"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully updated the bookmark",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request, invalid input",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "403": {
                        "description": "Forbidden, the user does not have permission to update this bookmark",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    }
                }
            }
        },
        "/api/user/logout": {
            "delete": {
                "description": "This endpoint allows a user to log out by deleting all active sessions associated with the user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "User logout",
                "responses": {
                    "200": {
                        "description": "Logout successful",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "description": "This endpoint allows a user to log in using their username and password. If already logged in, a conflict response is returned.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "User login",
                "parameters": [
                    {
                        "description": "Username and password",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/pkg.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request - Invalid input",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized - Invalid username or password or not verified",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict - User already logged in",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    }
                }
            }
        },
        "/user/password-reset/request": {
            "post": {
                "description": "This endpoint allows a user to request a password reset by providing their email.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Request a password reset",
                "parameters": [
                    {
                        "description": "User email",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.RequestPasswordReset"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Password reset email sent",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request - Invalid input",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "404": {
                        "description": "Email not found",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    }
                }
            }
        },
        "/user/password-reset/reset": {
            "post": {
                "description": "This endpoint allows a user to reset their password by providing a valid reset token and a new password.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Reset user password",
                "parameters": [
                    {
                        "description": "Reset token and new password",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.ResetPassword"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Password reset successfully",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request - Invalid input",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "404": {
                        "description": "Reset token not found or expired",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    }
                }
            }
        },
        "/user/register": {
            "post": {
                "description": "Register a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "Username, email and password",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    }
                }
            }
        },
        "/user/verify-email": {
            "post": {
                "description": "This endpoint allows a user to verify their email address by providing the email and verification code.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Verify email address",
                "parameters": [
                    {
                        "description": "verification code",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.EmailVerifyRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Email verified successfully",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request - Invalid input",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "404": {
                        "description": "Verification code not found",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    }
                }
            }
        },
        "/user/verify-email/request": {
            "post": {
                "description": "Resends a verification token for the provided username.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Resend Verification Token",
                "parameters": [
                    {
                        "description": "Request body containing the username",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.RequestVerificationToken"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Verification token resent successfully",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request, invalid input",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "403": {
                        "description": "User has verified email",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/pkg.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Bookmark": {
            "type": "object",
            "properties": {
                "icon_url": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "show_text": {
                    "type": "boolean"
                },
                "title": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "pkg.LoginResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "pkg.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "requests.EmailVerifyRequest": {
            "type": "object",
            "required": [
                "code"
            ],
            "properties": {
                "code": {
                    "type": "string",
                    "minLength": 6
                }
            }
        },
        "requests.LoginRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "requests.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 8
                },
                "username": {
                    "type": "string",
                    "minLength": 3
                }
            }
        },
        "requests.RequestPasswordReset": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "requests.RequestVerificationToken": {
            "type": "object",
            "required": [
                "username"
            ],
            "properties": {
                "username": {
                    "type": "string",
                    "minLength": 3
                }
            }
        },
        "requests.ResetPassword": {
            "type": "object",
            "required": [
                "password",
                "token"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
