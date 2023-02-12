# stop all services
pkill unnamed_plan_service_

path=$(pwd)

  # change dir, for run this script from any path
  cd "$(dirname "$0")" || exit

  # define function: add execute permission and run service background
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

cd "$path" || exit
