// ============================================================================
// Deploy nano-relams as a set of a container apps
// ============================================================================

targetScope = 'subscription'

@description('Base name for all resources, suffixes will be added where needed')
param baseName string = 'nano-realms'
@description('Resource group name')
param resGroupName string = 'nano-realms'

@description('Azure region for all resources')
param location string = deployment().location

@description('Frontend container image')
param imageFrontend string = 'ghcr.io/benc-uk/nano-realms-frontend:latest'
@description('Backend container image')
param imageBackend string = 'ghcr.io/benc-uk/nano-realms-backend:latest'
@description('Database container image')
param imageDatabase string = 'ghcr.io/benc-uk/nano-realms-database:latest'

// ===== Variables ============================================================

// ===== Modules & Resources ==================================================

resource resGroup 'Microsoft.Resources/resourceGroups@2021-04-01' = {
  name: resGroupName
  location: location
}

module logAnalytics './modules/monitoring/log-analytics.bicep' = {
  scope: resGroup
  name: 'monitoring'
  params: {
    name: baseName
  }
}

module network './modules/network/network-multi.bicep' = {
  scope: resGroup
  name: 'network'
  params: {
    name: baseName
    addressSpace: '10.75.0.0/16'
    subnets: [
      {
        name: 'controlplane'
        cidr: '10.75.0.0/21'
      }
      {
        name: 'apps'
        cidr: '10.75.8.0/21'
      }
    ]
  }
}

module containerAppEnv './modules/containers/app-env.bicep' = {
  scope: resGroup
  name: 'containerAppEnv'
  params: {
    name: baseName
    logAnalyticsName: logAnalytics.outputs.name
    logAnalyticsResGroup: resGroup.name
    controlPlaneSubnetId: network.outputs.subnets[0].id
    appsSubnetId: network.outputs.subnets[1].id
  }
}

module frontend './modules/containers/app.bicep' = {
  scope: resGroup
  name: 'frontend'

  params: {
    name: '${baseName}-frontend'
    environmentId: containerAppEnv.outputs.id
    image: imageFrontend
    ingressPort: 8001
    ingressExternal: true
    probePath: '/'
    probePort: 8001

    envs: [
      {
        name: 'API_ENDPOINT'
        value: 'https://${backend.outputs.fqdn}'
      }
    ]
  }
}

module backend './modules/containers/app.bicep' = {
  scope: resGroup
  name: 'backend'

  params: {
    name: '${baseName}-backend'
    environmentId: containerAppEnv.outputs.id
    image: imageBackend
    ingressPort: 8000
    ingressExternal: true
    probePath: '/health'
    probePort: 8000

    envs: [
      {
        name: 'DB_HOST'
        value: '${baseName}-database'
      }
    ]
  }
}

module database './modules/containers/app.bicep' = {
  scope: resGroup
  name: 'database'

  params: {
    name: '${baseName}-database'
    environmentId: containerAppEnv.outputs.id
    image: imageDatabase
    ingressPort: 7687
    ingressExternal: false
    ingressTransport: 'tcp'
    ingressExposedTcpPort: 7687
    memory: '1Gi'
    cpu: '0.5'

    envs: [
      {
        name: 'NEO4J_AUTH'
        value: 'none'
      }
    ]
  }
}

// ===== Outputs ==============================================================

output appURL string = 'https://${frontend.outputs.fqdn}'
