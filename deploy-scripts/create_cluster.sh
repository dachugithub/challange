#!/usr/bin/env bash
BASEDIR=$(dirname "$0")
export KOPS_STATE_STORE="precreated with terraform"
export CLUSTER_NAME="challange.revolut.com"

echo -e "Creating cluster now:\n"

kops create cluster                                \
   --state $KOPS_STATE_STORE                        \
   --zones eu-west-2a,eu-west-2b,eu-west-2c         \
   --master-zones eu-west-2a,eu-west-2b,eu-west-2c  \
   --master-count 5                                 \
   --node-count 6                                   \
   --cloud=aws                                      \
   --networking weave                               \
   --authorization RBAC                             \
   --bastion --topology=private                     \
   --name ${CLUSTER_NAME} --yes

# Waiting 8 minutes to the cluster stabilize
echo "Waiting ~8 minutes to stabilize the creation"
sleep 400

echo -e "\nDeploying helm...\n"
cat <<EOF | kubectl create -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tiller
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: tiller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: tiller
    namespace: kube-system
EOF

sleep 1
helm init --service-account tiller

echo -e "\nWaiting for helm\n"
sleep 30

echo -e "\nDeploying dashbord\n"
helm install stable/kubernetes-dashboard --name k8s-dashboard --set rbac.create=true

sleep 1

echo -e "\nDeploying external-dns\n"
# helm install --name external-dns stable/external-dns -f external-dns/values.yaml
cat <<EOF | helm install --name external-dns stable/external-dns -f -
rbac:
  create: true
domainFilters: [ "revolut.com." ]
aws:
  region: "eu-west-2"
policy: sync
sources:
  - service
EOF

echo -e "\nDeploying secrets\n"
# needs to be encrypted ie. git-crypt
helm install --namespace ${ENV} --name secrets Charts/birthday/secrets/ -f secret-values.yaml


echo -e "\nDeploying monitoring\n"
git clone -b master https://github.com/coreos/prometheus-operator.git prometheus-operator-temp;
cd prometheus-operator-temp/contrib/kube-prometheus
./hack/cluster-monitoring/self-hosted-deploy
cd -
rm -rf prometheus-operator-temp

sleep 3

# RBAC is needed to allow prometheus-k8s user access to the namesspace ie.
cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: prometheus-k8s
  namespace: challange
rules:
- apiGroups: [""]
  resources:
  - services
  - endpoints
  - pods
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: prometheus-k8s
  namespace: challange
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: prometheus-k8s
subjects:
- kind: ServiceAccount
  name: prometheus-k8s
  namespace: monitoring
EOF

# Deploy the endpoints/ingress for monitoring
cat <<EOF | kubectl apply -f -
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: monitoring-ingress
  namespace: monitoring
spec:
  rules:
  - host: alerts.${CLUSTER_NAME}
    http:
      paths:
      - backend:
          serviceName: alertmanager-main
          servicePort: 9093
        path: /
  - host: prom.${CLUSTER_NAME}
    http:
      paths:
      - backend:
          serviceName: prometheus-k8s
          servicePort: 9090
        path: /
  - host: grafana.${CLUSTER_NAME}
    http:
      paths:
      - backend:
          serviceName: grafana
          servicePort: 3000
        path: /
EOF

# Now we deploy the load balancer for monitoring
cat <<EOF | helm install --wait --namespace monitoring --name monitoring-lb stable/nginx-ingress -f -
rbac:
  create: true
controller:
  scope:
    enabled: true
    namespace: monitoring
  service:
    annotations:
      external-dns.alpha.kubernetes.io/hostname: alerts.${CLUSTER_NAME},prom.${CLUSTER_NAME},grafana.${CLUSTER_NAME}
    loadBalancerSourceRanges: [ "<restricted to operators network >" ]
EOF


