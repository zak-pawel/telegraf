# Kubernetes Inventory Input Plugin

This plugin gathers metrics from [Kubernetes][kubernetes] resources.

> [!NOTE]
> This plugin requires Kubernetes version 1.11+.

The gathered resources include for example daemon sets, deployments, endpoints,
ingress, nodes, persistent volumes and many more.

> [!CRITICAL]
> This plugin produces high cardinality data, which when not controlled for will
> cause high load on your database. Please make sure to [filter][filtering] the
> produced metrics or configure your database to avoid cardinality issues!

⭐ Telegraf v1.10.0
🏷️ containers
💻 all

[kubernetes]: https://kubernetes.io/
[filtering]: /docs/CONFIGURATION.md#metric-filtering

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md#plugins

## Configuration

```toml @sample.conf
# Read metrics from the Kubernetes api
[[inputs.kube_inventory]]
  ## URL for the Kubernetes API.
  ## If empty in-cluster config with POD's service account token will be used.
  # url = ""

  ## URL for the kubelet, if set it will be used to collect the pods resource metrics
  # url_kubelet = "http://127.0.0.1:10255"

  ## Namespace to use. Set to "" to use all namespaces.
  # namespace = "default"

  ## Node name to filter to. No filtering by default.
  # node_name = ""

  ## Use bearer token for authorization.
  ## Ignored if url is empty and in-cluster config is used.
  # bearer_token = "/var/run/secrets/kubernetes.io/serviceaccount/token"

  ## Set response_timeout (default 5 seconds)
  # response_timeout = "5s"

  ## Optional Resources to exclude from gathering
  ## Leave them with blank with try to gather everything available.
  ## Values can be - "daemonsets", deployments", "endpoints", "ingress",
  ## "nodes", "persistentvolumes", "persistentvolumeclaims", "pods", "services",
  ## "statefulsets"
  # resource_exclude = [ "deployments", "nodes", "statefulsets" ]

  ## Optional Resources to include when gathering
  ## Overrides resource_exclude if both set.
  # resource_include = [ "deployments", "nodes", "statefulsets" ]

  ## selectors to include and exclude as tags.  Globs accepted.
  ## Note that an empty array for both will include all selectors as tags
  ## selector_exclude overrides selector_include if both set.
  # selector_include = []
  # selector_exclude = ["*"]

  ## Optional TLS Config
  ## Trusted root certificates for server
  # tls_ca = "/path/to/cafile"
  ## Used for TLS client certificate authentication
  # tls_cert = "/path/to/certfile"
  ## Used for TLS client certificate authentication
  # tls_key = "/path/to/keyfile"
  ## Send the specified TLS server name via SNI
  # tls_server_name = "kubernetes.example.com"
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = false

  ## Uncomment to remove deprecated metrics.
  # fieldexclude = ["terminated_reason"]
```

## Kubernetes Permissions

If using [RBAC authorization][rbac], you will need to create a cluster role to
list "persistentvolumes" and "nodes". You will then need to make an [aggregated
ClusterRole][agg] that will eventually be bound to a user or group.

[rbac]: https://kubernetes.io/docs/reference/access-authn-authz/rbac/
[agg]: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#aggregated-clusterroles

```yaml
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: influx:cluster:viewer
  labels:
    rbac.authorization.k8s.io/aggregate-view-telegraf: "true"
rules:
  - apiGroups: [""]
    resources: ["persistentvolumes", "nodes"]
    verbs: ["get", "list"]

---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: influx:telegraf
aggregationRule:
  clusterRoleSelectors:
    - matchLabels:
        rbac.authorization.k8s.io/aggregate-view-telegraf: "true"
    - matchLabels:
        rbac.authorization.k8s.io/aggregate-to-view: "true"
rules: [] # Rules are automatically filled in by the controller manager.
```

Bind the newly created aggregated ClusterRole with the following config file,
updating the subjects as needed.

```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: influx:telegraf:viewer
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: influx:telegraf
subjects:
  - kind: ServiceAccount
    name: telegraf
    namespace: default
```

## Quickstart in k3s

