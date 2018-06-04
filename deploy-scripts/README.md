assumption:
iam accounts and kops state bucket are created (terraform ?)
helm is installed:
https://github.com/kubernetes/helm
kubectl is installed:
https://kubernetes.io/docs/tasks/tools/install-kubectl/

to create cluster run create_cluster.sh

to install deployment helm.deplopy

there are many missing parts:
1. where do we take docker image from (private registry?)
2. is psql created in rds and initialized ? (hopefully migrations)
etc.

 
