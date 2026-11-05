output "vm_public_ip" {
  description = "Public ip address of terraform vm"
  value       = azurerm_public_ip.main.ip_address
}

output "vm_fqdn" {
  description = "Fully qualidied DNS name of the terraform vm"
  value       = azurerm_public_ip.main.fqdn
}