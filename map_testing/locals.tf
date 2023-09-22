locals {
    pluginsLocal = {
    for k, v in var.pluginMap : regex(".*-(.*)-.*", k)[0] => v
    }
}
