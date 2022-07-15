pkill unnamed_plan_service_

path=$(pwd)

cd "$(dirname "$0")" || exit

function upgrade_service() {
  cd ./"$1" || exit
  chmod +x ./"$2"
  (./"$2" &)
  cd ..
}

upgrade_service "service_core" "unnamed_plan_service_core"
upgrade_service "service_gateway" "unnamed_plan_service_gateway"
upgrade_service "service_1_user" "unnamed_plan_service_1_user"
upgrade_service "service_2_cloud_file" "unnamed_plan_service_2_cloud_file"
upgrade_service "service_3_note" "unnamed_plan_service_3_note"
upgrade_service "service_4_task" "unnamed_plan_service_4_task"

cd "$path" || exit
