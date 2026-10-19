terraform import [options] octopusdeploy_process_templated_step.<name> "<process-id>:<step-id>"

# Resources in a space other than the provider's can be imported by prefixing
# the identifier with that space:
terraform import [options] octopusdeploy_process_templated_step.<name> "<space-id>:<process-id>:<step-id>"
