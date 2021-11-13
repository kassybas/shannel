# shapi:
#   vars:
#
#   input:
#     cluter_name: string
#     cluster_location: string
#     cluster_configPath: string

gcloud container clusters get-credentials "${cluster_name}" \
    --region "${cluster_location}" \
    --project=${gcpProject}
