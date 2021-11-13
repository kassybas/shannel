if [[ $GHE_PRIVATE_KEY_FILE == "" ]]; then
    echo "ERROR: GHE private key not found: please set GHE_PRIVATE_KEY_FILE variable"
    exit 1
fi
kubectl create secret generic git-creds \
    --namespace=config-management-system \
    --from-file="ssh=${GHE_PRIVATE_KEY_FILE}"
