terraform import [options] octopusdeploy_process.<name> <process-id>

# Resources in a space other than the provider's can be imported by prefixing
# the identifier with that space:
terraform import [options] octopusdeploy_process.<name> "<space-id>:<process-id>"
