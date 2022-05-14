# running instructions:
# from root of the project, execute me as
# pull latest reminder image, make sure ~/reminder directory exists, and run the tool
# `. ./scripts/run_via_docker.sh`
# run the tool (just run, without pulling image and other initialization steps)
# `. ./scripts/run_via_docker.sh fast`

MODE="${@:-default}"

if [[ ${MODE} != "fast" ]]; then
    # pull the image (or get the latest image)
    docker pull goyalmunish/reminder

    # make sure the directory for the data file exists on the host machine
    mkdir -p ~/reminder
fi

# spin up the container, with data file shared from the host machine
docker run -it --rm --name reminder -v ~/reminder:/root/reminder goyalmunish/reminder
