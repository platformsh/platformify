#################################################################################
# {{ .Assets.ServiceName }}-specific configuration

import base64
import json
import os
import sys
from urllib.parse import urlparse

# This variable should always match the primary database relationship name
{{ .Assets.EnvPrefix }}_DB_RELATIONSHIP = "{{ .Database }}"

# Helper function for decoding base64-encoded JSON variables.
def decode(variable):
    """Decodes a {{ .Assets.ServiceName }} environment variable.
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


# Import some {{ .Assets.ServiceName }} settings from the environment.
# Read more on {{ .Assets.ServiceName }} variables at {{ .Assets.Docs.Variables }}
if os.getenv("{{ .Assets.EnvPrefix }}_APPLICATION_NAME"):
    DEBUG = False

    if os.getenv("{{ .Assets.EnvPrefix }}_APP_DIR"):
        STATIC_ROOT = os.path.join(os.getenv("{{ .Assets.EnvPrefix }}_APP_DIR"), "static")

    if os.getenv("{{ .Assets.EnvPrefix }}_PROJECT_ENTROPY"):
        SECRET_KEY = os.getenv("{{ .Assets.EnvPrefix }}_PROJECT_ENTROPY")

    if os.getenv("{{ .Assets.EnvPrefix }}_ROUTES"):
        {{ lower .Assets.EnvPrefix }}_routes = decode(os.getenv("{{ .Assets.EnvPrefix }}_ROUTES"))
        ALLOWED_HOSTS = list(map(
            lambda key: urlparse(key).hostname,
            {{ lower .Assets.EnvPrefix }}_routes.keys(),
        ))

    # Database service configuration, post-build only.
    if os.getenv("{{ .Assets.EnvPrefix }}_RELATIONSHIPS") and {{ .Assets.EnvPrefix }}_DB_RELATIONSHIP:
        {{ lower .Assets.EnvPrefix }}_relationships = decode(os.getenv("{{ .Assets.EnvPrefix }}_RELATIONSHIPS"))
        if {{ .Assets.EnvPrefix }}_DB_RELATIONSHIP in {{ .Assets.EnvPrefix }}_relationships:
            db_settings = {{ lower .Assets.EnvPrefix }}_relationships[{{ .Assets.EnvPrefix }}_DB_RELATIONSHIP][0]
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
