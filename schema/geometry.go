package schema

// THIS FILE WAS COPIED BY HAND FROM https://raw.githubusercontent.com/whosonfirst/whosonfirst-json-schema/master/schema/geojson/geojson-geometry.json
// EVENTUALLY IT WILL BE CLONE (BY ROBOT) FROM THE SAME SOURCE. YOU SHOULD NOT UPDATE THIS FILE BY HAND.
// (20191227/thisisaaronland)

const GEOMETRY string = `{
    "$schema": "http://json-schema.org/draft-06/schema#",
    "$id": "geojson-geometry.json",
    "definitions": {
        "geometry": {
            "description": "Geometry as defined by GeoJSON",
            "type": "object",
            "required": [
                "type",
                "coordinates"
            ],
            "oneOf": [
                {
                    "title": "Point",
                    "properties": {
                        "type": {
                            "enum": [
                                "Point"
                            ]
                        },
                        "coordinates": {
                            "$ref": "#/definitions/position"
                        }
                    }
                },
                {
                    "title": "MultiPoint",
                    "properties": {
                        "type": {
                            "enum": [
                                "MultiPoint"
                            ]
                        },
                        "coordinates": {
                            "$ref": "#/definitions/positionArray"
                        }
                    }
                },
                {
                    "title": "LineString",
                    "properties": {
                        "type": {
                            "enum": [
                                "LineString"
                            ]
                        },
                        "coordinates": {
                            "$ref": "#/definitions/lineString"
                        }
                    }
                },
                {
                    "title": "MultiLineString",
                    "properties": {
                        "type": {
                            "enum": [
                                "MultiLineString"
                            ]
                        },
                        "coordinates": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/lineString"
                            }
                        }
                    }
                },
                {
                    "title": "Polygon",
                    "properties": {
                        "type": {
                            "enum": [
                                "Polygon"
                            ]
                        },
                        "coordinates": {
                            "$ref": "#/definitions/polygon"
                        }
                    }
                },
                {
                    "title": "MultiPolygon",
                    "properties": {
                        "type": {
                            "enum": [
                                "MultiPolygon"
                            ]
                        },
                        "coordinates": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/polygon"
                            }
                        }
                    }
                }
            ]
        },
        "position": {
            "description": "A single position",
            "type": "array",
            "minItems": 2,
            "items": [
                {
                    "type": "number"
                },
                {
                    "type": "number"
                }
            ],
            "additionalItems": false
        },
        "positionArray": {
            "description": "An array of positions",
            "type": "array",
            "items": {
                "$ref": "#/definitions/position"
            }
        },
        "lineString": {
            "description": "An array of two or more positions",
            "allOf": [
                {
                    "$ref": "#/definitions/positionArray"
                },
                {
                    "minItems": 2
                }
            ]
        },
        "linearRing": {
            "description": "An array of four positions where the first equals the last",
            "allOf": [
                {
                    "$ref": "#/definitions/positionArray"
                },
                {
                    "minItems": 4
                }
            ]
        },
        "polygon": {
            "description": "An array of linear rings",
            "type": "array",
            "items": {
                "$ref": "#/definitions/linearRing"
            }
        }
    }
}`
