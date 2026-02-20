#!/bin/bash

set -e                                                          # Salir si hay errores

SERVICE=""
HELP=false
DEPLOY=false
UNDO=false
DELETE=false
BUILD=false
DEBUG=false
BUILD_AND_DEPLOY=false
REPLICAS=1

# Parsear opciones
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --h | --help) HELP=true ;;                        # Activar la bandera si se proporciona --help
        --d | --deploy) DEPLOY=true ;;                    # Activar la bandera si se proporciona --deploy
        --u | --undo) UNDO=true ;;                        # Activar la bandera si se proporciona --undo
        --delete) DELETE=true ;;                          # Activar la bandera si se proporciona --delete
        --b | --build) BUILD=true ;;                      # Activar la bandera si se proporciona --build
        --bd | --build-deploy) BUILD_AND_DEPLOY=true ;;   # Activar ambas banderas
        --r | --replicas) REPLICAS="$2"; shift ;;         # Cambiar el número de réplicas
        --debug) DEBUG=true ;;                            # Activar la bandera si se proporciona --debug
        --s | --service) SERVICE="$2"; shift ;;           # Cambiar el nombre del servicio
        *) echo "Opción desconocida: $1"; exit 1 ;;
    esac
    shift
done

# Variables de configuración
COMPANY="celsia-internet"
CURRENT_DATETIME=$(date '+%Y%m%d%H%M%S')                        # Fecha y hora actual
OUTPUT_FILE="./deployments/oke.yml"                             # Archivo de salida
PORT=3300                                                       # Puerto de la aplicación
HOST="http://service-$SERVICE"                                  # Host de la aplicación
PATH_URL="/"                                                    # Ruta de la aplicación
CMD="$SERVICE"                                                  # Comando de la aplicación
APP="service-$SERVICE"                                          # Valor para reemplazar $ROLE
IMAGE="celsiainternet/services-$SERVICE"                        # Valor $IMAGE
VERSION=$(git describe --tags --abbrev=0)                       # Valor para reemplazar $VERSION obtenido de Git
RELEASE="$VERSION-$CURRENT_DATETIME"                            # Valor para reemplazar $RELEASE
BRANCH=$(git branch --show-current)                             # Valor para validar el namespace
PRD=false

if [ "$SERVICE" == "octopus" ]; then
  PATH_URL="/octopus"  
elif [ "$SERVICE" == "suspensiones" ]; then
  PATH_URL="/suspensiones"  
else
  PATH_URL="/octopus/$SERVICE"  
fi

if [ "$BRANCH" == "main" ]; then
    NAMESPACE="prd"    
    IMAGE_VERSION="$IMAGE:$VERSION"                             # Valor para reemplazar $IMAGE
    HISTORY_LIMIT=15                                            # Límite de historial de versiones
    PRODUCTION=true                                             # Bandera para indicar que es producción
    CPU_REQUEST="500m"                                          # CPU request
    CPU_LIMIT="500m"                                            # CPU limit
    MEMORY_REQUEST="512Mi"                                      # Memory request
    MEMORY_LIMIT="1024Mi"                                       # Memory limit
    MAX_PODS_AVAILABLE=1                                        # Número máximo de pods disponibles
    MAX_PODS_SURGE=1                                            # Número máximo de pods en exceso
    TEMPLATE_FILE="./deployments/oke-template.yml"              # Archivo de plantilla
    PRD=true                                                    # Flag par aindicar que el building es de prd
elif [ "$BRANCH" == "develop" ]; then
    NAMESPACE="qa"                                              # Valor para reemplazar $NS
    IMAGE_VERSION="$IMAGE:$VERSION-alpha"                       # Valor para reemplazar $IMAGE
    HISTORY_LIMIT=3                                             # Límite de historial de versiones
    PRODUCTION=false                                            # Bandera para indicar que es producción
    CPU_REQUEST="50m"                                           # CPU request
    CPU_LIMIT="100m"                                            # CPU limit
    MEMORY_REQUEST="64Mi"                                       # Memory request
    MEMORY_LIMIT="128Mi"                                        # Memory limit
    MAX_PODS_AVAILABLE=1                                        # Número máximo de pods disponibles
    MAX_PODS_SURGE=1                                            # Número máximo de pods en exceso    
    TEMPLATE_FILE="./deployments/oke-statefulset-qa.yml"      # Archivo de plantilla
else
    NAMESPACE="dev"                                             # Valor para reemplazar $NS
    IMAGE_VERSION="$IMAGE:$VERSION-beta"                        # Valor para reemplazar $IMAGE
    HISTORY_LIMIT=1                                             # Límite de historial de versiones    
    PRODUCTION=false                                            # Bandera para indicar que es producción
    CPU_REQUEST="50m"                                           # CPU request
    CPU_LIMIT="100m"                                            # CPU limit
    MEMORY_REQUEST="64Mi"                                       # Memory request
    MEMORY_LIMIT="128Mi"                                        # Memory limit
    MAX_PODS_AVAILABLE=1                                        # Número máximo de pods disponibles
    MAX_PODS_SURGE=1                                            # Número máximo de pods en exceso
    TEMPLATE_FILE="./deployments/oke-statefulset-dev.yml"      # Archivo de plantilla
fi

