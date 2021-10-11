cat << EOF > testfile
simple:
  topology:
    topology_name: simple
    resource_groups:
      - resource_group_name: new
        resource_group_type: openstack
        resource_definitions:
          - name: "${TFACON_YAML_PATH}"
            role: os_server
            flavor: {{ .flavorRef }}
            image: {{ .imageRef }}
            keypair: {{ .key_name }}
            networks: {{ .networks }}
EOF