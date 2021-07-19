echo "Importing Env Vars from .env file."
export $(grep -v '^#' .env | xargs) 
echo "Checking Env Vars: "

for ENV in $(cat .env | awk -F= '{print $1}')
do
    echo "Var Name: $ENV ----- Value: $(printenv $ENV)"
done