When monitoring [k3s](https://k3s.io) server instances one can re-use already
generated administration token. This is less secure than using the more
restrictive dedicated telegraf user but more convenient to set up.

```console
# replace `telegraf` with the user the telegraf process is running as
$ install -o telegraf -m400 /var/lib/rancher/k3s/server/token /run/telegraf-kubernetes-token
$ install -o telegraf -m400 /var/lib/rancher/k3s/server/tls/client-admin.crt /run/telegraf-kubernetes-cert
$ install -o telegraf -m400 /var/lib/rancher/k3s/server/tls/client-admin.key /run/telegraf-kubernetes-key
```

```toml
[kube_inventory]
bearer_token = "/run/telegraf-kubernetes-token"
tls_cert = "/run/telegraf-kubernetes-cert"
tls_key = "/run/telegraf-kubernetes-key"
```

## Metrics

- kubernetes_daemonset
  - tags:
    - daemonset_name
    - namespace
    - selector (\*varies)
  - fields:
    - generation
    - current_number_scheduled
    - desired_number_scheduled
    - number_available
    - number_misscheduled
    - number_ready
    - number_unavailable
    - updated_number_scheduled

- kubernetes_deployment
  - tags:
    - deployment_name
    - namespace
    - selector (\*varies)
  - fields:
    - replicas_available
    - replicas_unavailable
    - created

- kubernetes_endpoints
  - tags:
    - endpoint_name
    - namespace
    - hostname
    - node_name
    - port_name
    - port_protocol
    - kind (\*varies)
  - fields:
    - created
    - generation
    - ready
    - port

- kubernetes_ingress
  - tags:
    - ingress_name
    - namespace
    - hostname
    - ip
    - backend_service_name
    - path
    - host
  - fields:
    - created
    - generation
    - backend_service_port
    - tls

- kubernetes_node
  - tags:
    - node_name
    - status
    - condition
    - cluster_namespace
  - fields:
    - capacity_cpu_cores
    - capacity_millicpu_cores
    - capacity_memory_bytes
    - capacity_pods
    - allocatable_cpu_cores
    - allocatable_millicpu_cores
    - allocatable_memory_bytes
    - allocatable_pods
    - status_condition
    - spec_unschedulable
    - node_count

- kubernetes_persistentvolume
  - tags:
    - pv_name
    - phase
    - storageclass
  - fields:
    - phase_type (int, [see below](#pv-phase_type))

- kubernetes_persistentvolumeclaim
  - tags:
    - pvc_name
    - namespace
    - phase
    - storageclass
    - selector (\*varies)
  - fields:
    - phase_type (int, [see below](#pvc-phase_type))

- kubernetes_pod_container
  - tags:
    - container_name
    - namespace
    - node_name
    - pod_name
    - node_selector (\*varies)
    - phase
    - state
    - readiness
    - condition
  - fields:
    - restarts_total
    - state_code
    - state_reason
    - phase_reason
    - terminated_reason (string, deprecated in 1.15: use `state_reason` instead)
    - resource_requests_millicpu_units
    - resource_requests_memory_bytes
    - resource_limits_millicpu_units
    - resource_limits_memory_bytes
    - status_condition

- kubernetes_service
  - tags:
    - service_name
    - namespace
    - port_name
    - port_protocol
    - external_name
    - cluster_ip
    - selector (\*varies)
  - fields
    - created
    - generation
    - port
    - target_port

- kubernetes_statefulset
  - tags:
    - statefulset_name
    - namespace
    - selector (\*varies)
  - fields:
    - created
    - generation
    - replicas
    - replicas_current
    - replicas_ready
    - replicas_updated
    - spec_replicas
    - observed_generation

- kubernetes_resourcequota
  - tags:
    - resource
    - namespace
  - fields:
    - hard_cpu_limits
    - hard_cpu_requests
    - hard_memory_limit
    - hard_memory_requests
    - hard_pods
    - used_cpu_limits
    - used_cpu_requests
    - used_memory_limits
    - used_memory_requests
    - used_pods

- kubernetes_certificate
  - tags:
    - common_name
    - signature_algorithm
    - public_key_algorithm
    - issuer_common_name
    - san
    - verification
    - name
    - namespace
  - fields:
    - age
    - expiry
    - startdate
    - enddate
    - verification_code

### kubernetes node status `status`

The node status ready can mean 3 different values.

| Tag value | Corresponding field value | Meaning  |
| --------- | ------------------------- | -------- |
| ready     | 0                         | NotReady |
| ready     | 1                         | Ready    |
| ready     | 2                         | Unknown  |

### pv `phase_type`

The persistentvolume "phase" is saved in the `phase` tag with a correlated
numeric field called `phase_type` corresponding with that tag value.

| Tag value | Corresponding field value |
| --------- | ------------------------- |
| bound     | 0                         |
| failed    | 1                         |
| pending   | 2                         |
| released  | 3                         |
| available | 4                         |
| unknown   | 5                         |

### pvc `phase_type`

The persistentvolumeclaim "phase" is saved in the `phase` tag with a correlated
numeric field called `phase_type` corresponding with that tag value.

| Tag value | Corresponding field value |
| --------- | ------------------------- |
| bound     | 0                         |
| lost      | 1                         |
| pending   | 2                         |
| unknown   | 3                         |

## Example Output

```text
kubernetes_configmap,configmap_name=envoy-config,namespace=default,resource_version=56593031 created=1544103867000000000i 1547597616000000000
kubernetes_daemonset,daemonset_name=telegraf,selector_select1=s1,namespace=logging number_unavailable=0i,desired_number_scheduled=11i,number_available=11i,number_misscheduled=8i,number_ready=11i,updated_number_scheduled=11i,created=1527758699000000000i,generation=16i,current_number_scheduled=11i 1547597616000000000
kubernetes_deployment,deployment_name=deployd,selector_select1=s1,namespace=default replicas_unavailable=0i,created=1544103082000000000i,replicas_available=1i 1547597616000000000
kubernetes_node,host=vjain node_count=8i 1628918652000000000
kubernetes_node,condition=Ready,host=vjain,node_name=ip-172-17-0-2.internal,status=True status_condition=1i 1629177980000000000
kubernetes_node,cluster_namespace=tools,condition=Ready,host=vjain,node_name=ip-172-17-0-2.internal,status=True allocatable_cpu_cores=4i,allocatable_memory_bytes=7186567168i,allocatable_millicpu_cores=4000i,allocatable_pods=110i,capacity_cpu_cores=4i,capacity_memory_bytes=7291424768i,capacity_millicpu_cores=4000i,capacity_pods=110i,spec_unschedulable=0i,status_condition=1i 1628918652000000000
kubernetes_resourcequota,host=vjain,namespace=default,resource=pods-high hard_cpu=1000i,hard_memory=214748364800i,hard_pods=10i,used_cpu=0i,used_memory=0i,used_pods=0i 1629110393000000000
kubernetes_resourcequota,host=vjain,namespace=default,resource=pods-low hard_cpu=5i,hard_memory=10737418240i,hard_pods=10i,used_cpu=0i,used_memory=0i,used_pods=0i 1629110393000000000
kubernetes_persistentvolume,phase=Released,pv_name=pvc-aaaaaaaa-bbbb-cccc-1111-222222222222,storageclass=ebs-1-retain phase_type=3i 1547597616000000000
kubernetes_persistentvolumeclaim,namespace=default,phase=Bound,pvc_name=data-etcd-0,selector_select1=s1,storageclass=ebs-1-retain phase_type=0i 1547597615000000000
kubernetes_pod,namespace=default,node_name=ip-172-17-0-2.internal,pod_name=tick1 last_transition_time=1547578322000000000i,ready="false" 1547597616000000000
kubernetes_service,cluster_ip=172.29.61.80,namespace=redis-cache-0001,port_name=redis,port_protocol=TCP,selector_app=myapp,selector_io.kompose.service=redis,selector_role=slave,service_name=redis-slave created=1588690034000000000i,generation=0i,port=6379i,target_port=0i 1547597616000000000
kubernetes_pod_container,condition=Ready,host=vjain,pod_name=uefi-5997f76f69-xzljt,status=True status_condition=1i 1629177981000000000
kubernetes_pod_container,container_name=telegraf,namespace=default,node_name=ip-172-17-0-2.internal,node_selector_node-role.kubernetes.io/compute=true,pod_name=tick1,phase=Running,state=running,readiness=ready resource_requests_cpu_units=0.1,resource_limits_memory_bytes=524288000,resource_limits_cpu_units=0.5,restarts_total=0i,state_code=0i,state_reason="",phase_reason="",resource_requests_memory_bytes=524288000 1547597616000000000
kubernetes_statefulset,namespace=default,selector_select1=s1,statefulset_name=etcd replicas_updated=3i,spec_replicas=3i,observed_generation=1i,created=1544101669000000000i,generation=1i,replicas=3i,replicas_current=3i,replicas_ready=3i 1547597616000000000
```
