#################################################################################
# Platform.sh-specific configuration

import base64
import json
import os
import sys
from urllib.parse import urlparse

# This variable should always match the primary database relationship name,
#   configured in .platform.app.yaml.
PLATFORMSH_DB_RELATIONSHIP = "{{ .Database }}"

# Helper function for decoding base64-encoded JSON variables.
def decode(variable):
    """Decodes a Platform.sh environment variable.
    Args:
        variable (string):
            Base64-encoded JSON (the content of an environment variable).
    Returns:
        An dict (if representing a JSON object), or a scalar type.
    Raises:
        JSON decoding error.
    """
    try:
        if sys.version_info[1] > 5:
            return json.loads(base64.b64decode(variable))
        else:
            return json.loads(base64.b64decode(variable).decode("utf-8"))
    except json.decoder.JSONDecodeError:
        print("Error decoding JSON, code %d", json.decoder.JSONDecodeError)


# Import some Platform.sh settings from the environment.
# Read more on Platform.sh variables at https://docs.platform.sh/development/variables/use-variables.html#use-platformsh-provided-variables
if os.getenv("PLATFORM_APPLICATION_NAME"):
    DEBUG = False

    if os.getenv("PLATFORM_APP_DIR"):
        STATIC_ROOT = os.path.join(os.getenv("PLATFORM_APP_DIR"), "static")

    if os.getenv("PLATFORM_PROJECT_ENTROPY"):
        SECRET_KEY = os.getenv("PLATFORM_PROJECT_ENTROPY")

    if os.getenv("PLATFORM_ROUTES"):
        platform_routes = decode(os.getenv("PLATFORM_ROUTES"))
        ALLOWED_HOSTS = list(map(
            lambda key: urlparse(key).hostname,
            platform_routes.keys(),
        ))

    # Database service configuration, post-build only.
    if os.getenv("PLATFORM_RELATIONSHIPS") and PLATFORMSH_DB_RELATIONSHIP:
        platform_relationships = decode(os.getenv("PLATFORM_RELATIONSHIPS"))
        if PLATFORMSH_DB_RELATIONSHIP in platform_relationships:
            db_settings = platform_relationships[PLATFORMSH_DB_RELATIONSHIP][0]
            engine = None
            if (
                "mariadb" in db_settings["type"]
                or "mysql" in db_settings["type"]
                or "oracle-mysql" in db_settings["type"]
            ):
                engine = "django.db.backends.mysql"
            elif "postgresql" in db_settings["type"]:
                engine = "django.db.backends.postgresql"

            if engine:
                DATABASES = {
                    "default": {
                        "ENGINE": engine,
                        "NAME": db_settings["path"],
                        "USER": db_settings["username"],
                        "PASSWORD": db_settings["password"],
                        "HOST": db_settings["host"],
                        "PORT": db_settings["port"],
                    },
                }
