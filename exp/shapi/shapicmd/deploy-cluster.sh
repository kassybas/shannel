# shapi:
#   input:
#     cluter_name: string
#     cluster_location: string
#     cluster_configPath: string
gcloud beta container hub config-management apply \
    --project=${gcpProject} \
    --membership="${cluster_name}" \
    --config="${config}"
