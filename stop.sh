docker_start(){
    docker-compose up --build
}

killer(){
    docker-compose down
}
main(){
    docker_start &
    killer
    wait 

}

main