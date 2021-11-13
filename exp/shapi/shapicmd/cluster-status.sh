echo "Container hub config-management spec applied ✅"
echo "  Console: https://console.cloud.google.com/kubernetes/config_management?project=${gcpProject}"
echo "  Status: "
gcloud beta container hub config-management status
