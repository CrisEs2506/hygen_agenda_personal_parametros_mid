
#!/usr/bin/env bash

set -e
set -u
set -o pipefail

#if [ -n "${PARAMETER_STORE:-}" ]; then
#  export AGENDA_PERSONAL_PARAMETROS_API_MID_PGUSER="$(aws ssm get-parameter --name /${PARAMETER_STORE}agenda_personal_parametros_api_mid/db/username --output text --query Parameter.Value)"
#  export AGENDA_PERSONAL_PARAMETROS_API_MID_PGPASS="$(aws ssm get-parameter --with-decryption --name /${PARAMETER_STORE}/agenda_personal_parametros_api_mid/db/password --output text --query Parameter.Value)"

exec ./main "$@"