forceBuild='false'

CleanUpProject() {
    echo "-----Starting Clean Up-----"
    echo "Removing target directory..."
    rm -r target/
    echo "-----Finished Clean Up-----"
}

BuildProject() {
    echo $forceBuild
    if [ -d target/ ] && [ $forceBuild == 'false' ] ;
    then
        echo "Build already exists. Please clean up with the -c flag and rebuild."
        echo "Or force a build with the -f flag"
        exit 1
    fi

    echo "-----Starting Build-----"
    echo "Importing Env Vars from .env file."
    export $(grep -v '^#' .env | xargs) 
    echo "Checking Env Vars: "

    for ENV in $(cat .env | awk -F= '{print $1}')
    do
        echo "Var Name: $ENV ----- Value: $(printenv $ENV)"
    done

    echo "Beginning Build."
    mkdir "target"
    echo "Moving backend"
    cp -r backend target/
    echo "Moving Frontend."
    cp -r frontend target/
}

while getopts 'cfh' flag; do
  case "${flag}" in
    c) CleanUpProject 
        exit 0;;
    f) forceBuild='true';;
    h) 
      echo "CoinbaseBoard - attempt to build CoinbaseBoard"
      echo " "
      echo "To Run:"
      echo "bash build.sh [arguments]"
      echo " "
      echo "options:"
      echo "-h                        Show brief help"
      echo "-f                        Force a build while cleaning previous build"
      echo "-c,                       Clean Up A previous Build";;
    *) break ;;
  esac
done

BuildProject