# Mostrar las opciones elegidas
echo "Opciones elegidas:"
[[ "$DEPLOY" == true ]] && echo " - Deploy: Activado"
[[ "$UNDO" == true ]] && echo " - Undo: Activado"
[[ "$DELETE" == true ]] && echo " - Delete: Activado"
[[ "$BUILD" == true ]] && echo " - Build: Activado"
[[ "$REPLICAS" -gt 0 ]] && echo " - Réplicas: $REPLICAS"
[[ "$PRODUCTION" == true ]] && echo " - Producción: Activado"
[[ "$PRODUCTION" == false ]] && echo " - Producción: Desactivado"
[[ "$DEBUG" == true ]] && echo " - Debug: Activado"
[[ "$SERVICE" != "" ]] && echo " - Servicio: $SERVICE"

# Función para mostrar el mensaje de ayuda
help() {
  echo "Uso: deploy.sh [opciones]"
  echo "Opciones:"
  echo "  --h, --help: Muestra este mensaje de ayuda."
  echo "  --d, --deploy --s, --service <service>: Despliega la aplicación en el clúster de Kubernetes."
  echo "  --u, --undo --s, --service <service>: Deshace el despliegue de la aplicación en el clúster de Kubernetes."
  echo "  --delete --s, --service <service>: Elimina el despliegue de la aplicación en el clúster de Kubernetes."
  echo "  --b, --build --s, --service <service>: Construye la imagen de Docker de la aplicación."
  echo "  --bd, --build-deploy --s, --service <service>: Construye y despliega la aplicación."
  echo "  --r, --replicas --s, --service <service>: Cambia el número de réplicas de la aplicación."
  echo "  --debug: Activa el modo de depuración."
  exit 0
}

# Función para construir la imagen de Docker
build_image() {
  local platform=$1
  local image=$2
  local dockerfile=$3
  local tag_latest=$4

  # Si tag_latest es true, taggear como latest
  if [ "$tag_latest" = true ]; then
      docker buildx build --no-cache --platform "$platform" \
        -t "$IMAGE:latest" \
        -t "$image" \
        -f "$dockerfile" --push .
      echo "Imagen $image y $IMAGE:latest creadas con éxito."
  else
      docker buildx build --no-cache --platform "$platform" \
        -t "$image" \
        -f "$dockerfile" --push .
      echo "Imagen $image creada con éxito."
  fi
}

# Función para aplicar el archivo de configuración de Kubernetes
apply_k8s() {
  local namespace=$1
  local file=$2

  kubectl apply -f "$file"
  kubectl -n "$namespace" get pods

  echo "Deploy $file en el namespace $namespace."
}

# Reemplazar valores en el archivo de plantilla y guardar en el archivo de salida
create_manifest() {
sed -e "s#\$PORT#$PORT#g" \
    -e "s#\$HOST#$HOST#g" \
    -e "s#\$PATH_URL#$PATH_URL#g" \
    -e "s#\$REPLICAS#$REPLICAS#g" \
    -e "s#\$COMPANY#$COMPANY#g" \
    -e "s#\$PRODUCTION#$PRODUCTION#g" \
    -e "s#\$DEBUG#$DEBUG#g" \
    -e "s#\$ROLE#$APP#g" \
    -e "s#\$NS#$NAMESPACE#g" \
    -e "s#\$IMAGE#$IMAGE_VERSION#g" \
    -e "s#\$CPU_REQUEST#$CPU_REQUEST#g" \
    -e "s#\$CPU_LIMIT#$CPU_LIMIT#g" \
    -e "s#\$MAX_PODS_AVAILABLE#$MAX_PODS_AVAILABLE#g" \
    -e "s#\$MAX_PODS_SURGE#$MAX_PODS_SURGE#g" \
    -e "s#\$MEMORY_REQUEST#$MEMORY_REQUEST#g" \
    -e "s#\$MEMORY_LIMIT#$MEMORY_LIMIT#g" \
    -e "s#\$HISTORY_LIMIT#$HISTORY_LIMIT#g" \
    -e "s#\$RELEASE#$RELEASE#g" "$TEMPLATE_FILE" > "$OUTPUT_FILE"

echo "Archivo $OUTPUT_FILE generado con éxito."
}

if [ "$HELP" = true ]; then
  help
elif [ "$SERVICE" == "" ]; then
  echo "Error: Debes especificar el nombre del servicio con la opción --service"
  exit 1
elif [ "$UNDO" = true ]; then
  create_manifest
  kubectl rollout undo deployment "$APP" -n "$NAMESPACE"
  kubectl -n "$namespace" get pods
  echo "Desplegado deshecho."
elif [ "$BUILD" = true ]; then
  build_image "linux/amd64,linux/arm64" "$IMAGE_VERSION" "./cmd/$CMD/Dockerfile" $PRD
elif [ "$DELETE" = true ]; then
  create_manifest
  kubectl delete -f "$OUTPUT_FILE"
  kubectl -n "$namespace" get all
  echo "Despliegue eliminado."
elif [ "$DEPLOY" = true ]; then
  create_manifest
  apply_k8s "$NAMESPACE" "$OUTPUT_FILE"
elif [ "$BUILD_AND_DEPLOY" = true ]; then
  build_image "linux/amd64,linux/arm64" "$IMAGE_VERSION" "./cmd/$CMD/Dockerfile" $PRD
  create_manifest
  apply_k8s "$NAMESPACE" "$OUTPUT_FILE"
else
  help
fi

# Línea en blanco al final