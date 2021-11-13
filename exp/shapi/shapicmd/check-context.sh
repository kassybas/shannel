# shapi:
#   input:
#     cluter_name: string
#     cluster_location: string
#     cluster_configPath: string
kubectl config current-context | grep "${cluster_name}"
