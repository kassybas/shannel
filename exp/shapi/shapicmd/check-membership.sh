# shapi:
#   input:
#     cluter_name: string
#     cluster_location: string
#     cluster_configPath: string
gcloud container hub memberships list --project=${gcpProject} | grep "${cluster_name}"
