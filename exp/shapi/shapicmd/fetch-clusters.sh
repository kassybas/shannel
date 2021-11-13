# shapi:
#   input: null
#   output: ClusterInfo: string~'?P<cluster_name>[a-z\-]+'

gcloud container clusters list \
    --project="${gcpProject}" |
    grep "${clusterFilter}" |
    snl push -d "\n"
