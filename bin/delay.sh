# Don't blame me for you running my script kiddie tc copy pasta.
sudo tc qdisc del dev lo root

sudo tc qdisc add dev lo root handle 1: prio
sudo tc qdisc add dev lo parent 1:3 handle 30: netem delay 100ms

endpoints=$(cat config/endpoints.txt)

for port in ${endpoints//,/ }
do
  sudo tc filter add dev lo protocol ip parent 1:0 u32 \
    match ip sport $port 0xffff flowid 1:3
  sudo tc filter add dev lo protocol ip parent 1:0 u32 \
    match ip dport $port 0xffff flowid 1:3
done
