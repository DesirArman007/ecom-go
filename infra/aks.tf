resource "azurerm_subnet" "aks" {
  name                 = "ecom-aks-subnet"
  resource_group_name  = azurerm_resource_group.main.name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = ["10.0.4.0/22"]
}

resource "azurerm_user_assigned_identity" "aks" {
  name                = "ecom-aks-identity"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
}

resource "azurerm_role_assignment" "aks_network_contributor" {
  scope                = azurerm_subnet.aks.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.aks.principal_id
}

resource "azurerm_kubernetes_cluster" "aks" {
  name                = "ecom-aks"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  dns_prefix          = "ecom-aks"
  sku_tier            = "Free"

  node_provisioning_profile {
    mode = "Manual"
  }

  default_node_pool {
    name            = "system"
    node_count      = 1
    vm_size         = "Standard_B2s_v2"
    os_disk_size_gb = 30
    vnet_subnet_id  = azurerm_subnet.aks.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.aks.id]
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "standard"
    outbound_type     = "loadBalancer"
    service_cidr      = "10.2.0.0/16"
    dns_service_ip    = "10.2.0.10"
  }

  oidc_issuer_enabled       = true
  workload_identity_enabled = true

  depends_on = [
    azurerm_role_assignment.aks_network_contributor
  ]
}