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
            "url": "https://www.jkstack.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/agents": {
            "get": {
                "description": "获取节点列表",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "agents"
                ],
                "operationId": "/api/agents",
                "parameters": [
                    {
                        "type": "string",
                        "description": "节点类型,不指定则列出所有类型",
                        "name": "type",
                        "in": "query"
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "default": 1,
                        "description": "分页编号",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "minimum": 10,
                        "type": "integer",
                        "default": 20,
                        "description": "每页数量",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Success"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "payload": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/agents.info"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/agents/{id}": {
            "get": {
                "description": "获取某个节点信息",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "agents"
                ],
                "operationId": "/api/agents/info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "节点ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Success"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "payload": {
                                            "$ref": "#/definitions/agents.info"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/foo/{id}": {
            "get": {
                "description": "调用example类型的agent",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "foo"
                ],
                "operationId": "/api/foo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "节点ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Success"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "agents.info": {
            "type": "object",
            "properties": {
                "arch": {
                    "type": "string",
                    "enum": [
                        "i386",
                        "x86_64",
                        "..."
                    ],
                    "example": "操作系统位数"
                },
                "id": {
                    "type": "string",
                    "example": "agent_id"
                },
                "ip": {
                    "type": "string",
                    "example": "ip地址"
                },
                "mac": {
                    "type": "string",
                    "example": "mac地址"
                },
                "os": {
                    "type": "string",
                    "enum": [
                        "windows",
                        "linux"
                    ],
                    "example": "操作系统类型"
                },
                "platform": {
                    "type": "string",
                    "enum": [
                        "debian",
                        "centos",
                        "..."
                    ],
                    "example": "操作系统名称"
                },
                "type": {
                    "type": "string",
                    "example": "agent类型"
                },
                "version": {
                    "type": "string",
                    "example": "agent版本号"
                }
            }
        },
        "api.Success": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 0
                },
                "payload": {}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "agent-server",
	Description:      "jkagent management server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
