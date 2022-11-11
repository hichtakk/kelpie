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

Then you can call vSphere REST API. Kelpie accepts HTTP method as subcommand and API path as its argument.

Examples:
```
# get list of virtual machines
$ kelpie get /api/vcenter/vm
[                                                         
  {                                                       
    "cpu_count": 1,         
    "memory_size_MiB": 128,                               
    "name": "vCLS-cff254df-850f-4571-8ba1-628c021a0525",  
    "power_state": "POWERED_ON",
    "vm": "vm-1001"
  },
  {
    "cpu_count": 1,
    "memory_size_MiB": 128,
    "name": "vCLS-8eb757fb-efac-4dd6-b339-6d4946a91531",
    "power_state": "POWERED_ON",
    "vm": "vm-1002"
  },
  {
    "cpu_count": 1,
    "memory_size_MiB": 128,
    "name": "vCLS-81fe49e6-f922-4697-ae78-f500032a5309",
    "power_state": "POWERED_OFF",
    "vm": "vm-1003"
  }
]

# restart vm
$ kelpie post /api/vcenter/vm/${VM_ID}/power -q action=reset
```