# kube-proxy-healthz

The eBPF-based kube-proxy replacement can coexist with the original `kube-proxy`. eBPF handles traffic
earlier in the network stack, before it reaches the netfilter layer where `kube-proxy` rules are applied.
See the
Cilium [kube-proxy-hybrid-modes](https://docs.cilium.io/en/latest/network/kubernetes/kubeproxy-free/#kube-proxy-hybrid-modes)
for more information.

The `kube-proxy-healthz` is a simple health check utility designed to replace the default `kube-proxy`
binary in RKE1 (Rancher Kubernetes Engine) clusters. The `kube-proxy` container in RKE1 is hardcoded and
cannot be disabled. If the `kube-proxy` container does not report health, it is restarted, and the node is
marked as `NotReady`. This utility allows to run the RKE1 cluster without the default `kube-proxy` by
providing a health check that mimics the behavior of the `kube-proxy`.

There is an issue with `kube-apiserver` high availability (HA) that has been addressed
in the Cilium Issue [#37601](https://github.com/cilium/cilium/pull/37601) and documented
in [kubernetes-api-server-high-availability](https://docs.cilium.io/en/latest/network/kubernetes/kubeproxy-free/#kubernetes-api-server-high-availability).
However, in Rancher-managed clusters, there is
the [nginx-proxy](https://ranchermanager.docs.rancher.com/troubleshooting/kubernetes-components/troubleshooting-nginx-proxy)
on every node managed by the Rancher agent. The Kubernetes API is always accessible at `127.0.0.1:6443`.

## Migration Steps

1. Ensure kube-proxy replacement (e.g. Cilium) is deployed on all nodes.
2. Build the `kube-proxy-healthz` binary using the provided Makefile.
3. Deploy the `kube-proxy-healthz` binary to all RKE1 nodes where `kube-proxy` is running.
4. Bind mount the `kube-proxy-healthz` binary into the `kube-proxy` container to replace the default binary.
    ```yaml
    kubeproxy:
      extra_binds:
        - '/usr/local/bin/kube-proxy-healthz:/usr/local/bin/kube-proxy'
   ```
5. Restart the nodes to ensure that the old kube-proxy iptables rules are cleared. Alternatively, use the
   `kube-proxy --cleanup` command.

See the Rancher Terraform
provider [rke_config.services.kubeproxy.extra_binds](https://registry.terraform.io/providers/rancher/rke/latest/docs/resources/cluster#services-2)
resource for documentation.

## References

- RKE Issue [#1432](https://github.com/rancher/rke/issues/1432).
- RKE Cluster [BuildRKEConfigNodePlan](https://github.com/rancher/rke/blob/release/v1.7/cluster/plan.go#L122-L172).
- Kubernetes [kube-proxy healthz HTTP server](https://github.com/kubernetes/kubernetes/blob/release-1.32/pkg/proxy/healthcheck/proxier_health.go#L165-L236).
