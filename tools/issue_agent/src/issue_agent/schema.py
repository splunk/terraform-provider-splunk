"""Validate the JSON subsets stored next to this package."""

from __future__ import annotations

import json
from pathlib import Path
from typing import Any

SCHEMA_DIR = Path(__file__).resolve().parent / "schemas"


class SchemaError(ValueError):
    """Raised when an agent result does not match the strict schema."""


def load_schema(name: str) -> dict[str, Any]:
    path = SCHEMA_DIR / name
    return json.loads(path.read_text(encoding="utf-8"))


def validate(instance: Any, schema: dict[str, Any], path: str = "$") -> None:
    if "enum" in schema and instance not in schema["enum"]:
        raise SchemaError(f"{path} must be one of {schema['enum']}")

    expected = schema.get("type")
    if expected is not None and not _type_ok(instance, expected):
        raise SchemaError(f"{path} has the wrong type")

    if instance is None:
        return

    if schema.get("type") == "object" or (
        isinstance(expected, list) and "object" in expected and isinstance(instance, dict)
    ):
        _validate_object(instance, schema, path)
        return

    if schema.get("type") == "array" or (
        isinstance(expected, list) and "array" in expected and isinstance(instance, list)
    ):
        _validate_array(instance, schema, path)


def _validate_object(instance: dict[str, Any], schema: dict[str, Any], path: str) -> None:
    if not isinstance(instance, dict):
        raise SchemaError(f"{path} must be an object")
    required = schema.get("required", [])
    missing = [key for key in required if key not in instance]
    if missing:
        raise SchemaError(f"{path} is missing {missing}")
    properties = schema.get("properties", {})
    if schema.get("additionalProperties") is False:
        extra = [key for key in instance if key not in properties]
        if extra:
            raise SchemaError(f"{path} has unknown fields {extra}")
    for key, value in instance.items():
        if key in properties:
            validate(value, properties[key], f"{path}.{key}")


def _validate_array(instance: list[Any], schema: dict[str, Any], path: str) -> None:
    if not isinstance(instance, list):
        raise SchemaError(f"{path} must be an array")
    minimum = schema.get("minItems")
    if minimum is not None and len(instance) < minimum:
        raise SchemaError(f"{path} must contain at least {minimum} items")
    item_schema = schema.get("items")
    if item_schema is None:
        return
    for index, item in enumerate(instance):
        validate(item, item_schema, f"{path}[{index}]")


def _type_ok(instance: Any, expected: str | list[str]) -> bool:
    names = expected if isinstance(expected, list) else [expected]
    checks = {
        "object": lambda value: isinstance(value, dict),
        "array": lambda value: isinstance(value, list),
        "string": lambda value: isinstance(value, str),
        "boolean": lambda value: isinstance(value, bool),
        "integer": lambda value: isinstance(value, int) and not isinstance(value, bool),
        "null": lambda value: value is None,
    }
    return any(checks[name](instance) for name in names if name in checks)
