kelpie
======

Simple vSphere REST API client.

# Usage
Set vCenter endpoint as environment variables.

```
export KELPIE_VCENTER_SERVER=${YOUR_VCENTER_URL}
export KELPIE_VCENTER_USER=${YOUR_VCENTER_USERNAME}
export KELPIE_VCENTER_PASSWORD=${YOUR_VCENTER_USER_PASSWORD}
```

Then you can call vSphere REST API. Keipie accept HTTP method as subcommand and API path as its argument.

Examples:
```
# get list of virtual machines
$ kelpie get /api/vcenter/vm

# restart vm
$ kelpie post /api/vcenter/vm/${VM_ID}/power -q action=reset
```