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
            "name": "Anh Cao",
            "email": "anhcao4922@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/market-price": {
            "post": {
                "description": "Fetch the market spot price of electric in Finland in any times",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "market-price"
                ],
                "summary": "Retrieves the market price",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.PriceResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/market-price/today-tomorrow": {
            "get": {
                "description": "Returns the exchange price for today and tomorrow.\nIf tomorrow's price is not available yet, return empty struct.\nThen client needs to show readable information to indicate that data is not available yet.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "market-price"
                ],
                "summary": "Retrieves the market price for today and tomorrow",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.TodayTomorrowPrice"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/price-settings": {
            "get": {
                "description": "retrieves the price settings for specific user by identify through 'access token'.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "price-settings"
                ],
                "summary": "Retrieves the price settings for specific user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.PriceSettings"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a new price settings for new user by identify through 'access token'.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "price-settings"
                ],
                "summary": "Creates a new price settings for user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes the price settings for specific user by identify through 'access token'.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "price-settings"
                ],
                "summary": "Deletes the price settings for specific user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "patch": {
                "description": "Updates the price settings for specific user by identify through 'access token'.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "price-settings"
                ],
                "summary": "Updates the price settings for specific user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.DailyPrice": {
            "type": "object",
            "properties": {
                "available": {
                    "type": "boolean"
                },
                "prices": {
                    "$ref": "#/definitions/models.PriceSeries"
                }
            }
        },
        "models.Data": {
            "type": "object",
            "properties": {
                "includeVat": {
                    "description": "IncludeVat is legacy property that return string value and value \"0\" means no VAT included and string \"1\" is included",
                    "type": "string"
                },
                "isToday": {
                    "type": "boolean"
                },
                "orig_time": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "time": {
                    "type": "string"
                },
                "time_utc": {
                    "type": "string"
                },
                "vat_factor": {
                    "type": "number"
                }
            }
        },
        "models.PriceData": {
            "type": "object",
            "properties": {
                "group": {
                    "description": "Group represents 'hour', 'day', 'week', 'month', 'year'",
                    "type": "string"
                },
                "series": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.PriceSeries"
                    }
                }
            }
        },
        "models.PriceResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.PriceData"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.PriceSeries": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Data"
                    }
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models.PriceSettings": {
            "type": "object",
            "properties": {
                "margin": {
                    "type": "number"
                },
                "user_id": {
                    "type": "string"
                },
                "vat_included": {
                    "type": "boolean"
                }
            }
        },
        "models.TodayTomorrowPrice": {
            "type": "object",
            "properties": {
                "today": {
                    "$ref": "#/definitions/models.DailyPrice"
                },
                "tomorrow": {
                    "$ref": "#/definitions/models.DailyPrice"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:5001",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Stormbreaker API (electric service)",
	Description:      "Service for retrieving information about market electric price in Finland",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
