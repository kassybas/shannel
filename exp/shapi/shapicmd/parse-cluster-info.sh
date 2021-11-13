# shapi
#    input: cluster_info string

name="$(cut -d' ' -f1 <<<${cluster_info})"
path="config-sync/cluster-sync-configs/${name}"
location="$(cut -d' ' -f2 <<<${cluster_info})"
if [[ ! -f ${path} ]]; then
    echo "SKIP: no config found for cluster ${name}"
    exit 0
fi
snl push clusters .name=${name} .path=${path} .configPath=${path}
