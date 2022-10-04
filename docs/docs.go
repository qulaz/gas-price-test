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
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/graph": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get gas statistic graph",
                "operationId": "get-graph",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.GasResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.GasDayAverageResult": {
            "type": "object",
            "properties": {
                "average": {
                    "type": "number"
                },
                "date": {
                    "type": "string",
                    "example": "22-01-01 00:00"
                }
            }
        },
        "entity.GasHourFreqResult": {
            "type": "object",
            "properties": {
                "hour": {
                    "type": "integer",
                    "maximum": 24,
                    "minimum": 0
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "entity.GasMonthAmountResult": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "date": {
                    "type": "string",
                    "example": "22-01-01 00:00"
                }
            }
        },
        "entity.GasResult": {
            "type": "object",
            "properties": {
                "gasHourFreq": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.GasHourFreqResult"
                    }
                },
                "gasPerDayAverage": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.GasDayAverageResult"
                    }
                },
                "gasPerMonthAmount": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.GasMonthAmountResult"
                    }
                },
                "gasSpentTotal": {
                    "type": "number"
                }
            }
        },
        "v1.response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "message"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Gas Price Test Task",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}