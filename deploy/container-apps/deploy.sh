#!/bin/bash
set -e

cd "$(dirname "$0")"

which az > /dev/null || { echo "💥 Error! Azure CLI not found, please install https://aka.ms/azure-cli"; exit 1; }

for varName in AZURE_REGION AZURE_RESGRP AZURE_BASE_NAME IMAGE_TAG IMAGE_REG IMAGE_NAME; do
  varVal=$(eval echo "\${$varName}")
  [ -z "$varVal" ] && { echo "💥 Error! Required variable '$varName' is unset!"; varUnset=true; }
done
[ $varUnset ] && exit 1

SUB_NAME=$(az account show --query name -o tsv)
if [ -z "$SUB_NAME" ]; then
  echo -e "💥 \e[31mYou are not logged into Azure!"
  exit 1
fi
TENANT_ID=$(az account show --query tenantId -o tsv)

echo -e "\e[34mAzure details: \e[0m"
echo -e "  🔑 \e[34mSubscription: \e[33m$SUB_NAME\e[0m"
echo -e "  🌐 \e[34mTenant:       \e[33m$TENANT_ID\e[0m"

# Parameters
imageFrontend="$IMAGE_REG/$IMAGE_NAME-frontend:$IMAGE_TAG"
imageBackend="$IMAGE_REG/$IMAGE_NAME-backend:$IMAGE_TAG"
imageDatabase="$IMAGE_REG/$IMAGE_NAME-database:$IMAGE_TAG"
deploymentName="nano-realms"

echo -e "\n\e[34m🏆 Deployinging into '$AZURE_RESGRP' in region '$AZURE_REGION'...\e[0m"
echo -e "\e[34m🚀 Bicep deployment started..."
az deployment sub create                   \
  --template-file nano-realms.bicep        \
  --location "$AZURE_REGION"               \
  --name $deploymentName                   \
  --parameters baseName="$AZURE_BASE_NAME" \
     resGroupName="$AZURE_RESGRP"          \
     location="$AZURE_REGION"              \
     imageFrontend="$imageFrontend"        \
     imageBackend="$imageBackend"          \
     imageDatabase="$imageDatabase"                    

appURL=$(az deployment sub show --name nano-realms --query "properties.outputs.appURL.value" -o tsv)

echo -e "\e[34m🏅 Bicep deployment completed!"
echo -e "\e[34m🌍 Play the game here: \e[33m$appURL\e[0m"
